<template>
  <el-row class="breadcrumb">
    <el-col :span="21">
      <el-breadcrumb separator="/">
        <el-breadcrumb-item v-for="(item,index) of labelArr"
                            :key="index"
                            :to="item.path?item.path:''">{{$t(item.title)}}</el-breadcrumb-item>
      </el-breadcrumb>
    </el-col>
    <el-col :span="3"
            class="operate-btn">
      <slot></slot>
    </el-col>
  </el-row>
</template>
<script>
import { Breadcrumb, BreadcrumbItem } from 'element-ui'
import sidebarConfig from './sidebarConfig'
export default {
  components: {
    ElBreadcrumb: Breadcrumb,
    ElBreadcrumbItem: BreadcrumbItem
  },
  props: {
    route1: {
      type: Object,
      default () {
        return {}
      }
    }
  },
  data () {
    return {
      labelArr: []
    }
  },
  mounted () {
    const breadcrumbArr = window.sessionStorage.getItem('breadcrumbArr')
    if (breadcrumbArr) {
      this.labelArr = window.JSON.parse(breadcrumbArr)
    }
    this.getPath()
  },
  methods: {
    getPath () {
      const path = this.$route.path
      // 如果this.labelArr最后一个对象path等于this.$route.path则是页面更新
      if (this.labelArr[this.labelArr.length - 1] === path) {
        return
      }
      let pathArr = path.split('/').splice(1)
      let startIndex = 1
      if (pathArr[0] === 'manage' || pathArr[0] === 'model') {
        startIndex = 2
      }
      if (pathArr.length - 1 >= startIndex) {
        let existIndex = -1
        for (let i = startIndex - 1; i < this.labelArr.length; i++) {
          if (this.labelArr[i].path === path) {
            existIndex = i
            return
          }
        }
        if (existIndex === -1) {
          this.labelArr.push({ title: this.$route.meta.title, path: path })
        } else {
          this.labelArr = this.labelArr.splice(0, existIndex + 1)
        }
        this.setBreadcrumbArr(this.labelArr)
      } else {
        this.levelMenu(pathArr, path)
      }
    },
    levelMenu (pathArr, path) {
      const groupName = this.$route.meta.groupName
      let labelArr = []
      if (groupName) {
        for (let menu of sidebarConfig.menuList) {
          if (menu.subs && groupName === menu.index) {
            labelArr.push({ title: menu.title })
            for (let subMenu of menu.subs) {
              if (subMenu.index === path) {
                labelArr.push({ title: subMenu.title, path: path })
                break
              }
            }
          }
        }
      } else {
        for (let menu of sidebarConfig.menuList) {
          if (menu.index === path) {
            labelArr.push({
              title: menu.title
            })
            labelArr.push({
              title: menu.title,
              path: menu.index
            })
          }
        }
      }
      this.labelArr = labelArr
      this.setBreadcrumbArr(labelArr)
    },
    setBreadcrumbArr (labelArr) {
      window.sessionStorage.setItem('breadcrumbArr', window.JSON.stringify(labelArr))
    }
  }
}
</script>
<style lang="scss" >
// 面包屑
.breadcrumb {
  .el-breadcrumb__item {
    height: 40px;
    line-height: 40px;
    &:last-child {
      .el-breadcrumb__inner {
        font-weight: 600;
      }
    }
  }
  .el-breadcrumb__inner a,
  .el-breadcrumb__inner.is-link {
    font-weight: 400;
  }
  .operate-btn {
    height: 40px;
    line-height: 40px;
    text-align: right;
  }
}
</style>
