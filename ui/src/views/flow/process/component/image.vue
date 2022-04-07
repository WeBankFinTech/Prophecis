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
          {{ $t('DI.basicSettings') }}
        </div>
        <el-form-item :label="$t('flow.nodeName')"
                      prop="name">
          <el-input v-model="form.name"></el-input>
        </el-form-item>
        <el-form-item :label="$t('image.imageBase')"
                      prop="image"
                      maxlength="225">
          <el-input v-model="form.image"
                    :placeholder="$t('image.imageBasePro')" />
        </el-form-item>
        <el-form-item :label="$t('image.imageTag')"
                      prop="image_tag">
          <el-input v-model="form.image_tag"
                    :placeholder="$t('image.imageTagPro')" />
        </el-form-item>
        <el-form-item :label="$t('image.imageRemark')">
          <el-input v-model="form.remarks"
                    type="textarea"
                    :placeholder="$t('image.imageRemarkPro')" />
        </el-form-item>
      </div>
      <div v-show="activeName==='1'"
           key="step1">
        <div class="subtitle">
          {{ $t('serviceList.modelSettings') }}
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
                        prop="model_name">
            <el-select v-model="form.model_name"
                       :placeholder="$t('serviceList.modelNamePro1')"
                       @change="selectChangeModel">
              <el-option v-for="(item,index) of modelOptionList"
                         :label="item.model_name"
                         :value="item.id"
                         :key="index">
              </el-option>
            </el-select>
          </el-form-item>
          <div v-if="form.model_name">
            <el-form-item :label="$t('serviceList.modelVersion')"
                          prop="model_version_id">
              <el-select v-model="form.model_version_id"
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
                      placeholder=" "></el-input>
          </el-form-item>
        </div>
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
const initForm = {
  name: '',
  image: '',
  image_tag: '',
  remarks: '',
  group_id: '',
  model_name: '',
  model_version_id: '',
  model_type: '',
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
      groupList: [],
      modelOptionList: [],
      versionOptionList: [],
      form: { ...initForm },
      fromDss: false
    }
  },
  created () {
    const header = this.handlerFromDssHeader('get')
    this.getGroupList(header)
  },
  computed: {
    ruleValidate () {
      return {
        name: [
          { required: true, message: this.$t('flow.nodeNoEmpty') },
          { type: 'string', pattern: new RegExp(/^[0-9a-zA-Z-_]*$/), message: this.$t('flow.nodeNameTip') }
        ],
        image: [
          { required: true, message: this.$t('image.imageBaseReq') },
          { type: 'string', pattern: new RegExp(/^[a-z][a-z0-9-_./]*$/), message: this.$t('image.imageBaseFormat') }
        ],
        image_tag: [
          { required: true, message: this.$t('image.imageTagReq') },
          { type: 'string', pattern: new RegExp(/^[a-z0-9A-Z][a-z0-9A-Z-_./]*$/), message: this.$t('image.imageTagFormat') }
        ],
        group_id: [
          { required: true, message: this.$t('user.userGroupSettingsReq') }
        ],
        model_name: [
          { required: true, message: this.$t('modelList.modelNameReq') }
        ],
        model_version_id: [
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
  watch: {
    nodeData: {
      immediate: true,
      handler (val) {
        this.handleOriginNodeData(val, initForm)
      }
    }
  },
  methods: {
    handleImageNodeData (form) {
      if (form.param_type === 'model') {
        this.getVersionList(form.model_name)
        form.model_version_id = form.model_version_id + ''
        form.model_name = parseInt(form.model_name)
      } else {
        form.cutom_model_name = form.model_name
        form.cutom_version = form.model_version
        form.model_version_id = ''
        delete form.model_name
      }
      this.getModelList(form.group_id)
      this.form = { ...form }
    },
    judgmentFormTab () {
      return {
        '0': ['name', 'image', 'image_tag'],
        '1': ['param_type', 'group_id', 'model_name', 'model_version_id', 'model_type', 'cutom_model_name']
      }
    },
    save (evt, temSave = false) {
      this.$refs.formValidate.validate((valid, errObj) => {
        if (valid) {
          let param = {}
          const keyArr = ['name', 'image', 'image_tag', 'remarks', 'group_id', 'param_type']
          for (let key of keyArr) {
            param[key] = this.form[key]
          }
          if (this.form.param_type === 'model') {
            param.model_name = this.form.model_name
            param.model_type = this.form.model_type
            param.model_version_id = parseInt(this.form.model_version_id)
          } else {
            param.model_name = this.form.cutom_model_name
            param.model_version = this.form.cutom_version
          }
          param.image_name = this.form.image + ':' + this.form.image_tag
          param.user_name = localStorage.getItem('userId')
          this.handleCurrentNode(param)
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
    changeGroup (val) {
      this.form.model_name = ''
      this.form.model_version_id = ''
      this.versionOptionList = []
      this.getModelList(val)
    },
    getModelList (val) {
      const header = this.handlerFromDssHeader('get')
      this.FesApi.fetch(`/mf/v1/modelsByGroupID/${val}`, { page: 1, size: 1000 }, header).then((res) => {
        let modelOption = []
        for (let item of res.models) {
          if (item.model_type === 'CUSTOM') {
            modelOption.push(item)
          }
        }
        this.modelOptionList = this.formatResList(modelOption)
      })
    },
    selectChangeModel (val) {
      this.form.model_version_id = ''
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
      this.getVersionList(val)
    },
    getVersionList (val) {
      const header = this.handlerFromDssHeader('get')
      this.FesApi.fetch(`/mf/v1/modelversions/${val}`, { page: 1, size: 1000 }, header).then((res) => {
        let optionList = []
        if (res.modelVersions) {
          for (let item of res.modelVersions) {
            optionList.push({
              id: item.id + '',
              name: item.version,
              source: item.source
            })
          }
        }
        this.versionOptionList = Object.freeze(optionList)
      })
    }
  }
}
</script>
