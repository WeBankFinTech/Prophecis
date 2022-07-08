<template>
  <div>
    <div class="align-right">
      <el-button type="primary"
                 @click="openList">
        {{$t('DI.enterCodeDirectory')}}
      </el-button>
    </div>
    <el-dialog title="notebook"
               :visible.sync="dialogVisible"
               custom-class="dialog-style aide-list-dialog"
               append-to-body
               width="1200px">
      <el-table :data="dataList"
                border
                style="width: 100%"
                :empty-text="$t('common.noData')">
        <el-table-column prop="name"
                         :label="$t('DI.name')"
                         min-width="100" />
        <el-table-column prop="user"
                         :label="$t('user.user')" />
        <el-table-column prop="namespace"
                         :label="$t('ns.nameSpace')"
                         min-width="220" />
        <el-table-column prop="cpu"
                         :label="$t('AIDE.cpu')"
                         min-width="50" />
        <el-table-column prop="gpu"
                         :label="$t('AIDE.gpu')"
                         min-width="50" />
        <el-table-column prop="memory"
                         :label="$t('AIDE.memory1')"
                         min-width="70" />
        <el-table-column prop="image"
                         :label="$t('DI.image')"
                         min-width="300" />
        <el-table-column prop="status"
                         min-width="110"
                         :label="$t('DI.status')">
          <template slot-scope="scope">
            <status-tag :status="scope.row.status"
                        :statusObj="statusObj"></status-tag>
          </template>
        </el-table-column>
        <el-table-column prop="uptime"
                         :label="$t('DI.createTime')"
                         min-width="125" />
        <el-table-column :label="$t('common.operation')"
                         min-width="85">
          <template slot-scope="scope">
            <el-dropdown size="medium"
                         @command="handleCommand">
              <el-button type="text"
                         class="multi-operate-btn"
                         icon="el-icon-more"
                         round></el-button>
              <el-dropdown-menu slot="dropdown">
                <el-dropdown-item :command="beforeHandleCommand('visitNotebook',scope.row)"
                                  :disabled="scope.row.status!=='Ready'">{{$t('AIDE.access')}}</el-dropdown-item>
              </el-dropdown-menu>
            </el-dropdown>
          </template>
        </el-table-column>
      </el-table>
      <el-pagination @size-change="handleSizeChange"
                     @current-change="handleCurrentChange"
                     :current-page="pagination.pageNumber"
                     :page-sizes="sizeList"
                     :page-size="pagination.pageSize"
                     layout="total, sizes, prev, pager, next, jumper"
                     :total="pagination.totalPage"
                     background>
      </el-pagination>
    </el-dialog>
  </div>
</template>
<script>
import util from '../util/common'
import StatusTag from './StatusTag'
const headerObj = {
  'Content-Type': 'application/json',
  'X-Watson-Userinfo': 'bluemix-instance-id=test-user',
  'MLSS-AppTimestamp': '1603274993343',
  'MLSS-AppID': 'MLFLOW',
  'MLSS-Auth-Type': 'SYSTEM',
  'MLSS-APPSignature': '91a585402398b4c3f63ab26b30722cfc2a74e3b3'
}
export default {
  components: { StatusTag },
  props: {
    rootPath: {
      type: String,
      default: ''
    },
    storagePath: {
      type: String,
      default: ''
    }
  },
  data () {
    return {
      dialogVisible: false,
      dataList: [],
      userId: localStorage.getItem('userId'),
      statusObj: {
        Ready: 'success',
        Waiting: '',
        Terminate: 'danger'
      }
    }
  },
  methods: {
    openList () {
      this.dialogVisible = true
      this.getListData()
    },
    getListData () {
      let url = ''
      let paginationOption = {
        page: this.pagination.pageNumber,
        size: this.pagination.pageSize,
        workDir: '/' + this.rootPath
      }

      let superAdmin = localStorage.getItem('superAdmin')
      if (superAdmin === 'true') {
        url = `/aide/${this.FesEnv.aideApiVersion}/namespaces/null/user/null/notebooks`
      } else {
        url = `/aide/${this.FesEnv.aideApiVersion}/user/${this.userId}/notebooks`
      }
      let header = {
        method: 'get'
      }
      if (window.location.href.indexOf('contextID') > -1) {
        header.headers = headerObj
      }
      this.FesApi.fetch(url, paginationOption, header).then(rst => {
        if (rst && JSON.stringify(rst) !== '{}') {
          rst.list && rst.list.forEach(item => {
            let index = item.image.indexOf(':')
            item.image = item.image.substring(index + 1)
            item.gpu = item.gpu || 0
            item.uptime = item.uptime ? util.transDate(item.uptime) : ''
            if (item.queue) {
              item.cluster = true
            }
          })
          this.dataList = Object.freeze(rst.list)
          this.pagination.totalPage = rst.total || 0
        }
      }, () => {
        this.dataList = []
      })
    },
    visitNotebook (trData) {
      let url = `/cc/${this.FesEnv.ccApiVersion}/auth/access/namespaces/${trData.namespace}/notebooks/${trData.name}`
      let xhr = new XMLHttpRequest()
      // const mlssToken = util.getCookieVal(localStorage.getItem('cookieKey')) || 0
      xhr.open('GET', url, false)
      let header = {
        'Mlss-Userid': this.userId
      }
      if (window.location.href.indexOf('contextID') > -1) {
        header = headerObj
      }
      for (let key in header) {
        xhr.setRequestHeader(key, header[key])
      }
      const that = this
      xhr.onreadystatechange = function (e) {
        if (this.readyState === 4 && this.status === 200) {
          let rst = JSON.parse(this.responseText).result
          if (rst.notebookAddress) {
            that.dialogVisible = false
            const url = rst.notebookAddress + '/lab/tree/' + that.storagePath
            window.open(window.encodeURI(url))
          }
        } else {
          let rst = JSON.parse(this.responseText)
          if (this.status === 401 || this.status === 403) {
            this.fetchError[this.status](rst)
          } else {
            this.$Message.error(rst.message)
          }
        }
      }
      xhr.send()
    }
  }
}
</script>
<style lang="scss" scoped>
.align-right {
  text-align: right;
}
</style>
