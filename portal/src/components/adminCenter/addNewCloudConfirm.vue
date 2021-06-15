<template>
  <div>
    <el-row>
      <el-col>
        <h2>请确认待提交信息</h2>
        {{ formData }}
      </el-col>
    </el-row>
    <el-row>
      <el-button type="primary">点击提交</el-button>
      <el-button type="info" plain v-on:click="goBack"> 返回修改 </el-button>
    </el-row>
  </div>
</template>

<script>
export default {
  name: "addNewCloudConfirm",
  data() {
    return {
      formData: null
    };
  },
  created() {
    this.formData = this.$route.params.formData;
  },
  methods: {
    goBack() {
      this.$router.push({ name: "addNewCloud", params: { formData: this.formData } });
    }
  }
};
</script>

<style scoped>

</style>
