<template>
  <div>
    <breadcrumb-nav></breadcrumb-nav>
    <el-form ref="queryValidate"
             :model="query"
             class="query"
             label-width="0">
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
                     @click="resetQueryStrListData">
            {{ $t("common.reset") }}
          </el-button>
        </el-col>
      </el-row>
    </el-form>
    <el-table :data="dataList"
              border
              style="width: 100%"
              :empty-text="$t('common.noData')">
      <el-table-column prop="report_name"
                       :show-overflow-tooltip="true"
                       :label="$t('DI.name')"
                       min-width="90" />
      <el-table-column prop="user_name"
                       :show-overflow-tooltip="true"
                       :label="$t('user.user')"
                       min-width="80" />
      <el-table-column prop="group_name"
                       :show-overflow-tooltip="true"
                       :label="$t('user.userGroup')"
                       min-width="90" />
      <el-table-column prop="model_name"
                       :show-overflow-tooltip="true"
                       :label="$t('flow.factoryName')" />
      <el-table-column prop="latest_version"
                       width="105"
                       :show-overflow-tooltip="true"
                       :label="$t('report.latestVersion')" />
      <el-table-column prop="model_name"
                       :show-overflow-tooltip="true"
                       :label="$t('modelList.modelName')"
                       min-width="90" />
      <el-table-column prop="version"
                       :show-overflow-tooltip="true"
                       :label="$t('serviceList.modelVersion')"
                       min-width="90" />
      <el-table-column prop="model_type"
                       :show-overflow-tooltip="true"
                       :label="$t('serviceList.modelType')"
                       min-width="90" />
      <el-table-column prop="creation_timestamp"
                       :show-overflow-tooltip="true"
                       :label="$t('DI.createTime')"
                       min-width="90" />
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
              <el-dropdown-item :command="beforeHandleCommand('deleteItem',scope.row)">{{$t('common.delete')}}</el-dropdown-item>
              <el-dropdown-item :command="beforeHandleCommand('versionListFun',scope.row)">{{$t('modelList.versionList')}}</el-dropdown-item>
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
  </div>
</template>
<script type="text/ecmascript-6">
import util from '../../../util/common'
import modelMixin from '../../../util/modelMixin'
export default {
  mixins: [modelMixin],
  data () {
    return {
      query: {
        query_str: ''
      },
      dataList: [],
      loading: false,
      verPagination: {
        pageNumber: 1,
        pageSize: 10,
        pagination: 0
      }
    }
  },
  created () {
    this.getListData()
  },
  computed: {
    noDataText () {
      return !this.loading ? this.$t('common.noData') : ''
    }
  },
  methods: {
    resetQueryStrListData () {
      this.query.query_str = ''
      this.getListData()
    },
    getListData () {
      this.loading = true
      let url = '/mf/v1/reports'
      let params = {
        currentPage: this.pagination.pageNumber,
        pageSize: this.pagination.pageSize,
        queryStr: this.query.query_str
      }
      this.FesApi.fetch(url, params, 'get').then(res => {
        let reportsList = []
        for (let item of res.Reports) {
          let report = {
            id: item.id,
            report_name: item.latest_report_version.report_name,
            group_name: item.group.name,
            latest_version: item.latest_report_version.version,
            model_name: item.model.model_name,
            version: item.model_version.version,
            model_type: item.model.model_type,
            user_name: item.user.name,
            creation_timestamp: util.transDate(item.creation_timestamp),
            status: item.event.status
          }
          reportsList.push(report)
        }
        this.dataList = reportsList
        this.pagination.totalPage = res.total || 0
        this.loading = false
      }, () => {
        this.dataList = []
        this.loading = false
      })
    },
    deleteItem (trData) {
      const url = `/mf/v1/report/${trData.id}`
      this.deleteListItem(url)
    },
    versionListFun (trData) {
      this.$router.push(`/model/reportVersion/${trData.id}`)
    },
    verHandleSizeChange (size) {
      this.verPagination.pageSize = size
      this.getListData()
    },
    verHandleCurrentChange (current) {
      this.verPagination.pageNumber = current
      this.getListData()
    },
    initVersionPage () {
      this.verPagination.size = 10
      this.verPagination.pageNumber = 1
    }
  }
}
</script>
<style lang="scss" scoped>
.add-node-form {
  min-height: 250px;
}
</style>
