import Vue from 'vue'
import service from './api'
var api = {
  install () {
    Vue.prototype.FesApi = service
    Vue.prototype.FesEnv = {}
  }
}
export default api
