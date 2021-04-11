<template>
  <div style="text-align: left">
    <el-button @click="getUserInfo">获取用户信息</el-button><br />
    <JsonViewer :value="UserInfo"></JsonViewer>
  </div>
</template>

<script>
import JsonViewer from "vue-json-viewer";
import common from "@/api/common";

export default {
  name: "userInfo",
  components: {
    JsonViewer
  },
  data() {
    return {
      UserInfo: {}
    };
  },
  methods: {
    getUserInfo() {
      common.getUserInfo(this.$store.getters.token).then(resp => {
        this.UserInfo = resp;
      });
    }
  },
  created() {
    this.getUserInfo();
  }
};
</script>

<style scoped></style>
