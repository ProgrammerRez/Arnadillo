import json  # noqa: EXE002
import logging
import logging.config
import os

from dotenv import load_dotenv

# Loading Environment Variables
load_dotenv(".env")

# Getting Config Path from ENV VARS and trying to load the config
CONFIG_PATH: str = os.getenv("CONFIG_PATH", "config.json")
try:
    with open(CONFIG_PATH, "r", encoding="utf-8") as f:
        config = json.load(f)
        logging.config.dictConfig(config)
except json.JSONDecodeError:
    print("[ERROR]: Logging is disabled due to config issues")

# Setting Up the logger objects for importing
main_logger: logging.Logger = logging.getLogger("main")
rm_logger: logging.Logger = logging.getLogger("resource_management")
test_logger: logging.Logger = logging.getLogger("test_suite")


