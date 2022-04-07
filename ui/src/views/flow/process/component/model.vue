<template>
  <div>
    <el-form ref="formValidate"
             :model="form"
             :disabled="readonly"
             :rules="ruleValidate"
             label-position="top"
             class="node-parameter-bar form-height">
      <div class="subtitle">
        {{$t('modelList.basicModelInfo')}}
      </div>
      <el-form-item :label="$t('flow.nodeName')"
                    prop="name">
        <el-input v-model="form.name"></el-input>
      </el-form-item>
      <el-form-item :label="$t('modelList.modelName')"
                    prop="model_name">
        <el-input v-model="form.model_name"
                  :placeholder="$t('modelList.modelNamePro')"></el-input>
      </el-form-item>
      <el-form-item :label="$t('modelList.belongUserGroup')"
                    prop="group_id">
        <el-select v-model="form.group_id"
                   :filterable="true"
                   :placeholder="$t('modelList.belongUserGroupPro')">
          <el-option v-for="(item,index) of groupList"
                     :label="item.name"
                     :value="item.id"
                     :key="index">
          </el-option>
        </el-select>
      </el-form-item>
      <el-form-item :label="$t('serviceList.modelType')"
                    prop="model_type">
        <el-select v-model="form.model_type"
                   @change="changeModelType"
                   :disabled="modify"
                   :placeholder="$t('serviceList.modelTypePro')">
          <el-option label="Tensorflow"
                     value="TENSORFLOW">
          </el-option>
          <el-option label="SKLearn"
                     value="SKLEARN">
          </el-option>
          <!-- <el-option label="Pytorch"
                         value="PYTORCH">
              </el-option> -->
          <el-option label="XGBoost"
                     value="XGBOOST">
          </el-option>
          <el-option :label="$t('modelList.custom')"
                     value="CUSTOM">
          </el-option>
          <!-- <el-option :label="$t('serviceList.other')"
                         value="OTHER">
              </el-option> -->
        </el-select>
      </el-form-item>
      <div class="subtitle">
        {{$t('modelList.modelLocationinfo')}}
      </div>
      <el-form-item :label="$t('serviceList.modelLocation')"
                    prop="modelLocation">
        <el-select v-model="form.modelLocation"
                   :placeholder="$t('DI.uploadModeSettingsPro')">
          <el-option :label="$t('DI.manualUpload')"
                     value="codeFile">
          </el-option>
          <el-option :label="$t('serviceList.shareDirectory')"
                     value="HDFSPath">
          </el-option>
        </el-select>
      </el-form-item>
      <div v-if="form.modelLocation==='HDFSPath'"
           key="HDFSPath">
        <el-form-item :label="$t('modelList.storageRootPath')"
                      prop="root_path">
          <el-select v-model="form.root_path"
                     filterable
                     :placeholder="$t('modelList.storageRootPathPro')">
            <el-option v-for="(item,index) in localPathList"
                       :label="item"
                       :key="index"
                       :value="item">
            </el-option>
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('modelList.fileStorageSubpath')"
                      prop="child_path">
          <el-input v-model="form.child_path"
                    :placeholder="$t('modelList.fileStorageSubpathPro')" />
        </el-form-item>
      </div>
      <div v-if="form.modelLocation==='codeFile'"
           key="codeFile"
           class="tofs-upload">
        <upload v-if="form.model_type"
                :margin-left="0"
                :baseUrl="baseUrl"
                :file-size="2048"
                :params="params"
                :dssHeader="dssHeader"
                :allType="true"
                :transferParam="true"
                :codeFile='file_name'
                :readonly="readonly"
                @deleteFile="deleteFile"
                @uploadFile="uploadFile"></upload>
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
import Upload from '../../../../components/Upload'
import modelMixin from '../../../../util/modelMixin'
import nodeFormMixin from './nodeFormMixin'
const initForm = {
  name: '',
  model_name: '',
  group_id: '',
  model_type: '',
  modelLocation: 'codeFile',
  root_path: '',
  child_path: ''
}
export default {
  mixins: [modelMixin, nodeFormMixin],
  components: {
    Upload
  },
  props: {
    nodeData: {
      type: Object,
      default () {
        return {}
      }
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
      localPathList: [],
      form: { ...initForm },
      params: {
      }, // 上传文件参数
      dssHeader: {},
      modify: false,
      source: '',
      file_name: '',
      baseUrl: process.env.VUE_APP_BASE_SURL + '/mf/v1/model/uploadModel'
    }
  },
  created () {
    if (!this.$route.query.exp_id) {
      this.dssHeader = this.getSystemAuthHeader()
      const header = this.handlerFromDssHeader('get')
      this.getDIStoragePath(header)
      this.getGroupList(header)
    } else {
      this.getDIStoragePath()
      this.getGroupList()
    }
  },
  computed: {
    ruleValidate () {
      return {
        name: [
          { required: true, message: this.$t('flow.nodeNoEmpty') },
          { type: 'string', pattern: new RegExp(/^[0-9a-zA-Z-_]*$/), message: this.$t('flow.nodeNameTip') }
        ],
        model_name: [
          { required: true, message: this.$t('modelList.modelNameReq') },
          { type: 'string', pattern: new RegExp(/^[a-zaA-Z][a-zA-Z0-9-]*$/), message: this.$t('modelList.modelNameFormat') }
        ],
        group_id: [
          { required: true, message: this.$t('modelList.belongUserGroupReq') }
        ],
        model_type: [
          { required: true, message: this.$t('serviceList.modelTypeReq') }
        ],
        modelLocation: [
          { required: true, message: this.$t('serviceList.deploymentTypeReq') }
        ],
        root_path: [
          { required: true, message: this.$t('modelList.storageRootPathReq') }
        ],
        child_path: [
          { required: true, message: this.$t('modelList.fileStorageSubpathReq') },
          { type: 'string', pattern: new RegExp(/^\/[a-zA-Z][0-9a-zA-Z-._/]*$/), message: this.$t('modelList.fileStorageSubpathFormat') }
        ]
      }
    }
  },
  watch: {
    nodeData: {
      immediate: true,
      handler (val) {
        this.modify = false
        this.deleteFile()
        this.handleOriginNodeData(val, initForm)
      }
    }
  },
  methods: {
    save (evt, temSave = false) {
      this.$refs.formValidate.validate((valid) => {
        if (valid) {
          let param = { ...this.form }
          param.group_id = parseInt(param.group_id)
          if (this.form.modelLocation === 'codeFile') {
            if (temSave) {
              return
            }
            if (!this.source) {
              this.$message(this.$t('modelList.uploadFileReq'))
              return
            } else {
              param.s3_path = this.source
              param.file_name = this.file_name
              delete param.root_path
              delete param.child_path
            }
          }
          this.handleCurrentNode(param)
          this.$emit('saveNode', this.currentNode, temSave)
        }
      })
    },
    changeModelType (val) {
      this.params.modelType = val
    },
    deleteFile () {
      this.source = ''
      this.file_name = ``
    },
    uploadFile (file, key) {
      this.source = file.result.s3Path
      this.file_name = file.fileName
    },
    handleModelNodeData (form) {
      this.modify = true
      this.source = form.s3_path
      this.file_name = form.file_name
      this.params.modelType = form.model_type
      this.form = { ...form }
    }
  }
}
</script>
<style lang="scss" scoped>
</style>
