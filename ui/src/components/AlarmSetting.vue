<template>
  <el-form ref="formValidate"
           :model="form"
           :rules="ruleValidate"
           class="alarm-form"
           label-width="150px">
    <div class="subtitle alarm-set">
      {{ $t('DI.alarmSet') }}
      <i class="el-icon-circle-plus-outline add"
         @click="addAlarm" />
      <i class="el-icon-remove-outline delete"
         @click="deleteAlarm" />
    </div>
    <div v-for="(item,index) in form.setAlarm"
         :key="`setAlarm${index}`">
      <div class="subtitle "
           v-if="index>0">
        {{ $t('DI.alarmSet') }}
        <!-- <i v-if="index===0"
           class="el-icon-circle-plus-outline add"
           @click="addAlarm" /> -->
        <i class="el-icon-remove-outline delete"
           @click="deleteAlarm(index)" />
      </div>
      <el-row>
        <el-form-item :label="$t('DI.alarmTypes')"
                      :prop="`setAlarm[${index}].alarmType`">
          <el-select v-model="item.alarmType"
                     multiple
                     :placeholder="$t('DI.alarmTypesPro')"
                     @change="changeAlarmType(index)">
            <el-option :label="$t('DI.eventAlarm')"
                       value="1">
            </el-option>
            <el-option :label="$t('DI.designatedTimeout')"
                       value="2">
            </el-option>
            <el-option :label="$t('DI.hoursOvertime')"
                       value="3">
            </el-option>
          </el-select>
        </el-form-item>
      </el-row>
      <div v-if="item.alarmType.indexOf('1')>-1"
           key="eventChecker">
        <div class="subtitle">
          {{ $t('DI.eventAlarmSettings') }}
        </div>
        <el-row>
          <el-col :span="12">
            <el-form-item :label="$t('DI.listeningEnent')"
                          :prop="`setAlarm[${index}].event.event_checker`">
              <el-select v-model="item.event.event_checker"
                         :placeholder="$t('DI.listeningEnentPro')">
                <el-option :label="$t('DI.executeSuccess')"
                           value="COMPLETED">
                </el-option>
                <el-option :label="$t('DI.executeFail')"
                           value="FAILED">
                </el-option>
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item :label="$t('DI.alarmLevel')"
                          :prop="`setAlarm[${index}].event.alert_level`">
              <el-select v-model="item.event.alert_level"
                         :placeholder="$t('DI.alarmLevelPro')">
                <el-option v-for="(val,key) in alarmLevel"
                           :key="key"
                           :label="key"
                           :value="val"></el-option>
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row>
          <el-form-item :label="$t('DI.receiver')"
                        :prop="`setAlarm[${index}].event.receiver`">
            <el-input v-model="item.event.receiver"
                      :placeholder="$t('DI.receiverPro')" />
          </el-form-item>
        </el-row>
      </div>
      <div v-if="item.alarmType.indexOf('2')>-1"
           key="designatedTimeout">
        <div class="subtitle">
          {{ $t('DI.designatedTimeoutSetting') }}
        </div>
        <el-row>
          <el-col :span="12">
            <el-form-item :label="$t('DI.designatedTimeout')"
                          :prop="`setAlarm[${index}].fixTime.deadlineChecker`">
              <el-date-picker v-model="item.fixTime.deadlineChecker"
                              type="datetime"
                              :placeholder="$t('DI.designatedTimeoutPro')" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item :label="$t('DI.alarmFrequency')"
                          :prop="`setAlarm[${index}].fixTime.interval`">
              <el-input v-model="item.fixTime.interval"
                        :placeholder="$t('DI.alarmFrequencyPro')">
                <span slot="append">min</span>
              </el-input>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row>
          <el-col :span="12">
            <el-form-item :label="$t('DI.alarmLevel')"
                          :prop="`setAlarm[${index}].fixTime.alert_level`">
              <el-select v-model="item.fixTime.alert_level"
                         :placeholder="$t('DI.alarmLevelPro')">
                <el-option v-for="(val,key) in alarmLevel"
                           :key="key"
                           :label="key"
                           :value="val"></el-option>
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item :label="$t('DI.receiver')"
                          :prop="`setAlarm[${index}].fixTime.receiver`">
              <el-input v-model="item.fixTime.receiver"
                        :placeholder="$t('DI.receiverPro')" />
            </el-form-item>
          </el-col>
        </el-row>
      </div>
      <div v-if="item.alarmType.indexOf('3')>-1"
           key="overtimeChecker">
        <div class="subtitle">
          {{ $t('DI.timeoutTimeSetting') }}
        </div>
        <el-row>
          <el-col :span="12">
            <el-form-item :label="$t('DI.hoursOvertime')"
                          :prop="`setAlarm[${index}].overTime.overtime_checker`">
              <el-input v-model="item.overTime.overtime_checker"
                        :placeholder="$t('DI.timeoutTimePro')">
                <span slot="append">h</span>
              </el-input>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item :label="$t('DI.alarmFrequency')"
                          :prop="`setAlarm[${index}].overTime.interval`">
              <el-input v-model="item.overTime.interval"
                        :placeholder="$t('DI.alarmFrequencyPro')">
                <span slot="append">min</span>
              </el-input>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row>
          <el-col :span="12">
            <el-form-item :label="$t('DI.alarmLevel')"
                          :prop="`setAlarm[${index}].overTime.alert_level`">
              <el-select v-model="item.overTime.alert_level"
                         :placeholder="$t('DI.alarmLevelPro')">
                <el-option v-for="(val,key) in alarmLevel"
                           :key="key"
                           :label="key"
                           :value="val"></el-option>
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item :label="$t('DI.receiver')"
                          :prop="`setAlarm[${index}].overTime.receiver`">
              <el-input v-model="item.overTime.receiver"
                        :placeholder="$t('DI.receiverPro')" />
            </el-form-item>
          </el-col>
        </el-row>
      </div>
    </div>
  </el-form>
