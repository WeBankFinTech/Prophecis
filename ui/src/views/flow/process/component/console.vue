<template>
  <div ref="bottomPanel"
       v-if="openningNode"
       class="log-panel">
    <div class="workbench-tabs">
      <div class="workbench-tab-wrapper">
        <div class="workbench-tab">
          <div class="workbench-tab-item active">
            <span>{{$t('DI.viewlog')}}</span>
          </div>
        </div>
        <div class="workbench-tab-button">
          <i class="el-icon-close log-close"
             @click="closeConsole"></i>
        </div>
      </div>
      <div class="workbench-container">
        <!-- <log v-if="scriptViewState.showPanel == 'log'"
             :logs="script.log"
             :log-line="script.logLine"
             :script-view-state="scriptViewState" /> -->
        <div class="workbench-log-view"
             :style="{height:height+'px'}">
          <div class="log-tools">
            <div class="log-tools-control">
              <el-radio-group v-model="tabLabel">
                <el-radio-button label="all">All</el-radio-button>
                <el-radio-button label="error">Error</el-radio-button>
                <el-radio-button label="warning">Warning</el-radio-button>
                <el-radio-button label="info">Info</el-radio-button>
              </el-radio-group>
              <el-input v-model="searchText"
                        prefix-icon="el-icon-search"
                        :placeholder="$t('expExeRecord.enterSearchContent')"
                        size="mini"
                        class="log-search">
                <i class="el-icon-search"></i></el-input>
            </div>
          </div>
          <!-- <we-editor ref="logEditor"
                     :value="curLogs"
                     :style="{height:height-36+'px', overflow:'hidden'}"
                     @scrollChange="change"
                     @scrolling="scrolling"
                     @onload="editorOnload"
                     type="log" /> -->
          <div ref="logEditor"
               class="log-box"
               :style="{height:height-405+'px', 'overflow-y':'scroll'}">
            <!-- <ul>
              <li v-for="(item,index) of logs"
                  class="clearFloat"
                  :key="index">
                <div class="index">{{index}}</div><span>{{item}}</span>
              </li>
            </ul> -->
            <div class="col-banckground"></div>
            <table cellspacing="0"
                   class="log-table">
              <tbody>
                <tr v-for="(item,index) of curLog"
                    :key="index">
                  <td class="index">{{index+1}}</td>
                  <td class="content">{{item}}</td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
