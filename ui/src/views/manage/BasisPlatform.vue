<template>
  <div>
    <breadcrumb-nav></breadcrumb-nav>
    <div class="help">
      <div class="content-box"
           v-if="FesEnv.basisPlatform.dashboardUrl">
        <div class="title">
          {{ $t('basis.platformConsole') }}
        </div>
        <div class="resource-selection">
          <el-button type="primary"
                     @click="toApplication('dashboardUrl')">
            Kubernetes Dashboard
          </el-button>
        </div>
      </div>
      <div class="content-box"
           v-if="FesEnv.basisPlatform.grafanaUrl||FesEnv.basisPlatform.prometheusUrl">
        <div class="title">
          {{ $t('basis.platformMonitor') }}
        </div>
        <div class="common-operations">
          <el-button type="primary"
                     @click="toApplication('grafanaUrl')">
            Grafana
          </el-button>
          <el-button type="primary"
                     @click="toApplication('prometheusUrl')">
            Prometheus
          </el-button>
        </div>
      </div>
      <div class="content-box"
           v-if="FesEnv.basisPlatform.kibanaUrl">
        <div class="title">
          {{ $t('basis.logManagement') }}
        </div>
        <div class="content-doc">
          <el-button type="primary"
                     @click="toApplication('kibanaUrl')">
            Kibana
          </el-button>
        </div>
      </div>
    </div>
  </div>
</template>
<script>
export default {
  methods: {
    toApplication (url) {
      window.open(window.encodeURI(this.FesEnv.basisPlatform[url]), '_blank')
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
  .content-box:first-child {
    margin-top: 0px;
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
