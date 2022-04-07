<template>
  <div>
    <el-form ref="formValidate"
             :model="form"
             :disabled="readonly"
             :rules="ruleValidate"
             label-position="top"
             class="node-parameter-bar">
      <div v-show="activeName==='0'"
           key="step0">
        <div class="subtitle">
          {{$t('DI.basicSettings')}}
        </div>
        <el-form-item :label="$t('flow.nodeName')"
                      prop="name">
          <el-input v-model="form.name"></el-input>
        </el-form-item>
        <el-form-item :label="$t('serviceList.modelName')"
                      prop="service_name">
          <el-input v-model="form.service_name"
                    :placeholder="$t('serviceList.modelNamePro')" />
        </el-form-item>
        <el-form-item :label="$t('serviceList.modelDescription')">
          <el-input v-model="form.remark"
                    :placeholder="$t('serviceList.modelDescriptionPro')"
                    type="textarea" />
        </el-form-item>
      </div>
      <div v-show="activeName==='1'"
           key="step1">
        <div class="subtitle">
          {{$t('DI.computingResource')}}
        </div>
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
        <el-form-item label="CPU"
                      prop="cpu">
          <el-input v-model="form.cpu"
                    :placeholder="$t('DI.CPUPro')">
            <span slot="append">Core</span>
          </el-input>
        </el-form-item>
        <el-form-item :label="$t('AIDE.memory1')"
                      prop="memory">
          <el-input v-model="form.memory"
                    :placeholder="$t('DI.memoryPro')">
            <span slot="append">Gi</span>
          </el-input>
        </el-form-item>
        <el-form-item label="GPU"
                      prop="gpu">
          <el-input v-model="form.gpu"
                    :placeholder="$t('DI.GPUPro')">
            <span slot="append">{{$t('AIDE.block')}}</span>
          </el-input>
        </el-form-item>
      </div>
      <div v-show="activeName==='2'"
           key="step2">
        <div class="subtitle">
          {{$t('serviceList.deploymentInfoSingle')}}
        </div>
        <el-form-item :label="$t('image.parameterType')"
                      prop="param_type">
          <el-select v-model="form.param_type">
            <el-option :label="$t('image.modelList')"
                       value="model"></el-option>
            <el-option :label="$t('modelList.custom')"
                       value="custom"></el-option>
          </el-select>
        </el-form-item>
        <div v-if="form.param_type==='model'">
          <el-form-item :label="$t('user.userGroup')"
                        prop="group_id">
            <el-select v-model="form.group_id"
                       :filterable="true"
                       :placeholder="$t('user.userGroupSettingsPro')"
                       @change="changeGroup">
              <el-option v-for="(item,index) of groupList"
                         :label="item.name"
                         :value="item.id"
                         :key="index">
              </el-option>
            </el-select>
          </el-form-item>
          <el-form-item v-if="form.group_id"
                        :label="$t('modelList.modelName')"
                        prop="model_id">
            <el-select v-model="form.model_id"
                       :placeholder="$t('serviceList.modelNamePro1')"
                       @change="selectChangeModel">
              <el-option v-for="(item,index) of modelOptionList"
                         :label="item.model_name"
                         :value="item.id"
                         :key="index">
              </el-option>
            </el-select>
          </el-form-item>
          <div v-if="form.model_id">
            <el-form-item :label="$t('serviceList.modelVersion')"
                          prop="modelversion_id">
              <el-select v-model="form.modelversion_id"
                         @change="changeVersionId"
                         :placeholder="$t('serviceList.modelVersionPro')">
                <el-option v-for="(item,index) of versionOptionList"
                           :label="item.name"
                           :value="item.id"
                           :key="index">
                </el-option>
              </el-select>
            </el-form-item>
            <el-form-item :label="$t('serviceList.modelType')"
                          prop="model_type">
              <el-input v-model="form.model_type"
                        disabled></el-input>
            </el-form-item>
          </div>
          <div v-if="form.model_type">
            <el-form-item :label="$t('serviceList.serviceType')"
                          :placeholder="$t('serviceList.serviceTypePro')"
                          prop="endpoint_type">
              <el-select v-model="form.endpoint_type">
                <el-option value="REST"
                           label="REST">
                </el-option>
                <el-option v-if="form.model_type!=='TENSORFLOW'"
                           value="GRPC"
                           label="GRPC">
                </el-option>
              </el-select>
            </el-form-item>
          </div>
          <el-form-item v-if="form.model_type==='CUSTOM'&&form.modelversion_id"
                        :label="$t('DI.imageSelection')"
                        prop="image">
            <el-select v-model="form.image"
                       :placeholder="$t('DI.imagePro')">
              <el-option v-for="(item,index) of imageOptionList"
                         :label="item.image_name"
                         :value="item.image_name"
                         :key="index">
              </el-option>
            </el-select>
          </el-form-item>
          <div v-if="form.model_type==='TENSORFLOW'">
            <div class="subtitle">
              {{$t('serviceList.serviceParam')}}
              <i class="el-icon-circle-plus-outline add"
                 @click="addParam" />
              <i class="el-icon-remove-outline delete"
                 @click="deleteParam(index)" />
            </div>
            <div v-for="(item,index) in form.parameters"
                 :key="index">
              <el-form-item class="service-param"
                            :prop="`parameters[${index}].name`">
                <el-input v-model="item.name"
                          :placeholder="$t('serviceList.paramNamePro')"></el-input>
              </el-form-item>
              <el-form-item :prop="`parameters[${index}].type`">
                <el-input v-model="item.type"
                          :placeholder="$t('serviceList.paramTypePro')"></el-input>
              </el-form-item>
              <el-form-item :prop="`parameters[${index}].value`">
                <el-input v-model="item.value"
                          :placeholder="$t('serviceList.paramValuePro')"></el-input>
              </el-form-item>
            </div>
          </div>
        </div>
        <div v-else
             key="custom">
          <el-form-item :label="$t('user.userGroup')"
                        prop="group_id">
            <el-select v-model="form.group_id"
                       :filterable="true"
                       :placeholder="$t('user.userGroupSettingsPro')">
              <el-option v-for="(item,index) of groupList"
                         :label="item.name"
                         :value="item.id"
                         :key="index">
              </el-option>
            </el-select>
          </el-form-item>
          <el-form-item :label="$t('modelList.modelName')"
                        prop="cutom_model_name">
            <el-input v-model="form.cutom_model_name"
                      :placeholder="$t('modelList.modelNamePro')"></el-input>
          </el-form-item>
          <el-form-item :label="$t('flow.modelVersion')"
                        prop="cutom_version">
            <el-input v-model="form.cutom_version"
                      :placeholder="$t('flow.modelVersionPro')"></el-input>
          </el-form-item>
        </div>
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
const initForm = {
  name: '',
  service_name: '',
  remark: '',
  namespace: '',
  cpu: '',
  memory: '',
  gpu: '',
  // log_path: '',
  // subdirectory: '',
  group_id: '',
  model_id: '',
  modelversion_id: '',
  model_type: '',
  endpoint_type: '',
  parameters: [{
    name: '',
    type: '',
    value: ''
  }],
  param_type: 'model',
  cutom_model_name: '',
  cutom_version: 'v1'
}
export default {
  mixins: [modelMixin, nodeFormMixin],
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
    }
  },
  data () {
    return {
      currentNode: {},
      form: JSON.parse(JSON.stringify(initForm)),
      source: {},
      currentStep: 0,
      modelOptionList: [],
      spaceOptionList: [],
      groupList: [],
      btnDisabled: false,
      versionOptionList: [],
      imageOptionList: [],
      dialogVisible: false,
      serviceModels: {}
    }
  },
  created () {
    let header = this.handlerFromDssHeader('get')
    this.getGroupList(header)
    this.getDISpaceOption(header)
    this.getDIStoragePath(header)
  },
  computed: {
    ruleValidate () {
      return {
        name: [
          { required: true, message: this.$t('flow.nodeNoEmpty') },
          { type: 'string', pattern: new RegExp(/^[0-9a-zA-Z-_]*$/), message: this.$t('flow.nodeNameTip') }
        ],
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
        ],
        image: [
          { required: true, message: this.$t('DI.imageInputReq') }
        ],
        group_id: [
          { required: true, message: this.$t('user.userGroupSettingsReq') }
        ],
        model_id: [
          { required: true, message: this.$t('modelList.modelNameReq') }
        ],
        modelversion_id: [
          { required: true, message: this.$t('serviceList.modelVersionReq') }
        ],
        endpoint_type: [
          { required: true, message: this.$t('serviceList.serviceTypeReq') }
        ],
        'parameters[0].name': [
          { required: true, message: this.$t('serviceList.paramNameReq') }
        ],
        'parameters[0].type': [
          { required: true, message: this.$t('serviceList.paramTypeReq') }
        ],
        'parameters[0].value': [
          { required: true, message: this.$t('serviceList.paramValueReq') }
        ],
        param_type: [
          { required: true, message: this.$t('serviceList.paramTypeReq') }
        ],
        cutom_model_name: [
          { required: true, message: this.$t('modelList.modelNameReq') },
          { type: 'string', pattern: new RegExp(/^[a-zaA-Z][a-zA-Z0-9-]*$/), message: this.$t('modelList.modelNameFormat') }

        ],
        cutom_version: [
          { required: true, message: this.$t('flow.modelVersionReq') },
          { type: 'string', pattern: new RegExp(/^v\d*$/), message: this.$t('flow.modelVersionFormat') }
        ]
      }
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
  methods: {
    handleServiceNodeData (form) {
      this.getModelOption(form.group_id, true, form.model_id)
      let param = {}
      const serviceModel = form.service_post_models[0]
      if (form.param_type === 'model') {
        this.getVersionOption(form.model_id)
        param.model_id = form.model_id
        if (serviceModel.model_type === 'TENSORFLOW') {
          param.parameters = serviceModel.parameters
        }
        if (serviceModel.model_type === 'CUSTOM') {
          param.image = serviceModel.image
          this.changeVersionId(serviceModel.modelversion_id)
        }
        param.modelversion_id = serviceModel.modelversion_id
        param.model_type = serviceModel.model_type
        param.endpoint_type = serviceModel.endpoint_type
      } else {
        param.cutom_model_name = serviceModel.model_name
        param.cutom_version = serviceModel.model_version
        param.model_version_id = ''
        delete form.model_name
      }
      const keyArr = ['name', 'service_name', 'namespace', 'gpu', 'remark', 'cpu', 'memory', 'group_id', 'param_type']
      for (let key of keyArr) {
        param[key] = form[key]
      }
      for (let key in initForm) {
        if (param[key] === undefined) {
          param[key] = initForm[key]
        }
      }
      this.form = param
    },
    changeGroup (val) {
      this.form.model_id = ''
      this.form.modelversion_id = ''
      this.versionOptionList = []
      this.getModelOption(val)
    },
    getModelOption (val, isSetModel, modelId) {
      let header = this.handlerFromDssHeader('get')
      this.FesApi.fetch(`/mf/v1/modelsByGroupID/${val}`, { page: 1, size: 1000 }, header).then((res) => {
        this.modelOptionList = this.formatResList(res)
        if (isSetModel) {
          this.setServiceModels(modelId)
        }
      })
    },
    changeVersionId (versionId) {
      let header = this.handlerFromDssHeader('get')
      this.FesApi.fetch(`/mf/v1/images/${versionId}`, {}, header).then((res) => {
        this.imageOptionList = this.formatResList(res)
      })
    },
    selectChangeModel (val) {
      this.form.modelversion_id = ''
      this.imageOptionList = []
      this.setServiceModels(val)
      this.getVersionOption(val)
    },
    setServiceModels (val) {
      for (let item of this.modelOptionList) {
        if (item.id === val) {
          this.form.model_type = item.model_type
          this.serviceModels = {
            model_name: item.model_name,
            model_type: item.model_type,
            group_name: item.group.name
          }
          break
        }
      }
    },
    getVersionOption (val) {
      let header = this.handlerFromDssHeader('get')
      this.FesApi.fetch(`/mf/v1/modelversions/${val}`, { page: 1, size: 1000 }, header).then((res) => {
        let optionList = []
        for (let item of res.modelVersions) {
          optionList.push({
            id: item.id,
            name: item.version,
            source: item.source
          })
        }
        this.versionOptionList = this.formatResList(optionList)
      })
    },
    judgmentFormTab () {
      return {
        '0': ['name', 'service_name'],
        '1': ['namespace', 'cpu', 'memory', 'gpu'],
        '2': ['param_type', 'group_id', 'model_id', 'modelversion_id', 'model_type', 'endpoint_type', 'image', 'cutom_model_name']
      }
    },
    save (evt, temSave = false) {
      this.$refs.formValidate.validate((valid, errObj) => {
        if (valid) {
          this.sendRequest(temSave)
        } else {
          if (temSave) {
            return
          }
          const keyArr = this.judgmentFormTab()
          for (let key in this.ruleValidate) {
            if (key.indexOf('parameters[') > -1) {
              keyArr['2'].push(key)
            }
          }
          this.validateFormChangeTab(keyArr, errObj, this.activeName)
        }
      })
    },
    sendRequest (temSave) {
      let param = {
        cpu: parseFloat(this.form.cpu),
        type: 'single',
        memory: parseInt(this.form.memory),
        group_id: parseInt(this.form.group_id)
        // log_path: this.form.log_path + '/' + this.form.subdirectory
      }
      if (this.form.param_type === 'model') {
        let modelVersion = ''
        let source = ''
        for (let item of this.versionOptionList) {
          if (item.id === this.form.modelversion_id) {
            modelVersion = item.name
            source = item.source
            break
          }
        }
        param.service_post_models = [{
          ...this.serviceModels,
          model_version: modelVersion,
          modelversion_id: parseInt(this.form.modelversion_id),
          source,
          image: this.form.image,
          endpoint_type: this.form.endpoint_type
        }]
        if (this.form.model_type === 'TENSORFLOW') {
          param.service_post_models[0].parameters = this.form.parameters
        }
        param.model_id = this.form.model_id
      } else {
        let groupName = ''
        for (let item of this.groupList) {
          if (item.id === this.form.group_id) {
            groupName = item.name
          }
        }
        param.service_post_models = [{
          model_name: this.form.cutom_model_name,
          group_name: groupName,
          image: this.form.image
        }]
        param.model_version = this.form.cutom_version
      }

      const keyArr = ['name', 'service_name', 'namespace', 'gpu', 'remark', 'param_type']
      for (let key of keyArr) {
        param[key] = this.form[key]
      }
      this.handleCurrentNode(param)
      this.$emit('saveNode', this.currentNode, temSave)
    },
    addParam () {
      this.form.parameters.push({
        name: '',
        type: '',
        value: ''
      })
      this.addValidate()
    },
    addValidate () {
      const len = this.form.parameters.length - 1
      this.ruleValidate[`parameters[${len}].name`] = this.ruleValidate['parameters[0].name']
      this.ruleValidate[`parameters[${len}].type`] = this.ruleValidate['parameters[0].type']
      this.ruleValidate[`parameters[${len}].value`] = this.ruleValidate['parameters[0].value']
    },
    deleteParam (index) {
      if (this.form.parameters.length > 1) {
        this.form.parameters.splice(index, 1)
      }
    }
  }
}
</script>
<style lang="scss" scoped>
.add-form {
  min-height: 300px;
}
</style>
