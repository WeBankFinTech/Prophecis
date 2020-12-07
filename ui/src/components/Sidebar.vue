<template>
  <el-scrollbar>
    <el-menu class="sidebar-el-menu"
             :default-active="onRoutes"
             :collapse="isCollapse"
             :collapse-transition="false"
             :router="true">
      <template v-for="item in menuList">
        <template v-if="item.subs&&((item.index==='manage'&&superAdmin==='true')||(item.index!=='manage'))">
          <el-submenu :index="item.index"
                      :key="item.index">
            <template slot="title">
              <i class="menu-icon iconfont"
                 :class="item.icon"></i>
              <span slot="title">{{ $t(item.title) }}</span>
            </template>
            <template v-for="subItem in item.subs">
              <el-submenu v-if="subItem.subs"
                          :index="subItem.index"
                          :key="subItem.index">
                <template slot="title">{{ $t(subItem.title) }}</template>
                <el-menu-item v-for="(threeItem,i) in subItem.subs"
                              :key="i"
                              :index="threeItem.index">{{ $t(threeItem.title) }}</el-menu-item>
              </el-submenu>
              <el-menu-item v-else
                            :index="subItem.index"
                            :key="subItem.index">{{ $t(subItem.title) }}</el-menu-item>
            </template>
          </el-submenu>
        </template>
        <template v-if="!item.subs">
          <el-menu-item :index="item.index"
                        :key="item.index">
            <i class="menu-icon iconfont"
               :class="item.icon"></i>
            <span slot="title">{{ $t(item.title) }}</span>
          </el-menu-item>
        </template>
      </template>
    </el-menu>
  </el-scrollbar>
</template>

<script>
import logo from '../assets/images/logo.png'
import sidebarConfig from './sidebarConfig'
import { Scrollbar } from 'element-ui'
export default {
  components: {
    ElScrollbar: Scrollbar
  },
  props: {
    isCollapse: {
      type: Boolean,
      default: false
    }
  },
  data () {
    return {
      logo: logo,
      superAdmin: localStorage.getItem('superAdmin')
    }
  },
  computed: {
    onRoutes () {
      const detail = {
        'settingDataGroup': '/manage/data',
        'settingUserGroup': '/manage/user',
        'settingNamespaceGroup': '/manage/namespace'
      }
      const path = detail[this.$route.name] === undefined ? this.$route.path : detail[this.$route.name]
      return path
    },
    menuList () {
      return sidebarConfig.menuList
    }
  },
  methods: {
    doCollapse () {
      this.$emit('collapseAside')
    }
  }
}
</script>
<style lang="scss" scoped>
.el-menu {
  border-right: none;
  .menu-icon {
    padding-right: 10px;
    font-size: 22px;
  }
}
</style>
