import os

HOME_PATH = os.environ.get("LINKISMAGIC_CONF_DIR", "~/.linkismagic")
CONFIG_FILE = os.environ.get("LINKISMAGIC_CONF_FILE", "config.json")

LINKIS_REST_METHOD_EXECUTE = "/api/rest_j/v1/entrance/execute"
LINKIS_REST_METHOD_STATUS = "/api/rest_j/v1/entrance/${execID}/status"
LINKIS_REST_METHOD_LOG = "/api/rest_j/v1/entrance/${execID}/log"
LINKIS_REST_METHOD_PROGRESS = "/api/rest_j/v1/entrance/${execID}/progress"
LINKIS_REST_METHOD_KILL = "/api/rest_j/v1/entrance/{execID}/kill"

MAGICS_LOGGER_NAME = "linkisLogger"

LINKIS_JOB_ID = "execID"
LINKIS_JOB_TYPE = "runType"
LINKIS_ENGINE_TYPE = "executeApplicationName"
