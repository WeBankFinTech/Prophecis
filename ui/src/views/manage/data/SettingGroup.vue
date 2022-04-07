<template>
  <div>
    <breadcrumb-nav>
      <el-button type="primary"
                 icon="el-icon-plus"
                 size="small"
                 @click="dialogVisible=true">
        {{ $t('user.newGroup') }}
      </el-button>
    </breadcrumb-nav>
    <el-table :data="dataList"
              border
              style="width: 100%"
              :empty-text="$t('common.noData')">
      <el-table-column prop="path"
                       :label="$t('data.dataStorePath')" />
      <el-table-column prop="groupName"
                       :label="$t('user.userGroup')" />
      <el-table-column prop="permissionsT"
                       :label="$t('data.permission')"
                       :formatter="translationPermissions" />
      <el-table-column :label="$t('common.operation')">
        <template slot-scope="scope">
          <el-button @click="updateItem(scope.row)"
                     type="text"
                     size="small">{{$t('common.modify')}}</el-button>
          <el-button @click="deleteItem(scope.row)"
                     type="text"
                     size="small">{{$t('common.delete')}}</el-button>
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
               width="700px">
      <el-form ref="formValidate"
               :model="form"
               :rules="ruleValidate"
               label-width="170px">
        <el-form-item :label="$t('data.dataStorePath')">
          <el-input v-model="path"
                    :placeholder="$t('data.dataStorePro')"
                    disabled />
        </el-form-item>
        <el-form-item :label="$t('user.userGroupSettings')"
                      prop="groupId">
          <el-select v-model="form.groupId"
                     filterable
                     :placeholder="$t('user.userGroupSettingsPro')"
                     :disabled="modify">
            <el-option v-for="item in groupOptionList"
                       :label="item.name"
                       :key="item.id"
                       :value="item.id+'~'+item.groupType">
            </el-option>
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('data.permission')"
                      prop="permission">
          <el-select v-model="form.permissions"
                     :placeholder="$t('data.permissionPro')">
            <el-option :label="$t('AIDE.readOnly')"
                       value="r">
            </el-option>
            <el-option :label="$t('AIDE.readWrite')"
                       value="rw">
            </el-option>
          </el-select>
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
  </div>
</template>
<script type="text/ecmascript-6">
export default {
  props: {
    storageId: {
      type: String,
      default: ''
    }
  },
  data () {
    return {
      dialogVisible: false,
      path: '',
      form: {
        groupId: '',
        permissions: ''
      },
      modify: false,
      dataList: [],
      groupOptionList: []
    }
  },
  mounted () {
    this.path = this.$route.query.path
    this.getListData()
    this.getGroupOption()
  },
  computed: {
    title () {
      return this.modify ? this.$t('user.modifyGroup') : this.$t('user.newGroup')
    },
    getUrl () {
      return `/cc/${this.FesEnv.ccApiVersion}/groups/groupStorage/storage/${this.storageId}`
    },
    ruleValidate () {
      // 切换语言表单报错重置表单
      this.resetFormFields()
      return {
        roleId: [
          { required: true, message: this.$t('user.roleReq') }
        ],
        groupVal: [
          { required: true, message: this.$t('user.userGroupSettingsReq') }
        ]
      }
    }
  },
  methods: {
    translationPermissions (value, row) {
      return value.permissions === 'r' ? this.$t('AIDE.readOnly') : this.$t('AIDE.readWrite')
    },
    updateItem (trData) {
      this.dialogVisible = true
      this.modify = true
      setTimeout(() => {
        Object.assign(this.form, trData)
        this.form.groupId = trData.groupId + '~' + trData.groupType
      }, 0)
    },
    deleteItem (trData) {
      const url = `/cc/${this.FesEnv.ccApiVersion}/groups/groupStorage/id/${trData.id}`
      this.deleteListItem(url)
    },
    subInfo () {
      this.$refs.formValidate.validate((valid) => {
        if (valid) {
          this.setBtnDisabeld()
          let method = this.modify ? 'put' : 'post'
          let param = { ...this.form }
          const groupParam = param.groupId.split('~')
          param.groupId = +groupParam[0]
          param.groupType = groupParam[1]
          param.type = groupParam[1]
          param.path = this.path
          param.storageId = parseInt(this.storageId)
          this.FesApi.fetch(`/cc/${this.FesEnv.ccApiVersion}/groups/storages`, param, method).then(() => {
            this.getListData()
            this.toast()
            this.dialogVisible = false
          })
        }
      })
    },
    clearDialogForm () {
      this.form.groupId = ''
      this.form.permissions = ''
      this.modify = false
      this.$refs.formValidate.resetFields()
    }
  }
}
</script>
