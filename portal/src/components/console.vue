<template>
  <el-container>
    <el-aside width="200px">
      <el-menu default-active="1-1-1" class="el-menu-vertical-demo" @select="handleSelect">
        <el-submenu index="1">
          <template slot="title">
            <i class="el-icon-menu"></i>
            <span>工作台</span>
          </template>
          <el-submenu index="1-1">
            <template slot="title">集群管理</template>
            <el-menu-item index="1-1-1">概览</el-menu-item>
            <el-menu-item index="1-1-2">节点管理</el-menu-item>
            <el-menu-item index="1-1-3">项目管理</el-menu-item>
            <el-menu-item index="1-1-4">存储管理</el-menu-item>
          </el-submenu>
        </el-submenu>
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
        "1-1-1": "overview",
        "1-1-2": "nodeManagement",
        "1-1-3": "projectManagement",
        "1-1-4": "storageManagement"
      }
    };
  },
  methods: {
    handleSelect(key) {
      console.log(key);
      if (this.$route.path !== this.navList[key]) {
        this.$router.push({ path: `/console/${this.navList[key]}` });
      }
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
