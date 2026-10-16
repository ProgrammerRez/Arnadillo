'''This file will contain the main manager object. With the following features
1. Monitors System Resources
2. Calculates Rough Estimate of Resource Utilization for Jobs
3. Also Cleans up stuck Jobs
4. Says yes or no to incoming jobs
5. So more or less of a dynamic decision system

'''  # noqa: EXE002

from datetime import datetime
import asyncio
import os
import time
from typing import Any

import GPUtil
import psutil
from dotenv import load_dotenv
from logger import rm_logger as log

from models import Job

load_dotenv(".env")

# Loading the ENV VARS
try:
    SAFE_JOBS = int(os.getenv("SAFE_JOBS", 1))  # noqa: PLW1508
    MAX_JOBS = int(os.getenv("MAX_JOBS", 3))  # noqa: PLW1508
    SAFE_GPU_LIMIT = float(os.getenv("SAFE_GPU_LIMIT", 400.0))  # noqa: PLW1508
except Exception as e:  # noqa: BLE001
    log.error(f"Unable to load Environment Varaiables: {e}")

log.info("Loaded Environment Variables")

# The main System Resource Manager Object
class SystemResourcesManager:
    def __init__(self, max_ram: float, 
                max_cpu:float, 
                max_gpu: float, 
                selected_gpu: int, 
                job_timeout: int) -> None:
        
        # System Resource Management
        self.max_ram = max_ram
        self.max_cpu = max_cpu
        self.max_gpu = max_gpu

        # Managing GPU's
        self.available_gpus = GPUtil.getGPUs()
        self.stats: dict[int, dict[str, Any]] = self.gpu_stats()
        self.selected_gpu: int | None = self.stats[selected_gpu] 

        # Gathering Usage Data
        self.current_ram_percent: float = psutil.virtual_memory().percent
        self.current_cpu_percent: float = psutil.cpu_percent(interval=None)
        self.current_gpu_percent: float = self.selected_gpu["utilization"] if len(self.available_gpus) > 0 else 0.0

        # Job Timeout Limit
        self.job_timeout = job_timeout

        # Telemettery for `stuck` jobs
        self.active_jobs: dict[int, dict[str, Any]] = {}
        self.lock = asyncio.Lock()

        # Background Task Handling
        self._monitor_task: asyncio.Task | None = None 

        log.info("Initialized the System Resource Manager")

    def gpu_stats(self) -> dict[int, dict[str, Any]]:
        """This function skims the available gpu list and provides details based on the gpu"""
        stats: dict[int, dict[str, Any]] = {}
        for gpu in self.available_gpus:
            stats[gpu.id] = {
                "name": gpu.name,
                "utilization": gpu.load * 100,
                "memory_usage": gpu.memoryUsed,
                "total_memory": gpu.memoryTotal,
                "temp": gpu.temperature
            }
        log.info("Retrieved GPU Info")
        return stats 


    async def start (self) -> None:
        """Starts the background cleaning task"""
        self._monitor_task = asyncio.create_task(self._background_monitoring_loop())
        log.info("[ResourceManager] Started background monitoring and cleanup loop.")

    async def stop (self) -> None:
        """Stops the backgorund cleaning loop cleanly"""
        if self._monitor_task:
            self._monitor_task.cancel()
            try:
                await self._monitor_task
            except asyncio.CancelledError:
                log.info("[Resource Manager] Stopped background loop")
            except Exception as e:  # noqa: BLE001
                log.error(f"Unexpected Error Occurred when closing background loop: {e}")

    # This is the background monitoring loop
    async def _background_monitoring_loop(self):
        """
        1. Monitors System Resources at Regular Intervals
        2. Cleans Stuck Jobs
        """
        self.current_ram_percent: float = psutil.virtual_memory().percent
        self.current_cpu_percent: float = psutil.cpu_percent(interval=None)
        self.current_gpu_percent: float = self.selected_gpu["utilization"] if len(self.available_gpus) > 0 else 0.0
        while True:
            try:
                async with self.lock:
                    current_time: datetime = datetime.now()
                    stuck_jobs: list[int] = []
                    for id, job in self.active_jobs.items():
                        start_time: datetime = job.get("StartTime")
                        if start_time:
                            elapsed_time: float = (current_time - start_time).total_seconds()
                            if elapsed_time > self.job_timeout:
                                stuck_jobs.append(id)
                        else:
                            log.error(f"No Start Time Added for job id: {id}")
                    for id in stuck_jobs:
                        log.info(f"[Resource Manager] Warning: Job ID {id} exceeded")
                        del self.active_jobs[id]
            except Exception as e:  # noqa: BLE001
                log.error(f"[Resource Manager] Error in Monitor loop: {e}")

            await asyncio.sleep(3.0)        

    # This function calculates a job       
    def estimate_job_resources(self, job: Job) -> dict[str, float | Any] | None:

        # Just doing a check and hoping it doesn't mess up any test or create bugs
        training_params: dict[str, Any] | None = job.TrainingParameters

        if training_params == None:
            log.error("Empty Training Parameters Provided")
            return None

        model_type: dict[str, Any] | str = training_params.get("model_type", "traditional")
        model_params: dict[str, Any] = training_params.get("model_params", {})

        if model_params == {}:
            log.warning("No Training Parameters Provided for a Finetuned Pipeline")

        estimated_ram_mb = 20.0
        estimated_vram_mb = 20.0
        estimated_cpu_load = 20.0

        if model_type == "traditional":
            # Using traditional scaling factors like dataset limit, max_iter and estimators
            limit = model_params.get("limit", 1000)
            max_iter = model_params.get("max_iter", 100)
            n_estimators = model_params.get("n_estimators", 10)

            estimated_ram_mb = 200.0 + (limit + 0.005) + (n_estimators * 5.0)
            estimated_cpu_load: float = min(85.0, 20.0 + (max_iter * 0.2) + (n_estimators * 1.5))

        elif model_type == "deep":
            # Using Deep Learning Factors like architecture layers, neuorns and batch sizes
            layers = model_params.get("layers", []) 
            total_neurons = 0
            for layer in layers:
                total_neurons += layer[1]
            batch_size = model_params.get("batch_size", 32)

            estimated_ram_mb: float | Any = 800.0 + (total_neurons * 2.5)
            estimated_vram_mb = 500.0 + (total_neurons * batch_size * 0.2)
            estimated_cpu_load = 50.0

        return {
            "estimated_ram_mb": estimated_ram_mb,
            "estimated_vram_mb": estimated_vram_mb,
            "estimated_cpu_load": estimated_cpu_load
        }


    # This function accepts a Job
    async def can_accpet_job(self, job: Job) -> tuple[bool, str]:

        log.info("Checking for available system resources")

        # Just doing a check and hoping it doesn't mess up any test or create bugs
        training_params: dict[str, Any] = job.TrainingParameters

        if len(training_params) == 0:
            log.error("Empty Training Parameters Provided")
            return False, "Empty Training Parameters Provided"

        # Check Avaialable Resources
        if self.current_ram_percent >= self.max_ram:
            log.warning(f"System RAM is Saturated: {self.current_ram_percent}")
            return False, f"System RAM is Saturated: {self.current_ram_percent}"
        
        if self.current_cpu_percent >= self.max_cpu:
            log.warning(f"CPU utilization is too High: {self.current_cpu_percent}")
            return False, f"CPU utilization is too High: {self.current_cpu_percent}"
        
        if self.current_gpu_percent >= self.max_gpu:
            log.error(f"GPU utilization is too High: {self.current_gpu_percent}")
            return False, f"GPU utilization is too High: {self.current_gpu_percent}"

    
        # Estimate for Resources Required
        estimate: dict[str, float | Any] | None = self.estimate_job_resources(job)

        async with self.lock:
            activte_count: int = len(self.active_jobs)

            # First see if they are not above safe jobs threshold
            if activte_count >= SAFE_JOBS and estimate["estimated_vram_mb"] >= SAFE_GPU_LIMIT:
                log.warning("Rejected: Heavy deep learning workload already active.")
                return False, "Rejected: Heavy deep learning workload already active."
            
            if activte_count >= MAX_JOBS:
                log.warning("Rejected: Maximum Concurrent Job Slots Reached")
                return False, "Rejected: Maximum Concurrent Job Slots Reached"
            
        log.info("Job Approved")
        return True, "Approved"
    
    # Registers a job into the active tracking
    async def register_job_start(self, job_id: int, job: Job) -> None:
        log.info("Registering for Job Start")
        try:
            async with self.lock:

                estimate: dict[str, float | Any] | None = self.estimate_job_resources(job)
                
                self.active_jobs[job_id] = {
                    "ID": job_id,
                    "StartTime": datetime.now(),
                    "EstimatedResources": estimate,
                    "DataParameters": job.DataParameters,
                    "Training_Parameters": job.TrainingParameters,
                    "EvalParameters": job.EvalParameters,
                    "DoOverParameters": job.DoOverParameters
                }
        except Exception as e:  # noqa: BLE001
            log.error(f"Unexpected error occurred when starting job: {e}")
        
        finally:
            log.info("Job Registered")

    
    # Registers Job Completion and Removes the job from active tracking
    async def register_job_completion(self, job_id :int) -> None:
        log.info("Completing / Cancelling Job")
        try:
            async with self.lock:
                if job_id in self.active_jobs.keys():
                    for id, _ in self.active_jobs.items(): 
                        if id == job_id:
                            del self.active_jobs[id]
        except Exception as e:  # noqa: BLE001
            log.error(f"Unexpected Error Occurred when stopping job: {e}")
        finally:
            log.info("Job Stopped Successfully")