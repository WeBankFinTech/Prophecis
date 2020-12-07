import Vue from 'vue'
import { Message } from 'element-ui'
Vue.mixin({
  data () {
    return {
      btnDisabled: false,
      sizeList: [10, 50, 100],
      pagination: {
        pageSize: 10,
        pageNumber: 1,
        totalPage: 0
      }
    }
  },
  methods: {
    formatResList (res) {
      if (Array.isArray(res)) {
        return Object.freeze(res)
      }
      if (res.hasOwnProperty('models') && Array.isArray(res.models)) {
        return Object.freeze(res.models)
      }
      if (res.hasOwnProperty('list') && Array.isArray(res.list)) {
        return Object.freeze(res.list)
      }
      return []
    },
    getListData () {
      let paginationOption = {
        page: this.pagination.pageNumber,
        size: this.pagination.pageSize
      }
      this.FesApi.fetch(this.getUrl, paginationOption, 'get').then(rst => {
        this.dataList = this.formatResList(rst)
        this.pagination.totalPage = rst.total || 0
      })
    },
    handleSizeChange (size) {
      this.pagination.pageSize = size
      this.getListData()
    },
    handleCurrentChange (current) {
      this.pagination.pageNumber = current
      this.getListData()
    },
    clearDialogForm () {
      this.modify = false
      this.$refs.formValidate.resetFields()
    },
    deleteListItem (url, param = {}) {
      this.$confirm(this.$t('common.deletePro'), this.$t('common.prompt')).then((index) => {
        this.FesApi.fetch(url, param, 'delete').then(() => {
          this.deleteItemSuccess()
        })
      }).catch(() => { })
    },
    deleteItemSuccess () {
      if (this.dataList.length === 1 && this.pagination.pageNumber > 1) {
        --this.pagination.pageNumber
      }
      this.getListData()
      this.toast()
    },
    toast () {
      Message.success(this.$t('common.success'))
    },
    goBack () {
      this.$router.go(-1)
    },
    handleCommand (command) {
      if (typeof this[command.command] === 'function') {
        this[command.command](command.row)
      }
    },
    beforeHandleCommand (item, row) {
      return {
        'command': item,
        'row': row
      }
    },
    resetFormFields () {
      const that = this
      setTimeout(() => {
        that.$refs.formValidate && that.$refs.formValidate.resetFields()
      }, 0)
    },
    setBtnDisabeld () {
      this.btnDisabled = true
      setTimeout(() => {
        this.btnDisabled = false
      }, 2000)
    },
    getGroupOption () {
      this.FesApi.fetch(`/cc/${this.FesEnv.ccApiVersion}/groups`, 'get').then(rst => {
        this.groupOptionList = this.formatResList(rst)
      })
    },
    lastStep () {
      this.currentStep > 0 && this.currentStep--
    }
  }
})
