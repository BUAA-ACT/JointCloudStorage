<template>
  <div>
    <div class="mBanner releaseSuccessBanner">
      <p class="bannerTitle">
        需求发布成功
        <br />
        <span>云际计算平台根据您发布的参数，智能计算出以下方案</span>
        <br />
        <el-button class="releaseRequirementBtn" size="small" round @click="$router.push('/requirementList')">返回需求目录</el-button>
        <el-button class="releaseRequirementBtn" size="small" round @click="$router.push('/requirementRelease')">再次发布需求</el-button>
      </p>
    </div>
    <div class="manufacturer">
      <RecomList />
    </div>
  </div>
</template>

<script>
import RecomList from "@/components/requirement/recomList.vue";

export default {
  data() {
    return {};
  },
  components: { RecomList }
};
</script>

<style lang="scss" scoped>
.releaseSuccessBanner {
  .bannerTitle {
    text-align: center;
    padding-top: 60px;
    line-height: 4rem;
  }
}
</style>
