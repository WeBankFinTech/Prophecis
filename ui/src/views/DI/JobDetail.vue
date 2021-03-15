<template>
  <div>
    <el-form label-width="190px">
      <div class="subtitle">
        {{ $t('DI.basicSettings') }}
      </div>
      <el-row>
        <el-col :span="12">
          <el-form-item :label="$t('DI.trainingjobName')" maxlength="225">
            <el-input v-model="detailObj.name" disabled />
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item :label="$t('ns.nameSpace')">
            <el-select v-model="detailObj.namespace" disabled>
              <el-option v-for="item in spaceOptionList" :label="item" :key="item" :value="item">
              </el-option>
            </el-select>
          </el-form-item>
        </el-col>
      </el-row>
      <el-row>
        <el-form-item :label="$t('DI.description')">
          <el-input v-model="detailObj.description" type="textarea" disabled />
        </el-form-item>
      </el-row>
      <div class="subtitle">
        {{ $t('DI.imageSettings') }}
      </div>
      <el-row>
        <el-form-item :label="$t('DI.imageSelection')">
          <span>&nbsp;{{ defineImage }}&nbsp;&nbsp;</span>
          <el-input v-model="detailObj.image" disabled />
        </el-form-item>
      </el-row>
      <div class="subtitle">
        {{ $t('DI.computingResource') }}
      </div>
      <el-row>
        <el-col :span="12">
          <el-form-item :label="$t('DI.jobType')">
            <el-select v-model="detailObj.job_type" disabled>
              <el-option :label="$t('DI.Single')" value="Local">
              </el-option>
              <el-option :label="$t('DI.distributed')" value="dist-tf">
              </el-option>
            </el-select>
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item :label="$t('DI.cpu')">
            <el-input v-model="detailObj.cpus" disabled>
              <span slot="append">Core</span>
            </el-input>
          </el-form-item>
        </el-col>
      </el-row>
      <el-row>
        <el-col :span="12">
          <el-form-item :label="$t('DI.gpu')">
            <el-input v-model="detailObj.gpus" disabled>
              <span slot="append">{{ $t('DI.block') }}</span>
            </el-input>
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item :label="$t('DI.memory')">
            <el-input v-model="detailObj.memory" disabled>
              <span slot="append">Gb</span>
            </el-input>
          </el-form-item>
        </el-col>
      </el-row>
      <div v-if="detailObj.job_type==='dist-tf'">
        <el-row>
          <el-col :span="12">
            <el-form-item :label="$t('DI.learners')">
              <el-input v-model="detailObj.learners" disabled />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item :label="$t('DI.pss')">
              <el-input v-model="detailObj.pss" disabled />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row>
          <el-col :span="12">
            <el-form-item :label="$t('DI.ps_cpu')">
              <el-input v-model="detailObj.ps_cpu" disabled>
                <span slot="append">Core</span>
              </el-input>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item :label="$t('DI.ps_memory')">
              <el-input v-model="detailObj.ps_memory" disabled>
                <span slot="append">Gi</span>
              </el-input>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row>
          <el-col :span="24">
            <el-form-item :label="$t('DI.ps_image')">
              <span>&nbsp;{{ defineImage }}&nbsp;&nbsp;</span>
              <el-input v-model="detailObj.ps_image" disabled />
            </el-form-item>
          </el-col>
        </el-row>
      </div>
      <div class="subtitle">
        {{ $t('DI.trainingDirectory') }}
      </div>
      <el-row>
        <el-form-item :label="$t('DI.trainingDataStore')">
          <el-input v-model="detailObj.path" disabled />
        </el-form-item>
        <el-form-item :label="$t('DI.trainingData')">
          <el-input v-model="detailObj.trainingData" disabled />
        </el-form-item>
        <el-form-item :label="$t('DI.trainingResult')">
          <el-input v-model="detailObj.trainingResults" disabled />
        </el-form-item>
      </el-row>
      <div v-if="detailObj.job_alert">
        <div v-for="(item,index) in detailObj.job_alert" :key="index">
          <div class="subtitle ">
            {{ $t('DI.alarmSet') }}
          </div>
          <el-row>
            <el-form-item :label="$t('DI.alarmTypes')" prop="alarmType">
              <el-select v-model="item.alarmType" multiple :placeholder="$t('DI.alarmTypesPro')" disabled>
                <el-option :label="$t('DI.eventAlarm')" value="1">
                </el-option>
                <el-option :label="$t('DI.designatedTimeout')" value="2">
                </el-option>
                <el-option :label="$t('DI.hoursOvertime')" value="3">
                </el-option>
              </el-select>
            </el-form-item>
          </el-row>
          <div v-if="item.alarmType.indexOf('1')>-1" key="eventChecker">
            <div class="subtitle">
              {{ $t('DI.eventAlarmSettings') }}
            </div>
            <el-row>
              <el-col :span="12">
                <el-form-item :label="$t('DI.listeningEnent')" prop="eventChecker">
                  <el-select v-model="item.event.event_checker" :placeholder="$t('DI.listeningEnentPro')" disabled>
                    <el-option :label="$t('DI.executeSuccess')" value="COMPLETED">
                    </el-option>
                    <el-option :label="$t('DI.executeFail')" value="FAILED">
                    </el-option>
                  </el-select>
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item :label="$t('DI.alarmLevel')" prop="alertLevel">
                  <el-select v-model="item.event.alert_level" :placeholder="$t('DI.alarmLevelPro')" disabled>
                    <el-option label="Info" value="info">
                    </el-option>
                    <el-option label="Warning" value="warning">
                    </el-option>
                    <el-option label="Minor" value="minor">
                    </el-option>
                    <el-option label="Major" value="major">
                    </el-option>
                    <el-option label="Critical" value="critical">
                    </el-option>
                  </el-select>
                </el-form-item>
              </el-col>
            </el-row>
            <el-row>
              <el-form-item :label="$t('DI.receiver')" prop="receiver">
                <el-input v-model="item.event.receiver" :placeholder="$t('DI.receiver')" disabled />
              </el-form-item>
            </el-row>
          </div>
          <div v-if="item.alarmType.indexOf('2')>-1" key="designatedTimeout">
            <div class="subtitle">
              {{ $t('DI.designatedTimeoutSetting') }}
            </div>
            <el-row>
              <el-col :span="12">
                <el-form-item :label="$t('DI.designatedTimeout')" prop="deadlineChecker">
                  <el-input v-model="item.fixTime.deadline_checker" :placeholder="$t('DI.timeoutTimePro')" disabled />
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item :label="$t('DI.alarmFrequency')" prop="interval">
                  <el-input v-model="item.fixTime.interval" :placeholder="$t('DI.alarmFrequencyPro')" disabled>
                    <span slot="append">min</span>
                  </el-input>
                </el-form-item>
              </el-col>
            </el-row>
            <el-row>
              <el-col :span="12">
                <el-form-item :label="$t('DI.alarmLevel')" prop="alertLevel">
                  <el-select v-model="item.fixTime.alert_level" :placeholder="$t('DI.alarmLevelPro')" disabled>
                    <el-option label="Info" value="info">
                    </el-option>
                    <el-option label="Warning" value="warning">
                    </el-option>
                    <el-option label="Minor" value="minor">
                    </el-option>
                    <el-option label="Major" value="major">
                    </el-option>
                    <el-option label="Critical" value="critical">
                    </el-option>
                  </el-select>
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item :label="$t('DI.receiver')" prop="receiver">
                  <el-input v-model="item.fixTime.receiver" :placeholder="$t('DI.receiver')" disabled />
                </el-form-item>
              </el-col>
            </el-row>
          </div>
          <div v-if="item.alarmType.indexOf('3')>-1" key="overtimeChecker">
            <div class="subtitle">
              {{ $t('DI.timeoutTimeSetting') }}
            </div>
            <el-row>
              <el-col :span="12">
                <el-form-item :label="$t('DI.hoursOvertime')" prop="overtimeChecker">
                  <el-input v-model="item.overTime.overtime_checker" :placeholder="$t('DI.timeoutTimePro')" disabled>
                    <span slot="append">h</span>
                  </el-input>
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item :label="$t('DI.alarmFrequency')" prop="interval">
                  <el-input v-model="item.overTime.interval" :placeholder="$t('DI.alarmFrequencyPro')" disabled>
                    <span slot="append">min</span>
                  </el-input>
                </el-form-item>
              </el-col>
            </el-row>
            <el-row>
              <el-col :span="12">
                <el-form-item :label="$t('DI.alarmLevel')" prop="alertLevel">
                  <el-select v-model="item.overTime.alert_level" :placeholder="$t('DI.alarmLevelPro')" disabled>
                    <el-option label="Info" value="info">
                    </el-option>
                    <el-option label="Warning" value="warning">
                    </el-option>
                    <el-option label="Minor" value="minor">
                    </el-option>
                    <el-option label="Major" value="major">
                    </el-option>
                    <el-option label="Critical" value="critical">
                    </el-option>
                  </el-select>
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item :label="$t('DI.receiver')" prop="receiver">
                  <el-input v-model="item.overTime.receiver" :placeholder="$t('DI.receiver')" disabled />
                </el-form-item>
              </el-col>
            </el-row>
          </div>
        </div>
      </div>
      <div class="subtitle">
        {{ $t('DI.jobExecution') }}
      </div>
      <el-row>
        <el-form-item :label="$t('DI.entrance')">
          <el-input v-model="detailObj.command" :rows="4" type="textarea" disabled />
        </el-form-item>
        <el-form-item :label="$t('DI.executeCodeSettings')">
          <el-select v-model="detailObj.codeSettings" disabled>
            <el-option :label="$t('DI.manualUpload')" value="codeFile">
            </el-option>
            <el-option :label="$t('DI.shareDirectory')" value="storagePath">
            </el-option>
          </el-select>
        </el-form-item>
        <div class="upload-box" v-if="detailObj.codeSettings==='codeFile'">
          <span class="file-name">
            {{detailObj.fileName}}
          </span>
        </div>
        <el-form-item v-if="
              detailObj.codeSettings==='storagePath'" key="
              detailObjStoragePath" :label="$t('DI.shareDirectory')">
          <el-input v-model="detailObj.diStoragePath" disabled />
        </el-form-item>
      </el-row>
    </el-form>
  </div>
</template>

<script>
export default {
  props: {
    detailObj: {
      type: Object,
      default: function () {
        return {

        }
      }
    },
    spaceOptionList: {
      type: Array,
      default: function () {
        return []
      }
    }
  },
  computed: {
    defineImage: function () {
      return this.FesEnv.DI.defineImage
    }
  }
}
</script>

<style lang="scss" scoped>
.file-name {
  margin-left: 190px;
}
</style>
