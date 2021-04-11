<template>
  <el-card class="chartCard">
    <el-table style="padding:20px 0" :data="nodeMessage">
      <el-table-column prop="node" label="节点" width="200"> </el-table-column>
      <el-table-column prop="CPU" label="CPU使用率" width="120"> </el-table-column>
    </el-table>
  </el-card>
</template>

<script>
export default {
  name: "NodeMessage",
  data() {
    return {
      nodeMessage: [
        { node: "master 192.168.12.10", CPU: "16%" },
        { node: "node1 192.168.12.11", CPU: "12%" },
        { node: "node2 192.168.12.12", CPU: "7%" }
      ]
    };
  }
};
</script>

<style lang="scss"></style>
