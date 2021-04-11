<template>
  <div class="left" style="width: 500px">
    <el-form ref="form" :model="form" label-width="130px">
      <el-form-item label="供应商数目">
        <el-input-number :min="1" v-model="form.Vendor"></el-input-number>
      </el-form-item>
      <el-form-item label="存储成本限制">
        <el-input-number :min="0.0" :controls="false" v-model="form.StoragePrice" :precision="2" :step="0.1"></el-input-number>元/GB/月
      </el-form-item>
      <el-form-item label="流量成本限制">
        <el-input-number :min="0.0" :controls="false" v-model="form.TrafficPrice" :precision="2" :step="0.1"></el-input-number>元/GB
      </el-form-item>
      <el-form-item label="可用性要求">
        <el-input-number :min="0.0" :controls="false" v-model="form.Availability" :precision="2" :step="0.1" :max="1"></el-input-number>
      </el-form-item>
      <el-form-item label="允许存在下载延迟">
        <el-switch v-model="form.AllowDelay"> </el-switch>
      </el-form-item>
      <el-form-item>
        <el-button type="primary" @click="onSubmit" :loading="loading">更新个人偏好</el-button>
        <el-button>取消</el-button>
      </el-form-item>
    </el-form>
  </div>
</template>

<script>
import cloudStorage from "@/api/cloudStorage";

export default {
  data() {
    return {
      form: {
        Vendor: 1,
        StoragePrice: 1.0,
        TrafficPrice: 1.0,
        Availability: 0.5,
        AllowDelay: true,
        Latency: { beijing: 20 }
      },
      loading: false
    };
  },
  methods: {
    onSubmit() {
      this.loading = true;
      cloudStorage
        .updatePreference(this.form)
        .then(() => {
          this.$message.success("个人偏好更新成功");
          this.loading = false;
        })
        .catch(() => {
          this.loading = false;
        });
    }
  }
};
</script>

<style scoped>
.left {
  text-align: left;
}
</style>
