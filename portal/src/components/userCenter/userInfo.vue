<template>
  <div style="text-align: left">
    <el-button @click="getUserInfo" type="primary">获取用户信息</el-button><br />
    <!--    <JsonViewer :value="UserInfo"></JsonViewer>-->

    <el-card class="box-card user-info-container">
      <div slot="header" class="clearfix">
        <span>用户信息</span>
      </div>
      <el-form label-width="70px" label-position="left">
        <el-form-item label="用户名">{{ email }}</el-form-item>
        <el-form-item label="密码"><el-button @click="changePassword">修改密码</el-button></el-form-item>
      </el-form>
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
    <el-card class="box-card">
      <div slot="header" class="clearfix">
        <span>访问凭证-NEW</span>
      </div>
      <el-card class="inner-box-card" shadow="hover">
        <div slot="header" class="clearfix">
          <span></span>
        </div>
        <div class="text item">
          <el-form label-position="right">
            <el-form-item label="Access Key"> <el-input :value="ak" class="no-border" size="large" /> </el-form-item>
            <el-form-item label="Secret Key"> <el-input :value="sk" show-password class="no-border" size="large"/></el-form-item>
          </el-form>
        </div>
      </el-card>
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
      },
      ak: "application_key",
      sk: "SeCreTKey",
      type: "text"
    };
  },
  methods: {
    getUserInfo() {
      common.getUserInfo(this.$store.getters.token).then(resp => {
        this.UserInfo = resp.UserInfo;
      });
    },
    async changePassword() {
      // TODO
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
.user-info-container {
  .el-form {
    width: 50%;
  }
}
.box-card {
  margin: 10px 0;
  .inner-box-card {
    width: 30%;
    .no-border /deep/ .el-input__inner {
      border: 0;
    }
  }
}
</style>
