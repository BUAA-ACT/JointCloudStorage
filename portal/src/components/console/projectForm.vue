<template>
  <div class="projectInfo">
    <el-card class="chartCard" style="height:100px">
      <h4>
        项目
      </h4>
      <span class="tips">将根据项目资源进行分组, 可以按项目对资源进行查看管理</span>
    </el-card>
    <el-tabs type="border-card">
      <el-tab-pane label="用户项目">
        <el-table :data="nodeStatus">
          <el-table-column prop="name" label="名称" width="140"> </el-table-column>
          <el-table-column prop="status" label="状态" width="140"> </el-table-column>
          <el-table-column prop="space" label="企业空间" width="140"> </el-table-column>
          <el-table-column prop="cpu" label="CPU使用量" width="140"> </el-table-column>
          <el-table-column prop="ram" label="内存使用量" width="140"> </el-table-column>
          <el-table-column prop="pod" label="Pods数量" width="140"> </el-table-column>
        </el-table>
      </el-tab-pane>
      <el-tab-pane label="系统项目">
        <el-table :data="nodeStatus">
          <el-table-column prop="name" label="名称" width="140"> </el-table-column>
          <el-table-column prop="status" label="状态" width="140"> </el-table-column>
          <el-table-column prop="space" label="企业空间" width="140"> </el-table-column>
          <el-table-column prop="cpu" label="CPU使用量" width="140"> </el-table-column>
          <el-table-column prop="ram" label="内存使用量" width="140"> </el-table-column>
          <el-table-column prop="pod" label="Pods数量" width="140"> </el-table-column>
        </el-table>
      </el-tab-pane>
    </el-tabs>
  </div>
</template>

<script>
export default {
  name: "ProjectForm",
  data() {
    return {
      nodeStatus: [
        {
          name: "a111",
          status: "活跃",
          space: "",
          cpu: "0 m",
          ram: "0 Bytes",
          pod: "0"
        },
        {
          name: "node1",
          status: "活跃",
          space: "demo",
          cpu: "30 m",
          ram: "1.63 Gi",
          pod: "33"
        },
        {
          name: "aba",
          status: "活跃",
          space: "demo",
          cpu: "0 m",
          ram: "0 Bytes",
          pod: "0"
        },
        {
          name: "cert-manager",
          status: "活跃",
          space: "",
          cpu: "10 m",
          ram: "105.65 Mi",
          pod: "3"
        },
        {
          name: "a111",
          status: "活跃",
          space: "",
          cpu: "0 m",
          ram: "0 Bytes",
          pod: "0"
        },
        {
          name: "node1",
          status: "活跃",
          space: "demo",
          cpu: "30 m",
          ram: "1.63 Gi",
          pod: "33"
        },
        {
          name: "aba",
          status: "活跃",
          space: "demo",
          cpu: "0 m",
          ram: "0 Bytes",
          pod: "0"
        },
        {
          name: "cert-manager",
          status: "活跃",
          space: "",
          cpu: "10 m",
          ram: "105.65 Mi",
          pod: "3"
        },
        {
          name: "aba",
          status: "活跃",
          space: "demo",
          cpu: "0 m",
          ram: "0 Bytes",
          pod: "0"
        },
        {
          name: "cert-manager",
          status: "活跃",
          space: "",
          cpu: "10 m",
          ram: "105.65 Mi",
          pod: "3"
        }
      ]
    };
  }
};
</script>

<style lang="scss" scoped>
.projectInfo {
  text-align: left;
}
</style>
