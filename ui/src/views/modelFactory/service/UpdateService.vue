<template>
  <el-dialog :title="$t('serviceList.updateModelService')"
             :visible.sync="dialogVisible"
             @close="clearDialogForm"
             custom-class="dialog-style"
             width="1100px">
    <el-form ref="formValidate"
             :model="form"
             :rules="formValidate"
             class="add-form"
             label-width="135px">
      <div class="subtitle">
        {{$t('DI.basicSettings')}}
      </div>
      <el-form-item :label="$t('serviceList.modelName')"
                    prop="service_name">
        <el-input v-model="form.service_name"
                  disabled=""
                  :placeholder="$t('serviceList.modelNamePro')" />
      </el-form-item>
      <el-form-item :label="$t('serviceList.modelDescription')">
        <el-input v-model="form.remark"
                  :placeholder="$t('serviceList.modelDescriptionPro')"
                  type="textarea" />
      </el-form-item>
      <div class="subtitle">
        {{$t('DI.computingResource')}}
      </div>
      <el-row>
        <el-col :span="12">
          <el-form-item label="CPU"
                        prop="cpu">
            <el-input v-model="form.cpu"
                      :placeholder="$t('DI.CPUPro')">
              <span slot="append">Core</span>
            </el-input>
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item :label="$t('AIDE.memory1')"
                        prop="memory">
            <el-input v-model="form.memory"
                      :placeholder="$t('DI.memoryPro')">
              <span slot="append">Gi</span>
            </el-input>
          </el-form-item>
        </el-col>
      </el-row>
      <el-form-item label="GPU"
                    prop="gpu">
        <el-input v-model="form.gpu"
                  :placeholder="$t('DI.GPUPro')">
          <span slot="append">{{$t('AIDE.block')}}</span>
        </el-input>
      </el-form-item>
      <div class="subtitle">
        {{$t('serviceList.deploymentInfo')}}
      </div>
      <model-setting v-if="dialogVisible"
                     ref="modelSetting"
                     :group-list="groupList"
                     :form-data="form"
                     :modify="true"></model-setting>
    </el-form>
    <span slot="footer"
          class="dialog-footer">
      <el-button type="primary"
                 @click="subInfo"
                 :disabled="btnDisabled">
        {{ $t('common.save') }}
      </el-button>
      <el-button @click="dialogVisible=false">
        {{ $t('common.cancel') }}
      </el-button>
    </span>
  </el-dialog>
</template>
<script>
import ModelSetting from './ModelSetting'
export default {
  components: { ModelSetting },
  props: {
    groupList: {
      type: Array,
      default () {
        return []
      }
    }
  },
  data () {
    return {
      dialogVisible: false,
      form: {
        service_name: '',
        remark: '',
        cpu: '',
        memory: '',
        gpu: '',
        modelversion_id: '',
        modelversion: '',
        source: '',
        image: '',
        service_id: '',
        namespace: '',
        endpoint_type: ''
      },
      imageOptionList: [],
      resetSource: false,
      btnDisabled: false
    }
  },
  computed: {
    formValidate () {
      this.resetFormFields()
      return {
        cpu: [
          { required: true, message: this.$t('DI.CPUReq') },
          { pattern: new RegExp(/^((0{1})([.]\d{1})|([1-9]\d*)([.]\d{1})?)$/), message: this.$t('AIDE.CPUNumberReq') }
        ],
        memory: [
          { required: true, message: this.$t('DI.memoryReq') },
          { pattern: new RegExp(/^[1-9]\d*([.]0)*$/), message: this.$t('AIDE.memoryNumberReq') }
        ],
        gpu: [
          { pattern: new RegExp(/^(([1-9]{1}\d*([.]0)*)|([0]{1}))$/), message: this.$t('AIDE.GPUFormat') }
        ],
        modelversion_id: [{ required: true, message: this.$t('serviceList.modelVersionReq') }]
      }
    }
  },
  methods: {
    subInfo () {
      this.$refs.formValidate.validate((valid) => {
        if (valid) {
          const modelSetting = this.$refs.modelSetting
          modelSetting.$refs.formValidate.validate((valid) => {
            if (valid) {
              this.btnDisabled = true
              let param = { ...this.form }
              param.gpu = param.gpu + ''
              param.cpu = +Number(param.cpu).toFixed(1)
              param.memory = +param.memory
              const modelSettingFrom = modelSetting.form
              modelSetting.byVersionIdFindSource()
              const modelKeys = ['modelversion', 'image', 'source']
              for (let key of modelKeys) {
                param[key] = modelSettingFrom[key]
              }
              param.modelversion_id = parseInt(modelSettingFrom.modelversion_id)
              delete param.model_type
              this.FesApi.fetch(`/mf/v1/service/${this.form.service_id}`, param, 'put').then(() => {
                this.btnDisabled = false
                this.dialogVisible = false
                this.toast()
                this.$emit('getListData')
              }).finally(() => {
                this.btnDisabled = false
              })
            }
          })
        }
      })
    },
    changeVersionId (versionId) {
      this.FesApi.fetch(`/mf/v1/images/${versionId}`, {}, 'get').then((res) => {
        if (res.length > 0) {
          this.imageOptionList = this.formatResList(res)
        } else {
          this.form.image = ''
          this.imageOptionList = []
        }
      })
    },
    clearDialogForm () {
      this.btnDisabled = false
      this.resetSource = false
      for (let key in this.form) {
        this.form[key] = ''
      }
      this.$refs.formValidate.resetFields()
    }
  }
}
</script>
<style lang="scss" scoped>
.add-form {
  height: 500px;
}
</style>
