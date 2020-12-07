<template>
  <div>
    <breadcrumb-nav class="margin-b20">
      <el-button type="primary"
                 icon="el-icon-plus"
                 @click="dialogVisible = true"
                 size="small">
        {{ $t('label.addlabels') }}
      </el-button>
    </breadcrumb-nav>
    <el-table :data="dataList"
              border
              :empty-text="$t('common.noData')">
      <el-table-column prop="ipAddress"
                       min-width="50%"
                       :label="$t('label.ipAddress')" />
      <el-table-column prop="betaKubernetes"
                       :label="$t('label.betaKubernetes')">
        <template slot-scope="scope">
          <div class="machine-lable">
            <div v-for="(item,key) in scope.row.labels"
                 :key="key">
              <span class="lableKey">{{key}}: </span>
              <span class="lableValue">
                {{scope.row.labels[key]}}
              </span>
            </div>
          </div>
        </template>
      </el-table-column>
      <el-table-column :label="$t('common.operation')"
                       min-width="50%">
        <template slot-scope="scope">
          <el-button @click="updateItem(scope.row)"
                     type="text"
                     size="small">{{$t('common.modify')}}</el-button>
          <el-button @click="deleteItem(scope.row)"
                     type="text"
                     size="small">{{$t('common.delete')}}</el-button>
          <el-button @click="viewDetails(scope.row)"
                     type="text"
                     size="small">{{$t('ns.viewDetails')}}</el-button>
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
    <el-dialog :title="title"
               :visible.sync="dialogVisible"
               @close="clearDialogForm"
               custom-class="dialog-style">
      <el-form ref="formValidate"
               :model="form"
               :rules="ruleValidate"
               label-width="100px">
        <el-form-item :label="$t('ns.nameSpace')"
                      prop="nameSpace">
          <el-select v-model="form.nameSpace"
                     filterable
                     :placeholder="$t('ns.nameSpacePro')">
            <el-option v-for="item in spaceOptionList"
                       :label="item.namespace"
                       :key="item.id"
                       :value="item.namespace">
            </el-option>
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('label.ipAddress')"
                      prop="ipAddress">
          <el-input v-model="form.ipAddress"
                    :placeholder="$t('label.ipAddressPro')"
                    :disabled="modify" />
        </el-form-item>
        <el-form-item :label="$t('ns.lbGpuModel')"
                      prop="lbGpuModel">
          <el-input v-model="form.lbGpuModel"
                    :placeholder="$t('ns.lbGpuModelPro')" />
        </el-form-item>
        <el-form-item :label="$t('ns.lbBusType')"
                      prop="lbBusType">
          <el-input v-model="form.lbBusType"
                    :placeholder="$t('ns.lbBusTypePro')" />
        </el-form-item>
      </el-form>
      <span slot="footer"
            class="dialog-footer">
        <el-button type="primary"
                   :disabled="btnDisabled"
                   @click="subInfo ">{{ $t("common.save") }}</el-button>
        <el-button @click="dialogVisible=false">{{ $t("common.cancel") }}</el-button>
      </span>
    </el-dialog>
    <el-dialog width="1100px"
               :visible.sync="detailVisible"
               :title="$t('label.detailsMachine')"
               custom-class="dialog-style">
      <pre><code >{{lableDetail}}</code></pre>
    </el-dialog>
  </div>
</template>
<script type="text/ecmascript-6">
export default {
  data () {
    return {
      dialogVisible: false,
      detailVisible: false,
      form: {
        nameSpace: '',
        ipAddress: '',
        lbGpuModel: '',
        lbBusType: ''
      },
      modify: false,
      dataList: [],
      lableDetail: {},
      spaceOptionList: []
    }
  },
  mounted () {
    this.getListData()
    this.getSpaceOption()
  },
  computed: {
    title () {
      return this.modify ? this.$t('label.machineLabelModify') : this.$t('label.machineLabelSet')
    },
    getUrl () {
      return `/cc/${this.FesEnv.ccApiVersion}/resources/labels`
    },
    ruleValidate () {
      // 切换语言表单报错重置表单
      this.resetFormFields()
      return {
        nameSpace: [
          { required: true, message: this.$t('ns.nameSpaceReq') }
        ],
        ipAddress: [
          { required: true, message: this.$t('label.ipAddressReq') },
          { pattern: new RegExp(/^.*[^\s][\s\S]*$/), message: this.$t('label.ipAddressReq') }
        ]
      }
    }
  },
  methods: {
    getSpaceOption () {
      this.FesApi.fetch(`/cc/${this.FesEnv.ccApiVersion}/namespaces`, 'get').then(rst => {
        this.spaceOptionList = this.formatResList(rst)
      })
    },
    updateItem (trData) {
      this.dialogVisible = true
      this.modify = true
      setTimeout(() => {
        this.form.ipAddress = trData.ipAddress
        this.form.lbGpuModel = trData.labels['lb-gpu-model']
        this.form.lbBusType = trData.labels['lb-bus-type']
      }, 0)
    },
    deleteItem (trData) {
      this.$confirm(this.$t('common.deletePro'), this.$t('common.prompt')).then((index) => {
        let keyArr = []
        for (let key in trData.labels) {
          if (key.indexOf('lb-') >= 0) {
            keyArr.push(key)
          }
        }
        if (keyArr.length > 0) {
          let url = `/cc/${this.FesEnv.ccApiVersion}/nodes/${trData.labels['kubernetes.io/hostname']}/labels/${keyArr.join(',')}`
          this.FesApi.fetch(url, 'delete').then(() => {
            this.deleteItemSuccess()
          })
        } else {
          this.$message.error(this.$t('label.nolabelDeletePro'))
        }
      }).catch(() => { })
    },
    viewDetails (trData) {
      let url = `/cc/${this.FesEnv.ccApiVersion}/resources/${trData.labels['kubernetes.io/hostname']}`
      this.FesApi.fetch(url, 'get').then(rst => {
        this.lableDetail = rst || this.$t('common.noData')
        this.detailVisible = true
      })
    },
    subInfo () {
      this.$refs.formValidate.validate((valid) => {
        if (valid) {
          this.setBtnDisabeld()
          let method = this.modify ? 'put' : 'post'
          let parms = {
            namespace: this.form.nameSpace,
            ipAddress: this.form.ipAddress,
            lbGpuModel: this.form.lbGpuModel,
            lbBusType: this.form.lbBusType
          }
          this.FesApi.fetch(`/cc/${this.FesEnv.ccApiVersion}/resources/labels`, parms, method).then(() => {
            this.getListData()
            this.toast()
            this.dialogVisible = false
          })
        }
      })
    }
  }
}
</script>
<style lang="scss" scoped>
.machine-lable {
  min-width: 500px;
  .lableKey,
  .lableValue {
    display: inline-block;
    text-align: left;
  }
  .lableKey {
    text-align: left;
    min-width: 300px;
  }
  .lableValue {
    min-width: 200px;
  }
}
</style>
