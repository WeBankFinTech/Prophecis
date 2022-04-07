<template>
  <div class="page-process"
       :style="isFromDss? 'height:calc(100% + 36px)': 'height:calc(100% - 36px)'">
    <div class="process-tabs">
      <div v-if="!isFromDss"
           class="process-tab">
        <div v-for="(item, index) in tabs"
             :key="index"
             class="process-tab-item"
             :class="{active: index===active}"
             @click="choose(index)"
             @mouseenter.self="item.isHover = true"
             @mouseleave.self="item.isHover = false">
          <div>
            <img class="tab-icon"
                 :class="nodeImg[item.node.type].class"
                 :src="nodeImg[item.node.type].icon"
                 alt="">
            <div :title="$t(item.title)"
                 class="process-tab-name">{{ $t(item.title) }}</div>
            <span v-show="!item.isHover && item.node && item.node.isChange"
                  class="process-tab-unsave-icon fi-radio-on2" />
            <i v-if="item.isHover && item.close"
               class="el-icon-close"
               @click.stop="remove(index)" />
          </div>
        </div>
      </div>
      <div class="process-container">
        <template v-for="(item, index) in tabs">
          <Process ref="process"
                   v-if="item.type === 'Process'"
                   v-show="index===active"
                   :key="index"
                   :import-replace="false"
                   :flow-id="item.data.flowId"
                   :flow-name="item.data.name"
                   :version="item.data.version"
                   :readonly="item.data.readonly"
                   :tabs="tabs"
                   :open-files="openFiles"
                   :from-dss="isFromDss"
                   @isChange="isChange(index, arguments)"
                   @check-opened="checkOpened"
                   @saveBaseInfo="saveBaseInfo"
                   @reloadPage="reloadPage">
          </Process>
        </template>
      </div>
    </div>
  </div>
</template>
<script>
import Process from './module.vue'
import { NODETYPE, NODEICON } from '../../../services/nodeType'

export default {
  components: {
    Process
  },
  props: {
    query: {
      type: Object,
      default: () => { }
    }
  },
  data () {
    return {
      tabs: [{
        title: '',
        type: 'Process',
        close: false,
        data: this.query,
        node: {
          isChange: false,
          type: 'workflow.subflow'
        },
        key: 'flow.workflow',
        isHover: false
      }],
      isFromDss: false,
      active: 0,
      openFiles: {},
      nodeImg: NODEICON
    }
  },
  // dss 1.1版本
  created () {
    if (!this.$route.query.exp_id) {
      // dss跳转mlss
      // 提取参数
      let search = window.location.href
      const index = search.indexOf('?')
      search = decodeURIComponent(search.substring(index + 1))
      let searchArr = search.split('&')
      let param = {}
      for (let item of searchArr) {
        const paramKey = item.substring(0, item.indexOf('='))
        const paramVal = item.substring(item.indexOf('=') + 1)
        param[paramKey] = paramVal
      }
      console.log('param', param)
      this.tabs[0].data = {
        flowId: parseInt(param.expId),
        readonly: false,
        version: '',
        contextID: param.contextID,
        workspaceId: param.workspaceId,
        'bdp-user-ticket-id': param['bdp-user-ticket-id']
      }
      this.isFromDss = true
    } else {
      this.tabs[0].title = this.query.readonly === true ? 'flow.readOnlyMode' : 'flow.editMode'
    }
  },
  methods: {
    reloadPage () {
      this.reload()
    },
    choose (index) {
      this.active = index
    },
    remove (index) {
      // 删掉子工作流得删掉当前打开的子节点
      const currentTab = this.tabs[index]
      // 找到当前关闭项对应的子类
      const subArray = this.openFiles[currentTab.key] || []
      const changeList = this.tabs.filter((item) => {
        return subArray.includes(item.key) && item.node.isChange
      })
      if (changeList.length > 0 && currentTab.node.type === NODETYPE.FLOW) {
        let text = this.$t('flow.closeExitSubNodeWorkflow')
        if (currentTab.node.isChange) {
          text = this.$t('flow.closeUnSavedWorkflowTip')
        }
        this.$confirm(text, this.$t('common.prompt')).then((index) => {
          // 删除线先判断删除的是否是当前正在打开的tab，如果打开到最后一个tab，如果没有打开还是在当前的tab
          if (this.active === index) {
            // 删除的就是当前打开的
            this.tabs.splice(index, 1)
            this.choose(this.tabs.length - 1)
          } else {
            this.tabs.splice(index, 1)
            // this.choose(this.tabs.length - 1);
          }
        }).catch(() => { })
      } else {
        // 删除线先判断删除的是否是当前正在打开的tab，如果打开到最后一个tab，如果没有打开还是在当前的tab
        if (this.active >= index) {
          // 删除的就是当前打开的
          this.tabs.splice(index, 1)
          this.choose(index - 1)
        } else {
          this.tabs.splice(index, 1)
        }
      }
    },
    check (node) {
      if (node) {
        let boolean = true
        this.tabs.map((item) => {
          if (node.key === item.key) {
            boolean = true
          } else {
            if (this.tabs.length > 10) {
              boolean = false
              return
            }
            boolean = true
          }
        })
        if (!boolean) {
          this.$message.warning(this.$t('flow.openNodeTip'))
        }
        return boolean
      } else {
        if (this.tabs.length > 10) {
          this.$message.warning(this.$t('flow.openNodeTip'))
          return false
        }
        return true
      }
    },
    isChange (index, val) {
      if (this.tabs[index].node) {
        this.tabs[index].node.isChange = val[0]
      }
    },
    beforeLeaveHook () {
    },
    checkOpened (node, cb) {
    },
    saveBaseInfo (node) {
      this.tabs = this.tabs.map((item) => {
        if (item.key === node.key) {
          item.title = node.title
        }
        return item
      })
    }
  }
}
</script>
<style lang="scss" scoped>
@import "../../../assets/styles/variables";
.page-process {
  width: 100%;
  background-color: #f7f7f7;
  .bread-crumb {
    margin: 5px;
  }
}

