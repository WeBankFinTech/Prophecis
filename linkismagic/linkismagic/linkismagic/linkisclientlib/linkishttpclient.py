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
import requests

class LinkisHttpClient:
    def __init__(self, base_url, headers, configuration, login=False, username="", password=""):
        self.base_url = base_url
        self.headers = headers
        self.session = requests.Session()
        self.configuration = configuration
        if login:
            self.cookies = self.login(username, password)
        else:
            self.cookies = {}

    def login(self, username, password):
        login_url = self.base_url + "/api/rest_j/v1/user/login"
        login_data = {
            "userName": username,
            "password": password
        }
        login_result = self.session.post(login_url, data=json.dumps(login_data))
        return login_result

    def def_execute_entity(self, code, application, run_type=None, configuration=None):
        post_entity_map = {
            "method": "/api/rest_j/v1/entrance/execute",
            "executeApplicationName": application,
            "params": {
                "configuration": {
                        "special": {},
                        "runtime": {},
                        "startup": {},
                },
            },
            "executionCode": code,
        }
        if run_type is not None:
            post_entity_map['runType'] = run_type
        configuration_keys = list(configuration.keys())
        if "session_configs" in configuration_keys:
            spark_conf = configuration["session_configs"]
            spark_conf_keys = list(spark_conf.keys())
            startup = post_entity_map['params']['configuration']["startup"]
            for i in range(len(spark_conf_keys)):
                startup[spark_conf_keys[i]] = spark_conf[spark_conf_keys[i]]
        return post_entity_map

    def execute(self, application, code, run_type=None):
        execute_url = self.base_url + "/api/rest_j/v1/entrance/execute"
        post_entity = self.def_execute_entity(code, application, run_type, self.configuration)
        post_entity_json = json.dumps(post_entity)
        return self.session.post(execute_url, headers=self.headers, cookies=self.cookies, data=post_entity_json)

    def log(self, execID):
        log_url = self.base_url + "/api/rest_j/v1/entrance/" + execID + "/log?fromLine=0&size=10000"
        return self.session.get(log_url, cookies=self.cookies, headers=self.headers)

    def status(self, execID):
        status_url = self.base_url + "/api/rest_j/v1/entrance/" + execID + "/status"
        return self.session.get(status_url, cookies=self.cookies, headers=self.headers)

    def kill(self, execID):
        kill_url = self.base_url + "/api/rest_j/v1/entrance/" + execID + "/kill"
        return self.session.get(kill_url, cookies=self.cookies, headers=self.headers)

    def progress(self, execID):
        progress_url = self.base_url + "/api/rest_j/v1/entrance/" + execID + "/progress"
        return self.session.get(progress_url, cookies=self.cookies, headers=self.headers)

    def engines(self):
        progress_url = self.base_url + "/api/rest_j/v1/resourcemanager/engines"
        data = {}
        return self.session.post(progress_url, cookies=self.cookies, data=json.dumps(data), headers=self.headers)
    
    def engine_kill(self, ticket_id, module_name, engine_manager_instance, creator):
        engine_kill_url = self.base_url + "/api/rest_j/v1/resourcemanager/enginekill"
        data = [{
            "ticketId": ticket_id,
            "moduleName": module_name,
            "engineManagerInstance": engine_manager_instance,
            "creator": creator
        }]
        return self.session.post(engine_kill_url, cookies=self.cookies, data=json.dumps(data), headers=self.headers)

    def job_history(self):
        job_history_url = self.base_url + "/api/rest_j/v1/jobhistory/list?pageNow=1&pageSize=20"
        return self.session.get(job_history_url, cookies=self.cookies, headers=self.headers)

    # 获取单个Task信息和job_detail一致
    def job_detail(self, task_id):
        job_detail_url = self.base_url + "/api/rest_j/v1/jobhistory/"+str(task_id)+"/get"
        return self.session.get(job_detail_url, cookies=self.cookies, headers=self.headers)

    def dir_file_trees(self, result_location):
        dir_file_trees_url = self.base_url + "/api/rest_j/v1/filesystem/getDirFileTrees?path=" + result_location
        return self.session.get(dir_file_trees_url, cookies=self.cookies, headers=self.headers)

    def open_file(self, result_file_path):
        open_file_url = self.base_url + "/api/rest_j/v1/filesystem/openFile?path=" + result_file_path
        return self.session.get(open_file_url, cookies=self.cookies, headers=self.headers)

    def download_csv(self, file_path=""):
        excel_file_url = self.base_url + "/api/rest_j/v1/filesystem/resultsetToExcel?path=" + file_path \
                         + "&outputFileType=csv"
        return self.session.get(excel_file_url, headers=self.headers, cookies=self.cookies)

    def download(self, download_file_path, file_path=""):
        download_file_url = self.base_url + "/api/rest_j/v1/filesystem/download"
        post_data = {
            "path": download_file_path
        }
        return self.session.post(download_file_url, headers=self.headers, data=json.dumps(post_data))

