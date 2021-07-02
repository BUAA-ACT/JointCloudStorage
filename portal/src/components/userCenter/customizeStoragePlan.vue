<template>
  <el-card>
    <el-steps :active="curStep" finish-status="success">
      <el-step v-for="index in finishStep" :title="`步骤 ${index}`" :key="String(index)"></el-step>
    </el-steps>
    <el-form v-if="curStep === 0" title="请选择存储方式">
      <el-form-item label="存储方式">
        <el-radio v-model="curPlan.StorageMode" label="Replica">多副本</el-radio>
        <el-radio v-model="curPlan.StorageMode" label="EC">纠删码</el-radio>
      </el-form-item>
    </el-form>
    <el-button v-if="curStep > 0" @click="prevStep">上一步</el-button>
    <el-button @click="nextStep">{{ curStep === finishStep ? "完成" : "下一步" }}</el-button>
  </el-card>
</template>

<script>
export default {
  name: "customizeStoragePlan",
  data() {
    return {
      curStep: 0,
      curStepFinish: false,
      curPlan: {
        StorageMode: "Replica"
      }
    };
  },
  methods: {
    prevStep() {
      this.curStep -= 1;
      this.curStepFinish = true;
    },
    nextStep() {
      this.curStep += 1;
      this.curStepFinish = false;
    }
  },
  computed: {
    finishStep() {
      if (this.curPlan.StorageMode === "Replica") {
        return 3;
      }
      if (this.curPlan.StorageMode === "EC") {
        return 4;
      }
      return -1;
    }
  }
};
</script>

<style scoped></style>
