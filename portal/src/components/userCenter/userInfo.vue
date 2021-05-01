<template>
  <div style="text-align: left">
    <el-button @click="getUserInfo" type="primary">获取用户信息</el-button><br />
    <!--    <JsonViewer :value="UserInfo"></JsonViewer>-->

    <el-card class="box-card">
      <div slot="header" class="clearfix">
        <span>用户信息</span>
      </div>
      用户名：{{ email }}
    </el-card>
    <el-card class="box-card">
      <div slot="header" class="clearfix">
        <span>访问凭证</span>
      </div>
      <el-card class="inner-box-card" v-for="AccessCredential in UserInfo.AccessCredentials" :key="AccessCredential.CloudID" shadow="hover">
        <div slot="header" class="clearfix">
          <span>{{ AccessCredential.CloudID }}</span>
        </div>
        <div class="text item">
          访问用户：{{ AccessCredential.UserID }}<br />
          访问密码：{{ AccessCredential.Password }}<br />
        </div>
      </el-card>
    </el-card>
  </div>
</template>

<script>
import common from "@/api/common";

export default {
  name: "userInfo",
  data() {
    return {
      UserInfo: {
        AccessCredentials: []
      }
    };
  },
  methods: {
    getUserInfo() {
      common.getUserInfo(this.$store.getters.token).then(resp => {
        this.UserInfo = resp.UserInfo;
      });
    }
  },
  computed: {
    email() {
      return this.$store.getters.name;
    },
    AccessCredentials() {
      return this.UserInfo.AccessCredentials;
    }
  },
  created() {
    this.getUserInfo();
  }
};
</script>

<style scoped lang="scss">
.box-card {
  margin: 10px 0;
  .inner-box-card {
    width: 30%;
  }
}
</style>
