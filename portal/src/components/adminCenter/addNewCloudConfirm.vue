<template>
  <div class="confirm-container">
    <el-row>
      <el-col>
        <el-form ref="form" :model="cloud" label-width="100px" label-position="left">
          <el-form-item label="存储服务商">
            <el-input :value="cloud.ProviderName" class="read-only"> </el-input>
          </el-form-item>
          <el-form-item label="云存储名称">
            <el-input :value="cloud.CloudName" class="read-only"></el-input>
          </el-form-item>
          <el-form-item label="云际 id">
            <el-input :value="cloud.CloudID" class="read-only"></el-input>
          </el-form-item>
          <el-divider></el-divider>
          <el-form-item label="存储价格">
            <el-input :value="cloud.StoragePrice" class="read-only"></el-input>
            元/GB
          </el-form-item>
          <el-form-item label="流量价格">
            <el-input :value="cloud.TrafficPrice" class="read-only"></el-input>
            元/GB
          </el-form-item>
          <el-form-item label="可用性">
            <el-input :value="cloud.Availability" class="read-only"> </el-input>
          </el-form-item>
          <el-divider></el-divider>
          <el-form-item label="存储接入点">
            https://
            <el-input :value="cloud.Endpoint" class="read-only"></el-input>
          </el-form-item>
          <el-form-item label="Access Key">
            <el-input :value="cloud.AccessKey" class="read-only"></el-input>
          </el-form-item>
          <el-form-item label="Secret Key">
            <el-input :value="cloud.SecretKey" class="read-only"></el-input>
          </el-form-item>
          <el-form-item label="Bucket 名称">
            <el-input :value="cloud.Bucket" class="read-only"></el-input>
          </el-form-item>
          <el-divider></el-divider>
          <el-form-item label="云际地址">
            https://
            <el-input :value="cloud.Address" class="read-only"></el-input>
          </el-form-item>
          <el-form-item label="地理位置">
            <el-input :value="cloud.Location" class="read-only"></el-input>
          </el-form-item>
          <el-button type="primary" @click="submit" :loading="loading">点击提交</el-button>
          <el-button type="info" plain @click="goBack">返回修改</el-button>
        </el-form>
      </el-col>
    </el-row>
    <el-row> </el-row>
  </div>
</template>

<script>
import Clouds from "@/api/clouds";

export default {
  name: "addNewCloudConfirm",
  props: {
    cloud: {
      type: Object,
      required: true,
      default() {
        return {
          CloudID: "",
          Endpoint: "",
          AccessKey: "",
          SecretKey: "",
          StoragePrice: 0.01,
          TrafficPrice: 0.01,
          Availability: 0.995,
          Status: "DOWN",
          Location: "",
          Address: "",
          CloudName: "",
          ProviderName: "阿里",
          Bucket: ""
        };
      }
    }
  },
  data() {
    return {
      loading: false
    };
  },
  methods: {
    goBack() {
      this.$emit("cancel");
    },
    async submit() {
      this.loading = true;
      Clouds.addNewCloud(this.cloud)
        .then(resp => {
          if (resp) {
            this.$notify.success("成功添加云！请等待投票……");
            this.$emit("success");
          } else {
            this.$emit("fail");
          }
          this.loading = false;
        })
        .catch(() => {
          this.loading = false;
        });
    }
  }
};
</script>

<style scoped lang="scss">
.confirm-container {
  width: 25vw;
}
.read-only {
  width: 200px;
  /deep/ .el-input__inner {
    border: none;
    padding: 0;
  }
}
</style>
