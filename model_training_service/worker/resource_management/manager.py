'''This file will contain the main manager object. With the following features
1. Monitors System Resources
2. Calculates Rough Estimate of Resource Utilization for Jobs
3. Also Cleans up stuck Jobs
4. Says yes or no to incoming jobs
5. So more or less of a dynamic decision system

TODO: Create Custom Logger
'''

import asyncio
import time
import GPUtil
import psutil
from typing import Optional, Any, Dict, Tuple
from dotenv import load_dotenv

load_dotenv(".env")

# TODO: Put these in a .env file
SAFE_JOBS = 2
MAX_JOBS = 3
SAFE_GPU_LIMIT = 200.0


# The main System Resource Manager Object
class SystemResourcesManager:
    def __init__(self, max_ram: float, 
                max_cpu:float, 
                max_gpu: float, 
                selected_gpu: int, 
                job_timeout: int):
        
        # System Resource Management
        self.max_ram = max_ram
        self.max_cpu = max_cpu
        self.max_gpu = max_gpu

        # Managing GPU's
        self.available_gpus = GPUtil.getGPUs()
        self.stats = self.get_gpus()
        self.selected_gpu: Optional[int] = self.stats[selected_gpu] 

        # Job Timeout Limit
        self.job_timeout = job_timeout

        # Telemettery for `stuck` jobs
        self.active_jobs: Dict[str, Dict[str, Any]] = {}
        self.lock = asyncio.Lock()

        # Background Task Handling
        self._monitor_task = Optional[asyncio.Task] = None 

    def gpu_stats(self):
        """This function skims the available gpu list and provides details based on the gpu"""
        stats = {}
        for gpu in self.available_gpus:
            stats[gpu.id] = {
                "name": gpu.name,
                "utilization": gpu.utilization,
                "memory_usage": gpu.memoryUsed,
                "total_memory": gpu.memoryTotal,
                "temp": gpu.temperature
            }
        return stats 


    async def start (self):
        """Starts the background cleaning task"""
        self._monitor_task = asyncio.create_task(self._background_monitoring_loop)
        print("[ResourceManager] Started background monitoring and cleanup loop.")

    async def stop (self):
        """Stops the backgorund cleaning loop cleanly"""
        if self._monitor_task:
            self._monitor_task.cancel()
            try:
                await self._monitor_task
            except asyncio.CancelledError:
                print("[Resource Manager] Stopped background loop")

    # This is the background monitoring loop
    async def _background_monitoring_loop(self):
        """
        1. Monitors System Resources at Regular Intervals
        2. Cleans Stuck Jobs
        """
        while True:
            try:
                self.current_ram_percent: float = psutil.virtual_memory().percent
                self.current_cpu_percent: float = psutil.cpu_percent(interval=None)
                self.current_gpu_percent: float = self.selected_gpu["utilization"]

                async with self.lock:
                    current_time = time.time()
                    stuck_jobs = []
                    for job_id, meta in self.active_jobs.items():
                        elapsed_time = current_time - meta["start_time"]
                        if elapsed_time > self.job_timeout:
                            stuck_jobs.append(job_id)
                
                    for job_id in stuck_jobs:
                        print(f"[Resource Manager] Warning: Job ID {job_id} exceeded")
                        del self.active_jobs[job_id]
            except Exception as e:
                print(f"[Resource Manager] Error in Monitor loop: {e}")

            await asyncio.sleep(3.0)        

    # This function calculates a job       
    def estimate_job_resources(self, training_params: Dict[str, Dict[str, Any]]):

        model_type = training_params.get("mdoel_type", "traditional")
        model_params = training_params.get("model_params", {})

        estimated_ram_mb = 20.0
        estimated_vram_mb = 20.0
        estimated_cpu_load = 20.0

        if model_type == "traditional":
            # Using traditional scaling factors like dataset limit, max_iter and estimators
            limit = model_params.get("limit", 1000)
            max_iter = model_params.get("max_iter", 100)
            n_estimators = model_params.get("n_estimators", 10)

            estimated_ram_mb = 200.0 + (limit + 0.005) + (n_estimators * 5.0)
            estimated_cpu_load = min(85.0, 20.0 + (max_iter * 0.2) + (n_estimators * 1.5))

        elif model_type == "deep":
            # Using Deep Learning Factors like architecture layers, neuorns and batch sizes
            layers = model_params.get("layers", []) 
            num_layers = len(layers) if isinstance(layers, list) else 0

            total_neurons += [layer[1] for layer in layers]
            batch_size = model_params.get("batch_size", 32)

            estimated_ram_mb = 800.0 + (total_neurons * 2.5)
            estimated_vram_mb = 500.0 + (total_neurons * batch_size * 0.2)
            estimated_cpu_load = 50.0

        
        return {
            "estimated_ram_mb": estimated_ram_mb,
            "estimated_vram_mb": estimated_vram_mb,
            "estimated_cpu_load": estimated_cpu_load
        }


    # This function accepts a Job
    async def can_accpet_job(self, training_params: Dict[str, Any]) -> Tuple[bool, str]:

        # Check Avaialable Resources
        if self.current_ram_percent >= self.max_ram:
            return False, f"System RAM is Saturated: {self.current_ram_percent}"
        
        if self.current_cpu_percent >= self.max_cpu:
            return False, f"CPU utilization is too High: {self.current_cpu_percent}"
        
        if self.current_gpu_percent >= self.max_gpu:
            return False, f"GPU utilization is too High: {self.current_gpu_percent}"

        # Estimate for Resources Required

        estimate = self.estimate_job_resources(training_params)

        async with self.lock:
            activte_count = len(self.active_jobs)

            # First see if they are not above safe jobs threshold
            if activte_count >= SAFE_JOBS and estimate["estimated_vram_mb"] >= SAFE_GPU_LIMIT:
                return False, "Rejected: Heavy deep learning workload already active."
            
            if activte_count >= MAX_JOBS:
                return False, "Rejected: Maximum Concurrent Job Slots Reached"
            
        
        return True, "Approved"
    
    # Registers a job into the active tracking
    async def register_job_start(self, job_id: str, training_params: Dict[str, Any]):
        async with self.lock:
            self.active_jobs[str(job_id)] = {
                "start_time": time.time(),
                "estimate": self.estimate_job_resources(training_params)
            }

    
    # Registers Job Completion and Removes the job from active tracking
    async def register_job_completion(self, job_id :str):
        async with self.lock:
            if str(job_id) in self.active_jobs:
                del self.active_jobs[str(job_id)]
