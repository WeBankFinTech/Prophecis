<template>
  <div>
    <breadcrumb-nav></breadcrumb-nav>
    <el-form ref="queryValidate" :model="query" :rules="ruleValidate" class="query" label-width="85px">
      <el-row>
        <el-col :span="6">
          <el-form-item :label="$t('ns.nameSpace')" prop="namespace">
            <el-select v-model="query.namespace" filterable :placeholder="$t('ns.nameSpacePro')">
              <el-option v-for="item in spaceOptionList" :key="item" :value="item" :label="item">
              </el-option>
            </el-select>
          </el-form-item>
        </el-col>
        <el-col :span="6">
          <el-form-item :label="$t('user.user')" prop="userName">
            <el-select v-model="query.userName" filterable clearable :placeholder="$t('user.userPro')">
              <el-option v-for="item in userOptionList" :key="item.id" :value="item.name" :label="item.name">
              </el-option>
            </el-select>
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-button type="primary" class="margin-l30" icon="el-icon-search" @click="filterListData">
            {{ $t("common.filter") }}
          </el-button>
          <el-button type="warning" icon="el-icon-refresh-left" @click="resetListData">
            {{ $t("common.reset") }}
          </el-button>
          <el-button type="primary" icon="plus-circle-o" @click="showCreateDialog">
            {{ $t("AIDE.createNotebook") }}
          </el-button>
        </el-col>
      </el-row>
    </el-form>
    <el-table :data="dataList" border style="width: 100%" :empty-text="$t('common.noData')">
      <el-table-column prop="name" :label="$t('DI.name')" min-width="100" />
      <el-table-column prop="user" :label="$t('user.user')" />
      <el-table-column prop="namespace" :label="$t('ns.nameSpace')" min-width="220" />
      <el-table-column prop="cpu" :label="$t('AIDE.cpu')" min-width="50" />
      <el-table-column prop="gpu" :label="$t('AIDE.gpu')" min-width="50" />
      <el-table-column prop="memory" :label="$t('AIDE.memory1')" min-width="70" />
      <el-table-column prop="image" :label="$t('DI.image')" min-width="300" />
      <el-table-column prop="status" :label="$t('DI.status')">
        <template slot-scope="scope">
          <status-tag :status="scope.row.status" :statusObj="statusObj"></status-tag>
        </template></el-table-column>
      <el-table-column prop="uptime" :label="$t('DI.createTime')" min-width="125" />
      <el-table-column :label="$t('common.operation')" min-width="85">
        <template slot-scope="scope">
          <el-dropdown size="medium" @command="handleCommand">
            <el-button type="text" class="multi-operate-btn" icon="el-icon-more" round></el-button>
            <el-dropdown-menu slot="dropdown">
              <el-dropdown-item :command="beforeHandleCommand('visitNotebook',scope.row)" :disabled="scope.row.status!=='Ready'">{{$t('AIDE.access')}}</el-dropdown-item>
              <el-dropdown-item :command="beforeHandleCommand('deleteItem',scope.row)">{{$t('common.delete')}}</el-dropdown-item>
              <el-dropdown-item :command="beforeHandleCommand('copyJob',scope.row)">{{$t('common.copy')}}</el-dropdown-item>
              <el-dropdown-item :command="beforeHandleCommand('editYarn',scope.row)">{{$t('AIDE.editYarn')}}</el-dropdown-item>
            </el-dropdown-menu>
          </el-dropdown>
        </template>
      </el-table-column>
    </el-table>
    <el-pagination @size-change="handleSizeChange" @current-change="handleCurrentChange" :current-page="pagination.pageNumber" :page-sizes="sizeList" :page-size="pagination.pageSize" layout="total, sizes, prev, pager, next, jumper" :total="pagination.totalPage" background>
    </el-pagination>
    <el-dialog :title="$t('AIDE.createNotebook')" :visible.sync="dialogVisible" @close="clearDialogForm1" custom-class="dialog-style" :close-on-click-modal="false" width="1150px">
      <step :step="currentStep" :step-list="stepList" />
      <el-form ref="formValidate" :model="form" class="add-node-form" :rules="ruleValidate" label-width="155px">
        <div v-if="currentStep===0" key="step0">
          <div class="subtitle">
            {{ $t('DI.basicSettings') }}
          </div>
          <el-row>
            <el-col :span="24">
              <el-form-item :label="$t('AIDE.NotebookName')" prop="name" maxlength="225">
                <el-input v-model="form.name" :placeholder="$t('AIDE.NotebookPro')" />
              </el-form-item>
            </el-col>
          </el-row>
        </div>
        <div v-if="currentStep===1" key="step1">
          <div class="subtitle">
            {{ $t('DI.imageSettings') }}
          </div>
          <el-row>
            <el-form-item :label="$t('DI.imageType')" prop="image.imageType">
              <el-radio-group v-model="form.image.imageType" @change="changeImageType">
                <el-radio label="Standard">{{ $t('DI.standard') }}</el-radio>
                <el-radio label="Custom">{{ $t('DI.custom') }}</el-radio>
              </el-radio-group>
            </el-form-item>
          </el-row>
          <el-row>
            <el-form-item v-if="!isCustom" :label="$t('DI.imageSelection')" prop="image.imageOption">
              <span>&nbsp;{{ defineImage }}&nbsp;&nbsp;</span>
              <el-select v-model="form.image.imageOption" :placeholder="$t('DI.imagePro')">
                <el-option v-for="(item,index) in imageOptionList" :label="item" :key="index" :value="item">
                </el-option>
              </el-select>
            </el-form-item>
            <el-form-item v-if="isCustom" key="imageInput" :label="$t('DI.imageSelection')" prop="image.imageInput">
              <span>&nbsp;{{ defineImage }}&nbsp;&nbsp;</span>
              <el-input v-model="form.image.imageInput" :placeholder="$t('DI.imageInputPro')" />
            </el-form-item>
          </el-row>
        </div>
        <div v-if="currentStep===2" key="step2">
          <div class="subtitle">
            {{ $t('DI.computingResource') }}
          </div>
          <el-row>
            <el-col :span="12">
              <el-form-item :label="$t('ns.nameSpace')" prop="namespace">
                <el-select v-model="form.namespace" filterable :placeholder="$t('ns.nameSpacePro')">
                  <el-option v-for="item in spaceOptionList" :label="item" :key="item" :value="item">
                  </el-option>
                </el-select>
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item :label="$t('AIDE.cpu')" prop="cpu">
                <el-input v-model="form.cpu" :placeholder="$t('AIDE.CPUPro')">
                  <span slot="append">Core</span>
                </el-input>
              </el-form-item>
            </el-col>
          </el-row>
          <el-row>
            <el-col :span="12">
              <el-form-item :label="$t('AIDE.memory')" prop="memory">
                <el-input v-model="form.memory" :placeholder="$t('AIDE.memoryPro')">
                  <span slot="append">Gi</span>
                </el-input>
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item :label="$t('AIDE.gpu')" prop="extraResources">
                <el-input v-model="form.extraResources" :placeholder="$t('AIDE.gpuPro')">
                  <span slot="append">{{$t('AIDE.block')}}</span>
                </el-input>
              </el-form-item>
            </el-col>
          </el-row>
          <el-form-item :label="$t('AIDE.yarnClusterResources')">
            <el-switch v-model="form.cluster">
            </el-switch>
          </el-form-item>
          <el-row v-if="form.cluster===true">
            <el-form-item :label="$t('DI.YARNQueue')" prop="queue">
              <el-input v-model="form.queue" maxlength="225" :placeholder="$t('DI.queuePro')" />
            </el-form-item>
            <el-col :span="12">
              <el-form-item :label="$t('DI.driverMemory')" prop="driverMemory">
                <el-input v-model="form.driverMemory" maxlength="10" :placeholder="$t('DI.driverMemoryPro')">
                  <span slot="append">Gi</span>
                </el-input>
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item :label="$t('DI.executorInstances')" prop="executorCores">
                <el-input v-model="form.executorCores" maxlength="10" :placeholder="$t('DI.executorInstancesPro')">
                  <span slot="append">Core</span>
                </el-input>
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item :label="$t('DI.executorMemory')" prop="executorMemory">
                <el-input v-model="form.executorMemory" maxlength="10" :placeholder="$t('DI.executorMemoryPro')">
                  <span slot="append">Gi</span>
                </el-input>
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item :label="$t('DI.linkisInstance')" prop="executors">
                <el-input v-model="form.executors" maxlength="10" :placeholder="$t('DI.linkisInstancePro')" />
              </el-form-item>
            </el-col>
          </el-row>
        </div>
        <div v-if="currentStep===3" key="step3">
          <el-row>
            <el-form-item :label="$t('AIDE.settingType')">
              <el-radio-group v-model="storageDefalut">
                <el-radio label="true">{{ $t('AIDE.default') }}</el-radio>
                <el-radio label="false">{{ $t('DI.custom') }}</el-radio>
              </el-radio-group>
            </el-form-item>
          </el-row>
          <div v-if="storageDefalut=='true'">
            <div class="subtitle">
              {{ $t("AIDE.workspace") }}
            </div>
            <el-row>
              <el-form-item :label="$t('DI.localPath')" prop="workspaceVolume.localPath">
                <!-- <el-input v-model="item.localPath" :placeholder="$t('AIDE.workLocalPathNodePro')" /> -->
                <el-select v-model="form.workspaceVolume.localPath" filterable :placeholder="$t('AIDE.workLocalPathNodePro')">
                  <el-option v-for="(item,index) in localPathList" :label="item" :key="index" :value="item">
                  </el-option>
                </el-select>
              </el-form-item>
              <el-form-item :label="$t('AIDE.mountPath')" prop="workspaceVolume.mountPath">
                <el-input v-model="form.workspaceVolume.mountPath" :placeholder="$t('AIDE.mountPathPro')" disabled />
              </el-form-item>
            </el-row>
          </div>
          <div class="subtitle data-space" v-if="storageDefalut=='false'">
            {{ $t("AIDE.dataStorageSettings") }}
          </div>
          <div v-if="storageDefalut=='false'">
            <div v-for="(item,index) in form.dataVolume" :key="`dataVolume${index}`">
              <el-row>
                <el-form-item :label="$t('DI.localPath')" :prop="`dataVolume[${index}].localPath`">
                  <el-select v-model="item.localPath" filterable :placeholder="$t('AIDE.dataLocalPathNodePro')">
                    <el-option v-for="(item,index) in localPathList" :label="item" :key="index" :value="item">
                    </el-option>
                  </el-select>
                </el-form-item>
                <el-form-item :label="$t('AIDE.subpath')" :prop="`dataVolume[${index}].subPath`">
                  <el-input v-model="item.subPath" :placeholder="$t('AIDE.subpathPro')" />
                </el-form-item>
                <el-form-item :label="$t('AIDE.mountPath')" :prop="`dataVolume[${index}].mountPath`">
                  <el-input v-model="item.mountPath" :placeholder="$t('AIDE.dataMountPathPro')" />
                </el-form-item>
              </el-row>
              <el-row>
                <el-col :span="12">
                  <el-form-item :label="$t('AIDE.mountType')" :prop="`dataVolume[${index}].mountType`">
                    <el-select v-model="item.mountType" :placeholder="$t('AIDE.mountTypePro')">
                      <el-option value="New" :label="$t('AIDE.nodeDirectory')">
                      </el-option>
                    </el-select>
                  </el-form-item>
                </el-col>
                <el-col :span="12">
                  <el-form-item :label="$t('AIDE.accessMode')" :prop="`dataVolume[${index}].accessMode`">
                    <el-select v-model="item.accessMode" :placeholder="$t('AIDE.accessModePro')">
                      <el-option value="ReadOnlyMany" :label="$t('AIDE.readOnly')">
                      </el-option>
                      <el-option value="ReadWriteMany" :label="$t('AIDE.readWrite')">
                      </el-option>
                    </el-select>
                  </el-form-item>
                </el-col>
              </el-row>
            </div>
          </div>

        </div>
      </el-form>
      <span slot="footer" class="dialog-footer">
        <el-button v-if="currentStep>0" type="primary" @click="lastStep">
          {{ $t('common.lastStep') }}
        </el-button>
        <el-button v-if="currentStep<3" type="primary" @click="nextStep">
          {{ $t('common.nextStep') }}
        </el-button>
        <el-button v-if="currentStep===3" :disabled="btnDisabled" type="primary" @click="subInfo">
          {{ $t('common.save') }}
        </el-button>
        <el-button @click="dialogVisible=false">
          {{ $t('common.cancel') }}
        </el-button>
      </span>
    </el-dialog>

    <el-dialog :title="$t('AIDE.yarnClusterResources')" :visible.sync="yarnDialog" custom-class="dialog-style" :close-on-click-modal="false" width="1100px">
      <el-form ref="yarnForm" :rules="ruleValidate" :model="yarnForm" label-width="190px">
        <el-row>
          <el-form-item :label="$t('DI.YARNQueue')" prop="queue">
            <el-input v-model="yarnForm.queue" maxlength="225" :placeholder="$t('DI.queuePro')" />
          </el-form-item>
          <el-col :span="12">
            <el-form-item :label="$t('DI.driverMemory')" prop="driverMemory">
              <el-input v-model="yarnForm.driverMemory" maxlength="10" :placeholder="$t('DI.driverMemoryPro')">
                <span slot="append">Gi</span>
              </el-input>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item :label="$t('DI.executorInstances')" prop="executorCores">
              <el-input v-model="yarnForm.executorCores" maxlength="10" :placeholder="$t('DI.executorInstancesPro')">
                <span slot="append">Core</span>
              </el-input>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item :label="$t('DI.executorMemory')" prop="executorMemory">
              <el-input v-model="yarnForm.executorMemory" maxlength="10" :placeholder="$t('DI.executorMemoryPro')">
                <span slot="append">Gi</span>
              </el-input>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item :label="$t('DI.linkisInstance')" prop="executors">
              <el-input v-model="yarnForm.executors" maxlength="10" :placeholder="$t('DI.linkisInstancePro')" />
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
      <span slot="footer" class="dialog-footer">
        <el-button type="primary" @click="saveYarnForm">
          {{ $t('common.save') }}
        </el-button>
        <el-button @click="cannelSaveYarnForm">
          {{ $t('common.cancel') }}
        </el-button>
      </span>
    </el-dialog>
  </div>
