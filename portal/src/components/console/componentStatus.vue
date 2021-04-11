<template>
  <el-card class="chartCard" style="height: 150px">
    <el-table style="padding:20px 0" :data="componentStatus">
      <el-table-column prop="apiRequest" label="API每秒请求数" width="140"> </el-table-column>
      <el-table-column prop="delay" label="API请求延迟" width="140"> </el-table-column>
      <el-table-column prop="schedulerTimes" label="调度器调度次数" width="140"> </el-table-column>
      <el-table-column prop="schedulingPods" label="调度失败的容器组" width="140"> </el-table-column>
    </el-table>
  </el-card>
</template>

<script>
export default {
  name: "ComponentStatus",
  data() {
    return {
      componentStatus: [
        {
          apiRequest: "15.483 times/s",
          delay: "2.41 ms",
          schedulerTimes: "71",
          schedulingPods: "2"
        }
      ]
    };
  }
};
</script>

<style lang="scss"></style>
