from hdijupyterutils.log import Log
import linkismagic.utils.configuration as conf
from linkismagic.utils.constants import MAGICS_LOGGER_NAME


class LinkisLog(Log):
    def __init__(self, class_name):
        super(LinkisLog, self).__init__(MAGICS_LOGGER_NAME, conf.logging_config(), class_name)