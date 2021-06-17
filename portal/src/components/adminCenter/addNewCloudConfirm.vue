<template>
  <div>
    <el-row>
      <el-col>
        <el-card shadow="always" class="requirementForm kuberx">
          <el-form ref="form" :model="form" label-width="100px" label-position="left">
            <el-form-item label="存储服务商">
              <el-radio-group v-model="form.resource">
                <el-radio-button label="阿里"></el-radio-button>
                <el-radio-button label="华为"></el-radio-button>
                <el-radio-button label="腾讯"></el-radio-button>
                <el-radio-button label="百度"></el-radio-button>
                <el-radio-button label="金山"></el-radio-button>
              </el-radio-group>
            </el-form-item>
            <el-form-item label="云存储名称">
              <el-input v-model="form.cloudName" class="input"></el-input>
            </el-form-item>
            <el-form-item label="云际 id">
              <el-input v-model="form.cloudId" class="input"></el-input>
            </el-form-item>
            <el-divider></el-divider>
            <el-form-item label="存储价格">
              <el-input v-model="form.storagePrice" class="input"></el-input>
              元/GB
            </el-form-item>
            <el-form-item label="流量价格">
              <el-input v-model="form.trafficPrice" class="input"></el-input>
              元/GB
            </el-form-item>
            <el-form-item label="可用性">
              <el-radio-group v-model="form.availability">
                <el-radio-button label="0.99995"></el-radio-button>
                <el-radio-button label="0.9995"></el-radio-button>
                <el-radio-button label="0.995"></el-radio-button>
              </el-radio-group>
            </el-form-item>
            <el-divider></el-divider>
            <el-form-item label="存储接入点">
              https://
              <el-input v-model="form.endpoint" class="input"></el-input>
            </el-form-item>
            <el-form-item label="accessKey">
              <el-input v-model="form.accessKey" class="input"></el-input>
            </el-form-item>
            <el-form-item label="secretKey">
              <el-input v-model="form.secretKey" class="input"></el-input>
            </el-form-item>
            <el-form-item label="bucket 名称">
              <el-input v-model="form.bucket" class="input"></el-input>
            </el-form-item>
            <el-divider></el-divider>
            <el-form-item label="云际地址">
              https://
              <el-input v-model="form.address" class="input"></el-input>
            </el-form-item>
            <el-form-item label="地理位置">
            </el-form-item>
            <el-button type="primary">点击提交</el-button>
            <el-button type="info" plain v-on:click="goBack"> 返回修改 </el-button>
          </el-form>
        </el-card>
      </el-col>
    </el-row>
    <el-row>
    </el-row>
  </div>
</template>

<script>
export default {
  name: "addNewCloudConfirm",
  data() {
    return {
      form: null
    };
  },
  created() {
    this.form = this.$route.params.formData;
  },
  methods: {
    goBack() {
      this.$router.push({ name: "addNewCloud", params: { formData: this.form } });
    }
  }
};
</script>

<style scoped>

</style>
