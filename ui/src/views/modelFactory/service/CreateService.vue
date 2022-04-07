<template>
  <el-dialog :title="$t('serviceList.createModel')"
             :visible.sync="dialogVisible"
             @close="clearDialogForm"
             custom-class="dialog-style"
             width="1100px">
    <step :step="currentStep"
          :step-list="stepList" />
    <el-form ref="formValidate"
             :model="form"
             class="add-form"
             :rules="ruleValidate"
             label-width="135px">
      <div v-if="currentStep===0"
           key="step0">
        <div class="subtitle">
          {{$t('DI.basicSettings')}}
        </div>
        <el-form-item :label="$t('serviceList.modelName')"
                      prop="service_name">
          <el-input v-model="form.service_name"
                    :placeholder="$t('serviceList.modelNamePro')" />
        </el-form-item>
        <!-- <el-form-item label="部署类型"
                      prop="service_name">
          <el-select v-model="form.deployType">
            <el-option value="single"
                       label="单模型部署">
            </el-option>
            <el-option value="ABTest"
                       label="ABTest">
            </el-option>
          </el-select>
        </el-form-item> -->
        <el-form-item :label="$t('serviceList.modelDescription')">
          <el-input v-model="form.remark"
                    :placeholder="$t('serviceList.modelDescriptionPro')"
                    type="textarea" />
        </el-form-item>
      </div>
      <div v-if="currentStep===1"
           key="step1">
        <div class="subtitle">
          {{$t('DI.computingResource')}}
        </div>
        <div v-if="form.deployType==='single'||!form.deployType">
          <el-row>
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
            <el-col :span="12">
              <el-form-item label="CPU"
                            prop="cpu">
                <el-input v-model="form.cpu"
                          :placeholder="$t('DI.CPUPro')">
                  <span slot="append">Core</span>
                </el-input>
              </el-form-item>
            </el-col>
          </el-row>
          <el-row>
            <el-col :span="12">
              <el-form-item :label="$t('AIDE.memory1')"
                            prop="memory">
                <el-input v-model="form.memory"
                          :placeholder="$t('DI.memoryPro')">
                  <span slot="append">Gi</span>
                </el-input>
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="GPU"
                            prop="gpu">
                <el-input v-model="form.gpu"
                          :placeholder="$t('DI.GPUPro')">
                  <span slot="append">{{$t('AIDE.block')}}</span>
                </el-input>
              </el-form-item>
            </el-col>
          </el-row>
        </div>
        <div v-else>
          <el-row>
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
          <div v-for="(item,index) in form.ABTestSource"
               :key="index">
            <div class="subtitle">
              {{ABTestLable[index]}}
            </div>
            <el-row>
              <el-col :span="12">
                <el-form-item label="CPU">
                  <el-input v-model="form.ABTestSource[index].cpu"
                            :placeholder="$t('DI.CPUPro')">
                    <span slot="append">Core</span>
                  </el-input>
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item :label="$t('AIDE.memory1')">
                  <el-input v-model="form.ABTestSource[index].memory"
                            :placeholder="$t('DI.memoryPro')">
                    <span slot="append">Gi</span>
                  </el-input>
                </el-form-item>
              </el-col>
            </el-row>
            <el-row>
              <el-col :span="12">
                <el-form-item label="GPU">
                  <el-input v-model="form.ABTestSource[index].gpu"
                            :placeholder="$t('DI.GPUPro')">
                    <span slot="append">{{$t('AIDE.block')}}</span>
                  </el-input>
                </el-form-item>
              </el-col>
            </el-row>
          </div>
        </div>
      </div>
      <div v-show="currentStep===2"
           key="step2">
        <div class="subtitle">
          {{$t('serviceList.deploymentInfoSingle')}}
        </div>
        <model-setting v-if="dialogVisible"
                       ref="modelSetting"
                       :groupList="groupList"></model-setting>
      </div>
      <!-- <div v-if="currentStep===3"
           key="step3">
        <div class="subtitle">
          {{$t('serviceList.logDirectorySettings')}}
        </div>
        <el-form-item :label="$t('serviceList.logDirectory')"
                      prop="log_path">
          <el-select v-model="form.log_path"
                     filterable
                     :placeholder="$t('serviceList.logDirectoryPro')">
            <el-option v-for="(item,index) in localPathList"
                       :label="item"
                       :key="index"
                       :value="item">
            </el-option>
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('serviceList.subdirectory')"
                      prop="subdirectory">
          <el-input v-model="form.subdirectory"
                    :placeholder="$t('serviceList.subdirectoryPro')" />
        </el-form-item>
      </div> -->
    </el-form>
    <span slot="footer"
          class="dialog-footer">
      <el-button v-if="currentStep>0"
                 type="primary"
                 @click="lastStep">
        {{ $t('common.lastStep') }}
      </el-button>
      <el-button v-if="currentStep<2"
                 type="primary"
                 @click="nextStep">
        {{ $t('common.nextStep') }}
      </el-button>
      <el-button v-if="currentStep===2"
                 :disabled="btnDisabled"
                 type="primary"
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
import Step from '../../../components/Step'
import ModelSetting from './ModelSetting'
export default {
  components: {
    Step,
    ModelSetting
  },
  props: {
    spaceOptionList: {
      type: Array,
      default () {
        return []
      }
    },
    // localPathList: {
    //   type: Array,
    //   default () {
    //     return []
    //   }
    // },
    groupList: {
      type: Array,
      default () {
        return []
      }
    }
  },
  data () {
    return {
      form: {
        service_name: '',
        deployType: 'single',
        remark: '',
        namespace: '',
        cpu: '',
        memory: '',
        gpu: '',
        ABTestSource: [
          {
            cpu: '',
            memory: '',
            gpu: ''
          },
          {
            cpu: '',
            memory: '',
            gpu: ''
          }
        ]
        // log_path: '',
        // subdirectory: '',
      },
      parameters: [{
        name: '',
        type: '',
        value: ''
      }],
      source: {},
      currentStep: 0,

      btnDisabled: false,
      dialogVisible: false,
      ABTestLable: [this.$t('serviceList.modeln', { n: 'A' }), this.$t('serviceList.modeln', { n: 'B' })]
    }
  },
  computed: {
    stepList () { //, this.$t('serviceList.logDirectorySettings')
      return [this.$t('DI.basicSettings'), this.$t('DI.computingResource'), this.$t('serviceList.modelSettings')]
    },
    ruleValidate () {
      this.resetFormFields()
      return {
        service_name: [
          { required: true, message: this.$t('serviceList.modelNameReq') },
          { type: 'string', pattern: new RegExp(/^[a-z][a-z0-9-]*$/), message: this.$t('serviceList.modelNameFormat') }
        ],
        // version: [
        //   { required: true, message: this.$t('serviceList.modelVersionReq') }
        // ],
        namespace: [
          { required: true, message: this.$t('ns.nameSpaceReq') }
        ],
        cpu: [
          { required: true, message: this.$t('DI.CPUReq') },
          { pattern: new RegExp(/^((0{1})([.]\d{1})|([1-9]\d*)([.]\d{1})?)$/), message: this.$t('AIDE.CPUNumberReq') }
        ],
        memory: [
          { required: true, message: this.$t('DI.memoryReq') },
          { pattern: new RegExp(/^[1-9]\d*([.]0)*$/), message: this.$t('AIDE.memoryNumberReq') }
        ],
        gpu: [
          { pattern: new RegExp(/^(([1-9]{1}\d*([.]0)*)|([0]{1}))$/), message: this.$t('AIDE.GPUFormat') }
        ]
        // modelType: [
        //   { required: true, message: this.$t('serviceList.modelTypeReq') }
        // ],
        // deploymentType: [
        //   { required: true, message: this.$t('serviceList.deploymentTypeReq') }
        // ],
        // modelLocation: [
        //   { required: true, message: this.$t('serviceList.modelLocationReq') }
        // ],
        // log_path: [
        //   { required: true, message: this.$t('serviceList.logDirectoryReq') }
        // ],
        // subdirectory: [
        //   { required: true, message: this.$t('serviceList.subdirectoryReq') },
        //   { pattern: new RegExp(/^[a-zA-Z][0-9a-zA-Z-_/]*$/), message: this.$t('serviceList.subdirectoryFormat') }
        // ]
      }
    }
  },
  methods: {
    lastStep () {
      // 由于返回上一步告警数据会丢失，在返回上一步先获取模型参数
      if (this.currentStep === 3) {
        this.parameters = JSON.parse(JSON.stringify(this.$refs.serviceParam.form.parameters))
      }
      this.currentStep > 0 && this.currentStep--
    },
    nextStep () {
      this.$refs.formValidate.validate((valid) => {
        if (valid && this.currentStep < this.stepList.length - 1) {
          this.currentStep++
        }
      })
    },
    subInfo () {
      this.$refs.formValidate.validate((valid) => {
        if (valid) {
          if (this.form.deployType === 'single') {
            const modelSetting = this.$refs.modelSetting
            modelSetting.$refs.formValidate.validate((valid) => {
              if (valid) {
                this.btnDisabled = true
                modelSetting.byVersionIdFindSource()
                let modelSettingForm = modelSetting.form
                this.sendRequest(modelSettingForm)
              }
            })
          } else {

          }
        }
      })
    },
    sendRequest (modelSettingForm) {
      let param = {
        cpu: parseFloat(this.form.cpu),
        type: 'single',
        memory: parseInt(this.form.memory),
        group_id: parseInt(modelSettingForm.group_id),
        service_post_models: [{
          model_name: modelSettingForm.model_name,
          model_type: modelSettingForm.model_type,
          group_name: modelSettingForm.group_name,
          model_version: modelSettingForm.modelVersion,
          modelversion_id: parseInt(modelSettingForm.modelversion_id),
          source: modelSettingForm.source,
          image: modelSettingForm.image
        }]
      }
      param.service_post_models[0].endpoint_type = modelSettingForm.endpoint_type
      if (modelSettingForm.model_type === 'TENSORFLOW') {
        param.service_post_models[0].parameters = modelSettingForm.parameters
      }
      const keyArr = ['service_name', 'namespace', 'gpu', 'remark']
      for (let key of keyArr) {
        param[key] = this.form[key]
      }

      this.FesApi.fetch('/mf/v1/service', param, 'post').then(() => {
        this.dialogVisible = false
        this.toast()
        this.$emit('getListData')
        this.btnDisabled = false
      }).finally(() => {
        this.btnDisabled = false
      })
    },
    clearDialogForm () {
      this.btnDisabled = false
      this.currentStep = 0
      this.$refs.formValidate.resetFields()
      for (let key in this.form) {
        this.form[key] = ''
      }
      this.form.deployType = 'single'
    }
  }
}
</script>
<style lang="scss" scoped>
.add-form {
  min-height: 300px;
}
</style>