.process-tabs {
  position: relative;
  height: calc(100%);
  overflow: hidden;
  &.no-tab {
    .process-tab {
      display: none;
    }
    .process-container {
      height: 100%;
    }
  }
  .process-tab {
    display: flex;
    flex-direction: row;
    flex-wrap: nowrap;
    justify-content: flex-start;
    align-items: flex-start;
    height: 32px;
    padding: 0 20px;
    box-sizing: border-box;
    border-top: 1px solid #e3e8ee;
    .process-tab-item {
      position: relative;
      height: 30px;
      line-height: 30px;
      padding: 0 20px;
      border-right: 1px solid #e3e8ee;
      background: #f7f7f7;
      color: #17233d;
      cursor: pointer;
      min-width: 100px;
      max-width: 200px;
      overflow: hidden;
      text-align: center;
      &:first-child {
        border-left: 1px solid #e3e8ee;
      }
      &.active {
        margin-top: -1px;
        &:before {
          content: "";
          position: absolute;
          top: -1px;
          left: 0;
          right: 0;
          height: 3px;
          background: $color-primary;
        }
        height: 32px;
        line-height: 31px;
        background: #eff4f9;
        color: $color-primary;
      }
      .process-tab-name {
        width: 100%;
        font-size: 12px;
        overflow: hidden;
        white-space: nowrap;
        text-overflow: ellipsis;
        padding-right: 5px;
        padding-left: 5px;
      }
      .process-tab-unsave-icon,
      .ivu-icon {
        position: absolute;
        right: 10px;
        top: 10px;
      }
      .tab-icon {
        position: absolute;
        left: 7px;
        top: 8px;
        width: 15px;
        height: 15px;
      }
      .flow {
        width: 12px;
        height: 12px;
      }
      .hivesql {
        width: 16px;
        height: 16px;
      }
      .display {
        width: 13px;
        height: 13px;
      }
      .dashboard {
        width: 14px;
        height: 14px;
      }
    }
  }
  .process-container {
    height: calc(100% - 45px);
  }
}
</style>
