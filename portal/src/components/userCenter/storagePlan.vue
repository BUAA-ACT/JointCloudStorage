<template>
  <div style="text-align: left">
    <div class="storage-plan-header">
      <el-button @click="getStoragePlans">{{ `${havePlan ? "刷新" : "获取"}存储方案` }}</el-button>
      <el-button type="primary" @click="submit" :loading="submitLoading" :disabled="!plansLoaded" v-if="!havePlan">提交</el-button><br />
    </div>
    <div v-if="plansLoaded" class="plans-viewer-container">
      <div v-if="!havePlan" class="plans-selector">
        <el-row type="flex" class="row-bg" justify="space-around">
          <el-col :span="6">
            <el-card class="box-card">
              <div slot="header" class="clearfix">
                <el-radio v-model="storagePlanIndex" :label="0">存储价格优先</el-radio>
              </div>
              <div class="text item">
                存储模式： {{ modifyStorageMode(storagePlans.StoragePriceFirst) }}<br />
                存储价格： {{ formatPrice(storagePlans.StoragePriceFirst.StoragePrice) }}<br />
                流量价格： {{ formatPrice(storagePlans.StoragePriceFirst.TrafficPrice) }}<br />
                可用性：{{ storagePlans.StoragePriceFirst.Availability.toFixed(8) * 100 }}%
              </div>
            </el-card>
          </el-col>
          <el-col :span="6">
            <el-card class="box-card">
              <div slot="header" class="clearfix">
                <el-radio v-model="storagePlanIndex" :label="1">流量价格优先</el-radio>
              </div>
              <div class="text item">
                存储模式： {{ modifyStorageMode(storagePlans.TrafficPriceFirst) }}<br />
                存储价格： {{ formatPrice(storagePlans.TrafficPriceFirst.StoragePrice) }}<br />
                流量价格： {{ formatPrice(storagePlans.TrafficPriceFirst.TrafficPrice) }}<br />
                可用性：{{ storagePlans.TrafficPriceFirst.Availability.toFixed(8) * 100 }}%
              </div>
            </el-card>
          </el-col>
          <el-col :span="6">
            <el-card class="box-card">
              <div slot="header" class="clearfix">
                <el-radio v-model="storagePlanIndex" :label="2">自定义存储方案</el-radio>
              </div>
              <div class="text item">
                <el-button type="primary" :disabled="storagePlanIndex !== 2" @click="startCustomize">点击这里开始自定义</el-button>
              </div>
            </el-card>
          </el-col>
        </el-row>
      </div>
      <div v-else>
        <plan-detail />
      </div>
    </div>
    <location-viewer
      v-if="plansLoaded && storagePlanIndex !== 2"
      :clouds="candidates[storagePlanIndex].Clouds || []"
      :inactive-clouds="inactiveClouds"
      class="location-viewer"
    />
    <el-dialog :visible.sync="customizing">
      <customize-storage-plan @success="customSuccess" @failed="customizing = false" />
    </el-dialog>
  </div>
</template>

<script>
import Plan from "@/api/plan";
import Clouds from "@/api/clouds";
import locationViewer from "@/components/viewer/locationViewer.vue";
import CustomizeStoragePlan from "@/components/userCenter/customizeStoragePlan.vue";
import PlanDetail from "@/components/userCenter/planDetail.vue";

export default {
  name: "storagePlan",
  components: {
    PlanDetail,
    CustomizeStoragePlan,
    locationViewer
  },
  inject: ["formatPrice", "reload"],
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
      inactiveClouds: [],
      formattedClouds: [],
      plansLoaded: false,
      submitLoading: false,
      customizing: false
    };
  },
  methods: {
    modifyStorageMode(storagePlan) {
      return `${storagePlan.StorageMode}(N:${storagePlan.N}, K:${storagePlan.K})`;
    },
    async getStoragePlans() {
      this.plansLoaded = false;
      if (this.havePlan) {
        await this.$store.dispatch("updateInfo", "StoragePlan");
        this.candidates = [this.$store.getters.storagePlan];
      } else {
        await Plan.getStoragePlans().then(resp => {
          this.storagePlans = resp;
          // Object.keys(this.storagePlans).forEach(index => {
          //   console.log(index);
          //   if (index === "StoragePriceFirst" || index === "TrafficPriceFirst") {
          //     this.storagePlans[index].StorageMode += `(N:${this.storagePlans[index].N}, K:${this.storagePlans[index].K})`;
          //   }
          // });
          this.$log(resp);
          const { StoragePriceFirst, TrafficPriceFirst } = resp;
          StoragePriceFirst.Clouds = StoragePriceFirst.Clouds || [];
          TrafficPriceFirst.Clouds = TrafficPriceFirst.Clouds || [];
          this.candidates = [StoragePriceFirst, TrafficPriceFirst];
          // { name: "China", value: [104.195397, 35.86166, Caption] }
        });
      }
      this.candidates = this.candidates || [];
      this.formatClouds();
      await this.getAllCloud();
      this.plansLoaded = true;
    },
    /**
     * 获取所有云
     */
    async getAllCloud() {
      // this.inactiveClouds = Clouds.getAllCloud().clouds;
      Clouds.getAllClouds()
        .then(resp => {
          if (resp && resp.Clouds) this.inactiveClouds = resp.Clouds || [];
        })
        .catch(() => {
          this.inactiveClouds = [];
        });
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
            `存储价格：${this.formatPrice(value.StoragePrice)}元/GB/月<br/>
          流量价格：${this.formatPrice(value.TrafficPrice)}元/GB<br/>
          可用性：${value.Availability.toFixed(8) * 100}%<br />`
          ])
        };
      });
    },
    async submit() {
      this.submitLoading = true;
      Plan.changeStoragePlan(this.candidates[this.storagePlanIndex])
        .then(async resp => {
          if (resp) this.$message.success("更新存储方案成功！");
          this.submitLoading = false;
          await this.$store.dispatch("updateInfo", "StoragePlan");
          this.reload();
        })
        .catch(() => {
          this.submitLoading = false;
        });
    },
    startCustomize() {
      this.customizing = true;
    },
    customSuccess() {
      this.customizing = false;
      this.reload();
    }
  },
  watch: {
    storagePlanIndex() {
      // this.formatClouds();
    }
  },
  computed: {
    havePlan() {
      return this.$store.getters.haveStoragePlan;
    }
  },
  beforeMount() {
    this.getStoragePlans();
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
  height: 200px;
  display: inline-block;
  margin: 0 10px;
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
