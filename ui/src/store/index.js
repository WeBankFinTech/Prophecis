import Vue from 'vue'
import Vuex from 'vuex'

Vue.use(Vuex)
export default new Vuex.Store({
  state: {
    user: {
      isLogin: false,
      userId: ''
    },
    isCollapse: false
  },
  mutations: {
    UPDATE_INFO (state, data) {
      state.user = { ...state.user, ...data }
    },
    UPDATE_COLLAPSE (state, collapse) {
      state.isCollapse = collapse
    }
  },
  actions: {

  }
})
