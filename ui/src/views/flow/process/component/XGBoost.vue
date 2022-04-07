<template>
  <div>
    <div v-show="activeName==='0'"
         key="step0">
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
          <el-form-item :label="$t('ALG.n_estimators')"
                        prop="n_estimators">
            <el-input v-model="form.n_estimators"></el-input>
          </el-form-item>
          <el-form-item :label="$t('ALG.use_label_encoder')">
            <el-select v-model="form.use_label_encoder"
                       clearable>
              <el-option :label="$t('common.yes')"
                         :value="true">
              </el-option>
              <el-option :label="$t('common.no')"
                         :value="false">
              </el-option>
            </el-select>
          </el-form-item>
          <el-form-item :label="$t('ALG.max_depth')"
                        prop="max_depth">
            <el-input v-model="form.max_depth"></el-input>
          </el-form-item>
          <el-form-item :label="$t('ALG.learning_rate')"
                        prop="learning_rate">
            <el-input v-model="form.learning_rate"></el-input>
          </el-form-item>
          <el-form-item :label="$t('ALG.verbosity')"
                        prop="verbosity">
            <el-input v-model="form.verbosity"></el-input>
          </el-form-item>
          <el-form-item :label="$t('ALG.objective')">
            <el-select v-model="form.objective"
                       clearable>
              <el-option label="reg:linear"
                         value="reg:linear">
              </el-option>
              <el-option label="reg:logistic"
                         value="reg:logistic">
              </el-option>
              <el-option label="binary:logistic"
                         value="binary:logistic">
              </el-option>
              <el-option label="binary:logitraw"
                         value="binary:logitraw">
              </el-option>
              <el-option label="multi:softmax"
                         value="multi:softmax">
              </el-option>
              <el-option label="multi:softprob"
                         value="multi:softprob">
              </el-option>
            </el-select>
          </el-form-item>
          <el-form-item :label="$t('ALG.booster')">
            <el-select v-model="form.booster"
                       clearable>
              <el-option label="gbtree"
                         value="gbtree">
              </el-option>
              <el-option label="gblinear"
                         value="gblinear">
              </el-option>
            </el-select>
          </el-form-item>
          <el-form-item :label="$t('ALG.tree_method')">
            <el-select v-model="form.tree_method"
                       clearable>
              <el-option label="auto"
                         value="auto">
              </el-option>
              <el-option label="exact"
                         value="exact">
              </el-option>
              <el-option label="approx"
                         value="approx">
              </el-option>
            </el-select>
          </el-form-item>
          <el-form-item :label="$t('ALG.n_jobs')"
                        prop="n_jobs">
            <el-input v-model="form.n_jobs"></el-input>
          </el-form-item>
          <el-form-item :label="$t('ALG.min_child_weight')"
                        prop="min_child_weight">
            <el-input v-model="form.min_child_weight"></el-input>
          </el-form-item>
          <el-form-item :label="$t('ALG.max_delta_step')"
                        prop="max_delta_step">
            <el-input v-model="form.max_delta_step"></el-input>
          </el-form-item>
          <el-form-item :label="$t('ALG.subsample')"
                        prop="subsample">
            <el-input v-model="form.subsample"></el-input>
          </el-form-item>
          <el-form-item :label="$t('ALG.colsample_bytree')"
                        prop="colsample_bytree">
            <el-input v-model="form.colsample_bytree">
            </el-input>
          </el-form-item>
          <el-form-item :label="$t('ALG.reg_alpha')"
                        prop="reg_alpha">
            <el-input v-model="form.reg_alpha">
            </el-input>
          </el-form-item>
          <el-form-item :label="$t('ALG.reg_lambda')"
                        prop="reg_lambda">
            <el-input v-model="form.reg_lambda">
            </el-input>
          </el-form-item>
          <el-form-item :label="$t('ALG.scale_pos_weight')"
                        prop="scale_pos_weight">
            <el-input v-model="form.scale_pos_weight">
            </el-input>
          </el-form-item>
          <el-form-item :label="$t('ALG.base_score')"
                        prop="base_score">
            <el-input v-model="form.base_score">
            </el-input>
          </el-form-item>
          <el-form-item :label="$t('ALG.random_state')"
                        prop="random_state">
            <el-input v-model="form.random_state">
            </el-input>
          </el-form-item>
          <el-form-item :label="$t('ALG.missing')"
                        prop="missing">
            <el-input v-model="form.missing">
            </el-input>
          </el-form-item>
          <el-form-item :label="$t('ALG.importance_type')">
            <el-select v-model="form.importance_type"
                       clearable>
              <el-option label="weight"
                         value="weight"></el-option>
              <el-option label="gain"
                         value="gain"></el-option>
              <el-option label="cover"
                         value="cover"></el-option>
              <el-option label="total_gain"
                         value="total_gain"></el-option>
              <el-option label="total_cover"
                         value="total_cover"></el-option>
            </el-select>
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
  n_estimators: '',
  use_label_encoder: '',
  max_depth: '',
  learning_rate: '',
  verbosity: '',
  objective: '',
  booster: '',
  tree_method: '',
  n_jobs: '',
  min_child_weight: '',
  max_delta_step: '',
  subsample: '',
  colsample_bytree: '',
  reg_alpha: '',
  reg_lambda: '',
  scale_pos_weight: '',
  base_score: '',
  random_state: '',
  missing: '',
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
        n_estimators: [
          { type: 'string', pattern: new RegExp(/^(([1-9]{1}\d*)|(0{1}))$/), message: this.$t('ALG.integerFormat') }
        ],
        max_depth: [
          { type: 'string', pattern: new RegExp(/^(([1-9]{1}\d*)|(0{1}))$/), message: this.$t('ALG.integerFormat') }
        ],
        learning_rate: [
          { type: 'string', pattern: new RegExp(/^(([1-9]{1}\d*)|(0{1}))(\.\d{1,7})$/), message: this.$t('ALG.floatingFormat') }
        ],
        verbosity: [
          { type: 'string', pattern: new RegExp(/^(([1-9]{1}\d*)|(0{1}))$/), message: this.$t('ALG.integerFormat') }
        ],
        n_jobs: [
          { type: 'string', pattern: new RegExp(/^(([1-9-]{1}\d*)|(0{1}))$/), message: this.$t('ALG.integerFormat1') }
        ],
        min_child_weight: [
          { type: 'string', pattern: new RegExp(/^(([1-9]{1}\d*)|(0{1}))(\.\d{1,7})?$/), message: this.$t('ALG.integerOrFloat') }
        ],
        max_delta_step: [
          { type: 'string', pattern: new RegExp(/^(([1-9]{1}\d*)|(0{1}))(\.\d{1,7})?$/), message: this.$t('ALG.integerOrFloat') }
        ],
        subsample: [
          { type: 'string', pattern: new RegExp(/^(([1-9]{1}\d*)|(0{1}))(\.\d{1,7})?$/), message: this.$t('ALG.integerOrFloat') }
        ],
        colsample_bytree: [
          { type: 'string', pattern: new RegExp(/^(([1-9]{1}\d*)|(0{1}))(\.\d{1,7})?$/), message: this.$t('ALG.integerOrFloat') }
        ],
        reg_alpha: [
          { type: 'string', pattern: new RegExp(/^(([1-9]{1}\d*)|(0{1}))(\.\d{1,7})?$/), message: this.$t('ALG.integerOrFloat') }
        ],
        reg_lambda: [
          { type: 'string', pattern: new RegExp(/^(([1-9]{1}\d*)|(0{1}))(\.\d{1,7})?$/), message: this.$t('ALG.integerOrFloat') }
        ],
        scale_pos_weight: [
          { type: 'string', pattern: new RegExp(/^(([1-9]{1}\d*)|(0{1}))(\.\d{1,7})?$/), message: this.$t('ALG.integerOrFloat') }
        ],
        base_score: [
          { type: 'string', pattern: new RegExp(/^(([1-9]{1}\d*)|(0{1}))(\.\d{1,7})$/), message: this.$t('ALG.floatingFormat') }
        ],
        random_state: [
          { type: 'string', pattern: new RegExp(/^(([1-9]{1}\d*)|(0{1}))$/), message: this.$t('ALG.integerFormat') }
        ],
        missing: [
          { type: 'string', pattern: new RegExp(/^(([1-9]{1}\d*)|(0{1}))(\.\d{1,7})$/), message: this.$t('ALG.floatingFormat') }
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
        float: ['learning_rate', 'min_child_weight', 'max_delta_step', 'subsample', 'colsample_bytree', 'reg_alpha', 'reg_lambda', 'missing', 'scale_pos_weight', 'base_score'],
        int: ['n_estimators', 'max_depth', 'verbosity', 'n_jobs', 'random_state', 'early_stopping_rounds']
      }
    },
    save (evt, temSave = false) {
      this.algorithmSubmit(initForm, 'XGBoost', temSave)
    }
  }
}
</script>
