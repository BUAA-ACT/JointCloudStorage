<template>
  <el-container>
    <el-aside width="200px">
      <el-menu class="el-menu-vertical-demo" @select="handleSelect" :default-active="activeEntry">
        <el-submenu index="1">
          <template slot="title">
            <i class="el-icon-user"></i>
            <span>个人中心</span>
          </template>
          <el-menu-item index="1-7">用户信息</el-menu-item>
          <el-menu-item index="1-5">用户偏好</el-menu-item>
          <el-menu-item index="1-6">存储方案</el-menu-item>
        </el-submenu>
        <el-menu-item index="2-1"><i class="el-icon-folder-opened"></i>文件管理</el-menu-item>
        <el-menu-item index="3-1"><i class="el-icon-map-location"></i>数据分布</el-menu-item>
        <el-menu-item index="4-1"><i class="el-icon-guide"></i>存储迁移</el-menu-item>
      </el-menu>
    </el-aside>
    <el-main><router-view></router-view></el-main>
  </el-container>
</template>

<script>
export default {
  components: {},
  data() {
    return {
      navList: {
        "1-1": "overview",
        "1-2": "nodeManagement",
        "1-3": "projectManagement",
        "1-4": "storageManagement",
        "1-5": "userPreference",
        "1-6": "storagePlan",
        "1-7": "userInfo",
        "2-1": "fileManagement",
        "3-1": "dataDistribution",
        "4-1": "dataMigration"
      }
    };
  },
  methods: {
    handleSelect(key) {
      console.log(key);
      console.log(this.$route.path);
      if (this.$route.path !== this.navList[key]) {
        this.$router.push({ path: `/cloudStorage/${this.navList[key]}` });
      }
    }
  },
  computed: {
    activeEntry() {
      let currentRoute;
      Object.keys(this.navList).forEach(navListKey => {
        if (this.$route.path.endsWith(this.navList[navListKey])) {
          currentRoute = navListKey;
        }
      });
      return currentRoute;
    }
  }
};
</script>
<style scoped lang="scss">
.el-menu {
  border-right: 0;
  position: fixed;
  width: 199px;
  top: 63px;
}
.el-aside {
  overflow-x: hidden;
  border-right: solid 1px #e6e6e6;
}
</style>
