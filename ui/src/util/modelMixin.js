import axios from 'axios'
export default {
  methods: {
    // AIDE DI
    getDIUserOption () {
      this.FesApi.fetch(`/cc/${this.FesEnv.ccApiVersion}/users/myUsers`, 'get').then(rst => {
        this.userOptionList = this.formatResList(rst)
      })
    },
    getDISpaceOption (header = { method: 'get' }) {
      this.FesApi.fetch(`/cc/${this.FesEnv.ccApiVersion}/namespaces/myNamespace`, {}, header).then(rst => {
        this.spaceOptionList = this.formatResList(rst)
      })
    },
    getDIStoragePath (header = { method: 'get' }) {
      this.FesApi.fetch(`/cc/${this.FesEnv.diApiVersion}/groups/group/storage`, {}, header).then(rst => {
        this.localPathList = this.formatResList(rst)
      })
    },
    getGroupList (header = { method: 'get' }) {
      const userId = localStorage.getItem('userId')
      this.FesApi.fetch(`/cc/${this.FesEnv.ccApiVersion}/auth/access/usergroups/${userId}`, {}, header).then((res) => {
        this.groupList = this.formatResList(res)
      })
    },
    resetQueryStrListData () {
      this.pagination.pageNumber = 1
      this.query.query_str = ''
      this.getListData()
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
    },
    downFile (url, fileName) {
      axios({
        url: url,
        method: 'get',
        headers: {
          'Mlss-Userid': localStorage.getItem('userId')
        },
        responseType: 'arraybuffer'
      }).then((res) => {
        const blob = new Blob([res.data], { type: 'application/zip' })
        const url = URL.createObjectURL(blob)
        const a = document.createElement('a')
        a.style.display = 'none'
        a.href = url
        a.setAttribute('download', fileName)
        a.click()
      }, (err) => {
        this.uint8Msg(err)
      }).catch((err) => {
        this.uint8Msg(err)
      })
    },
    uint8Msg (err) {
      const uint8Msg = new Uint8Array(err.response.data)
      const decodedString = String.fromCharCode.apply(null, uint8Msg)
      const data = JSON.parse(decodedString)
      this.$message.error(data.result)
    }
  }
}
