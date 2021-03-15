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
    // hadoop 函数
    hdpTaskBaseObj () {
      return {
        codeSettings: 'codeFile',
        HDFSPath: '',
        fileName: ''
      }
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
          name: 'uat.sf.dockerhub.stgwebank/wedatasphere/prophecis',
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
    },
    validateFormChangeTab (keyArr, errObj, activeName) { // keyArr tab对应表单下所有字段，errObj 校验不通过返回的不通过字段，activeName 当前tab index
      let tabActiveArr = ['0', '1', '2'] // tab index总数组成数组
      let currentTab = false // 判断是否当前tab有字段未填
      const _this = this
      let isEmit = false // 是否切换过tab
      judgmentTba(keyArr[activeName], activeName)
      if (!currentTab) {
        tabActiveArr.splice(activeName, 1) // 如果不是当前tab表单字段不通过则循环其他tab字段
        for (let index of tabActiveArr) {
          if (!isEmit) {
            judgmentTba(keyArr[index], index)
          }
        }
      }
      function judgmentTba (formLabel, index) { // formLabel-> tab下表单字段; index-> formLabel对应哪个tab index
        for (let item of formLabel) {
          if (errObj[item]) {
            if (index !== activeName) {
              _this.$emit('changeTab', index)
              isEmit = true
            } else {
              currentTab = true
            }
            return
          }
        }
      }
    },
    mfflowBasicData () {
      return {
        method: '/api/rest_j/v1/entrance/execute',
        params: {
          ariable: {
            k1: 'v1'
          },
          configuration: {
            special: {
              k2: 'v2'
            },
            runtime: {
              'k3': 'v3'
            },
            startup: {
              k4: 'v4'
            }
          }
        },
        executeApplicationName: 'spark',
        runType: 'sql',
        source: {
          scriptPath: '/home/Linkis/Linkis.sql'
        }
      }
    },
    handleOriginNodeData (val, initForm) {
      this.currentNode = val
      if (this.currentNode.jobContent && this.currentNode.jobContent.ManiFest) {
        const form = this.currentNode.jobContent.ManiFest
        // let newFormData = this.basicInfoProcess(form)
        // this.handelTfosCopy(newFormData)
        // Object.assign(this.form, newFormData)
        if (form.name) {
          if (form.job_alert) {
            this.alarmData = form.job_alert
          }
          switch (this.currentNode.type) {
            case 'linkis.appjoint.mlflow.gpu':
              this.handleGpuNodeData(form)
              break
            case 'linkis.appjoint.mlflow.hadoop':
              this.handleHadoopNodeData(form)
              break
          }
        }
      } else {
        this.form = JSON.parse(JSON.stringify(initForm))
        this.form.name = this.currentNode.title
        setTimeout(() => {
          this.$refs.formValidate.clearValidate()
        }, 10)
      }
    }
  }
}
