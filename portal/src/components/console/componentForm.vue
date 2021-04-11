<template>
  <el-card class="chartCard" style="height: 1000px">
    <p>服务组件信息</p>
    <el-table :data="nodeStatus">
      <el-table-column prop="name" label="组件名称" width="300"> </el-table-column>
      <el-table-column prop="status" label="状态" width="140"> </el-table-column>
      <el-table-column prop="copy" label="副本数量" width="140"> </el-table-column>
      <el-table-column prop="runtime" label="运行时间" width="140"> </el-table-column>
    </el-table>
  </el-card>
</template>

<script>
export default {
  name: "ComponentForm",
  data() {
    return {
      nodeStatus: [
        {
          name: "ks-console",
          status: "健康",
          copy: "1/1",
          runtime: "4个月"
        },
        {
          name: "ks-apiserver",
          status: "健康",
          copy: "1/1",
          runtime: "4个月"
        },
        {
          name: "tower",
          status: "健康",
          copy: "1/1",
          runtime: "4个月"
        },
        {
          name: "redis",
          status: "健康",
          copy: "1/1",
          runtime: "4个月"
        },
        {
          name: "ks-controller-manager",
          status: "健康",
          copy: "1/1",
          runtime: "4个月"
        },
        {
          name: "etcd",
          status: "健康",
          copy: "1/1",
          runtime: "4个月"
        },
        {
          name: "mysql",
          status: "异常",
          copy: "1/1",
          runtime: "4个月"
        },
        {
          name: "kube-controller-manager-svc",
          status: "健康",
          copy: "1/1",
          runtime: "4个月"
        },
        {
          name: "coredns",
          status: "健康",
          copy: "1/1",
          runtime: "4个月"
        },
        {
          name: "metrics-server",
          status: "健康",
          copy: "1/1",
          runtime: "4个月"
        },
        {
          name: "kube-scheduler-svc",
          status: "健康",
          copy: "1/1",
          runtime: "4个月"
        },
        {
          name: "elasticsearch-logging-data",
          status: "健康",
          copy: "1/1",
          runtime: "4个月"
        },
        {
          name: "hyperpitrix",
          status: "异常",
          copy: "1/1",
          runtime: "4个月"
        },
        {
          name: "elasticsearch-logging-discovery",
          status: "异常",
          copy: "1/1",
          runtime: "4个月"
        },
        {
          name: "node-exporter",
          status: "健康",
          copy: "1/1",
          runtime: "4个月"
        }
      ]
    };
  }
};
</script>

<style lang="scss"></style>
