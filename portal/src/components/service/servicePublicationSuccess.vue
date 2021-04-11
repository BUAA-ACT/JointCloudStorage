<template>
  <div>
    <div class="mBanner releaseSuccessBanner">
      <p class="bannerTitle">
        服务发布成功
        <br />
        <span>服务已成功发布至服务目录</span>
        <br />
        <el-button class="releaseRequirementBtn" size="small" round @click="$router.push('/serviceList')">返回服务目录</el-button>
        <el-button class="releaseRequirementBtn" size="small" round @click="$router.push('/servicePublication')">再次发布服务</el-button>
      </p>
    </div>
  </div>
</template>

<script>
export default {
  data() {
    return {};
  }
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
