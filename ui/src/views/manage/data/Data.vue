<template>
  <div>
    <breadcrumb-nav class="margin-b20">
      <el-button type="primary"
                 icon="el-icon-plus"
                 @click="dialogVisible = true"
                 size="small">
        {{ $t('data.addDataStore') }}
      </el-button>
    </breadcrumb-nav>
    <el-table :data="dataList"
              border
              style="width: 100%"
              :empty-text="$t('common.noData')">
      <el-table-column prop="path"
                       :show-overflow-tooltip="true"
                       :label="$t('data.dataStorePath')" />
      <el-table-column prop="isActive"
                       :label="$t('user.isActive')"
                       :formatter="translationActive" />
      <el-table-column prop="remarks"
                       :show-overflow-tooltip="true"
                       :label="$t('DI.description')" />
      <el-table-column :label="$t('common.operation')">
        <template slot-scope="scope">
          <el-button @click="updateItem(scope.row)"
                     type="text"
                     size="small">{{$t('common.modify')}}</el-button>
          <el-button @click="deleteItem(scope.row)"
                     type="text"
                     size="small">{{$t('common.delete')}}</el-button>
          <el-button @click="settingGroup(scope.row)"
                     type="text"
                     size="small">{{$t('user.userGroupSettings')}}</el-button>
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
               width="800px">
      <el-form ref="formValidate"
               :model="form"
               :rules="ruleValidate"
               label-width="175px">
        <el-form-item :label="$t('data.dataStorePath')"
                      prop="path">
          <el-input v-model="form.path"
                    :placeholder="$t('data.dataStorePro')"
                    :disabled="modify" />
        </el-form-item>
        <el-form-item :label="$t('user.isActive')"
                      prop="isActive">
          <el-select v-model="form.isActive"
                     :placeholder="$t('user.activePro')">
            <el-option :label="$t('common.yes')"
                       :value="1">
            </el-option>
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('data.describe')"
                      prop="remarks">
          <el-input v-model="form.remarks"
                    type="textarea"
                    :placeholder="$t('data.describePro')" />
        </el-form-item>
      </el-form>
      <span slot="footer"
            class="dialog-footer">
        <el-button type="primary"
                   :disabled="btnDisabled"
                   @click="subInfo">{{ $t("common.save") }}</el-button>
        <el-button @click="dialogVisible=false">{{ $t("common.cancel") }}</el-button>
      </span>
    </el-dialog>
  </div>
</template>
<script type="text/ecmascript-6">
export default {
  data: function () {
    return {
      dialogVisible: false,
      form: {
        remarks: '',
        isActive: 1,
        path: ''
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
      return this.modify ? this.$t('data.modifyDataStore') : this.$t('data.newDataStore')
    },
    getUrl () {
      return `/cc/${this.FesEnv.ccApiVersion}/storages`
    },
    ruleValidate () {
      // 切换语言表单报错重置表单
      this.resetFormFields()
      return {
        path: [
          { required: true, message: this.$t('data.pathReq') },
          { pattern: new RegExp(/^\/[0-9a-zA-Z-_/*]*$/), message: this.$t('data.pathStartReq') }
        ],
        isActive: [
          { required: true, message: this.$t('user.activeReq') }
        ],
        remarks: [
          { pattern: new RegExp(/^.{0,200}$/), message: this.$t('data.remarksMaxLength') }
        ]
      }
    }
  },
  methods: {
    translationActive (value, row) {
      return this.$t('common.yes')
    },
    updateItem (trData) {
      this.dialogVisible = true
      this.modify = true
      setTimeout(() => {
        Object.assign(this.form, trData)
      }, 10)
    },
    deleteItem (trData) {
      const url = `/cc/${this.FesEnv.ccApiVersion}/storages/path`
      this.deleteListItem(url, { path: trData.path })
    },
    settingGroup (trData) {
      this.$router.push({
        name: 'settingDataGroup',
        params: {
          storageId: trData.id + ''
        },
        query: {
          path: trData.path
        }
      }
      )
    },
    subInfo () {
      this.$refs.formValidate.validate((valid) => {
        if (valid) {
          this.setBtnDisabeld()
          let method = this.modify ? 'put' : 'post'
          this.FesApi.fetch(`/cc/${this.FesEnv.ccApiVersion}/storages`, this.form, method).then(() => {
            this.getListData()
            this.toast()
            this.dialogVisible = false
          })
        }
      })
    }
  }
}
</script>
