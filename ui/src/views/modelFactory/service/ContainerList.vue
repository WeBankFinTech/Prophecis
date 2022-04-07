<template>
  <div>
    <breadcrumb-nav></breadcrumb-nav>
    <div class="container-box"
         v-for="(item,index) in containersList"
         :key="index">
      <el-form :model="item"
               class="add-form"
               :disabled="true"
               label-width="100px">
        <el-row>
          <el-col :span="10">
            <el-form-item :label="$t('serviceList.podName')">
              <el-input v-model="item.name"></el-input>
            </el-form-item>
          </el-col>
          <el-col :span="10">
            <el-form-item :label="$t('DI.status')">
              <el-input v-model="item.status"></el-input>
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
      <el-table :data="item.containers"
                border
                style="width: 100%;marginTop:10px;"
                :empty-text="$t('common.noData')">
        <el-table-column prop="container_name"
                         :show-overflow-tooltip="true"
                         :label="$t('serviceList.containerName')"
                         min-width="150" />
        <el-table-column prop="cpu"
                         :show-overflow-tooltip="true"
                         label="CPU"
                         width="80" />
        <el-table-column prop="gpu"
                         :show-overflow-tooltip="true"
                         label="GPU"
                         width="80" />
        <el-table-column prop="memory"
                         label="Memory"
                         :show-overflow-tooltip="true"
                         width="100" />
        <el-table-column prop="image"
                         :show-overflow-tooltip="true"
                         :label="$t('DI.image')"
                         min-width="250" />
        <el-table-column prop="status"
                         :show-overflow-tooltip="true"
                         :label="$t('DI.status')" />
        <el-table-column prop="started_time"
                         min-width="135"
                         :show-overflow-tooltip="true"
                         :label="$t('serviceList.startTime')" />
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
                <el-dropdown-item :command="beforeHandleCommand('getLog',scope.row)">{{$t('DI.viewlog')}}</el-dropdown-item>
              </el-dropdown-menu>
            </el-dropdown>
          </template>
        </el-table-column>
      </el-table>
    </div>
    <view-log ref="viewLog"
              :url="logUrl"></view-log>
  </div>
</template>
<script>
import util from '../../../util/common'
import ViewLog from '../../../components/ViewLog.vue'
export default {
  components: {
    ViewLog
  },
  data () {
    return {
      containersList: [],
      logUrl: ''
    }
  },
  created () {
    this.getListData()
  },
  methods: {
    getListData () {
      this.FesApi.fetch(`/mf/v1/service/containers/${this.$route.params.namespace}/${this.$route.params.serviceName}`, 'get').then((res) => {
        for (let pod of res) {
          for (let item of pod.containers) {
            item.gpu = item.gpu ? item.gpu : 0
            item.started_time = util.transDate(new Date(item.started_time).getTime())
            item.status = item.status.substring(16, item.status.indexOf(':&')) // 截取状态值
          }
        }
        this.containersList = res
      })
    },
    getLog (trData) {
      this.logUrl = `/mf/v1/service/${this.$route.params.serviceName}/namespace/${trData.namespace}/container/${trData.container_name}/log`
      this.$refs.viewLog.logDialog = true
    }
  }
}
</script>
<style lang="scss" scoped>
.container-box {
  margin: 10px 0 30px;
}
</style>
