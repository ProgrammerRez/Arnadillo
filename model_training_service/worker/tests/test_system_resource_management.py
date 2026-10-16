import asyncio
from typing import Any
from resource_management import SystemResourcesManager
from datetime import datetime, timedelta
from logger import test_logger as log
from dotenv import load_dotenv
from models import Job
import unittest
import os

load_dotenv(".env")
MAX_CPU = float(os.getenv("MAX_CPU", 90.0))
MAX_GPU = float(os.getenv("MAX_GPU", 90.0))
MAX_RAM = float(os.getenv("MAX_RAM", 90.0))
SELECTED_GPU = int(os.getenv("SELECTED_GPU", 0))
TIMEOUT = int(os.getenv("TIMEOUT", 30))

# This class will test the SystemResourceManager

class TestSystemResourcesManager(unittest.IsolatedAsyncioTestCase):

    @classmethod
    def setUpClass(cls) -> None:
        log.info("Setting up the test cases and the main object instance")

        cls.resource_manager = SystemResourcesManager(
            max_cpu=MAX_CPU,
            max_gpu=MAX_GPU,
            max_ram=MAX_RAM,
            selected_gpu=SELECTED_GPU,
            job_timeout=TIMEOUT
        )
        
        cls.job1 = Job(**{
            "ID": 1,
            "StartTime": None,
            "EstimatedResources": None,
            "DataParameters": None,
            "TrainingParameters":{
                "model_type": "traditional",
                "model_name": "LogisticRegression",
                "model_params": {
                    "max_iter": 1000
                },
            "EvalParameters": None,
            "DoOverParameters": None
            }
        })

        cls.job1_estimate: dict[str, float] = {
                                            'estimated_ram_mb': 1250.005, 
                                            'estimated_vram_mb': 20.0, 
                                            'estimated_cpu_load': 85.0
                                            }

        cls.stuckjob = Job(**{
            "ID": 1,
            "StartTime": None,
            "EstimatedResources": None,
            "DataParameters": None,
            "TrainingParameters":{
                "model_type": "traditional",
                "model_name": "LogisticRegression",
                "model_params": {
                    "max_iter": 1000
                },
            "EvalParameters": None,
            "DoOverParameters": None
            }
        })

    # This one tests the gpu function and it's output
    def test_gpu_stats(self) -> None:
        try:
            stats: dict[int, dict[str, Any]] = self.resource_manager.gpu_stats()
            self.assertIsInstance(stats, dict)
            for gpu_id, info in stats.items():
                self.assertIsInstance(gpu_id, int)
                self.assertIsInstance(info, dict)
            
        except Exception as e:
            log.error(f"GPU: Stats Error: {e}")
            raise e
        finally:
            log.info(f"GPU stats function has been tested: {stats}")

    # This one will test the estimated resources part
    def test_estimated_resources(self) -> None:
        try:
            estimated: dict[str, float | Any] | None = self.resource_manager.estimate_job_resources(self.job1)
            self.assertIsInstance(estimated, dict)
            for resource, estimate in estimated.items():
                self.assertIsInstance(resource, str)
                self.assertIsInstance(estimate, float)

            self.assertEqual(estimated, self.job1_estimate)
        except Exception as e:
            log.error(f"Resource Estimation Function Error: {e}")
            raise e
        finally:
            log.info(f"Estimation stats function has been tested: {estimated}")

    # This function will test the background loop
    async def test_background_loop(self) -> None:
        try:
            # This tests the start of the background process
            await self.resource_manager.start()
            self.assertIsInstance(self.resource_manager._monitor_task, asyncio.Task)
            self.assertFalse(self.resource_manager._monitor_task.done())

            # Testing the deletion of stuck jobs
            approved, reason = await self.resource_manager.can_accpet_job(self.stuckjob)
            if not approved:
                log.error(f"Background Loop Testing: Stuck Job Not Approved: {reason}")
                raise Exception(f"Background Loop Testing: Stuck Job Not Approved: {reason}")

            # if approved then register job
            await self.resource_manager.register_job_start(job=self.stuckjob, job_id=self.stuckjob.ID)
            self.assertIsInstance(self.resource_manager.active_jobs, dict)
            for id ,job in self.resource_manager.active_jobs.items():
                self.assertIsInstance(id, int)
                self.assertIsInstance(job, dict)
                for key, value in job.items():
                    self.assertIsInstance(key, str)
            
            # Checking if Job Exists
            self.assertTrue(len(self.resource_manager.active_jobs.items())==1)
            # Changing Time on the Stuck Job
            self.resource_manager.active_jobs[self.stuckjob.ID]["StartTime"] = (
                datetime.now() - timedelta(seconds=self.resource_manager.job_timeout + 5)
            )
            # Now checking whether job will be removed
            await asyncio.sleep(4.0)
            self.assertEqual(len(self.resource_manager.active_jobs), 0)

            # This tests the stop of the background process
            await self.resource_manager.stop()
            self.assertTrue(self.resource_manager._monitor_task.done() or self.resource_manager._monitor_task.cancelled())
        except Exception as e:
            log.error(f"Background Loop Testing Error: {e}")
            raise e
        finally:
            log.info("Background Loop Testing Successful")

    # Finally the Last Test would include the checking of job starting and completion
    async def test_job_start_and_completion(self):
        try:
            # First Register the Job
            self.assertIsInstance(self.job1, Job)
            await self.resource_manager.register_job_start(job=self.job1, job_id=self.job1.ID)
            self.assertEqual(len(self.resource_manager.active_jobs), 1)

            # Then Complete the Job
            await self.resource_manager.register_job_completion(self.job1.ID)
            self.assertEqual(len(self.resource_manager.active_jobs), 0)
        except Exception as e:
            log.error(f"Job Start and Completeion Testing Error: {e}")
            raise e
        finally:
            log.info("Job Start and Completion Testing Successfully Completed")

# Run the tests using this command: `python -m unittest <file_path>`
