import util from './common'
export default {
  methods: {
    basicInfoProcess (item) {
      let keyArry = ['model_id', 'name', 'user_id', 'description']
      let modelsItem = {}
      modelsItem.job_type = item.JobType
      modelsItem.image = item.framework.version
      modelsItem.status = item.training.training_status.status
      modelsItem.namespace = item.job_namespace
      modelsItem.submission_timestamp = util.transDate(item.submission_timestamp)
      modelsItem.completed_timestamp = util.transDate(item.completed_timestamp)
      if (item.data_stores) {
        modelsItem.trainingData = item.data_stores[0].connection.bucket
        modelsItem.trainingResults = item.data_stores[1].connection.bucket
        // modelsItem.data_stores = item.data_stores
        modelsItem.path = item.data_stores[1].connection.path
        // data_stores 如果存在第三个对象 执行代码为共享目录
        if (item.data_stores[2]) {
          modelsItem.codeSettings = 'storagePath'
          modelsItem.diStoragePath = item.data_stores[2].connection.codeWork
        } else {
          modelsItem.codeSettings = 'codeFile'
          modelsItem.fileName = item.fileName
          modelsItem.filePath = item.filePath
        }
      }
      if (item.JobAlert && item.JobAlert !== 'null') {
        modelsItem.job_alert = this.processAlertInfo(item.JobAlert)
      }
      if (item.JobType === 'tfos') {
        this.handleHadoopData(modelsItem, item)
      } else {
        this.handleGPUData(modelsItem, item)
      }
      for (let key of keyArry) {
        modelsItem[key] = item[key]
      }
      return modelsItem
    },
    handleGPUData (modelsItem, trdata) {
      modelsItem.gpus = trdata.training.gpus || 0
      const keyArry = ['cpus', 'memory', 'command', 'learners', 'ps_cpu', 'ps_image', 'ps_memory', 'pss']
      for (let key of keyArry) {
        if (key === 'cpus' || key === 'command' || key === 'memory' || key === 'learners') {
          modelsItem[key] = trdata.training[key]
        } else {
          modelsItem[key] = trdata[key]
        }
      }
      if (trdata.proxy_user && this.$refs.gpuDialog && this.$refs.gpuDialog.proxyUserOption) {
        const proxyUserOption = this.$refs.gpuDialog.proxyUserOption
        for (let item of proxyUserOption) {
          if (item.name === trdata.proxy_user) {
            modelsItem.proxy_user = trdata.proxy_user
            break
          }
        }
        modelsItem.haveProxy = true
      }
    },
    // hadoop数据提取
    handleHadoopData (modelsItem, trdata) {
      const taskBaseObj = {
        codeSettings: 'HDFSPath',
        HDFSPath: '',
        fileName: ''
      }
      const tfos = trdata.TFosRequest
      modelsItem.driver_memory = tfos.DriverMemory
      modelsItem.executor_cores = tfos.ExecutorCores
      modelsItem.executor_memory = tfos.ExecutorMemory
      modelsItem.queue = tfos.Queue
      modelsItem.executors = tfos.Executors
      modelsItem.command = tfos.Command
      // const keyArry = ['Queue', 'Executors', 'Command']
      // for (let key of keyArry) {
      //   modelsItem[key] = tfos[key]
      // }
      modelsItem.pyfile = tfos.py_file
      modelsItem.archives = tfos.Archives
      taskBaseObj.HDFSPath = tfos.EntryPoint.hdfs
      modelsItem.entrypoint = { ...taskBaseObj }
      let pythonUrl = tfos.TensorflowEnv.hdfs
      let curPyth = this.FesEnv.DI.pythonOption.find((x) => x.value.indexOf(pythonUrl) > -1)
      if (curPyth !== undefined) {
        modelsItem.pythonType = 'Standard'
        modelsItem.pythonOption = curPyth.value
      } else {
        modelsItem.pythonType = 'Custom'
        modelsItem.pythonInput = pythonUrl
      }
    },
    // 列表返回告警详情数据重新组装，由于后台告警信息直接存的是字符串
    processAlertInfo (alert, isFormat = true) {
      let alertObj = isFormat ? window.JSON.parse(alert) : alert
      let arr = []
      let jobAlert = []
      alertObj.deadline && arr.push(alertObj.deadline.length)
      alertObj.event && arr.push(alertObj.event.length)
      alertObj.overtime && arr.push(alertObj.overtime.length)
      arr.sort()
      for (let i = 0; i < arr[arr.length - 1]; i++) {
        let obj = {
          alarmType: []
        }
        if (alertObj.event && alertObj.event[i]) {
          obj.alarmType.push('1')
          obj.event = alertObj.event[i]
        }
        if (alertObj.deadline && alertObj.deadline[i]) {
          obj.alarmType.push('2')
          obj.fixTime = alertObj.deadline[i]
          obj.fixTime.deadlineChecker = obj.fixTime.deadline_checker
          delete obj.fixTime.deadline_checker
        }
        if (alertObj.overtime && alertObj.overtime[i]) {
          obj.alarmType.push('3')
          obj.overTime = alertObj.overtime[i]
        }
        jobAlert.push(obj)
      }
      console.log('jobAlert', jobAlert)
      return jobAlert
    }
  }
}
