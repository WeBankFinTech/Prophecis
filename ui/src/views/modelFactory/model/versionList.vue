<template>
  <div>
    <breadcrumb-nav></breadcrumb-nav>
    <el-table :data="versionList"
              border
              style="width: 100%;margin-top:15px"
              :empty-text="$t('common.noData')">
      <el-table-column prop="version"
                       :show-overflow-tooltip="true"
                       :label="$t('serviceList.modelVersion')"
                       width="105" />
      <el-table-column prop="model_name"
                       :show-overflow-tooltip="true"
                       :label="$t('modelList.modelName')"
                       width="250" />
      <el-table-column prop="model_type"
                       :show-overflow-tooltip="true"
                       width="105"
                       :label="$t('modelList.type')" />
      <el-table-column prop="username"
                       :show-overflow-tooltip="true"
                       width="105"
                       :label="$t('user.user')" />
      <el-table-column prop="groupName"
                       :show-overflow-tooltip="true"
                       :label="$t('user.userGroup')"
                       width="250" />
      <el-table-column prop="creation_timestamp"
                       :show-overflow-tooltip="true"
                       width="135"
                       :label="$t('DI.createTime')" />
      <el-table-column prop="source"
                       :show-overflow-tooltip="true"
                       :label="$t('modelList.filePath')">
      </el-table-column>
      <!-- <el-table-column prop="source"
                       :show-overflow-tooltip="true"
                       label="结果信息">
      </el-table-column> -->
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
              <el-dropdown-item :command="beforeHandleCommand('fileDownload',scope.row)">{{$t('common.export')}}</el-dropdown-item>
              <!-- <el-dropdown-item v-if="scope.row.model_type==='MLPipeline'"
                                :command="beforeHandleCommand('pushAgain',scope.row)">{{$t('report.PushAgain')}}</el-dropdown-item>
              <el-dropdown-item v-if="scope.row.model_type==='MLPipeline'"
                                :command="beforeHandleCommand('pushList',scope.row)">{{$t('report.pushList')}}</el-dropdown-item> -->
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
    <el-dialog :title="$t('report.pushParameters')"
               :visible.sync="dialogVisible"
               @close="clearDialogForm"
               custom-class="dialog-style"
               width="650px">
      <el-form ref="formValidate"
               :model="form"
               :rules="formValidate"
               label-width="155px">
        <el-form-item :label="$t('flow.factoryName')"
                      prop="factory_name">
          <el-input v-model="form.factory_name"
                    :placeholder="$t('flow.factoryNamePro')" />
        </el-form-item>
        <el-form-item :label="$t('serviceList.modelType')"
                      prop="model_type">
          <el-select v-model="form.model_type"
                     :placeholder="$t('serviceList.modelTypePro')">
            <el-option label="Logistic_Regression"
                       value="Logistic_Regression"></el-option>
            <el-option label="Decision_Tree"
                       value="Decision_Tree"></el-option>
            <el-option label="Random_Forest"
                       value="Random_Forest"></el-option>
            <el-option label="XGBoost"
                       value="XGBoost"></el-option>
            <el-option label="LightGBM"
                       value="LightGBM"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('modelList.algorithmType')"
                      prop="model_usage">
          <el-select v-model="form.model_usage"
                     :placeholder="$t('modelList.algorithmTypePro')">
            <el-option label="Classification"
                       value="Classification"></el-option>
            <el-option label="Regression"
                       value="Regression"></el-option>
          </el-select>
        </el-form-item>
      </el-form>
      <span slot="footer"
            class="dialog-footer">
        <el-button type="primary"
                   @click="pushAgain(currentTrData,true)">
          {{ $t('common.save') }}
        </el-button>
        <el-button @click="dialogVisible=false">
          {{ $t('common.cancel') }}
        </el-button>
      </span>
    </el-dialog>
    <push-list ref="pushList"
               type="model"
               :url="pushListUrl"></push-list>
  </div>