<script>
import util from '../../../../util/common'
import { debounce, forEach, filter } from 'lodash'
import { radioGroup, radioButton } from 'element-ui'
export default {
  components: {
    ElRadioGroup: radioGroup,
    ElRadioButton: radioButton
  },
  props: {
    node: {
      type: Object,
      default () {
        return {}
      }
    },
    logExecId: {
      type: [String, Number],
      default: ''
    },
    logType: {
      type: String,
      default: ''
    },
    jobStatus: {
      type: Object,
      default () {
        return {}
      }
    },
    openningNode: {
      type: Object,
      default: null
    }
  },
  data () {
    return {
      tabLabel: 'all',
      scriptViewState: {
        showPanel: 'log',
        height: 0,
        cacheLogScroll: 0
      },
      fromLine: 1,
      logs: [],
      logLine: 1,
      searchText: '',
      localLog: {
        log: { all: '', error: '', warning: '', info: '' },
        logLine: 1
      },
      // script: {
      //   result: {},
      //   steps: [],
      //   progress: {},
      //   log: { all: '', error: '', warning: '', info: '' },
      //   logLine: 1,
      //   resultList: null
      // },
      curLog: [],
      height: this.getHeight(),
      setTimeoutGetLog: null
    }
  },
  mounted () {
    this.scriptViewState = {
      ...this.scriptViewState,
      height: this.$el.clientHeight
    }
    this.byStatusGetLog()
  },
  watch: {
    // 如果发生拖动，对编辑器进行layout
    'height' (val) {
      if (val && this.$refs.logEditor.editor) {
        this.$refs.logEditor.editor.layout()
      }
    },
    logLine (val) {
      if (this.$refs.logEditor.editor) {
        this.$refs.logEditor.editor.revealLine(val)
      }
    },
    'scriptViewState' () {
      this.height = this.$parent.$el.clientHeight - 33
    },
    searchText () {
      this.getSearchList(this.localLog.log[this.tabLabel])
    },
    tabLabel () {
      this.curLog = Object.freeze(this.localLog.log[this.tabLabel].split('\n'))
    },
    logExecId () {
      this.localLog = {
        log: { all: '', error: '', warning: '', info: '' },
        logLine: 1
      }
      this.curLog = []
      this.byStatusGetLog()
    }
  },
  // computed: {
  //   // curLogs () {
  //   //   return this.formattedLogs()
  //   // }
  // },
  methods: {
    localLogShow () {
      if (!this.debounceLocalLogShow) {
        this.debounceLocalLogShow = debounce(() => {
          if (this.localLog) {
            this.logs = Object.freeze({ ...this.localLog.log })
            this.logLine = this.localLog.logLine
          }
        }, 1500)
      }
      this.debounceLocalLogShow()
    },
    closeConsole () {
      this.$emit('close-console')
      clearTimeout(this.setTimeoutGetLog)
    },
    // openPanel (type) {
    //   this.isLogShow = true
    // },
    byStatusGetLog () {
      if ('Inited,Scheduled,Running'.indexOf(this.jobStatus.state) > -1) {
        this.setTimeoutGetLog = setTimeout(() => {
          this.byStatusGetLog()
        }, 5000)
      } else {
        clearTimeout(this.setTimeoutGetLog)
      }
      this.getLogs()
    },
    getLogs () {
      // this.localLog = [
      //   'MODEL_DIR: /job/model-codedrwxr-xr-x 2 alexwu alexwu     4096 Mar 11 15:41 .ipynb_checkpointsdrwxr-xr-x 2 alexwu alexwu     4096 Mar 11 15:41 .ipynb_checkpointsdrwxr-xr-x 2 alexwu alexwu     4096 Mar 11 15:41 .ipynb_checkpointsdrwxr-xr-x 2 alexwu alexwu     4096 Mar 11 15:41 .ipynb_checkpointsdrwxr-xr-x 2 alexwu alexwu     4096 Mar 11 15:41 .ipynb_checkpointsdrwxr-xr-x 2 alexwu alexwu     4096 Mar 11 15:41 .ipynb_checkpointsdrwxr-xr-x 2 alexwu alexwu     4096 Mar 11 15:41 .ipynb_checkpointsdrwxr-xr-x 2 alexwu alexwu     4096 Mar 11 15:41 .ipynb_checkpoints',
      //   'TRAINING_JOB: ',
      //   ' RESULT_DIR: /mnt/results/result',
      //   'Storing trained model at:',
      //   'Contents of $MODEL_DIR'
      // ]
      // this.logs = this.localLog
      let url = this.logType === 'view' ? `/di/${this.FesEnv.diApiVersion}/experimentRun/${this.logExecId}/log` : `/di/${this.FesEnv.diApiVersion}/experimentJob/${this.logExecId}/log`
      this.FesApi.fetch(url, {
        from_line: this.fromLine,
        size: -1
      }, 'get')
        .then((rst) => {
          this.localLog = this.convertLogs(rst.log)
          // this.script.log = this.localLog.log
          const log = this.localLog.log[this.tabLabel]
          if (log) {
            this.curLog = Object.freeze(this.curLog.concat(log.split('\n')))
          }
          this.fromLine = rst.fromLine
          if (this.scriptViewState.showPanel === 'log') {
            this.localLogShow()
          }
          // this.updateNodeCache(['log'])
        })
    },
    convertLogs (logs) {
      const tmpLogs = this.localLog || {
        log: { all: '', error: '', warning: '', info: '' },
        logLine: 1
      }
      let hasLogInfo = Array.isArray(logs) ? logs.some((it) => it.length > 0) : logs
      if (!hasLogInfo) {
        return tmpLogs
      }
      const convertLogs = util.convertLog(logs)
      Object.keys(convertLogs).forEach((key) => {
        const convertLog = convertLogs[key]
        if (convertLog) {
          tmpLogs.log[key] += convertLog + '\n'
        }
        // if (key === 'all') {
        //   tmpLogs.logLine += convertLog.split('\n').length
        // }
      })
      return tmpLogs
    },
    change () {
      this.scriptViewState.cacheLogScroll = this.$refs.logEditor.editor.getScrollTop()
    },
    editorOnload () {
      this.$refs.logEditor.editor.setScrollPosition({ scrollTop: this.scriptViewState.cacheLogScroll || 0 })
    },
    // formattedLogs () {
    //   let logs = {
    //     all: '',
    //     error: '',
    //     warning: '',
    //     info: ''
    //   }
    //   Object.keys(this.logs).map((key) => {
    //     logs[key] = this.getSearchList(this.logs[key])
    //   })
    //   return this.getSearchList(this.logs)
    // },
    getSearchList (log) {
      let MatchText = ''
      const val = this.searchText
      if (!log) return MatchText
      if (val) {
        // 这部分代码是为了让正则表达式不报错，所以在特殊符号前面加上\
        let formatedVal = ''
        forEach(val, (o) => {
          if (/^[\w\u4e00-\u9fa5]$/.test(o)) {
            formatedVal += o
          } else {
            formatedVal += `\\${o}`
          }
        })
        // 全局和不区分大小写模式，正则是匹配searchText前后出了换行符之外的字符
        let regexp = new RegExp(`.*${formatedVal}.*`, 'gi')
        MatchText = filter(log.split('\n'), (item) => {
          return regexp.test(item)
        }).join('\n')
        regexp = null
      } else {
        MatchText = log
      }
      this.curLog = Object.freeze(MatchText.split('\n'))
    },
    scrolling (e, isBottom, isTop) {
      if (this.$parent && this.$parent.$el && ((isBottom && e.deltaY) || (e.deltaY < 0 && isTop))) {
        this.$parent.$el.scrollTop += e.deltaY
      }
    },
    getHeight () {
      return this.$parent.$el ? this.$parent.$el.clientHeight - 33 : 0
    }
  },
  beforeDestroy () {
    clearTimeout(this.setTimeoutGetLog)
  }
}
</script>
<style lang="scss" >
.log-panel {
  margin-top: 1px;
  border-top: 1px solid #d7dde4;
  background: #fff;
  &.full-screen {
    top: 40px;
    right: 0;
    bottom: 0;
    left: 0;
    position: fixed;
    z-index: 100;
    .workbench-tabs {
      height: auto;
    }
  }
  .workbench-tab-wrapper {
    border-top: none;
  }
  .workbench-tabs .workbench-tab .workbench-tab-item {
    text-align: center;
    border-top: none;
    background: #eff4f9;
    &.active {
      border-bottom: 1px solid #dcdee2;
      background: #f7f7f7;
      color: #17233d;
    }
  }
}
.workbench-tabs {
  position: relative;
  height: 100%;
  overflow: hidden;
  box-sizing: border-box;
  .workbench-tab-wrapper {
    display: flex;
    border-top: 1px solid #dcdcdc;
    border-bottom: 1px solid #dcdcdc;
    &.full-screen {
      position: fixed;
      left: 0;
      right: 0;
      top: 40px;
      z-index: 1200;
    }
    .workbench-tab {
      flex: 1;
      display: flex;
      flex-direction: row;
      flex-wrap: nowrap;
      justify-content: flex-start;
      align-items: flex-start;
      height: 32px;
      background-color: #fff;
      width: calc(100% - 45px);
      overflow: hidden;
      .workbench-tab-item {
        display: inline-block;
        height: 32px;
        line-height: 32px;
        background-color: #f7f7f7;
        color: #17233d;
        cursor: pointer;
        min-width: 100px;
        max-width: 200px;
        overflow: hidden;
        margin-right: 2px;
        border: 1px solid #eee;
        &.active {
          margin-top: 1px;
          height: 31px;
          background-color: #eff4f9;
          border-radius: 4px 4px 0 0;
          border: 1px solid #dcdee2;
          line-height: 30px;
        }
      }
    }
    .workbench-tab-button {
      flex: 0 0 30px;
      text-align: center;
      .log-close {
        font-size: 14px;
        margin-top: 8px;
        cursor: pointer;
      }
    }
  }
  .workbench-container {
    height: calc(100% - 36px);
    &.node {
      height: 100%;
    }
    .workbench-log-view {
      height: 100%;
      .log-tools {
        height: 36px;
        line-height: 36px;
        padding-left: 10px;
        background: #f7f7f7;
        position: relative;
        overflow: hidden;
        margin-bottom: -2px;
        .log-tools-control {
          display: inline-block;
          position: absolute;
          top: 2px;
          .log-search {
            width: 300px;
            position: absolute;
            left: 330px;
            top: -2px;
            font-size: 12px;
          }
        }
        .el-radio-button {
          &:first-child {
            .el-radio-button__inner {
              border-left: none;
              border-radius: 0;
            }
          }
        }
        .el-radio-button--small {
          .el-radio-button__inner {
            padding: 9px 24px;
            font-weight: 600;
          }
        }
        .el-radio-button__inner {
          background: transparent;
          border: none;
        }
        .el-radio-button__orig-radio {
          &:checked + .el-radio-button__inner {
            color: #409eff;
            -webkit-box-shadow: none;
            box-shadow: none;
          }
        }
      }
      .log-box {
        position: relative;
        .col-banckground {
          width: 50px;
          position: absolute;
          top: 0;
          left: 0;
          z-index: 3;
          background-color: #f2f3f6;
        }
        .log-table {
          position: relative;
          z-index: 10;
          .index {
            width: 50px;
            color: #2d8cf0;
            font-size: 12px;
            font-weight: 600;
            background: #f2f3f6;
            text-align: center;
            padding: 5px 0;
          }
          .content {
            padding-left: 5px;
          }
        }
      }
    }
  }
}
</style>
