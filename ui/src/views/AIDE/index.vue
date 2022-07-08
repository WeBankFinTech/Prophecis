<template>
  <div>
    <breadcrumb-nav></breadcrumb-nav>
    <el-form ref="queryValidate"
             :model="query"
             :rules="ruleValidate"
             class="query"
             label-width="85px">
      <el-row>
        <el-col :span="6">
          <el-form-item :label="$t('ns.nameSpace')">
            <el-select v-model="query.namespace"
                       filterable
                       clearable
                       :placeholder="$t('ns.nameSpacePro')">
              <el-option v-for="item in spaceOptionList"
                         :key="item"
                         :value="item"
                         :label="item">
              </el-option>
            </el-select>
          </el-form-item>
        </el-col>
        <el-col :span="6">
          <el-form-item :label="$t('user.user')"
                        prop="userName">
            <el-select v-model="query.userName"
                       filterable
                       clearable
                       :placeholder="$t('user.userPro')">
              <el-option v-for="item in userOptionList"
                         :key="item.id"
                         :value="item.name"
                         :label="item.name">
              </el-option>
            </el-select>
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-button type="primary"
                     class="margin-l30"
                     icon="el-icon-search"
                     @click="filterListData">
            {{ $t("common.filter") }}
          </el-button>
          <el-button type="warning"
                     icon="el-icon-refresh-left"
                     @click="resetListData">
            {{ $t("common.reset") }}
          </el-button>
          <el-button type="primary"
                     icon="plus-circle-o"
                     @click="showCreateDialog">
            {{ $t("AIDE.createNotebook") }}
          </el-button>
        </el-col>
      </el-row>
    </el-form>
    <el-table :data="dataList"
              border
              style="width: 100%"
              :empty-text="$t('common.noData')">
      <el-table-column prop="name"
                       :show-overflow-tooltip="true"
                       :label="$t('DI.name')"
                       min-width="100" />
      <el-table-column prop="user"
                       :show-overflow-tooltip="true"
                       :label="$t('user.user')" />
      <el-table-column prop="namespace"
                       :show-overflow-tooltip="true"
                       :label="$t('ns.nameSpace')"
                       min-width="200" />
      <el-table-column prop="cpu"
                       :show-overflow-tooltip="true"
                       :label="$t('AIDE.cpu')"
                       min-width="50" />
      <el-table-column prop="gpu"
                       :show-overflow-tooltip="true"
                       :label="$t('AIDE.gpu')"
                       min-width="50" />
      <el-table-column prop="memory"
                       :show-overflow-tooltip="true"
                       :label="$t('AIDE.memory1')"
                       min-width="70" />
      <el-table-column prop="image"
                       :show-overflow-tooltip="true"
                       :label="$t('DI.image')"
                       min-width="200" />
      <el-table-column prop="status"
                       :show-overflow-tooltip="true"
                       :label="$t('DI.status')">
        <template slot-scope="scope">
          <el-tag :class="{pointer:scope.row.status!=='Stop'}"
                  :type="statusObj[scope.row.status]"
                  @click="nodeBookStatus(scope.row)">{{scope.row.status}}</el-tag>
        </template></el-table-column>
      <el-table-column prop="uptime"
                       :show-overflow-tooltip="true"
                       :label="$t('DI.createTime')"
                       min-width="100" />
      <el-table-column :label="$t('common.operation')"
                       min-width="75">
        <template slot-scope="scope">
          <el-dropdown size="medium"
                       @command="handleCommand">
            <el-button type="text"
                       class="multi-operate-btn"
                       icon="el-icon-more"
                       round></el-button>
            <el-dropdown-menu slot="dropdown">
              <el-dropdown-item :command="beforeHandleCommand('visitNotebook',scope.row)"
                                :disabled="scope.row.status!=='Ready'">{{$t('AIDE.access')}}</el-dropdown-item>
              <el-dropdown-item :command="beforeHandleCommand('deleteJob',scope.row)">{{$t('common.delete')}}</el-dropdown-item>
              <el-dropdown-item :command="beforeHandleCommand('copyJob',scope.row)">{{$t('common.copy')}}</el-dropdown-item>
              <el-dropdown-item v-if="scope.row.status==='Ready'"
                                :command="beforeHandleCommand('getLog',scope.row)">{{$t('DI.viewlog')}}</el-dropdown-item>
              <el-dropdown-item v-if="scope.row.status!=='Terminated'&&scope.row.status!=='Creating'&&scope.row.status!=='Stopping'&&scope.row.status!=='Stopped'"
                                :command="beforeHandleCommand('editYarn',scope.row)">{{$t('AIDE.setting')}}</el-dropdown-item>
              <el-dropdown-item v-if="scope.row.status==='Ready'||scope.row.status==='Waiting'"
                                :command="beforeHandleCommand('disableJob',scope.row)">{{$t('AIDE.stop')}}</el-dropdown-item>
              <el-dropdown-item v-if="scope.row.status==='Stopped'"
                                :command="beforeHandleCommand('enableJob',scope.row)">{{$t('AIDE.enable')}}</el-dropdown-item>
            </el-dropdown-menu>
          </el-dropdown>
        </template>
      </el-table-column>
    </el-table>
    <el-pagination @size-change="handleSizeChange"
                   @current-change="handleCurrentChange"
                   :current-page="pagination.pageNumber"
                   :page-sizes="sizeList"
                   :page-size="pagination.pageSize"
                   layout="total, sizes, prev, pager, next, jumper"
                   :total="pagination.totalPage"
                   background>
    </el-pagination>
    <el-dialog :title="$t('AIDE.createNotebook')"
               :visible.sync="dialogVisible"
               @close="clearDialogForm1"
               custom-class="dialog-style"
               :close-on-click-modal="false"
               width="1150px">
      <step :step="currentStep"
            :step-list="stepList" />
      <el-form ref="formValidate"
               :model="form"
               class="add-node-form"
               :rules="ruleValidate"
               label-width="155px">
        <div v-if="currentStep===0"
             key="step0">
          <div class="subtitle">
            {{ $t('DI.basicSettings') }}
          </div>
          <el-form-item :label="$t('AIDE.NotebookName')"
                        prop="name"
                        maxlength="225">
            <el-input v-model="form.name"
                      :placeholder="$t('AIDE.NotebookPro')" />
          </el-form-item>
        </div>
        <div v-if="currentStep===1"
             key="step1">
          <div class="subtitle">
            {{ $t('DI.imageSettings') }}
          </div>
          <el-row>
            <el-form-item :label="$t('DI.imageType')"
                          prop="image.imageType">
              <el-radio-group v-model="form.image.imageType"
                              @change="changeImageType">
                <el-radio label="Standard">{{ $t('DI.standard') }}</el-radio>
                <el-radio label="Custom">{{ $t('DI.custom') }}</el-radio>
              </el-radio-group>
            </el-form-item>
          </el-row>
          <el-row>
            <el-form-item v-if="!isCustom"
                          :label="$t('DI.imageSelection')"
                          prop="image.imageOption">
              <span>&nbsp;{{ defineImage }}&nbsp;&nbsp;</span>
              <el-select v-model="form.image.imageOption"
                         filterable
                         :placeholder="$t('DI.imagePro')">
                <el-option v-for="(item,index) in imageOptionList"
                           :label="item"
                           :key="index"
                           :value="item">
                </el-option>
              </el-select>
            </el-form-item>
            <el-form-item v-if="isCustom"
                          key="imageInput"
                          :label="$t('DI.imageSelection')"
                          prop="image.imageInput">
              <span>&nbsp;{{ defineImage }}&nbsp;&nbsp;</span>
              <el-input v-model="form.image.imageInput"
                        :placeholder="$t('DI.imageInputPro')" />
            </el-form-item>
          </el-row>
        </div>
        <div v-if="currentStep===2"
             key="step2">
          <div class="subtitle">
            {{ $t('DI.computingResource') }}
          </div>
          <el-row>
            <el-col :span="12">
              <el-form-item :label="$t('ns.nameSpace')"
                            prop="namespace">
                <el-select v-model="form.namespace"
                           filterable
                           :placeholder="$t('ns.nameSpacePro')">
                  <el-option v-for="item in spaceOptionList"
                             :label="item"
                             :key="item"
                             :value="item">
                  </el-option>
                </el-select>
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item :label="$t('AIDE.cpu')"
                            prop="cpu">
                <el-input v-model="form.cpu"
                          :placeholder="$t('AIDE.CPUPro')">
                  <span slot="append">Core</span>
                </el-input>
              </el-form-item>
            </el-col>
          </el-row>
          <el-row>
            <el-col :span="12">
              <el-form-item :label="$t('AIDE.memory')"
                            prop="memory">
                <el-input v-model="form.memory"
                          :placeholder="$t('AIDE.memoryPro')">
                  <span slot="append">Gi</span>
                </el-input>
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item :label="$t('AIDE.gpu')"
                            prop="extraResources">
                <el-input v-model="form.extraResources"
                          :placeholder="$t('AIDE.gpuPro')">
                  <span slot="append">{{$t('AIDE.block')}}</span>
                </el-input>
              </el-form-item>
            </el-col>
          </el-row>
          <el-form-item :label="$t('AIDE.sparkResourceSettings')">
            <el-switch v-model="form.cluster">
            </el-switch>
          </el-form-item>
          <el-row v-if="form.cluster===true">
            <el-row>
              <el-col :span="12">
                <el-form-item :label="$t('AIDE.instanceNumber')"
                              prop="sparkSessionNum">
                  <el-input v-model="form.sparkSessionNum"
                            maxlength="225"
                            :placeholder="$t('AIDE.instanceNumberPro')" />
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item :label="$t('AIDE.queueSettings')"
                              prop="queue">
                  <el-input v-model="form.queue"
                            maxlength="225"
                            :placeholder="$t('DI.queuePro')" />
                </el-form-item>
              </el-col>
            </el-row>
            <el-col :span="12">
              <el-form-item :label="$t('DI.driverMemory')"
                            prop="driverMemory">
                <el-input v-model="form.driverMemory"
                          maxlength="10"
                          :placeholder="$t('DI.driverMemoryPro')">
                  <span slot="append">Gi</span>
                </el-input>
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item :label="$t('DI.executorInstances')"
                            prop="executorCores">
                <el-input v-model="form.executorCores"
                          maxlength="10"
                          :placeholder="$t('DI.executorInstancesPro')">
                  <span slot="append">Core</span>
                </el-input>
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item :label="$t('DI.executorMemory')"
                            prop="executorMemory">
                <el-input v-model="form.executorMemory"
                          maxlength="10"
                          :placeholder="$t('DI.executorMemoryPro')">
                  <span slot="append">Gi</span>
                </el-input>
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item :label="$t('DI.linkisInstance')"
                            prop="executors">
                <el-input v-model="form.executors"
                          maxlength="10"
                          :placeholder="$t('DI.linkisInstancePro')" />
              </el-form-item>
            </el-col>
          </el-row>
        </div>
        <div v-if="currentStep===3"
             key="step3">
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
              <el-form-item :label="$t('DI.localPath')"
                            prop="workspaceVolume.localPath">
                <!-- <el-input v-model="item.localPath" :placeholder="$t('AIDE.workLocalPathNodePro')" /> -->
                <el-select v-model="form.workspaceVolume.localPath"
                           :disabled="haveProxy"
                           filterable
                           :placeholder="$t('AIDE.workLocalPathNodePro')">
                  <el-option v-for="(item,index) in localPathList"
                             :label="item"
                             :key="index"
                             :value="item">
                  </el-option>
                </el-select>
              </el-form-item>
              <el-form-item :label="$t('AIDE.mountPath')"
                            prop="workspaceVolume.mountPath">
                <el-input v-model="form.workspaceVolume.mountPath"
                          :placeholder="$t('AIDE.mountPathPro')"
                          disabled />
              </el-form-item>
            </el-row>
          </div>
          <div class="subtitle data-space"
               v-if="storageDefalut=='false'">
            {{ $t("AIDE.dataStorageSettings") }}
          </div>
          <div v-if="storageDefalut=='false'">
            <div v-for="(item,index) in form.dataVolume"
                 :key="`dataVolume${index}`">
              <el-row>
                <el-form-item :label="$t('DI.localPath')"
                              :prop="`dataVolume[${index}].localPath`">
                  <el-select v-model="item.localPath"
                             filterable
                             :disabled="haveProxy"
                             :placeholder="$t('AIDE.dataLocalPathNodePro')">
                    <el-option v-for="(item,index) in localPathList"
                               :label="item"
                               :key="index"
                               :value="item">
                    </el-option>
                  </el-select>
                </el-form-item>
                <el-form-item :label="$t('AIDE.subpath')"
                              :prop="`dataVolume[${index}].subPath`">
                  <el-input v-model="item.subPath"
                            :placeholder="$t('AIDE.subpathPro')" />
                </el-form-item>
                <el-form-item :label="$t('AIDE.mountPath')"
                              :prop="`dataVolume[${index}].mountPath`">
                  <el-input v-model="item.mountPath"
                            :placeholder="$t('AIDE.dataMountPathPro')" />
                </el-form-item>
              </el-row>
              <el-row>
                <el-col :span="12">
                  <el-form-item :label="$t('AIDE.mountType')"
                                :prop="`dataVolume[${index}].mountType`">
                    <el-select v-model="item.mountType"
                               :placeholder="$t('AIDE.mountTypePro')">
                      <el-option value="New"
                                 :label="$t('AIDE.nodeDirectory')">
                      </el-option>
                    </el-select>
                  </el-form-item>
                </el-col>
                <el-col :span="12">
                  <el-form-item :label="$t('AIDE.accessMode')"
                                :prop="`dataVolume[${index}].accessMode`">
                    <el-select v-model="item.accessMode"
                               :placeholder="$t('AIDE.accessModePro')">
                      <el-option value="ReadOnlyMany"
                                 :label="$t('AIDE.readOnly')">
                      </el-option>
                      <el-option value="ReadWriteMany"
                                 :label="$t('AIDE.readWrite')">
                      </el-option>
                    </el-select>
                  </el-form-item>
                </el-col>
              </el-row>
            </div>
          </div>

        </div>
      </el-form>
      <span slot="footer"
            class="dialog-footer">
        <el-button v-if="currentStep>0"
                   type="primary"
                   @click="lastStep">
          {{ $t('common.lastStep') }}
        </el-button>
        <el-button v-if="currentStep<3"
                   type="primary"
                   @click="nextStep">
          {{ $t('common.nextStep') }}
        </el-button>
        <el-button v-if="currentStep===3"
                   :disabled="btnDisabled"
                   type="primary"
                   @click="subInfo">
          {{ $t('common.save') }}
        </el-button>
        <el-button @click="dialogVisible=false">
          {{ $t('common.cancel') }}
        </el-button>
      </span>
    </el-dialog>

    <el-dialog :title="$t('AIDE.editYarn')"
               :visible.sync="yarnDialog"
               custom-class="dialog-style"
               :close-on-click-modal="false"
               width="1100px">
      <el-form ref="yarnForm"
               :rules="ruleValidate"
               :model="yarnForm"
               label-width="190px">
        <div class="subtitle">
          {{ $t('DI.imageSettings') }}
        </div>
        <el-form-item :label="$t('DI.imageType')"
                      prop="image.imageType">
          <el-radio-group v-model="yarnForm.image.imageType">
            <el-radio label="Standard">{{ $t('DI.standard') }}</el-radio>
            <el-radio label="Custom">{{ $t('DI.custom') }}</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-row>
          <el-form-item v-if="yarnForm.image.imageType==='Standard'"
                        :label="$t('DI.imageSelection')"
                        prop="image.imageOption">
            <span>&nbsp;{{ defineImage }}&nbsp;&nbsp;</span>
            <el-select v-model="yarnForm.image.imageOption"
                       filterable
                       :placeholder="$t('DI.imagePro')">
              <el-option v-for="(item,index) in imageOptionList"
                         :label="item"
                         :key="index"
                         :value="item">
              </el-option>
            </el-select>
          </el-form-item>
          <el-form-item v-else
                        key="imageInput"
                        :label="$t('DI.imageSelection')"
                        prop="image.imageInput">
            <span>&nbsp;{{ defineImage }}&nbsp;&nbsp;</span>
            <el-input v-model="yarnForm.image.imageInput"
                      :placeholder="$t('DI.imageInputPro')" />
          </el-form-item>
        </el-row>
        <div class="subtitle">
          {{ $t('DI.computingResource') }}
        </div>
        <el-row>
          <el-col :span="12">
            <el-form-item :label="$t('AIDE.memory')"
                          prop="memory">
              <el-input v-model="yarnForm.memory"
                        :placeholder="$t('AIDE.memoryPro')">
                <span slot="append">Gi</span>
              </el-input>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item :label="$t('AIDE.cpu')"
                          prop="cpu">
              <el-input v-model="yarnForm.cpu"
                        :placeholder="$t('AIDE.CPUPro')">
                <span slot="append">Core</span>
              </el-input>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row>
          <el-col :span="12">
            <el-form-item :label="$t('AIDE.gpu')"
                          prop="extraResources">
              <el-input v-model="yarnForm.extraResources"
                        :placeholder="$t('AIDE.gpuPro')">
                <span slot="append">{{$t('AIDE.block')}}</span>
              </el-input>
            </el-form-item>
          </el-col>
        </el-row>
        <div class="subtitle">
          {{$t('AIDE.sparkResourceSettings')}}
        </div>
        <el-row>
          <el-col :span="12">
            <el-form-item :label="$t('AIDE.instanceNumber')"
                          prop="sparkSessionNum">
              <el-input v-model="yarnForm.sparkSessionNum"
                        maxlength="225"
                        :placeholder="$t('AIDE.instanceNumberPro')" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item :label="$t('AIDE.queueSettings')"
                          prop="queue">
              <el-input v-model="yarnForm.queue"
                        maxlength="225"
                        :placeholder="$t('DI.queuePro')" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row>
          <el-col :span="12">
            <el-form-item :label="$t('DI.driverMemory')"
                          prop="driverMemory">
              <el-input v-model="yarnForm.driverMemory"
                        maxlength="10"
                        :placeholder="$t('DI.driverMemoryPro')">
                <span slot="append">Gi</span>
              </el-input>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item :label="$t('DI.executorInstances')"
                          prop="executorCores">
              <el-input v-model="yarnForm.executorCores"
                        maxlength="10"
                        :placeholder="$t('DI.executorInstancesPro')">
                <span slot="append">Core</span>
              </el-input>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item :label="$t('DI.executorMemory')"
                          prop="executorMemory">
              <el-input v-model="yarnForm.executorMemory"
                        maxlength="10"
                        :placeholder="$t('DI.executorMemoryPro')">
                <span slot="append">Gi</span>
              </el-input>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item :label="$t('DI.linkisInstance')"
                          prop="executors">
              <el-input v-model="yarnForm.executors"
                        maxlength="10"
                        :placeholder="$t('DI.linkisInstancePro')" />
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
      <span slot="footer"
            class="dialog-footer">
        <el-button type="primary"
                   @click="saveYarnForm">
          {{ $t('common.save') }}
        </el-button>
        <el-button @click="cannelSaveYarnForm">
          {{ $t('common.cancel') }}
        </el-button>
      </span>
    </el-dialog>
    <el-dialog :title="$t('AIDE.notebookStatus')"
               :visible.sync="waitingStatusDialog"
               custom-class="dialog-style"
               :close-on-click-modal="false"
               width="920px">
      <div class="status-box">
        <el-row>
          <el-col class="label"
                  :span="5">{{$t('AIDE.resourceStaus')}}</el-col>
          <el-col class="content"
                  :span="19">{{waitingStatus.notebook_status}}</el-col>
        </el-row>
        <el-row class="status-height">
          <el-col class="label "
                  :span="5">{{$t('AIDE.statusMessage')}}</el-col>
          <el-col class="content"
                  :span="19">{{waitingStatus.notebook_status_info}}</el-col>
        </el-row>
        <el-row class="container-height">
          <el-col class="label"
                  :span="5">{{$t('AIDE.containerSatus')}}</el-col>
          <el-col class="content "
                  :span="19">
            <div class="detail-arr">[
              <div v-for="(item,index) in waitingStatus.containers_status_info"
                   :key="index"
                   class="detail-obj">{
                <div v-for="(val,key) in item"
                     :key="key"
                     class="detail-info">
                  "{{key}}"："{{val}}",
                </div>
                },
              </div>
              ]
            </div>
            <!-- <pre><code >{{waitingStatus.containers_status_info}}</code></pre> -->
          </el-col>
        </el-row>
      </div>
    </el-dialog>
    <view-log ref="viewLog"
              :url="logUrl"></view-log>
  </div>
