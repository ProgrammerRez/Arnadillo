from dotenv import load_dotenv
import logging.config
import logging
import json
import os

# Loading Environment Variables
load_dotenv(".env")

# Getting Config Path from ENV VARS and trying to load the config
CONFIG_PATH: str = os.getenv("CONFIG_PATH", "config.json")
try:
    with open(CONFIG_PATH, "r", encoding="utf-8") as f:
        config = json.load(f)
        logging.config.dictConfig(config)
except json.JSONDecodeError as e:
    print("[ERROR]: Logging is disabled due to config issues")

# Setting Up the logger objects for importing
main_logger: logging.Logger = logging.getLogger("main")
rm_logger: logging.Logger = logging.getLogger("resource_management")


