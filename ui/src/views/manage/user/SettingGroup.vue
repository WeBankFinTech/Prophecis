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
              :empty-text="$t('common.noData')">
      <el-table-column prop="username"
                       :label="$t('user.userName')" />
      <el-table-column prop="groupName"
                       :label="$t('user.userGroup')" />
      <el-table-column prop="roleName"
                       :label="$t('user.role')" />
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
               width="600px">
      <el-form ref="formValidate"
               :model="form"
               :rules="ruleValidate"
               label-width="185px">
        <el-form-item :label="$t('user.userName')">
          <el-input v-model="userName"
                    :placeholder="$t('user.userNamePro')"
                    disabled />
        </el-form-item>
        <el-form-item :label="$t('user.userGroupSettings')"
                      prop="groupVal">
          <el-select v-model="form.groupVal"
                     filterable
                     :placeholder="$t('user.userGroupSettingsPro')"
                     :disabled="modifyGroup"
                     @change="changeGroupId">
            <el-option v-for="item in groupOptionList"
                       :label="item.name"
                       :key="item.id"
                       :value="item.id+'-'+item.groupType">
            </el-option>
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('user.roleSet')"
                      prop="roleId">
          <el-select v-model="form.roleId"
                     :placeholder="$t('user.roleSetPro')"
                     :disabled="roleDisable">
            <el-option v-for="item in roleOptionList"
                       :label="item.name"
                       :key="item.id"
                       :value="item.id">
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
  data () {
    return {
      userName: '',
      dialogVisible: false,
      form: {
        id: 0,
        groupId: '',
        groupVal: '',
        roleId: '',
        remarks: ''
      },
      roleDisable: false,
      modifyGroup: false,
      beforeModify: false, // 刚进入修改，用户组变化，角色不清空
      dataList: [],
      roleOptionList: [],
      groupOptionList: []
    }
  },
  mounted () {
    this.userName = this.$route.query.userName
    this.userId = Number(this.$route.params.userId)
    this.getListData()
    this.getGroupOption()
    this.getRoleOption()
  },
  computed: {
    title () {
      return this.modifyGroup ? this.$t('user.modifyGroup') : this.$t('user.newGroup')
    },
    getUrl () {
      return `/cc/${this.FesEnv.ccApiVersion}/groups/userGroup/user/${Number(this.$route.params.userId)}`
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
    // formatResList (res) {
    //   if (Array.isArray(res)) {
    //     return res
    //   }
    //   if (res.hasOwnProperty('models') && Array.isArray(res.models)) {
    //     return res.models
    //   }
    //   if (res.hasOwnProperty('list') && Array.isArray(res.list)) {
    //     return res.list
    //   }
    //   return []
    // },
    // getGroupOption () {
    //   this.FesApi.fetch(`/cc/${this.FesEnv.ccApiVersion}/groups`, 'get').then(rst => {
    //     this.groupOptionList = this.formatResList(rst)
    //   })
    // },
    getRoleOption () {
      this.FesApi.fetch(`/cc/${this.FesEnv.ccApiVersion}/roles`, 'get').then(rst => {
        this.roleOptionList = this.formatResList(rst)
      })
    },
    changeGroupId (val) {
      if (val) {
        let index = val.indexOf('-')
        let groupType = val.substring(index + 1)
        this.form.groupId = val.substring(0, index)
        if (groupType === 'PRIVATE') {
          this.form.roleId = 1
          this.roleDisable = true
        } else {
          if (this.beforeModify) {
            this.beforeModify = false
            return
          }
          this.form.roleId = ''
          this.roleDisable = false
        }
      }
    },
    updateItem (trData) {
      this.beforeModify = true
      this.modifyGroup = true
      this.roleDisable = false
      setTimeout(() => {
        this.form.id = trData.id
        this.form.roleId = trData.roleId
        this.form.groupId = trData.groupId
        this.form.groupVal = trData.groupId + '-' + trData.groupType

        if (trData.groupType === 'PRIVATE' && trData.roleId === 1) {
          this.roleDisable = true
        }
      }, 0)
      this.dialogVisible = true
    },
    deleteItem (trData) {
      const url = `/cc/${this.FesEnv.ccApiVersion}/groups/userGroup/id/${trData.id}`
      this.deleteListItem(url)
    },
    subInfo () {
      this.$refs.formValidate.validate((valid) => {
        if (valid) {
          this.setBtnDisabeld()
          let method = this.modifyGroup ? 'put' : 'post'
          this.form.groupId = this.form.groupVal.substring(0, this.form.groupVal.indexOf('-'))
          let param = { ...this.form }
          param.username = this.userName
          param.userId = this.userId
          delete param.groupVal
          !this.modifyGroup && delete param.username
          !this.modifyGroup && delete param.id
          param.groupId = parseInt(param.groupId)
          param.roleId = parseInt(param.roleId)
          this.FesApi.fetch(`/cc/${this.FesEnv.ccApiVersion}/groups/users`, param, method).then(() => {
            this.getListData()
            this.toast()
            this.dialogVisible = false
          })
        }
      })
    },
    clearDialogForm () {
      this.form.groupVal = ''
      this.form.roleId = ''
      this.modifyGroup = false
      this.roleDisable = false
      this.$refs.formValidate.resetFields()
    }
  }
}
</script>
