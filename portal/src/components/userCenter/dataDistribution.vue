<template>
  <div>
    <location-viewer class="location-viewer" :clouds="cloudStat" :format-function="formatClouds"></location-viewer>
  </div>
</template>

<script>
import locationViewer from "@/components/viewer/locationViewer.vue";
import Utils from "@/utils/other";

export default {
  name: "dataDistribution",
  components: {
    locationViewer
  },
  data() {
    return {
      cloudUsage: {}
    };
  },
  computed: {
    cloudStat() {
      return this.$store.getters.dataStats;
    }
  },
  methods: {
    formatClouds: Clouds => {
      const { Volume } = this.cloudStat;
      return Clouds.map(value => {
        return {
          name: value.CloudID,
          value: value.Location.split(",").concat([
            `存储用量：${Utils.formatBytes(Volume)}<br/>
             上传流量：${Utils.formatBytes(value.UploadTraffic)}<br/>
             下载流量：${Utils.formatBytes(value.DownloadTraffic)}<br/>`
          ])
        };
      });
    }
  }
};
</script>

<style scoped>
.location-viewer {
  width: 50vw;
  height: 50vh;
}
</style>
