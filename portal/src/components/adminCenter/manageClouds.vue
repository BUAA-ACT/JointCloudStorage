<template>
  <div>
    <cloud-list :clouds="allClouds">
      <template #operations="props">
        <el-button type="primary" @click="showChangeCloud(props)">修改云信息</el-button>
      </template>
    </cloud-list>
    <el-dialog title="修改云信息" :visible.sync="changeCloudDiagVis"
      ><add-new-cloud :cloud="modifyingCloud" @success="hideChangeCloud" @cancel="hideChangeCloud"
    /></el-dialog>
  </div>
</template>

<script>
import CloudList from "@/components/adminCenter/cloudList.vue";
import Clouds from "@/api/clouds";
import AddNewCloud from "@/components/adminCenter/addNewCloud.vue";

export default {
  name: "manageClouds",
  components: { AddNewCloud, CloudList },
  data() {
    return {
      allClouds: [],
      modifyingCloud: {},
      changeCloudDiagVis: false
    };
  },
  methods: {
    async getAllClouds() {
      Clouds.getAllClouds()
        .then(resp => {
          if (resp) {
            this.allClouds = resp.Clouds || [];
          }
        })
        .catch(() => {});
    },
    hideChangeCloud() {
      this.changeCloudDiagVis = false;
      this.modifyingCloud = {};
    },
    showChangeCloud(cloud) {
      this.modifyingCloud = cloud;
      this.changeCloudDiagVis = true;
    }
  },
  beforeMount() {
    this.getAllClouds();
  }
};
</script>

<style scoped></style>
