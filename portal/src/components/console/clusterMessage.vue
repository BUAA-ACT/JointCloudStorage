<template>
  <el-card class="chartCard" style="height: 150px">
    <el-table style="padding:20px 0" :data="clusterMessage">
      <el-table-column prop="agent" label="服务商" width="140"> </el-table-column>
      <el-table-column prop="version" label="KubeX版本" width="140"> </el-table-column>
      <el-table-column prop="visible" label="集群可见性" width="140"> </el-table-column>
    </el-table>
  </el-card>
</template>

<script>
export default {
  name: "ClusterMessage",
  data() {
    return {
      clusterMessage: [{ agent: "KubeX", version: "v0.1", visible: "部分可见" }]
    };
  }
};
</script>

<style lang="scss">
.el-table::before {
  height: 0;
}
</style>
