<template>
  <div class="confirm-container">
    <el-row>
      <el-col>
        <el-form ref="form" :model="cloud" label-width="100px" label-position="left">
          <el-form-item label="存储服务商">
            <el-input :value="cloud.resource" class="read-only"> </el-input>
          </el-form-item>
          <el-form-item label="云存储名称">
            <el-input :value="cloud.cloudName" class="read-only"></el-input>
          </el-form-item>
          <el-form-item label="云际 id">
            <el-input :value="cloud.cloudId" class="read-only"></el-input>
          </el-form-item>
          <el-divider></el-divider>
          <el-form-item label="存储价格">
            <el-input :value="cloud.storagePrice" class="read-only"></el-input>
            元/GB
          </el-form-item>
          <el-form-item label="流量价格">
            <el-input :value="cloud.trafficPrice" class="read-only"></el-input>
            元/GB
          </el-form-item>
          <el-form-item label="可用性">
            <el-input :value="cloud.availability" class="read-only"> </el-input>
          </el-form-item>
          <el-divider></el-divider>
          <el-form-item label="存储接入点">
            https://
            <el-input :value="cloud.endpoint" class="read-only"></el-input>
          </el-form-item>
          <el-form-item label="accessKey">
            <el-input :value="cloud.accessKey" class="read-only"></el-input>
          </el-form-item>
          <el-form-item label="secretKey">
            <el-input :value="cloud.secretKey" class="read-only"></el-input>
          </el-form-item>
          <el-form-item label="bucket 名称">
            <el-input :value="cloud.bucket" class="read-only"></el-input>
          </el-form-item>
          <el-divider></el-divider>
          <el-form-item label="云际地址">
            https://
            <el-input :value="cloud.address" class="read-only"></el-input>
          </el-form-item>
          <el-form-item label="地理位置">
            <el-input :value="cloud.location" class="read-only"></el-input>
          </el-form-item>
          <el-button type="primary" @click="submit" :loading="loading">点击提交</el-button>
          <el-button type="info" plain @click="goBack"> 返回修改 </el-button>
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
          storagePrice: "0.01",
          cloudName: "",
          cloudId: "",
          location: "",
          resource: "阿里",
          trafficPrice: "0.01",
          availability: "0.9995",
          endPoint: "",
          accessKey: "",
          secretKey: "",
          bucket: "",
          address: ""
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
