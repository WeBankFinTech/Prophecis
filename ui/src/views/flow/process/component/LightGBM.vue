<template>
  <div>
    <div v-show="activeName==='0'"
         key="step0">
      <data-setting ref="dataSetting"
                    :node-type="'LightGBM'"
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
        <el-form-item :label="'API'+$t('modelList.type')"
                      prop="API_type">
          <el-select v-model="form.API_type">
            <el-option label="Classifier"
                       value="Classifier">
            </el-option>
            <el-option label="Train"
                       value="Train">
            </el-option>
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('image.parameterType')">
          <el-select v-model="form.param_type">
            <el-option :label="$t('ALG.list')"
                       value="list">
            </el-option>
            <el-option :label="$t('modelList.custom')"
                       value="custom">
            </el-option>
          </el-select>
        </el-form-item>
        <div v-show="form.param_type==='list'">
          <el-form-item :label="$t('ALG.boosting_type')">
            <el-select v-model="form.boosting_type"
                       clearable>
              <el-option label="gbdt"
                         value="gbdt">
              </el-option>
              <el-option label="rf"
                         value="rf">
              </el-option>
              <el-option label="dart"
                         value="dart">
              </el-option>
              <el-option label="goss"
                         value="goss">
              </el-option>
            </el-select>
          </el-form-item>
          <el-form-item :label="$t('ALG.num_leaves')"
                        prop="num_leaves">
            <el-input v-model="form.num_leaves"></el-input>
          </el-form-item>
          <el-form-item :label="$t('ALG.max_depth')"
                        prop="max_depth">
            <el-input v-model="form.max_depth"></el-input>
          </el-form-item>
          <el-form-item :label="$t('ALG.learning_rate')"
                        prop="learning_rate">
            <el-input v-model="form.learning_rate"></el-input>
          </el-form-item>
          <el-form-item :label="$t('ALG.n_estimators1')"
                        prop="n_estimators">
            <el-input v-model="form.n_estimators"></el-input>
          </el-form-item>
          <el-form-item label="subsample_for_bin"
                        prop="subsample_for_bin">
            <el-input v-model="form.subsample_for_bin"></el-input>
          </el-form-item>
          <el-form-item :label="$t('ALG.objective')">
            <el-select v-model="form.objective"
                       clearable>
              <el-option label="regression"
                         value="regression"></el-option>
              <el-option label="regression_l1"
                         value="regression_l1"></el-option>
              <el-option label="huber"
                         value="huber"></el-option>
              <el-option label="fair"
                         value="fair"></el-option>
              <el-option label="poisson"
                         value="poisson"></el-option>
              <el-option label="quantile"
                         value="quantile"></el-option>
              <el-option label="quantile_l2"
                         value="quantile_l2"></el-option>
              <el-option label="binary"
                         value="binary"></el-option>
              <el-option label="multiclass"
                         value="multiclass"></el-option>
              <el-option label="multiclassova"
                         value="multiclassova"></el-option>
              <el-option label="xentropy"
                         value="xentropy"></el-option>
              <el-option label="xentlambda"
                         value="xentlambda"></el-option>
              <el-option label="lambdarank"
                         value="lambdarank"></el-option>
            </el-select>
          </el-form-item>
          <el-form-item :label="$t('ALG.is_unbalance')">
            <el-select v-model="form.is_unbalance"
                       clearable>
              <el-option :label="$t('common.yes')"
                         :value="true"></el-option>
              <el-option :label="$t('common.no')"
                         :value="false"></el-option>
            </el-select>
          </el-form-item>
          <el-form-item :label="$t('ALG.min_split_gain')"
                        prop="min_split_gain">
            <el-input v-model="form.min_split_gain"></el-input>
          </el-form-item>
          <el-form-item :label="$t('ALG.min_child_weight')"
                        prop="min_child_weight">
            <el-input v-model="form.min_child_weight"></el-input>
          </el-form-item>
          <el-form-item label="subsample"
                        prop="subsample">
            <el-input v-model="form.subsample"></el-input>
          </el-form-item>
          <el-form-item label="subsample_freq"
                        prop="subsample_freq">
            <el-input v-model="form.subsample_freq"></el-input>
          </el-form-item>
          <el-form-item label="colsample_bytree"
                        prop="colsample_bytree">
            <el-input v-model="form.colsample_bytree"></el-input>
          </el-form-item>
          <el-form-item label="reg_alpha"
                        prop="reg_alpha">
            <el-input v-model="form.reg_alpha"></el-input>
          </el-form-item>
          <el-form-item label="reg_lambda"
                        prop="reg_lambda">
            <el-input v-model="form.reg_lambda"></el-input>
          </el-form-item>
          <el-form-item label="random_state"
                        prop="random_state">
            <el-input v-model="form.random_state"></el-input>
          </el-form-item>
          <el-form-item label="n_jobs"
                        prop="n_jobs">
            <el-input v-model="form.n_jobs"></el-input>
          </el-form-item>
          <el-form-item label="importance_type"
                        prop="importance_type">
            <el-input v-model="form.importance_type"></el-input>
          </el-form-item>
        </div>
        <el-form-item v-show="form.param_type==='custom'"
                      :label="$t('ALG.customParameters')"
                      prop="custom_param">
          <el-input v-model="form.custom_param"
                    :rows="8"
                    @change="checkCustomParam"
                    :placeholder="$t('ALG.trainingParametersPro')"
                    type="textarea"></el-input>
        </el-form-item>
        <el-form-item :label="$t('ALG.early_stopping_rounds')"
                      prop="early_stopping_rounds">
          <el-input v-model="form.early_stopping_rounds"></el-input>
        </el-form-item>
        <el-form-item :label="$t('ALG.sample_weight_eval_set')"
                      prop="sample_weight_eval_set">
          <el-input v-model="form.sample_weight_eval_set"></el-input>
        </el-form-item>
        <el-form-item :label="$t('ALG.eval_metric')"
                      prop="eval_metric">
          <el-input v-model="form.eval_metric"></el-input>
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
import CustomModel from './CustomModel'
import nodeFormMixin from './nodeFormMixin'
const initForm = {
  API_type: '',
  param_type: 'list',
  boosting_type: '',
  num_leaves: '',
  max_depth: '',
  learning_rate: '',
  n_estimators: '',
  subsample_for_bin: '',
  objective: '',
  is_unbalance: '',
  min_split_gain: '',
  min_child_weight: '',
  subsample: '',
  subsample_freq: '',
  colsample_bytree: '',
  reg_alpha: '',
  reg_lambda: '',
  random_state: '',
  n_jobs: '',
  importance_type: '',
  early_stopping_rounds: '',
  sample_weight_eval_set: '',
  eval_metric: '',
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
        API_type: [
          { required: true, message: this.$t('ALG.APITypeReq') }
        ],
        num_leaves: [
          { type: 'string', pattern: new RegExp(/^(([1-9]{1}\d*)|(0{1}))(\.\d{1,7})?$/), message: this.$t('ALG.integerOrFloat') }
        ],
        max_depth: [
          { type: 'string', pattern: new RegExp(/^(([1-9-]{1}\d*)|(0{1}))$/), message: this.$t('ALG.integerFormat1') }
        ],
        learning_rate: [
          { type: 'string', pattern: new RegExp(/^(([1-9]{1}\d*)|(0{1}))(\.\d{1,7})$/), message: this.$t('ALG.floatingFormat') }
        ],
        n_estimators: [
          { type: 'string', pattern: new RegExp(/^(([1-9]{1}\d*)|(0{1}))$/), message: this.$t('ALG.integerFormat') }
        ],
        subsample_for_bin: [
          { type: 'string', pattern: new RegExp(/^(([1-9]{1}\d*)|(0{1}))$/), message: this.$t('ALG.integerFormat') }
        ],
        min_split_gain: [
          { type: 'string', pattern: new RegExp(/^(([1-9]{1}\d*)|(0{1}))(\.\d{1,7})?$/), message: this.$t('ALG.integerOrFloat') }
        ],
        min_child_weight: [
          { type: 'string', pattern: new RegExp(/^(([1-9]{1}\d*)|(0{1}))(\.\d{1,7})$/), message: this.$t('ALG.floatingFormat') }
        ],
        subsample: [
          { type: 'string', pattern: new RegExp(/^(([1-9]{1}\d*)|(0{1}))(\.\d{1,7})$/), message: this.$t('ALG.floatingFormat') }
        ],
        subsample_freq: [
          { type: 'string', pattern: new RegExp(/^(([1-9]{1}\d*)|(0{1}))$/), message: this.$t('ALG.integerFormat') }
        ],
        colsample_bytree: [
          { type: 'string', pattern: new RegExp(/^(([1-9]{1}\d*)|(0{1}))(\.\d{1,7})$/), message: this.$t('ALG.floatingFormat') }
        ],
        reg_alpha: [
          { type: 'string', pattern: new RegExp(/^(([1-9]{1}\d*)|(0{1}))(\.\d{1,7})$/), message: this.$t('ALG.floatingFormat') }
        ],
        reg_lambda: [
          { type: 'string', pattern: new RegExp(/^(([1-9]{1}\d*)|(0{1}))(\.\d{1,7})$/), message: this.$t('ALG.floatingFormat') }
        ],
        random_state: [
          { type: 'string', pattern: new RegExp(/^(([1-9]{1}\d*)|(0{1}))$/), message: this.$t('ALG.integerFormat') }
        ],
        n_jobs: [
          { type: 'string', pattern: new RegExp(/^(([1-9-]{1}\d*)|(0{1}))$/), message: this.$t('ALG.integerFormat1') }
        ],
        importance_type: [
          { type: 'string', pattern: new RegExp(/^[a-zA-Z]*$/), message: this.$t('ALG.onlyEnglish') }
        ],
        early_stopping_rounds: [
          { type: 'string', pattern: new RegExp(/^(([1-9]{1}\d*)|(0{1}))$/), message: this.$t('ALG.integerFormat') }
        ],
        sample_weight_eval_set: [
          { type: 'string', pattern: new RegExp(/^[^\u4e00-\u9fa5]*$/), message: this.$t('ALG.noCCharacters') }
        ],
        eval_metric: [
          { type: 'string', pattern: new RegExp(/^[^\u4e00-\u9fa5]*$/), message: this.$t('ALG.noCCharacters') }
        ],
        sample_weight: [
          { type: 'string', pattern: new RegExp(/^[^\u4e00-\u9fa5]*$/), message: this.$t('ALG.noCCharacters') }
        ],
        custom_param: [
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
    checkCustomParam () {
      try {
        let customParamStr = this.form.custom_param
        customParamStr = customParamStr.replace(/":False/g, '":false')
        customParamStr = customParamStr.replace(/":True/g, '":true')
        let customParam = JSON.parse(customParamStr)
        if (typeof (customParam) === 'string') {
          let customParam1 = JSON.parse(customParam)
          if (typeof (customParam) !== 'object' || Array.isArray(customParam1)) {
            this.$message.warning(this.$t('ALG.jsonPro'))
            this.form.custom_param = ''
          } else {
            this.form.custom_param = JSON.stringify(customParam)
          }
        } else {
          if (typeof (customParam) !== 'object' || Array.isArray(customParam)) {
            this.$message.warning(this.$t('ALG.jsonPro'))
            this.form.custom_param = ''
          } else {
            this.form.custom_param = JSON.stringify(customParam)
          }
        }
      } catch (err) {
        this.$message.warning(this.$t('ALG.jsonPro'))
        this.form.custom_param = ''
      }
    },
    algorithmDataFormatKey () {
      return {
        float: ['num_leaves', 'learning_rate', 'min_split_gain', 'min_child_weight', 'subsample', 'colsample_bytree', 'reg_alpha', 'reg_lambda'],
        int: ['max_depth', 'n_estimators', 'subsample_for_bin', 'subsample_freq', 'random_state', 'n_jobs', 'early_stopping_rounds']
      }
    },
    save (evt, temSave = false) {
      this.algorithmSubmit(initForm, 'LightGBM', temSave)
    }
  }
}
</script>
