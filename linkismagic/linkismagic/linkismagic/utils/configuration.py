from linkismagic.utils.constants import *
from hdijupyterutils.utils import join_paths
from hdijupyterutils.configuration import override as _override
from hdijupyterutils.configuration import override_all as _override_all
from hdijupyterutils.configuration import with_override

d = {}
path = join_paths(HOME_PATH, CONFIG_FILE)

_with_override = with_override(d, path)

# @_with_override
# def custom_headers():
#     return {
#         "TokenUser" = ,
#         "TokenCode" = ,
#         "Content-Type": "application/json"
#     }