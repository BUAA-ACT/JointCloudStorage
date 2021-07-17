<template>
  <div>
    <plan-advice v-if="!migrating" />
    <div class="no-new-plan" v-else>
      <i class="el-icon-info tip-icon"></i><br />
      数据迁移中……
      <!-- TODO: 迁移进度     -->
    </div>
  </div>
</template>

<script>
import PlanAdvice from "@/components/userCenter/planAdvice.vue";

export default {
  name: "dataMigration",
  components: { PlanAdvice },
  data() {
    return {
      intervalId: null
    };
  },
  computed: {
    migrating() {
      return this.$store.getters.status === "FORBIDDEN";
    }
  },
  methods: {
    checkCurStatus() {
      this.$store.dispatch("getInfo");
    }
  },
  mounted() {
    this.intervalId = setInterval(this.checkCurStatus, 2000);
  },
  beforeDestroy() {
    if (this.intervalId) {
      clearInterval(this.intervalId);
    }
  }
};
</script>

<style scoped lang="scss">
.tip-icon {
  font-size: 100px;
  color: #909399;
}
</style>
