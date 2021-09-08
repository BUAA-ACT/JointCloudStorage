<template>
  <div>
    <div v-if="newPlanAvailable && !migrating" class="migration-operations">
      <el-button type="primary" @click="accept" :loading="submitLoading">迁移</el-button><br />
      <!--      <el-button type="primary" @click="cancel" :loading="cancelLoading" :disabled="!plansLoaded">我再想想</el-button><br />-->
      <el-button type="info" @click="cancel" :loading="submitLoading">取消</el-button><br />
    </div>
    <div v-if="migrating || newPlanAvailable">
      <div class="plans-viewer-container">
        <plan-detail :plan="Advices.StoragePlanOld">
          旧存储方案
        </plan-detail>
        <div class="migration-progress">
          <div v-if="migrating" class="migration-progress-bar">
            <el-progress
              :percentage="curPercent"
              :stroke-width="10"
              :status="status"
              type="circle"
              :format="formatPercent"
              :width="200"
            ></el-progress>
          </div>
        </div>
        <plan-detail :plan="Advices.StoragePlanNew">
          新存储方案
        </plan-detail>
      </div>
    </div>
    <location-viewer
      :clouds="Advices.CloudsOld"
      :new-clouds="Advices.CloudsNew"
      :dynamic="migrating"
      :inactive-clouds="allClouds"
      class="location-viewer"
      ref="viewer"
      v-if="migrating || newPlanAvailable"
    />
    <div class="no-new-plan" v-if="!newPlanAvailable && !migrating">
      <i class="el-icon-success tip-icon"></i><br />
      暂时没有为您找到更优的存储方案！
    </div>
  </div>
</template>

<script>
import MyWS from "@/api/WebSocket";
import Plan from "@/api/plan";
import Clouds from "@/api/clouds";
import locationViewer from "@/components/viewer/locationViewer.vue";
import PlanDetail from "@/components/userCenter/planDetail.vue";

export default {
  name: "dataMigration",
  components: { PlanDetail, locationViewer },
  inject: ["formatPrice", "reload"],
  data() {
    return {
      intervalId: null,
      ws: null,
      curPercent: 0,
      status: null,
      Advices: {},
      allClouds: [],
      submitLoading: false,
      newPlanAvailable: false
    };
  },
  computed: {
    migrating() {
      return this.$store.getters.status === "FORBIDDEN";
    }
  },
  methods: {
    checkCurStatus() {
      this.$store.dispatch("updateInfo", "Status");
    },
    initialMsg() {
      this.ws.send({ AccessToken: this.$store.getters.token });
    },
    updateProgress(event) {
      let data = { TaskID: "", TaskType: "", TaskState: "", TaskStartTime: "", UserID: "", Progress: 0 };
      const resp = JSON.parse(event.data);
      if (resp) {
        data = resp.data;
      }
      this.$log(data);
      this.curPercent = Number(data.Progress.toFixed(2));
      if (data.TaskState === "FINISH") {
        this.status = "success";
        this.getNewAdvice();
      }
    },
    closeConnection() {
      this.removeActions();
      this.ws.close();
      this.ws = null;
      this.getNewAdvice();
    },
    retryConnection() {
      this.$log("CONNECTION CLOSED due to ERROR");
      this.removeActions();
      this.ws.close();
      this.$log("Trying to establish new WebSocket Connection...");
      this.ws = new MyWS("/task/getMigration");
    },
    addActions() {
      this.ws.addAction(this.initialMsg, "open");
      this.ws.addAction(this.updateProgress, "message");
      this.ws.addAction(this.closeConnection, "close");
      this.ws.addAction(this.retryConnection, "error");
    },
    removeActions() {
      this.ws.removeAction(this.initialMsg, "open");
      this.ws.removeAction(this.updateProgress, "message");
      this.ws.removeAction(this.closeConnection, "close");
      this.ws.removeAction(this.retryConnection, "error");
    },
    async getNewAdvice() {
      Plan.getNewAdvice().then(resp => {
        if (resp) {
          if (resp.Advices && resp.Advices.length > 0) {
            [this.Advices] = resp.Advices;
            this.newPlanAvailable = this.Advices.Status === "PENDING";
          } else {
            this.Advices = {};
          }
        }
      });
    },
    async accept() {
      this.submitLoading = true;
      await Plan.acceptAdvice().then(resp => {
        if (resp) {
          this.$notify.success({
            title: "数据迁移已开始！",
            message: (
              <span>
                请耐心等待迁移完成……
                <br />
                在此期间，您的数据将暂时不可用。
              </span>
            )
          });
          this.$emit("accept");
        }
      });
      this.submitLoading = false;
    },
    async cancel() {
      this.submitLoading = true;
      await Plan.abandonAdvice();
      this.$emit("cancel");
      this.submitLoading = false;
    },
    async getAllClouds() {
      Clouds.getAllClouds()
        .then(res => {
          if (res && res.Clouds) {
            this.allClouds = res.Clouds;
          }
        })
        .catch(() => {})
        .then(() => {
          this.allClouds = this.allClouds || [];
        });
    },
    formatPercent(per) {
      return per === 100 ? "迁移已完成！" : `迁移中:${per}%`;
    }
  },
  mounted() {
    this.intervalId = setInterval(this.checkCurStatus, 2000);
    this.getNewAdvice();
    this.getAllClouds();
    if (this.migrating && !this.ws) {
      this.ws = new MyWS("/task/getMigration");
      this.addActions();
    }
  },
  watch: {
    migrating(newVal) {
      if (newVal && !this.ws) {
        this.ws = new MyWS("/task/getMigration");
        this.addActions();
      }
    }
  },
  beforeDestroy() {
    if (this.intervalId) {
      clearInterval(this.intervalId);
    }
    if (this.ws) {
      this.removeActions();
      this.ws.close();
      this.ws = null;
    }
  }
};
</script>

<style scoped lang="scss">
.tip-icon {
  font-size: 100px;
  color: #909399;
}
.no-new-plan {
  .tip-icon {
    font-size: 100px;
    color: #67c23a;
  }
}
.migration-operations {
  display: flex;
  width: 200px;
  justify-content: space-between;
}
.plans-viewer-container {
  display: -webkit-box;
  display: -webkit-flex;
  display: -ms-flexbox;
  display: flex;
  margin-top: 30px;
  -webkit-box-pack: justify;
  -webkit-justify-content: space-between;
  -ms-flex-pack: justify;
  justify-content: space-between;
  width: 1400px;
}
.location-viewer {
  width: 1400px;
  height: 400px;
}
.text {
  font-size: 14px;
  text-align: left;
}

.item {
  margin-bottom: 18px;
}

.clearfix:before,
.clearfix:after {
  display: table;
  content: "";
}
.clearfix:after {
  clear: both;
}

.box-card {
  width: 250px;
  height: 200px;
  display: inline-block;
  margin: 0 10px;
}
.migration-progress {
  width: 600px;
  .migration-progress-text {
    margin: 20px 0;
  }
  .migration-progress-bar {
    margin: 50px;
    /deep/ i {
      font-size: 50px;
      font-weight: 1000;
    }
  }
}
</style>
