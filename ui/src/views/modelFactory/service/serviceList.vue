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
                     @click="resetQueryStrListData">
            {{ $t("common.reset") }}
          </el-button>
          <el-button type="primary"
                     icon="plus-circle-o"
                     @click="openModel">
            {{ $t("serviceList.createModel") }}
          </el-button>
        </el-col>
      </el-row>
    </el-form>
    <el-table :data="dataList"
              border
              style="width: 100%"
              :empty-text="$t('common.noData')">
      <el-table-column prop="service_name"
                       :show-overflow-tooltip="true"
                       :label="$t('serviceList.modelName')" />
      <el-table-column prop="model_name"
                       :show-overflow-tooltip="true"
                       :label="$t('modelList.modelName')" />
      <el-table-column prop="version"
                       :show-overflow-tooltip="true"
                       :label="$t('serviceList.modelVersion')" />
      <el-table-column prop="cpu"
                       :show-overflow-tooltip="true"
                       label="CPU" />
      <el-table-column prop="gpu"
                       :show-overflow-tooltip="true"
                       label="GPU" />
      <el-table-column prop="memory1"
                       :show-overflow-tooltip="true"
                       label="Memory" />
      <el-table-column prop="user_name"
                       :show-overflow-tooltip="true"
                       :label="$t('user.user')" />
      <el-table-column prop="status"
                       :show-overflow-tooltip="true"
                       :label="$t('DI.status')" />
      <el-table-column prop="last_updated_timestamp"
                       :show-overflow-tooltip="true"
                       :label="$t('serviceList.updateTime')" />
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
              <el-dropdown-item v-if="scope.row.status==='Stop'"
                                :command="beforeHandleCommand('startRun',scope.row)">{{$t('serviceList.run')}}</el-dropdown-item>
              <el-dropdown-item v-else
                                :command="beforeHandleCommand('stoptRun',scope.row)">{{$t('serviceList.stop')}}</el-dropdown-item>
              <el-dropdown-item :command="beforeHandleCommand('deleteItem',scope.row)">{{$t('common.delete')}}</el-dropdown-item>
              <el-dropdown-item :command="beforeHandleCommand('detailModel',scope.row)">{{$t('ns.viewDetails')}}</el-dropdown-item>
              <el-dropdown-item v-if="scope.row.status!=='Stop'"
                                :command="beforeHandleCommand('updateModelService',scope.row)">{{$t('serviceList.updateModelService')}}</el-dropdown-item>
              <el-dropdown-item v-if="scope.row.status!=='Stop'"
                                :command="beforeHandleCommand('containerList',scope.row)">{{$t('serviceList.containerList')}}</el-dropdown-item>
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
    <create-service ref="modelDialog"
                    :spaceOptionList="spaceOptionList"
                    :groupList="groupList"
                    @getListData="getListData">
    </create-service>
    <view-service ref="detailModel"
                  :form="modelItem"
                  :groupList="groupList"
                  :spaceOptionList="spaceOptionList">
    </view-service>
    <update-service ref="updateModelg"
                    :groupList="groupList"
                    :spaceOptionList="spaceOptionList"
                    @getListData="getListData">
    </update-service>
  </div>
