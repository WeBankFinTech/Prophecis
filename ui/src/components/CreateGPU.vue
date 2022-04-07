<template>
  <el-dialog :title="this.$t('DI.createTraining')"
             :visible.sync="dialogVisible"
             @close="clearDialogForm"
             custom-class="dialog-style"
             :close-on-click-modal="false"
             width="1240px">
    <step :step="currentStep"
          :step-list="stepList" />
    <el-form ref="formValidate"
             :model="form"
             :rules="ruleValidate"
             class="add-form"
             label-width="190px">
      <div v-if="currentStep===0"
           key="step0">
        <div class="subtitle">
          {{ $t('DI.basicSettings') }}
        </div>
        <el-row>
          <el-col :span="24">
            <el-form-item :label="$t('DI.trainingjobName')"
                          prop="name">
              <el-input v-model="form.name"
                        maxlength="225"
                        :placeholder="$t('DI.trainingjobNamePro')" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row>
          <el-form-item :label="$t('DI.description')"
                        prop="description">
            <el-input v-model="form.description"
                      type="textarea"
                      :placeholder="$t('DI.jobDescriptionPro')" />
          </el-form-item>
        </el-row>
      </div>
      <div v-if="currentStep===1"
           key="step1">
        <div class="subtitle">
          {{ $t('DI.imageSettings') }}
        </div>
        <el-row>
          <el-col :span="12">
            <el-form-item :label="$t('DI.imageType')"
                          prop="imageType">
              <el-radio-group v-model="form.imageType">
                <el-radio label="Standard">
                  {{ $t('DI.standard') }}
                </el-radio>
                <el-radio label="Custom">
                  {{ $t('DI.custom') }}
                </el-radio>
              </el-radio-group>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row>
          <el-form-item v-if="form.imageType === 'Custom'"
                        key="imageInput"
                        :label="$t('DI.imageSelection')"
                        prop="imageInput">
            <span>&nbsp;{{ defineImage }}&nbsp;&nbsp;</span>
            <el-input v-model="form.imageInput"
                      :placeholder="$t('DI.imageInputPro')" />
          </el-form-item>
          <el-form-item v-else
                        key="imageOption"
                        :label="$t('DI.imageSelection')"
                        prop="imageOption">
            <span>&nbsp;{{ defineImage }}&nbsp;&nbsp;</span>
            <el-select v-model="form.imageOption"
                       :placeholder="$t('DI.imagePro')">
              <el-option v-for="(item,index) in imageOptionList"
                         :label="item"
                         :key="index"
                         :value="item">
              </el-option>
            </el-select>
          </el-form-item>
        </el-row>
      </div>
      <div v-if="currentStep===2"
           key="step2">
        <div class="subtitle">
          {{ $t('DI.computingResource') }}
        </div>
        <el-row>
          <el-col :span="12">
            <el-form-item :label="$t('DI.jobType')"
                          prop="job_type">
              <el-select v-model="form.job_type"
                         :placeholder="$t('DI.jobTypePro')">
                <el-option :label="$t('DI.Single')"
                           value="Local">
                </el-option>
                <el-option :label="$t('DI.distributed')"
                           value="dist-tf">
                </el-option>
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item :label="$t('ns.nameSpace')"
                          prop="namespace">
              <el-select v-model="form.namespace"
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
        </el-row>
        <el-row>
          <el-col :span="12">
            <el-form-item :label="$t('DI.cpu')"
                          prop="cpus">
              <el-input v-model="form.cpus"
                        :placeholder="$t('DI.CPUPro')">
                <span slot="append">Core</span>
              </el-input>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item :label="$t('DI.gpu')"
                          prop="gpus">
              <el-input v-model="form.gpus"
                        :placeholder="$t('DI.GPUPro')">
                <span slot="append">{{ $t('DI.block') }}</span>
              </el-input>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row>
          <el-col :span="12">
            <el-form-item :label="$t('DI.memory')"
                          prop="memory">
              <el-input v-model="form.memory"
                        :placeholder="$t('DI.memoryPro')">
                <span slot="append">Gb</span>
              </el-input>
            </el-form-item>
          </el-col>
          <el-col :span="12"
                  v-if="form.job_type==='dist-tf'">
            <el-form-item :label="$t('DI.learners')"
                          prop="learners">
              <el-input v-model="form.learners"
                        :placeholder="$t('DI.learnersPro')" />
            </el-form-item>
          </el-col>
        </el-row>
        <div v-if="form.job_type==='dist-tf'">
          <el-row>
            <el-col :span="12">
              <el-form-item :label="$t('DI.pss')"
                            prop="pss">
                <el-input v-model="form.pss"
                          :placeholder="$t('DI.pssPro')" />
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item :label="$t('DI.ps_cpu')"
                            prop="ps_cpu">
                <el-input v-model="form.ps_cpu"
                          :placeholder="$t('DI.ps_cpuPro')">
                  <span slot="append">Core</span>
                </el-input>
              </el-form-item>
            </el-col>
          </el-row>
          <el-row>
            <el-col :span="12">
              <el-form-item :label="$t('DI.ps_memory')"
                            prop="ps_memory">
                <el-input v-model="form.ps_memory"
                          :placeholder="$t('DI.ps_memoryPro')">
                  <span slot="append">Gi</span>
                </el-input>
              </el-form-item>
            </el-col>
          </el-row>
          <el-row>
            <el-col :span="24">
              <el-form-item :label="$t('DI.ps_imageType')"
                            prop="ps_imageType">
                <el-radio-group v-model="form.ps_imageType">
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
                <span>&nbsp;{{ defineImage }}&nbsp;&nbsp;</span>
                <el-input v-model="form.ps_imageInput"
                          :placeholder="$t('DI.imageInputPro')" />
              </el-form-item>
              <el-form-item v-else
                            key="psImageSelect"
                            :label="$t('DI.ps_image')"
                            prop="ps_image">
                <span>&nbsp;{{ defineImage }}&nbsp;&nbsp;</span>
                <el-select v-model="form.ps_image"
                           :placeholder="$t('DI.imagePro')">
                  <el-option v-for="(item,index) in imageOptionList"
                             :label="item"
                             :key="index"
                             :value="item">
                  </el-option>
                </el-select>
              </el-form-item>
            </el-col>
          </el-row>
        </div>
      </div>
      <div v-if="currentStep===3"
           key="step3">
        <div class="subtitle">
          {{ $t('DI.trainingDirectory') }}
        </div>
        <el-row>
          <el-form-item :label="$t('DI.trainingDataStore')"
                        prop="path">
            <el-select v-model="form.path"
                       filterable
                       :placeholder="$t('DI.trainingDataStorePro')">
              <el-option v-for="(item,index) in localPathList"
                         :label="item"
                         :key="index"
                         :value="item">
              </el-option>
            </el-select>
          </el-form-item>
          <el-form-item :label="$t('DI.trainingData')"
                        prop="trainingData">
            <el-input v-model="form.trainingData"
                      :placeholder="$t('DI.dataSubdirectoryPro')" />
          </el-form-item>
          <el-form-item :label="$t('DI.trainingResult')"
                        prop="trainingResults">
            <el-input v-model="form.trainingResults"
                      :placeholder="$t('DI.resultSubdirectoryPro')" />
          </el-form-item>
        </el-row>
      </div>
      <div v-if="currentStep===4"
           key="step4">
        <div class="subtitle">
          {{ $t('DI.jobExecution') }}
        </div>
        <el-row>
          <el-form-item :label="$t('DI.entrance')"
                        prop="command">
            <el-input v-model="form.command"
                      :rows="4"
                      type="textarea"
                      :placeholder="$t('DI.entrancePro')" />
          </el-form-item>
        </el-row>
        <el-row>
          <el-form-item :label="$t('DI.executeCodeSettings')"
                        prop="codeSettings">
            <el-select v-model="form.codeSettings"
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
               class="upload-box">
            <upload :marginLeft="190"
                    :baseUrl="baseUrl"
                    :code-file="fileName"
                    @deleteFile="deleteFile"
                    @uploadFile="uploadFile"></upload>
          </div>
          <el-form-item v-else
                        key="storagePath"
                        :label="$t('DI.shareDirectory')"
                        prop="diStoragePath">
            <el-input v-model="form.diStoragePath"
                      :placeholder="$t('DI.shareDirectoryPro')" />
          </el-form-item>
          <storage-path-aide v-if="form.codeSettings==='storagePath'"
                             :root-path="form.path"
                             :storage-path="form.diStoragePath"></storage-path-aide>
        </el-row>
      </div>
    </el-form>
    <span slot="footer"
          class="dialog-footer">
      <el-button v-if="currentStep>0"
                 type="primary"
                 @click="lastStep">
        {{ $t('common.lastStep') }}
      </el-button>
      <el-button v-if="currentStep<4"
                 type="primary"
                 @click="nextStep">
        {{ $t('common.nextStep') }}
      </el-button>
      <el-button v-if="currentStep===4"
                 type="primary"
                 :disabled="disabled"
                 @click="subInfo">
        {{ $t('common.save') }}
      </el-button>
      <el-button @click="dialogVisible=false">
        {{ $t('common.cancel') }}
      </el-button>
    </span>
  </el-dialog>
