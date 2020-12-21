import Vue from 'vue'
import App from './App.vue'
import router from './router'
import VueI18n from 'vue-i18n'
import api from './services'
import messages from '../i18n-config.json'
import initFunc from './app'
import BreadcrumbNav from './components/BreadcrumbNav'
import lang from 'element-ui/lib/locale/lang/en'
import cn from 'element-ui/lib/locale/lang/zh-CN'
import locale from 'element-ui/lib/locale'
import './mixin'
import {
  Pagination,
  Dialog,
  Dropdown,
  DropdownMenu,
  DropdownItem,
  Menu,
  Submenu,
  MenuItem,
  MenuItemGroup,
  Input,
  Radio,
  RadioGroup,
  Switch,
  Select,
  DatePicker,
  Option,
  Button,
  Table,
  TableColumn,
  Form,
  FormItem,
  Icon,
  Row,
  Col,
  Loading,
  Message,
  MessageBox,
  Upload,
  Tabs,
  TabPane
} from 'element-ui'
import './assets/icon/iconfont.css'
// import './assets/flow-iconfont/iconfont.css'
Vue.config.productionTip = false
Vue.component('BreadcrumbNav', BreadcrumbNav)
Vue.use(api)
Vue.use(VueI18n)
Vue.use(Pagination)
Vue.use(Dialog)
Vue.use(Dropdown)
Vue.use(DropdownMenu)
Vue.use(DropdownItem)
Vue.use(Menu)
Vue.use(Submenu)
Vue.use(MenuItem)
Vue.use(MenuItemGroup)
Vue.use(Input)
Vue.use(Radio)
Vue.use(RadioGroup)
Vue.use(Select)
Vue.use(Switch)
Vue.use(DatePicker)
Vue.use(Option)
Vue.use(Button)
Vue.use(Table)
Vue.use(TableColumn)
Vue.use(Form)
Vue.use(FormItem)
Vue.use(Icon)
Vue.use(Row)
Vue.use(Col)
Vue.use(Loading.directive)
Vue.use(Upload)
Vue.use(Tabs)
Vue.use(TabPane)

Vue.prototype.$loading = Loading.service
Vue.prototype.$msgbox = MessageBox
Vue.prototype.$alert = MessageBox.alert
Vue.prototype.$confirm = MessageBox.confirm
Vue.prototype.$message = Message
Vue.prototype.$ELEMENT = { size: 'small', zIndex: 3000 }
Object.assign(messages['en'], lang)
Object.assign(messages['zh-cn'], cn)

const i18n = new VueI18n({
  locale: 'zh-cn',
  fallbackLocale: 'zh-cn',
  messages
})
locale.i18n((key, value) => i18n.t(key, value))

window.$i18n = i18n
// 设置语言
locale.use(lang)
var app = new Vue({
  router,
  i18n,
  render: h => h(App)
}).$mount('#app')

initFunc.call(app)
