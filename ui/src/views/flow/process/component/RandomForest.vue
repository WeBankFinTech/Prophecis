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
        <el-form-item :label="$t('ALG.n_estimators')"
                      prop="n_estimators">
          <el-input v-model="form.n_estimators"></el-input>
        </el-form-item>
        <el-form-item :label="$t('ALG.criterion')">
          <el-select v-model="form.criterion"
                     clearable>
            <el-option label="gini"
                       value="gini">
            </el-option>
            <el-option label="entropy"
                       value="entropy">
            </el-option>
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('ALG.max_depth')"
                      prop="max_depth">
          <el-input v-model="form.max_depth"></el-input>
        </el-form-item>
        <el-form-item :label="$t('ALG.min_samples_split')"
                      prop="min_samples_split">
          <el-input v-model="form.min_samples_split"></el-input>
        </el-form-item>
        <el-form-item :label="$t('ALG.min_samples_leaf')"
                      prop="min_samples_leaf">
          <el-input v-model="form.min_samples_leaf"></el-input>
        </el-form-item>
        <el-form-item :label="$t('ALG.min_weight_fraction_leaf')"
                      prop="min_weight_fraction_leaf">
          <el-input v-model="form.min_weight_fraction_leaf"></el-input>
        </el-form-item>
        <el-form-item :label="$t('ALG.max_features')"
                      prop="max_features">
          <el-input v-model="form.max_features"></el-input>
        </el-form-item>
        <el-form-item :label="$t('ALG.max_leaf_nodes')"
                      prop="max_leaf_nodes">
          <el-input v-model="form.max_leaf_nodes"></el-input>
        </el-form-item>
        <el-form-item :label="$t('ALG.min_impurity_decrease')"
                      prop="min_impurity_decrease">
          <el-input v-model="form.min_impurity_decrease"></el-input>
        </el-form-item>
        <el-form-item :label="$t('ALG.min_impurity_split')"
                      prop="min_impurity_split">
          <el-input v-model="form.min_impurity_split"></el-input>
        </el-form-item>
        <el-form-item :label="$t('ALG.bootstrap')">
          <el-select v-model="form.bootstrap"
                     clearable>
            <el-option :label="$t('common.yes')"
                       :value="true"></el-option>
            <el-option :label="$t('common.no')"
                       :value="false"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('ALG.oob_score')"
                      v-if="form.bootstrap===true">
          <el-select v-model="form.oob_score"
                     clearable>
            <el-option :label="$t('common.yes')"
                       :value="true"></el-option>
            <el-option :label="$t('common.no')"
                       :value="false"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('ALG.n_jobs')"
                      prop="n_jobs">
          <el-input v-model="form.n_jobs"></el-input>
        </el-form-item>
        <el-form-item :label="$t('ALG.random_state')"
                      prop="random_state">
          <el-input v-model="form.random_state"></el-input>
        </el-form-item>
        <el-form-item :label="$t('ALG.verbose')"
                      prop="verbose">
          <el-input v-model="form.verbose"></el-input>
        </el-form-item>
        <el-form-item :label="$t('ALG.warm_start')">
          <el-select v-model="form.warm_start"
                     clearable>
            <el-option :label="$t('common.yes')"
                       :value="true"></el-option>
            <el-option :label="$t('common.no')"
                       :value="false"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('ALG.class_weight1')">
          <el-select v-model="form.class_weight"
                     clearable>
            <el-option label="dict"
                       value="dict"></el-option>
            <el-option label="balanced"
                       value="balanced"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('ALG.ccp_alpha')"
                      prop="ccp_alpha">
          <el-input v-model="form.ccp_alpha"></el-input>
        </el-form-item>
        <el-form-item :label="$t('ALG.max_samples')"
                      prop="max_samples">
          <el-input v-model="form.max_samples"></el-input>
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
  n_estimators: '',
  criterion: '',
  max_depth: '',
  min_samples_split: '',
  min_samples_leaf: '',
  min_weight_fraction_leaf: '',
  max_features: '',
  max_leaf_nodes: '',
  min_impurity_decrease: '',
  min_impurity_split: '',
  bootstrap: '',
  oob_score: '',
  n_jobs: '',
  random_state: '',
  verbose: '',
  warm_start: '',
  class_weight: '',
  ccp_alpha: '',
  max_samples: '',
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
        n_estimators: [
          { type: 'string', pattern: new RegExp(/^(([1-9]{1}\d*)|(0{1}))$/), message: this.$t('ALG.integerFormat') }
        ],
        max_depth: [
          { type: 'string', pattern: new RegExp(/^(([1-9]{1}\d*)|(0{1}))$/), message: this.$t('ALG.integerFormat') }
        ],
        min_samples_split: [
          { type: 'string', pattern: new RegExp(/^(([1-9]{1}\d*)|(0{1}))(\.\d{0,7})?$/), message: this.$t('ALG.integerOrFloat') }
        ],
        min_samples_leaf: [
          { type: 'string', pattern: new RegExp(/^(([1-9]{1}\d*)|(0{1}))(\.\d{0,7})?$/), message: this.$t('ALG.integerOrFloat') }
        ],
        min_weight_fraction_leaf: [
          { type: 'string', pattern: new RegExp(/^(([1-9]{1}\d*)|(0{1}))(\.\d{0,7})$/), message: this.$t('ALG.floatingFormat') }
        ],
        max_features: [
          { type: 'string', pattern: new RegExp(/^[a-zA-Z0-9.]*$/), message: this.$t('ALG.intFloatString') }
        ],
        max_leaf_nodes: [
          { type: 'string', pattern: new RegExp(/^(([1-9]{1}\d*)|(0{1}))$/), message: this.$t('ALG.integerFormat') }
        ],
        min_impurity_decrease: [
          { type: 'string', pattern: new RegExp(/^(([1-9]{1}\d*)|(0{1}))(\.\d{0,7})$/), message: this.$t('ALG.floatingFormat') }
        ],
        min_impurity_split: [
          { type: 'string', pattern: new RegExp(/^(([1-9]{1}\d*)|(0{1}))(\.\d{0,7})$/), message: this.$t('ALG.floatingFormat') }
        ],
        n_jobs: [
          { type: 'string', pattern: new RegExp(/^(([1-9-]{1}\d*)|(0{1}))$/), message: this.$t('ALG.integerFormat1') }
        ],
        random_state: [
          { type: 'string', pattern: new RegExp(/^(([1-9]{1}\d*)|(0{1}))$/), message: this.$t('ALG.integerFormat') }
        ],
        verbose: [
          { type: 'string', pattern: new RegExp(/^(([1-9]{1}\d*)|(0{1}))$/), message: this.$t('ALG.integerFormat') }
        ],
        ccp_alpha: [
          { type: 'string', pattern: new RegExp(/^(([1-9]{1}\d*)|(0{1}))(\.\d{0,7})$/), message: this.$t('ALG.floatingFormat') }
        ],
        max_samples: [
          { type: 'string', pattern: new RegExp(/^(([1-9]{1}\d*)|(0{1}))(\.\d{0,7})?$/), message: this.$t('ALG.integerOrFloat') }
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
        float: ['min_samples_split', 'min_samples_leaf', 'min_weight_fraction_leaf', 'min_impurity_decrease', 'min_impurity_split', 'max_samples', 'ccp_alpha'],
        int: ['n_estimators', 'max_depth', 'max_leaf_nodes', 'n_jobs', 'random_state', 'verbose']
      }
    },
    save (evt, temSave = false) {
      this.algorithmSubmit(initForm, 'RandomForest', temSave)
    }
  }
}
</script>
