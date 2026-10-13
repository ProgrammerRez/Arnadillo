from typing import Any
from resource_management import SystemResourcesManager
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

class TestSystemResourcesManager(unittest.TestCase):

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
            "StartTime": ,
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
            estimated = self.resource_manager.estimate_job_resources(self.job1)
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


# Run the tests using this command: `python -m unittest <file_path>`
