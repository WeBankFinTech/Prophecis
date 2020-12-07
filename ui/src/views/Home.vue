<template>
  <div>
    <breadcrumb-nav></breadcrumb-nav>
    <div class="home-page">
      <el-row class="content-box header">
        <div class="title">
          {{ $t('home.title') }}
        </div>
        <div class="introduce">
          {{ $t('home.introduction') }}
        </div>
      </el-row>
      <div class="label">
        {{ $t('home.modelTrain') }}
      </div>
      <div class="content-box">
        <el-row>
          <el-col :span="18"
                  class="notebooks">
            Notebooks<span>MLLabis</span>
          </el-col>
          <el-col :span="6"
                  class="btn-right">
            <el-button type="primary"
                       @click="goNotebook">
              {{ $t('home.instanceList') }}
            </el-button>
          </el-col>
        </el-row>
        <el-row class="center stauts-row">
          <el-col :span="4">
            {{ AIDEList.nbTotal }}
          </el-col>
          <el-col :span="4">
            {{ AIDEList.nbRunning }}
          </el-col>
          <el-col :span="4">
            {{ AIDEList.gpuCount }}
          </el-col>
        </el-row>
        <el-row class="center img-row">
          <el-col :span="4">
            <img :src="taskImg">
          </el-col>
          <el-col :span="4">
            <img :src="taskImg">
          </el-col>
          <el-col :span="4">
            <img :src="cardImg">
          </el-col>
        </el-row>
        <el-row class="center nape-text">
          <el-col :span="4">
            {{ $t('home.totalInstance') }}
          </el-col>
          <el-col :span="4">
            {{ $t('home.runInstance') }}
          </el-col>
          <el-col :span="4">
            {{ $t('home.cardNumber') }}
          </el-col>
        </el-row>
      </div>
    </div>
  </div>
</template>
<script type="text/ecmascript-6">
import task from '../assets/images/task.png'
import card from '../assets/images/card.png'
export default {
  data: function () {
    return {
      taskImg: task,
      cardImg: card,
      intervalAIDEFunc: '',
      AIDEList: {
        gpuCount: 0,
        nbRunning: 0,
        nbTotal: 0
      }
    }
  },
  mounted () {
    this.getAIDEList()
  },
  methods: {
    getAIDEList () {
      this.getAIDEFunc()
      this.intervalAIDEFunc = setInterval(() => {
        this.getAIDEFunc()
      }, 10000)
    },
    getAIDEFunc () {
      let url = `/aide/${this.FesEnv.aideApiVersion}/dashboards`
      this.FesApi.fetch(url, 'get').then(rst => {
        this.AIDEList = rst
      }, () => {
        this.AIDEList = {
          gpuCount: 0,
          nbRunning: 0,
          nbTotal: 0
        }
      })
    },
    goNotebook () {
      this.$router.push('/AIDE')
    }
  },
  destroyed: function () {
    clearInterval(this.intervalAIDEFunc)
  }
}
</script>
<style lang="scss" scoped>
.home-page {
  min-block-size: 100%;
  .header {
    .title {
      font-size: 24px;
      font-weight: 600;
      color: #333333;
      height: 60px;
      line-height: 60px;
    }
    .introduce {
      color: #657180;
      padding: 5px 0;
      font-size: 14px;
    }
  }
  .label {
    margin-top: 20px;
    padding-left: 10px;
    font-size: 18px;
    font-weight: bold;
  }
  .content-box {
    margin-top: 20px;
    padding: 10px 20px;
    box-shadow: 0 1px 10px rgba(0, 0, 0, 0.2);
    background-color: #fff;
    .btn-right {
      text-align: right;
    }
    .center {
      text-align: center;
    }
    .nape-text {
      color: #333333;
      font-size: 12px;
      font-weight: 600;
    }
    .distributed,
    .notebooks {
      color: #333333;
      font-size: 16px;
      font-weight: 600;
    }
    .distributed span,
    .notebooks span {
      margin-left: 5px;
      color: #999;
    }
    .stauts-row {
      color: #999;
      font-size: 30px;
      font-weight: 600;
    }
    .img-row img {
      width: 40px;
    }
  }
  .content-box.header {
    margin-top: 0px;
  }
}
</style>
