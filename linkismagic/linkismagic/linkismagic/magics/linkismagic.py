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
from IPython.core.magic import Magics, magics_class
from hdijupyterutils.ipythondisplay import IpythonDisplay
from IPython.core.magic import line_cell_magic, needs_local_scope, line_magic
from IPython.core.magic_arguments import parse_argstring
from IPython.core.magic_arguments import argument, magic_arguments
from linkismagic.linkisclientlib.linkisclient import LinkisClient


@magics_class
class LinkisMagic(Magics):
    def __init__(self, shell, data=None, widget=None):
        super(LinkisMagic, self).__init__(shell)
        self.ipython_display = IpythonDisplay()
        self.data = data
        #        if widget is None:
        #            widget = MagicsControllerWidget(self.spark_controller, IpyWidgetFactory(), self.ipython_display)
        #        self.manage_widget = widget
        self.linkis_client = LinkisClient()

    @magic_arguments()
    @line_cell_magic
    @needs_local_scope
    @argument("-o", "--output", type=str, default=None, help="Output of Job ")
    @argument("-p", "--path", type=str, default=None, help="Download output in path ")
    @argument("-q", "--quiet", type=str, default=True, help="Do not display result on console")
    def spark(self, line, cell="", local_ns=None):
        user_input = parse_argstring(self.spark, line)

        code = cell
        status, exec_id, task_id = self.linkis_client.execute("spark", code, run_type="spark")
        if not status:
            raise Exception("HTTPException")
        else:
            status, exec_status, log, result = self.linkis_client.get_execute_result(exec_id, task_id)
            if status:
                if exec_status == "Failed":
                    print(log["keyword"])
                    print(log["log"])
                else:
                    if user_input.output is not None:
                        self.shell.user_ns[user_input.output] = result
                    if user_input.path is not None:
                        status, result = self.linkis_client.download_by_pipeline_engine(task_id, user_input.path)
                        if status and result == "None":
                            raise Exception("Save Error, Result dir is None")
                        elif not status or result != "Success":
                            raise Exception("Save Error, Result: " + result)
                    if user_input.quiet == "False" or user_input.quiet == "false":
                        return result
            else:
                raise Exception("HTTPException: get_execute_result error")

    @magic_arguments()
    @line_cell_magic
    @needs_local_scope
    @argument("-o", "--output", type=str, default=None, help="Output of Job ")
    @argument("-p", "--path", type=str, default=None, help="Download output in path ")
    @argument("-q", "--quiet", type=str, default=True, help="Do not display result on console")
    @argument("-v", "--var", type=str, default=None, help="transport var from spark cluster to local python ")
    @argument("-u", "--upload", type=str, default=None, help="transport var from local python to spark cluster")
    def pyspark(self, line, cell="", local_ns=None):
        user_input = parse_argstring(self.pyspark, line)
        pyspark_code = cell

        if user_input.upload is not None:
            pyspark_code = self.linkis_client.pyspark_load_pickle_code(user_input.upload) + pyspark_code
            self.linkis_client.save_pickle_file(user_input.upload, self.shell.user_ns[user_input.upload])

        if user_input.var is not None:
            pyspark_code = pyspark_code + self.linkis_client.define_pickel_code(user_input.var)

        status, exec_id, task_id = self.linkis_client.execute("spark", pyspark_code, run_type="python")
        if not status:
            raise Exception("HTTPException: execute job error")
        else:
            status, exec_status, log, result = self.linkis_client.get_execute_result(exec_id, task_id)
            if status:
                if exec_status == "Failed":
                    print(log["keyword"])
                    print(log["log"])
                else:
                    if user_input.output is not None:
                        self.shell.user_ns[user_input.output] = result
                    if user_input.path is not None:
                        status, result = self.linkis_client.download_by_pipeline_engine(task_id, user_input.path)
                        if status and result == "None":
                            raise Exception("Save Error, Result dir is None")
                        elif not status or result != "Success":
                            raise Exception("Save Error, Result: " + result)
                    if user_input.var is not None:
                        self.shell.user_ns[user_input.var] = self.linkis_client.load_pickle_var(user_input.var)
                    if user_input.upload is not None:
                        self.linkis_client.delete_upload_var(user_input.upload)
                    # self.ipython_display.display(result)
                    if user_input.quiet == "False" or user_input.quiet == "false":
                        return result
            else:
                raise Exception("HTTPException: get_execute_result error")

    @magic_arguments()
    @line_cell_magic
    @needs_local_scope
    @argument("-o", "--output", type=str, default=None, help="Output of Job ")
    @argument("-p", "--path", type=str, default=None, help="Download output in path ")
    @argument("-q", "--quiet", type=str, default=True, help="Do not display result on console")
    def sql(self, line, cell="", local_ns=None):
        user_input = parse_argstring(self.sql, line)
        sql_code = 'if "hiveContext" not in locals().keys():\n \
        \tfrom pyspark.sql import HiveContext\n\
        \thiveContext = HiveContext(sc)\n'
        cell_list = cell.split("\n")
        for i in range(len(cell_list)):
            if "" != cell_list[i]:
                sql_code = sql_code + 'hiveContext.sql("' + cell_list[i] + '").show()' + '\n'
        print(sql_code)
        status, exec_id, task_id = self.linkis_client.execute("spark", sql_code, run_type="python")
        if not status:
            raise Exception("HTTPException")
        else:
            status, exec_status, log, result = self.linkis_client.get_execute_result(exec_id, task_id)
            if status:
                if exec_status == "Failed":
                    print(log["keyword"])
                    print(log["log"])
                else:
                    if user_input.output is not None:
                        self.shell.user_ns[user_input.output] = result
                    if user_input.path is not None:
                        status, result = self.linkis_client.download_by_pipeline_engine(task_id)
                        if status and result == "None":
                            raise Exception("Save Error, Result dir is None")
                        elif not status or result != "Success":
                            raise Exception("Save Error, Result: " + result)
                    # self.ipython_display.display(result)
                    if user_input.quiet == "False" or user_input.quiet == "false":
                        return result
            else:
                raise Exception("HTTPException")

    @magic_arguments()
    @line_cell_magic
    @needs_local_scope
    @argument("-o", "--output", type=str, default=None, help="output var of Job ")
    @argument("-p", "--path", type=str, default=None, help="Download output in path ")
    @argument("-q", "--quiet", type=str, default=True, help="Do not display result on console")
    def sparksql(self, line, cell="", local_ns=None):
        user_input = parse_argstring(self.sparksql, line)
        code = cell
        status, exec_id, task_id = self.linkis_client.execute("spark", code, run_type="sql")
        if not status:
            raise Exception("HTTPException")
        else:
            status, exec_status, log, result = self.linkis_client.get_execute_result(exec_id, task_id)
            if status:
                if exec_status == "Failed":
                    print(log["keyword"])
                    print(log["log"])
                else:
                    if user_input.output is not None:
                        self.shell.user_ns[user_input.output] = result
                    if user_input.path is not None:
                        # status, result = self.linkis_client.download_csv(task_id, user_input.path)
                        status, result = self.linkis_client.download_by_pipeline_engine(task_id, user_input.path)
                        if status != True or result != "Success":
                            raise Exception("Save Error, result: " + result)
                    # self.ipython_display.display(result)
                    if user_input.quiet == "False" or user_input.quiet == "false":
                        return result
            else:
                raise Exception("HTTPException")

    @magic_arguments()
    @line_cell_magic
    @needs_local_scope
    @argument("-o", "--output", type=str, default=None, help="Output of Job ")
    def listjob(self, line, cell="", local_ns=None):
        user_input = parse_argstring(self.listjob, line)
        job_list = self.linkis_client.job_history()
        if user_input.output is not None:
            self.shell.user_ns[user_input.output] = job_list
        else:
            return job_list

    @magic_arguments()
    @line_cell_magic
    @needs_local_scope
    @argument("-i", "--id", type=str, default=None, help="Exec ID of Job ")
    def progress(self, line, cell="", local_ns=None):
        user_input = parse_argstring(self.kill, line)
        status, result = self.linkis_client.progress(user_input.id)
        if status:
            print(result)
        else:
            raise Exception("HTTPException")

    @magic_arguments()
    @line_cell_magic
    @needs_local_scope
    @argument("-i", "--id", type=str, default=None, help="Exec ID of Job ")
    def kill(self, line, cell="", local_ns=None):
        kill_input = parse_argstring(self.kill, line)
        #        print(kill_input.id)
        status, result = self.linkis_client.kill(kill_input.id)
        print(result)
        if status:
            print("Succeed")
        else:
            print("ERROR")

    @magic_arguments()
    @line_cell_magic
    @needs_local_scope
    @argument("-i", "--id", type=str, default=None, help="Exec ID of Job ")
    def log(self, line, cell="", local_ns=None):
        user_input = parse_argstring(self.kill, line)
        status, result = self.linkis_client.log(user_input.id)
        if status:
            print(result)
        else:
            raise Exception("HTTPException")

    @magic_arguments()
    @line_cell_magic
    @needs_local_scope
    @argument("-i", "--id", type=str, default=None, help="Exec ID of Job ")
    def status(self, line, cell="", local_ns=None):
        user_input = parse_argstring(self.kill, line)
        status, result = self.linkis_client.status(user_input.id)
        if status:
            print(result)
        else:
            raise Exception("HTTPException")

    @magic_arguments()
    @line_cell_magic
    @needs_local_scope
    @argument("-o", "--output", type=str, default=None, help="Output of Job ")
    def listengine(self, line, cell="", local_ns=None):
        user_input = parse_argstring(self.listengine, line)
        status, engine_list = self.linkis_client.engines()
        if not status:
            raise Exception("Http Exception")
        if user_input.output is not None:
            self.shell.user_ns[user_input.output] = engine_list
        else:
            return engine_list

    @magic_arguments()
    @line_cell_magic
    @needs_local_scope
    @argument("-i", "--instance", type=str, default=None, help="Instance of Engine ")
    def enginekill(self, line, cell="", local_ns=None):
        user_input = parse_argstring(self.enginekill, line)
        status, result = self.linkis_client.engine_kill(user_input.instance)
        if status:
            print("Success")
        else:
            raise Exception("HTTPException")

    def load_ipython_extension(ip):
        ip.register_magics(LinkisMagic)

    # 优化log显示，分为详细log和关键性息
    def log_detail(self):
        pass

    @magic_arguments()
    @line_cell_magic
    @needs_local_scope
    def flashcookies(self, line, cell="", local_ns=None):
        user_input = parse_argstring(self.listengine, line)
        self.linkis_client.refresh_cookies()
        self.ipython_display.display("Refresh Cookies Successful.")


ip = get_ipython()
ip.register_magics(LinkisMagic)
