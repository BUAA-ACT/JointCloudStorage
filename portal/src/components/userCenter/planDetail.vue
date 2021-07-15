<template>
  <el-card class="out-card">
    <div slot="header">
      <span>存储方案概览</span>
    </div>
    <el-form inline label-position="top">
      <el-form-item label="存储模式" class="summary">{{ storageMode }}</el-form-item>
      <el-form-item label="存储价格" class="summary">{{ storagePrice }}元/GB/月</el-form-item>
      <el-form-item label="流量价格" class="summary">{{ trafficPrice }}元/GB</el-form-item>
      <el-form-item label="可用性" class="summary">{{ availability * 100 }}%</el-form-item>
      <!--      <el-divider direction="vertical" />-->
      <div class="detail">
        <el-form-item v-if="replicaClouds" label="正在使用以下云节点存储多副本">
          <td v-for="cloud in replicaClouds" :key="cloud.CloudID">
            <el-card class="inner-card" shadow="hover">{{ cloud.CloudName }}</el-card>
          </td>
        </el-form-item>
        <el-form-item v-if="ECKClouds" label="纠删码正在使用以下云节点存储数据分块" class="detail-item">
          <td v-for="cloud in ECKClouds" :key="cloud.CloudID">
            <el-card class="inner-card" shadow="hover">{{ cloud.CloudName }}</el-card>
          </td>
        </el-form-item>
        <el-divider direction="vertical" class="splitter" v-if="ECKClouds" />
        <el-form-item v-if="ECNKClouds" label="纠删码正在使用以下云节点存储校验分块" class="detail-item">
          <td v-for="cloud in ECNKClouds" :key="cloud.CloudID">
            <el-card class="inner-card" shadow="hover">{{ cloud.CloudName }}</el-card>
          </td>
        </el-form-item>
      </div>
    </el-form>
  </el-card>
</template>

<script>
const modeMap = {
  Replica: "多副本",
  EC: "纠删码"
};

export default {
  name: "planDetail",
  data() {
    return {
      storagePlan: {}
    };
  },
  methods: {
    getStoragePlan() {
      if (!this.$store.getters.ready) {
        setTimeout(() => {
          this.getStoragePlan();
        }, 50);
        return;
      }
      this.storagePlan = this.$store.getters.storagePlan;
    }
  },
  beforeMount() {
    this.getStoragePlan();
  },
  computed: {
    storagePrice() {
      return this.storagePlan.StoragePrice;
    },
    trafficPrice() {
      return this.storagePlan.TrafficPrice;
    },
    availability() {
      return this.storagePlan.Availability.toFixed(8);
    },
    storageMode() {
      return modeMap[this.storagePlan.StorageMode];
    },
    replicaClouds() {
      return this.storagePlan.StorageMode === "Replica" && this.storagePlan.Clouds;
    },
    ECKClouds() {
      return this.storagePlan.StorageMode === "EC" && this.storagePlan.Clouds.slice(0, this.storagePlan.K);
    },
    ECNKClouds() {
      return this.storagePlan.StorageMode === "EC" && this.storagePlan.Clouds.slice(this.storagePlan.K - this.storagePlan.N);
    }
  }
};
</script>

<style scoped lang="scss">
.out-card {
  width: 1200px;
  margin-bottom: 10px;
  .summary {
    width: 150px;
  }
  .detail {
    display: inline-block;
    .splitter {
      padding-bottom: 100%; //内边距为父div高度，导致子div高度向下延伸
      margin-bottom: -100%; //外边距为父div高度，导致下方div向上挤压子div，但是由于父div的存在挡住了挤压，子div得以填充父div
    }
    .detail-item {
      width: 350px;
    }
    .inner-card {
      width: fit-content;
      display: flex;
      /deep/ .el-card__body {
        padding: 10px;
      }
    }
  }
}
</style>
