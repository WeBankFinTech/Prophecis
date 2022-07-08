<template>
  <div>
    <breadcrumb-nav></breadcrumb-nav>
    <el-form ref="queryValidate"
             :model="query"
             class="query"
             :rules="ruleValidate"
             label-width="85px">
      <el-row>
        <el-col :span="6">
          <el-form-item :label="$t('ns.nameSpace')"
                        prop="namespace">
            <el-select v-model="query.namespace"
                       filterable
                       :placeholder="$t('ns.nameSpacePro')">
              <el-option v-for="item in spaceOptionList"
                         :label="item"
                         :key="item"
                         :value="item">
              </el-option>
            </el-select>
          </el-form-item>
        </el-col>
        <el-col :span="6">
          <el-form-item :label="$t('user.user')"
                        prop="userid">
            <el-select v-model="query.userid"
                       clearable
                       filterable
                       :placeholder="$t('user.userPro')">
              <el-option v-for="item in userOptionList"
                         :label="item.name"
                         :key="item.id"
                         :value="item.name">
              </el-option>
            </el-select>
          </el-form-item>
        </el-col>
        <el-col :span="12">
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
          <el-button v-show="false">
            {{currentExpRunId}}
          </el-button>
        </el-col>
      </el-row>
    </el-form>
    <el-table :data="dataList"
              border
              style="width: 100%"
              :empty-text="$t('common.noData')">
      <el-table-column prop="model_id"
                       :show-overflow-tooltip="true"
                       :label="$t('DI.taskExecutionID')"
                       min-width="120" />
      <el-table-column prop="name"
                       :show-overflow-tooltip="true"
                       :label="$t('DI.jobName')"
                       min-width="120" />
      <el-table-column prop="expRunId"
                       :show-overflow-tooltip="true"
                       :label="$t('DI.experimentExecutionID')"
                       min-width="120" />
      <el-table-column prop="expName"
                       :show-overflow-tooltip="true"
                       :label="$t('DI.experimentName')"
                       min-width="100" />
      <el-table-column prop="submission_timestamp"
                       :show-overflow-tooltip="true"
                       :label="$t('DI.submissionTime')"
                       min-width="100" />
      <el-table-column prop="completed_timestamp"
                       :show-overflow-tooltip="true"
                       :label="$t('DI.endTime')"
                       min-width="100" />
      <el-table-column prop="user_id"
                       :show-overflow-tooltip="true"
                       :label="$t('DI.committer')"
                       min-width="80" />
      <el-table-column prop="status"
                       :show-overflow-tooltip="true"
                       :label="$t('DI.executionStatus')"
                       min-width="90">
        <template slot-scope="scope">
          <status-tag :status="scope.row.status"
                      :statusObj="statusObj"></status-tag>
        </template>
      </el-table-column>
      <el-table-column :label="$t('common.operation')"
                       min-width="85">
        <template slot-scope="scope">
          <el-dropdown size="medium"
                       @command="handleCommand">
            <el-button type="text"
                       class="multi-operate-btn"
                       icon="el-icon-more"
                       round></el-button>
            <el-dropdown-menu slot="dropdown">
              <el-dropdown-item v-if="scope.row.JobType !== 'tfos'"
                                :command="beforeHandleCommand('taskDetail',scope.row)">{{$t('DI.taskDetail')}}</el-dropdown-item>
              <el-dropdown-item v-if="scope.row.JobType!=='MLPipeline'&&scope.row.JobType !== 'tfos'"
                                :command="beforeHandleCommand('copyJob',scope.row)">{{$t('common.copy')}}</el-dropdown-item>
              <el-dropdown-item :command="beforeHandleCommand('deleteItem',scope.row)">{{$t('common.delete')}}</el-dropdown-item>
              <el-dropdown-item v-if="scope.row.JobType !== 'tfos'"
                                :command="beforeHandleCommand('viewLogFun',scope.row)">{{$t('DI.viewlog')}}</el-dropdown-item>
              <el-dropdown-item v-if="(scope.row.status==='PENDING'||scope.row.status==='QUEUED'||scope.row.status==='RUNNING')&&scope.row.JobType !== 'tfos'"
                                :command="beforeHandleCommand('downloadJob',scope.row)">{{$t('common.export')}}</el-dropdown-item>
              <el-dropdown-item v-if="scope.row.status==='PENDING'||scope.row.status==='QUEUED'||scope.row.status==='RUNNING'"
                                :command="beforeHandleCommand('killJob',scope.row)">{{$t('DI.termination')}}</el-dropdown-item>
              <el-dropdown-item v-if="scope.row.JobType !== 'tfos'"
                                :command="beforeHandleCommand('repullJob',scope.row)">{{$t('DI.retry')}}</el-dropdown-item>
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
    <el-dialog :title="detailTitle"
               :visible.sync="detailVisible"
               custom-class="dialog-style"
               width="1100px">
      <job-detail :detail-obj="detailObj"
                  :space-option-list="spaceOptionList" />
    </el-dialog>
    <el-dialog :title="detailTitle"
               :visible.sync="hadoopDetailVisible"
               custom-class="dialog-style"
               width="1150px">
    </el-dialog>
    <!-- <el-dialog :visible.sync="storageVisible"
               :title="title"
               width="900px">
      <el-table :data="storageList"
                border
                :no-data-text="$t('common.noData')">
        <el-table-column prop="bucket"
                         :label="$t('DI.bucket')" />
        <el-table-column prop="name"
                         :label="$t('DI.name')" />
        <el-table-column prop="path"
                         :label="$t('DI.path')"
                         width="300%" />
        <el-table-column prop="type"
                         :label="$t('DI.type')" />
      </el-table>
    </el-dialog> -->
    <create-gpu ref="gpuDialog"
                :spaceOptionList="spaceOptionList"
                :localPathList="localPathList"
                @getListData="getListData"></create-gpu>
  </div>
