<template>
  <el-form ref="formValidate"
           :model="form"
           :disabled="readonly"
           :rules="ruleValidate"
           label-position="top"
           class="node-parameter-bar">
    <el-form-item :label="$t('flow.nodeName')"
                  prop="name">
      <el-input v-model="form.name"
                :placeholder="$t('flow.nodeInfoPro')"></el-input>
    </el-form-item>
    <el-form-item :label="$t('flow.dataFormat')"
                  prop="data_type">
      <el-select v-model="form.data_type">
        <el-option label="Pickle"
                   value="pickle"></el-option>
        <el-option v-if="nodeType==='LightGBM'"
                   label="DataSet"
                   value="DataSet"></el-option>
      </el-select>
    </el-form-item>
    <div class="subtitle">
      {{$t('flow.directorySet')}}
    </div>
    <el-form-item :label="$t('flow.storageRootDirectory')"
                  prop="path">
      <el-select v-model="form.path"
                 filterable
                 :placeholder="$t('DI.trainingDataStorePro')">
        <el-option v-for="(item,index) in localPathList"
                   :label="item"
                   :key="index"
                   :value="item">
        </el-option>
      </el-select>
    </el-form-item>
    <div class="subtitle">
      {{$t('flow.trainingSet')}}
    </div>
    <el-form-item :label="$t('flow.data')"
                  prop="training_data_path">
      <el-input v-model="form.training_data_path"
                :placeholder="$t('flow.dataDirectoryPro')" />
    </el-form-item>
    <el-form-item :label="$t('flow.label')"
                  prop="training_label_path">
      <el-input v-model="form.training_label_path"
                :placeholder="$t('flow.labelDirectoryPro')" />
    </el-form-item>
    <div class="subtitle">
      {{$t('flow.testSet')}}
    </div>
    <el-form-item :label="$t('flow.data')"
                  prop="testing_data_path">
      <el-input v-model="form.testing_data_path"
                :placeholder="$t('flow.dataDirectoryPro')" />
    </el-form-item>
    <el-form-item :label="$t('flow.label')"
                  prop="testing_label_path">
      <el-input v-model="form.testing_label_path"
                :placeholder="$t('flow.labelDirectoryPro')" />
    </el-form-item>
    <div class="subtitle">
      {{$t('flow.validationSet')}}
    </div>
    <el-form-item :label="$t('flow.data')"
                  prop="validation_data_path">
      <el-input v-model="form.validation_data_path"
                :placeholder="$t('flow.dataDirectoryPro')" />
    </el-form-item>
    <el-form-item :label="$t('flow.label')"
                  prop="validation_label_path">
      <el-input v-model="form.validation_label_path"
                :placeholder="$t('flow.labelDirectoryPro')" />
    </el-form-item>
  </el-form>
</template>
<script>
import modelMixin from '../../../../util/modelMixin'
import nodeFormMixin from './nodeFormMixin'
const initForm = {
  name: '',
  data_type: 'pickle',
  dataFormat: '',
  path: '',
  training_data_path: '',
  training_label_path: '',
  testing_data_path: '',
  testing_label_path: '',
  validation_data_path: '',
  validation_label_path: ''
}
export default {
  mixins: [modelMixin, nodeFormMixin],
  data () {
    return {
      localPathList: [],
      groupList: [],
      modelOptionList: [],
      form: { ...initForm }
    }
  },
  props: {
    readonly: {
      type: Boolean,
      default: false
    },
    nodeType: {
      type: String,
      default: ''
    }
  },
  computed: {
    ruleValidate: function () {
      let validateRule = {
        name: [
          { required: true, message: this.$t('flow.nodeNoEmpty') },
          { type: 'string', pattern: new RegExp(/^[0-9a-zA-Z-_]*$/), message: this.$t('flow.nodeNameTip') }
        ],
        data_type: [
          { required: true, message: this.$t('flow.dataFormatReq') }
        ],
        path: [
          { required: true, message: this.$t('DI.trainingDataStoreReq') },
          { pattern: new RegExp(/^\/[a-zA-Z][0-9a-zA-Z-_/.]*$/), message: this.$t('DI.trainingDataStoreFormat') }
        ],
        training_data_path: [
          { required: true, message: this.$t('DI.dataDirectoryReq') },
          { pattern: new RegExp(/^[a-zA-Z][0-9a-zA-Z-_/.]*$/), message: this.$t('DI.dataDirectoryFormat') }
        ],
        testing_data_path: [
          { pattern: new RegExp(/^[a-zA-Z][0-9a-zA-Z-_/.]*$/), message: this.$t('DI.dataDirectoryFormat') }
        ],
        testing_label_path: [
          { pattern: new RegExp(/^[a-zA-Z][0-9a-zA-Z-_/.]*$/), message: this.$t('flow.labelSubdirectoryFormat') }
        ],
        validation_data_path: [
          { pattern: new RegExp(/^[a-zA-Z][0-9a-zA-Z-_/.]*$/), message: this.$t('DI.dataDirectoryFormat') }
        ],
        validation_label_path: [
          { pattern: new RegExp(/^[a-zA-Z][0-9a-zA-Z-_/.]*$/), message: this.$t('flow.labelSubdirectoryFormat') }
        ]
      }
      if (this.nodeType === 'LightGBM') {
        validateRule.training_label_path = [
          { pattern: new RegExp(/^[a-zA-Z][0-9a-zA-Z-_/.]*$/), message: this.$t('flow.labelSubdirectoryFormat') }
        ]
      } else {
        validateRule.training_label_path = [
          { required: true, message: this.$t('flow.labelSubdirectoryReq') },
          { pattern: new RegExp(/^[a-zA-Z][0-9a-zA-Z-_/.]*$/), message: this.$t('flow.labelSubdirectoryFormat') }
        ]
      }
      return validateRule
    }
  },
  created () {
    const header = this.handlerFromDssHeader('get')
    this.getDIStoragePath(header)
    this.getGroupList(header)
  },
  methods: {
    initFormData (currentNode) {
      this.form = { ...initForm }
      this.form.name = currentNode.title
      setTimeout(() => {
        this.$refs.formValidate.clearValidate()
      }, 0)
    },
    save () {
      return new Promise(resolve => {
        this.$refs.formValidate.validate(valid => {
          if (valid) {
            let dataPath = {}
            let pathArr = ['data_type', 'training_data_path', 'testing_data_path', 'validation_data_path', 'training_label_path', 'testing_label_path', 'validation_label_path']
            for (let key of pathArr) {
              dataPath[key] = this.form[key]
            }
            let dataSetting = {
              name: this.form.name,
              data_stores: [
                {
                  id: 'hostmount',
                  type: 'mount_volume',
                  training_data: {
                    container: 'data'
                  },
                  training_results: {
                    container: 'result'
                  },
                  connection: {
                    type: 'host_mount',
                    name: 'host-mount',
                    path: this.form.path
                  },
                  training_workspace: {
                    container: 'data'
                  }
                }
              ],
              data_set: dataPath
            }
            resolve(dataSetting)
          } else {
            resolve(false)
          }
        })
      })
    }
  }
}
</script>
