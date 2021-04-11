<template>
  <div class="requirement">
    <div class="sBanner sBanner-s2">
      <p class="bannerTitle">
        需求发布
        <br />
        <span>免费发布供求信息，国内大型的云服务平台</span>
      </p>
    </div>
    <div class="bg-gray">
      <el-tabs class="dataCenter" v-model="activeName">
        <el-tab-pane label="自定义需求发布" name="first"><SingleRequirement /></el-tab-pane>
        <el-tab-pane label="综合需求推荐" name="second"><ComprehensiveRequirement /></el-tab-pane>
        <el-tab-pane label="数据需求发布" name="third"><DataRequirement /></el-tab-pane>
      </el-tabs>
    </div>
  </div>
</template>
<script>
import SingleRequirement from "./singleRequirement.vue";
import ComprehensiveRequirement from "./comprehensiveRequirement.vue";
import DataRequirement from "./dataRequirement.vue";

export default {
  name: "RequirementForm",
  data() {
    return {
      form: {},
      marks: {
        0: "0Mbps",
        30: "30Mbps",
        120: "120Mbps",
        200: "200Mbps"
      },
      activeName: "first"
    };
  },
  components: { SingleRequirement, ComprehensiveRequirement, DataRequirement },
  methods: {
    onSubmit() {}
  }
};
</script>
<style lang="scss">
.requirement .dataCenter {
  width: 86%;
  max-width: 1200px;
  margin: 0 auto;
  margin-top: 20px;
  .el-tabs__header {
    padding: 0 20px;
    margin: 0;
  }
}
</style>
<style lang="scss" scoped>
.bg-gray {
  padding-bottom: 50px;
}
</style>
