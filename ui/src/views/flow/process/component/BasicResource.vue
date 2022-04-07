<template>
  <el-form ref="formValidate"
           :model="form"
           :rules="ruleValidate"
           label-position="top"
           :disabled="readonly"
           class="node-parameter-bar">
    <el-form-item :label="$t('DI.imageSettings')"
                  prop="imageType">
      <el-radio-group v-model="form.imageType">
        <el-radio label="Standard">
          {{ $t('DI.standard') }}
        </el-radio>
        <el-radio label="Custom">
          {{ $t('DI.custom') }}
        </el-radio>
      </el-radio-group>
    </el-form-item>
    <el-form-item v-if="form.imageType === 'Custom'"
                  key="imageInput"
                  :label="$t('DI.imageSelection')"
                  prop="imageInput">
      <span>&nbsp;{{ FesEnv.DI.defineImage }}&nbsp;&nbsp;</span>
      <el-input v-model="form.imageInput"
                :placeholder="$t('DI.imageInputPro')" />
    </el-form-item>
    <el-form-item v-else
                  key="imageOption"
                  :label="$t('DI.imageSelection')"
                  prop="imageOption">
      <span>&nbsp;{{ FesEnv.DI.defineImage }}&nbsp;&nbsp;</span>
      <el-select v-model="form.imageOption"
                 :placeholder="$t('DI.imagePro')">
        <el-option v-for="(item,index) in imageOptionList"
                   :label="item"
                   :key="index"
                   :value="item">
        </el-option>
      </el-select>
    </el-form-item>
    <div class="subtitle">
      {{ $t('DI.computingResource') }}
    </div>
    <el-form-item :label="$t('ns.nameSpace')"
                  prop="namespace">
      <el-select v-model="form.namespace"
                 filterable
                 :placeholder="$t('ns.nameSpacePro')">
        <el-option v-for="item in spaceOptionList"
                   :label="item"
                   :key="item"
                   :value="item">
        </el-option>
      </el-select>
    </el-form-item>
    <el-form-item label="CPU"
                  prop="cpus">
      <el-input v-model="form.cpus"
                :placeholder="$t('DI.CPUPro')">
        <span slot="append">Core</span>
      </el-input>
    </el-form-item>
    <el-form-item :label="$t('AIDE.memory1')"
                  prop="memory">
      <el-input v-model="form.memory"
                :placeholder="$t('DI.memoryPro')">
        <span slot="append">Gb</span>
      </el-input>
    </el-form-item>
    <el-form-item label="GPU"
                  prop="gpus">
      <el-input v-model="form.gpus"
                :placeholder="$t('DI.GPUPro')">
        <span slot="append">{{ $t('DI.block') }}</span>
      </el-input>
    </el-form-item>
  </el-form>
</template>
<script>
import modelMixin from '../../../../util/modelMixin'
import nodeFormMixin from './nodeFormMixin'
const initForm = {
  imageType: 'Standard',
  defineImage: '',
  imageOption: '',
  namespace: '',
  cpus: '',
  gpus: '',
  memory: ''
}
export default {
  mixins: [modelMixin, nodeFormMixin],
  data () {
    return {
      imageOptionList: this.FesEnv.DI.imageOption,
      spaceOptionList: [],
      form: { ...initForm }
    }
  },
  props: {
    readonly: {
      type: Boolean,
      default: false
    }
  },
  computed: {
    ruleValidate: function () {
      return {
        namespace: [
          { required: true, message: this.$t('ns.nameSpaceReq') }
        ],
        imageType: [
          { required: true, message: this.$t('DI.imageTypeReq') }
        ],
        imageOption: [
          { required: true, message: this.$t('DI.imageColumnReq') }
        ],
        imageInput: [
          { required: true, message: this.$t('DI.imageInputReq') },
          { pattern: new RegExp(/^[[a-zA-Z0-9][a-zA-Z0-9-._]*$/), message: this.$t('DI.imageInputFormat') }
        ],
        cpus: [
          { required: true, message: this.$t('DI.CPUReq') },
          { pattern: new RegExp(/^((0{1})([.]\d{1})|([1-9]\d*)([.]\d{1})?)$/), message: this.$t('AIDE.CPUNumberReq') }
        ],
        gpus: [
          { pattern: new RegExp(/^(([1-9]{1}\d*([.]0)*)|([0]{1}))$/), message: this.$t('DI.GPUFormat') }
        ],
        memory: [
          { required: true, message: this.$t('DI.memoryReq') },
          { pattern: new RegExp(/^[1-9]\d*([.]0)*$/), message: this.$t('AIDE.memoryNumberReq') }
        ]
      }
    }
  },
  created () {
    const header = this.handlerFromDssHeader('get')
    this.getDISpaceOption(header)
  },
  methods: {
    initFormData () {
      this.form = { ...initForm }
      setTimeout(() => {
        this.$refs.formValidate.clearValidate()
      }, 0)
    },
    save () {
      return new Promise(resolve => {
        this.$refs.formValidate.validate((valite) => {
          if (valite) {
            let param = {
              framework: {
                name: this.FesEnv.DI.defineImage,
                version: this.form.imageType === 'Standard' ? this.form.imageOption : this.form.imageInput
              },
              gpus: +this.form.gpus,
              cpus: +this.form.cpus,
              memory: this.form.memory + 'Gb'
            }
            let keyArr = ['imageType', 'namespace']
            for (let key of keyArr) {
              param[key] = this.form[key]
            }
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
