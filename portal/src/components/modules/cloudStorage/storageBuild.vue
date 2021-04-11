<template>
  <div>
    <div class="mBanner mBanner-m4">
      <p class="bannerTitle">
        云际存储
        <br />
        <span>海量，安全，低成本，高可靠的云存储服务</span>
      </p>
    </div>
    <div class="bg-gray">
      <div class="manufacturer">
        <SelectServer v-if="!monitorDisplay" @selectedServerSumbit="onSubmit" />
        <DynamicMonitor @back="monitorDisplay = false" v-if="monitorDisplay" />
      </div>
    </div>
    <!-- <el-dialog width="800px" title="" :visible.sync="dialogFormVisible">
      <img class="gif" src="@/assets/modules/cloudStorage/ppt.gif"/>
    </el-dialog> -->
  </div>
</template>
<script>
import DynamicMonitor from "./components/dynamicMonitor.vue";
import SelectServer from "./components/selectServer.vue";

export default {
  name: "storageBuild",
  data() {
    return {
      form: {},
      monitorDisplay: false,
      dialogFormVisible: false
    };
  },
  components: {
    SelectServer,
    DynamicMonitor
  },
  methods: {
    onSubmit() {
      // this.dialogFormVisible = true;
      // setTimeout(() => {
      // this.dialogFormVisible = false;
      this.monitorDisplay = true;
      // }, 4000);
    },
    DynamicMonitoronSubmit() {}
  }
};
</script>
<style lang="scss" scoped>
.gif {
  width: 100%;
  display: block;
}
</style>
