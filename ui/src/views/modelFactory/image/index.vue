<template>
  <div>
    <breadcrumb-nav></breadcrumb-nav>
    <el-form ref="queryValidate"
             :model="query"
             :rules="ruleValidate"
             class="query"
             label-width="0">
      <el-row>
        <el-col :span="8">
          <el-form-item>
            <el-input v-model="query.query_str"
                      :placeholder="$t('expExeRecord.enterSearchContent')"></el-input>
          </el-form-item>
        </el-col>
        <el-col :span="16">
          <el-button type="primary"
                     class="margin-l30"
                     icon="el-icon-search"
                     @click="filterListData">
            {{ $t("common.filter") }}
          </el-button>
          <el-button type="warning"
                     icon="el-icon-refresh-left"
                     @click="resetQueryStrListData">
            {{ $t("common.reset") }}
          </el-button>
          <el-button type="primary"
                     icon="plus-circle-o"
                     @click="dialogVisible=true">
            {{$t('image.createImage')}}
          </el-button>
        </el-col>
      </el-row>
    </el-form>
    <el-table :data="dataList"
              border
              style="width: 100%"
              :empty-text="$t('common.noData')">
      <el-table-column prop="image"
                       :show-overflow-tooltip="true"
                       :label="$t('image.imageBase')"
                       min-width="150" />
      <el-table-column prop="image_tag"
                       :show-overflow-tooltip="true"
                       :label="$t('image.imageName')"
                       min-width="100" />
      <el-table-column prop="user_name"
                       :show-overflow-tooltip="true"
                       :label="$t('user.user')"
                       min-width="90" />
      <el-table-column prop="model_name"
                       :show-overflow-tooltip="true"
                       :label="$t('modelList.modelName')" />
      <el-table-column prop="version"
                       width="80"
                       :show-overflow-tooltip="true"
                       :label="$t('serviceList.modelVersion')" />
      <el-table-column prop="creation_timestamp"
                       :show-overflow-tooltip="true"
                       :label="$t('DI.createTime')"
                       min-width="90" />
      <el-table-column prop="last_updated_timestamp"
                       :show-overflow-tooltip="true"
                       :label="$t('serviceList.updateTime')"
                       min-width="90" />
      <el-table-column prop="status"
                       :show-overflow-tooltip="true"
                       min-width="90"
                       :label="$t('DI.status')">
        <template slot-scope="scope">
          <!-- <el-button v-if="scope.row.status==='FAIL'"
                     type="primary"
                     size="mini"
                     round
                     plain
                     @click="promptFail(scope.row)">
            {{scope.row.status}}
          </el-button>
          <span v-else>
            {{scope.row.status}}
          </span> -->
          <status-tag :status="scope.row.status"
                      :statusObj="statusObj">
          </status-tag>
        </template>
      </el-table-column>
      <el-table-column :label="$t('common.operation')"
                       width="85">
        <template slot-scope="scope">
          <el-dropdown size="medium"
                       @command="handleCommand">
            <el-button type="text"
                       class="multi-operate-btn"
                       icon="el-icon-more"
                       round></el-button>
            <el-dropdown-menu slot="dropdown">
              <el-dropdown-item v-if="scope.row.status==='COMPLETE'||scope.row.status==='FAIL'"
                                :command="beforeHandleCommand('modifyFun',scope.row)">{{$t('common.modify')}}</el-dropdown-item>
              <el-dropdown-item :command="beforeHandleCommand('deleteItem',scope.row)">{{$t('common.delete')}}</el-dropdown-item>
              <el-dropdown-item v-if="scope.row.status==='FAIL'"
                                :command="beforeHandleCommand('promptFail',scope.row)">{{$t('image.failureReason')}}</el-dropdown-item>
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
    <el-dialog :title="title"
               :visible.sync="dialogVisible"
               @close="clearDialogForm"
               custom-class="dialog-style"
               :close-on-click-modal="false"
               width="900px">
      <step :step="currentStep"
            :step-list="stepList" />
      <el-form ref="formValidate"
               :model="form"
               class="add-node-form"
               :rules="ruleValidate"
               label-width="130px">
        <div v-if="currentStep===0"
             key="step0">
          <div class="subtitle">
            {{ $t('DI.basicSettings') }}
          </div>
          <el-form-item :label="$t('image.imageBase')"
                        prop="image"
                        maxlength="225">
            <el-input v-model="form.image"
                      :placeholder="$t('image.imageBasePro')" />
          </el-form-item>
          <el-form-item :label="$t('image.imageTag')"
                        prop="image_tag">
            <el-input v-model="form.image_tag"
                      :placeholder="$t('image.imageTagPro')" />
          </el-form-item>
          <el-form-item :label="$t('image.imageRemark')">
            <el-input v-model="form.remarks"
                      type="textarea"
                      :placeholder="$t('image.imageRemarkPro')" />
          </el-form-item>
        </div>
        <div v-if="currentStep===1"
             key="step1">
          <div class="subtitle">
            {{ $t('serviceList.modelSettings') }}
          </div>
          <el-row>
            <el-col :span="12">
              <el-form-item :label="$t('user.userGroup')"
                            prop="group_id">
                <el-select v-model="form.group_id"
                           :filterable="true"
                           :placeholder="$t('user.userGroupSettingsPro')"
                           @change="changeGroup">
                  <el-option v-for="(item,index) of groupList"
                             :label="item.name"
                             :value="item.id"
                             :key="index">
                  </el-option>
                </el-select>
              </el-form-item>
            </el-col>
            <el-col :span="12"
                    v-if="form.group_id">
              <el-form-item :label="$t('modelList.modelName')"
                            prop="model_name">
                <el-select v-model="form.model_name"
                           :placeholder="$t('serviceList.modelNamePro1')"
                           @change="selectChangeModel">
                  <el-option v-for="(item,index) of modelOptionList"
                             :label="item.model_name"
                             :value="item.id"
                             :key="index">
                  </el-option>
                </el-select>
              </el-form-item>
            </el-col>
          </el-row>
          <el-row v-if="form.model_name">
            <el-col :span="12">
              <el-form-item :label="$t('serviceList.modelVersion')"
                            prop="model_version_id">
                <el-select v-model="form.model_version_id"
                           :placeholder="$t('serviceList.modelVersionPro')">
                  <el-option v-for="(item,index) of versionOptionList"
                             :label="item.name"
                             :value="item.id"
                             :key="index">
                  </el-option>
                </el-select>
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item :label="$t('serviceList.modelType')"
                            prop="model_type">
                <el-input v-model="form.model_type"
                          disabled></el-input>
              </el-form-item>
            </el-col>
          </el-row>
        </div>
      </el-form>
      <span slot="footer"
            class="dialog-footer">
        <el-button v-if="currentStep>0"
                   type="primary"
                   @click="lastStep">
          {{ $t('common.lastStep') }}
        </el-button>
        <el-button v-if="currentStep<1"
                   type="primary"
                   @click="nextStep">
          {{ $t('common.nextStep') }}
        </el-button>
        <el-button v-if="currentStep===1"
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
  </div>
