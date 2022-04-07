<template>
  <div>
    <breadcrumb-nav></breadcrumb-nav>
    <el-table :data="versionList"
              border
              style="width: 100%;margin-top:15px"
              :empty-text="$t('common.noData')">
      <el-table-column prop="report_name"
                       :show-overflow-tooltip="true"
                       :label="$t('flow.reportName')" />
      <el-table-column prop="version"
                       :show-overflow-tooltip="true"
                       :label="$t('report.reportVersion')"
                       width="110" />
      <el-table-column prop="creation_timestamp"
                       :show-overflow-tooltip="true"
                       :label="$t('DI.createTime')" />
      <el-table-column prop="source"
                       :show-overflow-tooltip="true"
                       :label="$t('report.filePath')" />
      <el-table-column prop="hash_value"
                       :show-overflow-tooltip="true"
                       :label="$t('report.FPSHashValue')" />
      <el-table-column prop="status"
                       :show-overflow-tooltip="true"
                       :label="$t('DI.status')"
                       width="300">
      </el-table-column>
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
              <el-dropdown-item :command="beforeHandleCommand('pushAgain',scope.row)">{{$t('report.PushAgain')}}</el-dropdown-item>
              <el-dropdown-item :command="beforeHandleCommand('pushList',scope.row)">{{$t('report.pushList')}}</el-dropdown-item>
            </el-dropdown-menu>
          </el-dropdown>
        </template>
      </el-table-column>
    </el-table>
    <el-dialog :title="$t('report.pushParameters')"
               :visible.sync="dialogVisible"
               @close="clearDialogForm"
               custom-class="dialog-style"
               width="600px">
      <el-form ref="formValidate"
               :model="form"
               :rules="formValidate"
               label-width="155px">
        <el-form-item :label="$t('flow.factoryName')"
                      prop="factory_name">
          <el-input v-model="form.factory_name"
                    :placeholder="$t('flow.factoryNamePro')" />
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
               type="report"></push-list>
  </div>
</template>
<script>
import util from '../../../util/common'
import pushList from '../../../components/pushList'
export default {
  components: { pushList },
  data () {
    return {
      versionList: [],
      reportId: this.$route.params.reportId,
      dialogVisible: false,
      form: {
        factory_name: ''
      },
      currentTrData: {}
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
        ]
      }
    }
  },
  methods: {
    getListData () {
      this.FesApi.fetch(`/mf/v1/reportversions/${this.reportId}`, {
      }, 'get').then((res) => {
        for (let item of res) {
          item.hash_value = item.event.hash_value
          item.status = item.event.status
          item.creation_timestamp = util.transDate(item.creation_timestamp)
        }
        this.versionList = this.formatResList(res)
      })
    },
    pushAgain (trData, addInfo = false) {
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
        if (trData.event.params) {
          param = JSON.parse(trData.event.params)
        } else {
          this.currentTrData = trData
          this.dialogVisible = true
          return
        }
      }
      this.FesApi.fetch(`/mf/v1/reportVersion/Push/${trData.id}`, param, 'post').then((res) => {
        if (this.dialogVisible) {
          this.dialogVisible = false
        }
        this.toast()
        this.getListData()
      })
    },
    pushList (trData) {
      this.$refs.pushList.pushListUrl = `/mf/v1/reportVersion/Push/${trData.id}`
      this.$refs.pushList.openDialog()
    },
    clearDialogForm () {
      this.$refs.formValidate.resetFields()
    }
  }
}

</script>
