<template>
  <el-form ref="formValidate"
           :model="form"
           :disabled="readonly"
           :rules="ruleValidate"
           label-position="top"
           class="node-parameter-bar">
    <div class="subtitle">
      {{$t('flow.modelParameters')}}
    </div>
    <el-form-item :label="$t('user.userGroup')"
                  prop="group_id">
      <el-select v-model="form.group_id"
                 :filterable="true"
                 :placeholder="$t('user.userGroupSettingsPro')">
        <el-option v-for="(item,index) of groupList"
                   :label="item.name"
                   :value="item.id"
                   :key="index">
        </el-option>
      </el-select>
    </el-form-item>
    <el-form-item :label="$t('modelList.modelName')"
                  prop="model_name">
      <el-input v-model="form.model_name"
                :placeholder="$t('modelList.modelNamePro')"></el-input>
    </el-form-item>
    <el-form-item :label="$t('flow.modelPush')">
      <el-radio v-model="form.model_push"
                :label="true">{{$t('common.yes')}}</el-radio>
      <el-radio v-model="form.model_push"
                :label="false">{{$t('common.no')}}</el-radio>
    </el-form-item>
    <!-- <el-switch v-model="form.value1"
               inactive-text="模型推送">
    </el-switch> -->
    <el-form-item v-if="form.model_push"
                  :label="$t('flow.factoryName')"
                  prop="factory_name">
      <el-input v-model="form.factory_name"
                :placeholder="$t('flow.factoryNamePro')"></el-input>
    </el-form-item>
  </el-form>
</template>
<script>
import modelMixin from '../../../../util/modelMixin'
import nodeFormMixin from './nodeFormMixin'
const initForm = {
  group_id: '',
  model_name: '',
  model_push: false,
  factory_name: ''
}
export default {
  mixins: [modelMixin, nodeFormMixin],
  data () {
    return {
      groupList: [],
      form: { ...initForm }
    }
  },
  props: {
    readonly: {
      type: Boolean,
      default: false
    }
  },
  created () {
    const header = this.handlerFromDssHeader('get')
    this.getGroupList(header)
  },
  computed: {
    ruleValidate () {
      return {
        group_id: [
          { required: true, message: this.$t('user.userGroupSettingsReq') }
        ],
        model_name: [
          { required: true, message: this.$t('modelList.modelNameReq') },
          { type: 'string', pattern: new RegExp(/^[a-zaA-Z][a-zA-Z0-9-]*$/), message: this.$t('modelList.modelNameFormat') }
        ],
        model_push: [
          { required: true, message: this.$t('flow.modelPushReq') }
        ],
        factory_name: [
          { required: true, message: this.$t('flow.factoryNameReq') },
          { type: 'string', pattern: new RegExp(/^[\u4e00-\u9fa5_0-9a-zA-Z-]*$/), message: this.$t('flow.factoryNameFormat') }
        ]
      }
    }
  },
  methods: {
    initFormData () {
      this.form = { ...initForm }
      setTimeout(() => {
        this.$refs.formValidate.clearValidate()
      }, 0)
    },
    save () {
      return new Promise((resolve) => {
        this.$refs.formValidate.validate((valid) => {
          if (valid) {
            let param = { ...this.form }
            if (!param.model_push) {
              param.factory_name = ''
            }
            delete param.model_push
            resolve(param)
          } else {
            resolve(false)
          }
        })
      })
    }
  }
}
</script>
