<template>
  <div>
    <breadcrumb-nav></breadcrumb-nav>
    <el-form ref="queryValidate"
             :model="query"
             class="query">
      <el-row>
        <el-col :span="8">
          <el-form-item>
            <el-input v-model="query.query_str"
                      :placeholder="$t('expExeRecord.enterSearchContent')"></el-input>
          </el-form-item>
        </el-col>
        <el-col :span="16">
          <el-button type="primary"
                     class="margin-l30"
                     icon="el-icon-search"
                     @click="filterListData">
            {{ $t("common.filter") }}
          </el-button>
          <el-button type="warning"
                     icon="el-icon-refresh-left"
                     @click="resetListData">
            {{ $t("common.reset") }}
          </el-button>
          <el-button type="primary"
                     icon="plus-circle-o"
                     @click="dialogVisible=true">
            {{$t('distribute.creationExpriment')}}
          </el-button>
          <el-button type="primary"
                     icon="plus-circle-o"
                     @click="showGPU">
            {{$t('distribute.createDLTask')}} (GPU)
          </el-button>
        </el-col>
      </el-row>
    </el-form>
    <el-table :data="dataList"
              border
              style="width: 100%"
              :empty-text="$t('common.noData')">
      <el-table-column prop="exp_name"
                       :show-overflow-tooltip="true"
                       :label="$t('distribute.experimentName')" />
      <el-table-column prop="group_name"
                       :show-overflow-tooltip="true"
                       :label="$t('user.userGroup')" />
      <el-table-column prop="exp_desc"
                       :show-overflow-tooltip="true"
                       :label="$t('distribute.experimentDesc')" />
      <el-table-column prop="tag_str"
                       :show-overflow-tooltip="true"
                       :label="$t('distribute.experimentTag')" />
      <el-table-column prop="dss_dss_project_name"
                       :show-overflow-tooltip="true"
                       :label="$t('distribute.projectName')" />
      <el-table-column prop="dss_dss_flow_name"
                       :show-overflow-tooltip="true"
                       :label="$t('distribute.workflowName')" />
      <el-table-column prop="create_type"
                       width="85"
                       :show-overflow-tooltip="true"
                       :label="$t('DI.type')" />
      <el-table-column prop="dss_dss_flow_version"
                       :show-overflow-tooltip="true"
                       :label="$t('distribute.workflowVersion')" />
      <el-table-column prop="exp_create_time"
                       :show-overflow-tooltip="true"
                       :label="$t('distribute.createTime')" />
      <el-table-column prop="exp_modify_time"
                       :show-overflow-tooltip="true"
                       :label="$t('distribute.updateTime')" />
      <el-table-column prop="createUserName"
                       :show-overflow-tooltip="true"
                       :label="$t('distribute.creator')" />
      <el-table-column prop="modifyUserName"
                       :show-overflow-tooltip="true"
                       :label="$t('distribute.modifier')" />
      <el-table-column :label="$t('common.operation')"
                       width="85">
        <template slot-scope="scope">
          <el-dropdown size="medium"
                       @command="handleCommand">
            <el-button type="text"
                       class="multi-operate-btn"
                       icon="el-icon-more"
                       round></el-button>
            <el-dropdown-menu slot="dropdown">
              <el-dropdown-item :command="beforeHandleCommand('editeFlow',scope.row)">{{ $t('distribute.experimentEdit') }}</el-dropdown-item>
              <el-dropdown-item :command="beforeHandleCommand('runExperiment',scope.row)">{{ $t('distribute.experimentExecution') }}</el-dropdown-item>
              <el-dropdown-item :command="beforeHandleCommand('deleteItem',scope.row)">{{ $t('distribute.deleteExperiment') }}</el-dropdown-item>
              <el-dropdown-item :command="beforeHandleCommand('downloadFlow',scope.row)">{{ $t('distribute.exportExperiment') }}</el-dropdown-item>
              <el-dropdown-item :command="beforeHandleCommand('toExperimentHistory',scope.row)">{{ $t('distribute.executionHistory') }}</el-dropdown-item>
              <el-dropdown-item :command="beforeHandleCommand('updateItem',scope.row)">{{ $t('distribute.informationEdit') }}</el-dropdown-item>
              <el-dropdown-item :command="beforeHandleCommand('addRecord',scope.row)">{{$t('distribute.tracking')}}</el-dropdown-item>
              <!-- <el-dropdown-item :command="beforeHandleCommand('stopExperiment',scope.row)">实验执行列表</el-dropdown-item> -->
            </el-dropdown-menu>
          </el-dropdown>
        </template>
      </el-table-column>
    </el-table>
    <el-pagination @size-change="handleSizeChange"
                   @current-change="handleCurrentChange"
                   :current-page="pagination.pageNumber"
                   :page-sizes="sizeList"
                   :page-size="pagination.pageSize"
                   layout="total, sizes, prev, pager, next, jumper"
                   :total="pagination.totalPage"
                   background>
    </el-pagination>
    <el-dialog :title="modify?$t('distribute.modifyExpriment'):$t('distribute.creationExpriment')"
               :visible.sync="dialogVisible"
               @close="clearDialogForm"
               custom-class="dialog-style"
               width="600px">
      <el-form ref="formValidate"
               :model="form"
               :rules="ruleValidate"
               class="add-form"
               label-width="120px">
        <el-form-item :label="$t('distribute.creationWay')"
                      prop="type">
          <el-select v-model="form.type"
                     :placeholder="$t('distribute.creationWayTip')">
            <el-option value="flow"
                       :label="$t('distribute.visualEditing')"></el-option>
            <el-option value="import"
                       :label="$t('distribute.externalImport')"></el-option>
          </el-select>
        </el-form-item>
        <div v-if="form.type==='flow'">
          <el-form-item :label="$t('distribute.experimentName')"
                        :placeholder="$t('distribute.experimentNameTip')"
                        prop="exp_name">
            <el-input v-model="form.exp_name" />
          </el-form-item>
          <el-form-item :label="$t('user.userGroup')"
                        placeholder="请选择项目组"
                        prop="group_name">
            <el-select v-model="form.group_name"
                       filterable>
              <el-option v-for="(item,index) in groupList"
                         :key="index"
                         :value="item.name"
                         :label="item.name"></el-option>
            </el-select>
          </el-form-item>
          <el-form-item :label="$t('distribute.experimentTag')">
            <el-tag :key="index"
                    v-for="(tag,index) in form.tag_list"
                    closable
                    :disable-transitions="false"
                    @close="handleClose(tag)">
              {{tag.exp_tag}}
            </el-tag>
            <el-input class="input-new-tag"
                      v-if="inputVisible"
                      v-model="inputValue"
                      ref="saveTagInput"
                      size="small"
                      @keyup.enter.native="handleInputConfirm"
                      @blur="handleInputConfirm">
            </el-input>
            <el-button v-else
                       class="button-new-tag"
                       size="mini"
                       @click="showInput">+ {{$t('label.addlabels')}}</el-button>
          </el-form-item>
          <el-form-item :label="$t('distribute.experimentDesc')"
                        prop="exp_desc">
            <el-input type="textarea"
                      :placeholder="$t('distribute.experimentDescTip')"
                      v-model="form.exp_desc" />
          </el-form-item>
        </div>
        <div v-else>
          <el-upload class="upload"
                     drag
                     :action="baseUrl"
                     :headers="{'Mlss-Userid':userId}"
                     :data="uploadParam"
                     :before-upload="beforeUpload"
                     :on-success="uploadSuccess"
                     :on-error="uploadError">
            <i class="el-icon-upload"></i>
            <div class="el-upload__text">{{$t('distribute.dragFile')}}<em>{{$t('distribute.clickUpload')}}</em></div>
            <div class="el-upload__tip"
                 slot="tip">{{$t('distribute.uploadFileTip')}}</div>
          </el-upload>
        </div>
      </el-form>
      <span slot="footer"
            class="dialog-footer">
        <el-button type="primary"
                   :disabled="btnDisable||form.type==='import'"
                   @click="subInfo">{{ $t('common.save') }}</el-button>
        <el-button @click="dialogVisible=false">{{ $t('common.cancel') }}</el-button>
      </span>
    </el-dialog>
    <create-gpu ref="gpuDialog"
                :spaceOptionList="spaceOptionList"
                :localPathList="localPathList"></create-gpu>
  </div>
