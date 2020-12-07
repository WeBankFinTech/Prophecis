export default {
  methods: {
    // AIDE DI
    getDIUserOption () {
      this.FesApi.fetch(`/cc/${this.FesEnv.ccApiVersion}/users/myUsers`, 'get').then(rst => {
        this.userOptionList = this.formatResList(rst)
      })
    },
    getDISpaceOption () {
      this.FesApi.fetch(`/cc/${this.FesEnv.ccApiVersion}/namespaces/myNamespace`, {}, 'get').then(rst => {
        this.spaceOptionList = this.formatResList(rst)
      })
    },
    getDIStoragePath () {
      this.FesApi.fetch(`/cc/${this.FesEnv.ccApiVersion}/groups/group/storage`, {}, 'get').then(rst => {
        this.localPathList = this.formatResList(rst)
      })
    },
    resetQueryFields () {
      setTimeout(() => {
        this.$refs.queryValidate && this.$refs.queryValidate.resetFields()
      }, 10)
    },
    filterListData () {
      this.$refs.queryValidate.validate((valid) => {
        if (valid) {
          this.pagination.pageNumber = 1
          this.getListData()
        }
      })
    },
    resetListData () {
      this.pagination.pageNumber = 1
      this.$refs.queryValidate.resetFields()
      this.getListData()
    }
  }
}
