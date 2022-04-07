<template>
  <el-dialog :title="$t('ns.viewDetails')"
             :visible.sync="dialogVisible"
             custom-class="dialog-style"
             width="1100px">
    <el-form ref="formValidate"
             :model="form"
             :disabled="disable"
             class="add-form"
             label-width="135px">
      <div class="subtitle">
        {{$t('DI.basicSettings')}}
      </div>
      <el-form-item :label="$t('serviceList.modelName')">
        <el-input v-model="form.service_name" />
      </el-form-item>
      <el-form-item :label="$t('serviceList.modelDescription')">
        <el-input v-model="form.remark"
                  type="textarea" />
      </el-form-item>
      <div class="subtitle">
        {{$t('DI.computingResource')}}
      </div>
      <el-row>
        <el-col :span="12">
          <el-form-item :label="$t('ns.nameSpace')">
            <el-select v-model="form.namespace">
              <el-option v-for="item in spaceOptionList"
                         :label="item"
                         :key="item"
                         :value="item">
              </el-option>
            </el-select>
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item label="CPU">
            <el-input v-model="form.cpu">
              <span slot="append">Core</span>
            </el-input>
          </el-form-item>
        </el-col>
      </el-row>
      <el-row>
        <el-col :span="12">
          <el-form-item :label="$t('AIDE.memory1')">
            <el-input v-model="form.memory">
              <span slot="append">Gi</span>
            </el-input>
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item label="GPU">
            <el-input v-model="form.gpu">
              <span slot="append">{{$t('AIDE.block')}}</span>
            </el-input>
          </el-form-item>
        </el-col>
      </el-row>
      <div class="subtitle">
        {{$t('serviceList.deploymentInfo')}}
      </div>
      <model-setting v-if="dialogVisible"
                     :form-data="form"
                     :view-detail="true"></model-setting>
    </el-form>
    <!-- <div class="subtitle">
        {{$t('serviceList.logDirectorySettings')}}
      </div>
      <el-form-item :label="$t('serviceList.logDirectory')"
                    prop="logDirectory">
        <el-input v-model="form.log_path">
        </el-input>
      </el-form-item> -->
  </el-dialog>
</template>
<script>
import ModelSetting from './ModelSetting'
export default {
  components: { ModelSetting },
  props: {
    form: {
      type: Object,
      default () {
        return {
          service_name: '',
          version: '',
          remark: '',
          namespace: '',
          cpu: '',
          memory: '',
          gpu: '',
          model_type: '',
          type: '',
          parameters: []
          // filepath: '',
          // log_path: ''
        }
      }
    },
    groupList: {
      type: Array,
      default () {
        return []
      }
    },
    disable: {
      type: Boolean,
      default: true
    },
    spaceOptionList: {
      type: Array,
      default () {
        return []
      }
    }
  },
  data () {
    return {
      dialogVisible: false
    }
  }
}
</script>
<style lang="scss" scoped>
.padding-label {
  padding-left: 135px;
  margin-bottom: 18px;
}
</style>
