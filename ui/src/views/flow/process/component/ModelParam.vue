<template>
  <el-form ref="formValidate"
           :model="form"
           :disabled="readonly"
           :rules="ruleValidate"
           label-position="top"
           class="node-parameter-bar">
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
      </div>
    </div>
    <div v-else
         key="custom">
      <el-form-item :label="$t('user.userGroup')"
                    prop="group_id">
        <el-select v-model="form.group_id"
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
  </el-form>
</template>
<script>
import modelMixin from '../../../../util/modelMixin'
import nodeFormMixin from './nodeFormMixin'
const initForm = {
  group_id: '',
  model_id: '',
  modelversion_id: '',
  model_type: '',
  param_type: 'model',
  cutom_model_name: '',
  cutom_version: 'v1'
}
export default {
  mixins: [modelMixin, nodeFormMixin],
  data () {
    return {
      form: { ...initForm },
      modelOptionList: [],
      groupList: [],
      btnDisabled: false,
      versionOptionList: [],
      imageOptionList: [],
      serviceModels: {}
    }
  },
  props: {
    readonly: {
      type: Boolean,
      default: false
    }
  },
  created () {
    let header = this.handlerFromDssHeader('get')
    this.getGroupList(header)
    this.getDIStoragePath(header)
  },
  computed: {
    ruleValidate () {
      return {
        group_id: [
          { required: true, message: this.$t('user.userGroupSettingsReq') }
        ],
        model_id: [
          { required: true, message: this.$t('modelList.modelNameReq') }
        ],
        modelversion_id: [
          { required: true, message: this.$t('serviceList.modelVersionReq') }
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
  methods: {
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
    save () {
      return new Promise((resolve) => {
        this.$refs.formValidate.validate((valid) => {
          if (valid) {
            let versionName = ''
            if (this.form.param_type === 'model') {
              versionName = this.versionOptionList.find((item) => {
                if (item.id === this.form.modelversion_id) {
                  return item.name
                }
              })
            }
            let param = {
              ...this.form,
              versionName,
              model_name: this.serviceModels.model_name
            }
            resolve(param)
          } else {
            resolve(false)
          }
        })
      })
    }
  }
}
</script>
