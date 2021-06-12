<template>
  <el-container id="app">
    <el-main class="main">
      <HelloWorld msg="云际存储"/>
      <div class="blank"></div>
      <JcsPanIndex></JcsPanIndex>
      <div class="blank"></div>
    </el-main>
    <el-footer class="footer center">
      国家重点研发计划——云际计算项目组    版权所有 &copy; 北京航空航天大学 <br>
      <el-link href="http://beian.miit.gov.cn/">ICP备案号：京ICP备2021011189号-1</el-link>
    </el-footer>
  </el-container>
</template>

<script>
import HelloWorld from './components/HelloWorld.vue'
import JcsPanIndex from "@/JcsPanIndex";
export default {
  name: 'App',
  components: {
    JcsPanIndex,
    HelloWorld
  },
}
</script>

<style scoped>
.blank {
  margin-top: 50px;
}
.main {
  min-height: calc(100vh - 100px);
}
.center {
  text-align: center;
}
.footer {
  background: azure;
  height: 100px;
}

</style>