</template>
<script type="text/ecmascript-6">
import util from '../../util/common'
import Step from '../../components/Step.vue'
import modelMixin from '../../util/modelMixin'
import StatusTag from '../../components/StatusTag'
export default {
  mixins: [modelMixin],
  components: {
    Step,
    StatusTag
  },
  data () {
    return {
      dialogVisible: false,
      query: {
        namespace: '',
        userName: ''
      },
      currentStep: 0,
      form: {
        name: '',
        cpu: '',
        extraResources: '',
        memory: '',
        queue: '',
        namespace: '',
        driverMemory: '',
        executorCores: '',
        executorMemory: '',
        executors: '',
        cluster: false,
        image: {
          'imageType': 'Standard',
          'imageOption': '',
          'imageInput': ''
        },
        workspaceVolume: {
          accessMode: '',
          localPath: '',
          subPath: '',
          mountPath: '',
          mountType: '',
          size: 0
        },
        dataVolume: [{
          accessMode: '',
          localPath: '',
          subPath: '',
          mountPath: '',
          mountType: '',
          size: 0
        }]
      },
      localPathList: [],
      loading: false,
      dataList: [],
      userOptionList: [],
      spaceOptionList: [],
      yarnDialog: false,
      yarnForm: {
        cluster: true,
        queue: '',
        driverMemory: '',
        executorCores: '',
        executorMemory: '',
        executors: ''
      },
      storageDefalut: 'true',
      statusObj: {
        Ready: 'success',
        Waiting: '',
        Terminate: 'danger'
      },
      setIntervalList: null
    }
  },
  mounted () {
    this.getListData()
    this.getDIUserOption()
    this.getDISpaceOption()
    this.getDIStoragePath()
    this.setIntervalList = setInterval(this.getListData, 60000)
  },
  computed: {
    noDataText () {
      return !this.loading ? this.$t('common.noData') : ''
    },
    stepList () {
      return [this.$t('DI.basicSettings'), this.$t('DI.imageSettings'), this.$t('DI.computingResource'), this.$t('AIDE.storageSettings')]
    },
    defineImage () {
      return this.FesEnv.AIDE.defineImage
    },
    imageOptionList () {
      return this.FesEnv.AIDE.imageOption
    },
    userId () {
      return localStorage.getItem('userId')
    },
    isCustom () {
      return this.form.image.imageType === 'Custom'
    },
    ruleValidate () {
      // 切换语言表单报错重置表单
      this.resetQueryFields()
      this.resetFormFields()
      return {
        name: [
          { required: true, message: this.$t('AIDE.NotebookNameReq') },
          { type: 'string', pattern: new RegExp(/^[a-z][a-z0-9-]*$/), message: this.$t('AIDE.NotebookNameFormatReq') }
        ],
        namespace: [
          { required: true, message: this.$t('ns.nameSpaceReq') }
        ],
        'image.imageType': [
          { required: true, message: this.$t('DI.imageTypeReq') }
        ],
        'image.imageOption': [
          { required: true, message: this.$t('DI.imageColumnReq') }
        ],
        'image.imageInput': [
          { required: true, message: this.$t('DI.imageColumnReq') },
          { pattern: new RegExp(/^[a-zA-Z0-9][a-zA-Z0-9-._]*$/), message: this.$t('DI.imageInputFormat') }
        ],
        cpu: [
          { required: true, message: this.$t('DI.CPUReq') },
          { pattern: new RegExp(/^((0{1})([.]\d{1})|([1-9]\d*)([.]\d{1})?)$/), message: this.$t('AIDE.CPUNumberReq') }
        ],
        memory: [
          { required: true, message: this.$t('DI.memoryReq') },
          { pattern: new RegExp(/^[1-9]\d*([.]0)*$/), message: this.$t('AIDE.memoryNumberReq') }
        ],
        queue: [
          { required: true, message: this.$t('DI.queueReq') },
          { pattern: new RegExp(/^[a-zA-Z0-9][a-zA-Z0-9_.]*$/), message: this.$t('DI.queueFormat') }
        ],
        driverMemory: [
          { pattern: new RegExp(/^[1-9]\d*([.]0)*$/), message: this.$t('DI.driverMemoryFormat') },
          { pattern: new RegExp(/^.{0,10}$/), message: this.$t('AIDE.driverMemoryMax') }
        ],
        executorCores: [
          { pattern: new RegExp(/^[1-9]\d*([.]0)*$/), message: this.$t('DI.executorInstancesFormat') },
          { pattern: new RegExp(/^.{0,10}$/), message: this.$t('AIDE.executorInstancesMax') }
        ],
        executors: [
          { pattern: new RegExp(/^[1-9]\d*([.]0)*$/), message: this.$t('AIDE.executorFormat') },
          { pattern: new RegExp(/^.{0,10}$/), message: this.$t('AIDE.executorsMax') }
        ],
        executorMemory: [
          { pattern: new RegExp(/^[1-9]\d*([.]0)*$/), message: this.$t('DI.executorMemoryFormat') },
          { pattern: new RegExp(/^.{0,10}$/), message: this.$t('AIDE.executorMemoryMax') }
        ],
        'workspaceVolume.localPath': [
          { required: true, message: this.$t('DI.localPathReq') },
          { pattern: new RegExp(/^\/\w[0-9a-zA-Z-_/]*$/), message: this.$t('AIDE.localPathFormat') }
        ],
        'workspaceVolume.mountPath': [
          { required: true, message: this.$t('AIDE.mountPathReq') }
        ],
        'dataVolume[0].localPath': [
          { required: true, message: this.$t('DI.localPathReq') },
          { pattern: new RegExp(/^\/\w[0-9a-zA-Z-_/]*$/), message: this.$t('AIDE.localPathFormat') }
        ],
        'dataVolume[0].subPath': [
          { pattern: new RegExp(/^[a-zA-Z][0-9a-zA-Z-_/]*$/), message: this.$t('AIDE.subpathReq') }
        ],
        'dataVolume[0].mountPath': [
          { required: true, message: this.$t('AIDE.mountPathReq') },
          { pattern: new RegExp(/^\/\w[0-9a-zA-Z-_/]*$/), message: this.$t('AIDE.dataMountPathFormat') }
        ],
        'dataVolume[0].mountType': [
          { required: true, message: this.$t('AIDE.mountTypeReq') }
        ],
        'dataVolume[0].accessMode': [
          { required: true, message: this.$t('AIDE.accessModeReq') }
        ],
        extraResources: [
          { pattern: new RegExp(/^(([1-9]{1}\d*([.]0)*)|([0]{1}))$/), message: this.$t('AIDE.GPUFormat') }
        ]
      }
    }
  },
  methods: {
    showCreateDialog () {
      this.form.workspaceVolume.mountPath = `/home/${this.userId}/workspace`
      this.dialogVisible = true
    },
    visitNotebook (trData) {
      let url = `${process.env.VUE_APP_BASE_SURL}/cc/${this.FesEnv.ccApiVersion}/auth/access/namespaces/${trData.namespace}/notebooks/${trData.name}`
      let xhr = new XMLHttpRequest()
      // const mlssToken = util.getCookieVal(localStorage.getItem('cookieKey')) || 0
      xhr.open('GET', url, false)
      xhr.setRequestHeader('Mlss-Userid', this.userId)
      // xhr.setRequestHeader('Mlss-Token', mlssToken)
      xhr.onreadystatechange = function (e) {
        if (this.readyState === 4 && this.status === 200) {
          let rst = JSON.parse(this.responseText).result
          const url = process.env.VUE_APP_BASE_SURL + window.encodeURI(rst.notebookAddress)
          rst.notebookAddress && window.open(url)
        } else {
          let rst = JSON.parse(this.responseText)
          if (this.status === 401 || this.status === 403) {
            this.fetchError[this.status](rst)
          } else {
            this.$Message.error(rst.message)
          }
        }
      }
      xhr.send()
    },
    deleteItem (trData) {
      const url = `/aide/${this.FesEnv.aideApiVersion}/namespaces/${trData.namespace}/notebooks/${trData.name}`
      this.deleteListItem(url)
    },
    copyJob (trData) {
      const job = {}
      const formKey = ['cpu', 'namespace', 'queue', 'cluster', 'executorCores', 'executorMemory', 'executors', 'driverMemory']
      for (let item of formKey) {
        job[item] = trData[item]
      }
      job.memory = parseInt(trData.memory)
      job.extraResources = trData.gpu + ''
      const imageType = this.imageOptionList.indexOf(trData.image) > -1 ? 'Standard' : 'Custom'
      job.image = {
        'imageType': imageType,
        imageOption: '',
        imageInput: ''
      }

      job.executorMemory = isNaN(parseInt(job.executorMemory)) ? '' : parseInt(job.executorMemory) + ''
      job.driverMemory = isNaN(parseInt(job.driverMemory)) ? '' : parseInt(job.driverMemory) + ''
      imageType === 'Standard' ? job.image.imageOption = trData.image : job.image.imageInput = trData.image
      setTimeout(() => {
        Object.assign(this.form, job)
        if (trData.dataVolume) {
          const dataVolume = this.dataVolumeInfo()
          for (let index in trData.dataVolume) {
            if (index > 0) {
              this.addValidateItem(dataVolume.column, index, 'dataVolume')
            }
            this.form.dataVolume[index] = trData.dataVolume[index]
          }
          // 存储设置类型 为 自定义类型
          this.storageDefalut = 'false'
          this.form.workspaceVolume.mountPath = `/home/${this.userId}/workspace`
        } else {
          this.form.workspaceVolume = trData.workspaceVolume
          // 存储设置类型 为 默认
          this.storageDefalut = 'true'
        }
      }, 100)
      this.dialogVisible = true
    },
    nextStep () {
      this.$refs.formValidate.validate((valid) => {
        if (valid && this.currentStep < this.stepList.length - 1) {
          this.currentStep++
        }
      })
    },
    getListData () {
      this.loading = true
      let url = ''
      let paginationOption = {
        page: this.pagination.pageNumber,
        size: this.pagination.pageSize
      }
      if (this.query.namespace && this.query.userName) {
        url = `/aide/${this.FesEnv.aideApiVersion}/namespaces/${this.query.namespace}/user/${this.query.userName}/notebooks`
      } else if (this.query.namespace && !this.query.userName) {
        url = `/aide/${this.FesEnv.aideApiVersion}/namespaces/${this.query.namespace}/notebooks`
      } else {
        let superAdmin = localStorage.getItem('superAdmin')
        if (superAdmin === 'true') {
          url = `/aide/${this.FesEnv.aideApiVersion}/namespaces/null/user/null/notebooks`
        } else {
          let userName = this.query.userName ? this.query.userName : this.userId
          url = `/aide/${this.FesEnv.aideApiVersion}/user/${userName}/notebooks`
        }
      }
      this.FesApi.fetch(url, paginationOption, 'get').then(rst => {
        if (rst && JSON.stringify(rst) !== '{}') {
          rst.list && rst.list.forEach(item => {
            let index = item.image.indexOf(':')
            item.image = item.image.substring(index + 1)
            item.gpu = item.gpu || 0
            item.uptime = item.uptime ? util.transDate(item.uptime) : ''
            if (item.queue) {
              item.cluster = true
            }
          })
          this.dataList = Object.freeze(rst.list)
          this.pagination.totalPage = rst.total || 0
        }
        this.loading = false
      }, () => {
        this.dataList = []
        this.loading = false
      })
    },
    subInfo () {
      this.$refs.formValidate.validate((valid) => {
        if (valid) {
          this.setBtnDisabeld()
          let param = {}
          let keyArr = ['name', 'namespace']
          keyArr = this.form.cluster ? keyArr.concat(['queue', 'executorCores', 'executorMemory', 'executors', 'driverMemory']) : keyArr
          for (let key of keyArr) {
            if (key === 'name' || key === 'namespace') {
              param[key] = this.form[key]
            } else {
              if (this.form[key]) {
                param[key] = this.form[key]
              }
            }
          }
          param.cpu = parseFloat(this.form.cpu)
          param.memory = {
            memoryAmount: parseInt(this.form.memory),
            memoryUnit: 'Gi'
          }
          param.extraResources = JSON.stringify({ 'nvidia.com/gpu': this.form.extraResources })
          let imageColumn = this.form.image.imageOption || this.form.image.imageInput
          let imageName = this.defineImage + ':' + imageColumn
          param.image = { 'imageName': imageName }
          param.image.imageType = this.form.image.imageType
          // 设置类型为自定义类型时
          if (this.storageDefalut === 'true') {
            param.workspaceVolume = { ...this.form.workspaceVolume, ...{ accessMode: 'ReadWriteMany', mountType: 'New' } }
          } else {
            param.dataVolume = { ...this.form.dataVolume[0] }
          }
          this.FesApi.fetch(`/aide/${this.FesEnv.aideApiVersion}/namespaces/namespace/notebooks`, param, 'post').then(() => {
            this.getListData()
            this.toast()
            this.dialogVisible = false
          })
        }
      })
    },
    dataVolumeInfo () {
      return {
        lg: this.form.dataVolume.length,
        column: ['accessMode', 'localPath', 'subPath', 'mountPath', 'mountType']
      }
    },
    changeImageType (value) {
      // 镜像类型为自定义时，清空镜像列下拉
      if (value === 'Custom') {
        this.form.image.imageOption = ''
      } else {
        this.form.image.imageInput = ''
      }
    },
    addValidateItem (column, columnLength, objKey) {
      for (let item of column) {
        this.ruleValidate[`${objKey}[${columnLength}].${item}`] = this.ruleValidate[`${objKey}[${columnLength - 1}].${item}`]
      }
    },
    deleteValidateItem (column, columnLength, objKey) {
      for (let item of column) {
        delete this.ruleValidate[`${objKey}[${columnLength - 1}].${item}`]
      }
    },
    clearDialogForm1 () {
      const dataVolume = this.dataVolumeInfo()
      if (dataVolume.lg > 1) {
        for (let i = dataVolume.lg; i > 1; i--) {
          this.deleteValidateItem(dataVolume.column, i, 'dataVolume')
        }
      }
      this.$refs.formValidate.resetFields()
      const formKey = ['name', 'cpu', 'extraResources', 'memory', 'namespace', 'queue', 'driverMemory', 'executorCores', 'executorMemory', 'executors']
      for (let item of formKey) {
        this.form[item] = ''
      }
      for (let key in this.form.image) {
        this.form.image[key] = key === 'imageType' ? 'Standard' : ''
      }
      this.form.workspaceVolume = {
        accessMode: '',
        localPath: '',
        subPath: '',
        mountPath: `/home/${this.userId}/workspace`,
        mountType: '',
        size: 0
      }
      this.form.dataVolume = [{
        accessMode: '',
        localPath: '',
        subPath: '',
        mountPath: '',
        mountType: '',
        size: 0
      }]
      this.form.cluster = false
      setTimeout(() => {
        this.currentStep = 0
        this.storageDefalut = 'true'
        console.log('this.form', this.form)
      }, 1000)
    },
    editYarn (trData) {
      this.yarnDialog = true
      this.yarnForm = {
        name: trData.name,
        namespace: trData.namespace,
        queue: trData.queue || '',
        driverMemory: isNaN(parseInt(trData.driverMemory)) ? '' : parseInt(trData.driverMemory) + '',
        executorCores: isNaN(parseInt(trData.executorCores)) ? '' : parseInt(trData.executorCores) + '',
        executorMemory: isNaN(parseInt(trData.executorMemory)) ? '' : parseInt(trData.executorMemory) + '',
        executors: isNaN(parseInt(trData.executors)) ? '' : parseInt(trData.executors) + ''
      }
    },
    saveYarnForm () {
      const _this = this
      let parms = { ...this.yarnForm }
      parms.driverMemory = Number(parms.driverMemory) > 0 ? parms.driverMemory + 'g' : ''
      parms.executorMemory = Number(parms.executorMemory) > 0 ? parms.executorMemory + 'g' : ''
      this.$refs.yarnForm.validate((valid) => {
        if (valid) {
          this.FesApi.fetch(`/aide/${this.FesEnv.aideApiVersion}/namespaces/${_this.yarnForm.namespace}/notebooks`, parms, {
            method: 'patch',
            headers: {
              'Mlss-Userid': this.userId
            }
          }).then(() => {
            this.getListData()
            this.toast()
            this.yarnDialog = false
          })
        }
      })
    },
    cannelSaveYarnForm () {
      this.yarnDialog = false
      this.$refs.yarnForm.resetFormFields()
    }
  },
  beforeDestroy () {
    clearInterval(this.setIntervalList)
  }
}
</script>
<style lang="scss" scoped>
.add-node-form {
  min-height: 500px;
}
</style>
