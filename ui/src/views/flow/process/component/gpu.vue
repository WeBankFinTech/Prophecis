<template>
  <div>
    <el-form ref="formValidate"
             :model="form"
             :rules="ruleValidate"
             label-position="top"
             class="node-parameter-bar">
      <div v-show="activeName==='0'"
           key="step0">
        <el-form-item :label="$t('flow.nodeName')"
                      prop="name">
          <el-input v-model="form.name"
                    :disabled="readonly"></el-input>
        </el-form-item>

        <el-form-item :label="$t('flow.nodeDesc')"
                      prop="description">
          <el-input v-model="form.description"
                    type="textarea"
                    :disabled="readonly"></el-input>
        </el-form-item>
      </div>
      <div v-show="activeName==='1'"
           key="step1">
        <el-form-item :label="$t('DI.imageSettings')"
                      prop="imageType">
          <el-radio-group v-model="form.imageType"
                          :disabled="readonly">
            <el-radio label="Standard">
              {{ $t('DI.standard') }}
            </el-radio>
            <el-radio label="Custom">
              {{ $t('DI.custom') }}
            </el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item v-if="form.imageType === 'Custom'"
                      key="imageInput"
                      :label="$t('DI.imageSelection')"
                      prop="imageInput">
          <span>&nbsp;{{ FesEnv.DI.defineImage }}&nbsp;&nbsp;</span>
          <el-input v-model="form.imageInput"
                    :disabled="readonly"
                    :placeholder="$t('DI.imageInputPro')" />
        </el-form-item>
        <el-form-item v-else
                      key="imageOption"
                      :label="$t('DI.imageSelection')"
                      prop="imageOption">
          <span>&nbsp;{{ FesEnv.DI.defineImage }}&nbsp;&nbsp;</span>
          <el-select v-model="form.imageOption"
                     :placeholder="$t('DI.imagePro')">
            <el-option v-for="(item,index) in imageOptionList"
                       :label="item"
                       :key="index"
                       :value="item"
                       :disabled="readonly">
            </el-option>
          </el-select>
        </el-form-item>
        <div class="subtitle">
          {{ $t('DI.computingResource') }}
        </div>
        <el-form-item :label="$t('DI.jobType')"
                      prop="job_type">
          <el-select v-model="form.job_type"
                     :placeholder="$t('DI.jobTypePro')"
                     :disabled="readonly">
            <el-option :label="$t('DI.Single')"
                       value="Local">
            </el-option>
            <el-option :label="$t('DI.distributed')"
                       value="dist-tf">
            </el-option>
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('ns.nameSpace')"
                      prop="namespace">
          <el-select v-model="form.namespace"
                     filterable
                     :disabled="readonly"
                     :placeholder="$t('ns.nameSpacePro')">
            <el-option v-for="item in spaceOptionList"
                       :label="item"
                       :key="item"
                       :value="item">
            </el-option>
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('DI.cpu')"
                      prop="cpus">
          <el-input v-model="form.cpus"
                    :disabled="readonly"
                    :placeholder="$t('DI.CPUPro')">
            <span slot="append">Core</span>
          </el-input>
        </el-form-item>
        <el-form-item :label="$t('DI.gpu')"
                      prop="gpus">
          <el-input v-model="form.gpus"
                    :disabled="readonly"
                    :placeholder="$t('DI.GPUPro')">
            <span slot="append">{{ $t('DI.block') }}</span>
          </el-input>
        </el-form-item>
        <el-form-item :label="$t('DI.memory')"
                      prop="memory">
          <el-input v-model="form.memory"
                    :disabled="readonly"
                    :placeholder="$t('DI.memoryPro')">
            <span slot="append">Gb</span>
          </el-input>
        </el-form-item>
        <div v-if="form.job_type==='dist-tf'">
          <el-form-item :label="$t('DI.learners')"
                        prop="learners">
            <el-input v-model="form.learners"
                      :disabled="readonly"
                      :placeholder="$t('DI.learnersPro')" />
          </el-form-item>
          <el-form-item :label="$t('DI.pss')"
                        prop="pss">
            <el-input v-model="form.pss"
                      :disabled="readonly"
                      :placeholder="$t('DI.pssPro')" />
          </el-form-item>
          <el-form-item :label="$t('DI.ps_cpu')"
                        prop="ps_cpu">
            <el-input v-model="form.ps_cpu"
                      :disabled="readonly"
                      :placeholder="$t('DI.ps_cpuPro')">
              <span slot="append">Core</span>
            </el-input>
          </el-form-item>
          <el-form-item :label="$t('DI.ps_memory')"
                        prop="ps_memory">
            <el-input v-model="form.ps_memory"
                      :disabled="readonly"
                      :placeholder="$t('DI.ps_memoryPro')">
              <span slot="append">Gi</span>
            </el-input>
          </el-form-item>
          <el-form-item :label="$t('DI.ps_imageType')"
                        prop="ps_imageType">
            <el-radio-group v-model="form.ps_imageType"
                            :disabled="readonly">
              <el-radio label="Standard">
                {{ $t('DI.standard') }}
              </el-radio>
              <el-radio label="Custom">
                {{ $t('DI.custom') }}
              </el-radio>
            </el-radio-group>
          </el-form-item>
          <el-form-item v-if="form.ps_imageType==='Custom'"
                        key="psImageInput"
                        :label="$t('DI.ps_image')"
                        prop="ps_imageInput">
            <span>&nbsp;{{ FesEnv.DI.defineImage }}&nbsp;&nbsp;</span>
            <el-input v-model="form.ps_imageInput"
                      :disabled="readonly"
                      :placeholder="$t('DI.imageInputPro')" />
          </el-form-item>
          <el-form-item v-else
                        key="psImageSelect"
                        :label="$t('DI.ps_image')"
                        prop="ps_image">
            <span>&nbsp;{{ FesEnv.DI.defineImage }}&nbsp;&nbsp;</span>
            <el-select v-model="form.ps_image"
                       :disabled="readonly"
                       :placeholder="$t('DI.imagePro')">
              <el-option v-for="(item,index) in imageOptionList"
                         :label="item"
                         :key="index"
                         :value="item">
              </el-option>
            </el-select>
          </el-form-item>
        </div>
        <div class="subtitle">
          {{ $t('flow.directorySet') }}
        </div>
        <el-form-item :label="$t('flow.storageRootDirectory')"
                      prop="path">
          <el-select v-model="form.path"
                     filterable
                     :disabled="readonly"
                     :placeholder="$t('DI.trainingDataStorePro')">
            <el-option v-for="(item,index) in localPathList"
                       :label="item"
                       :key="index"
                       :value="item">
            </el-option>
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('flow.dataSubdirectory')"
                      prop="trainingData">
          <el-input v-model="form.trainingData"
                    :disabled="readonly"
                    :placeholder="$t('DI.dataSubdirectoryPro')" />
        </el-form-item>
        <el-form-item :label="$t('flow.resultSubdirectory')"
                      prop="trainingResults">
          <el-input v-model="form.trainingResults"
                    :disabled="readonly"
                    :placeholder="$t('DI.resultSubdirectoryPro')" />
        </el-form-item>
      </div>
      <div v-show="activeName==='2'"
           key="step2">
        <div class="subtitle">
          {{ $t('DI.jobExecution') }}
        </div>
        <el-form-item :label="$t('DI.entrance')"
                      prop="command">
          <el-input v-model="form.command"
                    :rows="4"
                    type="textarea"
                    :disabled="readonly"
                    :placeholder="$t('DI.entrancePro')" />
        </el-form-item>
        <el-form-item :label="$t('DI.executeCodeSettings')"
                      prop="codeSettings">
          <el-select v-model="form.codeSettings"
                     :disabled="readonly"
                     :placeholder="$t('DI.executeCodeSettingsPro')">
            <el-option :label="$t('DI.manualUpload')"
                       value="codeFile">
            </el-option>
            <el-option :label="$t('DI.shareDirectory')"
                       value="storagePath">
            </el-option>
          </el-select>
        </el-form-item>
        <div v-if="form.codeSettings==='codeFile'"
             key="codeFile"
             class="code-upload">
          <upload :marginLeft="0"
                  :baseUrl="baseUrl"
                  :dss-header="dssHeader"
                  :readonly="readonly"
                  :code-file="fileName"
                  @deleteFile="deleteFile"
                  @uploadFile="uploadFile"></upload>
        </div>
        <el-form-item v-if="form.codeSettings==='storagePath'"
                      key="storagePath"
                      :label="$t('DI.shareDirectory')"
                      prop="diStoragePath">
          <el-input v-model="form.diStoragePath"
                    :disabled="readonly"
                    :placeholder="$t('DI.shareDirectoryPro')" />
        </el-form-item>
        <storage-path-aide v-if="form.codeSettings==='storagePath'"
                           :root-path="form.path"
                           :storage-path="form.diStoragePath"></storage-path-aide>
      </div>
    </el-form>
    <div class="save-button"
         v-if="!readonly">
      <el-button @click="save"
                 :disabled="nodeSaveBtnDisable">{{ $t('common.save') }}</el-button>
    </div>
  </div>
