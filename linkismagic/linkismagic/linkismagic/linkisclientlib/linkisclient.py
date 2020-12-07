"""
 * Copyright 2020 WeBank
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
"""

import json
import os
import getpass
import pickle
import time
from linkismagic.linkisclientlib.linkishttpclient import LinkisHttpClient
import datetime


class LinkisClient:

    def __init__(self):
        # 判断是否使用token还是直接登陆获取token，如果错误返回错误
        self.user = getpass.getuser()
        self.config = {}
        self.base_url = ""
        self.load_config()
        self.headers = {
            "Token-User": self.user,
            "Content-Type": "application/json",
            "Token-Code": self.config["Token-Code"]
        }

        # init directory config
        dir_name = str(datetime.datetime.now().timestamp())
        self.var_store_dir_linkis = self.config["var_store_dir_linkis"] + "/" + dir_name
        self.var_store_dir_jupyter = self.config["var_store_dir_jupyter"] + "/" + dir_name
        self.file_store_dir_linkis = self.config["file_store_dir_linkis"]
        os.mkdir(self.var_store_dir_jupyter)

        if "username" in self.config.keys() and "password" in self.config.keys():
            self.cookies = self.login(self.config["username"], self.config["password"])
        else:
            self.cookies = {}

    def login(self, username, password):
        response = self.linkishttpclient.login(username, password)
        if response.status_code != 200:
            raise Exception("HTTPException:" + response.text)
        cookies = response.cookies
        return cookies

    # load config and return configuration
    def load_config(self):
        filepath = os.path.join("/etc/linkismagic/linkismagic.json")
        if not os.path.exists(filepath):
            print("config not exisit")
        with open(filepath, 'r') as config_file:
            self.config = json.load(config_file)
        self.base_url = self.config["base_url"]
        if "Token-Code" not in self.config.keys():
            self.config["Token-Code"] = "PROPHECIS"
        if "linkis_dir" not in self.config.keys():
            self.config["linkis_dir"] = "/data/bdap-ss/mlss-data/"
        if "jupyter_dir" not in self.config.keys():
            self.config["jupyter_dir"] = "/workspace"
        if "host_dir" not in self.config.keys():
            self.config["host_dir"] = "/data/bdap-ss/mlss-data/"

        # jupyter_dir path should equals linkis_dir + user path
        if not ("var_store_dir_linkis" in self.config.keys()):
            self.config["var_store_dir_linkis"] = self.config["linkis_dir"] + self.user + "/.mlss_var"
        if not ("var_store_dir_jupyter" in self.config.keys()):
            self.config["var_store_dir_jupyter"] = self.config["jupyter_dir"] + "/.mlss_var"
        if not ("file_store_dir_host" in self.config.keys()):
            self.config["file_store_dir_host"] = self.config["host_dir"] + self.user
        self.linkishttpclient = LinkisHttpClient(self.base_url, self.headers, self.config)

    # 执行、判断结果，返回结果，task
    def execute(self, application, code, run_type=None):
        response = self.linkishttpclient.execute(application, code, run_type)
        if response.status_code != 200:
            raise Exception("HTTPException:" + response.text)
        else:
            result = json.loads(response.text)
            return True, result['data']['execID'], result['data']['taskID']

    # 封装为关键信息和关键报错
    def log(self, execID):
        response = self.linkishttpclient.log(execID)
        if response.status_code != 200:
            raise Exception("HTTPException:" + response.text)
        else:
            result = json.loads(response.text)
            return True, result['data']

    def status(self, execID):
        response = self.linkishttpclient.status(execID)
        if response.status_code != 200:
            raise Exception("HTTPException:" + response.text)
        else:
            result = json.loads(response.text)
            return True, result['data']['status']

    def kill(self, execID):
        response = self.linkishttpclient.kill(execID)
        if response.status_code != 200:
            raise Exception("HTTPException:" + response.text)
        else:
            result = json.loads(response.text)
            return True, result["message"]

    def progress(self, execID):
        response = self.linkishttpclient.progress(execID)
        if response.status_code != 200:
            raise Exception("HTTPException:" + response.text)
        else:
            progress = json.loads(response.text)
            return True, progress['data']['progress']

    def job_detail(self, taskID):
        return self.linkishttpclient.job_detail(taskID)

    '''
        return http_status_code, has_result_dir, Path
    '''
    def has_result_dir(self, task_id):
        response = self.job_detail(task_id)
        if response.status_code != 200:
            raise Exception("HTTPException:" + response.text)
        else:
            job_detail = json.loads(response.text)
            if job_detail['data']['task']['resultLocation'] is None:
                return True, False, None
            else:
                return True, True, job_detail['data']['task']['resultLocation']

    def get_result_path(self, result_dir):
        response = self.linkishttpclient.dir_file_trees(result_dir)
        if 200 != response.status_code:
            raise Exception("HTTPException:" + response.text)
        dir_file_trees = json.loads(response.text)
        # 文件可能未保存完整
        if dir_file_trees['data']['dirFileTrees'] is None:
            import time
            time.sleep(1)
            response = self.linkishttpclient.dir_file_trees(result_dir)
            dir_file_trees = json.loads(response.text)
        if dir_file_trees['data']['dirFileTrees'] is None:
            return True, None
        return True, dir_file_trees['data']['dirFileTrees']['children']

    def open_result(self, result_file_path):
        response = self.linkishttpclient.open_file(result_file_path)
        # 400通常是文件尚未保存完
        while response.status_code == 400:
            #            print("open request body: " + str(response.reason))
            response = self.linkishttpclient.open_file(result_file_path)
        if response.status_code == 200:
            open_file = json.loads(response.text)
            return True, open_file['data']["fileContent"]
        else:
            return False, None

    '''
        return status_code, download_status
    '''
    def download_csv(self, task_id, store_path=""):
        status, has_result_dir, file_path = self.has_result_dir(task_id)
        if not status:
            return False, "Failed"
        elif has_result_dir is None:
            return True, "None"
        status, file_list = self.get_result_path(file_path)
        if not status:
            return False, "Failed"
        elif has_result_dir is None:
            return True, "None"
        elif file_list is None:
            return True, "None"

        # 多结果集遍历保存
        for i in range(len(file_list)):
            print(file_list[i]['path'])
            r = self.linkishttpclient.download_csv(file_path=file_list[i]['path'])
            #        filename = file_path[file_path.rindex("/") + 1:]
            filename = store_path + "-" + str(i) + ".csv"
            if store_path == "":
                path = os.path.join(os.path.expanduser('~'), filename)
            else:
                path = store_path + filename
            print(path)
            with open(path, "wb") as f:
                for chunk in r.iter_content(chunk_size=1024):
                    if chunk:
                        f.write(chunk)

        return True, "Success"

    def upload(self, upload_file_path):
        pass

    def job_history(self):
        json_history = json.loads(self.linkishttpclient.job_history().text)
        return json_history['data']

    '''
    return status_code, exec_status, log, result, var
    '''
    def get_execute_result(self, exec_id, task_id):
        # Wait Job Complete
        status, exec_status = self.status(exec_id)
        while status and exec_status != "Succeed":
            if exec_status == "Failed":
                status, log = self.log(exec_id)
                job_detail = self.job_detail(task_id)
                error_log = {"log": log, "keyword": json.loads(job_detail.text)['data']['task']['errDesc']}
                return True, exec_status, error_log, None
            elif not status:
                # raise Exception
                return False, None, None, None
            time.sleep(5)
            status, exec_status = self.status(exec_id)
        status, has_result_idr, result_dir = self.has_result_dir(task_id)
        if not status:
            return False, None, None, None
        #        print("get result dir...")
        status, result_path_list = self.get_result_path(result_dir)
        if not status:
            return False, exec_status, None, None
        if result_path_list is None:
            return True, exec_status, None, None
        #        print("open result...")

        result_dict = {}
        for i in range(len(result_path_list)):
            status, result = self.open_result(result_path_list[i]['path'])
            if not status:
                return False, exec_status, None, None
            result_dict["result" + str(i)] = result
        return True, exec_status, None, result_dict

    def log_key_info(self, execID):
        response = self.linkishttpclient.log(execID)
        if response.status_code != 200:
            raise Exception("HTTPException:" + response.text)
        else:
            result = json.loads(response.text)
            return True, result['data']['log']

    def engines(self):
        response = self.linkishttpclient.engines()
        if response.status_code != 200:
            raise Exception("HTTPException:" + response.text)
        else:
            result = json.loads(response.text)
            if result['data'] == {}:
                return True, []
            return True, result['data']['engines']

    def engine(self, engine_manager_instance):
        status, engine_list = self.engines()
        if not status:
            raise Exception("Get Engine Fail")
        for i in range(len(engine_list)):
            if engine_list[i]["engineManagerInstance"] == engine_manager_instance:
                return engine_list[i]["ticketId"], engine_list[i]["moduleName"], engine_list[i]["creator"]
        return None, None, None

    def engine_kill(self, engine_manager_instance):
        ticket_id, application_name, creator = self.engine(engine_manager_instance)
        response = self.linkishttpclient.engine_kill(ticket_id, application_name, engine_manager_instance, creator)
        if response.status_code != 200:
            raise Exception("HTTPException:" + response.text)
        else:
            result = json.loads(response.text)
            return True, result['data']

    def download_by_pipeline_engine(self, task_id, store_path):
        status, has_result_dir, file_path = self.has_result_dir(task_id)
        if not status:
            return False, "Failed"
        elif has_result_dir is None:
            return True, "None"
        status, file_list = self.get_result_path(file_path)
        if not status:
            return False, "Failed"
        elif has_result_dir is None:
            return True, "None"
        elif file_list is None:
            return True, "None"

        # 多结果集遍历保存
        for i in range(len(file_list)):
            source = file_list[i]['path']
            target = self.file_store_dir_linkis
            store_code = "from " + source + " to " + target + store_path
            if store_code.endswith(".csv"):
                store_code = store_code + ".csv"
            print(store_code)
            status, exec_id, task_id = self.execute("pipeline", code=store_code, run_type="")
            status, exec_status, log, result = self.get_execute_result(exec_id, task_id)
            if status:
                if exec_status == "Failed":
                    print(log["keyword"])
                    print(log["log"])
                else:
                    continue
            else:
                raise Exception("HTTPException")
        return True, "Success"

    # Run in Linkis Site
    def define_pickel_code(self, var):
        file_path = self.var_store_dir_linkis + "/" + var + ".pkl"
        var_code = "\nimport pickle\n" + \
                   "with open(\"" + file_path + "\",\"wb\") as f:\n" + \
                   "\tpickle.dump(" + var + ", f )"
        return var_code

    def load_pickle_var(self, var):
        file_path = self.var_store_dir_jupyter + "/" + var + ".pkl"
        with open(file_path, "rb+") as f:
            var_instance = pickle.load(f)
        os.remove(file_path)
        return var_instance

    def pyspark_load_pickle_code(self, var_name):
        file_path = self.var_store_dir_linkis + "/" + var_name + ".pkl"
        var_code = "\nimport pickle\n" + \
                   "with open(\"" + file_path + "\",\"rb+\") as f:\n" + \
                   "\t" + var_name + " = pickle.load(f)\n"
        return var_code

    def save_pickle_file(self, var_name, value):
        file_path = self.var_store_dir_jupyter + "/" + var_name + ".pkl"
        with open(file_path, "wb") as f:
            pickle.dump(value, f)

    def delete_upload_var(self, var):
        os.remove(self.var_store_dir_jupyter + "/" + var + ".pkl")

    def refresh_cookies(self):
        self.cookies = self.login(self.config["username"], self.config["password"])
