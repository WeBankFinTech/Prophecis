<template>
  <el-dialog :title="$t('flow.alertSetting')"
             :visible.sync="alarmVisible"
             custom-class="dialog-style"
             :close-on-click-modal="false"
             @close="closeAlarm"
             width="1100px">
    <alarm-setting ref="alarmSetting"
                   :alarmData="alarmData"></alarm-setting>
    <span slot="footer"
          class="dialog-footer">
      <el-button type="primary"
                 @click="save">
        {{ $t('common.save') }}
      </el-button>
      <el-button @click="alarmVisible=false">
        {{ $t('common.cancel') }}
      </el-button>
    </span>
  </el-dialog>
</template>
<script>
import AlarmSetting from '../../../../components/AlarmSetting'
import handleDIDetailMixin from '../../../../util/handleDIDetailMixin'
export default {
  components: {
    AlarmSetting
  },
  mixins: [handleDIDetailMixin],
  props: {
    nodeData: {
      type: Object,
      default () {
        return {}
      }
    }
  },
  data () {
    return {
      alarmVisible: false,
      currentNode: {},
      alarmData: [{
        alarmType: []
      }]
    }
  },
  mounted () {
    this.handleOriginAlarmData()
  },
  methods: {
    handleOriginAlarmData () {
      this.currentNode = this.nodeData
      const jobAlert = this.currentNode.jobContent && this.currentNode.jobContent.ManiFest && this.currentNode.jobContent.ManiFest.job_alert
      if (jobAlert) {
        const alarmData = this.processAlertInfo(jobAlert, false)
        this.alarmData = alarmData
      }
    },
    save () {
      this.$refs.alarmSetting.$refs.formValidate.validate((valid) => {
        if (valid) {
          this.$refs.alarmSetting.handleAlarmParam(this.currentNode.jobContent.ManiFest)
          this.$emit('saveNode', this.currentNode)
          this.alarmVisible = false
        }
      })
    },
    closeAlarm () {
      this.$emit('closeAlarm')
    }
  }
}
</script>
