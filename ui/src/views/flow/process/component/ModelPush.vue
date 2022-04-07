<template>
  <div>
    <el-form ref="formValidate"
             :model="form"
             :disabled="readonly"
             :rules="ruleValidate"
             label-position="top"
             class="node-parameter-bar">
      <div v-show="activeName==='0'"
           :key="'step0'+activeName">
        <model-param ref="modelParam"
                     :readonly="readonly"></model-param>
      </div>
      <div v-show="activeName==='1'"
           :key="'step1'+activeName">
        <el-form-item :label="$t('flow.FPSUsers')">
          <el-select v-model="form.fpsUser">
            <el-option label="hduser"
                       value="hduser"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('flow.RMBDerviceId')">
          <el-select v-model="form.rmbUser">
            <el-option label="hduser"
                       value="hduser"></el-option>
          </el-select>
        </el-form-item>
      </div>
    </el-form>
    <div class="save-button"
         v-if="!readonly">
      <el-button>{{ $t('common.save') }}</el-button>
    </div>
  </div>
</template>
<script>
import modelMixin from '../../../../util/modelMixin'
import ModelParam from './ModelParam'
import nodeFormMixin from './nodeFormMixin'
const initForm = {
  path: '',
  filePath: '',
  fpsUser: '',
  rmbUser: ''
}
export default {
  mixins: [modelMixin, nodeFormMixin],
  components: { ModelParam },
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
      localPathList: [],
      form: { ...initForm }
    }
  },
  created () {
  },
  computed: {
    ruleValidate: function () {
      return {
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
  }
}
</script>
