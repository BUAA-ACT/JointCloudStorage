<template>
  <div :class="!mini ? 'el-header' : 'mini-header'">
    <el-row :gutter="10">
      <el-col :xs="16" :sm="16" :md="16" :lg="14">
        <a href="/"><img src="@/assets/logo.png" class="logoImg" alt="logo"/></a>
        <span class="logoImg menu-header">
          <!--          <span>{{ curCloudName }}</span>-->
          <svg viewBox="0 0 300 26" class="svgBox">
            <defs>
              <linearGradient id="changer" x1="0" y1="0" x2="100%" y2="0">
                <stop offset="0" stop-color="rgb(43, 111, 193)" />
                <stop offset="50%" stop-color="rgb(49, 140, 179)" />
                <stop offset="100%" stop-color="rgb(56, 172, 164)" />
              </linearGradient>
            </defs>
            <text text-anchor="start" fill="url(#changer)" x="0" y="1em">{{ curCloudName }}</text>
          </svg>
        </span>
      </el-col>
      <el-col :xs="8" :sm="8" :md="8" :lg="10" class="insideFloatRight">
        <el-link v-if="!userName" href="#/register">注册</el-link>
        <el-link v-if="!userName" href="#/login">登录</el-link>
        <!-- <el-input
          class="search"
          placeholder="搜索"
          suffix-icon="el-icon-search"
          v-model="search">
        </el-input> -->
        <div v-if="userName">
          <el-dropdown class="el-dropdown" @command="handleCommand">
            <span class="el-dropdown-link"> {{ userName }}<i class="el-icon-arrow-down el-icon--right"></i> </span>
            <el-dropdown-menu slot="dropdown">
              <el-dropdown-item command="userCenter">个人中心</el-dropdown-item>
              <el-dropdown-item command="manageCenter" v-if="isAdmin">管理中心</el-dropdown-item>
              <el-dropdown-item command="logOut">退出登录</el-dropdown-item>
            </el-dropdown-menu>
          </el-dropdown>
        </div>
      </el-col>
      <el-col :md="24" :lg="24" v-if="!mini">
        <el-menu :default-active="activeIndex" class="el-menu" mode="horizontal" @select="handleSelect">
          <el-menu-item index="1">首页</el-menu-item>
          <el-menu-item index="2">需求目录</el-menu-item>
          <el-menu-item index="3">服务目录</el-menu-item>
          <el-menu-item index="4">集群服务</el-menu-item>
          <el-menu-item index="5">解决方案</el-menu-item>
          <el-menu-item index="6">技术成果</el-menu-item>
        </el-menu>
      </el-col>
    </el-row>
  </div>
</template>

<script>
import Cloud from "@/api/clouds";

export default {
  name: "headNav",
  props: {
    mini: {
      type: Boolean,
      default: false
    }
  },
  data() {
    return {
      search: "",
      activeIndex: "1",
      navList: {
        1: "/",
        2: "/requirementList",
        3: "/serviceList",
        4: "/clusterService",
        5: "/solution",
        6: "/404"
      },
      curCloudName: "阿里云青岛"
    };
  },
  created() {
    Object.keys(this.navList).forEach(element => {
      if (this.$route.path === this.navList[element]) {
        this.activeIndex = element;
      }
    });
  },
  computed: {
    userName() {
      return this.$store.state.name;
    },
    isAdmin() {
      return this.$store.getters.isAdmin;
    }
  },
  methods: {
    handleSelect(key) {
      if (this.$route.path !== this.navList[key]) {
        this.$router.push(this.navList[key]);
      }
    },
    handleCommand(command) {
      if (command !== "logOut") {
        this.$router.push({ path: `/${command}`, query: this.$route.query });
      } else {
        this.$store.dispatch("logout");
        this.$router.push({ path: "/", query: { redirect: this.$route.path } });
      }
    },
    getCurCloudName() {
      Cloud.getCurCloudName().then(resp => {
        if (resp) {
          this.curCloudName = resp.CloudName;
        }
      });
    }
  },
  beforeMount() {
    this.getCurCloudName();
  }
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped lang="scss">
.el-menu-item {
  border: 0;
}
.el-menu-item {
  padding: 0;
}
.el-menu--horizontal > .el-menu-item {
  margin-right: 80px;
  color: #333333;
  font-size: 14px;
  height: 50px;
  line-height: 50px;
  @media screen and (max-width: 1200px) {
    margin-right: 50px;
  }
}
.el-menu--horizontal > .el-menu-item:last-child {
  margin-right: 0;
}
.el-menu--horizontal > .el-menu-item.is-active,
.el-menu--horizontal > .el-submenu.is-active .el-submenu__title {
  border: 0;
  color: #419eff;
}

.el-menu {
  float: left;
  border: 0;
}
// .el-dropdown{
//   line-height: 40px;
// }
.insideFloatRight .el-dropdown {
  margin-top: 5px;
  line-height: 30px;
  float: right;
}
.el-dropdown-menu {
  width: 120px;
}

h3 {
  margin: 40px 0 0;
}
ul {
  list-style-type: none;
  padding: 0;
}
li {
  display: inline-block;
  margin: 0 10px;
}
// a {
//   color: #42b983;
// }
.el-input__inner {
  border-radius: 50%;
}

.menu-header {
  font-size: 20px;
  color: rgb(43, 111, 193);
  letter-spacing: 5px;
  float: left;
  margin-top: 32px;
  height: 22px;
}

.svgBox {
  width: fit-content;
  height: 22px;
  vertical-align: baseline;
  position: relative;
  top: 2px;
  text-align: left;
}
</style>
