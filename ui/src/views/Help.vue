<template>
  <div class="help">
    <breadcrumb-nav></breadcrumb-nav>
    <div class="content-box margin-top-0">
      <div class="title">
        {{ $t('help.commonDocuments') }}
      </div>
      <div class="content-doc">
        <el-button type="primary"
                   @click="helpDoc">
          {{ $t('help.helpDocument') }}
        </el-button>
      </div>
    </div>
  </div>
</template>
<script>
export default {
  methods: {
    helpDoc () {
      let url = process.env.VUE_APP_BASE_SURL + '/_book/index.html'
      window.open(window.encodeURI(url))
    }
  }
}
</script>
<style lang="scss" scoped>
.help {
  min-block-size: 100%;
  .content-box {
    margin-top: 30px;
    padding: 15px 20px;
    box-shadow: 0 1px 10px rgba(0, 0, 0, 0.2);
    background-color: #fff;
    .title {
      font-size: 18px;
      font-weight: 600;
      color: #333333;
      padding-bottom: 20px;
    }
  }
  .resource-selection,
  .common-operations {
    button {
      margin-right: 15px;
    }
  }
  .margin-top-0 {
    margin-top: 0;
  }
  .content-doc {
    div {
      padding: 5px 0;
    }
  }
}
</style>
