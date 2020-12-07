import Vue from 'vue'
import VueI18n from 'vue-i18n'
import messages from '../i18n-config.json'
import lang from 'element-ui/lib/locale/lang/en'
import cn from 'element-ui/lib/locale/lang/zh-CN'
import locale from 'element-ui/lib/locale'
Vue.use(VueI18n)
Object.assign(messages['en'], lang)
Object.assign(messages['zh-cn'], cn)

locale && locale.i18n((key, value) => VueI18n.t(key, value))
// 设置语言
locale && locale.use(lang)

let currentLanguage = sessionStorage.getItem('currentLanguage')
if (!currentLanguage) {
  let language = (navigator.language || navigator.browserLanguage).toLowerCase()
  if (language !== 'zh-cn') {
    currentLanguage = 'en'
    sessionStorage.setItem('currentLanguage', 'en')
  } else {
    currentLanguage = 'zh-cn'
    sessionStorage.setItem('currentLanguage', 'zh-cn')
  }
}

export default new VueI18n({
  locale: currentLanguage,
  messages
})
