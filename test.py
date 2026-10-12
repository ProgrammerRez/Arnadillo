import GPUtil as gpu

gpus = gpu.getGPUs()

for gpu in gpus:
  print(f"--- GPU {gpu.id}: {gpu.name} ---")
  print(f"Utilization : {gpu.load * 100:.1f}%")
  print(
      f"Memory Usage: {gpu.memoryUsed} MB / {gpu.memoryTotal} MB"
      f" ({gpu.memoryUtil * 100:.1f}%)"
  )
  print(f"Temperature : {gpu.temperature} °C")