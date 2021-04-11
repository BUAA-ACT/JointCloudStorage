<template>
  <div>
    <div class="mBanner mBanner-m5">
      <p class="bannerTitle">
        {{ $t("disasterReco.title") }}
        <br />
        <span>{{ $t("disasterReco.bannerText") }}</span>
      </p>
    </div>
    <div class="manufacturer">
      <h2>{{ $t("disasterReco.technicalSpot") }}</h2>
      <el-divider></el-divider>
      <div class="text" v-html="$t('disasterReco.technicalSpotHtml')"></div>
    </div>
  </div>
</template>
<script>
export default {
  name: "technicalSpot",
  data() {
    return {};
  },
  watch: {
    evaluationTime() {
      // console.log(e);
    }
  },
  methods: {
    handleSelectionChange(e) {
      // console.log(e);
      this.selectedList = [...e];
    },
    compare() {
      if (this.selectedList.length < 2) {
        this.$message({
          message: "请选择至少两条进行对比",
          type: "warning"
        });
      } else {
        // console.log('compare');
      }
    }
  }
};
</script>
<style lang="scss">
.text {
  p {
    text-indent: 2em;
    line-height: 2em;
    -webkit-margin-before: 0em;
    margin-block-start: 0em;
    -webkit-margin-after: 0em;
    margin-block-end: 0em;
  }
}
</style>
