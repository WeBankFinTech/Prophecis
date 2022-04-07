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
        <el-form-item :label="$t('flow.nodeName')"
                      prop="name">
          <el-input v-model="form.name"
                    :placeholder="$t('flow.nodeNamePro')"></el-input>
        </el-form-item>
        <el-form-item :label="$t('flow.reportName')"
                      prop="report_name">
          <el-input v-model="form.report_name"
                    :placeholder="$t('flow.reportNamePro')">
          </el-input>
        </el-form-item>
        <el-form-item :label="$t('flow.modelFactory')"
                      prop="factory_name">
          <el-input v-model="form.factory_name"
                    :placeholder="$t('flow.modelFactoryPro')"></el-input>
        </el-form-item>
        <div class="subtitle">
          {{$t('flow.directorySet')}}
        </div>
        <el-form-item :label="$t('flow.storageRootDirectory')"
                      prop="path">
          <el-select v-model="form.path"
                     filterable
                     :placeholder="$t('DI.trainingDataStorePro')">
            <el-option v-for="(item,index) in localPathList"
                       :label="item"
                       :value="item"
                       :key="index">
            </el-option>
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('modelList.filePath')"
                      prop="file_path">
          <el-input v-model="form.file_path"
                    maxlength="400"
                    :placeholder="$t('flow.filePathPro')" />
        </el-form-item>
      </div>
      <div v-show="activeName==='1'"
           key="step1">
        <model-param ref="modelParam"
                     :readonly="readonly"></model-param>
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
import ModelParam from './ModelParam'
import nodeFormMixin from './nodeFormMixin'
const initForm = {
  name: '',
  report_name: '',
  factory_name: '',
  path: '',
  file_path: ''
}
export default {
  mixins: [modelMixin, nodeFormMixin],
  components: { ModelParam },
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
      localPathList: [],
      form: { ...initForm }
    }
  },
  created () {
    const header = this.handlerFromDssHeader('get')
    this.getDIStoragePath(header)
  },
  computed: {
    ruleValidate: function () {
      return {
        name: [
          { required: true, message: this.$t('flow.nodeNoEmpty') },
          { type: 'string', pattern: new RegExp(/^[0-9a-zA-Z-_]*$/), message: this.$t('flow.nodeNameTip') }
        ],
        report_name: [
          { required: true, message: this.$t('flow.reportNameReq') },
          { type: 'string', pattern: new RegExp(/^[0-9a-zA-Z-_]*$/), message: this.$t('flow.reportNameFormat') }
        ],
        factory_name: [
          { required: true, message: this.$t('flow.modelFactoryReq') },
          { type: 'string', pattern: new RegExp(/^[\u4e00-\u9fa5_0-9a-zA-Z-]*$/), message: this.$t('flow.modelFactoryFormat') }
        ],
        path: [
          { required: true, message: this.$t('DI.trainingDataStoreReq') },
          { pattern: new RegExp(/^\/[a-zA-Z][0-9a-zA-Z-_/]*$/), message: this.$t('DI.trainingDataStoreFormat') }
        ],
        file_path: [
          { required: true, message: this.$t('flow.filePathReq') },
          { pattern: new RegExp(/^[a-zA-Z][0-9a-zA-Z-_./]*$/), message: this.$t('flow.filePathFormat') }
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
    save (evt, temSave = false) {
      this.$refs.formValidate.validate(async (valid) => {
        if (valid) {
          let modelParam = await this.$refs.modelParam.save()
          if (modelParam) {
            let param = { ...this.form }
            param.filepath = this.form.path + '/' + this.form.file_path
            param.param_type = modelParam.param_type
            param.group_id = modelParam.group_id
            if (modelParam.param_type === 'model') {
              param.model_name = modelParam.model_name
              param.model_id = modelParam.model_id
              param.model_version_id = modelParam.modelversion_id
              param.model_version = modelParam.versionName
            } else {
              param.model_name = modelParam.cutom_model_name
              param.model_version = modelParam.cutom_version
            }
            this.handleCurrentNode(param)
            this.$emit('saveNode', this.currentNode, temSave)
          } else {
            if (temSave) {
              return
            }
            this.$emit('changeTab', '1')
          }
        } else {
          if (temSave) {
            return
          }
          this.$emit('changeTab', '0')
        }
      })
    },
    handleReportPushNodeData (form) {
      let formData = {}
      for (let key in form) {
        formData[key] = form[key]
      }
      this.form = { ...formData }
      this.$nextTick(() => {
        let modelForm = this.$refs.modelParam.form
        this.$refs.modelParam.getModelOption(form.group_id, true, form.model_id)
        modelForm.group_id = form.group_id
        modelForm.param_type = form.param_type
        if (form.param_type === 'model') {
          this.$refs.modelParam.getVersionOption(form.model_id)
          modelForm.model_id = form.model_id
          modelForm.modelversion_id = form.model_version_id
        } else {
          modelForm.cutom_model_name = form.model_name
          modelForm.cutom_version = form.model_version
        }
      })
    }
  }
}
</script>
