/**
 * 操作Api
 */
import * as util from '../util/type'
import axios from 'axios'
import { Message } from 'element-ui'
import Router from '../router'
var trim = function (obj) {
  for (let p in obj) {
    if (util.isString(obj[p])) {
      obj[p] = obj[p].trim()
    } else if (util.isPlainObject(obj[p])) {
      trim(obj[p])
    } else if (util.isArray(obj[p])) {
      trim(obj[p])
    }
  }
}

const instance = axios.create({
  baseURL: process.env.VUE_APP_BASE_SURL,
  timeout: 10000,
  withCredentials: true
})
const api = {
  instance: instance,
  constructionOfResponse: {
    codePath: 'code',
    successCode: '200',
    messagePath: 'msg',
    resultPath: 'result'
  }
}

const success = function (response) {
  let resultFormat = (response.config && response.config.resultFormat) || api.constructionOfResponse
  if (util.isNull(resultFormat.codePath) || util.isNull(resultFormat.successCode) ||
    util.isNull(resultFormat.messagePath) || util.isNull(resultFormat.resultPath)) {
    return
  }

  var data
  if (util.isString(response.data)) {
    data = JSON.parse(response.data)
  } else if (util.isObject(response.data)) {
    data = response.data
  } else {
    throw new Error('后台服务异常')
  }
  let { code, message, result } = getData(data, resultFormat)
  if (resultFormat.successCode === '*') {
    return result || {}
  }
  code += '' // 兼容数字跟字符串
  if (code !== resultFormat.successCode) {
    if (api.error[code]) {
      api.error[code].forEach((fn) => fn(response))
      throw new Error('')
    } else {
      throw new Error(message || '后台服务异常')
    }
  }
  return result || {}
}

const getData = function (data, resultFormat) {
  let _arr = ['codePath', 'messagePath', 'resultPath']
  let arr = []; let rst = {}
  for (let i = 0; i < _arr.length; i++) {
    let pathArray = resultFormat[_arr[i]].split('.')
    let pathLength = pathArray.length
    let result
    if (pathLength === 1 && pathArray[0] === '*') {
      result = data
    } else {
      result = data[pathArray[0]]
    }
    for (let j = 1; j < pathLength; j++) {
      result = result[pathArray[j]]
      if (!result) {
        if (j < pathLength - 1) {
          // eslint-disable-next-line
          console.error('【FEX】ConstructionOfResponse配置错误：' + _arr[i] + '拿到的值是undefined，请检查配置')
        }
        break
      }
    }
    arr.push(result)
  }
  rst.code = arr[0]
  rst.message = arr[1]
  rst.result = arr[2]
  return rst
}
// 数据没有按规定格式返回
const successUT = function (response) {
  if (response.headers && response.headers.csrftoken) {

  }
  let data
  if (util.isString(response.data) && response.data) {
    data = JSON.parse(response.data)
  } else if (util.isObject(response.data)) {
    data = response.data
  }
  return data || {}
}

/**
 * 配置错误响应
 * @param option
 */
api.error = {
  // 使用sso登录，当401时跳转到统一登录地址
  401 (response) {
    if (window.location.hash.indexOf('login') === -1) {
      Router.push('/login')
      localStorage.clear()
    }
  },
  403 (response) {
    setTimeout(function () {
      Message({
        message: window.$i18n.t('common.403Pro'),
        type: 'error'
      })
    }, 100)
  }
}

const fail = function (error) {
  let _message = ''
  let errorPro = window.$i18n.t('common.networkException')
  let response = error.response
  if (response && api.error[response.status]) {
    api.error[response.status](response)
  } else {
    _message = errorPro
    try {
      if (response && response.data) {
        var data
        if (util.isString(response.data)) {
          data = JSON.parse(response.data)
        } else if (util.isObject(response.data)) {
          data = response.data
        }
        if (data) {
          _message = data.msg || data.message || data.error
        }
      }
    } catch (e) { }
  }
  error.message = _message
  throw error
}

const param = function (url, data, option) {
  let method = 'post'
  if (util.isNull(url)) {
    let errorPro = window.$i18n.t('common.urlPro')
    return Message({
      message: errorPro,
      type: 'error'
    })
  } else if (!util.isNull(url) && util.isNull(data) && util.isNull(option)) {
    option = {
      method: method
    }
  } else if (!util.isNull(url) && !util.isNull(data) && util.isNull(option)) {
    option = {
      method: method
    }
    if (util.isString(data)) {
      option.method = data
    } else if (util.isObject(data)) {
      option.data = data
    }
  } else if (!util.isNull(url) && !util.isNull(data) && !util.isNull(option)) {
    if (!util.isObject(data)) {
      data = {}
    }
    if (util.isString(option)) {
      option = {
        method: option
      }
    } else if (util.isObject(option)) {
      option.method = option.method || method
    } else {
      option = {
        method: method
      }
    }
    if (option.method === 'get' || option.method === 'delete' || option.method === 'head' || option.method === 'options') {
      option.params = data
    }
    if (option.method === 'post' || option.method === 'put' || option.method === 'patch') {
      option.data = data
    }
  }
  let config = {
    'Authorization': 'Basic dGVkOndlbGNvbWUx',
    'Mlss-Userid': localStorage.getItem('userId') || '',
    'Auth_type': 'LDAP'
  }
  if (util.isObject(option.headers)) {
    Object.assign(option.headers, config)
  } else {
    Object.assign(option.headers = {}, config)
  }
  // 如果传入button
  if (option && option.button) {
    option.button.currentDisabled = true
  }
  let _data = option.params || option.data
  if (_data && util.isObject(_data)) {
    trim(_data)
  }
  option.url = url

  return instance.request(option)
}

api.fetch = function (url, data, option) {
  return param(url, data, option)
    .then(success, fail)
    .catch((error) => {
      error.message && Message({
        message: error.message,
        type: 'error'
      })
      throw error
    }).finally(() => {
      // 如果传入button
      if (option && option.button) {
        option.button.currentDisabled = false
      }
    })
}

api.fetchUT = function (url, data, option) {
  return param(url, data, option)
    .then(successUT, fail)
    .catch((error) => {
      error.message && Message({
        message: error.message,
        type: 'error'
      })
      throw error
    }).finally(() => {
      // 如果传入button
      if (option && option.button) {
        option.button.currentDisabled = false
      }
    })
}

export default {
  fetch: api.fetch,
  fetchUT: api.fetchUT,
  fetchError: api.error
}
