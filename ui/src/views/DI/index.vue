<template>
  <div>
    <breadcrumb-nav></breadcrumb-nav>
    <el-form ref="queryValidate" :model="query" class="query" :rules="ruleValidate" label-width="85px">
      <el-row>
        <el-col :span="6">
          <el-form-item :label="$t('ns.nameSpace')" prop="namespace">
            <el-select v-model="query.namespace" filterable clearable :placeholder="$t('ns.nameSpacePro')">
              <el-option v-for="item in spaceOptionList" :label="item" :key="item" :value="item">
              </el-option>
            </el-select>
          </el-form-item>
        </el-col>
        <el-col :span="6">
          <el-form-item :label="$t('user.user')" prop="userid">
            <el-select v-model="query.userid" clearable filterable :placeholder="$t('user.userPro')">
              <el-option v-for="item in userOptionList" :label="item.name" :key="item.id" :value="item.name">
              </el-option>
            </el-select>
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-button type="primary" class="margin-l30" icon="el-icon-search" @click="filterListData">
            {{ $t("common.filter") }}
          </el-button>
          <el-button type="warning" icon="el-icon-refresh-left" @click="resetListData">
            {{ $t("common.reset") }}
          </el-button>
          <el-button v-show="false">
            {{currentExpRunId}}
          </el-button>
          <el-button type="primary" icon="plus-circle-o" @click="showGPU">
            {{ $t("DI.createGPU") }}
          </el-button>
          <!-- <el-button type="primary" icon="plus-circle-o" @click="showHadoop">
            {{ $t("DI.createHadoop") }}
          </el-button> -->
          <!-- <el-button type="primary"
                     icon="plus-circle-o"
                     @click="goMlFlow">
            {{$t("DI.createMlFlow")}}
          </el-button> -->
        </el-col>
      </el-row>
    </el-form>
    <el-table :data="dataList" border style="width: 100%" :empty-text="$t('common.noData')">
      <el-table-column prop="model_id" :label="$t('DI.modelId')" min-width="120" />
      <el-table-column prop="name" :label="$t('DI.jobName')" min-width="120" />
      <el-table-column prop="user_id" :label="$t('user.user')" min-width="80" />
      <el-table-column prop="namespace" :label="$t('ns.nameSpace')" min-width="155" />
      <el-table-column prop="cpus" :label="$t('DI.cpu')" min-width="65" />
      <el-table-column prop="gpus" :label="$t('DI.gpu')" min-width="65" />
      <el-table-column prop="image" :label="$t('DI.image')" min-width="150" />
      <el-table-column prop="trainingData" :label="$t('DI.dataStorage')" min-width="80">
        <template slot-scope="scope">
          <el-button type="text" @click="storageDetail(scope.row,7)">
            {{scope.row.trainingData}}
          </el-button>
        </template>
      </el-table-column>
      <el-table-column prop="trainingResults" :label="$t('DI.resultStorage')" min-width="80">
        <template slot-scope="scope">
          <el-button type="text" @click="storageDetail(scope.row,8)">
            {{scope.row.trainingResults}}
          </el-button>
        </template>
      </el-table-column>
      <el-table-column prop="status" :label="$t('DI.status')" min-width="90" />
      <el-table-column prop="submission_timestamp" :label="$t('DI.createTime')" min-width="105" />
      <el-table-column :label="$t('common.operation')" min-width="85">
        <template slot-scope="scope">
          <el-dropdown trigger="click" @command="handleCommand">
            <el-button type="text">{{$t("DI.operationList")}}</el-button>
            <el-dropdown-menu slot="dropdown">
              <el-dropdown-item :command="beforeHandleCommand('taskDetail',scope.row)">{{$t('DI.taskDetail')}}</el-dropdown-item>
              <el-dropdown-item :command="beforeHandleCommand('deleteItem',scope.row)">{{$t('common.delete')}}</el-dropdown-item>
              <el-dropdown-item :command="beforeHandleCommand('viewLogFun',scope.row)">{{$t('DI.viewlog')}}</el-dropdown-item>
              <el-dropdown-item :command="beforeHandleCommand('copyJob',scope.row)">{{$t('common.copy')}}</el-dropdown-item>
              <el-dropdown-item :command="beforeHandleCommand('downloadJob',scope.row)">{{$t('common.export')}}</el-dropdown-item>
            </el-dropdown-menu>
          </el-dropdown>
        </template>
      </el-table-column>
    </el-table>
    <el-pagination @size-change="handleSizeChange" @current-change="handleCurrentChange" :current-page="pagination.pageNumber" :page-sizes="sizeList" :page-size="pagination.pageSize" layout="total, sizes, prev, pager, next, jumper" :total="pagination.totalPage" background>
    </el-pagination>
    <el-dialog :title="detailTitle" :visible.sync="detailVisible" custom-class="dialog-style" width="1100px">
      <job-detail :detail-obj="detailObj" :space-option-list="spaceOptionList" />
    </el-dialog>
    <el-dialog :visible.sync="storageVisible" :title="title" width="900px">
      <el-table :data="storageList" border :no-data-text="$t('common.noData')">
        <el-table-column prop="bucket" :label="$t('DI.bucket')" />
        <el-table-column prop="name" :label="$t('DI.name')" />
        <el-table-column prop="path" :label="$t('DI.path')" width="300%" />
        <el-table-column prop="type" :label="$t('DI.type')" />
      </el-table>
    </el-dialog>
    <create-gpu ref="gpuDialog" :spaceOptionList="spaceOptionList" :localPathList="localPathList" @getListData="getListData"></create-gpu>
  </div>
