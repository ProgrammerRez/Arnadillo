from resource_management import SystemResourcesManager
from logger import main_logger as log

if __name__ == "__main__":
    obj = SystemResourcesManager(
        max_ram = 90.0,
        max_cpu = 90.0,
        max_gpu = 90.0,
        selected_gpu = 0,
        job_timeout = 30
    )

    log.info("Started Logging")

