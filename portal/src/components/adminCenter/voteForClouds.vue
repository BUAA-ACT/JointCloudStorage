<template>
  <div>
    <el-card>
      <el-table :data="pendingClouds" style="width: 100%">
        <el-table-column type="expand">
          <template slot-scope="props">
            <el-form label-position="left" inline class="demo-table-expand">
              <el-form-item label="访问端点">
                <span>{{ props.row.EndPoint }}</span>
              </el-form-item>
              <el-form-item label="Access Key">
                <span>{{ props.row.AccessKey }}</span>
              </el-form-item>
              <el-form-item label="Secret Key">
                <span>{{ props.row.SecretKey }}</span>
              </el-form-item>
              <el-form-item label="地理位置">
                <span>{{ props.row.Location }}</span>
              </el-form-item>
              <el-form-item label="存储价格">
                <span>{{ props.row.StoragePrice }}</span>
              </el-form-item>
              <el-form-item label="流量价格">
                <span>{{ props.row.TrafficPrice }}</span>
              </el-form-item>
              <el-form-item label="可用性">
                <span>{{ props.row.Availability * 100 }}%</span>
              </el-form-item>
            </el-form>
          </template>
        </el-table-column>
        <el-table-column prop="CloudID" label="云节点ID" width="180"> </el-table-column>
        <el-table-column prop="CloudName" label="云节点名称" width="180"> </el-table-column>
        <el-table-column prop="ProviderName" label="供应商"> </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<script>
import Clouds from "@/api/clouds";

export default {
  name: "voteForClouds",
  data() {
    return {
      BMap: null,
      pendingClouds: []
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
    }
  }
};
</script>

<style scoped></style>