</template>
<script type="text/ecmascript-6">
import handleDIDetailMixin from '../../util/handleDIDetailMixin'
import modelMixin from '../../util/modelMixin'
import JobDetail from './JobDetail'
import CreateGpu from '../../components/CreateGPU'
import axios from 'axios'
export default {
  mixins: [modelMixin, handleDIDetailMixin],
  components: {
    JobDetail,
    CreateGpu
  },
  data: function () {
    return {
      detailVisible: false,
      storageVisible: false,
      query: {
        namespace: '',
        userid: ''
      },
      localPathList: [],
      detailObj: {}, // 详情对象
      dataModel: true, // true为数据存储 false为结果存储
      loading: false,
      dataList: [],
      storageList: [],
      userOptionList: [],
      spaceOptionList: [],
      trainingId: '',
      disabled: false,
      uiServer: {
        'bdp-ui': 'bdp',
        'bdap-ui': 'bdap',
        'bdap-safe': 'bdapsafe'
      },
      statusObj: {
        COMPLETED: 'success',
        PENDING: '',
        QUEUED: 'info',
        FAILED: 'danger'
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
    showGPU () {
      this.$refs.gpuDialog.dialogVisible = true
    },
    taskDetail (trData) {
      this.trainingId = trData.model_id
      this.detailObj = {}
      Object.assign(this.detailObj, trData)
      this.detailVisible = true
    },
    viewLogFun (trData) {
      this.$router.push({
        name: 'DILogDetail',
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
      this.downFile(url, trData.fileName)
    },
    downFile (url, fileName) {
      axios({
        url: url,
        method: 'get',
        headers: {
          'Mlss-Userid': localStorage.getItem('userId'),
          'Authorization': 'Basic dGVkOndlbGNvbWUx'
        },
        responseType: 'arraybuffer'
      }).then((res) => {
        const blob = new Blob([res.data], { type: 'application/zip' })
        const url = URL.createObjectURL(blob)
        const a = document.createElement('a')
        a.style.display = 'none'
        a.href = url
        a.setAttribute('download', fileName)
        a.click()
      }, (err) => {
        this.uint8Msg(err)
      }).catch((err) => {
        this.uint8Msg(err)
      })
    },
    // 数据存储路径详情
    storageDetail (trData, index) {
      if (trData.data_stores) {
        switch (index) {
          case 7: // 数据存储
            this.storageVisible = true
            this.storageList = [trData.data_stores[0].connection]
            break
          case 8: // 结果存储
            this.storageVisible = true
            this.storageList = [trData.data_stores[1].connection]
            break
          default:
            break
        }
      }
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

      this.judgeImageType('imageType', 'imageOption', 'imageInput', trData.image, trData)
      if (trData.ps_image) {
        this.judgeImageType('ps_imageType', 'ps_image', 'ps_imageInput', trData.ps_image, trData)
      }
      setTimeout(() => {
        const gpuObj = this.$refs.gpuDialog
        Object.assign(gpuObj.form, trData)
        if (trData.filePath) {
          gpuObj.fileName = trData.fileName
          gpuObj.codeFile = trData.filePath
        }
      }, 100)
      this.showGPU()
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
      if (this.FesEnv.filterUiServer) {
        search.clusterName = this.uiServer[this.FesEnv.uiServer]
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
