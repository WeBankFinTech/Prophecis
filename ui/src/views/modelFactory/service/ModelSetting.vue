<template>
  <div>
    <el-form ref="formValidate"
             :model="form"
             class="add-form"
             :rules="ruleValidate"
             :disabled="viewDetail"
             label-width="135px">
      <el-row>
        <el-col :span="12">
          <el-form-item :label="$t('user.userGroup')"
                        prop="group_id">
            <el-input v-if="viewDetail||modify"
                      v-model="form.group_name"
                      :disabled="modify"></el-input>
            <el-select v-else
                       v-model="form.group_id"
                       :placeholder="$t('user.userGroupSettingsPro')"
                       :disabled="modify"
                       :filterable="true"
                       @change="changeGroup">
              <el-option v-for="(item,index) of groupList"
                         :label="item.name"
                         :value="item.id"
                         :key="index">
              </el-option>
            </el-select>
          </el-form-item>
        </el-col>
        <el-col :span="12"
                v-if="form.group_id">
          <el-form-item :label="$t('modelList.modelName')"
                        prop="model_id">
            <el-input v-if="viewDetail||modify"
                      v-model="form.model_name"
                      :disabled="modify"></el-input>
            <el-select v-else
                       v-model="form.model_id"
                       :disabled="modify"
                       :placeholder="$t('serviceList.modelNamePro1')"
                       @change="selectChangeModel">
              <el-option v-for="(item,index) of modelOptionList"
                         :label="item.model_name"
                         :value="item.id"
                         :key="index">
              </el-option>
            </el-select>
          </el-form-item>
        </el-col>
      </el-row>
      <el-row v-if="form.model_id">
        <el-col :span="12">
          <el-form-item :label="$t('serviceList.modelVersion')"
                        prop="modelversion_id">
            <el-select v-model="form.modelversion_id"
                       @change="changeVersionId"
                       :placeholder="$t('serviceList.modelVersionPro')">
              <el-option v-for="(item,index) in versionOptionList"
                         :label="item.name"
                         :value="item.id"
                         :key="index">
              </el-option>
            </el-select>
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item :label="$t('serviceList.modelType')"
                        prop="model_type">
            <el-input v-model="form.model_type"
                      disabled></el-input>
          </el-form-item>
        </el-col>
      </el-row>
      <el-row v-if="form.model_type">
        <el-col :span="12">
          <el-form-item :label="$t('serviceList.serviceType')"
                        :placeholder="$t('serviceList.serviceTypePro')"
                        prop="endpoint_type">
            <el-select v-model="form.endpoint_type"
                       :disabled="modify">
              <el-option value="REST"
                         label="REST">
              </el-option>
              <el-option v-if="form.model_type!=='TENSORFLOW'"
                         value="GRPC"
                         label="GRPC">
              </el-option>
            </el-select>
          </el-form-item>
        </el-col>
      </el-row>
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
      <div v-if="form.model_type==='TENSORFLOW'&&!modify">
        <div class="subtitle">
          {{$t('serviceList.serviceParam')}} <i v-if="!viewDetail"
             class="el-icon-circle-plus-outline add"
             @click="addParam" />
        </div>
        <el-row class="padding-label"
                v-for="(item,index) in form.parameters"
                :key="index"
                :gutter="20">
          <el-col :span="7">
            <el-form-item class="service-param"
                          label-width="0px"
                          :prop="`parameters[${index}].name`">
              <el-input v-model="item.name"
                        :placeholder="$t('serviceList.paramNamePro')"></el-input>
            </el-form-item>
          </el-col>
          <el-col :span="7"
                  class="service-param">
            <el-form-item :prop="`parameters[${index}].type`"
                          label-width="0px">
              <el-input v-model="item.type"
                        :placeholder="$t('serviceList.paramTypePro')"></el-input>
            </el-form-item>
          </el-col>
          <el-col :span="8"
                  class="service-param">
            <el-form-item :prop="`parameters[${index}].value`"
                          label-width="0px">
              <el-input v-model="item.value"
                        :placeholder="$t('serviceList.paramValuePro')"></el-input>
            </el-form-item>
          </el-col>
          <el-col :span="2"
                  class="subtitle">
            <i class="el-icon-remove-outline delete"
               v-if="!viewDetail"
               @click="deleteParam(index)" />
          </el-col>
        </el-row>
      </div>
    </el-form>
  </div>
</template>
<script>
export default {
  props: {
    formData: {
      type: Object,
      default () {
        return {}
      }
    },
    groupList: {
      type: Array,
      default () {
        return []
      }
    },
    modify: {
      type: Boolean,
      default: false
    },
    viewDetail: {
      type: Boolean,
      default: false
    }
  },
  data () {
    return {
      modelOptionList: [],
      imageOptionList: [],
      versionOptionList: [],
      form: {
        group_id: '',
        model_id: '',
        modelversion_id: '',
        model_type: '',
        endpoint_type: '',
        image: '',
        parameters: [{
          name: '',
          type: '',
          value: ''
        }]
      }
    }
  },
  computed: {
    ruleValidate () {
      return {
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
        ]
      }
    }
  },
  created () {
    if (this.modify || this.viewDetail) {
      Object.assign(this.form, this.formData)
      console.log('this.form', this.form)
      // this.getModelOption(this.formData.group_id)
      this.getVersionOption(this.formData.model_id)
      if (this.formData.model_type === 'CUSTOM') {
        this.getImageOption(this.formData.modelversion_id)
      }
    }
  },
  methods: {
    changeGroup (val) {
      console.log('1111')
      this.form.model_id = ''
      this.form.modelversion_id = ''
      this.versionOptionList = []
      this.getModelOption(val)
    },
    getModelOption (val) {
      this.FesApi.fetch(`/mf/v1/modelsByGroupID/${val}`, { page: 1, size: 1000 }, 'get').then((res) => {
        this.modelOptionList = this.formatResList(res)
      })
    },
    changeVersionId (versionId) {
      if (this.form.model_type === 'CUSTOM') {
        this.form.image = ''
        this.getImageOption(versionId)
      }
    },
    getImageOption (versionId) {
      this.FesApi.fetch(`/mf/v1/images/${versionId}`, {}, 'get').then((res) => {
        this.imageOptionList = this.formatResList(res)
      })
    },
    selectChangeModel (val) {
      this.form.modelversion_id = ''
      this.imageOptionList = []
      for (let item of this.modelOptionList) {
        if (item.id === val) {
          this.form.model_type = item.model_type
          this.form.group_name = item.group_name
          this.form.model_name = item.model_name
          break
        }
      }
      this.getVersionOption(val)
    },
    getVersionOption (val) {
      this.FesApi.fetch(`/mf/v1/modelversions/${val}`, { page: 1, size: 1000 }, 'get').then((res) => {
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
    byVersionIdFindSource () {
      let modelVersion = ''
      let source = ''
      for (let item of this.versionOptionList) {
        if (item.id === this.form.modelversion_id) {
          modelVersion = item.name
          source = item.source
          break
        }
      }
      this.form.modelVersion = modelVersion
      this.form.source = source
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
.padding-label {
  padding-left: 135px;
}
.delete {
  line-height: 32px;
}
</style>
