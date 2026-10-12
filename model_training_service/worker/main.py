import resource_management as rm

if __name__ == "__main__":
    obj = rm.SystemResourcesManager(
        max_ram = 90.0,
        max_cpu = 90.0,
        max_gpu = 90.0,
        selected_gpu = 0,
        job_timeout = 30
    )