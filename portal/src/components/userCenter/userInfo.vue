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
        <el-form-item label="密码"><el-button @click="passDiagVis = true">修改密码</el-button></el-form-item>
      </el-form>
    </el-card>
    <key-management />
    <el-dialog title="修改密码" :visible.sync="passDiagVis" width="25%" :before-close="resetForm" class="pass-form-container">
      <el-form label-position="right" ref="passForm" :model="passChangeForm" :rules="rules" class="pass-form">
        <el-form-item label="原密码" prop="oldPass"> <el-input v-model="passChangeForm.oldPass" show-password /> </el-form-item>
        <el-form-item label="新密码" prop="newPass"> <el-input v-model="passChangeForm.newPass" show-password /> </el-form-item>
        <el-form-item label="确认密码" prop="conPass">
          <el-input v-model="passChangeForm.conPass" show-password />
        </el-form-item>
      </el-form>
      <span slot="footer" class="dialog-footer">
        <el-button
          @click="
            passDiagVis = false;
            resetForm();
          "
          >取消</el-button
        >
        <el-button type="primary" @click="doChangePass">确定</el-button>
      </span>
    </el-dialog>
  </div>
</template>

<script>
import common from "@/api/common";
import KeyManagement from "./userInfo/keyManagement.vue";

export default {
  name: "userInfo",
  components: {
    KeyManagement
  },
  data() {
    return {
      UserInfo: {},
      rules: {
        conPass: [
          {
            trigger: "blur",
            validator: (rule, value) =>
              new Promise((resolve, reject) => {
                if (value !== this.passChangeForm.newPass) {
                  reject(new Error("两次密码不一致！"));
                } else {
                  resolve();
                }
              })
          }
        ]
      },
      passChangeForm: {
        oldPass: "",
        newPass: "",
        conPass: ""
      },
      type: "text",
      passDiagVis: false
    };
  },
  methods: {
    getUserInfo() {
      common.getUserInfo(this.$store.getters.token).then(resp => {
        this.UserInfo = resp.UserInfo;
      });
    },
    async doChangePass() {
      this.$refs.passForm
        .validate()
        .then(() => {
          const { oldPass, newPass } = this.passChangeForm;
          common.changePassword(oldPass, newPass, this.$store.getters.token).then(resp => {
            if (resp) {
              this.$log(resp);
              this.$message({
                message: "密码修改成功！",
                type: "success"
              });
              this.passDiagVis = false;
            }
          });
        })
        .catch(() => {});
    },
    resetForm(done) {
      this.passChangeForm = {
        oldPass: "",
        newPass: "",
        conPass: ""
      };
      this.$refs.passForm.resetFields();
      if (done) done();
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
  beforeMount() {
    this.getUserInfo();
  }
};
</script>

<style scoped lang="scss">
.pass-form-container {
  min-width: 300px;
  .pass-form {
    margin: auto;
    width: 300px;
  }
}

.user-info-container {
  .el-form {
    width: 50%;
  }
}

.box-card {
  margin: 10px 0;
  .inner-box-card {
    width: 30%;
  }
}
</style>