</template>
<script>
import pushList from '../../../components/pushList'
import axios from 'axios'
export default {
  components: { pushList },
  data () {
    return {
      versionList: [],
      pushListUrl: '',
      dialogVisible: false,
      currentTrData: {},
      form: {
        factory_name: '',
        model_type: '',
        model_usage: ''
      },
      modelId: this.$route.params.modelId
    }
  },
  created () {
    this.getListData()
  },
  computed: {
    formValidate () {
      return {
        factory_name: [
          { required: true, message: this.$t('flow.factoryNameReq') },
          { type: 'string', pattern: new RegExp(/^[\u4e00-\u9fa5_0-9a-zA-Z-]*$/), message: this.$t('flow.factoryNameFormat') }
        ],
        model_type: [
          { required: true, message: this.$t('serviceList.modelTypeReq') }
        ],
        model_usage: [
          { required: true, message: this.$t('modelList.algorithmTypeReq') }
        ]
      }
    }
  },
  methods: {
    getListData () {
      this.FesApi.fetch(`/mf/v1/modelversions/${this.modelId}`, {
        page: this.pagination.pageNumber,
        size: this.pagination.pageSize
      }, 'get').then((res) => {
        for (let item of res.modelVersions) {
          item.model_name = item.model.model_name
          item.model_type = item.model.model_type
          item.username = item.user.name
          item.groupName = item.group.name
        }
        this.versionList = this.formatResList(res.modelVersions)
        this.pagination.totalPage = res.total || 0
      })
    },
    fileDownload (trData) {
      // responseType: 'blob'
      axios({
        url: `${process.env.VUE_APP_BASE_SURL}/mf/v1/modelVersionDownload/${trData.id}`,
        method: 'get',
        headers: {
          'Mlss-Userid': localStorage.getItem('userId')
        },
        responseType: 'arraybuffer'
      }).then(res => {
        let suffix = ''
        let downloadFileName = ''
        if (trData.model_type === 'TENSORFLOW' || trData.model_type === 'CUSTOM') {
          suffix = 'zip'
          downloadFileName = 'model_name'
        } else {
          let fileName = trData.file_name
          suffix = fileName.substring(fileName.length - fileName.lastIndexOf('.') - 1)
          downloadFileName = trData.file_name
        }
        console.log('suffix', suffix)
        const blob = new Blob([res.data], { type: `application/${suffix}` })
        const url = URL.createObjectURL(blob)
        // const blob = new Blob([res.data], { type: 'application/zip' })
        let link = document.createElement('a')
        // link.setAttribute('href',
        //   'data:application/zip;base64,' + rst.data
        // )
        link.href = url
        link.style.display = 'none'
        link.setAttribute('download', downloadFileName)
        document.body.appendChild(link)
        link.click()
      }, () => {
        this.$message.error(this.$t('serviceList.downloadFailed'))
      })
    },
    async pushAgain (trData, addInfo) {
      // addInfo为true手动填写工厂名称
      let param = {}
      if (addInfo) {
        let validate = false
        this.$refs.formValidate.validate(valid => {
          validate = valid
        })
        if (validate) {
          param = { ...this.form }
        } else {
          return
        }
      } else {
        const pushInfo = await this.FesApi.fetch(`/mf/v1/modelVersion/push/${trData.id}`, 'get')
        if (!pushInfo || (pushInfo && !pushInfo.events) || (pushInfo && pushInfo.events.length === 0) || (pushInfo && pushInfo.events && !pushInfo.events[0].event.params)) {
          this.currentTrData = trData
          this.dialogVisible = true
          return
        } else {
          let index = 0
          // 获取最近的一个event。param
          for (let i = 0; i < pushInfo.events.length - 1; i++) {
            let current = new Date(pushInfo.events[i].event.creation_timestamp).getTime()
            let latest = new Date(pushInfo.events[index].event.creation_timestamp).getTime()
            if (current > latest && pushInfo.events[index].event.params) {
              index = i
            }
          }
          param = JSON.parse(pushInfo.events[index].event.params)
        }
      }
      // const params = JSON.parse(pushInfo.events[index].event.params)
      this.FesApi.fetch(`/mf/v1/modelVersion/push/${trData.id}`, param, { method: 'post', timeout: 25000 }).then((res) => {
        if (this.dialogVisible) {
          this.dialogVisible = false
        }
        this.toast()
      })
    },
    pushList (trData) {
      this.$refs.pushList.pushListUrl = `/mf/v1/modelVersion/push/${trData.id}`
      this.$refs.pushList.openDialog()
    }
  }
}
</script>
