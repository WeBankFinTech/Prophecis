<template>
  <el-tag :type="statusObj[status]">{{status}}</el-tag>
</template>
<script>
import { Tag } from 'element-ui'
export default {
  components: {
    ElTag: Tag
  },
  props: {
    statusObj: {
      type: Object,
      default () {
        return {
          success: 'success',
          info: 'info',
          warning: 'warning',
          danger: 'danger'
        }
      }
    },
    status: {
      type: String,
      default: 'success'
    }
  },
  data () {
    return {}
  }
}
</script>
<style lang="scss" scoped>
.el-tag {
  border-radius: 10px;
}
</style>
