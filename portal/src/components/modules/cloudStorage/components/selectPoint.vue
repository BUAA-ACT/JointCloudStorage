<template>
  <div>
    <!-- <h4>用户服务实时监控</h4> -->
    <el-row :gutter="30">
      <el-col :span="8">
        <StatusLineChart />
      </el-col>
      <el-col :span="16">
        <ServerPercentage />
      </el-col>
      <el-col :span="24">
        <el-card class="requirementForm">
          <p>云厂商服务状态</p>
          <!-- <div class="radioGroup">
            <el-radio v-model="radio" label="1" border>
              东京服务器: d-auiufqiwuqwehquihiuqgefiquwd</el-radio>
            <el-radio v-model="radio" label="2" border>
              东京服务器: d-auiufqiwuqwehquihiuqgefiquwd</el-radio>
            <el-radio v-model="radio" label="3" border>
              东京服务器: d-auiufqiwuqwehquihiuqgefiquwd</el-radio>
          </div> -->
          <el-table ref="multipleTable" :data="tableData" highlight-current-row tooltip-effect="dark" style="width: 100%">
            <el-table-column prop="name" label="服务器名称" show-overflow-tooltip> </el-table-column>
            <el-table-column prop="delayTime" label="延时" width="120"> </el-table-column>
            <el-table-column prop="status" label="状态" width="120"> </el-table-column>
          </el-table>
          <!-- <el-button class="submitBtn" type="primary" @click="onSubmit">确认提交</el-button> -->
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script>
import StatusLineChart from "./statusLineChart.vue";
import ServerPercentage from "./serverPercentage.vue";

export default {
  name: "SelectPoint",
  data() {
    return {
      radio: "",
      tableData: [
        {
          name: "阿里云-北京",
          delayTime: "11ms",
          status: "online"
        },
        {
          name: "百度云-上海",
          delayTime: "26ms",
          status: "online"
        },
        {
          name: "金山云-广州",
          delayTime: "43ms",
          status: "online"
        }
      ]
    };
  },
  components: {
    StatusLineChart,
    ServerPercentage
  },

  methods: {
    onSubmit() {
      this.$emit("selectedPointSumbit", this.radio);
    }
  }
};
</script>
<style lang="scss"></style>
