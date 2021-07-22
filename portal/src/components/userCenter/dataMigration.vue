<template>
  <div>
    <plan-advice v-if="!migrating" />
    <div class="no-new-plan" v-else>
      <i class="el-icon-info tip-icon"></i><br />
      数据迁移中……
      <!-- TODO: 迁移进度     -->
      <div>
        <span>当前进度：</span>
        <el-progress :percentage="curPercent" :stroke-width="20" text-inside :status="status"></el-progress>
      </div>
      <location-viewer
        :clouds="Advices.CloudsOld"
        :new-clouds="Advices.CloudsNew"
        :dynamic="migrating"
        :inactive-clouds="allClouds"
        class="location-viewer"
        ref="viewer"
      />
    </div>
  </div>
</template>

<script>
import PlanAdvice from "@/components/userCenter/planAdvice.vue";
import MyWS from "@/api/WebSocket";
import Plan from "@/api/plan";
import Clouds from "@/api/clouds";
import locationViewer from "@/components/viewer/locationViewer.vue";

export default {
  name: "dataMigration",
  components: { locationViewer, PlanAdvice },
  data() {
    return {
      intervalId: null,
      ws: null,
      curPercent: 0,
      status: null,
      Advices: {},
      allClouds: []
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
      this.curPercent = (data.Progress * 100).toFixed(2);
      if (data.TaskState === "FINISH") {
        this.status = "success";
      }
    },
    closeConnection() {
      this.removeActions();
      this.ws.close();
      this.ws = null;
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
      this.ws.addAction(this.retryConnection, "close");
      this.ws.addAction(this.retryConnection, "error");
    },
    removeActions() {
      this.ws.removeAction(this.initialMsg, "open");
      this.ws.removeAction(this.updateProgress, "message");
      this.ws.removeAction(this.retryConnection, "close");
      this.ws.removeAction(this.retryConnection, "error");
    },
    async getNewAdvice() {
      Plan.getNewAdvice().then(resp => {
        if (resp) {
          if (resp.Advices && resp.Advices.length > 0) {
            [this.Advices] = resp.Advices;
          } else {
            this.Advices = {};
          }
        }
      });
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
    }
  },
  mounted() {
    this.intervalId = setInterval(this.checkCurStatus, 2000);
    this.getNewAdvice();
    this.getAllClouds();
  },
  watch: {
    migrating(newVal) {
      if (newVal && !this.ws) {
        this.ws = new MyWS("/task/getMigration");
        this.ws.addAction(this.initialMsg, "open");
        this.ws.addAction(this.updateProgress, "message");
        this.ws.addAction(this.closeConnection, "close");
        this.ws.addAction(this.retryConnection, "error");
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
</style>
