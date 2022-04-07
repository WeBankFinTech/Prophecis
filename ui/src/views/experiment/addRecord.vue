<template>
  <div class="add-content">
    <iframe :src="iframeUrl"
            class="iframe"
            id="addIframe"></iframe>
  </div>
</template>
<script>
export default {
  data () {
    return {
      iframeUrl: process.env.VUE_APP_BASE_SURL + '/mlflow#experiments/' + this.$route.params.mlflow_exp_id
    }
  },
  mounted () {
  },
  methods: {
  }
}
</script>
<style lang="scss" scoped>
.add-content {
  width: 100%;
  height: 100%;
  .iframe {
    border: 1px solid #e6e6e6;
    width: 100%;
    height: 100%;
  }
}
</style>
