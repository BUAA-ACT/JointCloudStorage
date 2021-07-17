<template>
  <div class="warning-container">
    <el-alert title="请勿玩火" type="error" class="warning-title" effect="dark" center :closable="false"> </el-alert>
    <i class="el-icon-remove" style="color: red; font-size: 20vh" />
    <h1>401</h1>
    <p>你没有权限！</p>
    <p>You have no permission!</p>
    <a @click="goBack">
      返回上一页<br />
      Go Back
    </a>
  </div>
</template>
<script>
export default {
  name: "Page401",
  methods: {
    goBack() {
      this.$router.back();
    }
  }
};
</script>
<style lang="scss">
.warning-container {
  padding-top: 10px;
  .warning-title {
    margin: 10px 5vw;
    width: 90vw;
  }
}
</style>
