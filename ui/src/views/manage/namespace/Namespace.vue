<template>
  <div>
    <breadcrumb-nav class="margin-b20">
      <el-button type="primary"
                 icon="el-icon-plus"
                 @click="dialogVisible = true"
                 size="small">
        {{ $t("ns.addNameSpace") }}
      </el-button>
    </breadcrumb-nav>
    <el-table :data="dataList"
              border
              :empty-text="$t('common.noData')">
      <el-table-column prop="namespace"
                       :show-overflow-tooltip="true"
                       :label="$t('ns.nameSpaceName')" />
      <el-table-column prop="isActive"
                       min-width="80px"
                       :show-overflow-tooltip="true"
                       :label="$t('user.isActive')"
                       :formatter="translationActive" />
      <el-table-column prop="remarks"
                       :show-overflow-tooltip="true"
                       :label="$t('DI.description')" />
      <el-table-column :label="$t('common.operation')"
                       width="80px">
        <template slot-scope="scope">
          <el-dropdown size="medium"
                       @command="handleCommand">
            <el-button type="text"
                       class="multi-operate-btn"
                       icon="el-icon-more"
                       round></el-button>
            <el-dropdown-menu slot="dropdown">
              <el-dropdown-item :command="beforeHandleCommand('updateItem',scope.row)">{{$t("common.modify")}}</el-dropdown-item>
              <el-dropdown-item :command="beforeHandleCommand('deleteItem',scope.row)">{{$t('common.delete')}}</el-dropdown-item>
              <el-dropdown-item :command="beforeHandleCommand('resourceQuota',scope.row)">{{$t('ns.resourceQuota')}}</el-dropdown-item>
              <el-dropdown-item :command="beforeHandleCommand('viewDetails',scope.row)">{{$t('ns.viewDetails')}}</el-dropdown-item>
              <el-dropdown-item :command="beforeHandleCommand('settingGroup',scope.row)">{{$t('user.userGroupSettings')}}</el-dropdown-item>
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
    <el-dialog :title="title"
               :visible.sync="dialogVisible"
               @close="clearDialogForm"
               custom-class="dialog-style"
               width="900px">
      <el-form ref="formValidate"
               :model="form"
               :rules="ruleValidate"
               label-width="160px">
        <el-form-item :label="$t('ns.nameSpaceName')"
                      prop="namespace">
          <el-input v-model="form.namespace"
                    :placeholder="$t('ns.nsNamePro')"
                    :disabled="modify" />
        </el-form-item>
        <el-form-item :label="$t('ns.platformNameSpace')">
          <el-input v-model="platformNamespace"
                    :placeholder="$t('ns.platformNsPro')"
                    disabled />
        </el-form-item>
        <div class="subtitle">
          {{ $t('ns.otherAnnotations') }}:
        </div>
        <el-row>
          <el-col :span="12">
            <el-form-item :label="$t('ns.lbGpuModel')"
                          prop="lbGpuModel">
              <el-input v-model="form.lbGpuModel"
                        :placeholder="$t('ns.lbGpuModelPro')" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item :label="$t('ns.lbBusType')"
                          prop="lbBusType">
              <el-input v-model="form.lbBusType"
                        :placeholder="$t('ns.lbBusTypePro')" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item :label="$t('user.isActive')"
                      prop="isActive">
          <el-select v-model="form.isActive"
                     :placeholder="$t('user.activePro')">
            <el-option :label="$t('common.yes')"
                       :value="1">
            </el-option>
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('ns.nsDescribe')"
                      prop="remarks">
          <el-input v-model="form.remarks"
                    type="textarea"
                    :placeholder="$t('ns.nsDescribePro')" />
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
    <el-dialog :title="$t('ns.resourceQuota')"
               :visible.sync="resourceVisible"
               @close="clearResource"
               custom-class="dialog-style"
               width="650px">
      <el-form ref="resourceForm"
               :model="resource"
               :rules="ruleValidate"
               label-width="140px">
        <el-form-item :label="$t('ns.nameSpaceName')"
                      prop="namespace">
          <el-input v-model="resource.namespace"
                    :placeholder="$t('ns.nsNamePro')"
                    disabled />
        </el-form-item>
        <el-form-item :label="$t('AIDE.cpu')"
                      prop="cpu">
          <el-input v-model="resource.cpu"
                    :placeholder="$t('ns.CPUPro')">
            <span slot="append">Core</span>
          </el-input>
        </el-form-item>
        <el-form-item :label="$t('AIDE.gpu')"
                      prop="gpu">
          <el-input v-model="resource.gpu"
                    :placeholder="$t('ns.GPUPro')">
            <span slot="append">{{ $t('DI.block') }}</span>
          </el-input>
        </el-form-item>
        <el-form-item :label="$t('AIDE.memory')"
                      prop="memoryAmount">
          <el-input v-model="resource.memoryAmount"
                    :placeholder="$t('ns.memoryPro')">
            <span slot="append">Mi</span>
          </el-input>
        </el-form-item>
      </el-form>
      <span slot="footer"
            class="dialog-footer">
        <el-button type="primary"
                   :disabled="btnDisabled"
                   @click="subResourceQuota">{{ $t("common.save") }}</el-button>
        <el-button @click="resourceVisible=false">{{ $t("common.cancel") }}</el-button>
      </span>
    </el-dialog>
    <el-dialog :title="$t('ns.nsInstance')"
               :visible.sync="instanceVisible"
               custom-class="instance-dialog"
               width="800px">
      <div v-if="instanceName"
           class="instance">
        <span class="title">
          {{ $t('ns.nameSpaceName') }}:
        </span>
        <span>
          {{ instanceName }}
        </span>
      </div>
      <div v-if="instanceTime"
           class="instance">
        <span class="title">
          {{ $t('DI.createTime') }}:
        </span>
        <span>
          {{ instanceTime }}
        </span>
      </div>
      <el-table :data="instanceList"
                border
                :empty-text="$t('common.noData')">
        <el-table-column prop="nameSpace"
                         :show-overflow-tooltip="true"
                         :label="$t('ns.attributesName')" />
        <el-table-column prop="additionalProp"
                         :show-overflow-tooltip="true"
                         :label="$t('ns.attributesValue')" />
      </el-table>
    </el-dialog>
  </div>
