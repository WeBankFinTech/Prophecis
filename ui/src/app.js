// import axios from 'axios'
export default function () {
  // 将环境配置赋值到FesEnv
  Object.assign(this.FesEnv, window.mlssEnvConfig[process.env.NODE_ENV])
  // 使用钩子函数对路由进行权限跳转
  this.$router.beforeEach((to, from, next) => {
    const userId = localStorage.getItem('userId')
    if (to.name !== 'login' && !userId) {
      next({ name: 'login' })
    } else {
      const superAdmin = localStorage.getItem('superAdmin')
      if (superAdmin === 'true') {
        next()
      } else {
        if (to.meta.permission) {
          this.$message.error(this.$t('noAccess'))
        } else {
          next()
        }
      }
    }
  })
  setTimeout(() => {
    let currentLanguage = sessionStorage.getItem('currentLanguage')
    if (currentLanguage) {
      this.$i18n.locale = currentLanguage
    } else {
      let language = (navigator.language || navigator.browserLanguage).toLowerCase()
      if (language !== 'zh-cn') {
        this.$i18n.locale = 'en'
        sessionStorage.setItem('currentLanguage', 'en')
      } else {
        sessionStorage.setItem('currentLanguage', 'zh-cn')
      }
    }
  }, 0)
}
