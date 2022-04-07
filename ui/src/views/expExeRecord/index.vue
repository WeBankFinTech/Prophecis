<template>
  <div>
    <breadcrumb-nav></breadcrumb-nav>
    <el-form ref="queryValidate"
             :model="query"
             class="query">
      <!-- <el-row>
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
          </el-form>
        </el-col> -->
      <el-row>
        <el-col :span="6">
          <el-form-item>
            <el-input v-model="query.query_str"
                      :placeholder="$t('expExeRecord.enterSearchContent')"></el-input>
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
            {{currentExpId}}
          </el-button>
        </el-col>
      </el-row>
    </el-form>
    <el-table :data="dataList"
              border
              style="width: 100%"
              :empty-text="$t('common.noData')">
      <el-table-column prop="id"
                       :show-overflow-tooltip="true"
                       :label="$t('expExeRecord.executionId')"
                       min-width="120" />
      <el-table-column prop="exp_name"
                       :show-overflow-tooltip="true"
                       :label="$t('expExeRecord.exprimentName')"
                       min-width="120" />
      <el-table-column prop="exp_run_create_time"
                       :show-overflow-tooltip="true"
                       :label="$t('expExeRecord.creationTime')"
                       min-width="100" />
      <el-table-column prop="exp_run_end_time"
                       :label="$t('expExeRecord.endTime')"
                       :show-overflow-tooltip="true"
                       min-width="100" />
      <el-table-column prop="userName"
                       :label="$t('expExeRecord.committer')"
                       :show-overflow-tooltip="true"
                       min-width="80" />
      <el-table-column prop="exp_exec_type"
                       :show-overflow-tooltip="true"
                       :label="$t('expExeRecord.excutionType')"
                       min-width="120" />
      <el-table-column prop="exp_exec_status"
                       :show-overflow-tooltip="true"
                       :label="$t('expExeRecord.excutionStatus')"
                       min-width="60">
        <template slot-scope="scope">
          <status-tag :status="scope.row.exp_exec_status"
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
              <el-dropdown-item :command="beforeHandleCommand('viewFlow',scope.row)">{{$t('DI.taskDetail')}}</el-dropdown-item>
              <el-dropdown-item :command="beforeHandleCommand('deleteItem',scope.row)">{{$t('common.delete')}}</el-dropdown-item>
              <!-- <el-dropdown-item :command="beforeHandleCommand('viewLogFun',scope.row)">{{$t('DI.viewlog')}}</el-dropdown-item> -->
              <el-dropdown-item v-if="scope.row.exp_exec_status==='Inited'||scope.row.exp_exec_status==='WaitForRetry'||scope.row.exp_exec_status==='Scheduled'||scope.row.exp_exec_status==='Running'"
                                :command="beforeHandleCommand('stopExperiment',scope.row)">{{$t('expExeRecord.end')}}</el-dropdown-item>
              <el-dropdown-item :command="beforeHandleCommand('toJobExeRecord',scope.row)">{{$t('expExeRecord.taskList')}}</el-dropdown-item>
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
    <el-dialog :title="$t('expExeRecord.experimentLog')"
               :visible.sync="dialogVisible"
               @close="closeLog"
               custom-class="dialog-style"
               width="900px">
      <div class="log-box"
           v-loading="loading">
        <div v-if="expLog.length>0">
          <p v-for="(item ,index) of expLog"
             :key="index">{{item}}</p>
        </div>
        <div v-else
             class="no-data">
          {{this.$t('common.noData')}}
        </div>
      </div>
    </el-dialog>
  </div>