</template>
<script type="text/ecmascript-6">
import util from '../../../util/common'
export default {
  data () {
    return {
      dialogVisible: false,
      resourceVisible: false,
      instanceVisible: false,
      form: {
        id: '',
        namespace: '',
        lbGpuModel: '',
        lbBusType: '',
        remarks: '',
        isActive: 1,
        enableFlag: 1
      },
      resource: {
        namespace: '',
        cpu: '',
        gpu: '',
        memoryAmount: ''
      },
      modify: false,
      dataList: [],
      instanceList: [],
      instanceName: '',
      instanceTime: '',
      platformNamespace: ''
    }
  },
  mounted () {
    this.getListData()
    this.platformNamespace = this.FesEnv.ns.platformNamespace ? this.FesEnv.ns.platformNamespace : 'default'
  },
  computed: {
    title () {
      return this.modify ? this.$t('ns.modifyNs') : this.$t('ns.newNs')
    },
    getUrl () {
      return `/cc/${this.FesEnv.ccApiVersion}/namespaces`
    },
    ruleValidate () {
      // 切换语言表单报错重置表单
      this.resetFormFields()
      this.resetResourceFields()
      return {
        namespace: [
          { required: true, message: this.$t('ns.nsNameReq') },
          { pattern: new RegExp(/^.*[^\s][\s\S]*$/), message: this.$t('ns.nsNameReq') }
        ],
        isActive: [
          { required: true, message: this.$t('user.activeReq') }
        ],
        cpu: [
          { required: true, message: this.$t('DI.CPUReq') },
          { pattern: new RegExp(/^((0{1})([.]\d{1})|([1-9]\d*)([.]\d{1})?)|([0]{1})$/), message: this.$t('ns.CPUFormat') }
        ],
        gpu: [
          { required: true, message: this.$t('ns.GPUReq') },
          { pattern: new RegExp(/^(([1-9]{1}\d*([.]0)*)|([0]{1}))$/), message: this.$t('ns.GPUFormat') }
        ],
        memoryAmount: [
          { required: true, message: this.$t('DI.memoryReq') },
          { pattern: new RegExp(/^(([1-9]{1}\d*([.]0)*)|([0]{1}))*$/), message: this.$t('ns.memoryFormat') }
        ]
      }
    }
  },
  methods: {
    resetResourceFields () {
      setTimeout(() => {
        this.$refs.resourceForm && this.$refs.resourceForm.resetFields()
      }, 10)
    },
    translationActive (value, row) {
      return this.$t('common.yes')
    },
    updateItem (trData) {
      this.dialogVisible = true
      this.modify = true
      setTimeout(() => {
        Object.assign(this.form, trData)
        this.form.lbGpuModel = (trData.annotations === undefined || trData.annotations['lb-gpu-model'] === undefined) ? '' : trData.annotations['lb-gpu-model']
        this.form.lbBusType = (trData.annotations === undefined || trData.annotations['lb-bus-type'] === undefined) ? '' : trData.annotations['lb-bus-type']
      }, 0)
    },
    deleteItem (trData) {
      const url = `/cc/${this.FesEnv.ccApiVersion}/namespaces/namespace/${trData.namespace}`
      this.deleteListItem(url)
    },
    resourceQuota (trData) {
      this.resourceVisible = true
      this.$nextTick(() => {
        this.resource.namespace = trData.namespace
        if (trData.hard) {
          this.resource.cpu = trData.hard['limits.cpu'] + ''
          this.resource.gpu = trData.hard['requests.nvidia.com/gpu'] + ''
          const memory = trData.hard['limits.memory']
          this.resource.memoryAmount = memory.indexOf('Gi') > -1 ? parseInt(trData.hard['limits.memory']) * 1024 + '' : parseInt(trData.hard['limits.memory']) + ''
        }
      })
    },
    viewDetails (trData) {
      let url = `/cc/${this.FesEnv.ccApiVersion}/namespaces/namespace/${trData.namespace}`
      this.FesApi.fetch(url, 'get').then((rst) => {
        if (rst && JSON.stringify(rst) !== '{}') {
          let annotations = rst.metadata.annotations
          let metadata = rst.metadata
          let list = []
          if (metadata && annotations && JSON.stringify(annotations) !== '{}') {
            this.instanceName = rst.metadata.namespace
            this.instanceTime = util.transDate(metadata.creationTimestamp.millis)
            for (let key in annotations) {
              if (annotations[key]) {
                let item = {
                  'nameSpace': key,
                  'additionalProp': annotations[key]
                }
                list.push(item)
              }
            }
          }
          this.instanceList = list
        } else {
          this.instanceName = ''
          this.instanceTime = ''
          this.instanceList = []
        }
        this.instanceVisible = true
      })
    },
    settingGroup (trData) {
      this.$router.push({
        name: 'settingNamespaceGroup',
        params: {
          namespace: trData.namespace
        },
        query: {
          nid: trData.id
        }
      })
    },
    subResourceQuota () {
      this.$refs.resourceForm.validate((valid) => {
        if (valid) {
          this.setBtnDisabeld()
          let param = {
            cpu: this.resource.cpu + '',
            gpu: this.resource.gpu + '',
            namespace: this.resource.namespace,
            memory: {
              memoryAmount: this.resource.memoryAmount + '',
              memoryUnit: 'Mi'
            }
          }
          this.FesApi.fetchUT(`/cc/${this.FesEnv.ccApiVersion}/resources`, param, 'put').then(() => {
            this.getListData()
            this.toast()
            this.resourceVisible = false
          })
        }
      })
    },
    subInfo () {
      this.$refs.formValidate.validate((valid) => {
        if (valid) {
          this.setBtnDisabeld()
          let param = {}
          const arry = ['id', 'namespace', 'remarks']
          for (let key of arry) {
            param[key] = this.form[key]
          }
          param.platformNamespace = this.platformNamespace
          param.annotations = {}
          param.annotations['lb-gpu-model'] = this.form.lbGpuModel && this.form.lbGpuModel.trim() ? this.form.lbGpuModel : null
          param.annotations['lb-bus-type'] = this.form.lbBusType && this.form.lbBusType.trim() ? this.form.lbBusType : null
          let method = this.modify ? 'put' : 'post'
          !this.modify && delete param.id
          this.FesApi.fetch(`/cc/${this.FesEnv.ccApiVersion}/namespaces`, param, method).then(() => {
            this.getListData()
            this.toast()
            this.dialogVisible = false
          })
        }
      })
    },
    clearResource () {
      this.$refs.resourceForm.resetFields()
    }
  }
}
</script>
<style lang="scss" >
.instance-dialog {
  .el-dialog__header {
    border-bottom: 1px solid #e3e8ee;
    .el-dialog__title {
      font-size: 16px;
      font-weight: 600;
    }
  }
  .el-dialog__body {
    padding: 10px 20px 25px;
    .instance {
      padding: 5px 0;
      .title {
        font-weight: bold;
      }
    }
  }
}
</style>
