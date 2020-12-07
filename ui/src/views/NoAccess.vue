<template>
  <div />
</template>
<script>
export default {
  mounted: function () {
    let prompt = this.$route.query.noBADPAccess ? 'noBDAPAccessUser' : 'noAccessUser'
    this.$alert(this.$t(prompt), this.$t('common.prompt')).then(() => {
      this.logout()
    })
    setTimeout(() => {
      this.logout()
    }, 10000)
  },
  methods: {
    logout () {
      this.FesApi.fetchUT(`/cc/${this.FesEnv.ccApiVersion}/logout`, 'get').then(() => {
        window.location.href = window.encodeURI(`${this.FesEnv.ssoLogoutUrl}?service=${window.location.origin + '/index.html'}`)
      }, (error) => {
        error.message && this.$message.error(error.message)
      })
    }
  }
}
</script>
