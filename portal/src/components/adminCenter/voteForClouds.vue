<template>
  <div>
    <el-card>
      <cloud-list :clouds="pendingClouds">
        <template #operations="props">
          <el-button type="success" @click="accept(props.cloud.CloudID)">同意</el-button>
          <el-button type="warning" @click="reject(props.cloud.CloudID)">反对</el-button>
        </template>
      </cloud-list>
    </el-card>
  </div>
</template>

<script>
import Clouds from "@/api/clouds";
import CloudList from "./cloudList.vue";

export default {
  name: "voteForClouds",
  components: {
    CloudList
  },
  data() {
    return {
      BMap: null,
      pendingClouds: [
        {
          CloudID: "test",
          Endpoint: "localhost",
          AccessKey: "asdjfghkasjdhfkljahs",
          SecretKey: "ajksdhjflkjashdlfjkh",
          StoragePrice: 0.01,
          TrafficPrice: 0.01,
          Availability: 0.995,
          Status: "DOWN",
          Location: "114.51419,19.810",
          Address: "localhost",
          CloudName: "新主楼G513",
          ProviderName: "阿里",
          Bucket: ""
        }
      ],
      cloud2Cities: {},
      status2Type: {
        UP: "success",
        DOWN: "danger"
      }
    };
  },
  methods: {
    accept(cloudID) {
      this.vote(cloudID, true);
    },
    reject(cloudID) {
      this.vote(cloudID, false);
    },
    /**
     * Vote for Cloud
     * @param {string} cloudID
     * @param {Boolean} option
     * @returns {Promise<void>}
     */
    async vote(cloudID, option) {
      Clouds.voteForCloud(cloudID, option)
        .then(resp => {
          if (resp) {
            this.$message.success("投票成功！");
          }
        })
        .catch(() => {});
    },
    /**
     * getPendingClouds
     * @see getVoteRequests
     */
    async getPendingClouds() {
      Clouds.getVoteRequests().then(resp => {
        if (resp) {
          this.pendingClouds = resp.Clouds || [];
        }
      });
    }
  },
  beforeMount() {
    // this.getPendingClouds();
  }
};
</script>

<style scoped lang="scss">
.el-button--success {
  color: #fff;
  background-color: #67c23a;
  border-color: #67c23a;
}

.el-button--warning {
}
</style>