</template>
<script>
import Step from './Step'
import yaml from 'js-yaml'
import handleDIJobMixin from '../util/handleDIJobMixin'
import StoragePathAide from './StoragePathAide'
import Upload from './Upload'
export default {
  mixins: [handleDIJobMixin],
  components: {
    Step,
    StoragePathAide,
    Upload
  },
  props: {
    spaceOptionList: {
      type: Array,
      default: () => {
        return []
      }
    },
    localPathList: {
      type: Array,
      default: () => {
        return []
      }
    }
  },
  data () {
    return {
      disabled: false,
      dialogVisible: false,
      currentStep: 0,
      form: {
        name: '',
        cpus: '',
        gpus: '',
        memory: '',
        namespace: '',
        imageType: 'Standard',
        imageOption: '',
        imageInput: '',
        description: '',
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
      },
      fileName: '',
      codeFile: '',
      defineImage: this.FesEnv.DI.defineImage,
      imageOptionList: this.FesEnv.DI.imageOption,
      baseUrl: `${process.env.VUE_APP_BASE_SURL}/di/v1/codeUpload`
    }
  },
  computed: {
    stepList: function () {
      return [this.$t('DI.basicSettings'), this.$t('DI.imageSettings'), this.$t('DI.computingResource'), this.$t('DI.trainingDirectory'), this.$t('DI.jobExecution')]
    },
    ruleValidate: function () {
      // 切换语言表单报错重置表单
      this.resetFormFields()
      return {
        name: [
          { required: true, message: this.$t('DI.trainingjobNameReq') },
          { type: 'string', pattern: new RegExp(/^[0-9a-zA-Z-]*$/), message: this.$t('DI.trainingjobNameFormat') }
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
        description: [
          { required: true, message: this.$t('DI.descriptionReq') }
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
        codeSettings: [
          { required: true, message: this.$t('DI.executeCodeSettingsReq') }
        ],
        diStoragePath: [
          { required: true, message: this.$t('DI.shareDirectoryReq') },
          { pattern: new RegExp(/^[a-zA-Z][0-9a-zA-Z-_/]*$/), message: this.$t('DI.shareDirectoryFormat') }
        ],
        learners: [
          { required: true, message: this.$t('DI.learnersReq') },
          { pattern: new RegExp(/^[1-9]\d*([.]0)*$/), message: this.$t('DI.learnersFormat') }
        ],
        pss: [
          { required: true, message: this.$t('DI.pssReq') },
          { pattern: new RegExp(/^[1-9]\d*([.]0)*$/), message: this.$t('DI.pssFormat') }
        ]
      }
    }
  },
  methods: {
    nextStep () {
      this.$refs.formValidate.validate((valid) => {
        if (valid && this.currentStep < this.stepList.length - 1) {
          if (this.currentStep === 4) {
            if (this.form.codeSettings === 'codeFile' && !this.codeFile) {
              this.$message.error(this.$t('common.uploadReq'))
              return
            }
          }
          this.currentStep++
        }
      })
    },
    lastStep () {
      // 由于返回上一步告警数据会丢失，在返回上一步是获取告警数据
      this.currentStep > 0 && this.currentStep--
    },
    subInfo () {
      this.$refs.formValidate.validate((valid) => {
        if (valid) {
          this.disabled = true
          let param = this.handleGPUSubParam()
          if (this.form.job_type === 'dist-tf') {
            param.learners = parseInt(this.form.learners)
            param.pss = this.form.pss
            param.ps_cpu = parseFloat(this.form.ps_cpu)
            param.ps_image = this.form.ps_imageType === 'Standard' ? this.form.ps_image : this.form.ps_imageInput
            param.ps_memory = this.form.ps_memory + 'Gi'
          }
          let formData = new FormData()
          if (this.form.codeSettings === 'codeFile') {
            param.code_path = this.codeFile
            param.fileName = this.fileName
          } else {
            param.data_stores[0].training_workspace = {}
            param.data_stores[0].training_workspace.container = this.form.diStoragePath
          }
          console.log('param', param)
          let ymlFile = yaml.safeDump(param)
          console.log('ymlFile', ymlFile)
          let blob = new Blob([ymlFile], { type: 'text / yaml' })
          formData.append('manifest', blob)
          this.FesApi.fetchUT(`/di/${this.FesEnv.diApiVersion}/models?version=200`, formData, {
            method: 'post',
            headers: {
              'Content-Type': 'multipart/form-data'
            }
          }).then(() => {
            this.toast()
            this.$parent.getListData()
            this.dialogVisible = false
            this.disabled = false
          }, () => {
            this.disabled = false
          })
        }
      })
    },
    deleteFile () {
      this.codeFile = ''
      this.fileName = ''
    },
    uploadFile (file, key) {
      this.codeFile = file.result.s3Path
      this.fileName = file.fileName
    },
    clearDialogForm () {
      for (let key in this.form) {
        this.form[key] = ''
      }
      this.form.imageType = 'Standard'
      this.form.ps_imageType = 'Standard'
      this.form.codeSettings = 'codeFile'
      this.form.job_type = 'Local'
      this.codeFile = ''
      this.fileName = ''
      this.gpuDeleteFile()
      this.$refs.formValidate.clearValidate()
      setTimeout(() => {
        this.currentStep = 0
      }, 1000)
    }
  }
}
</script>
<style lang="scss" scoped>
.el-date-editor {
  width: 100%;
}
.add-form {
  min-height: 500px;
}
</style>
