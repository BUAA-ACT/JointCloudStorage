<template>
  <div class="manufacturer">
    <el-card shadow="always" class="requirementForm">
      <p class="indexTitle">数据需求发布</p>
      <el-form ref="form" :model="form" label-width="100px" label-position="left">
        <el-form-item label="数据名称">
          <el-input v-model="form.desc"></el-input>
        </el-form-item>
        <el-form-item label="数据需求类型">
          <el-select v-model="form.region" placeholder="请选择数据需求类型">
            <el-option label="需求类型" value="shanghai"></el-option>
            <el-option label="需求类型2" value="beijing"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="数据需求描述">
          <el-input type="textarea" v-model="form.desc"></el-input>
        </el-form-item>
        <el-form-item label="行业类型">
          <el-radio-group v-model="form.resource">
            <el-radio label="地产"></el-radio>
            <el-radio label="电子商务"></el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="对接联系人">
          <el-input v-model="form.desc"></el-input>
        </el-form-item>
        <el-form-item label="对接联系方式">
          <el-input v-model="form.desc"></el-input>
        </el-form-item>
      </el-form>
      <el-button class="submitBtn" type="primary" @click="onSubmit">确认提交</el-button>
    </el-card>
  </div>
</template>

<script>
export default {
  data() {
    return {
      form: {}
    };
  },
  methods: {
    onSubmit() {
      this.$router.push({ path: "/requirementReleaseSuccess" });
    }
  }
};
</script>

<style lang="scss" scoped></style>
