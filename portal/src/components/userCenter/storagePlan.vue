<template>
  <div style="text-align: left">
    <el-button @click="getStoragePlans">获取存储方案</el-button>
    <el-button type="primary" @click="submit" :loading="submitLoading" :disabled="!plansLoaded">提交</el-button><br />
    <div v-if="plansLoaded" class="plans-viewer-container">
      <el-card class="box-card">
        <div slot="header" class="clearfix">
          <el-radio v-model="storagePlanIndex" :label="0">存储价格优先</el-radio>
        </div>
        <div class="text item">
          存储模式： {{ modifyStorageMode(storagePlans.StoragePriceFirst) }}<br />
          存储价格： {{ storagePlans.StoragePriceFirst.StoragePrice }}<br />
          流量价格： {{ storagePlans.StoragePriceFirst.TrafficPrice }}<br />
          可用性：{{ storagePlans.StoragePriceFirst.Availability }}
        </div>
      </el-card>
      <el-card class="box-card">
        <div slot="header" class="clearfix">
          <el-radio v-model="storagePlanIndex" :label="1">流量价格优先</el-radio>
        </div>
        <div class="text item">
          存储模式： {{ modifyStorageMode(storagePlans.TrafficPriceFirst) }}<br />
          存储价格： {{ storagePlans.TrafficPriceFirst.StoragePrice }}<br />
          流量价格： {{ storagePlans.TrafficPriceFirst.TrafficPrice }}<br />
          可用性：{{ storagePlans.TrafficPriceFirst.Availability }}
        </div>
      </el-card>
    </div>
    <location-viewer v-if="plansLoaded" :clouds="candidates[storagePlanIndex].Clouds" :inactive-clouds="inactiveClouds" class="location-viewer" />
  </div>
</template>

<script>
import Plan from "@/api/plan";
import Clouds from "@/api/clouds";
import locationViewer from "@/components/viewer/locationViewer.vue";

export default {
  name: "storagePlan",
  components: {
    locationViewer
  },
  data() {
    return {
      storagePlanIndex: 0,
      storagePlanPredefine: ["StoragePriceFirst", "TrafficPriceFirst"],
      // { "ID": "txyun-chongqing", "Location": "116.381252,39.906569" }
      candidates: [],
      storagePlans: {
        StoragePriceFirst: {
          StorageMode: "Replica",
          N: 1,
          K: 1,
          Clouds: [],
          StoragePrice: 0.12,
          TrafficPrice: 0.5,
          Availability: 0.99995
        },
        TrafficPriceFirst: {
          StorageMode: "Replica",
          N: 2,
          K: 1,
          Clouds: [],
          StoragePrice: 0.24,
          TrafficPrice: 0.4,
          Availability: 0.999999975
        }
      },
      // { name: "China", value: [104.195397, 35.86166, 550] }
      inactiveClouds: {},
      formattedClouds: [],
      plansLoaded: false,
      submitLoading: false
    };
  },
  methods: {
    modifyStorageMode(storagePlan) {
      return `${storagePlan.StorageMode}(N:${storagePlan.N}, K:${storagePlan.K})`;
    },
    async getStoragePlans() {
      this.plansLoaded = false;
      await Plan.getStoragePlans().then(resp => {
        this.storagePlans = resp;
        // Object.keys(this.storagePlans).forEach(index => {
        //   console.log(index);
        //   if (index === "StoragePriceFirst" || index === "TrafficPriceFirst") {
        //     this.storagePlans[index].StorageMode += `(N:${this.storagePlans[index].N}, K:${this.storagePlans[index].K})`;
        //   }
        // });
        const { StoragePriceFirst, TrafficPriceFirst } = resp;
        this.candidates = [StoragePriceFirst, TrafficPriceFirst];
        // { name: "China", value: [104.195397, 35.86166, Caption] }
      });
      this.candidates = this.candidates || [];
      this.formatClouds();
      this.plansLoaded = true;
      this.getAllCloud();
    },
    /**
     * 获取所有云
     *
     */ async getAllCloud() {
      this.inactiveClouds = Clouds.getAllCloud().clouds;
      this.$log(this.inactiveClouds);
    },
    /**
     * @method formatClouds
     * @summary Format resp clouds into form: [Longitude, Latitude, Caption]
     * @input clouds: {
                CloudID: "ksyun-beijing",
                Endpoint: "ks3-cn-beijing.ksyun.com",
                StoragePrice: 0.12,
                TrafficPrice: 0.4,
                Availability: 0.9995,
                Status: "UP",
                Location: "116.381252,39.906569",
                Address: "localhost:8182"
              }
     * @output clouds: [116.381252, 39.906569, StoragePrice + TrafficPrice + Availability]
     */
    formatClouds() {
      this.formattedClouds = this.candidates[this.storagePlanIndex].Clouds.map(value => {
        return {
          name: value.CloudID,
          value: value.Location.split(",").concat([
            `存储价格：${value.StoragePrice}元/GB/月<br/>
          流量价格：${value.TrafficPrice}元/GB<br/>
          可用性：${value.Availability}<br />`
          ])
        };
      });
    },
    async submit() {
      this.submitLoading = true;
      Plan.changeStoragePlan(this.candidates[this.storagePlanIndex])
        .then(resp => {
          if (resp) this.$message.success("更新存储方案成功！");
          this.submitLoading = false;
        })
        .catch(() => {
          this.submitLoading = false;
        });
    }
  },
  watch: {
    storagePlanIndex() {
      // this.formatClouds();
    }
  }
};
</script>

<style scoped>
.text {
  font-size: 14px;
}

.item {
  margin-bottom: 18px;
}

.clearfix:before,
.clearfix:after {
  display: table;
  content: "";
}
.clearfix:after {
  clear: both;
}

.box-card {
  width: 250px;
  display: inline-block;
  margin: 0 100px;
}
.plans-viewer-container {
  display: -webkit-box;
  display: -webkit-flex;
  display: -ms-flexbox;
  display: flex;
  margin-top: 30px;
  -webkit-box-pack: justify;
  -webkit-justify-content: space-between;
  -ms-flex-pack: justify;
  justify-content: space-between;
  width: 50vw;
  min-width: 800px;
}
.location-viewer {
  width: 50vw;
  min-width: 800px;
  height: 400px;
}
</style>