</template>
<script type="text/ecmascript-6">
import util from '../../../util/common'
import Step from '../../../components/Step.vue'
import modelMixin from '../../../util/modelMixin'
import StatusTag from '../../../components/StatusTag'
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
        query_str: ''
      },
      dataList: [],
      groupList: [],
      modify: false,
      modelOptionList: [],
      versionOptionList: [],
      currentStep: 0,
      form: {
        image: '',
        image_tag: '',
        remarks: '',
        group_id: '',
        model_name: '',
        model_version_id: '',
        model_type: ''
      },
      loading: false,
      statusObj: {
        Creating: '',
        BUILDING: '',
        PUSHING: '',
        COMPLETE: 'success',
        FAIL: 'danger'
      }
    }
  },
  mounted () {
    this.getListData()
    this.getGroupList()
  },
  computed: {
    noDataText () {
      return !this.loading ? this.$t('common.noData') : ''
    },
    stepList () {
      return [this.$t('image.imageBaseInfo'), this.$t('serviceList.modelSettings')]
    },
    title () {
      return this.modify ? this.$t('image.modifyImage') : this.$t('image.createImage')
    },
    ruleValidate () {
      // 切换语言表单报错重置表单
      this.resetQueryFields()
      this.resetFormFields()
      return {
        image: [
          { required: true, message: this.$t('image.imageBaseReq') },
          { type: 'string', pattern: new RegExp(/^[a-z][a-z0-9-_./]*$/), message: this.$t('image.imageBaseFormat') }
        ],
        image_tag: [
          { required: true, message: this.$t('image.imageTagReq') },
          { type: 'string', pattern: new RegExp(/^[a-z0-9A-Z][a-z0-9A-Z-_./]*$/), message: this.$t('image.imageTagFormat') }
        ],
        group_id: [
          { required: true, message: this.$t('user.userGroupSettingsReq') }
        ],
        model_name: [
          { required: true, message: this.$t('modelList.modelNameReq') }
        ],
        model_version_id: [
          { required: true, message: this.$t('serviceList.modelVersionReq') }
        ]
      }
    }
  },
  methods: {
    deleteItem (trData) {
      const url = `/mf/v1/image/${trData.id}`
      this.deleteListItem(url)
    },
    modifyFun (trData) {
      this.modify = true
      this.getModelList(trData.group_id)
      let param = {}
      for (let key in this.form) {
        param[key] = trData[key]
      }
      param.id = trData.id
      let modelId = trData.model_version.model.id
      this.getVersionList(modelId)
      param.model_type = trData.model_version.model.model_type
      param.model_name = modelId
      param.model_version_id += ''
      this.dialogVisible = true
      setTimeout(() => {
        this.form = param
      }, 0)
    },
    nextStep () {
      this.$refs.formValidate.validate((valid) => {
        if (valid && this.currentStep < this.stepList.length - 1) {
          this.currentStep++
        }
      })
    },
    promptFail (trData) {
      this.$alert(trData.msg ? trData.msg : this.$t('common.noData'), this.$t('image.failureReason'))
    },
    getListData () {
      this.loading = true
      let url = '/mf/v1/images'
      let params = {
        page: this.pagination.pageNumber,
        size: this.pagination.pageSize,
        query_str: this.query.query_str
      }
      this.FesApi.fetch(url, params, 'get').then(rst => {
        rst.images && rst.images.forEach(item => {
          let index = item.image_name.indexOf(':')
          item.image = item.image_name.substring(0, index)
          item.image_tag = item.image_name.substring(index + 1)
          item.creation_timestamp = item.creation_timestamp ? util.transDate(item.creation_timestamp) : ''
          item.last_updated_timestamp = item.last_updated_timestamp ? util.transDate(item.last_updated_timestamp) : ''
          item.user_name = item.user.name
          item.model_name = item.model_version.model.model_name
          item.version = item.model_version.version
        })
        this.dataList = Object.freeze(rst.images)
        this.pagination.totalPage = rst.total || 0
        this.loading = false
      }, () => {
        this.dataList = []
        this.loading = false
      })
    },
    subInfo () {
      this.$refs.formValidate.validate((valid) => {
        if (valid) {
          let param = {
            model_version_id: parseInt(this.form.model_version_id),
            remarks: this.form.remarks
          }
          let method = 'post'
          let url = '/mf/v1/image'
          if (this.modify) {
            method = 'put'
            url = `/mf/v1/image/${this.form.id}`
          }
          param.group_id = this.form.group_id
          param.image_name = this.form.image + ':' + this.form.image_tag
          param.user_name = localStorage.getItem('userId')
          this.setBtnDisabeld()
          this.FesApi.fetch(url, param, method).then(() => {
            this.getListData()
            this.toast()
            this.dialogVisible = false
          })
        }
      })
    },
    changeGroup (val) {
      this.form.model_name = ''
      this.form.model_version_id = ''
      this.versionOptionList = []
      this.getModelList(val)
    },
    getModelList (val) {
      this.FesApi.fetch(`/mf/v1/modelsByGroupID/${val}`, { page: 1, size: 1000 }, 'get').then((res) => {
        let modelOption = []
        for (let item of res.models) {
          if (item.model_type === 'CUSTOM') {
            modelOption.push(item)
          }
        }
        this.modelOptionList = this.formatResList(modelOption)
      })
    },
    selectChangeModel (val) {
      this.form.model_version_id = ''
      for (let item of this.modelOptionList) {
        if (item.id === val) {
          this.form.model_type = item.model_type
          this.serviceModels = {
            model_name: item.model_name,
            model_type: item.model_type,
            group_name: item.group.name
          }
          break
        }
      }
      this.getVersionList(val)
    },
    getVersionList (val) {
      this.FesApi.fetch(`/mf/v1/modelversions/${val}`, { page: 1, size: 1000 }, 'get').then((res) => {
        let optionList = []
        if (res.modelVersions) {
          for (let item of res.modelVersions) {
            optionList.push({
              id: item.id + '',
              name: item.version,
              source: item.source
            })
          }
        }
        this.versionOptionList = Object.freeze(optionList)
      })
    },
    clearDialogForm () {
      this.currentStep = 0
      this.modelOptionList = []
      this.versionOptionList = []
      this.modify = false
      for (let key in this.form) {
        this.form[key] = ''
      }
      setTimeout(() => {
        this.$refs.formValidate.resetFields()
      }, 0)
    }
  }
}
</script>
<style lang="scss" scoped>
.add-node-form {
  min-height: 250px;
}
</style>
