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
    <div class="manufacturer">
      <p class="text">{{ $t("disasterReco.intro") }}</p>
      <img src="@/assets/modules/cloudBookkeeping/img.png" alt="" />
    </div>
  </div>
</template>

<script>
export default {
  name: "overView",
  data() {
    return {};
  },
  methods: {}
};
</script>

<style lang="scss" scoped>
.text {
  font-weight: normal;
  line-height: 2rem;
  font-size: 14px;
  color: #000000;
  text-indent: 2rem;
}
.manufacturer img {
  width: 100%;
  margin: auto;
}
</style>
