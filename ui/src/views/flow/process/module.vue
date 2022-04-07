<template>
  <div class="process-module"
       v-loading="loading">
    <vueProcess ref="process"
                :shapes="shapes"
                :value="originalData"
                :ctx-menu-options="nodeMenuOptions"
                :view-options="viewOptions"
                :disabled="workflowIsExecutor"
                @change="change"
                @message="message"
                @node-click="click"
                @node-baseInfo="saveNodeBaseInfo"
                @node-param="saveNodeParam"
                @node-delete="nodeDelete"
                @add="addNode"
                @ctx-menu-mycopy="copyNode"
                @ctx-menu-mypaste="pasteNode"
                @ctx-menu-allDelete="allDelete"
                @ctx-menu-console="openConsole"
                @ctx-menu-alarm="openAlarm">
      <template v-if="!fromDss">
        <div v-if="!workflowExeteId"
             :title="$t('flow.run')"
             class="button"
             @click="clickswitch">
          <span class="icon el-icon-video-play"></span>
          <span>{{ $t('flow.run') }}</span>
        </div>
        <div v-if="workflowExeteId&&'Inited,Scheduled,Running'.indexOf(flowStatus.state)>-1"
             :title="$t('flow.stop')"
             class="button"
             @click="clickswitch">
          <span class="icon el-icon-video-pause"
                style="color: red;"></span>
          <span>{{ $t('flow.stop') }}</span>
        </div>
        <div class="devider" />
        <div v-if="(!workflowExeteId||workflowIsExecutor)&&!readonly"
             :title="$t('common.save')"
             class="button"
             @click="handleSave">
          <svg class="icon"
               viewBox="0 0 1024 1024"
               version="1.1"
               xmlns="http://www.w3.org/2000/svg"
               p-id="4564">
            <path d="M176.64 1024c-97.28 0-176.64-80.384-176.64-178.688V178.688c0-98.816 79.36-178.688 176.64-178.688h670.72c97.28 0 176.64 80.384 176.64 178.688V845.312c0 98.816-79.36 178.688-176.64 178.688h-670.72z m0-936.96c-50.688 0-91.648 41.472-91.648 92.16V845.312c0 50.688 40.96 92.16 91.648 92.16h670.72c50.688 0 91.648-41.472 91.648-92.16V178.688c0-50.688-40.96-92.16-91.648-92.16h-3.584v437.248h-663.04v-437.248h-4.096z m581.632 350.208v-350.208h-492.544v350.208h492.544z m-160.768-35.328v-240.128h84.992v240.128h-84.992z"
                  p-id="4565"></path>
          </svg>
          <span>{{$t('common.save')}}</span>
        </div>
        <div class="devider" />
      </template>
    </vueProcess>
    <div class="process-module-param"
         v-show="nodebaseinfoShow"
         @click="clickBaseInfo">
      <div class="process-module-param-modal-content">
        <i class="el-icon-close"
           @click="closeClickNode"></i>
        <el-tabs v-model="activeName"
                 v-if="nodebaseinfoShow"
                 type="card">
          <el-tab-pane v-for="(item,index) in nodeTypeRefName[clickCurrentNode.type].tabPanes"
                       :label="item"
                       :name="index+''"
                       :key="index">
          </el-tab-pane>
          <component :is="nodeTypeRefName[clickCurrentNode.type].component"
                     :ref="nodeTypeRefName[clickCurrentNode.type].ref"
                     :node-data="clickCurrentNode"
                     :readonly="readonly"
                     :active-name="activeName"
                     :form-dss="fromDss"
                     :node-save-btn-disable="nodeSaveBtnDisable"
                     :proxyUserOption="proxyUserOption"
                     @changeTab="changeTab"
                     @saveNode="saveNode"></component>
        </el-tabs>
      </div>
    </div>
    <alarm-setting ref="AlarmSetting"
                   v-if="alarmVisible"
                   :nodeData="clickCurrentNode"
                   @closeAlarm="closeAlarm"></alarm-setting>
    <el-dialog :visible.sync="repetitionNameShow"
               :title="$t('flow.nodeNameNotice')"
               custom-class="dialog-style repetition-name">
      {{$t('flow.repeatNode')}}{{ repeatTitles.join(', ') }}?
      <div slot="footer">
        <el-button type="primary"
                   @click="repetitionName">{{$t('common.comfirm')}}</el-button>
      </div>
    </el-dialog>
    <template>
      <console v-if="openningNode"
               key="item.createTime"
               class="process-console"
               :style="{'width': `100% `}"
               :log-exec-id="logExecId"
               :log-type="logType"
               :job-status="clickLogJobStatus"
               :openningNode="openningNode"
               @close-console="closeConsole"></console>
    </template>
  </div>
