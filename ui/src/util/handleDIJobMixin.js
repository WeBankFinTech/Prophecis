import util from './common'
export default {
  methods: {
    // gpu函数
    gpuBeforeUpload (file) {
      const fileTpye = file.type.indexOf('zip') > -1
      const isLt512M = file.size / 1024 / 1024 < 512
      const limit = this.files === ''
      if (!fileTpye) {
        this.$message.error(this.$t('common.fileFormat'))
      }
      if (!isLt512M) {
        this.$message.error(this.$t('common.fileSize'))
      }
      if (!limit) {
        this.$message.error(this.$t('common.fileLimit'))
        return false
      }
      return fileTpye && isLt512M
    },
    gpuDeleteFile () {
      this.files = ''
      this.fileName = ''
    },
    gpuUploadFile (files) {
      if (files.file) {
        this.fileName = files.file.name
        this.files = files.file
      }
    },
    handleGPUSubParam () {
      let param = {
        name: this.form.name,
        description: this.form.description,
        version: '1.0',
        gpus: +this.form.gpus,
        cpus: +this.form.cpus,
        memory: this.form.memory + 'Gb',
        namespace: this.form.namespace,
        code_selector: this.form.codeSettings,
        job_type: this.form.job_type,
        data_stores: [{
          id: 'hostmount',
          type: 'mount_volume',
          training_data: {
            container: this.form.trainingData
          },
          training_results: {
            container: this.form.trainingResults
          },
          connection: {
            type: 'host_mount',
            name: 'host-mount',
            path: this.form.path
          }
        }],
        framework: {
          name: this.FesEnv.DI.defineImage,
          version: this.form.imageType === 'Standard' ? this.form.imageOption : this.form.imageInput,
          command: this.form.command
        },
        evaluation_metrics: {
          type: 'fluent-bit'
        }
      }
      return param
    },
    hdpAddExecutorsTask (objKey) {
      const baseObj = this.hdpTaskBaseObj()
      const leng = this.form[objKey].length
      if (leng >= 1) {
        this.hdpAddExecutorsValidate(objKey, leng)
      }
      this.form[objKey].push(baseObj)
    },
    hdpDeleteExecutorsTask (objKey) {
      const leng = this.form[objKey].length
      this.form[objKey].splice(leng - 1, 1)
      if (leng > 1) {
        this.hdpDeleteExecutorsValidate(objKey, leng)
      }
    },
    hdpAddExecutorsValidate (objKey, index) {
      this.ruleValidate[`${objKey}[${index}].HDFSPath`] = this.ruleValidate[`${objKey}[0].HDFSPath`]
    },
    hdpDeleteExecutorsValidate (objKey, index) {
      delete this.ruleValidate[`${objKey}[${index - 1}].HDFSPath`]
    },
    // changeCodeSettings (objKey, index) {
    //   let modules = {}
    //   modules = this.form[objKey][index]
    //   if (this.form[objKey].codeSettings === 'codeFile') {
    //     modules.HDFSPath = ''
    //   } else {
    //     modules.fileName = ''
    //   }
    // },
    hdpSubmitUpload (objKey, index) {
      const key = `${objKey}Upload${index}`
      this.$refs[key][0].submit()
      this.currentUplod.uploadModule = objKey
      this.currentUplod.index = index
    },
    hdpBeforeUpload (file) {
      const fileTpye = file.type.indexOf('zip') > -1
      const isLt512M = file.size / 1024 / 1024 < 512
      const currentUplod = this.currentUplod.uploadModule
      let limit = ''
      const currentIndex = this.currentUplod.index
      limit = this.form[currentUplod][currentIndex].fileName === ''
      if (!fileTpye) {
        this.$message.error(this.$t('common.fileFormat'))
      }
      if (!isLt512M) {
        this.$message.error(this.$t('common.fileSize'))
      }
      if (!limit) {
        this.$message.error(this.$t('common.fileLimit'))
        return false
      }
      return fileTpye && isLt512M
    },
    hdpDeleteFile (objKey, index) {
      this.form[objKey][index].fileName = ''
    },
    hadUploadFile (response, files) {
      if (response.data) {
        const currentUplod = this.currentUplod.uploadModule
        const currentIndex = this.currentUplod.index
        this.form[currentUplod][currentIndex].fileName = files.raw.name
      }
    },
    handleHdpSubParam () {
      let param = {
        version: '1.0',
        job_type: 'tfos',
        code_selector: 'storagePath',
        framework: {
          name: 'uat.sf.dockerhub.stgwebank/webank/mlss-di',
          version: 'tfosexecutor-1.6.0'
        },
        evaluation_metrics: {
          type: 'fluent-bit'
        }
      }
      const arr = ['name', 'namespace', 'description']
      for (let item of arr) {
        param[item] = this.form[item]
      }
      this.handelHdpTfos(param)
      const tfosArr = ['queue', 'executor_cores', 'executor_memory', 'driver_memory', 'executors', 'command']
      for (let key of tfosArr) {
        param.tfos_request[key] = this.form[key]
      }
      return param
    },
    handelHdpTfos (param) {
      let arr = ['pyfile', 'archives']
      param.tfos_request = {
        pyfile: [],
        archives: []
      }
      const time = util.getTime()
      const definePython = this.FesEnv.DI.definePython
      for (let key of arr) {
        for (let item of this.form[key]) {
          if (item.codeSettings === 'codeFile' && item.fileName) {
            let hdfs = `${definePython}/${this.userId}/${time}/${item.fileName}`
            param.tfos_request[key].push({ hdfs: hdfs })
          } else if (item.codeSettings === 'HDFSPath' && item.HDFSPath) {
            param.tfos_request[key].push({ hdfs: item.HDFSPath })
          }
        }
      }
      let envHDFSPath = this.form.pythonType === 'Standard' ? this.form.pythonOption : this.form.pythonInput
      param.tfos_request.tensorflow_env = { hdfs: envHDFSPath }
    },
    validateForm () {
      let isValid = false
      this.$refs.formValidate.validate((valid) => {
        isValid = valid
      })
      return isValid
    }
  }
}
