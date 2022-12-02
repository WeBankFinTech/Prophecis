let util = {
  transDate (timestap, format) {
    if (!timestap) return ''
    format = format || 'yyyy-MM-dd hh:mm:ss'
    if (/^[0-9]*$/.test(timestap)) {
      timestap = Number(timestap)
    } else {
      timestap = new Date(timestap).getTime()
    }
    let time = new Date(timestap)
    var obj = {
      'y+': time.getFullYear(),
      'M+': time.getMonth() + 1,
      'd+': time.getDate(),
      'h+': time.getHours(),
      'm+': time.getMinutes(),
      's+': time.getSeconds()
    }

    if (new RegExp('(y+)').test(format)) {
      format = format.replace(RegExp.$1, obj['y+'])
    }
    for (var j in obj) {
      if (new RegExp('(' + j + ')').test(format)) {
        format = format.replace(RegExp.$1, (RegExp.$1.length === 1) ? (obj[j]) : (('00' + obj[j]).substr(('' + obj[j]).length)))
      }
    }
    return format
  },
  getTime () {
    const current = new Date()
    const year = current.getFullYear() + ''
    let month = current.getMonth() + 1
    month = month > 9 ? month : '0' + month
    let day = current.getDate()
    day = day > 9 ? day : '0' + day
    return year + month + day
  },
  // 上传校验
  convertLog (logs) {
    let logMap = {
      all: '',
      error: '',
      warning: '',
      info: ''
    }
    let newMap = {}
    if (typeof logs === 'string') {
      newMap = {
        all: logs
      }
    } else if (Array.isArray(logs)) {
      let keysArr = ['error', 'warning', 'info', 'all']
      logs.forEach((log, index) => {
        newMap[keysArr[index]] = log
      })
    } else if (_.isPlainObject(logs)) {
      newMap = logs
    }

    return Object.assign(logMap, newMap)
  }
}
export default util
