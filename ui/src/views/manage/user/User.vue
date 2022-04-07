<template>
  <div>
    <breadcrumb-nav class="margin-b20">
      <el-button type="primary"
                 icon="el-icon-plus"
                 @click="dialogVisible = true"
                 size="small">
        {{ $t('user.addUser') }}
      </el-button>
    </breadcrumb-nav>
    <el-table :data="dataList"
              style="width: 100%"
              border
              :empty-text="$t('common.noData')">
      <el-table-column prop="name"
                       :label="$t('user.userName')" />
      <el-table-column prop="uid"
                       :label="$t('user.uid')" />
      <el-table-column prop="gid"
                       :label="$t('user.gid')" />
      <el-table-column prop="type"
                       :label="$t('user.userType')" />
      <el-table-column prop="enableFlag"
                       :label="$t('user.isActive')"
                       :formatter="translationActive" />
      <el-table-column prop="guidCheckT"
                       :label="$t('user.isGuidCheck')"
                       :formatter="translationGuidCheck" />
      <el-table-column prop="remarks"
                       :label="$t('user.describe')"
                       width="350px" />
      <el-table-column :label="$t('common.operation')"
                       width="120px">
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
              <el-dropdown-item :command="beforeHandleCommand('settingGroup',scope.row)">
                {{$t('user.userGroupSettings')}}
              </el-dropdown-item>
              <el-dropdown-item :command="beforeHandleCommand('settingProxy',scope.row)">
                {{$t('user.proxyUserSetting')}}
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
               width="920px">
      <el-form ref="formValidate"
               :model="form"
               :rules="ruleValidate"
               label-width="120px">
        <el-row>
          <el-col :span="12">
            <el-form-item :label="$t('user.userName')"
                          prop="name">
              <el-input v-model="form.name"
                        :placeholder="$t('user.userNamePro')"
                        :disabled="modify" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item :label="$t('user.userType')"
                          prop="type">
              <el-select v-model="form.type"
                         :placeholder="$t('user.typePro')">
                <el-option value="USER">
                  USER
                </el-option>
                <el-option value="SYSTEM">
                  SYSTEM
                </el-option>
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row>
          <el-col :span="12">
            <el-form-item :label="$t('user.uid')"
                          prop="uid">
              <el-input v-model="form.uid"
                        :placeholder="$t('user.uidPro')" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item :label="$t('user.gid')"
                          prop="gid">
              <el-input v-model="form.gid"
                        :placeholder="$t('user.gidDPro')" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row>
          <el-col :span="12">
            <el-form-item :label="$t('user.isActive')"
                          prop="enableFlag">
              <el-select v-model="form.enableFlag"
                         :placeholder="$t('user.activePro')">
                <el-option :label="$t('common.yes')"
                           :value="1">
                </el-option>
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item :label="$t('user.guidCheck')"
                          prop="guidCk">
              <el-select v-model="form.guidCk"
                         :placeholder="$t('user.guidCheckPro')">
                <el-option :label="$t('common.yes')"
                           :value="1">
                </el-option>
                <el-option :label="$t('common.no')"
                           :value="0">
                </el-option>
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item :label="$t('user.describe')"
                      prop="remarks">
          <el-input v-model="form.remarks"
                    type="textarea"
                    :placeholder="$t('user.describePro')" />
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
      dialogVisible: false,
      form: {
        name: '',
        gid: '',
        uid: '',
        guidCk: '',
        guidCheck: '',
        type: '',
        enableFlag: 1,
        remarks: ''
      },
      modify: false,
      dataList: []
    }
  },
  mounted () {
    this.getListData()
  },
  computed: {
    title () {
      return this.modify ? this.$t('user.modifyUser') : this.$t('user.newUser')
    },
    getUrl () {
      return `/cc/${this.FesEnv.ccApiVersion}/users`
    },
    ruleValidate () {
      // 切换语言表单报错重置表单
      this.resetFormFields()
      return {
        name: [
          { required: true, message: this.$t('user.userNameReq') },
          { pattern: new RegExp(/^.*[^\s][\s\S]*$/), message: this.$t('user.userNameReq') }
        ],
        gid: [
          { required: true, message: this.$t('user.gidReq') },
          { pattern: new RegExp(/^\d*$/), message: this.$t('user.gidNumberPro') }
        ],
        uid: [
          { required: true, message: this.$t('user.uidReq') },
          { pattern: new RegExp(/^\d*$/), message: this.$t('user.uidNumberPro') }
        ],
        guidCk: [
          { required: true, message: this.$t('user.guidCheckReq') }
        ],
        enableFlag: [
          { required: true, message: this.$t('user.activeReq') }
        ],
        type: [
          { required: true, message: this.$t('user.typeReq') }
        ],
        remarks: [
          { pattern: new RegExp(/^.{0,200}$/), message: this.$t('user.remarksMaxLength') }
        ]
      }
    }
  },
  methods: {
    translationActive (value, row) {
      return value.enableFlag ? this.$t('common.yes') : this.$t('common.no')
    },
    translationGuidCheck (value, row) {
      return value.guidCheck === '1' ? this.$t('common.yes') : this.$t('common.no')
    },
    updateItem (trData) {
      this.dialogVisible = true
      this.modify = true
      setTimeout(() => {
        trData.guidCk = trData.guidCheck === '1' ? 1 : 0
        Object.assign(this.form, trData)
      }, 0)
    },
    deleteItem (trData) {
      const url = `/cc/${this.FesEnv.ccApiVersion}/users/id/${trData.id}`
      this.deleteListItem(url)
    },
    settingGroup (trData) {
      this.$router.push({
        name: 'settingUserGroup',
        params: {
          userId: trData.id
        },
        query: {
          userName: trData.name
        }
      })
    },
    subInfo () {
      this.form.guidCheck = !!parseInt(this.form.guidCk)
      this.$refs.formValidate.validate((valid) => {
        if (valid) {
          this.setBtnDisabeld()
          let method = 'post'
          let url = `/cc/${this.FesEnv.ccApiVersion}/users`
          if (this.modify) {
            method = 'put'
          }
          let data = { ...this.form }
          data.gid = parseInt(data.gid)
          data.uid = parseInt(data.uid)
          data.guidCheck = this.form.guidCheck ? '1' : '0'
          delete data.guidCk
          this.FesApi.fetch(url, data, method).then(() => {
            this.getListData()
            this.toast()
            this.dialogVisible = false
          })
        }
      })
    },
    settingProxy (trData) {
      this.$router.push({
        name: 'settingUserProxy',
        query: {
          username: trData.name,
          userId: trData.id
        }
      })
    }
  }
}
</script>
