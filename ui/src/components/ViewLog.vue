<template>
  <el-dialog :title="$t('DI.viewlog')"
             :visible.sync="logDialog"
             @close="closeLogDialog"
             custom-class="dialog-style"
             :close-on-click-modal="false"
             width="920px">
    <div class="log-box">
      <div class="detail-arr"
           v-if="logDetail.length>0">
        <div v-for="(item,index) in logDetail"
             :key="index">
          <div v-if="isObject(item.log)">
            <div v-for="(val,key) in item.log"
                 class="detail-info"
                 :key="key">
              {{key}}:{{val}}
            </div>
          </div>
          <div v-else
               class="detail-info">
            {{item.log}}
          </div>
        </div>
      </div>
      <div v-else
           class="no-data">
        {{$t('common.noData')}}
      </div>

    </div>
    <el-pagination v-if="logPagination.totalPage>0"
                   @size-change="handleLogSizeChange"
                   @current-change="handleLogCurrentChange"
                   :current-page="logPagination.pageNumber"
                   :page-sizes="sizeList"
                   :page-size="logPagination.pageSize"
                   layout="total, sizes, prev, pager, next, jumper"
                   :total="logPagination.totalPage"
                   background>
    </el-pagination>
  </el-dialog>
</template>
<script>
export default {
  props: {
    url: {
      type: String,
      default: ''
    }
  },
  data () {
    return {
      logDialog: false,
      logDetail: [],
      sizeList: [20, 50, 100],
      logPagination: {
        pageSize: 20,
        pageNumber: 1,
        totalPage: 0
      }
    }
  },
  watch: {
    logDialog (val) {
      if (val) {
        this.getLog()
      }
    }
  },
  methods: {
    getLog () {
      this.FesApi.fetch(this.url, { currentPage: this.logPagination.pageNumber, pageSize: this.logPagination.pageSize }, 'get').then((res) => {
        for (let item of res.log_list) {
          try {
            if (typeof JSON.parse(item.log) === 'object') {
              item.log = JSON.parse(item.log)
            }
          } catch (error) {
          }
        }
        this.logDetail = Object.freeze(res.log_list)
        this.logPagination.totalPage = res.total || 0
        this.logDialog = true
      })
    },
    isObject (value) {
      const type = typeof value
      return !!value && type === 'object'
    },
    handleLogSizeChange (size) {
      this.logPagination.pageSize = size
      this.getLog()
    },
    handleLogCurrentChange (current) {
      this.logPagination.pageNumber = current
      this.getLog()
    },
    closeLogDialog () {
      this.logPagination.pageNumber = 1
    }
  }
}
</script>
<style lang="scss" scoped>
.detail-arr {
  padding: 5px;
}
.detail-info {
  padding: 5px 2px;
}
.log-box {
  height: 600px;
  overflow-y: scroll;
  .no-data {
    color: #979da9;
    font-size: 18px;
    text-align: center;
  }
}
</style>