</template>
<script>
import CreateService from './CreateService'
import ViewService from './ViewService'
import UpdateService from './UpdateService'
import modelMixin from '../../../util/modelMixin'
import util from '../../../util/common'
export default {
  mixins: [modelMixin],
  components: {
    CreateService,
    ViewService,
    UpdateService
  },
  data () {
    return {
      query: {
        query_str: ''
      },
      spaceOptionList: [],
      dataList: [],
      groupList: [],
      modify: false,
      setIntervalList: null,
      modelItem: {}
    }
  },
  mounted () {
    this.getListData()
    this.getDISpaceOption()
    this.getGroupList()
    this.setIntervalList = setInterval(this.getListData, 60000)
  },
  methods: {
    // 打开创建模型
    openModel () {
      this.$refs.modelDialog.dialogVisible = true
    },
    getListData () {
      let param = {
        page: this.pagination.pageNumber,
        size: this.pagination.pageSize,
        query_str: this.query.query_str
      }
      let url = '/mf/v1/services'
      this.FesApi.fetch(url, param, 'get').then(rst => {
        for (let item of rst.models) {
          item.last_updated_timestamp = util.transDate(item.last_updated_timestamp)
          item.version = item.modelversion.version
          item.memory1 = item.memory + 'Gi'
          item.gpu = item.gpu ? item.gpu : 0
        }
        this.dataList = this.formatResList(rst)
        this.pagination.totalPage = rst.total || this.dataList.length
      })
    },
    deleteItem (trData) {
      const url = `/mf/v1/service/${trData.service_id}`
      this.deleteListItem(url)
    },
    startRun (trData) {
      let param = {}
      const keyArr = ['service_id', 'service_name', 'type', 'namespace', 'creation_timestamp', 'last_updated_timestamp', 'cpu', 'memory', 'group_id', 'log_path', 'remark', 'image']
      for (let key of keyArr) {
        param[key] = trData[key]
      }
      param.Gpu = trData.gpu + ''
      param.service_post_models = [{
        modelversion_id: trData.modelversion_id,
        model_name: trData.model_name,
        model_type: trData.model_type,
        group_name: trData.group_name,
        source: trData.modelversion.source,
        model_version: trData.modelversion.version
      }]
      if (trData.model_type === 'TENSORFLOW') {
        param.service_post_models[0].parameters = JSON.parse(trData.modelversion.params)
      }
      this.FesApi.fetch(`/mf/v1/serviceRun`, param, 'post').then(rst => {
        setTimeout(() => {
          this.toast()
          this.getListData()
        }, 500)
      })
    },
    stoptRun (trData) {
      this.FesApi.fetch(`/mf/v1/serviceStop/${trData.namespace}/${trData.service_name}/${trData.service_id}`, 'get').then(rst => {
        setTimeout(() => {
          this.toast()
          this.getListData()
        }, 500)
      })
    },
    getFormKey (param, trData) {
      param.group_id = trData.group_id
      param.model_id = trData.modelversion.model_id + ''
      param.modelversion = trData.modelversion.version
      return ['cpu', 'gpu', 'memory', 'modelversion_id', 'service_id', 'namespace', 'service_name', 'remark', 'model_type', 'endpoint_type', 'group_name', 'model_name', 'image']
    },
    updateModelService (trData) {
      const updateModelg = this.$refs.updateModelg
      let param = {}
      const keyArr = this.getFormKey(param, trData)
      for (let item of keyArr) {
        param[item] = trData[item]
      }
      param.source = trData.modelversion.source
      updateModelg.form = param
      updateModelg.changeVersionId(trData.modelversion_id)
      updateModelg.dialogVisible = true
    },
    detailModel (trData) {
      let param = {}
      const keyArr = this.getFormKey(param, trData)
      for (let item of keyArr) {
        param[item] = trData[item]
      }
      param.model_id = trData.modelversion.model_id
      param.modelversion = trData.modelversion.version
      if (trData.model_type === 'TENSORFLOW') {
        param.parameters = trData.modelversion.params && JSON.parse(trData.modelversion.params)
      }
      // param.log_path = trData.log_path
      this.modelItem = { ...param }
      this.$refs.detailModel.dialogVisible = true
    },
    containerList (trData) {
      this.$router.push({
        path: `/model/serviceList/${trData.namespace}/${trData.service_name}`
      })
    }
    // fileDownload (trData) {
    //   this.FesApi.fetchUT(`/tds/v1/s3/download`, { filePath: trData.filepath }, {
    //     method: 'get',
    //     responseType: 'blob'
    //   }).then(rst => {
    //     // const source = new Blob([rst.data])
    //     let url = window.URL.createObjectURL(rst)
    //     let link = document.createElement('a')
    //     link.style.display = 'none'
    //     link.href = url
    //     link.setAttribute('download', 'modelFile.zip')
    //     document.body.appendChild(link)
    //     link.click()
    //   }, () => {
    //     this.$message.error(this.$t('serviceList.downloadFailed'))
    //   })
    // }
  },
  beforeDestroy () {
    this.setIntervalList = null
  }
}
</script>
