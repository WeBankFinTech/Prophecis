<template>
  <el-dialog :title="modify?$t('modelList.modifyModel'):$t('modelList.createModel')"
             :visible.sync="dialogVisible"
             @close="clearDialogForm"
             custom-class="dialog-style"
             width="1100px">
    <el-form ref="formValidate"
             :model="form"
             :rules="formValidate"
             label-width="155px">
      <div class="subtitle">
        {{$t('modelList.basicModelInfo')}}
      </div>
      <el-form-item :label="$t('modelList.modelName')"
                    prop="model_name">
        <el-input v-model="form.model_name"
                  :disabled="modify"
                  :placeholder="$t('modelList.modelNamePro')"></el-input>
      </el-form-item>
      <el-row>
        <el-col :span="12">
          <el-form-item :label="$t('modelList.belongUserGroup')"
                        prop="group_id">
            <el-select v-model="form.group_id"
                       filterable
                       :placeholder="$t('modelList.belongUserGroupPro')">
              <el-option v-for="(item,index) of groupList"
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
        </el-col>
      </el-row>
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
      <div v-if="form.modelLocation==='codeFile'&&dialogVisible"
           key="codeFile"
           class="tofs-upload">
        <upload v-if="form.model_type"
                :margin-left="155"
                :baseUrl="baseUrl"
                :file-size="2048"
                :params="params"
                :allType="true"
                :transferParam="true"
                :codeFile='file_name'
                @deleteFile="deleteFile"
                @uploadFile="uploadFile"></upload>
      </div>
    </el-form>
    <span slot="footer"
          class="dialog-footer">
      <el-button type="primary"
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
import Upload from '../../../components/Upload'
export default {
  components: {
    Upload
  },
  props: {
    modelId: {
      type: Number,
      default: 0
    },
    localPathList: {
      type: Array,
      default () {
        return []
      }
    },
    groupList: {
      type: Array,
      default () {
        return []
      }
    }
  },
  data () {
    return {
      dialogVisible: false,
      form: {
        model_name: '',
        group_id: '',
        model_type: '',
        modelLocation: 'codeFile',
        root_path: '',
        child_path: ''
      },
      params: {
      }, // 上传文件参数
      modify: false,
      changeFile: false,
      source: '',
      file_name: '',
      baseUrl: process.env.VUE_APP_BASE_SURL + '/mf/v1/model/uploadModel'
    }
  },
  computed: {
    formValidate () {
      this.resetFormFields()
      return {
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
  methods: {
    subInfo () {
      this.$refs.formValidate.validate((valid) => {
        if (valid) {
          let param = { ...this.form }
          param.group_id = parseInt(param.group_id)
          delete param.modelLocation
          if (this.form.modelLocation === 'codeFile') {
            if (!this.source) {
              this.$message(this.$t('modelList.uploadFileReq'))
              return
            } else {
              // if (!this.modify || (this.modify && this.changeFile)) {
              param.s3_path = this.source
              param.file_name = this.file_name
              // }
              delete param.root_path
              delete param.child_path
            }
          }
          let method = 'post'
          let url = '/mf/v1/model'
          if (this.modify) {
            method = 'put'
            url += '/' + this.modelId
          }
          this.FesApi.fetch(url, param, method).then((res) => {
            this.toast()
            this.$emit('getListData')
            this.dialogVisible = false
          })
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
      if (this.modify) {
        this.changeFile = true
      }
      this.source = file.result.s3Path
      this.file_name = file.fileName
    },
    clearDialogForm () {
      this.deleteFile()
      this.modify = false
      this.changeFile = false
      for (let key in this.form) {
        console.log('key', key)
        this.form[key] = ''
      }
      this.deleteFile()
      setTimeout(() => {
        this.$refs.formValidate.resetFields()
        this.form.modelLocation = 'codeFile'
      }, 0)
    }
  }
}
</script>
<style lang="scss" scoped>
.el-form {
  height: 320px;
}
</style>
