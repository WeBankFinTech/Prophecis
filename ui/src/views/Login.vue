<template>
  <div class="login" @keyup.enter.stop.prevent="handleSubmit('loginForm')">
    <div class="login-img">
      <img src="../assets/images/login.svg" />
    </div>
    <div class="login-form">
      <div class="login-title">
        <img src="../assets/images/logo1.png" />
      </div>
      <el-form ref="loginForm" :model="loginForm" :rules="ruleInline">
        <el-form-item prop="username">
          <el-input v-model="loginForm.username" prefix-icon="el-icon-user" type="text" :placeholder="$t('login.username')" size="large">
          </el-input>
        </el-form-item>
        <el-form-item prop="password">
          <el-input v-model="loginForm.password" prefix-icon="el-icon-lock" type="password" :placeholder="$t('login.password')" size="large" />
        </el-form-item>
      </el-form>
      <div class="center">
        <el-button :loading="loading" type="primary" :disabled="disabled" @click="handleSubmit('loginForm')">{{$t('login.signIn')}}</el-button>
      </div>
    </div>
  </div>
</template>
<script>
import loginSvg from '../assets/images/login.svg'
import Jsencrypt from 'jsencrypt/bin/jsencrypt.min.js'
export default {
  data () {
    return {
      loading: false,
      loginForm: {
        username: '',
        password: ''
      },
      disabled: false,
      loginSvg: loginSvg,
      publicKeyData: '',
      ruleInline: {
        username: [
          { required: true, message: this.$t('login.usernameReq'), trigger: 'blur' }
        ],
        password: [
          { required: true, message: this.$t('login.passwordReq'), trigger: 'blur' }
        ]
      }
    }
  },
  created () {
    this.getPublicKey()
  },
  methods: {
    // 获取公钥接口
    getPublicKey () {
      this.FesApi.fetch(`/cc/${this.FesEnv.ccApiVersion}/getrsapubkey`, 'get').then((res) => {
        this.publicKeyData = res
      })
    },
    handleSubmit () {
      this.$refs.loginForm.validate((valid) => {
        if (valid) {
          this.disabled = true
          let params = {}
          if (this.publicKeyData) {
            const encryptor = new Jsencrypt()
            encryptor.setPublicKey(this.publicKeyData)
            const password = encryptor.encrypt(this.loginForm.password)
            params = {
              username: this.loginForm.username,
              password: password
            }
          } else {
            params = {
              username: this.loginForm.username,
              password: this.loginForm.password
            }
          }
          this.FesApi.fetch(`/cc/${this.FesEnv.ccApiVersion}/LDAPlogin`, params, 'post').then((res) => {
            this.disabled = false
            localStorage.setItem('userId', res.userName)
            res.isSuperadmin && localStorage.setItem('superAdmin', res.isSuperadmin)
            this.$router.push('/home')
          }).catch((res) => {
            this.disabled = false
            if (res.response.status === 401) {
              this.$message.error(this.$t('login.notPermissions'))
            }
            localStorage.clear()
          })
        }
      })
    }
  }
}
</script>
<style lang="scss" scoped>
.login {
  position: absolute;
  left: 0;
  right: 0;
  top: 0;
  bottom: 0;
  background: linear-gradient(83deg, #2e7bd9 50%, #234dca);
  display: flex;
  align-items: center;
  justify-content: flex-end;
  .login-img {
    margin-right: 15%;
    width: 40%;
  }
  .login-form {
    width: 380px;
    background-color: #fff;
    margin-right: 13%;
    padding: 20px;
    border-radius: 15px;
    box-shadow: 0px 0px 15px 0px rgba(0, 0, 0, 0.2);
    z-index: 8;
    .login-title {
      text-align: center;
      img {
        width: 200px;
        height: auto;
      }
    }
    input {
      // border: none;
      // height: 44px;
      // border-radius: 22px;
      padding-left: 30px;
      background-color: #f8f8f9;
      box-shadow: 0 0 0px 500px #f8f8f9 inset;
      color: #515a6e;
    }
    .center {
      text-align: center;

      button {
        width: 100%;
        height: 40px;
        margin: 30px 0;
        background: #5088f1;
        border-radius: 10px;
      }
    }
  }
  @media screen and (max-width: 1500px) {
    .login-img {
      width: 45%;
      margin-right: 12%;
    }
    .login-form {
      width: 320px;
      margin-right: 9%;
    }
  }
}
</style>
