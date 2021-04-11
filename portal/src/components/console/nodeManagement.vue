<template>
  <div>
    <el-card class="chartCard">
      <h4>
        集群节点
      </h4>
      <span class="tips">集群节点提供了当前集群下节点的运行状态，以及可以编辑删除节点</span>
      <el-collapse v-model="activeName" accordion>
        <el-collapse-item title="集群节点的类型？" name="1">
          <div>节点分为主控 (Master) 节点和工作 (Worker) 节点</div>
        </el-collapse-item>
        <el-collapse-item title="什么是节点污点？" name="2">
          <div>
            节点污点 (Taints) 可以阻止某些容器组 (Pod) 副本部署至该节点中, 与容忍度 (Tolerations) 一起工作确保容器组不会被调度到不合适的节点上
          </div>
        </el-collapse-item>
      </el-collapse>
    </el-card>
    <el-card class="chartCard" style="height:150px;width:600px">
      <el-table style="margin-top:25px" :data="nodeReliefMessage">
        <el-table-column prop="totalNumber" label="节点数量"> </el-table-column>
        <el-table-column prop="masterNumber" label="控制节点数量"> </el-table-column>
        <el-table-column prop="nodeNumber" label="工作节点数量"> </el-table-column>
      </el-table>
    </el-card>
    <NodeStatus />
  </div>
</template>

<script>
import NodeStatus from "@/components/console/nodeStatus.vue";

export default {
  name: "",
  components: { NodeStatus },
  data() {
    return {
      activeName: "1",
      nodeReliefMessage: [{ totalNumber: "3", masterNumber: "1", nodeNumber: "3" }]
    };
  }
};
</script>

<style lang="scss">
.chartCard {
  text-align: left;
  span {
    margin-bottom: 10px;
  }
}
</style>
