<template>
  <div>
    <el-tabs class="dataCenter" v-model="activeName">
      <el-tab-pane :label="$t('disasterReco.caseTab1')" name="first">
        <el-card class="caseRecom">
          <div class="leftTitle">
            方案一
          </div>
          <div class="content">
            <p>高性能CPU方案 总价：<span class="red">300¥/月</span></p>
            <List border :columns="caseMap" getListAction="serviceRecommendation.score"></List>
          </div>
        </el-card>
      </el-tab-pane>
      <el-tab-pane :label="$t('disasterReco.caseTab2')" name="second">
        <h4>{{ $t("disasterReco.sourceList") }}</h4>
        <List
          ref="multipleTable"
          :columns="selectedCaseMap"
          getListAction="serviceRecommendation.score"
          :query="{ mark: 'page' }"
          tooltip-effect="dark"
          @selection-change="handleSelectionChange"
        ></List>
      </el-tab-pane>
    </el-tabs>
    <el-button @click="back">{{ $t("disasterReco.prev") }}</el-button>
    <el-button type="primary" class="submitBtn" @click="onSubmit"> {{ $t("disasterReco.next") }}</el-button>
  </div>
</template>

<script>
import List from "@/components/list.vue";
import { caseMap, selectedCaseMap } from "@/scripts/map/caseMap";

export default {
  data() {
    return {
      caseMap,
      selectedCaseMap,
      activeName: "first"
    };
  },
  components: { List },
  methods: {
    handleSelectionChange() {},
    back() {
      this.$router.go(-1);
    },
    onSubmit() {
      this.$router.push({ path: "fileView" });
    }
  }
};
</script>
<style lang="scss">
.caseRecom {
  margin-bottom: 20px;
  background: #f9f9f9;
  padding-right: 20px;
  .el-card__body {
    padding: 0;
    padding-top: 20px;
    width: 100%;
    height: 100%;
    position: relative;
    p {
      margin-top: 0;
    }
    .el-table th {
      background: #ffffff;
    }
    .leftTitle {
      width: 25px;
      height: 100%;
      background: #409eff;
      color: #fff;
      position: absolute;
      padding-top: 10%;
      padding-left: 5px;
      top: 0;
      left: 0;
    }
    .content {
      margin-left: 50px;
    }
  }
}
</style>
