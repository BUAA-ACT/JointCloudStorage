<template>
  <el-card class="chartCard">
    <p>节点详细信息</p>
    <el-table :data="nodeStatus">
      <el-table-column prop="name" label="名称" width="140"> </el-table-column>
      <el-table-column prop="status" label="状态" width="140"> </el-table-column>
      <el-table-column prop="role" label="角色" width="140"> </el-table-column>
      <el-table-column prop="cpu" label="CPU" width="140"> </el-table-column>
      <el-table-column prop="ram" label="内存" width="140"> </el-table-column>
      <el-table-column prop="pod" label="Pods" width="140"> </el-table-column>
      <el-table-column prop="allocatedCPU" label="已分配CPU" width="180"> </el-table-column>
      <el-table-column prop="allocatedRAM" label="已分配内存" width="180s"> </el-table-column>
    </el-table>
  </el-card>
</template>

<script>
export default {
  name: "NodeStatus",
  data() {
    return {
      nodeStatus: [
        {
          name: "master",
          status: "运行中",
          role: "master, worker",
          cpu: "17%",
          ram: "48%",
          pod: "50%",
          allocatedCPU: "2.877 Core (79%)",
          allocatedRAM: "2.177 Gi (30%)"
        },
        {
          name: "node1",
          status: "运行中",
          role: "worker",
          cpu: "13%",
          ram: "31%",
          pod: "25%",
          allocatedCPU: "1.612 Core (44%)",
          allocatedRAM: "2.06 Gi (29%)"
        },
        {
          name: "node2",
          status: "运行中",
          role: "worker",
          cpu: "9%",
          ram: "50%",
          pod: "23%",
          allocatedCPU: "1.752 Core (48%)",
          allocatedRAM: "4.703 Gi (66%)"
        }
      ]
    };
  }
};
</script>

<style lang="scss"></style>
