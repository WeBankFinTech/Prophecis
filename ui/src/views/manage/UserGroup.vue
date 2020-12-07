<template>
  <div>
    <breadcrumb-nav class="margin-b20">
      <el-button type="primary"
                 icon="el-icon-plus"
                 @click="dialogVisible = true"
                 size="small">
        {{ $t("group.addGroup") }}
      </el-button>
    </breadcrumb-nav>
    <el-table :data="dataList"
              border
              style="width: 100%"
              :empty-text="$t('common.noData')">
      <el-table-column prop="name"
                       :label="$t('group.groupName')" />
      <el-table-column prop="departmentName"
                       :label="$t('group.departmentName')" />
      <el-table-column prop="groupTypeT"
                       :label="$t('group.groupType')"
                       :formatter="translationGroupType" />
      <el-table-column prop="isActive"
                       :label="$t('user.isActive')"
                       :formatter="translationActive" />
      <el-table-column prop="remarks"
                       :label="$t('DI.description')"
                       min-width="130%" />
      <el-table-column :label="$t('common.operation')">
        <template slot-scope="scope">
          <el-dropdown size="medium"
                       @command="handleCommand">
            <el-button type="text"
                       class="multi-operate-btn"
                       icon="el-icon-more"
                       round></el-button>
            <el-dropdown-menu slot="dropdown">
              <el-dropdown-item :command="beforeHandleCommand('updateItem',scope.row)">
                {{$t('common.modify')}}
              </el-dropdown-item>
              <el-dropdown-item :command="beforeHandleCommand('deleteItem',scope.row)">
                {{$t('common.delete')}}
              </el-dropdown-item>
              <el-dropdown-item :command="beforeHandleCommand('viewNameSpace',scope.row)">
                {{$t('group.viewNameSpace')}}
              </el-dropdown-item>
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
               width="900px">
      <el-form ref="formValidate"
               :model="form"
               :rules="ruleValidate"
               label-width="180px">
        <el-row :gutter="24">
          <el-col :span="12">
            <el-form-item :label="$t('group.groupName')"
                          prop="name">
              <el-input v-model="form.name"
                        :placeholder="$t('group.groupNamePro')"
                        :disabled="modify" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item :label="$t('group.groupType')"
                          prop="groupType">
              <el-select v-model="form.groupType"
                         :disabled="modify"
                         :placeholder="$t('group.groupTypePro')">
                <el-option :label="$t('group.system')"
                           value="SYSTEM">
                </el-option>
                <el-option :label="$t('group.private')"
                           value="PRIVATE">
                </el-option>
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="24">
          <el-col :span="12">
            <el-form-item :label="$t('group.departmentName')"
                          prop="departmentName">
              <el-input v-model="form.departmentName"
                        :placeholder="$t('group.departmentNamePro')"
                        :maxlength="255" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item :label="$t('group.departmentId')"
                          prop="departmentId">
              <el-input v-model="form.departmentId"
                        :placeholder="$t('group.departmentIdPro')"
                        :maxlength="5" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item :label="$t('user.isActive')"
                      prop="isActive">
          <el-select v-model="form.isActive"
                     :placeholder="$t('user.activePro')">
            <el-option :label="$t('common.yes')"
                       :value="1">
            </el-option>
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('group.groupDescribe')"
                      prop="remarks">
          <el-input v-model="form.remarks"
                    type="textarea"
                    :placeholder="$t('group.groupDescribePro')" />
        </el-form-item>
      </el-form>
      <span slot="footer"
            class="dialog-footer">
        <el-button type="primary"
                   :disabled="btnDisabled"
                   @click="subInfo ">{{ $t("common.save") }}</el-button>
        <el-button @click="dialogVisible=false">{{ $t("common.cancel") }}</el-button>
      </span>
    </el-dialog>
    <el-dialog :title="$t('group.nameSpaceList')"
               :visible.sync="spaceDialogVisible"
               custom-class="dialog-style">
      <el-table :data="spaceList"
                border
                :empty-text="$t('common.noData')">
        <el-table-column prop="namespace"
                         :label="$t('group.spaceName')" />
      </el-table>
    </el-dialog>
  </div>
</template>
<script type="text/ecmascript-6">
export default {
  data () {
    return {
      dialogVisible: false,
      spaceDialogVisible: false,
      form: {
        name: '',
        departmentId: '',
        departmentName: '',
        groupType: 'SYSTEM',
        isActive: 1,
        enableFlag: 1,
        remarks: ''
      },
      modify: false,
      dataList: [],
      spaceList: []
    }
  },
  mounted () {
    this.getListData()
  },
  computed: {
    title () {
      return this.modify ? this.$t('group.modifyGroup') : this.$t('group.newGroup')
    },
    getUrl () {
      return `/cc/${this.FesEnv.ccApiVersion}/groups`
    },
    ruleValidate () {
      // 切换语言表单报错重置表单
      this.resetFormFields()
      return {
        name: [
          { required: true, message: this.$t('group.groupNameReq') },
          { pattern: new RegExp(/^.*[^\s][\s\S]*$/), message: this.$t('group.groupNameReq') }
        ],
        groupType: [
          { required: true, message: this.$t('group.groupTypeReq') }
        ],
        isActive: [
          { required: true, message: this.$t('user.activeReq') }
        ],
        departmentId: [
          { pattern: new RegExp(/^\d*$/), message: this.$t('group.departmentIdFormat') }
        ],
        departmentName: [
          { required: true, message: this.$t('group.departmentNameReq') },
          { pattern: new RegExp(/^.*[^\s][\s\S]*$/), message: this.$t('group.departmentNameReq') }
        ],
        remarks: [
          { pattern: new RegExp(/^.{0,200}$/), message: this.$t('group.remarksMaxLength') }
        ]
      }
    }
  },
  methods: {
    translationActive (value, row) {
      return this.$t('common.yes')
    },
    translationGroupType (value, row) {
      return value.groupType === 'PRIVATE' ? this.$t('group.private') : this.$t('group.system')
    },
    updateItem (trData) {
      this.dialogVisible = true
      this.modify = true
      setTimeout(() => {
        Object.assign(this.form, trData)
      }, 0)
    },
    deleteItem (trData) {
      const url = `/cc/${this.FesEnv.ccApiVersion}/groups/id/${trData.id}`
      this.deleteListItem(url)
    },
    viewNameSpace (trData) {
      this.getSpaceList(trData.id)
      this.spaceDialogVisible = true
    },
    subInfo () {
      this.form.guidCheck = !!parseInt(this.form.guidCk)
      this.$refs.formValidate.validate((valid) => {
        if (valid) {
          this.setBtnDisabeld()
          let method = this.modify ? 'put' : 'post'
          let parms = { ...this.form }
          delete parms.guidCheck
          delete parms.isActive
          parms.departmentId = parseInt(parms.departmentId)
          this.FesApi.fetch(`/cc/${this.FesEnv.ccApiVersion}/groups`, parms, method).then(() => {
            this.getListData()
            this.toast()
            this.dialogVisible = false
          })
        }
      })
    },
    getSpaceList (id) {
      let url = `/cc/${this.FesEnv.ccApiVersion}/groups/${id}/namespaces`
      this.FesApi.fetch(url, {}, 'get').then(rst => {
        this.spaceList = rst
      })
    }
  }
}
</script>
