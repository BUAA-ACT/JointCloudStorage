<template>
  <el-progress :percentage="curPercent" :stroke-width="20" text-inside></el-progress>
</template>

<script>
import MyWS from "@/api/WebSocket";

export default {
  name: "migrationProcess",
  data() {
    return {
      ws: null,
      curPercent: 0
    };
  },
  methods: {
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
      this.removeActions();
      this.ws.close();
      this.$log("CONNECTION CLOSED");
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
    }
  },
  beforeMount() {
    this.ws = new MyWS("/task/getMigration");
    this.ws.addAction(this.initialMsg, "open");
    this.ws.addAction(this.updateProgress, "message");
    this.ws.addAction(this.closeConnection, "close");
    this.ws.addAction(this.retryConnection, "error");
  },
  beforeDestroy() {
    if (this.ws) {
      this.removeActions();
      this.ws.close();
      this.ws = null;
    }
  }
};
</script>

<style scoped></style>