</template>
<script>
import util from '../util/common'
export default {
  props: {
    alarmData: {
      type: Array,
      default: () => {
        return [{
          alarmType: []
        }]
      }
    }
  },
  data () {
    return {
      alarmLevel: {
        Info: 'info',
        Warning: 'warning',
        Minor: 'minor',
        Major: 'major',
        Critical: 'critical'
      },
      form: {
        setAlarm: [{
          alarmType: []
        }]
      },
      userId: localStorage.getItem('userId')
    }
  },
  computed: {
    setAlarmLength: function () {
      return this.form.setAlarm.length
    },
    ruleValidate: function () {
      // 切换语言表单报错重置表单
      this.resetFormFields()
      return {
        'setAlarm[0].alarmType': [
          { required: true, message: this.$t('DI.alarmTypesReq') }
        ],
        'setAlarm[0].event.event_checker': [
          { required: true, message: this.$t('DI.listeningEnentReq') }
        ],
        'setAlarm[0].event.alert_level': [
          { required: true, message: this.$t('DI.alarmLevelReq') }
        ],
        'setAlarm[0].event.receiver': [
          { required: true, message: this.$t('DI.receiverReq') }
        ],
        'setAlarm[0].fixTime.deadlineChecker': [
          { required: true, message: this.$t('DI.designatedTimeoutReq') }
        ],
        'setAlarm[0].fixTime.alert_level': [
          { required: true, message: this.$t('DI.alarmLevelReq') }
        ],
        'setAlarm[0].fixTime.receiver': [
          { required: true, message: this.$t('DI.receiverReq') }
        ],
        'setAlarm[0].fixTime.interval': [
          { required: true, message: this.$t('DI.alarmFrequencyReq') },
          { pattern: new RegExp(/^[1-9]\d*$/), message: this.$t('DI.alarmFrequencyFormat') }
        ],
        'setAlarm[0].overTime.overtime_checker': [
          { required: true, message: this.$t('DI.timeoutTimeReq') },
          { pattern: new RegExp(/^[1-9]\d*$/), message: this.$t('DI.timeoutTimeFormat') }
        ],
        'setAlarm[0].overTime.alert_level': [
          { required: true, message: this.$t('DI.alarmLevelReq') }
        ],
        'setAlarm[0].overTime.receiver': [
          { required: true, message: this.$t('DI.receiverReq') }
        ],
        'setAlarm[0].overTime.interval': [
          { required: true, message: this.$t('DI.alarmFrequencyReq') },
          { pattern: new RegExp(/^[1-9]\d*$/), message: this.$t('DI.alarmFrequencyFormat') }
        ]
      }
    }
  },
  created () {
    this.form.setAlarm = JSON.parse(JSON.stringify(this.alarmData))
    if (this.setAlarmLength > 1) {
      this.addAlarmValidate(this.setAlarmLength)
    }
  },
  methods: {
    addAlarm () {
      this.form.setAlarm.push({
        alarmType: []
      })
      this.addAlarmValidate(this.form.setAlarm.length)
    },
    deleteAlarm (index) {
      this.deleteAlarmValidate()
      if (this.form.setAlarm.length > 0) {
        this.form.setAlarm.splice(index, 1)
      }
    },
    // 告警类型变化时，对应告警类型数据处理
    changeAlarmType (index) {
      if (!this.form.setAlarm[index]) {
        return
      }
      let alarmObj = this.form.setAlarm[index]
      let alarmType = alarmObj.alarmType
      let basis = {
        'alert_level': '',
        'interval': '',
        'receiver': this.userId
      }
      if (alarmType.indexOf('1') === -1) {
        delete alarmObj.event
      } else {
        if (!alarmObj.event) {
          this.$set(alarmObj, 'event', {
            'alert_level': '',
            'receiver': this.userId,
            'event_checker': ''
          })
        }
      }
      if (alarmType.indexOf('2') === -1) {
        delete alarmObj.fixTime
      } else {
        if (!alarmObj.fixTime) {
          let fixTime = { ...basis }
          fixTime.deadlineChecker = new Date().getTime()
          this.$set(alarmObj, 'fixTime', fixTime)
        }
      }
      if (alarmType.indexOf('3') === -1) {
        delete alarmObj.overTime
      } else {
        if (!alarmObj.overTime) {
          let overTime = { ...basis }
          overTime.overtime_checker = ''
          this.$set(alarmObj, 'overTime', overTime)
        }
      }
    },
    addAlarmValidate (length) {
      if (length <= 1) {
        return
      }
      this.ruleValidate[`setAlarm[${length - 1}].alarmType`] = this.ruleValidate['setAlarm[0].alarmType']
      const eventArr = ['alert_level', 'receiver', 'event_checker']
      this.handelAlarmValidate('event', eventArr, length)
      const fixTimeArr = ['alert_level', 'receiver', 'interval', 'deadlineChecker']
      this.handelAlarmValidate('fixTime', fixTimeArr, length)
      const overTimeArr = ['alert_level', 'receiver', 'interval', 'overtime_checker']
      this.handelAlarmValidate('overTime', overTimeArr, length)
    },
    handelAlarmValidate (label, arr, length) {
      for (let item of arr) {
        this.ruleValidate[`setAlarm[${length - 1}].${label}.${item}`] = this.ruleValidate[`setAlarm[0].${label}.${item}`]
      }
    },
    deleteAlarmValidate () {
      if (this.setAlarmLength < 1) {
        return
      }
      for (let key in this.ruleValidate) {
        if (key.indexOf(`setAlarm[${this.setAlarmLength - 1}]`) > -1) {
          delete this.ruleValidate[key]
        }
      }
    },
    // 提交时告警数据处理
    handleAlarmParam (param) {
      const setAlarm = this.form.setAlarm
      if (setAlarm.length > 0) {
        param.job_alert = {
          event: [],
          deadline: [],
          overtime: []
        }
        let jobAlert = param.job_alert
        for (let item of setAlarm) {
          if (item.event) {
            jobAlert.event.push(item.event)
          }
          if (item.fixTime) {
            let fixTime = {}
            Object.assign(fixTime, item.fixTime)
            fixTime.deadline_checker = util.transDate(item.fixTime.deadlineChecker)
            delete fixTime.deadlineChecker
            jobAlert.deadline.push(fixTime)
          }
          if (item.overTime) {
            jobAlert.overtime.push(item.overTime)
          }
        }
      }
      return param
    },
    clearDialogForm () {
      for (let key in this.ruleValidate) {
        if (key.indexOf('form.setAlarm') > -1 && key.indexOf('form.setAlarm[0]') < 0) {
          delete this.ruleValidate[key]
        }
      }
      this.form.setAlarm = [{
        alarmType: []
      }]
      setTimeout(() => {
        this.$refs.formValidate.resetFields()
      }, 0)
    }
  }
}
</script>
<style lang="scss" scoped>
.alarm-form {
  min-height: 300px;
}
.el-date-editor {
  width: 100%;
}
</style>