</template>
<script>
import vueProcess from '../../../components/vue-process'
import console from './component/console'
import AlarmSetting from './component/AlarmSetting'
import shapes from './shape.js'
import { debounce } from 'lodash'
import { setTimeout, clearTimeout } from 'timers'
import Gpu from './component/gpu'
import Model from './component/model'
import ModelImage from './component/image'
import Service from './component/service'
import ReportPush from './component/ReportPush'
import nodeFormMixin from './component/nodeFormMixin'
import LogisticRegression from './component/LogisticRegression'
import DecisionTree from './component/DecisionTree'
import RandomForest from './component/RandomForest'
import XGBoost from './component/XGBoost'
import LightGBM from './component/LightGBM'
export default {
  mixins: [nodeFormMixin],
  components: {
    vueProcess,
    console,
    AlarmSetting,
    Gpu,
    Model,
    ModelImage,
    // ModelPush,
    ReportPush,
    LogisticRegression,
    DecisionTree,
    RandomForest,
    XGBoost,
    LightGBM
  },
  props: {
    flowId: {
      type: [String, Number],
      default: ''
    },
    flowName: {
      type: String,
      default: ''
    },
    version: {
      type: [String, Number],
      default: ''
    },
    readonly: {
      type: [String, Boolean],
      default: false
    },
    projectVersionID: {
      type: [String, Number],
      default: ''
    },
    openFiles: {
      type: Object,
      default: () => { }
    },
    tabs: {
      type: Array,
      default: () => []
    },
    fromDss: {
      type: Boolean,
      default: false
    }
  },
  data () {
    return {
      activeName: '0',
      name: '',
      versionId: '',
      shapes,
      // 原始数据
      originalData: null,
      // 插件返回的json数据
      json: null,
      // 工作流级别的参数
      props: [
        { 'user.to.proxy': '' }
      ],
      // 工作流级别的资源
      resources: [],
      saveModel: {
        comment: ''
      },
      // 是否有改变
      jsonChange: false,
      loading: false,
      ops: [], // 记录子工作流操作
      repetitionNameShow: false, // 节点名称重复提示
      repeatTitles: [],
      nodebaseinfoShow: false, // 自定义节点信息弹窗展示
      nodeSaveBtnDisable: false, // 节点表单是否填写标识，如果填到最后一个tabs显示保存按钮
      clickCurrentNode: {}, // 当前点击的节点
      timerClick: '',
      viewOptions: {
        showBaseInfoOnAdd: false, // 不显示默认的拖拽添加节点弹出的基础信息面板
        shapeView: true, // 左侧shape列表
        control: true
      },
      paramsIsChange: false, // 判断参数是否有在做操作改变，是否自动保存
      addNodeShow: false, // 创建节点的弹窗显示
      cacheNode: null,
      addNodeTitle: this.$t('flow.createSubWorkflow'), // 创建节点时弹窗的title
      workflowIsExecutor: false, // 当前工作流是否再执行
      openningNode: null, // 上一次打开控制台的节点
      shapeWidth: 180, // 流程图插件左侧工具栏的宽度
      workflowExeteId: '',
      workflowTaskId: '',
      excuteTimer: '',
      executorStatusTimer: '',
      // 新建widget的弹框中 当前选中的view
      selectBindView: null,
      // 当前选中的view
      viewName: '',
      // 当前选中的view
      viewContent: '',
      logExecId: '',
      logType: '',
      // 一个节点的key值与disabled属性的对应表，widget节点创建时，会添加一个 key值:true 到这个对象中，在请求成功或者失败后 delete keyDisabled[key]，因为widget节点的创建接口可能会pending很久（因为要等引擎启动），所以这个主要是用来控制在接口的pending中，不能重复请求创建接口
      keyDisabled: {},
      workflowExecutorCache: [],
      contextID: '',
      // pubulishShow: false,
      flowVersion: '',
      flowStatus: {
        state: ''
      }, // 记录工作流状态
      nodeStatus: {
        state: ''
      }, // 记录节点状态
      clickLogJobStatus: {}, // 查看日志任务运行状态
      alarmVisible: false,
      proxyUserOption: [] // 联合分析隐射
    }
  },
  computed: {
    reportTab () {
      return [this.$t('flow.datasetSettings'), this.$t('flow.modelSettings')]
    },
    algorithmTab () {
      return [this.$t('flow.dataSettings'), this.$t('flow.parameterSettings'), this.$t('flow.resourceSettings')]
    },
    // 动态节点配置
    nodeTypeRefName () {
      return {
        'linkis.appconn.mlflow.gpu': {
          ref: 'gpuNodeDialog',
          component: Gpu,
          tabPanes: [this.$t('flow.basicInformation'), this.$t('flow.resourceSetting'), this.$t('flow.codeSetting')]// 表单tab名称
        },
        'linkis.appconn.mlflow.model': {
          ref: 'modelNodeDialog',
          component: Model,
          tabPanes: []// 表单tab名称
        },
        'linkis.appconn.mlflow.image': {
          ref: 'imageNodeDialog',
          component: ModelImage,
          tabPanes: [this.$t('flow.basicInformation'), this.$t('flow.modelSettings')]
        },
        'linkis.appconn.mlflow.service': {
          ref: 'serviceNodeDialog',
          component: Service,
          tabPanes: [this.$t('flow.basicInformation'), this.$t('flow.resourceSetting'), this.$t('flow.modelSettings')]
        },
        'linkis.appconn.mlflow.reportPush': {
          ref: 'reportPushNodeDialog',
          component: ReportPush,
          tabPanes: [this.$t('flow.datasetSettings'), this.$t('flow.modelSettings')]
        },
        'linkis.appconn.mlflow.LogisticRegression': {
          ref: 'logisticRegressionNodeDialog',
          component: LogisticRegression,
          tabPanes: [this.$t('flow.dataSettings'), this.$t('flow.parameterSettings'), this.$t('flow.resourceSettings')]
        },
        'linkis.appconn.mlflow.DecisionTree': {
          ref: 'decisionTreeNodeDialog',
          component: DecisionTree,
          tabPanes: [this.$t('flow.dataSettings'), this.$t('flow.parameterSettings'), this.$t('flow.resourceSettings')]
        },
        'linkis.appconn.mlflow.RandomForest': {
          ref: 'randomForestNodeDialog',
          component: RandomForest,
          tabPanes: [this.$t('flow.dataSettings'), this.$t('flow.parameterSettings'), this.$t('flow.resourceSettings')]
        },
        'linkis.appconn.mlflow.XGBoost': {
          ref: 'XGBoostNodeDialog',
          component: XGBoost,
          tabPanes: [this.$t('flow.dataSettings'), this.$t('flow.parameterSettings'), this.$t('flow.resourceSettings')]
        },
        'linkis.appconn.mlflow.LightGBM': {
          ref: 'LightGBMNodeDialog',
          component: LightGBM,
          tabPanes: [this.$t('flow.dataSettings'), this.$t('flow.parameterSettings'), this.$t('flow.resourceSettings')]
        }
      }
    },
    type () {
      return 'flow'
    },
    nodeMenuOptions () {
      const _this = this
      return {
        defaultMenu: {
          config: false, // 不展示默认的基础信息菜单项
          param: false, // 不展示默认的参数配置菜单项
          copy: false, // 不使用默认menu复制
          delete: false// 不使用默认menu删除
        },
        userMenu: [],
        beforeShowMenu: (node, arr, type) => {
          // type : 'node' | 'link' | 'view' 分别是节点右键，边右键，画布右键
          // 如果有runState说明已经执行过,显示日志
          const execId = _this.$route.query.exec_id
          if (execId) {
            _this.logType = type
            if (type === 'node') {
              if (node.runState.taskID && node.runState) {
                _this.logExecId = node.runState.taskID
                _this.clickLogJobStatus = _this.nodeStatus
                // 当存在runState是证明该节点已运行，可以查看状态
                arr.push({
                  text: this.$t('DI.viewlog'),
                  value: 'console',
                  img: require('./images/menu/xitongguanlitai.svg')
                })
              } else {
                this.$message.info(this.$t('flow.nodeInitNoLog'))
              }
            } else {
              _this.logExecId = execId
              _this.clickLogJobStatus = _this.flowStatus
              arr.push({
                text: this.$t('DI.viewlog'),
                value: 'console',
                img: require('./images/menu/xitongguanlitai.svg')
              })
            }
          }
          if (type === 'node' && !_this.readonly) {
            if (_this.nodeCopy(node)) {
              arr.push({
                text: _this.$t('common.delete'),
                value: 'delete',
                img: require('./images/menu/delete.svg')
              })
              arr.push({
                text: _this.$t('common.copy'),
                value: 'mycopy',
                img: require('./images/menu/fuzhi.svg')
              })
            }
            if (['linkis.appconn.mlflow.gpu'].indexOf(node.type) > -1) {
              arr.push({
                text: _this.$t('flow.alertSetting'),
                value: 'alarm',
                img: require('./images/menu/fuzhi.svg')
              })
            }
          }
          if (type === 'view' && !_this.readonly) {
            arr.push({
              text: _this.$t('common.paste'),
              value: 'mypaste',
              img: require('./images/menu/zhantie.svg')
            })
            arr.push({
              text: _this.$t('flow.deleteAll'),
              value: 'allDelete',
              img: require('./images/menu/delete.svg')
            })
          }
          return arr
        }
      }
    }
  },
  watch: {
    jsonChange (val) {
      this.$emit('isChange', val)
    }
  },
  created () {
    if (this.fromDss) {
      // 获取用户信息
      this.fromDssGetUser()
      this.contextID = this.tabs[0].data.contextID
    }
    this.viewOptions.shapeView = !this.readonly
  },
  mounted () {
    if (this.readonly) {
      this.workflowExeteId = this.$route.query.exec_id
    }
    if (!this.fromDss) {
      // 基础信息
      this.getBaseInfo()
      this.getProxyUser()
    }
  },
  beforeDestroy () {
    if (this.timer) {
      clearInterval(this.timer)
    }
    if (this.excuteTimer) {
      clearTimeout(this.excuteTimer)
    }
    if (this.executorStatusTimer) {
      clearTimeout(this.executorStatusTimer)
    }
  },
  methods: {
    getProxyUser () {
      let headerParam = this.handlerFromDssHeader('get')
      const userId = this.getUserName()
      this.FesApi.fetch(`/cc/${this.FesEnv.ccApiVersion}/proxyUsers/${userId}`, {}, headerParam).then((res) => {
        this.proxyUserOption = Object.freeze(res)
      })
    },
    fromDssGetUser () {
      const param = this.tabs[0].data
      this.FesApi.fetch('di/v1/dssUserInfo', {
        'id': param['bdp-user-ticket-id'], // encodeURI(param['bdp-user-ticket-id'])
        workspaceId: param.workspaceId
      }, 'get').then((res) => {
        localStorage.setItem('userId', res.userName)
        this.getBaseInfo()
        this.getProxyUser()
      }, (err) => {
        this.$message.error(err.msg)
        localStorage.removeItem('userId')
      })
    },
    getUserName () {
      const username = localStorage.getItem('userId')
      return username
    },
    nodeCopy (node) {
      let flag = false
      this.shapes.forEach((item) => {
        if (item.children) {
          item.children.forEach((subItem) => {
            if (subItem.type === node.type) {
              flag = subItem.enableCopy
              return flag
            }
          })
        }
      })
      return flag
    },
    updateOriginData (node) {
      this.json.nodes = this.json.nodes.map((item) => {
        if (item.key === node.key) {
          item.jobContent = node.jobContent
          item.resources = node.resources
          item.params = node.params
        }
        return item
      })
      this.originalData = this.json
    },
    getBaseInfo () {
      this.loading = true
      this.getOriginJson()
    },
    initAction (json) {
      // 创建工作流之后就有值
      this.contextID = json.contextID ? json.contextID : ''
      // 保存节点才有的值
      if (json && json.nodes) {
        this.originalData = this.json = JSON.parse(JSON.stringify(json))
        this.resources = json.resources
        if (json.props) {
          this.props = json.props
        }
        // this.getSqlNodeList(json)
      }
      this.$nextTick(() => {
        this.init = true
        this.loading = false
      })
    },
    getOriginJson () {
      this.loading = false
      const header = this.handlerFromDssHeader('get')
      this.FesApi.fetch(`/di/v1/experiment/${this.flowId}`, {}, header).then((res) => {
        let json = this.convertJson(res)
        this.initAction(json)
        // 如果工作流已执行，获取工作流及节点状态
        if (this.workflowExeteId) {
          this.queryWorkflowExecutor(this.workflowExeteId)
        }
      })
    },
    convertJson (flow) {
      this.name = flow.prophecis_experiment.name
      let json = flow.flow_json

      // JSON:  优先缓存 > 通过id去查关联JSON数据
      if (json) {
        json = JSON.parse(json)
        if (json.nodes && Array.isArray(json.nodes)) {
          // 需要把 id -> key, jobType -> type
          json.nodes = json.nodes.map((node) => {
            // node.key = node.id;
            // 由于复制多个节点不知道为啥node.jobType
            if (node.jobType) {
              node.type = node.jobType
            } else {
              node.jobType = node.type
            }
            node = this.bindNodeBasicInfo(node)
            return node
          })
        }
      }
      return json
    },
    getProjectList () {
      const baseInfo = this.storage.get('baseInfo')
      return baseInfo.applications
    },
    change (obj) {
      this.json = obj
      if (this.init) {
        this.jsonChange = true
      }
    },
    initNode (arg) {
      arg = this.bindNodeBasicInfo(arg)
      this.clickCurrentNode = JSON.parse(JSON.stringify(arg))
    },
    click (arg) {
      // 当弹框已经打开又点击其他节点，需要初始化nodeSaveBtnDisable
      // this.nodeSaveBtnDisable = true
      clearTimeout(this.timerClick)
      if (this.workflowIsExecutor) return
      if (this.nodebaseinfoShow) { // 将上一个节点数据保存到flow
        this.temSaveNodeInfo()
      }
      this.activeName = '0'
      const _this = this
      this.timerClick = setTimeout(() => {
        _this.initNode(arg)
        _this.nodebaseinfoShow = true
        // this.$emit('node-click', arg)
      }, 200)
    },

    async saveNodeBaseInfo (arg, temSave = false) {
      // 当保存子流程节点的基础信息时，如果子流程节点没有 embeddedFlowId:"flow_id" 则先创建子流程节点
      let node = arg

      // 为了表单校验，基础信息弹窗保存的节点已不再是响应式，需重新赋值给json
      this.json.nodes = this.json.nodes.map((item) => {
        if (item.key === node.key) {
          if (['linkis.appconn.mlflow.model', 'linkis.appconn.mlflow.image', 'linkis.appconn.mlflow.service', 'linkis.appconn.mlflow.reportPush'].indexOf(node.type) > -1) {
            item.title = node.jobContent.content.name
          } else {
            item.title = node.jobContent.ManiFest.name
          }
          item.resources = node.resources || []
          item.jobType = node.type
          item.jobContent = node.jobContent
        }
        return item
      })
      this.originalData = this.json
      this.jsonChange = true
      // 保存工作流
      if (temSave) {
        return
      }
      this.autoSave('paramsSave', false)
    },
    saveNodeParam (...arg) {
      this.$emit('saveParam', arg)
    },
    /**
         * 显示保存视图
         */
    handleSave: debounce(function () {
      // 如果右侧弹框开打，保存前先将节点信息保存flow
      if (this.nodebaseinfoShow) {
        this.temSaveNodeInfo()
        setTimeout(() => {
          this.autoSave('手动保存', false)
        }, 200)
      } else {
        this.autoSave('手动保存', false)
      }
    }, 1500),
    /**
         * 保存工作流
         */
    save () {
      this.$refs.formSave.validate((valid) => {
        // 注释未输入
        if (!valid) {
          return
        }
        // 检查当前json是否有子节点未保存
        const subArray = this.openFiles[this.name] || []
        const changeList = this.tabs.filter((item) => {
          return subArray.includes(item.key) && item.node.isChange
        })
        // 保存时关闭控制台
        this.openningNode = null
        if (changeList.length > 0) {
          this.$Modal.confirm({
            title: this.$t('message.process.cancelNotice'),
            content: this.$t('message.process.noSaveHtml'),
            okText: this.$t('message.process.confirmSave'),
            cancelText: this.$t('message.newConst.cancel'),
            onOk: () => {
              this.saveModal = false
              let json = JSON.parse(JSON.stringify(this.json))
              json.nodes.forEach((node) => {
                this.$refs.process.setNodeRunState(node.key, {
                  borderColor: '#666'
                })
              })
              this.autoSave(this.saveModel.comment, false)
            },
            onCancel: () => {
            }
          })
        } else {
          this.saveModal = false
          let json = JSON.parse(JSON.stringify(this.json))
          json.nodes.forEach((node) => {
            this.$refs.process.setNodeRunState(node.key, {
              borderColor: '#666'
            })
          })
          this.autoSave(this.saveModel.comment, false)
        }
      })
    },
    autoSave (comment, f) {
      // 检查JSON
      if (!this.validateJSON()) {
        this.loading = false
        return false
      }
      this.nodeSaveBtnDisable = true
      // 需要把 key -> id,  type -> jobType title -> id
      let json = JSON.parse(JSON.stringify(this.json))
      json.edges.forEach((edge) => {
        let targets = json.nodes.filter((node) => {
          return node.key === edge.target
        })
        if (targets.length > 0) {
          let target = targets[0]
          edge.target = target.key
        }
      })
      // let flage = false
      json.nodes.forEach((node) => {
        // const reg = /^[a-zA-Z][a-zA-Z0-9_]*$/
        if (node.jobContent && node.jobContent.ManiFest) {
          node.jobContent.ManiFest.exp_id = this.flowId
          node.jobContent.ManiFest.exp_name = this.flowName
        } else if (node.jobContent && node.jobContent.content) {
          node.jobContent.content.exp_id = this.flowId
          node.jobContent.content.exp_name = this.flowName
        }
        node.id = node.key
        node.jobType = node.type
        node.selected = false // 保存之前初始化选中状态
        delete node.type
        delete node.menu // 删除菜单配置
        // 将用户保存的resources值为空字符串转为空数组
        if (!node.resources) {
          node.resources = []
        }
        // 保存之前删掉执行的状态信息
        if (node.runState) {
          delete node.runState
        }
      })
      return this.saveRequest(json, comment, f)
    },
    // 保存请求
    saveRequest (json, comment, f, cb) {
      const updateTime = Date.now()
      // 如果保存的时候代理用户为空加上默认用户
      if (!this.props[0]['user.to.proxy']) {
        this.props[0]['user.to.proxy'] = this.getUserName()
      }
      const paramsJson = JSON.parse(JSON.stringify(Object.assign(json, {
        comment: comment,
        type: this.type,
        updateTime,
        props: this.props,
        resources: this.resources,
        scheduleParams: this.scheduleParams,
        contextID: this.contextID,
        labels: {
          route: 'dev'
        }
      })))

      paramsJson.nodes = paramsJson.nodes.map((node) => {
        // 删除节点里的contextID
        if (node.contextID) {
          delete node.contextID
        }
        delete node.jobParams
        return node
      })
      const header = this.handlerFromDssHeader('put')
      return this.FesApi.fetch('/di/v1/experiment', {
        experiment_id: +this.flowId,
        flow_json: JSON.stringify(paramsJson)
      }, header).then((res) => {
        this.loading = false
        this.nodebaseinfoShow = false
        this.$message.success(this.$t('flow.saveWorkflowCuccess'))
        this.nodeSaveBtnDisable = false
        this.jsonChange = false
        if (cb) {
          cb()
        }
        // 保存成功后去更新tab的工作流数据
        return res
      }).catch(() => {
        this.nodeSaveBtnDisable = false
        this.loading = false
        if (cb) {
          cb()
        }
      })
    },
    /**
         * 检查JSON，是否符合规范
         * @return {Boolean}
         */
    validateJSON () {
      if (!this.json) {
        this.$message.warning(this.$t('flow.dragNode'))
        return false
      }
      let edges = this.json.edges
      let nodes = this.json.nodes
      if (!nodes || nodes.length === 0) {
        this.$message.warning(this.$t('flow.dragNode'))
        return false
      }
      let headers = []
      let footers = []
      let titles = []
      let repeatTitles = []
      nodes.forEach((node) => {
        if (titles.includes(node.title)) {
          repeatTitles.push(node.title)
        } else {
          titles.push(node.title)
        }
        if (!edges.some((edge) => {
          return edge.target === node.key
        })) {
          headers.push(node)
        }
        if (!edges.some((edge) => {
          return edge.source === node.key
        })) {
          footers.push(node)
        }
      })
      // 后台会把名称当做id处理，所以名称必须唯一
      if (repeatTitles.length > 0) {
        this.repeatTitles = repeatTitles
        this.repetitionNameShow = true
        return false
      }
      return true
    },
    changeTab (index) {
      this.activeName = index
    },
    closeClickNode () {
      // if (e.target.className === 'node-box-content' || e.target.className === 'node-box') return
      if (this.nodebaseinfoShow) {
        this.temSaveNodeInfo()
        setTimeout(() => {
          this.nodebaseinfoShow = false
        }, 200)
      }
    },
    temSaveNodeInfo () {
      const currentNodeRef = this.nodeTypeRefName[this.clickCurrentNode.type].ref
      this.$refs[currentNodeRef].save({}, true) // true表示临时保存数据不提交到后端
    },
    nodeDelete (node) {
      // 如果删除的是当前修改参数的节点，关闭侧边栏
      if (this.clickCurrentNode.key === node.key) {
        this.clickCurrentNode = {}
        this.nodebaseinfoShow = false
      }

      // 删除事件比jsonchange时机早
      const timeId = setTimeout(() => {
        this.autoSave('deleteSave', false)
        clearTimeout(timeId)
      }, 500)
    },
    repetitionName () {
      this.repetitionNameShow = false
    },
    // 单击节点出来的右边的弹框的保存事件
    saveNode (node, temSave = false) { // 保存节点参数配置
      this.saveNodeFunction(node, temSave)
    },
    // saveNode函数里面可以复用的代码
    saveNodeFunction (node, temSave) {
      // if (this.readonly) return this.$message.warning(this.$t('message.process.readonly'))
      this.saveNodeBaseInfo(node, temSave)
    },
    paramsChange (val) {
      this.paramsIsChange = val
    },
    handleTitle (clickCurrentNode) {
      switch (clickCurrentNode.type) {
        case 'linkis.appconn.mlflow.model':
          clickCurrentNode.title = 'MODEL'
          break
        case 'linkis.appconn.mlflow.image':
          clickCurrentNode.title = 'IMAGE'
          break
        case 'linkis.appconn.mlflow.service':
          clickCurrentNode.title = 'SERVICE'
          break
        case 'linkis.appconn.mlflow.reportPush':
          clickCurrentNode.title = 'reportPush'
          break
        case 'linkis.appconn.mlflow.LogisticRegression':
          clickCurrentNode.title = 'LogisticRegression'
          break
        case 'linkis.appconn.mlflow.DecisionTree':
          clickCurrentNode.title = 'DecisionTree'
          break
        case 'linkis.appconn.mlflow.RandomForest':
          clickCurrentNode.title = 'RandomForest'
          break
        case 'linkis.appconn.mlflow.XGBoost':
          clickCurrentNode.title = 'XGBoost'
          break
        case 'linkis.appconn.mlflow.LightGBM':
          clickCurrentNode.title = 'LightGBM'
          break
      }
      return clickCurrentNode.title + '-' + Math.round(Math.random() * 10000)
    },
    addNode (node) {
      // 关闭右侧弹窗
      this.nodebaseinfoShow = false
      // 右侧表单保存按钮disable
      // this.nodeSaveBtnDisable = true
      // 新拖入的节点，自动生成新的不重复名称,给名称后面加四位随机数
      node = this.bindNodeBasicInfo(node)
      this.clickCurrentNode = JSON.parse(JSON.stringify(node))
      this.clickCurrentNode.title = this.handleTitle(this.clickCurrentNode)
      // 弹窗提示由后台控制
      if (node.shouldCreationBeforeNode) {
        this.addNodeShow = true
        this.addNodeTitle = this.$t('message.process.createNode')
      } else {
        // 还得同步更新json中的node
        this.json.nodes = this.json.nodes.map((subItem) => {
          if (subItem.key === this.clickCurrentNode.key) {
            subItem.title = this.clickCurrentNode.title
            // Object.assign(subItem, this.clickCurrentNode)
          }
          return subItem
        })
        this.originalData = this.json
        // this.autoSave('新增节点', false)
      }
    },
    validatorName (rule, value, callback) {
      if (value === `${this.name}_`) {
        callback(new Error(this.$t('flow.nodeNameValid')))
      } else {
        callback()
      }
    },
    addFlowCancel () {
      this.addNodeShow = false
      // 删除未创建成功的节点
      this.json.nodes = this.json.nodes.filter((subItem) => {
        return this.clickCurrentNode.key !== subItem.key
      })
      this.originalData = this.json
    },

    // addFlowOk函数里可以复用的操作
    addFlowOkFunction () {
      this.addNodeShow = false
      // if (this.readonly) return this.$message.warning(this.$t('message.process.readonlyNoCeated'))
      this.saveNodeBaseInfo(this.clickCurrentNode)
    },
    copyNode (node) {
      node = this.bindNodeBasicInfo(node)
      if (node.enable_copy === false) return this.$message.warning(this.$t('flow.noCopy'))
      // if (this.readonly) return this.$message.warning(this.$t('message.process.readonlyNoCopy'))
      this.cacheNode = JSON.parse(JSON.stringify(node))
    },
    pasteNode (e) {
      // if (this.readonly) return this.$message.warning(this.$t('message.process.noPaste'))
      if (!this.cacheNode) {
        return this.$message.warning(this.$t('flow.firstCopy'))
      }
      // 获取屏幕的缩放值
      let pageSize = this.$refs.process.getState().baseOptions.pageSize
      const key = '' + new Date().getTime() + Math.ceil(Math.random() * 100)
      this.cacheNode.key = key
      this.cacheNode.createTime = Date.now()
      this.cacheNode.layout = {
        height: this.cacheNode.height,
        width: this.cacheNode.width,
        x: (e.offsetX / pageSize),
        y: (e.offsetY / pageSize)
      }
      this.json.nodes.push(this.cacheNode)
      this.nodeSelectedFalse(this.cacheNode)
      this.originalData = this.json
    },
    // 由于插件的selected不是响应式，所以得手动改变
    nodeSelectedFalse (node = {}) {
      this.json.nodes = this.json.nodes.map((subItem) => {
        if (node.key && node.key === subItem.key) {
          subItem.selected = true
        } else {
          subItem.selected = false
        }
        return subItem
      })
      this.originalData = this.json
    },
    clickBaseInfo () {
      // 插件组件内的节点outSide事件后执行，会将chosing清空，所以加个异步后执行，但是页面能看到节点样式抖动，后面想办法优化
      const teimer = setTimeout(() => {
        this.nodeSelectedFalse(this.clickCurrentNode)
        clearTimeout(teimer)
      }, 0)
    },
    // 批量删除选中节点
    allDelete () {
      if (this.readonly) return this.$message.warning(this.$t('message.process.noDelete'))
      const selectNodes = this.$refs.process.getSelectedNodes().map((item) => item.key)
      this.json.nodes = this.json.nodes.filter((item) => {
        if (selectNodes.includes(item.key)) {
          this.json.edges = this.json.edges.filter((link) => {
            return !(link.source === item.key || link.target === item.key)
          })
        } else {
          item.selected = false
          return true
        }
      })
      this.originalData = this.json
    },
    openConsole (node) {
      this.nodeSelectedFalse(node)
      this.openningNode = node
      this.shapeWidth = this.$refs.process && this.$refs.process.state.shapeOptions.viewWidth
    },
    closeConsole () {
      this.openningNode = null
    },
    openAlarm (node) {
      // 判断基本信息是否填写下
      if (node.jobContent && node.jobContent.ManiFest.name) {
        this.alarmVisible = true
        this.clickCurrentNode = node
        setTimeout(() => {
          this.$refs.AlarmSetting.alarmVisible = true
        }, 200)
      } else {
        this.$message.error(this.$t('flow.setBasicInfoPro'))
        if (!this.nodebaseinfoShow) {
          this.click(node)
        }
      }
    },
    closeAlarm () {
      this.alarmVisible = false
    },
    // 根据节点类型将后台节点基础信息加入
    bindNodeBasicInfo (node) {
      const shapes = JSON.parse(JSON.stringify(this.shapes))
      shapes.map((item) => {
        if (item.children.length > 0) {
          item.children.map((subItem) => {
            if (subItem.type === node.type) {
              node = Object.assign(subItem, node)
            }
          })
        }
      })
      return node
    },
    // 点击节流
    clickswitch: debounce(
      function () {
        this.workflowExeteId ? this.workflowStop() : this.workflowRun()
      }, 1000
    ),
    async workflowRun () {
      // this.dispatch('IndexedDB:clearNodeCache')
      // 重新执行清掉上次的计时器
      clearTimeout(this.excuteTimer)
      clearTimeout(this.executorStatusTimer)

      this.openningNode = null
      let fillNode = false
      if (this.json && this.json.nodes) {
        for (let item of this.json.nodes) {
          if (item.jobContent && (item.jobContent.ManiFest || item.jobContent.content)) {
            fillNode = true
          }
          item.id = item.key
          item.jobType = item.type
          delete item.editParam
          delete item.editBaseInfo
          delete item.submitToScheduler
          delete item.enableCopy
          delete item.shouldCreationBeforeNode
          delete item.supportJump
        }
      }
      if (!fillNode) {
        this.$message.error(this.$t('flow.fillNodePro'))
        return
      }
      let json = JSON.parse(JSON.stringify(this.json))
      json.type = 'flow'
      json.resources = []
      json.comment = 'paramsSave'
      json.contextID = this.contextID
      this.FesApi.fetch(`/di/${this.FesEnv.diApiVersion}/experimentRun/${this.flowId}`, {
        exp_exec_type: 'MLFLOW',
        flow_json: JSON.stringify(json)
      }, 'post').then((res) => {
        this.workflowExeteId = res.id
        this.viewOptions = {
          showBaseInfoOnAdd: false,
          shapeView: false,
          control: true
        }
        this.$router.replace({
          query: {
            exp_name: this.$route.query.exp_name,
            readonly: true,
            exp_id: this.flowId,
            exec_id: res.id
          }
        })
        this.$emit('reloadPage')
        // this.queryWorkflowExecutor(res.id)
      }).catch(() => {
        this.workflowIsExecutor = false
      })
      this.workflowIsExecutor = true
    },
    workflowStop () {
      clearTimeout(this.excuteTimer)
      clearTimeout(this.executorStatusTimer)
      this.FesApi.fetch(`/di/${this.FesEnv.diApiVersion}/experimentRun/${this.workflowExeteId}/kill`, 'get').then((res) => {
        this.workflowIsExecutor = false
        this.flowStatus.state = 'Cancelled'
        this.flowExecutorNode(this.workflowExeteId, true)
      }).catch(() => {
        this.workflowIsExecutor = false
      })
    },
    queryWorkflowExecutor (execID) {
      this.FesApi.fetch(`/di/${this.FesEnv.diApiVersion}/experimentRun/${execID}/status`, 'get').then((res) => {
        this.flowExecutorNode(execID)
        // 根据执行状态判断是否轮询
        const status = res.status
        this.flowStatus.state = status
        if (['Inited', 'Scheduled', 'Running'].includes(status)) {
          if (this.excuteTimer) {
            clearTimeout(this.excuteTimer)
            this.excuteTimer = null
          }
          this.excuteTimer = setTimeout(() => {
            this.queryWorkflowExecutor(execID)
          }, 1000)
        } else {
          // Succees, Failed, Cancelled, Timeout
          this.workflowIsExecutor = false
          if (status === 'Succeed') {
            this.$message.success(this.$t('flow.workflowRunSuccess'))
          }
          if (status === 'Failed') {
            this.$message.error(this.$t('flow.workflowRunFail'))
            this.flowExecutorNode(execID, true)
          }
          if (status === 'Cancelled') {
            this.$message.error(this.$t('flow.workflowRunCanceled'))
            this.flowExecutorNode(execID, true)
          }
          if (status === 'Timeout') {
            this.$message.error(this.$t('flow.workflowRunOvertime'))
            this.flowExecutorNode(execID, true)
          }
        }
      })
    },
    flowExecutorNode (execID, end = false) {
      this.FesApi.fetch(`/di/${this.FesEnv.diApiVersion}/experimentRun/${execID}/execution`, {}, 'get').then((res) => {
        // 【0：未执行；1：运行中；2：已成功；3：已失败；4：已跳过】
        const actionStatus = {
          pendingJobs: {
            color: 'rgba(102, 102, 102, 0.3)',
            status: 0,
            iconType: '',
            colorClass: { 'executor-icon': true },
            isShowTime: false,
            title: this.$t('flow.waitExecution'),
            showConsole: false
          },
          runningJobs: {
            color: 'rgba(69, 166, 239, 1)',
            status: 1,
            iconType: 'el-icon-loading',
            colorClass: { 'executor-loading': true, 'el-icon-loading': 'el-icon-loading', 'executor-icon': true },
            isShowTime: true,
            title: this.$t('flow.inExecution'),
            showConsole: true
          },
          succeedJobs: {
            color: 'rgba(78, 225, 78, 0.4)',
            status: 2,
            iconType: 'el-icon-circle-check',
            colorClass: { 'executor-success': true, 'el-icon-circle-check': 'el-icon-circle-check', 'executor-icon': true },
            isShowTime: false,
            title: this.$t('flow.executeSuccess'),
            showConsole: true
          },
          failedJobs: {
            color: 'rgba(231, 77, 98, 0.4)',
            status: 3,
            iconType: 'el-icon-circle-close',
            colorClass: { 'executor-faile': true, 'el-icon-circle-close': 'el-icon-circle-close', 'executor-icon': true },
            isShowTime: false,
            title: this.$t('flow.executFailure'),
            showConsole: true
          },
          skippedJobs: {
            color: 'rgba(236, 152, 47, 0.8)',
            status: 4,
            iconType: 'el-icon-warning-outline',
            colorClass: { 'executor-skip': true },
            isShowTime: false,
            title: this.$t('flow.skip'),
            showConsole: false
          }
        }
        const stateMapping = {
          pendingJobs: 'Inited',
          runningJobs: 'Running',
          succeedJobs: 'Succeed',
          failedJobs: 'Failed',
          skippedJobs: 'Cancelled'
        }
        // 获取节点的状态，如果没有执行完成继续查询
        const data = res
        Object.keys(data).forEach((key) => {
          if (!Array.isArray(data[key])) {
            return
          }
          if (data[key].length > 0) {
            this.nodeStatus.state = stateMapping[key]
          }
          // 如果当前工作流已经执行结束，还得获取状态到没有执行的节点为止
          if (end && key === 'runningJobs' && data[key].length > 0) {
            // 手动停掉执行和切换页面停止调接口
            this.executorStatusTimer = setTimeout(() => {
              this.flowExecutorNode(execID, true)
            }, 2000)
          }
          data[key].forEach((node) => {
            if (this.$refs.process) {
              const currentTime = new Date().getTime()
              const startTime = node.startTime !== 0 ? node.startTime : new Date().getTime()
              this.$refs.process.setNodeRunState(node.nodeID, {
                time: this.timeTransition(currentTime - startTime),
                status: actionStatus[key].status,
                borderColor: actionStatus[key].color,
                iconType: actionStatus[key].iconType,
                colorClass: actionStatus[key].colorClass,
                isShowTime: actionStatus[key].isShowTime,
                title: actionStatus[key].title,
                showConsole: actionStatus[key].showConsole,
                execID: node.execID,
                taskID: node.taskID
              })
            }
          })
        })
      })
    },
    message (obj) {
      let type = {
        warning: 'warn',
        error: 'error',
        info: 'info'
      }
      this.$Message[type[obj.type]]({
        content: obj.msg,
        duration: 2
      })
    },
    timeTransition (time) {
      time = Math.floor(time / 1000)
      let hour = 0
      let minute = 0
      let second = 0
      // if (time < 0) return str;
      if (time > 60) {
        minute = Math.floor(time / 60)
        second = Math.floor(time % 60)
        if (minute >= 60) {
          hour = Math.floor(minute / 60)
          minute = Math.floor(minute % 60)
        } else {
          hour = 0
        }
      } else {
        hour = 0
        if (time === 60) {
          minute = 1
          second = 0
        } else {
          minute = 0
          second = time
        }
      }
      const addZero = (num) => {
        let result = num
        if (num < 10) {
          result = `0${num}`
        }
        return result
      }
      const timeResult = `${addZero(hour)}: ${addZero(minute)}: ${addZero(second)}`
      time = 0
      return timeResult
    }
  }
}
</script>
<style src="./index.scss" lang="scss"></style>
