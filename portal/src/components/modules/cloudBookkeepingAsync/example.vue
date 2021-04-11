<template>
  <div>
    <div class="mBanner mBanner-m6">
      <p class="bannerTitle">
        云际异步记账
        <br />
        <span>安全、灵活、高吞吐量和低延迟相结合</span>
      </p>
    </div>
    <h2>云际异步记账</h2>
    <div class="manufacturer cloudComputation">
      <iframe src="http://106.53.150.192:18888/" style="border: 0"></iframe>
    </div>
  </div>
</template>
<script>
export default {
  name: "example",
  components: {},
  data() {
    return {
      interval: "",
      activeName: "0",
      form: {
        env: "python"
      },
      args: "{}",
      result: "",
      percentage: 0,
      second: 0
    };
  },
  methods: {}
};
</script>
<style lang="scss">
iframe {
  width: 100%;
  height: 500px;
  margin-top: 30px;
}
</style>
