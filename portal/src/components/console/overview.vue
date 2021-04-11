<template>
  <div class="overview">
    <h4>hostHost 集群</h4>
    <span class="tips">
      Automatically created by KubeX, we encourage you to use host cluster for clusters management only,deploy workloads to member clusters.</span
    >
    <el-row :gutter="20">
      <el-col :span="12">
        <h5>集群信息</h5>
        <ClusterMessage />
      </el-col>
      <el-col :span="12">
        <h5>KubeX组件状态</h5>
        <ComponentStatus />
      </el-col>
      <el-col :span="16">
        <h5>集群资源使用情况</h5>
        <ResourceUsage />
      </el-col>
      <el-col :span="8">
        <h5>节点用量 Top3</h5>
        <NodeMessage />
      </el-col>
    </el-row>
  </div>
</template>

<script>
import ClusterMessage from "@/components/console/clusterMessage.vue";
import ComponentStatus from "@/components/console/componentStatus.vue";
import NodeMessage from "@/components/console/nodeMessage.vue";
import ResourceUsage from "@/components/console/resourceUsage.vue";

export default {
  data() {
    return {
      key: ""
    };
  },
  components: {
    ClusterMessage,
    ComponentStatus,
    NodeMessage,
    ResourceUsage
  }
};
</script>

<style lang="scss" scoped>
.overview {
  text-align: left;
}
</style>
