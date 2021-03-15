<template>
  <div>
    <breadcrumb-nav></breadcrumb-nav>
    <div class="logs-detail">
      <div class="operate-btn-box ">
        <div class="model-id left">
          {{ $t('DI.modelId') }}：{{ modelId }}
        </div>
      </div>
      <div class="log-box">
        <div class="no-log" v-if="logList.length<1">
          <!-- <div v-if="loading"
               v-loading="loading"
               key="loading">
          </div> -->
          <div>{{$t('DI.nolog')}}</div>
        </div>
        <ul v-else>
          <li v-for="(item,index) of logList" :key="index">
            <pre><code>{{item}}</code></pre>
          </li>
        </ul>
      </div>
    </div>
  </div>
</template>
<script>
export default {
  data: function () {
    return {
      modelId: '',
      logList: [],
      websock: {
        readyState: null
      },
      count: 0,
      loading: false
    }
  },
  mounted () {
    this.modelId = this.$route.params.modelId
    this.loading = true
    this.initWebSocket()
  },
  // watch: {
  //   'websock.readyState': {
  //     handler (oldStatus, newStatus) {
  //       console.log('newStatus', newStatus)
  //       console.log('this.loading1', this.loading)
  //       if (oldStatus.readyState) {
  //         this.loading = false
  //         console.log('this.loading2', this.loading)
  //       }
  //     }
  //   }
  // },
  methods: {
    initWebSocket () {
      this.count++
      if (this.count <= 2) {
        let userId = localStorage.getItem('userId')
        let host = window.location.host
        let webSocketUrl = `ws://${host}/di/${this.FesEnv.diApiVersion}/models/${this.modelId}/logs?follow=false&version=2017-02-13&mlss-userid=${userId}`
        this.websock = new WebSocket(webSocketUrl)
        this.websock.onmessage = this.websocketonmessage
        this.websock.onerror = this.websocketonerror
        // if (this.websock.readyState === 1 || this.websock.readyState === 3) {
        //   this.loading = false
        // } else {
        //   setTimeout(() => { this.loading = false }, 100)
        // }
      }
    },
    websocketonerror () { // 连接建立失败重连
      this.initWebSocket()
    },
    websocketonmessage (e) { // 数据接收
      this.logList.push(e.data.trim())
    }
  },
  destroyed: function () {
    this.websock.close()
  }
}
</script>
<style lang="scss" scoped>
.logs-detail {
  font-family: Arial;
  box-shadow: 0 1px 10px rgba(0, 0, 0, 0.2);
  height: 100%;
  .model-id {
    color: #333333;
    font-size: 14px;
    font-weight: 600;
    padding: 15px 0 5px 15px;
  }
  .log-box {
    color: #606266;
    ul {
      margin: 0;
      li {
        list-style: none;
      }
    }
    pre {
      code {
        white-space: pre-wrap;
        word-break: break-word;
        line-height: 20px;
      }
    }
  }
  .loading {
    text-align: center;
    padding: 100px 0;
  }
  .no-log {
    text-align: center;
    padding: 100px 0;
    font-size: 18px;
    font-weight: 600;
  }
}
</style>
