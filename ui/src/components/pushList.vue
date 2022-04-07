<template>
  <el-dialog :title="$t('report.pushList')"
             :visible.sync="dialogVisible"
             custom-class="dialog-style"
             width="1150px">
    <el-table :data="pushList"
              border
              style="width: 100%;margin-top:15px"
              :empty-text="$t('common.noData')">
      <el-table-column prop="report_name"
                       v-if="type === 'report'"
                       :show-overflow-tooltip="true"
                       :label="$t('flow.reportName')" />
      <el-table-column prop="version"
                       :show-overflow-tooltip="true"
                       :label="$t('report.pushVersionNumber')"
                       width="100" />
      <el-table-column prop="user_name"
                       :show-overflow-tooltip="true"
                       :label="$t('user.user')"
                       width="95" />
      <el-table-column prop="fps_file_id"
                       :show-overflow-tooltip="true"
                       :label="$t('report.FPSFileID')" />
      <el-table-column prop="hash_value"
                       :show-overflow-tooltip="true"
                       :label="$t('report.FPSHashValue')" />
      <el-table-column prop="file_type"
                       :show-overflow-tooltip="true"
                       :label="$t('report.pushType')"
                       width="85" />
      <el-table-column prop="rmb_resp"
                       :show-overflow-tooltip="true"
                       :label="$t('report.RMBInfor')" />
      <el-table-column prop="status"
                       :show-overflow-tooltip="true"
                       width="75"
                       :label="$t('DI.status')">
      </el-table-column>
      <el-table-column prop="creation_timestamp"
                       :show-overflow-tooltip="true"
                       width="135"
                       :label="$t('DI.createTime')">
      </el-table-column>
      <el-table-column :label="$t('common.operation')"
                       width="85">
        <template slot-scope="scope">
          <el-dropdown size="medium"
                       @command="handleCommand">
            <el-button type="text"
                       class="multi-operate-btn"
                       icon="el-icon-more"
                       round></el-button>
            <el-dropdown-menu slot="dropdown">
              <el-dropdown-item :command="beforeHandleCommand('rmbLogDownload',scope.row)">{{$t('report.downloadLog')}}</el-dropdown-item>
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
</template>
<script>
import axios from 'axios'
import util from '../util/common'
export default {
  props: {
    type: {
      type: String,
      default: ''
    }
  },
  data () {
    return {
      dialogVisible: false,
      pushList: [],
      pushListUrl: ''
    }
  },
  methods: {
    openDialog () {
      this.dialogVisible = true
      this.getListData()
    },
    getListData () {
      this.FesApi.fetch(this.pushListUrl, {
        currentPage: this.pagination.pageNumber,
        pageSize: this.pagination.pageSize
      }, 'get').then((res) => {
        if (this.type === 'report') {
          this.handleReportData(res)
        } else {
          this.handleModelData(res)
        }
      })
    },
    handleReportData (res) {
      if (res.events) {
        for (let item of res.events) {
          item.report_name = res.report.report_name
          item.version = item.event.version
          item.eventId = item.event.id
          item.status = item.event.status
          item.fps_file_id = item.event.fps_file_id
          item.file_type = item.event.file_type
          item.hash_value = item.event.hash_value
          item.rmb_resp = item.event.params
          item.creation_timestamp = item.event.creation_timestamp ? util.transDate(new Date(item.event.creation_timestamp).getTime()) : ''
          item.user_name = item.user.name
        }
      }
      this.pushList = res.events
      this.pagination.totalPage = res.total || 0
    },
    handleModelData (res) {
      if (res.events) {
        for (let item of res.events) {
          item.version = item.event.version
          item.status = item.event.status
          item.fps_file_id = item.event.fps_file_id
          item.eventId = item.event.id
          item.file_type = item.event.file_type
          item.hash_value = item.event.hash_value
          item.rmb_resp = item.event.params
          item.creation_timestamp = item.event.creation_timestamp ? util.transDate(new Date(item.event.creation_timestamp).getTime()) : ''
          item.user_name = item.user.name
        }
        this.pushList = res.events
        this.pagination.totalPage = res.total || 0
      }
    },
    rmbLogDownload (trData) {
      if (!trData.eventId) {
        this.$message.error(this.$t('report.pushParameterAbnormal'))
        return
      }
      axios({
        url: `${process.env.VUE_APP_BASE_SURL}/mf/v1/rmbLog/Download/${trData.eventId}`,
        method: 'get',
        headers: {
          'Mlss-Userid': localStorage.getItem('userId')
        }
      }).then(rst => {
        // const source = new Blob([rst.data])
        let link = document.createElement('a')
        link.setAttribute('href',
          'data:text/plain;charset=utf-8,' + rst.data
        )
        link.style.display = 'none'
        link.setAttribute('download', this.type + 'RmbLog.txt')
        document.body.appendChild(link)
        link.click()
      }, () => {
        this.$message.error(this.$t('serviceList.downloadFailed'))
      })
    }
  }
}
</script>
