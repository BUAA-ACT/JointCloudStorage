<template>
  <div>
    <List
      ref="multipleTable"
      :columns="fileViewMap"
      getListAction="serviceRecommendation.score"
      :query="{ mark: 'page' }"
      tooltip-effect="dark"
    ></List>
    <el-button @click="back">上一步</el-button>
  </div>
</template>

<script>
import List from "@/components/list.vue";
import { fileViewMap } from "@/scripts/map/caseMap";

export default {
  data() {
    return {
      fileViewMap
    };
  },
  components: { List },
  methods: {
    back() {
      this.$router.go(-1);
    },
    handleSelectionChange() {}
  }
};
</script>
