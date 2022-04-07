export default {
  methods: {
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
      if (this.currentNode.jobContent && (this.currentNode.jobContent.ManiFest || this.currentNode.jobContent.content)) {
        const form = this.currentNode.jobContent.ManiFest || this.currentNode.jobContent.content
        if (form.name) { // 由于算法节点没有节点名称，如果cpus存在则证明存在节点信息
          if (form.job_alert) {
            this.alarmData = form.job_alert
          }
          switch (this.currentNode.type) {
            case 'linkis.appconn.mlflow.gpu':
              this.handleGpuNodeData(form)
              break
            case 'linkis.appconn.mlflow.model':
              this.handleModelNodeData(form)
              break
            case 'linkis.appconn.mlflow.image':
              this.handleImageNodeData(form)
              break
            case 'linkis.appconn.mlflow.service':
              this.handleServiceNodeData(form)
              break
            case 'linkis.appconn.mlflow.reportPush':
              this.handleReportPushNodeData(form)
              break
            default:
              this.handleAlgorithmNodeData(form, initForm)
              break
          }
        }
      } else {
        const algorithmType = ['linkis.appconn.mlflow.LogisticRegression', 'linkis.appconn.mlflow.DecisionTree', 'linkis.appconn.mlflow.RandomForest', 'linkis.appconn.mlflow.XGBoost', 'linkis.appconn.mlflow.LightGBM']
        this.form = JSON.parse(JSON.stringify(initForm))
        if (algorithmType.indexOf(this.currentNode.type) > -1) {
          this.$nextTick(() => {
            this.$refs.dataSetting.initFormData(this.currentNode)
            this.$refs.customModel.initFormData()
            this.$refs.basicResource.initFormData()
          })
        } else {
          this.form.name = this.currentNode.title
        }
        setTimeout(() => {
          this.$refs.formValidate.clearValidate()
        }, 10)
      }
    },
    handleCurrentNode (param) {
      const mlflowJobType = {
        'linkis.appconn.mlflow.model': 'ModelStorage',
        'linkis.appconn.mlflow.image': 'ImageBuild',
        'linkis.appconn.mlflow.service': 'ModelDeploy',
        'linkis.appconn.mlflow.reportPush': 'ReportPush'
      }
      const maniFestNode = ['linkis.appconn.mlflow.gpu', 'linkis.appconn.mlflow.LogisticRegression', 'linkis.appconn.mlflow.DecisionTree', 'linkis.appconn.mlflow.RandomForest', 'linkis.appconn.mlflow.XGBoost', 'linkis.appconn.mlflow.LightGBM']
      this.currentNode.jobContent = this.mfflowBasicData()
      if (maniFestNode.indexOf(this.currentNode.type) > -1) {
        this.currentNode.jobContent.ManiFest = param
      } else if (['linkis.appconn.mlflow.model', 'linkis.appconn.mlflow.image', 'linkis.appconn.mlflow.service', 'linkis.appconn.mlflow.reportPush'].indexOf(this.currentNode.type) > -1) {
        this.currentNode.jobContent.content = param
      }
      if (mlflowJobType[this.currentNode.type]) {
        this.currentNode.jobContent.mlflowJobType = mlflowJobType[this.currentNode.type]
      } else {
        this.currentNode.jobContent.mlflowJobType = 'DistributedModel'
      }
    },
    // 系统配置header
    getSystemAuthHeader () {
      return {
        'Content-Type': 'application/json',
        'X-Watson-Userinfo': 'bluemix-instance-id=test-user',
        'MLSS-AppTimestamp': '1603274993343',
        'MLSS-AppID': 'MLFLOW',
        'MLSS-Auth-Type': 'SYSTEM',
        'MLSS-APPSignature': '91a585402398b4c3f63ab26b30722cfc2a74e3b3'
      }
    },
    // 如果从dss跳转过来改成系统鉴权，需要增加header
    handlerFromDssHeader (method) {
      let header = {
        method: method
      }
      if (!this.$route.query.exp_id) {
        header.headers = this.getSystemAuthHeader()
      }
      return header
    },
    validateFormChangeTab (keyArr, errObj, activeName) { // keyArr tab对应表单下所有字段，errObj 校验不通过返回的不通过字段，activeName 当前tab index
      let tabActiveArr = [] // tab index总数组成数组
      for (let key in keyArr) {
        tabActiveArr.push(key)
      }
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
    // 处理算法节点表单参数
    getAlgorithmParam (initForm) {
      return new Promise(resolve => {
        this.$refs.formValidate.validate(valid => {
          if (valid) {
            let param = {}
            // 由于参数都是非必填，如果校验type类型为number int、float如果为空会报错，type为string提交前在格式化
            const formatKey = this.algorithmDataFormatKey()
            // const maniFest = this.currentNode.jobContent ? this.currentNode.jobContent.ManiFest : {}
            // const jobParams = maniFest.name ? JSON.parse(maniFest.job_params) : ''
            for (let key in this.form) {
              if (this.form[key] === '') {
                continue
              }
              this.formatStringNumber(formatKey, param, key)
            }
            // param.version = this.form.version
            resolve(param)
          } else {
            resolve(false)
          }
        })
      })
    },
    // 将string格式化成number
    formatStringNumber (formatKey, param, key) {
      if (formatKey.float.indexOf(key) > -1) {
        param[key] = isNaN(parseFloat(this.form[key])) ? '' : parseFloat(this.form[key])
      } else if (formatKey.int.indexOf(key) > -1) {
        console.log('isNaN(parseInt(this.form[key]))', isNaN(parseInt(this.form[key])))
        param[key] = isNaN(parseInt(this.form[key])) ? '' : parseInt(this.form[key])
      } else {
        param[key] = this.form[key]
      }
    },
    // 算法节点提交校验
    async algorithmSubmit (initForm, algorithmType, temSave) {
      let dataSetting = await this.$refs.dataSetting.save()
      if (dataSetting) {
        let jobParmas = await this.getAlgorithmParam(initForm)
        let customModel = await this.$refs.customModel.save()
        if (jobParmas && customModel) {
          let basicResource = await this.$refs.basicResource.save()
          if (basicResource) {
            let param = { ...dataSetting, ...basicResource }
            param.mf_model = { ...customModel }
            const fitKey = ['early_stopping_rounds', 'sample_weight_eval_set', 'eval_metric', 'sample_weight']
            let fitParams = {}
            for (let key of fitKey) {
              if (jobParmas[key] !== undefined) {
                fitParams[key] = jobParmas[key]
                delete jobParmas[key]
              }
            }
            const APIType = jobParmas.API_type
            if (APIType) {
              delete jobParmas.API_type
            }
            const paramType = jobParmas.param_type
            if (this.currentNode.type === 'linkis.appconn.mlflow.XGBoost' || this.currentNode.type === 'linkis.appconn.mlflow.LightGBM') {
              debugger
              delete jobParmas.param_type
              if (paramType === 'custom') {
                param.job_params = jobParmas.custom_param
              } else {
                delete jobParmas.custom_param
                param.job_params = JSON.stringify(jobParmas)
              }
            } else {
              param.job_params = JSON.stringify(jobParmas)
            }
            param.fit_params = JSON.stringify(fitParams)
            param.evaluation_metrics = {
              type: 'fluent-bit'
            }
            param.code_selector = 'storagePath'
            this.handleCurrentNode(param)
            if (APIType) {
              this.currentNode.jobContent.ManiFest.API_type = APIType
            }
            if (paramType) {
              this.currentNode.jobContent.ManiFest.param_type = paramType
            }
            this.currentNode.jobContent.ManiFest.algorithm = algorithmType
            this.currentNode.jobContent.ManiFest.job_type = 'MLPipeline'
            this.$emit('saveNode', this.currentNode, temSave)
          } else {
            if (temSave) {
              return
            }
            this.$emit('changeTab', '2')
          }
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
    },
    handleAlgorithmNodeData (form, initForm) {
      this.$nextTick(() => {
        // 回写数据设置
        let dataSetting = {
          name: form.name,
          path: form.data_stores[0].connection.path,
          ...form.data_set
        }
        this.$refs.dataSetting.form = { ...dataSetting }
        this.$refs.customModel.form = { ...form.mf_model, model_push: !!form.mf_model.factory_name }
        // 算法节点
        if (form.job_params || form.fit_params) {
          let algorithmParam = {}
          const fitParams = form.fit_params ? JSON.parse(form.fit_params) : {}
          const algorithm = form.job_params ? JSON.parse(form.job_params) : {}
          if (!form.param_type || form.param_type === 'list') {
            Object.assign(algorithm, fitParams)
            const formatKey = this.algorithmDataFormatKey()
            for (let key in initForm) {
              if (algorithm[key] !== undefined) { // 判断节点是否有该字段
                // 将number类型转化为string
                if (formatKey.float.indexOf(key) > -1) {
                  let num = algorithm[key] + ''
                  // 判断数据是否为整数，因为提交将字符串格式化成number，1.0会变成1
                  if (num.indexOf('.') > -1) {
                    algorithmParam[key] = num
                  } else {
                    algorithmParam[key] = Number(algorithm[key]).toFixed(1)
                  }
                } else if (formatKey.int.indexOf(key) > -1) {
                  algorithmParam[key] = algorithm[key] + ''
                } else {
                  algorithmParam[key] = algorithm[key]
                }
              } else {
                algorithmParam[key] = initForm[key]
              }
            }
          } else {
            algorithmParam = { ...initForm }
            for (let key in fitParams) {
              if (key === 'early_stopping_rounds') {
                algorithmParam[key] = fitParams[key] + ''
              } else {
                algorithmParam[key] = fitParams[key]
              }
            }
            algorithmParam.custom_param = form.job_params
          }
          if (form.API_type) {
            algorithmParam.API_type = form.API_type
          }
          if (form.param_type) {
            algorithmParam.param_type = form.param_type
          }
          this.form = { ...algorithmParam }
        }
        // 回写资源信息
        let resource = {}
        if (form.imageType === 'Standard') {
          resource.imageOption = form.framework.version
        } else {
          resource.imageInput = form.framework.version
        }
        let keyArr = ['imageType', 'namespace', 'gpus', 'cpus']
        for (let key of keyArr) {
          resource[key] = form[key]
        }
        resource.memory = parseInt(form.memory)
        this.$refs.basicResource.form = { ...resource }
      })
    }
  }
}