</template>
<script>
import mixin from '../../util/modelMixin'
import { Tag } from 'element-ui'
import CreateGpu from '../../components/CreateGPU'
import util from '../../util/common'
export default {
  mixins: [mixin],
  components: {
    ElTag: Tag,
    CreateGpu
  },
  data () {
    return {
      query: {
        query_str: ''
      },
      dialogVisible: false,
      spaceOptionList: [],
      localPathList: [],
      dataList: [],
      groupList: [],
      form: {
        type: 'flow',
        group_name: '',
        tag_list: [],
        exp_name: '',
        exp_desc: '',
        id: ''
      },
      uploadParam: {},
      deleteTagArr: [],
      inputVisible: false,
      inputValue: '',
      fileSize: 100,
      fileName: '',
      btnDisable: false,
      modify: false,
      userId: localStorage.getItem('userId')
    }
  },
  computed: {
    getUrl () {
      return `/di/${this.FesEnv.diApiVersion}/experiments`
    },
    baseUrl () {
      return `${process.env.VUE_APP_BASE_SURL}/di/${this.FesEnv.diApiVersion}/experiment/import`
    },
    ruleValidate () {
      this.resetFormFields()
      return {
        type: [
          { required: true, message: this.$t('distribute.creationWayNoEmpty') }
        ],
        exp_name: [
          { required: true, message: this.$t('distribute.experimentNameNoEmpty') },
          { pattern: new RegExp(/^[a-zA-Z][0-9a-zA-Z-_]*$/), message: this.$t('distribute.experimentNameReg') }
        ],
        exp_desc: [
          { required: true, message: this.$t('distribute.experimentDescReq') }
        ],
        group_name: [
          { required: true, message: this.$t('ns.userGroupReq') }
        ]
      }
    }
  },
  mounted () {
    this.getListData()
    this.getDISpaceOption()
    this.getDIStoragePath()
    this.getGroupList()
  },
  methods: {
    showGPU () {
      this.$refs.gpuDialog.dialogVisible = true
    },
    getListData () {
      let query = {
        page: this.pagination.pageNumber,
        size: this.pagination.pageSize,
        username: localStorage.getItem('userId')
      }
      if (this.query.query_str) {
        query.query_str = this.query.query_str
      }
      this.FesApi.fetch(this.getUrl, query, 'get').then(rst => {
        rst.experiments && rst.experiments.forEach(item => {
          if (item.tag_list) {
            item.tag_str = item.tag_list[0].exp_tag
            for (let i = 1; i < item.tag_list.length; i++) {
              item.tag_str += ',' + item.tag_list[i].exp_tag
            }
          }
          item.exp_create_time = util.transDate(new Date(item.exp_create_time).getTime())
          item.exp_modify_time = util.transDate(new Date(item.exp_modify_time).getTime())
          item.createUserName = item.create_user.Name
          item.modifyUserName = item.modify_user.Name
        })
        this.dataList = Object.freeze(rst.experiments)
        this.pagination.totalPage = rst.total || 0
      })
    },
    subInfo () {
      this.$refs.formValidate.validate((valid) => {
        if (valid) {
          this.btnDisable = true
          let method = 'post'
          let param = {
            exp_name: this.form.exp_name,
            group_name: this.form.group_name,
            exp_desc: this.form.exp_desc,
            tag_list: this.form.tag_list
          }
          let url = `/di/${this.FesEnv.diApiVersion}/experiment`
          if (this.modify) {
            method = 'put'
            url += `/${this.form.id}`
            // param.tag_list.concat(this.deleteTagArr)
          }
          this.FesApi.fetch(url, param, method).then((res) => {
            this.toast()
            this.dialogVisible = false
            this.btnDisable = false
            if (!this.modify) {
              this.$router.push({
                name: 'experimentFlow',
                query: {
                  readonly: false,
                  exp_id: res.id,
                  exp_name: res.exp_name
                }
              })
            } else {
              this.getListData()
            }
          }, () => {
            this.btnDisable = false
          })
        }
      })
    },
    deleteItem (trData) {
      const url = `/di/${this.FesEnv.diApiVersion}/experiment/${trData.id}`
      this.deleteListItem(url)
    },
    updateItem (trData) {
      this.dialogVisible = true
      this.modify = true
      this.uploadParam.experimentId = trData.id
      setTimeout(() => {
        let param = {}
        Object.keys(this.form).forEach((item) => {
          param[item] = trData[item]
        })
        param.type = 'flow'
        let tagList = []
        if (trData.tag_list) {
          for (let item of trData.tag_list) {
            tagList.push({ exp_id: item.exp_id, exp_tag: item.exp_tag })
          }
        }
        param.tag_list = tagList
        this.form = param
      }, 0)
    },
    editeFlow (trData) {
      this.$router.push({
        name: 'experimentFlow',
        query: {
          readonly: trData.create_type === 'WTSS',
          exp_id: trData.id,
          exp_name: trData.exp_name
        }
      })
    },
    addRecord (trData) {
      this.$router.push('/experiment/addRecord/' + trData.mlflow_exp_id)
    },
    handleClose (tag) {
      let tagList = this.form.tag_list
      for (let i = 0; i < tagList.length; i++) {
        if (tag.exp_tag === tagList[i].exp_tag) {
          this.form.tag_list.splice(i, 1)
        }
        // if (tag.exp_id !== -1) {
        //   if (tag.exp_id === tagList[i].exp_id) {
        //     this.form.tag_list.splice(i, 1)
        //   }
        // } else {
        //   if (tag.exp_tag === tagList[i].exp_tag) {
        //     this.form.tag_list.splice(i, 1)
        //   }
        // }
      }
      // if (this.modify && tag.exp_id !== -1) {
      //   this.deleteTagArr.push({ exp_id: -tag.exp_id, exp_tag: tag.exp_tag })
      // }
    },
    showInput () {
      this.inputVisible = true
      this.$nextTick(_ => {
        this.$refs.saveTagInput.$refs.input.focus()
      })
    },
    handleInputConfirm () {
      let inputValue = this.inputValue
      if (inputValue) {
        this.form.tag_list.push({ exp_id: -1, exp_tag: inputValue })
      }
      this.inputVisible = false
      this.inputValue = ''
    },
    uploadFile (file) {
      this.$message.success(this.$t('distribute.uploadSuccess'))
    },
    downloadFlow (trData) {
      const url = `${process.env.VUE_APP_BASE_SURL}/di/v1/experiment/${trData.id}/export`
      this.downFile(url, trData.exp_name)
    },
    beforeUpload (file) {
      const fileTpye = file.type.indexOf('zip') > -1
      const isLt512M = file.size / 1024 / 1024 < this.fileSize
      let limit = this.fileName === ''
      if (!fileTpye) {
        this.$message.error(this.$t('common.fileFormat'))
      }
      if (!isLt512M) {
        this.$message.error(this.$t('common.fileSize' + this.fileSize + 'BM'))
      }
      if (!limit) {
        this.$message.error(this.$t('common.fileLimit'))
        return false
      }
      if (fileTpye && isLt512M && isLt512M) {
        this.loading = true
      }
      return fileTpye && isLt512M
    },
    runExperiment (trData) {
      this.FesApi.fetch(`/di/${this.FesEnv.diApiVersion}/experimentRun/${trData.id}`, { exp_exec_type: 'MLFLOW' }, 'post').then((res) => {
        this.$router.push({
          name: 'experimentFlow',
          query: {
            exp_name: trData.exp_name,
            readonly: true,
            exp_id: trData.id,
            exec_id: res.id
          }
        })
      })
    },
    uploadSuccess (file) {
      this.$router.push({
        name: 'experimentFlow',
        query: {
          exp_name: file.result.exp_name,
          readonly: false,
          exp_id: file.result.id
        }
      })
    },
    uploadError (file) {
      const message = JSON.parse(file.message)
      this.$message.error(message.result)
    },
    toExperimentHistory (trData) {
      this.$router.push({
        name: 'expExeRecord',
        query: {
          exp_id: trData.id
        }
      })
    },
    clearDialogForm () {
      this.modify = false
      this.form.tag_list = []
      this.form.exp_desc = ''
      delete this.form.exp_id
      setTimeout(() => {
        this.$refs.formValidate.resetFields()
      }, 0)
    },
    resetListData () {
      this.pagination.pageNumber = 1
      this.query.query_str = ''
      this.$refs.queryValidate.resetFields()
      this.getListData()
    }
  }
}
</script>
<style lang="scss" scoped>
.el-tag + .el-tag {
  margin-left: 10px;
}
.button-new-tag {
  margin-left: 10px;
  height: 32px;
  line-height: 30px;
  padding-top: 0;
  padding-bottom: 0;
  border-radius: 20px;
}
.input-new-tag {
  width: 90px;
  margin-left: 10px;
  vertical-align: bottom;
}
.upload {
  margin-left: 120px;
}
.el-dropdown-menu {
  a {
    text-decoration: none;
    color: #606266;
    &:hover {
      color: #66b1ff;
    }
  }
}
</style>