</template>
<script>
import modelMixin from '../../../../util/modelMixin'
import nodeFormMixin from './nodeFormMixin'
import handleDIJobMixin from '../../../../util/handleDIJobMixin'
// import { Upload } from 'element-ui'
import Upload from '../../../../components/Upload'
import StoragePathAide from '../../../../components/StoragePathAide'
const initForm = {
  name: '',
  description: '',
  imageType: 'Standard',
  imageOption: '',
  imageInput: '',
  nameSpace: '',
  cpus: '',
  gpus: '',
  memory: '',
  path: '',
  trainingData: '',
  trainingResults: '',
  command: '',
  codeSettings: 'codeFile',
  diStoragePath: '',
  job_type: 'Local',
  learners: '',
  pss: '',
  ps_cpu: '',
  ps_imageType: 'Standard',
  ps_image: '',
  ps_imageInput: '',
  ps_memory: ''
}
export default {
  mixins: [modelMixin, handleDIJobMixin, nodeFormMixin],
  components: {
    Upload,
    StoragePathAide
  },
  props: {
    nodeData: {
      type: Object,
      default () {
        return {}
      }
    },
    activeName: {
      type: String,
      default: '0'
    },
    readonly: {
      type: Boolean,
      default: false
    },
    nodeSaveBtnDisable: {
      type: Boolean,
      default: false
    },
    proxyUserOption: {
      type: Array,
      default () {
        return []
      }
    }
  },
  data () {
    return {
      currentNode: {},
      imageOptionList: this.FesEnv.DI.imageOption,
      spaceOptionList: [],
      localPathList: [],
      form: initForm,
      fileName: '',
      codeFile: '',
      alarmData: null,
      dssHeader: {},
      baseUrl: `${process.env.VUE_APP_BASE_SURL}/di/v1/codeUpload`,
      userId: localStorage.getItem('userId')
    }
  },
  computed: {
    ruleValidate: function () {
      // 切换语言表单报错重置表单
      // util.resetQueryFields.apply(this)
      return {
        name: [
          { required: true, message: this.$t('flow.nodeNoEmpty') },
          { type: 'string', pattern: new RegExp(/^[0-9a-zA-Z-_]*$/), message: this.$t('flow.nodeNameTip') }
        ],
        description: [
          { required: true, message: this.$t('DI.descriptionReq') }
        ],
        namespace: [
          { required: true, message: this.$t('ns.nameSpaceReq') }
        ],
        imageType: [
          { required: true, message: this.$t('DI.imageTypeReq') }
        ],
        imageOption: [
          { required: true, message: this.$t('DI.imageColumnReq') }
        ],
        imageInput: [
          { required: true, message: this.$t('DI.imageInputReq') },
          { pattern: new RegExp(/^[[a-zA-Z0-9][a-zA-Z0-9-._]*$/), message: this.$t('DI.imageInputFormat') }
        ],
        cpus: [
          { required: true, message: this.$t('DI.CPUReq') },
          { pattern: new RegExp(/^((0{1})([.]\d{1})|([1-9]\d*)([.]\d{1})?)$/), message: this.$t('AIDE.CPUNumberReq') }
        ],
        gpus: [
          { required: true, message: this.$t('DI.GPUReq') },
          { pattern: new RegExp(/^(([1-9]{1}\d*([.]0)*)|([0]{1}))$/), message: this.$t('DI.GPUFormat') }
        ],
        memory: [
          { required: true, message: this.$t('DI.memoryReq') },
          { pattern: new RegExp(/^[1-9]\d*([.]0)*$/), message: this.$t('AIDE.memoryNumberReq') }
        ],
        job_type: [
          { required: true, message: this.$t('DI.jobTypeReq') }
        ],
        ps_cpu: [
          { required: true, message: this.$t('DI.CPUReq') },
          { pattern: new RegExp(/^((0{1})([.]\d{1})|([1-9]\d*)([.]\d{1})?)$/), message: this.$t('AIDE.CPUNumberReq') }
        ],
        ps_memory: [
          { required: true, message: this.$t('DI.memoryReq') },
          { pattern: new RegExp(/^[1-9]\d*([.]0)*$/), message: this.$t('AIDE.memoryNumberReq') }
        ],
        ps_imageType: [
          { required: true, message: this.$t('DI.imageTypeReq') }
        ],
        ps_image: [
          { required: true, message: this.$t('DI.imageColumnReq') }
        ],
        ps_imageInput: [
          { required: true, message: this.$t('DI.imageInputReq') },
          { pattern: new RegExp(/^[[a-zA-Z0-9][a-zA-Z0-9-._]*$/), message: this.$t('DI.imageInputFormat') }
        ],
        learners: [
          { required: true, message: this.$t('DI.learnersReq') },
          { pattern: new RegExp(/^[1-9]\d*([.]0)*$/), message: this.$t('DI.learnersFormat') }
        ],
        pss: [
          { required: true, message: this.$t('DI.pssReq') },
          { pattern: new RegExp(/^[1-9]\d*([.]0)*$/), message: this.$t('DI.pssFormat') }
        ],
        trainingData: [
          { required: true, message: this.$t('DI.dataDirectoryReq') },
          { pattern: new RegExp(/^[a-zA-Z][0-9a-zA-Z-_/]*$/), message: this.$t('DI.dataDirectoryFormat') }
        ],
        trainingResults: [
          { required: true, message: this.$t('DI.resultDirectoryReq') },
          { pattern: new RegExp(/^[a-zA-Z][0-9a-zA-Z-_/]*$/), message: this.$t('DI.resultDirectoryFormat') }
        ],
        path: [
          { required: true, message: this.$t('DI.trainingDataStoreReq') },
          { pattern: new RegExp(/^\/[a-zA-Z][0-9a-zA-Z-_/]*$/), message: this.$t('DI.trainingDataStoreFormat') }
        ],
        command: [
          { required: true, message: this.$t('DI.entranceReq') },
          { pattern: new RegExp(/^(.*[^\s])/), message: this.$t('DI.entranceReq') }
        ],
        codeSettings: [
          { required: true, message: this.$t('DI.executeCodeSettingsReq') }
        ],
        diStoragePath: [
          { required: true, message: this.$t('DI.shareDirectoryReq') },
          { pattern: new RegExp(/^[a-zA-Z][0-9a-zA-Z-:._/]*$/), message: this.$t('DI.shareDirectoryFormat') }
        ]
      }
    }
  },
  created () {
    const header = this.handlerFromDssHeader('get')
    this.dssHeader = this.getSystemAuthHeader()
    this.getDISpaceOption(header)
    this.getDIStoragePath(header)
  },
  watch: {
    nodeData: {
      immediate: true,
      handler (val) {
        this.alarmData = null
        this.serviceModels = {}
        this.deleteFile()
        this.handleOriginNodeData(val, initForm)
      }
    }
  },
  methods: {
    deleteFile () {
      this.codeFile = ''
      this.fileName = ''
    },
    uploadFile (file, key) {
      this.codeFile = file.result.s3Path
      this.fileName = file.fileName
    },
    judgmentFormTab () {
      return {
        '0': ['name', 'description'],
        '1': ['imageType', 'imageOption', 'imageInput', 'namespace', 'cpus', 'gpus', 'memory', 'path', 'trainingResults', 'trainingData'],
        '2': ['command', 'codeSettings', 'diStoragePath']
      }
    },
    save (evt, temSave = false) { // temSave 为true 数据保存到工作flow 不提交到后台
      this.$refs.formValidate.validate((valid, errObj) => {
        if (valid) {
          if (this.form.codeSettings === 'codeFile' && !this.codeFile) {
            if (!temSave) {
              this.$message.error(this.$t('common.uploadReq'))
            }
            return
          }
          this.disabled = true
          const param = this.handleGPUSubParam()
          if (this.form.codeSettings === 'storagePath') {
            param.data_stores[0].training_workspace = {}
            param.data_stores[0].training_workspace.container = this.form.diStoragePath
          } else {
            param.code_path = this.codeFile
            param.fileName = this.fileName
          }
          param.job_type = this.form.job_type
          if (this.form.job_type === 'dist-tf') {
            param.learners = parseInt(this.form.learners)
            param.pss = this.form.pss
            param.ps_cpu = parseFloat(this.form.ps_cpu)
            param.ps_image = this.form.ps_imageType === 'Standard' ? this.form.ps_image : this.form.ps_imageInput
            param.ps_imageType = this.form.ps_imageType
            param.ps_memory = this.form.ps_memory + 'Gi'
          }
          this.handleCurrentNode(param)
          if (this.alarmData) {
            this.currentNode.jobContent.ManiFest.job_alert = this.alarmData
          }
          this.$emit('saveNode', this.currentNode, temSave)
        } else {
          if (temSave) {
            return
          }
          const keyArr = this.judgmentFormTab()
          this.validateFormChangeTab(keyArr, errObj, this.activeName)
          // this.$message.warning(this.$t('message.process.nodeParameter.BCSB'))
        }
      })
    },
    handleGpuNodeData (param) {
      let formData = {}
      const arr = ['name', 'description', 'gpus', 'cpus', 'namespace', 'codeSettings', 'job_type']
      for (let key of arr) {
        formData[key] = param[key]
      }
      const psKeyArr = ['learners', 'pss', 'ps_cpu', 'ps_imageType', 'ps_image']
      for (let psKey of psKeyArr) {
        if (param.job_type === 'dist-tf') {
          formData[psKey] = param[psKey]
        } else {
          formData[psKey] = initForm[psKey]
        }
      }

      formData.ps_memory = param.ps_memory && parseInt(param.ps_memory)
      formData.memory = parseInt(param.memory)
      formData.trainingData = param.data_stores[0].training_data.container
      formData.trainingResults = param.data_stores[0].training_results.container
      formData.path = param.data_stores[0].connection.path
      formData.command = param.framework.command
      let curPyth = this.imageOptionList.indexOf(param.framework.version) > -1
      if (curPyth !== undefined) {
        formData.imageType = 'Standard'
        formData.imageOption = param.framework.version
      } else {
        formData.imageType = 'Custom'
        formData.imageInput = param.framework.version
      }
      this.initCodeSettings(param, formData)
      this.form = formData
    },
    initCodeSettings (param, formData) {
      formData.codeSettings = 'storagePath'
      if (param.data_stores[0].training_workspace && param.data_stores[0].training_workspace.container) {
        formData.diStoragePath = param.data_stores[0].training_workspace.container
      } else {
        formData.codeSettings = 'codeFile'
        this.codeFile = param.code_path
        this.fileName = param.fileName
      }
    },
    changeProxySet (val) {
      if (!val) {
        this.form.path = ''
      }
    },
    changeProxyUser (val) {
      for (let item of this.proxyUserOption) {
        if (item.name === val) {
          this.form.path = item.path
        }
      }
    }
  }
}
</script>
<style lang="scss" scoped>
</style>
