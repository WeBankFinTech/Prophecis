<template>
  <div>
    <div v-show="activeName==='0'"
         key="'step0'">
      <data-setting ref="dataSetting"
                    :readonly="readonly"></data-setting>
    </div>
    <div v-show="activeName==='1'"
         key="step1">
      <custom-model ref="customModel"
                    :readonly="readonly"></custom-model>
      <el-form ref="formValidate"
               :model="form"
               :disabled="readonly"
               :rules="ruleValidate"
               label-position="top"
               class="node-parameter-bar">
        <div class="subtitle">
          {{$t('flow.trainingParameters')}}
        </div>
        <el-form-item :label="$t('ALG.penalty')">
          <el-select v-model="form.penalty"
                     clearable>
            <el-option label="l1"
                       value="l1">
            </el-option>
            <el-option label="l2"
                       value="l2">
            </el-option>
            <el-option label="elasticnet"
                       value="elasticnet">
            </el-option>
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('ALG.dual')">
          <el-select v-model="form.dual"
                     clearable>
            <el-option :label="$t('common.yes')"
                       :value="true"></el-option>
            <el-option :label="$t('common.no')"
                       :value="false"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('ALG.tol')"
                      prop="tol">
          <el-input v-model="form.tol"></el-input>
        </el-form-item>
        <el-form-item :label="$t('ALG.C')"
                      prop="C">
          <el-input v-model="form.C"></el-input>
        </el-form-item>
        <el-form-item :label="$t('ALG.fit_intercept')">
          <el-select v-model="form.fit_intercept"
                     clearable>
            <el-option :label="$t('common.yes')"
                       :value="true"></el-option>
            <el-option :label="$t('common.no')"
                       :value="false"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('ALG.intercept_scaling')"
                      v-show="form.solver==='liblinear'&&form.fit_intercept===true"
                      prop="intercept_scaling">
          <el-input v-model="form.intercept_scaling"></el-input>
        </el-form-item>
        <el-form-item :label="$t('ALG.class_weight')">
          <el-select v-model="form.class_weight"
                     clearable>
            <el-option label="dict"
                       value="dict"></el-option>
            <el-option label="balanced"
                       value="falbalancedse"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('ALG.solver')">
          <el-select v-model="form.solver"
                     clearable>
            <el-option label="newton-cg"
                       value="newton-cg"></el-option>
            <el-option label="lbfgs"
                       value="lbfgs"></el-option>
            <el-option label="liblinear"
                       value="liblinear"></el-option>
            <el-option label="sag"
                       value="sag"></el-option>
            <el-option label="saga"
                       value="saga"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('ALG.random_state')"
                      v-if="form.solver=='sag'||form.solver=='saga'||form.solver=='liblinear'">
          <el-input v-model="form.random_state"></el-input>
        </el-form-item>
        <el-form-item :label="$t('ALG.max_iter')"
                      prop="max_iter">
          <el-input v-model="form.max_iter"></el-input>
        </el-form-item>
        <el-form-item :label="$t('ALG.multi_class')">
          <el-select v-model="form.multi_class"
                     clearable>
            <el-option label="auto"
                       value="auto"></el-option>
            <el-option label="ovr"
                       value="ovr"></el-option>
            <el-option label="multinomial"
                       value="multinomial"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('ALG.verbose')"
                      prop="verbose">
          <el-input v-model="form.verbose"></el-input>
        </el-form-item>
        <el-form-item :label="$t('ALG.n_jobs')"
                      prop="n_jobs">'
          <el-input v-model="form.n_jobs">
          </el-input>
        </el-form-item>
        <el-form-item :label="$t('ALG.l1_ratio')"
                      v-show="form.penalty==='elasticnet'"
                      prop="l1_ratio">
          <el-input v-model="form.l1_ratio">
          </el-input>
        </el-form-item>
        <el-form-item :label="$t('ALG.sample_weight')"
                      prop="sample_weight">
          <el-input v-model="form.sample_weight"></el-input>
        </el-form-item>
      </el-form>
    </div>
    <div v-show="activeName==='2'"
         key="step2">
      <basic-resource ref="basicResource"
                      :readonly="readonly"></basic-resource>
    </div>
    <div class="save-button"
         v-if="!readonly">
      <el-button @click="save"
                 :disabled="nodeSaveBtnDisable">{{ $t('common.save') }}</el-button>
    </div>
  </div>
</template>
<script>
import DataSetting from './DataSetting'
import BasicResource from './BasicResource'
import nodeFormMixin from './nodeFormMixin'
import CustomModel from './CustomModel'
const initForm = {
  penalty: '',
  dual: '',
  tol: '',
  C: '',
  fit_intercept: '',
  intercept_scaling: '',
  class_weight: '',
  random_state: '',
  solver: '',
  max_iter: '',
  multi_class: '',
  verbose: '',
  n_jobs: '',
  l1_ratio: '',
  sample_weight: ''
}
export default {
  mixins: [nodeFormMixin],
  components: {
    DataSetting,
    BasicResource,
    CustomModel
  },
  props: {
    nodeData: {
      type: Object,
      default () {
        return {}
      }
    },
    activeName: {
      type: String,
      default: '0'
    },
    readonly: {
      type: Boolean,
      default: false
    },
    nodeSaveBtnDisable: {
      type: Boolean,
      default: false
    }
  },
  data () {
    return {
      form: { ...initForm }
    }
  },
  computed: {
    ruleValidate () {
      return {
        // 由于这些参数都不是必填，值可能被置空，type设置为string，提交的时候在转化成number
        tol: [
          { type: 'string', pattern: new RegExp(/^(([1-9]{1}\d*)|(0{1}))(\.\d{0,7})$/), message: this.$t('ALG.floatingFormat') }
        ],
        C: [
          { type: 'string', pattern: new RegExp(/^(([1-9]{1}\d*)|(0{1}))(\.\d{0,7})$/), message: this.$t('ALG.floatingFormat') }
        ],
        max_iter: [
          { type: 'string', pattern: new RegExp(/^(([1-9]{1}\d*)|(0{1}))$/), message: this.$t('ALG.integerFormat') }
        ],
        intercept_scaling: [
          { type: 'string', pattern: new RegExp(/^(([1-9]{1}\d*)|(0{1}))(\.\d{0,7})?$/), message: this.$t('ALG.integerOrFloat') }
        ],
        verbose: [
          { type: 'string', pattern: new RegExp(/^(([1-9]{1}\d*)|(0{1}))$/), message: this.$t('ALG.integerFormat') }
        ],
        n_jobs: [
          { type: 'string', pattern: new RegExp(/^(([1-9-]{1}\d*)|(0{1}))$/), message: this.$t('ALG.integerFormat1') }
        ],
        l1_ratio: [
          { type: 'string', pattern: new RegExp(/^(((1{1})|(0{1})|(1.0{1}))|((0{1}))(\.\d{0,7}))$/), message: this.$t('ALG.maximum1') }
        ],
        sample_weight: [
          { type: 'string', pattern: new RegExp(/^[^\u4e00-\u9fa5]*$/), message: this.$t('ALG.noCCharacters') }
        ]
      }
    }
  },
  watch: {
    nodeData: {
      immediate: true,
      handler (val) {
        this.handleOriginNodeData(val, initForm)
      }
    }
  },
  methods: {
    algorithmDataFormatKey () {
      return {
        float: ['tol', 'C', 'intercept_scaling', 'l1_ratio'],
        int: ['max_iter', 'verbose', 'n_jobs']
      }
    },
    save (evt, temSave = false) {
      this.algorithmSubmit(initForm, 'LogisticRegression', temSave)
    }
  }
}
</script>
