<template>
  <div>
    <breadcrumb-nav>
      <el-button type="primary"
                 icon="el-icon-plus"
                 size="small"
                 @click="dialogVisible=true">
        {{$t('user.addProxyUser')}}
      </el-button>
    </breadcrumb-nav>
    <el-table :data="dataList"
              border
              :empty-text="$t('common.noData')">
      <el-table-column prop="name"
                       :label="$t('user.proxyUserName')" />
      <el-table-column prop="uid"
                       :label="$t('user.uid')" />
      <el-table-column prop="gid"
                       :label="$t('user.gid')" />
      <el-table-column prop="path"
                       :label="$t('user.storagePath')" />
      <el-table-column prop="isActivated"
                       :label="$t('user.isActive')"
                       :formatter="translationActive" />
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
                        :disabled="modifyProxy"
                        :placeholder="$t('user.userNamePro')" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item :label="$t('user.isActive')"
                          prop="isActivated">
              <el-select v-model="form.isActivated"
                         :placeholder="$t('user.activePro')">
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
        <el-form-item :label="$t('user.storagePath')"
                      prop="path">
          <el-input v-model="form.path"
                    :placeholder="$t('user.storagePathPro')"></el-input>
        </el-form-item>
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
      userName: '',
      dialogVisible: false,
      modifyProxy: false,
      form: {
        user_id: '',
        name: '',
        gid: '',
        uid: '',
        path: '',
        isActivated: 1,
        remarks: ''
      },
      proxyId: '',
      dataList: []
    }
  },
  mounted () {
    this.form.user_id = Number(this.$route.query.userId)
    this.getListData()
  },
  computed: {
    title () {
      return this.modifyProxy ? this.$t('user.modifyProxyUser') : this.$t('user.addProxyUser')
    },
    getUrl () {
      return `/cc/${this.FesEnv.ccApiVersion}/proxyUsers/${this.$route.query.username}`
    },
    ruleValidate () {
      // 切换语言表单报错重置表单
      this.resetFormFields()
      return {
        name: [
          { required: true, message: this.$t('user.userNameReq') },
          { pattern: new RegExp(/^[a-z0-9][a-z0-9_]*$/), message: this.$t('user.proxyNameFormat') }
        ],
        gid: [
          { required: true, message: this.$t('user.gidReq') },
          { pattern: new RegExp(/^\d*$/), message: this.$t('user.gidNumberPro') }
        ],
        uid: [
          { required: true, message: this.$t('user.uidReq') },
          { pattern: new RegExp(/^\d*$/), message: this.$t('user.uidNumberPro') }
        ],
        isActivated: [
          { required: true, message: this.$t('user.activeReq') }
        ],
        remarks: [
          { pattern: new RegExp(/^.{0,200}$/), message: this.$t('user.remarksMaxLength') }
        ],
        path: [
          { required: true, message: this.$t('user.storagePathReq') },
          { pattern: new RegExp(/^\/[0-9a-zA-Z-_/*]*$/), message: this.$t('user.StoragePathFormat') }
        ]
      }
    }
  },
  methods: {
    translationActive (value, row) {
      return value.isActivated ? this.$t('common.yes') : this.$t('common.no')
    },
    updateItem (trData) {
      this.modifyProxy = true
      this.proxyId = trData.id
      setTimeout(() => {
        Object.assign(this.form, trData)
      }, 0)
      this.dialogVisible = true
    },
    deleteItem (trData) {
      const url = `/cc/${this.FesEnv.ccApiVersion}/proxyUser/${trData.id}`
      this.deleteListItem(url)
    },
    subInfo () {
      this.$refs.formValidate.validate((valid) => {
        if (valid) {
          this.setBtnDisabeld()
          let method = 'post'
          let url = `/cc/${this.FesEnv.ccApiVersion}/proxyUser`
          if (this.modifyProxy) {
            method = 'put'
            url += '/' + this.proxyId
          }
          let param = { ...this.form }
          param.gid = +param.gid
          param.uid = +param.uid
          this.FesApi.fetch(url, param, method).then(() => {
            this.getListData()
            this.toast()
            this.dialogVisible = false
          })
        }
      })
    },
    clearDialogForm () {
      this.modifyProxy = false
      this.$refs.formValidate.resetFields()
    }
  }
}
</script>
