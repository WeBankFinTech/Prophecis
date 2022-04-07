<template>
  <div class="form-content">
    <el-form ref="formValidate"
             :model="form"
             :rules="ruleValidate"
             size="mini"
             label-position="top"
             :disabled="readonly"
             class="node-parameter-bar">
      <div v-show="activeName==='0'"
           key="step0">
        <el-form-item :label="$t('flow.nodeName')"
                      prop="name">
          <el-input v-model="form.name"></el-input>
        </el-form-item>
        <el-form-item :label="$t('flow.nodeDesc')"
                      prop="description">
          <el-input v-model="form.description"
                    type="textarea"></el-input>
        </el-form-item>
      </div>
      <div v-show="activeName==='1'"
           key="step1">
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
        <el-form-item :label="$t('DI.pythonType')"
                      prop="pythonType">
          <el-radio-group v-model="form.pythonType">
            <el-radio label="Standard">
              {{ $t('DI.standard') }}
            </el-radio>
            <el-radio label="Custom">
              {{ $t('DI.custom') }}
            </el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item v-if="form.pythonType === 'Custom'"
                      key="pythonInput"
                      :label="$t('DI.HDFSPath')"
                      prop="pythonInput">
          <el-input v-model="form.pythonInput"
                    :placeholder="$t('DI.HDFSPathPro')" />
        </el-form-item>
        <el-form-item v-else
                      key="pythonSelect"
                      :label="$t('DI.pythonSelect')"
                      prop="pythonOption">
          <el-select v-model="form.pythonOption"
                     :placeholder="$t('DI.pythonSelectPro')">
            <el-option v-for="(item,index) in pythonOptionList"
                       :label="item.key"
                       :key="index"
                       :value="item.value">
            </el-option>
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('DI.driverMemory')"
                      prop="driver_memory">
          <el-input v-model="form.driver_memory"
                    maxlength="10"
                    :placeholder="$t('DI.driverMemoryPro')">
            <span slot="append">Gi</span>
          </el-input>
        </el-form-item>
        <el-form-item :label="$t('DI.executorMemory')"
                      prop="executor_memory">
          <el-input v-model="form.executor_memory"
                    maxlength="10"
                    :placeholder="$t('DI.executorMemoryPro')">
            <span slot="append">Gi</span>
          </el-input>
        </el-form-item>
        <el-form-item :label="$t('DI.executorInstances')"
                      prop="executor_cores">
          <el-input v-model="form.executor_cores"
                    maxlength="10"
                    :placeholder="$t('DI.executorInstancesPro')">
            <span slot="append">Core</span>
          </el-input>
        </el-form-item>
        <el-form-item :label="$t('DI.linkisInstance')"
                      prop="executors">
          <el-input v-model="form.executors"
                    maxlength="10"
                    :placeholder="$t('DI.linkisInstancePro')" />
        </el-form-item>
        <el-form-item :label="$t('DI.queue')"
                      prop="queue">
          <el-input v-model="form.queue"
                    maxlength="225"
                    :placeholder="$t('DI.queuePro')" />
        </el-form-item>
      </div>
      <div v-show="activeName==='2'"
           key="step2">
        <div class="subtitle ">
          {{ $t('DI.dateSettings') }}
          <i class="el-icon-circle-plus-outline add"
             @click="hdpAddExecutorsTask('archives')" />
          <i class="el-icon-remove-outline delete"
             @click="hdpDeleteExecutorsTask('archives')" />
        </div>
        <el-row v-for="(item,index) of form.archives"
                :key="`archives${index}`">
          <el-form-item :label="$t('DI.uploadModeSettings')"
                        :prop="`archives[${index}]codeSettings`">
            <el-select v-model="item.codeSettings"
                       :placeholder="$t('DI.uploadModeSettingsPro')">
              <el-option :label="$t('DI.manualUpload')"
                         value="codeFile">
              </el-option>
              <el-option :label="$t('DI.HDFSPath')"
                         value="HDFSPath">
              </el-option>
            </el-select>
          </el-form-item>
          <div v-if="item.codeSettings==='codeFile'"
               key="codeFile"
               class="tofs-upload">
            <span v-if="item.fileName"
                  class="file-name">
              {{ item.fileName }}
              <i class="el-icon-close delete-file"
                 @click="hdpDeleteFile('archives',index)" />
            </span>
            <el-upload :ref="`archivesUpload${index}`"
                       :action="baseUrl"
                       :headers="header"
                       :on-success="hadUploadFile"
                       class="upload"
                       :before-upload="hdpBeforeUpload">
              <el-button type="primary"
                         @click="hdpSubmitUpload('archives',index)">
                {{ $t('DI.manualUpload') }}
              </el-button>
            </el-upload>
          </div>
          <el-form-item v-if="item.codeSettings==='HDFSPath'"
                        key="HDFSPath"
                        :label="$t('DI.HDFSPath')"
                        :prop="`archives[${index}].HDFSPath`">
            <el-input v-model="item.HDFSPath"
                      :placeholder="$t('DI.HDFSPathPro')" />
          </el-form-item>
        </el-row>

        <div class="subtitle alarm-set">
          {{ $t('DI.codeSettings') }}
        </div>
        <el-row v-for="(item,index) of form.pyfile"
                :key="`pyfile${index}`">
          <el-form-item :label="$t('DI.uploadModeSettings')"
                        :prop="`pyfile[${index}].codeSettings`">
            <el-select v-model="item.codeSettings"
                       :placeholder="$t('DI.uploadModeSettingsPro')">
              <el-option :label="$t('DI.manualUpload')"
                         value="codeFile">
              </el-option>
              <el-option :label="$t('DI.HDFSPath')"
                         value="HDFSPath">
              </el-option>
            </el-select>
          </el-form-item>
          <div v-if="item.codeSettings==='codeFile'"
               key="codeFile"
               class="tofs-upload">
            <span v-if="item.fileName"
                  class="file-name">
              {{ item.fileName }}
              <i class="el-icon-close delete-file"
                 @click="hdpDeleteFile('pyfile',index)" />
            </span>
            <el-upload :ref="`pyfileUpload${index}`"
                       :action="baseUrl"
                       :headers="header"
                       :on-success="hadUploadFile"
                       class="upload"
                       :before-upload="hdpBeforeUpload">
              <el-button type="primary"
                         @click="hdpSubmitUpload('pyfile',index)">
                {{ $t('DI.manualUpload') }}
              </el-button>
            </el-upload>
          </div>
          <el-form-item v-if="item.codeSettings==='HDFSPath'"
                        key="HDFSPath"
                        :label="$t('DI.HDFSPath')"
                        :prop="`pyfile[${index}].HDFSPath`">
            <el-input v-model="item.HDFSPath"
                      :placeholder="$t('DI.HDFSPathPro')" />
          </el-form-item>
        </el-row>
        <el-form-item :label="$t('DI.entrance')"
                      prop="command">
          <el-input v-model="form.command"
                    :rows="4"
                    type="textarea"
                    :placeholder="$t('DI.entrancePro')" />
        </el-form-item>
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
import handleDIJobMixin from '../../../../util/handleDIJobMixin'
import nodeFormMixin from './nodeFormMixin'
import { Upload } from 'element-ui'
import modelMixin from '../../../../util/modelMixin'
const initForm = {
  name: '',
  description: '',
  pythonType: 'Standard',
  pythonOption: '',
  pythonInput: '',
  driver_memory: '',
  executor_memory: '',
  executor_cores: '',
  executors: '',
  queue: '',
  command: '',
  namespace: '',
  pyfile: [{
    codeSettings: 'codeFile',
    HDFSPath: '',
    fileName: ''
  }],
  archives: [{
    codeSettings: 'codeFile',
    HDFSPath: '',
    fileName: ''
  }]
}
export default {
  mixins: [modelMixin, handleDIJobMixin, nodeFormMixin],
  components: {
    ElUpload: Upload
  },
  props: {
    nodeData: {
      type: Object,
      default: () => {
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
    fromDss: {
      type: Boolean,
      default: false
    }
  },
  data () {
    return {
      currentNode: {}, // 保存节点信息
      pythonOptionList: this.FesEnv.DI.pythonOption,
      form: JSON.parse(JSON.stringify(initForm)),
      alarmData: null,
      header: {
        'Mlss-Userid': localStorage.getItem('userId')
      },
      spaceOptionList: []
    }
  },
  watch: {
    nodeData: {
      immediate: true,
      handler (val) {
        this.handleOriginNodeData(val, initForm)
      }
    }
  },
  mounted () {
    const header = this.handlerFromDssHeader('get')
    this.header.headers = { ...header }
    delete this.header.method
    this.getDISpaceOption(header)
  },
  computed: {
    baseUrl () {
      return `${process.env.VUE_APP_BASE_SURL}/linkis/upload`
    },
    currentUplod () {
      return {
        uploadModule: '',
        index: 0
      }
    },
    ruleValidate: function () {
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
        pythonType: [
          { required: true, message: this.$t('DI.pythonTypeReq') }
        ],
        pythonOption: [
          { required: true, message: this.$t('DI.pythonSelectReq') }
        ],
        pythonInput: [
          { required: true, message: this.$t('DI.HDFSPathReq') },
          { pattern: new RegExp(/^[0-9a-zA-Z-_.:/]*$/), message: this.$t('DI.HDFSPathFormat') }
        ],
        driver_memory: [
          { required: true, message: this.$t('DI.driverMemoryReq') },
          { pattern: new RegExp(/^[1-9]\d*([.]0)*$/), message: this.$t('DI.driverMemoryFormat') }
        ],
        executor_cores: [
          { required: true, message: this.$t('DI.executorInstancesReq') },
          { pattern: new RegExp(/^[1-9]\d*([.]0)*$/), message: this.$t('DI.executorInstancesFormat') }
        ],
        executor_memory: [
          { required: true, message: this.$t('DI.executorMemoryReq') },
          { pattern: new RegExp(/^[1-9]\d*([.]0)*$/), message: this.$t('DI.executorMemoryFormat') }
        ],
        executors: [
          { required: true, message: this.$t('DI.linkisInstanceReq') },
          { pattern: new RegExp(/^[1-9]\d*([.]0)*$/), message: this.$t('DI.linkisInstanceFormat') }
        ],
        queue: [
          { required: true, message: this.$t('DI.queueReq') },
          { pattern: new RegExp(/^[a-zA-Z0-9][a-zA-Z0-9_.]*$/), message: this.$t('DI.queueFormat') }
        ],
        'pyfile[0].codeSettings': [
          { required: true, message: this.$t('DI.uploadModeSettingsReq') }
        ],
        'pyfile[0].HDFSPath': [
          { required: true, message: this.$t('DI.HDFSPathReq') },
          { pattern: new RegExp(/^[0-9a-zA-Z-_:./]*$/), message: this.$t('DI.HDFSPathFormat') }
        ],
        'archives[0].codeSettings': [
          { required: true, message: this.$t('DI.uploadModeSettingsReq') }
        ],
        'archives[0].HDFSPath': [
          { pattern: new RegExp(/^[0-9a-zA-Z-_:./]*$/), message: this.$t('DI.HDFSPathFormat') }
        ],
        command: [
          { required: true, message: this.$t('DI.entranceReq') },
          { pattern: new RegExp(/^(.*[^\s])/), message: this.$t('DI.entranceReq') }
        ]
      }
    }
  },
  methods: {
    judgmentFormTab () {
      return {
        '0': ['name', 'description'],
        '1': ['namespace', 'pythonType', 'pythonOption', 'pythonInput', 'driver_memory', 'executor_cores', 'executor_memory', 'executors', 'queue'],
        '2': ['command']
      }
    },
    save (evt, temSave = false) {
      this.$refs.formValidate.validate((valid, errObj) => {
        if (valid) {
          this.disabled = true
          let param = this.handleHdpSubParam()
          this.handleCurrentNode(param)
          if (this.alarmData) {
            this.currentNode.jobContent.ManiFest.job_alert = this.alarmData
          }
          this.$emit('saveNode', this.currentNode, temSave)
        } else {
          if (temSave) {
            return
          }
          if (this.activeName === '2') {
            for (let key in errObj) {
              if (key.indexOf('pyfile[0]') > -1 || key.indexOf('archives[') > -1) {
                return
              }
            }
            if (this.form.pyfile[0].codeSettings === 'codeFile' && !this.form.pyfile[0].fileName) {
              this.$message.error(this.$t('DI.uploadModenotFile'))
              return
            }
          }
          const keyArr = this.judgmentFormTab()
          this.validateFormChangeTab(keyArr, errObj, this.activeName)
          // this.$message.warning(this.$t('message.process.nodeParameter.BCSB'))
        }
      })
    },
    handleHadoopNodeData (param) {
      let formData = {}
      const arr = ['name', 'namespace', 'description']
      for (let item of arr) {
        formData[item] = param[item]
      }
      const tfosArr = ['queue', 'executor_cores', 'executor_memory', 'driver_memory', 'executors', 'command']
      for (let key of tfosArr) {
        formData[key] = param.tfos_request[key]
      }
      formData = this.handelPythonData(formData, param)
      this.handelTfosOriginData(formData, param)
      this.form = formData
    },
    handelPythonData (formData, param) {
      let pythonUrl = param.tfos_request.tensorflow_env.hdfs
      let curPyth = this.FesEnv.DI.pythonOption.find((x) => x.value.indexOf(pythonUrl) > -1)
      if (curPyth !== undefined) {
        formData.pythonType = 'Standard'
        formData.pythonOption = curPyth.value
      } else {
        formData.pythonType = 'Custom'
        formData.pythonInput = pythonUrl
      }
      return formData
    },
    handelTfosOriginData (formData, param) {
      const tfosKey = ['pyfile', 'archives']
      for (let item of tfosKey) {
        this.handleHadoopOrigin(item, formData, param.tfos_request)
      }
    },
    handleHadoopOrigin (moduleKey, formData, param) {
      const taskBaseObj = {
        codeSettings: 'HDFSPath',
        HDFSPath: '',
        fileName: ''
      }
      let moduleArr = []
      if (param[moduleKey].length > 0) {
        for (let j = 0; j < param[moduleKey].length; j++) {
          if (j > 0) {
            this.hdpAddExecutorsValidate(moduleKey, j) // 数据文件设置表单校验规则动态增加
          }
          let item = { ...taskBaseObj }
          item.HDFSPath = param[moduleKey][j].hdfs
          moduleArr.push(item)
        }
      } else {
        moduleArr.push({ ...taskBaseObj })
      }
      formData[moduleKey] = moduleArr
      // this.$refs.hadoopDialog.form[moduleKey] = moduleArr
    }
  }
}
</script>
<style lang="scss" scoped>
.el-form {
  .subtitle {
    padding-left: 0;
  }
}
.tofs-upload {
  .file-name {
    float: left;
    margin-left: 0px;
    padding: 10px 20px;
  }
}

// .save-button {
//   position: fixed;
//   bottom: 0;
//   right: 150px;
// }
// .form-content {
//   padding-bottom: 40px;
// }
</style>