</template>
<script type="text/ecmascript-6">
import util from '../../util/common'
import Step from '../../components/Step.vue'
import ViewLog from '../../components/ViewLog.vue'
import modelMixin from '../../util/modelMixin'
import { Tag } from 'element-ui'
export default {
  mixins: [modelMixin],
  components: {
    Step,
    ElTag: Tag,
    ViewLog
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
        sparkSessionNum: 0,
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
      originForm: {},
      yarnForm: {
        image: {
          'imageType': 'Standard',
          'imageOption': '',
          'imageInput': ''
        },
        cpu: '',
        memory: '',
        extraResources: '',
        sparkSessionNum: 0,
        queue: '',
        driverMemory: '',
        executorCores: '',
        executorMemory: '',
        executors: ''
      },
      userId: localStorage.getItem('userId'),
      storageDefalut: 'true',
      statusObj: {
        Ready: 'success',
        Waiting: '',
        Creating: '',
        Stopping: 'warning',
        Stopped: 'warning',
        Terminate: 'danger'
      },
      setIntervalList: null,
      waitingStatusDialog: false,
      waitingStatus: {},
      logUrl: ''
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
        sparkSessionNum: [
          { pattern: new RegExp(/^(([0]|([1-9]\d*))([.]0)*)$/), message: this.$t('AIDE.instanceNumberFormat') }
        ],
        queue: [
          // { required: true, message: this.$t('DI.queueReq') },
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
      let url = `/cc/${this.FesEnv.ccApiVersion}/auth/access/namespaces/${trData.namespace}/notebooks/${trData.name}`
      let xhr = new XMLHttpRequest()
      // const mlssToken = util.getCookieVal(localStorage.getItem('cookieKey')) || 0
      xhr.open('GET', url, false)
      xhr.setRequestHeader('Mlss-Userid', this.userId)
      // xhr.setRequestHeader('Mlss-Token', mlssToken)
      xhr.onreadystatechange = function (e) {
        if (this.readyState === 4 && this.status === 200) {
          let rst = JSON.parse(this.responseText).result
          rst.notebookAddress && window.open(window.encodeURI(rst.notebookAddress))
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
    getLog (trData) {
      this.logUrl = `/aide/v1/namespaces/${trData.namespace}/notebooks/${trData.name}/log`
      this.$refs.viewLog.logDialog = true
    },
    deleteJob (trData) {
      const url = `/aide/${this.FesEnv.aideApiVersion}/notebooks/${trData.id}`
      this.deleteListItem(url)
    },
    enableJob (trData) {
      this.FesApi.fetch(`/aide/${this.FesEnv.aideApiVersion}/notebooks/${trData.id}/start`, 'post').then(() => {
        this.toast()
        this.getListData()
      })
    },
    disableJob (trData) {
      this.FesApi.fetch(`/aide/${this.FesEnv.aideApiVersion}/notebooks/${trData.id}/stop`, 'delete').then(() => {
        this.toast()
        this.getListData()
      })
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
      job.sparkSessionNum = trData.sparkSessionNum ? trData.sparkSessionNum : 0
      job.executorMemory = isNaN(parseInt(job.executorMemory)) ? '' : parseInt(job.executorMemory) + ''
      job.driverMemory = isNaN(parseInt(job.driverMemory)) ? '' : parseInt(job.driverMemory) + ''
      imageType === 'Standard' ? job.image.imageOption = trData.image : job.image.imageInput = trData.image
      setTimeout(() => {
        Object.assign(this.form, job)
        if (trData.dataVolume && trData.dataVolume[0].localPath) {
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
          if (this.haveProxy) {
            this.form.workspaceVolume.localPath = trData.dataVolume[0].localPath
          }
        } else {
          this.form.workspaceVolume = trData.workspaceVolume
          if (this.haveProxy) {
            this.form.dataVolume[0].localPath = trData.workspaceVolume.localPath
          }
          // 存储设置类型 为 默认
          this.storageDefalut = 'true'
        }
      }, 100)
      this.dialogVisible = true
    },
    nodeBookStatus (trData) {
      const status = trData.status.toLowerCase()
      if (status !== 'stop') {
        this.waitingStatus = {}
        this.FesApi.fetch(`/aide/v1/notebooks/${trData.id}/status`, 'get').then((res) => {
          this.waitingStatus = res
          this.waitingStatusDialog = true
        })
      }
    },
    nextStep () {
      this.$refs.formValidate.validate((valid) => {
        if (valid && this.currentStep < this.stepList.length - 1) {
          if (this.currentStep === 2 && this.form.cluster && (parseInt(this.form.sparkSessionNum) > this.FesEnv.AIDE.SparkSessionCount)) {
            this.$message.warning(this.$t('AIDE.instancesNumberFormatN', { n: this.FesEnv.AIDE.SparkSessionCount }))
            return
          }
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
              if (this.form[key] !== '') {
                param[key] = this.form[key]
              }
            }
          }
          if (this.form.cluster) {
            param.sparkSessionNum = this.form.sparkSessionNum ? parseInt(this.form.sparkSessionNum) : 0
          }

          param.cpu = parseFloat(this.form.cpu)
          param.memory = {
            memoryAmount: parseInt(this.form.memory),
            memoryUnit: 'Gi'
          }
          param.extraResources = JSON.stringify({ 'nvidia.com/gpu': this.form.extraResources })
          let imageColumn = this.form.image.imageType === 'Standard' ? this.form.image.imageOption : this.form.image.imageInput
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
      this.haveProxy = false
      this.$refs.formValidate.resetFields()
      const formKey = ['name', 'cpu', 'extraResources', 'memory', 'namespace', 'queue', 'driverMemory', 'executorCores', 'executorMemory', 'executors']
      for (let item of formKey) {
        this.form[item] = ''
      }
      for (let key in this.form.image) {
        this.form.image[key] = key === 'imageType' ? 'Standard' : ''
      }
      this.form.sparkSessionNum = 0
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
      debugger
      this.yarnDialog = true
      const imageType = this.imageOptionList.indexOf(trData.image) > -1 ? 'Standard' : 'Custom'
      let image = {
        'imageType': imageType,
        imageOption: '',
        imageInput: ''
      }
      imageType === 'Standard' ? image.imageOption = trData.image : image.imageInput = trData.image
      this.yarnForm = {
        id: trData.id,
        image: image,
        memory: parseInt(trData.memory),
        cpu: trData.cpu,
        extraResources: trData.gpu + '',
        sparkSessionNum: isNaN(parseInt(trData.sparkSessionNum)) ? 0 : trData.sparkSessionNum + '',
        queue: trData.queue || '',
        driverMemory: isNaN(parseInt(trData.driverMemory)) ? '' : parseInt(trData.driverMemory) + '',
        executorCores: isNaN(parseInt(trData.executorCores)) ? '' : parseInt(trData.executorCores) + '',
        executorMemory: isNaN(parseInt(trData.executorMemory)) ? '' : parseInt(trData.executorMemory) + '',
        executors: isNaN(parseInt(trData.executors)) ? '' : parseInt(trData.executors) + ''
      }
      this.originForm = {
        memoryAmount: parseInt(trData.memory),
        cpu: parseFloat(trData.cpu),
        extraResources: trData.gpu + ''
      }
      this.originForm.image = { ...image }
    },
    saveYarnForm () {
      this.$refs.yarnForm.validate((valid) => {
        if (valid) {
          let parms = {}
          parms.queue = this.yarnForm.queue
          parms.executorCores = this.yarnForm.executorCores
          parms.executors = this.yarnForm.executors
          parms.cpu = parseFloat(this.yarnForm.cpu)
          parms.memoryAmount = parseInt(this.yarnForm.memory)
          parms.memoryUnit = 'Gi'
          parms.extraResources = JSON.stringify({ 'nvidia.com/gpu': this.yarnForm.extraResources })
          let imageColumn = this.yarnForm.image.imageType === 'Standard' ? this.yarnForm.image.imageOption : this.yarnForm.image.imageInput
          parms.imageName = this.defineImage + ':' + imageColumn
          parms.imageType = this.yarnForm.image.imageType
          parms.sparkSessionNum = this.yarnForm.sparkSessionNum === '' ? 0 : parseInt(this.yarnForm.sparkSessionNum)
          parms.driverMemory = Number(this.yarnForm.driverMemory) > 0 ? this.yarnForm.driverMemory + 'g' : ''
          parms.executorMemory = Number(this.yarnForm.executorMemory) > 0 ? this.yarnForm.executorMemory + 'g' : ''
          if (parseInt(parms.sparkSessionNum) > this.FesEnv.AIDE.SparkSessionCount) {
            this.$message.warning(this.$t('AIDE.instancesNumberFormatN', { n: this.FesEnv.AIDE.SparkSessionCount }))
            return
          }
          // 校验资源是否修改
          let resourceChange = false
          if (this.originForm.image.imageType === 'Standard') {
            resourceChange = this.originForm.image.imageOption !== this.yarnForm.image.imageOption
          } else {
            resourceChange = this.originForm.image.imageInput !== this.yarnForm.image.imageInput
          }
          const resourceKey = ['cpu', 'memoryAmount']
          for (let key of resourceKey) {
            if (this.originForm[key] !== parms[key]) {
              resourceChange = true
            }
          }
          if (this.originForm.extraResources !== this.yarnForm.extraResources) {
            resourceChange = true
          }
          if (resourceChange) {
            this.$alert(this.$t('AIDE.resourceModifyFormat'), this.$t('common.prompt'), {
              confirmButtonText: this.$t('common.comfirm'),
              callback: action => {
                this.submitYarnForm(parms)
              }
            })
          } else {
            this.submitYarnForm(parms)
          }
        }
      })
    },
    submitYarnForm (parms) {
      this.FesApi.fetch(`/aide/${this.FesEnv.aideApiVersion}/notebooks/${this.yarnForm.id}`, parms, {
        method: 'patch',
        headers: {
          'Mlss-Userid': this.userId
        }
      }).then(() => {
        this.getListData()
        this.toast()
        this.yarnDialog = false
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
.pointer {
  cursor: pointer;
}
.status-box {
  .el-row {
    height: 32px;
    line-height: 20px;
    margin-bottom: 20px;
  }
  .content {
    height: 100%;
    padding: 5px;
    color: #979da9;
    border: 1px solid #dcdfe6;
    overflow-y: scroll;
  }
  .status-height {
    min-height: 60px;
    overflow: hidden;
  }
  .container-height {
    min-height: 350px;
    overflow: hidden;
  }
  .label {
    text-align: right;
    padding-right: 20px;
  }
}
.detail-arr {
  padding: 5px;
}
.detail-obj {
  padding: 0 5px;
}
.detail-info {
  padding: 2px 0;
}
</style>
