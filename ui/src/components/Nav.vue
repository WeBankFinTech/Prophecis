<template>
  <div class="top-nav">
    <img :src="prophecisImg"
         class="logo"
         @click="goHome" />
    <div class="operate-col">
      <el-dropdown size="medium"
                   @command="handleCommand">
        <el-button type="text"
                   class="user-name"
                   round>{{userName}}</el-button>
        <el-dropdown-menu slot="dropdown">
          <el-dropdown-item :command="beforeHandleCommand('logout')">{{$t('loginOut')}}</el-dropdown-item>
          <el-dropdown-item :command="beforeHandleCommand('languageSwitching')">{{$t('languageSwitching')}}</el-dropdown-item>
        </el-dropdown-menu>
      </el-dropdown>
    </div>
  </div>
</template>
<script>
import logo from '../assets/images/logo.png'
export default {
  data () {
    return {
      userName: localStorage.getItem('userId'),
      prophecisImg: logo
    }
  },
  methods: {
    goHome () {
      this.$router.push('/home')
    },
    languageSwitching () {
      let language = sessionStorage.getItem('currentLanguage', 'zh-cn')
      if (language === 'zh-cn') {
        this.$i18n.locale = 'en'
        sessionStorage.setItem('currentLanguage', 'en')
      } else {
        this.$i18n.locale = 'zh-cn'
        sessionStorage.setItem('currentLanguage', 'zh-cn')
      }
    },
    logout () {
      this.FesApi.fetchUT(`/cc/${this.FesEnv.ccApiVersion}/logout`, 'get').then(() => {
        this.$router.push('/login')
        localStorage.clear()
      }, (error) => {
        error.message && this.$message.error(error.message)
      })
    }
  }
}
</script>
<style lang="scss"  >
.top-nav {
  .logo {
    margin: 11px 20px 0;
    float: left;
    cursor: pointer;
    width: 140px;
  }
  .user-name {
    color: #fff;
    font-size: 18px;
  }
  .operate-col {
    float: right;
    padding-top: 5px;
  }
}
</style>
