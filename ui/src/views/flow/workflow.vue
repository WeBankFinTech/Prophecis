<template>
  <div class="flow-page">
    <breadcrumb-nav></breadcrumb-nav>
    <div class="workflow-wrap">
      <div class="workflow-tabs">
        <div v-for="(item, index) in tabList"
             :key="index"
             class="workflow-tabs-item"
             @click="onTabClick(item)"
             :class="{active: item.version ? (currentVal.name===item.name && currentVal.version === item.version) : currentVal.name===item.name }"
             @mouseenter.self="item.isHover = true"
             @mouseleave.self="item.isHover = false">
          <div :title="item.label"
               class="workflow-tabs-name">{{ item.label }}</div>
          <i v-if="item.close"
             class="workflow-tabs-close el-icon-close"
             @click.stop="onTabRemove(item)" />
        </div>
      </div>
      <template v-for="(item, index) in tabList">
        <process :key="'flow'+index"
                 v-if="currentVal.name===item.name"
                 v-show="item.version ? (currentVal.name===item.name && currentVal.version === item.version) : currentVal.name===item.name"
                 :query="item.query"></process>
      </template>
    </div>
  </div>
</template>
<script>
import Process from './process'
export default {
  components: {
    process: Process.component
  },
  data () {
    return {
      tabList: [],
      currentVal: { name: 'view', version: 'index' },
      lastVal: null
    }
  },
  mounted () {
    // 折叠右侧菜单
    // if (!this.$store.state.isCollapse) {
    //   this.$store.commit('UPDATE_COLLAPSE', true)
    // }
    const query = this.$route.query
    this.currentVal = {
      close: false,
      flowId: 0,
      isHover: true,
      label: query.exp_name,
      name: query.exp_name,
      query: {
        flowId: query.exp_id,
        name: query.exp_name,
        readonly: query.readonly === true || query.readonly === 'true',
        version: ''
      }
    }
    this.tabList = [
      this.currentVal
    ]
  },
  methods: {
    onTabRemove (item) {
      const index = this.tabList.findIndex((subitem) => item.name === subitem.name && item.version === subitem.version)
      this.tabList.splice(index, 1)
      // 判断删除的是否是当前打开的，不是不用动
      if (item.name === this.currentVal.name && item.version === this.currentVal.version) {
        this.currentVal = this.lastVal
      }
    },
    onTabClick (item) {
      this.lastVal = this.currentVal
      this.currentVal = {
        name: item.name,
        version: item.version,
        ...item
      }
      const index = this.tabList.findIndex((subitem) => item.name === subitem.name && item.version === subitem.version)
      if (index < 0) {
        this.tabList.push(this.currentVal)
      }
    }
  },
  beforeDestroy () {
    // const collapse = sessionStorage.getItem('menuCollapse') === 'true'
    // this.$store.commit('UPDATE_COLLAPSE', collapse)
  }
}
</script>
<style lang="scss" scoped>
.flow-page {
  height: 100%;
}
.workflow-wrap {
  width: calc(100% - 10px);
  height: calc(100% - 5px);
  border: 1px solid #dcdee2;
  background: #f8f8f9;
  .workflow-tabs {
    width: 100%;
    height: 32px;
    margin: 4px 0;
    white-space: nowrap;
    line-height: 1.5;
    font-size: 14px;
    position: relative;
    color: #515a6e;
    border-bottom: 1px solid #dcdee2;
    padding: 0 5px;
    .workflow-tabs-item {
      margin: 0;
      height: 24px;
      padding: 5px 16px 4px;
      border-bottom: 1px solid #dcdee2;
      border-radius: 4px 4px 0 0;
      background: #f8f8f9;
      display: inline-block;
      cursor: pointer;
      position: relative;
      &.active {
        height: 24px;
        padding-bottom: 5px;
        background: #fff;
        transform: translateZ(0);
        border: 1px solid #dcdee2;
        border-bottom: 1px solid #fff;
        color: #2d8cf0;
      }
      .workflow-tabs-name {
        display: inline-block;
      }
      .workflow-tabs-close {
        width: 22px;
        margin-right: -6px;
        height: 22px;
        font-size: 14px;
        color: #999;
        text-align: right;
        vertical-align: middle;
        overflow: hidden;
        position: relative;
        top: 4px;
        transform-origin: 100% 50%;
        transition: all 0.3s ease-in-out;
        cursor: pointer;
      }
    }
  }
}
</style>