</template>
<script>
import StatusTag from '../../components/StatusTag'
import mixin from '../../util/modelMixin'
import util from '../../util/common'
export default {
  mixins: [mixin],
  components: {
    StatusTag
  },
  data () {
    return {
      dialogVisible: false,
      dataList: [],
      fromLine: 1,
      loading: false,
      expLog: [],
      query: {
        query_str: ''
      },
      statusObj: {
        Inited: 'info',
        WaitForRetry: 'info',
        Scheduled: '',
        Running: '',
        Succeed: 'success',
        Failed: 'danger',
        Timeout: 'danger',
        Cancelled: 'warning'
      },
      flowStatus: '',
      historyExpId: this.$route.query.exp_id
    }
  },
  mounted () {
    this.getListData()
  },
  computed: {
    noDataText: function () {
      return !this.loading ? this.$t('common.noData') : ''
    },
    currentExpId () {
      // 监听$route.query.exp_id 如果id有值后面没有值证明先从实验点击执行历史 在点击的实验执行历史，应重新请求实验执行全部历史记录
      if (this.historyExpId && !this.$route.query.exp_id) {
        this.resetListData()
      }
      return this.$route.query.exp_id
    }
  },
  methods: {
    getListData () {
      let query = {
        page: this.pagination.pageNumber,
        size: this.pagination.pageSize,
        username: localStorage.getItem('userId')
      }
      if (this.query.query_str) {
        query.query_str = this.query.query_str
      }
      let url = `/di/${this.FesEnv.diApiVersion}/experimentRuns`
      if (this.$route.query.exp_id) {
        url = `/di/${this.FesEnv.diApiVersion}/experimentRunsHistory/${this.$route.query.exp_id}`
      }
      this.FesApi.fetch(url, query, 'get').then(rst => {
        rst.experiment_runs && rst.experiment_runs.forEach(item => {
          item.exp_run_create_time = util.transDate(new Date(item.exp_run_create_time).getTime())
          item.exp_run_end_time = item.exp_run_end_time === '0001-01-01T00:00:00.000Z' ? '' : util.transDate(new Date(item.exp_run_end_time).getTime())
          item.userName = item.user.Name
          item.exp_name = item.experiment.exp_name
        })
        this.dataList = Object.freeze(rst.experiment_runs)
        this.pagination.totalPage = rst.total || 0
      })
    },
    viewFlow (trData) {
      this.$router.push({
        name: 'experimentFlow',
        query: {
          readonly: true,
          exp_id: trData.exp_id,
          exp_name: trData.exp_name,
          exec_id: trData.id
        }
      })
    },
    toJobExeRecord (trData) {
      this.$router.push({
        name: 'jobExeRecord',
        query: {
          expRunId: trData.id
        }
      })
    },
    viewLogFun (trData) {
      this.expLog = []
      this.dialogVisible = true
      this.loading = true
      this.flowStatus = trData.exp_exec_status
      this.byStatusGetLog(trData)
    },
    closeLog () {
      this.expLog = []
      this.setTimeoutGetLog && clearTimeout(this.setTimeoutGetLog)
    },
    byStatusGetLog (trData) {
      if ('Inited,Scheduled,Running'.indexOf(this.flowStatus) > -1) {
        this.FesApi.fetch(`/di/${this.FesEnv.diApiVersion}/experimentRun/${trData.id}/status`, 'get').then((res) => {
          this.flowStatus = res.status
        })
        this.setTimeoutGetLog = setTimeout(() => {
          this.byStatusGetLog(trData)
        }, 5000)
      } else {
        this.setTimeoutGetLog && clearTimeout(this.setTimeoutGetLog)
      }
      this.getLog(trData)
    },
    getLog (trData) {
      this.FesApi.fetch(`/di/${this.FesEnv.diApiVersion}/experimentRun/${trData.id}/log`, {
        from_line: this.fromLine,
        size: -1
      }, 'get')
        .then((rst) => {
          if (rst.log[3]) {
            this.expLog = Object.freeze(this.expLog.concat(rst.log[3].split('\n')))
          }
          this.loading = false
        }, () => {
          this.loading = false
        })
    },
    stopExperiment (trData) {
      this.$confirm(this.$t('expExeRecord.endExperimentTip'), this.$t('common.prompt')).then((index) => {
        this.FesApi.fetch(`/di/${this.FesEnv.diApiVersion}/experimentRun/${trData.id}/kill`, 'get').then(() => {
          this.deleteItemSuccess()
        })
      }).catch(() => { })
    },
    deleteItem (trData) {
      const url = `/di/${this.FesEnv.diApiVersion}/experimentRun/${trData.id}`
      this.deleteListItem(url)
    },
    resetListData () {
      this.pagination.pageNumber = 1
      this.query.query_str = ''
      this.$refs.queryValidate.resetFields()
      this.getListData()
    }
  },
  beforeDestroy () {
    this.setTimeoutGetLog && clearTimeout(this.setTimeoutGetLog)
  }
}
</script>
<style lang="scss" scoped>
.log-box {
  height: 600px;
  line-height: 20px;
  overflow-y: scroll;
}
.no-data {
  margin-top: 100px 0;
  text-align: center;
  font-size: 18px;
  font-weight: 600;
}
</style>
