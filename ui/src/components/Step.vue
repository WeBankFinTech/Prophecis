<template>
  <div class="step-steps">
    <ul>
      <li v-for="(item,index) in stepList"
          :key="index">
        <span :class="{'steps-dot':'steps-dot','steps-current':index===step,'steps-finish':step>index}" /> {{ item }}
      </li>
    </ul>
  </div>
</template>
<script>
export default {
  props: {
    step: {
      type: Number,
      default () {
        return 0
      }
    },
    stepList: {
      type: Array,
      default () {
        return []
      }
    }
  }
}
</script>
<style lang="scss" scoped>
.step-steps {
  .steps-dot {
    display: inline-block;
    width: 16px;
    height: 16px;
    margin-right: 8px;
    border-radius: 50%;
    border: 2px solid #a4afbb;
    background-color: #fff;
    vertical-align: middle;
  }
  .steps-current {
    border-color: #55bc8a;
    -webkit-box-shadow: 0 4px 8px 0 rgba(85, 188, 138, 0.32);
    box-shadow: 0 4px 8px 0 rgba(85, 188, 138, 0.32);
  }
  .steps-finish {
    background: #55bc8a;
    border-color: #55bc8a;
  }
  ul {
    text-align: center;
    padding-bottom: 12px;
    > li {
      display: inline-block;
      vertical-align: middle;
      margin-bottom: 12px;
      + li {
        position: relative;
        margin-left: 75px;
      }
      + li:after {
        content: "";
        position: absolute;
        top: 50%;
        left: -58px;
        width: 30px;
        height: 0.5px;
        border: 1px dashed #ccd3db;
        -webkit-transform: translateY(-50%);
        -ms-transform: translateY(-50%);
        transform: translateY(-50%);
      }
    }
  }
}
</style>
