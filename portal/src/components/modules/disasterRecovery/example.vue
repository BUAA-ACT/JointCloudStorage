<template>
  <div>
    <div class="mBanner mBanner-m5">
      <p class="bannerTitle">
        {{ $t("disasterReco.title") }}
        <br />
        <span>{{ $t("disasterReco.bannerText") }}</span>
      </p>
    </div>
    <h2>{{ $t("disasterReco.title") }}</h2>
    <div class="manufacturer disasterRecovery">
      <el-steps :active="activeNum" simple style="margin-bottom:30px">
        <el-step :title="$t('disasterReco.tab1')" icon="el-icon-thumb"></el-step>
        <el-step :title="$t('disasterReco.tab2')" icon="el-icon-cloudy"></el-step>
        <el-step :title="$t('disasterReco.tab3')" icon="el-icon-circle-check"></el-step>
        <el-step :title="$t('disasterReco.tab4')" icon="el-icon-document"></el-step>
      </el-steps>
      <router-view></router-view>
    </div>
  </div>
</template>
<script>
export default {
  name: "Example",
  components: {},
  data() {
    return {
      // activeNum: 0,
    };
  },
  computed: {
    activeNum() {
      const url = this.$route.path.toString().split("/disasterRecovery/example")[1];
      switch (url) {
        case "/selectSolution":
          return 1;
        case "/caseConfirm":
          return 2;
        case "/fileView":
          return 3;
        default:
          return 0;
      }
    }
  },
  methods: {}
};
</script>
<style lang="scss">
.disasterRecovery .el-table {
  margin: 30px 0;
}
</style>
