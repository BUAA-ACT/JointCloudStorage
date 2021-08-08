<template>
  <div class="login-container">
    <div class="headerBg">
      <div class="el-header">
        <el-row :gutter="50">
          <el-col :md="14" :lg="14">
            <a href="/"><img src="@/assets/logo.png" class="logoImg" alt="logoImg"/></a>
            <span class="title">国家重点研发计划</span>
          </el-col>
          <el-col v-if="true" :md="10" :lg="10" class="insideFloatRight">
            <el-link href="#/login">用户登录</el-link>
          </el-col>
        </el-row>
      </div>
    </div>
    <div class="manufacturer bg-white">
      <el-form ref="registerForm" :model="regForm" :rules="regRules" class="register-form" autocomplete="on" label-position="left" label-width="80px">
        <p class="title">欢迎注册云际计算平台</p>

        <el-form-item prop="Email" label="用户名">
          <el-input
            ref="username"
            v-model="regForm.Email"
            placeholder="您在本系统的唯一标识"
            name="username"
            type="text"
            tabindex="1"
            autocomplete="on"
          />
        </el-form-item>

        <!--        <el-form-item prop="Nickname" label="昵称">-->
        <!--          <el-input ref="nickname" v-model="regForm.Nickname" placeholder="默认为用户名" name="nickname" type="text" tabindex="2" autocomplete="on" />-->
        <!--        </el-form-item>-->

        <!--         <el-form-item prop="email" label="邮箱">
          <el-input ref="email" v-model="regForm.email" placeholder="" name="email" type="text" tabindex="1" autocomplete="on" />
        </el-form-item>

         <el-form-item prop="emailCode" label="验证码">
          <el-input
            ref="emailCode"
            v-model="regForm.emailCode"
            placeholder=""
            name="emailCode"
            type="text"
            tabindex="1"
            autocomplete="on"
          />
        </el-form-item> -->

        <el-tooltip :value="passwordTip && capsTooltip" content="大写锁定已打开" placement="right" manual>
          <el-form-item prop="Password" label="密码">
            <el-input
              ref="password"
              v-model="regForm.Password"
              placeholder="6-16个字符组成，区分大小写"
              name="password"
              tabindex="3"
              autocomplete="off"
              @keyup.native="checkCapslock"
              @focus="passwordTip = true"
              @blur="passwordTip = false"
              show-password
            />
          </el-form-item>
        </el-tooltip>

        <!--                <el-tooltip v-model="capsTooltip" content="Caps lock is On" placement="right" manual>
                  <el-form-item prop="confirmPassword" label="确认密码">
                    <el-input
                      ref="confirmPassword"
                      v-model="regForm.password"
                      placeholder="请再输入一遍密码"
                      name="password"
                      tabindex="4"
                      autocomplete="on"
                      @keyup.native="checkCapslock"
                      @blur="capsTooltip = false"
                      show-password
                    />
                  </el-form-item>
                </el-tooltip>-->

        <el-tooltip :value="confirmTip && capsTooltip" content="大写锁定已打开" placement="right" manual>
          <el-form-item prop="confirmPassword" label="确认密码">
            <el-input
              ref="confirmPassword"
              placeholder="请再输入一遍密码"
              name="confirmPassword"
              v-model="regForm.confirmPassword"
              tabindex="4"
              autocomplete="off"
              @keyup.native="checkCapslock"
              @focus="confirmTip = true"
              @blur="confirmTip = false"
              show-password
            />
          </el-form-item>
        </el-tooltip>

        <el-tooltip :value="!agreeCheck" content="您必须勾选此项才能继续注册" placement="left" manual>
          <el-checkbox class="margin-20" v-model="agreeCheck" tabindex="5"
            >我已阅读并同意
            <el-link type="primary" href="/legal-agreement.html">《云际计算平台注册协议》</el-link>
          </el-checkbox>
        </el-tooltip>

        <el-button :loading="loading" :disabled="!agreeCheck" type="primary" style="width:100%;margin-bottom:30px;" @click="handleRegister"
          >注册</el-button
        >

        <div style="position:relative">
          <div class="tips">
            <p><a href="/">返回首页</a></p>
          </div>
        </div>
      </el-form>
    </div>
    <!--    <Footer type="mini" />-->
  </div>
</template>

<script>
import Common from "@/api/common";

export default {
  name: "Register",
  components: {},
  data() {
    return {
      regForm: {
        Email: "",
        Password: ""
        // Nickname: ""
      },
      regRules: {
        Email: [{ required: true, message: "请输入用户名", trigger: "blur" }],
        // Nickname: [
        //   {
        //     type: "string",
        //     trigger: "blur",
        //     validator: (rule, value) => {
        //       return new Promise((resolve, reject) => {
        //         if (value === "") {
        //           this.regForm.Nickname = this.regForm.Email;
        //           if (this.regForm.Email === "") {
        //             reject(new Error("请输入昵称"));
        //           } else {
        //             reject(new Error("已使用用户名作为昵称，你可以继续修改"));
        //           }
        //         } else {
        //           resolve();
        //         }
        //       });
        //     }
        //   },
        //   { required: true, message: "请输入昵称", trigger: "blur" }
        // ],
        Password: [
          { required: true, message: "请输入密码", trigger: "blur" }
          // { min: 6, max: 16, message: "长度在 6 到 16 个字符", trigger: "blur" }
        ],
        confirmPassword: [
          { required: true, message: "请输入确认密码", trigger: "blur" },
          {
            type: "string",
            trigger: "blur",
            validator: (rule, value) => {
              return new Promise((resolve, reject) => {
                if (value !== this.regForm.Password) {
                  reject(new Error("两次输入的密码不一致"));
                } else {
                  resolve();
                }
              });
            }
          }
        ]
      },
      agreeCheck: false,
      capsTooltip: false,
      passwordTip: false,
      confirmTip: false,
      loading: false,
      redirect: undefined,
      otherQuery: {}
    };
  },
  watch: {
    $route: {
      handler(route) {
        const { query } = route;
        if (query) {
          this.redirect = query.redirect;
          // this.otherQuery = this.getOtherQuery(query);
        }
      },
      immediate: true
    }
  },
  created() {
    // window.addEventListener('storage', this.afterQRScan)
  },
  mounted() {
    this.$refs.username.focus();
  },
  destroyed() {
    // window.removeEventListener('storage', this.afterQRScan)
  },
  methods: {
    checkCapslock(e) {
      const { key } = e;
      this.capsTooltip = key && key.length === 1 && key >= "A" && key <= "Z";
    },
    async handleRegister() {
      this.loading = true;
      if (!this.agreeCheck) {
        this.loading = false;
        return;
      }
      this.$refs.registerForm
        .validate()
        .then(() => {
          Common.register(this.regForm)
            .then(resp => {
              if (resp) {
                // console.log(resp);
                this.$store.dispatch("login", {
                  username: this.regForm.Email,
                  password: this.regForm.Password
                });
                this.$message.success("注册成功");
                this.$router.push({ path: this.redirect || "/" });
              }
            })
            .catch(() => {
              // console.log(e);
            });
          this.loading = false;
        })
        .catch(() => {
          this.loading = false;
        });
    }
  }
};
</script>
<style scoped lang="scss">
.register-form {
  width: 300px;
  margin: 0 auto;
  padding: 40px 0;
  .title {
    font-size: 28px;
  }
}
.bg-white {
  background-color: #ffffff;
  margin-bottom: 70px;
}
</style>
