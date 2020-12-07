<template>
  <el-container class="container"
                direction="vertical">
    <el-header style="height:48px;">
      <layout-header />
    </el-header>
    <el-container>
      <el-aside :width="sidebWidth+'px'"
                :class="{'hide':!isChrome}">
        <sidebar :is-collapse="isCollapse"
                 @collapseAside="collapseAside" />
        <div class="collapse-btn"
             @click="collapseAside()"
             :style="'left:'+sidebWidth+'px'">
          <i :class="isCollapse?'el-icon-arrow-right':'el-icon-arrow-left'"></i>
        </div>
      </el-aside>
      <el-main :style="`padding: ${mainPadding};`">
        <router-view v-if="isRouterAlive" />
      </el-main>
    </el-container>
  </el-container>
</template>

<script>
import { Container, Main, Header, Aside } from 'element-ui'
import LayoutHeader from '../components/Nav'
import Sidebar from '../components/Sidebar'
export default {
  components: {
    ElContainer: Container,
    ElMain: Main,
    ElHeader: Header,
    ElAside: Aside,
    LayoutHeader,
    Sidebar
  },
  provide () {
    return {
      reload: this.reload
    }
  },
  data () {
    return {
      sidebWidth: 220,
      mainPadding: '28px 40px',
      isChrome: true,
      isRouterAlive: true,
      isCollapse: false
    }
  },
  created () {
    const collapse = sessionStorage.getItem('menuCollapse') === 'true'
    if (collapse) {
      this.isCollapse = collapse
      this.changeMinPadding(collapse)
    }
    this.isChrome = navigator.userAgent.indexOf('AppleWebKit') > -1
  },
  methods: {
    collapseAside () {
      this.isCollapse = !this.isCollapse
      this.changeMinPadding(this.isCollapse)
      sessionStorage.setItem('menuCollapse', this.isCollapse)
    },
    changeMinPadding (collapse) {
      if (collapse) {
        this.sidebWidth = 56
        this.mainPadding = '28px 40px'
      } else {
        this.sidebWidth = 220
        this.mainPadding = '28px 40px'
      }
    },
    reload () {
      this.isRouterAlive = false
      this.$nextTick(function () {
        this.isRouterAlive = true
      })
    }
  }
}
</script>

<style lang="scss" scoped>
@import "../assets/styles/variables.scss";
.container {
  height: 100%;
}
.el-header {
  background-color: #2e3846;
}
.el-aside {
  background-color: #fff;
  border-right: solid 1px #e6e6e6;
}
.hide {
  overflow: hidden;
}
.collapse-btn {
  position: absolute;
  top: 55px;
  width: 20px;
  cursor: pointer;
  i {
    padding: 3px 0;
    border: 1px solid #dcdfe6;
    background: #fff;
  }
}
</style>
