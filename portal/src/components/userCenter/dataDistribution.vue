<template>
  <div>
    <location-viewer
      class="location-viewer"
      :clouds="cloudStat.cloudsDetails"
      :format-function="formatClouds"
      :inactive-clouds="inactiveClouds"
      :format-in-active-cloud-function="formatInactiveClouds"
    >
    </location-viewer>
  </div>
</template>

<script>
import locationViewer from "@/components/viewer/locationViewer.vue";
import Utils from "@/utils/other";
import Clouds from "@/api/clouds";

export default {
  name: "dataDistribution",
  components: {
    locationViewer
  },
  data() {
    return {
      cloudUsage: {},
      inactiveClouds: []
    };
  },
  computed: {
    cloudStat() {
      return this.$store.getters.dataStats;
    }
  },
  methods: {
    formatClouds(clouds) {
      const { Volume } = this.cloudStat;
      return clouds.map(value => {
        return {
          name: value.CloudID,
          value: value.Location.split(",").concat([
            `${value.CloudID}<br />
             存储用量：${Utils.formatBytes(Volume)}<br/>
             上传流量：${Utils.formatBytes(value.UploadTraffic)}<br/>
             下载流量：${Utils.formatBytes(value.DownloadTraffic)}<br/>`
          ])
        };
      });
    },
    formatInactiveClouds(clouds) {
      return clouds
        .filter(item => {
          this.$log(this.clouds);
          // eslint-disable-next-line no-restricted-syntax
          return !this.clouds.some(value => {
            return value.CloudID === item.CloudID;
          });
        })
        .map(value => {
          return {
            name: value.CloudID,
            value: value.Location.split(",") // longitude ,latitude
              .concat([
                `存储价格：${value.StoragePrice}元/GB/月<br/>
          流量价格：${value.TrafficPrice}元/GB<br/>
          可用性：${value.Availability}<br />`
              ])
          };
        });
    },
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
    }
  },
  async mounted() {
    await this.$store.dispatch("getInfo");
  }
};
</script>

<style scoped>
.location-viewer {
  width: 50vw;
  height: 50vh;
}
</style>
