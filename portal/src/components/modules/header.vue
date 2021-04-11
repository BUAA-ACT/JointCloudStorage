<template>
  <div class="el-header s-header">
    <el-row :gutter="10">
      <el-col :xs="10" :sm="10" :md="10" :lg="10">
        <a href="/"><img src="@/assets/logo.png" class="logoImg" alt="logoImg"/></a>
        <span class="title">{{ title }}</span>
        <el-button size="mini" style="margin-top:10px;margin-left:10px;" @click="changeLangEvent">
          {{ $t("base.lang") }}
        </el-button>
      </el-col>
      <el-col :xs="14" :sm="14" :md="14" :lg="14">
        <el-menu :default-active="activeIndex" class="el-menu" mode="horizontal" @select="handleSelect">
          <el-menu-item v-for="(item, index) in navList" :key="'menu' + index" :index="index + ''">{{ item.name }}</el-menu-item>
        </el-menu>
      </el-col>
    </el-row>
  </div>
</template>

<script>
export default {
  name: "Header",
  data() {
    return {
      activeIndex: "0",
      lang: "zh-CN"
    };
  },
  props: {
    navList: Array,
    title: String
  },
  created() {
    const arr = this.$route.path.split("/");
    if (arr.length === 3) {
      this.navList.forEach((element, index) => {
        if (arr[2] === element.url) {
          this.activeIndex = index.toString();
        }
      });
    }
  },
  methods: {
    handleSelect(key) {
      if (this.$route.path !== `/${this.$route.path.split("/")[1]}/${this.navList[key].url}`) {
        this.$router.push({ path: `/${this.$route.path.split("/")[1]}/${this.navList[key].url}` });
      }
    },
    changeLangEvent() {
      this.$confirm("确定切换语言吗?", "提示", {
        confirmButtonText: "确定",
        cancelButtonText: "取消",
        type: "warning"
      })
        .then(() => {
          if (this.lang === "zh-CN") {
            this.lang = "en-US";
            this.$i18n.locale = this.lang; // 关键语句
          } else {
            this.lang = "zh-CN";
            this.$i18n.locale = this.lang; // 关键语句
          }
        })
        .catch(e => {
          console.log(e);
        });
    }
  }
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped lang="scss">
.s-header .el-menu {
  float: right;
}
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
    margin-right: 35px;
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
.el-dropdown {
  line-height: 60px;
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
.s-header {
  height: 50px;
  margin-top: 0;
  .logoImg {
    margin: 5px 0;
    height: 40px;
  }
  .title {
    margin: 5px 0;
    border-left: 0;
    font-size: 20px;
  }
}
</style>
