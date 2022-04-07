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
      <el-table-column prop="namespace"
                       :label="$t('ns.nameSpaceName')" />
      <el-table-column prop="groupName"
                       :label="$t('user.userGroup')" />
      <el-table-column :label="$t('common.operation')"
                       width="80px">
        <template slot-scope="scope">
          <el-button @click="deleteItem(scope.row)"
                     type="text"
                     size="small">{{$t('common.delete')}}</el-button>
        </template>
      </el-table-column>
    </el-table>
    <el-dialog :title="$t('user.newGroup')"
               :visible.sync="dialogVisible"
               @close="clearDialogForm"
               custom-class="dialog-style"
               width="600px">
      <el-form ref="formValidate"
               :model="form"
               :rules="ruleValidate"
               label-width="160px">
        <el-form-item :label="$t('ns.nameSpaceName')">
          <el-input v-model="namespace"
                    :placeholder="$t('ns.nsNamePro')"
                    disabled />
        </el-form-item>
        <el-form-item :label="$t('user.userGroupSettings')"
                      prop="groupId">
          <el-select v-model="form.groupId"
                     filterable
                     :placeholder="$t('user.userGroupSettingsPro')">
            <el-option v-for="item in groupOptionList"
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
  props: {
    namespace: {
      type: String,
      default: ''
    }
  },
  data () {
    return {
      dialogVisible: false,
      namespaceId: '',
      form: {
        namespace: '',
        groupId: '',
        namespaceId: ''
      },
      dataList: [],
      groupOptionList: []
    }
  },
  mounted () {
    this.namespaceId = this.$route.query.nid
    this.getListData()
    this.getGroupOption()
  },
  computed: {
    getUrl () {
      return `/cc/${this.FesEnv.ccApiVersion}/groups/groupNamespace/${this.namespaceId}`
    },
    ruleValidate () {
      // 切换语言表单报错重置表单
      this.resetFormFields()
      return {
        groupId: [
          { required: true, message: this.$t('user.userGroupSettingsReq') }
        ]
      }
    }
  },
  methods: {
    getListData () {
      this.FesApi.fetch(this.getUrl, 'get').then(rst => {
        this.dataList = rst.models || []
      })
    },
    deleteItem (trData) {
      const url = `/cc/${this.FesEnv.ccApiVersion}/groups/namespaces/id/${trData.id}`
      this.deleteListItem(url)
    },
    subInfo () {
      this.$refs.formValidate.validate((valid) => {
        if (valid) {
          this.setBtnDisabeld()
          let param = { ...this.form }
          param.namespace = this.namespace
          param.namespaceId = parseInt(this.namespaceId)
          this.FesApi.fetch(`/cc/${this.FesEnv.ccApiVersion}/groups/namespaces`, param, 'post').then(() => {
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
