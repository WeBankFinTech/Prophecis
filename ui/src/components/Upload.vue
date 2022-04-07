<template>
  <div key="codeFile"
       class="upload-box">
    <span v-if="fileName"
          class="file-name"
          :style="'margin-left:' +marginLeft+'px;'">
      {{ fileName }}
      <i class="el-icon-close delete-file"
         v-if="!readonly"
         @click="deleteFile()" />
    </span>
    <el-upload v-if="baseUrl"
               ref="upload"
               :action="uploadUrl"
               :disabled="loading || readonly"
               :headers="header"
               :on-success="uploadFile"
               :data="param"
               class="upload"
               :before-upload="beforeUpload"
               :on-error="errrUpload">
      <el-button type="primary"
                 :loading="loading">
        {{ $t('DI.manualUpload') }}
      </el-button>
    </el-upload>
    <el-upload v-else
               action="#"
               :disabled="loading|| readonly"
               :http-request="uploadFile"
               class="upload"
               :before-upload="beforeUpload"
               :on-error="errrUpload">
      <el-button type="primary"
                 :loading="loading">
        {{ $t('DI.manualUpload') }}
      </el-button>
    </el-upload>
  </div>
</template>
<script>
export default {
  props: {
    baseUrl: {
      type: String,
      default: ''
    },
    fileKey: {
      type: String,
      default: ''
    },
    marginLeft: {
      type: Number,
      default: 178
    },
    codeFile: {
      type: String,
      default: ''
    },
    fileSize: {
      type: Number,
      default: 512
    },
    dssHeader: {
      type: Object,
      default () {
        return {}
      }
    },
    allType: {
      type: Boolean,
      default: false
    },
    transferParam: {
      type: Boolean,
      default: false
    },
    params: {
      type: Object,
      default: function () {
        return {}
      }
    },
    readonly: {
      type: Boolean,
      default: false
    }
  },
  data () {
    return {
      fileName: '',
      loading: false,
      uploadUrl: this.baseUrl,
      header: {
        'Mlss-Userid': localStorage.getItem('userId')
      },
      param: {}
    }
  },
  watch: {
    params (val) {
      if (val.modelType) {
        this.fileName = ''
        this.$emit('deleteFile', this.fileKey)
      }
    }
  },
  mounted () {
    this.fileName = this.codeFile
    if (this.dssHeader['MLSS-AppID']) {
      Object.assign(this.header, this.dssHeader)
      delete this.header['Content-Type']
    }
  },
  methods: {
    deleteFile () {
      this.fileName = ''
      this.$emit('deleteFile', this.fileKey)
    },
    uploadFile (response, files) {
      this.loading = false
      if (response) {
        let file = ''
        if (this.baseUrl) {
          file = response
          file.fileName = files.name
          this.fileName = files.name
        } else {
          this.fileName = response.file.name
          file = response.file
        }
        this.$emit('uploadFile', file, this.fileKey)
      }
    },
    beforeUpload (file) {
      let fileTpye = true
      let fileCheck = false
      if (!this.allType) {
        fileTpye = file.type.indexOf('zip') > -1
        if (!fileTpye) {
          this.$message.error(this.$t('common.fileFormat'))
        }
        fileCheck = fileTpye
      } else {
        fileCheck = true
      }
      const isLt512M = file.size / 1024 / 1024 < this.fileSize
      let limit = this.fileName === ''
      if (!isLt512M) {
        this.$message.error(this.$t('common.fileSize' + this.fileSize + 'BM'))
      }
      if (!limit) {
        this.$message.error(this.$t('common.fileLimit'))
        return false
      }
      if (fileCheck && isLt512M && isLt512M) {
        if (this.transferParam) { // 是否需要传文件名
          this.param.fileName = file.name
          Object.assign(this.param, this.params)
        }
        this.loading = true
      }
      return fileCheck && isLt512M
    },
    errrUpload (error) {
      const msg = JSON.parse(error.message)
      this.$message.error(msg)
      this.loading = false
    }
  }
}
</script>
<style lang="scss" scoped>
.upload-box {
  padding: 15px 0 20px 45px;
  height: 50px;
  text-align: right;
  padding: 0px;
  .upload {
    display: inline;
  }
  .file-name {
    float: left;
    padding: 10px 30px;
    font-size: 14px;
    background: #f5f7fa;
    border-radius: 5px;
  }
  .delete-file {
    position: absolute;
    right: 0;
    top: 0;
    cursor: pointer;
  }
  .el-upload-list {
    display: none;
  }
}
</style>
