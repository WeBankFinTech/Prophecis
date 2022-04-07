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
            {{ $t("modelList.createModel") }}
          </el-button>
        </el-col>
      </el-row>
    </el-form>
    <el-table :data="dataList"
              border
              style="width: 100%"
              :empty-text="$t('common.noData')">
      <el-table-column prop="model_name"
                       :show-overflow-tooltip="true"
                       :label="$t('modelList.modelName')" />
      <el-table-column prop="version"
                       :show-overflow-tooltip="true"
                       :label="$t('serviceList.modelVersion')" />
      <el-table-column prop="model_type"
                       :show-overflow-tooltip="true"
                       :label="$t('modelList.type')" />
      <el-table-column prop="username"
                       :show-overflow-tooltip="true"
                       :label="$t('user.user')" />
      <el-table-column prop="groupName"
                       :show-overflow-tooltip="true"
                       :label="$t('user.userGroup')" />
      <el-table-column prop="update_timestamp"
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
              <el-dropdown-item v-if="scope.row.model_type!=='MLPipeline'"
                                :command="beforeHandleCommand('updateModel',scope.row)">{{$t('common.modify')}}</el-dropdown-item>
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
    <create-model ref="createModelg"
                  :localPathList="localPathList"
                  :groupList="groupList"
                  :modelId="modelId"
                  @getListData="getListData">
    </create-model>
  </div>
</template>
<script>
import CreateModel from './CreateModel'
// import ViewModel from './ViewModel'
// import UpdateModel from './Update'
import modelMixin from '../../../util/modelMixin'
import util from '../../../util/common'
export default {
  mixins: [modelMixin],
  components: {
    CreateModel
  },
  data () {
    return {
      query: {
        query_str: ''
      },
      userOptionList: [],
      localPathList: [],
      groupList: [],
      dataList: [],
      modelId: 0,
      verPagination: {
        pageNumber: 1,
        pageSize: 10,
        pagination: 0
      }
    }
  },
  mounted () {
    this.getListData()
    this.getDIUserOption()
    this.getDIStoragePath()
    this.getGroupList()
  },
  methods: {
    // 打开创建模型
    openModel () {
      this.$refs.createModelg.dialogVisible = true
    },
    getListData () {
      let param = {
        page: this.pagination.pageNumber,
        size: this.pagination.pageSize,
        query_str: this.query.query_str
      }
      this.FesApi.fetch('/mf/v1/models', param, 'get').then(rst => {
        for (let item of rst.models) {
          item.update_timestamp = util.transDate(item.update_timestamp)
          item.username = item.user.name
          item.groupName = item.group.name
          item.groupId = item.group.id
          item.version = item.model_latest_version.version
        }
        this.dataList = this.formatResList(rst)
        this.pagination.totalPage = rst.total || 0
      })
    },
    deleteItem (trData) {
      const url = `/mf/v1/model/${trData.id}`
      this.deleteListItem(url)
    },
    updateModel (trData) {
      const createModelg = this.$refs.createModelg
      createModelg.modify = true
      this.modelId = trData.id
      let param = {}
      for (let key in createModelg.form) {
        param[key] = trData[key]
      }
      const lastModel = trData.model_latest_version
      if (lastModel.filepath) {
        param.modelLocation = 'HDFSPath'
        for (let item of this.localPathList) {
          if (lastModel.filepath.indexOf(item) > -1) {
            param.root_path = item
            param.child_path = lastModel.filepath.substring(item.length)
            break
          }
        }
      } else {
        createModelg.params.modelType = trData.model_type
        createModelg.source = lastModel.source
        createModelg.file_name = lastModel.file_name
        param.modelLocation = 'codeFile'
      }
      createModelg.dialogVisible = true
      setTimeout(() => {
        createModelg.form = param
      })
    },
    versionListFun (trData) {
      this.$router.push(`/model/versionList/${trData.id}`)
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