</template>
<script type="text/ecmascript-6">
import handleDIDetailMixin from '../../util/handleDIDetailMixin'
import modelMixin from '../../util/modelMixin'
import JobDetail from './JobDetail'
import StatusTag from '../../components/StatusTag'
import CreateGpu from '../../components/CreateGPU'
export default {
  mixins: [modelMixin, handleDIDetailMixin],
  components: {
    JobDetail,
    StatusTag,
    CreateGpu
  },
  data: function () {
    return {
      detailVisible: false,
      // storageVisible: false,
      hadoopDetailVisible: false,
      query: {
        namespace: '',
        userid: ''
      },
      localPathList: [],
      detailObj: {}, // 详情对象
      dataModel: true, // true为数据存储 false为结果存储
      loading: false,
      dataList: [],
      // storageList: [],
      userOptionList: [],
      spaceOptionList: [],
      trainingId: '',
      disabled: false,
      statusObj: {
        COMPLETED: 'success',
        PENDING: '',
        RUNNING: '',
        QUEUED: 'info',
        FAILED: 'danger',
        CANCELLED: 'warning'
      },
      historyExpRunId: this.$route.query.expRunId
    }
  },
  mounted () {
    this.getListData()
    this.getDIUserOption()
    this.getDISpaceOption()
    this.getDIStoragePath()
  },
  computed: {
    title: function () {
      return this.dataModel ? this.$t('DI.dataStorage') : this.$t('DI.resultStorage')
    },
    detailTitle: function () {
      let title = ''
      if (this.$t('DI.alarmDetail') === '告警详情') {
        title = `任务${this.trainingId}详情`
      } else {
        title = this.$t('DI.alarmDetail') + this.trainingId
      }
      return title
    },
    noDataText: function () {
      return !this.loading ? this.$t('common.noData') : ''
    },
    currentExpRunId () {
      // 监听$route.query.expRunId 如果id有值后面没有值证明先从实验执行记录点击执行任务列表 在点击的任务执行记录，应重新请求任务执行全部记录
      if (this.historyExpRunId && !this.$route.query.expRunId) {
        this.resetListData()
      }
      return this.$route.query.exp_id
    },
    imageOptionList: function () {
      return this.FesEnv.DI.imageOption
    },
    ruleValidate: function () {
      // 切换语言表单报错重置表单
      this.resetQueryFields()
      return {
        name: [
          { required: true, message: this.$t('DI.trainingjobNameReq') },
          { type: 'string', pattern: new RegExp(/^[0-9a-zA-Z-]*$/), message: this.$t('DI.trainingjobNameFormat') }
        ]
      }
    }
  },
  methods: {
    // goMlFlow () {
    //   window.open(window.encodeURI(this.FesEnv.DI.mlFlowUrl), '_blank')
    // },
    showGPU () {
      this.$refs.gpuDialog.dialogVisible = true
    },
    taskDetail (trData) {
      console.log('taskDetail', trData)
      this.trainingId = trData.model_id
      this.detailObj = {}
      Object.assign(this.detailObj, trData)
      if (trData.job_type === 'tfos') {
        this.hadoopDetailVisible = true
      } else {
        this.detailVisible = true
      }
    },
    viewLogFun (trData) {
      this.$router.push({
        name: 'jobExeRecordDetail',
        params: {
          modelId: trData.model_id
        }
      })
    },
    deleteItem (trData) {
      this.$confirm(this.$t('DI.terminateJobPro'), this.$t('common.prompt')).then((index) => {
        this.FesApi.fetchUT(`/di/${this.FesEnv.diApiVersion}/models/${trData.model_id}`, { version: '200' }, 'delete').then(() => {
          this.deleteItemSuccess()
        })
      }).catch(() => { })
    },
    downloadJob (trData) {
      const url = `${process.env.VUE_APP_BASE_SURL}/di/v1/models/${trData.model_id}/export`
      this.downFile(url, trData.fileName || 'jobFile')
    },
    killJob (trData) {
      this.FesApi.fetchUT(`/di/${this.FesEnv.diApiVersion}/models/${trData.model_id}/kill`, 'get').then(() => {
        this.getListData()
        this.toast()
      })
    },
    copyJob (trData) {
      if (trData.job_alert) {
        // if (trData.job_alert.length > 1) {
        //   for (let i = 1; i <= trData.job_alert.length; i++) {
        //     this.addAlarmValidate(i)
        //   }
        // }
        trData.job_alert.forEach((item, index) => {
          if (item.fixTime && item.fixTime.deadline_checker) {
            trData.job_alert[index].fixTime.deadlineChecker = new Date(item.fixTime.deadline_checker).getTime()
          }
        })
        trData.setAlarm = trData.job_alert
      }
      if (trData.ps_image) {
        this.judgeImageType('ps_imageType', 'ps_image', 'ps_imageInput', trData.ps_image, trData)
      }
      this.judgeImageType('imageType', 'imageOption', 'imageInput', trData.image, trData)
      setTimeout(() => {
        const gpuObj = this.$refs.gpuDialog
        // delete form.job_alert
        Object.assign(gpuObj.form, trData)
        if (trData.filePath) {
          gpuObj.fileName = trData.fileName
          gpuObj.codeFile = trData.filePath
        }
      }, 100)
      this.showGPU()
    },
    repullJob (trData) {
      this.FesApi.fetchUT(`/di/v1/models/${trData.model_id}/retry`, { version: '200' }, 'get').then((res) => {
        this.toast()
        this.getListData()
      })
    },
    judgeImageType (imageType, image, imageInput, trDataImage, trData) {
      if (this.imageOptionList.indexOf(trDataImage) > -1) {
        trData[imageType] = 'Standard'
        trData[image] = trDataImage
      } else {
        trData[imageType] = 'Custom'
        trData[imageInput] = trDataImage
      }
    },
    getListData () {
      let search = {
        version: '200',
        page: this.pagination.pageNumber,
        size: this.pagination.pageSize
      }
      if (this.query.namespace || this.query.userid) {
        for (let key in this.query) {
          if (this.query[key]) {
            search[key] = this.query[key]
          }
        }
      }

      if (this.$route.query.expRunId) {
        search.expRunId = this.$route.query.expRunId
      }
      this.FesApi.fetchUT(`/di/${this.FesEnv.diApiVersion}/models`, search, 'get').then(rst => {
        if (rst && rst.models) {
          let models = []
          rst.models.forEach(item => {
            let modelsItem = this.basicInfoProcess(item)
            modelsItem.expRunId = item.expRunId
            modelsItem.expName = item.expName
            modelsItem.JobType = item.JobType
            models.push(modelsItem)
          })
          this.dataList = Object.freeze(models)
          this.pagination.totalPage = rst.total || 0
        }
      })
    }
  }
}
</script>
<style lang="scss" scoped>
.storage-td {
  color: #409eff;
}
</style>
