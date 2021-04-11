<template>
  <div>
    <div class="mBanner mBanner-m1">
      <p class="bannerTitle">
        云服务推荐
        <br />
        <span>针对您的需求推荐最适合的云服务器。</span>
      </p>
    </div>
    <h2>云服务推荐</h2>
    <div class="manufacturer">
      <div class="processImg">
        <img src="@/assets/modules/serviceRecommendation/process2.png" alt="" />
      </div>
    </div>
    <h2>最佳实例推荐</h2>
    <div class="manufacturer">
      <!-- <el-table
        ref="multipleTable"
        :data="tableData"
        highlight-current-row
        tooltip-effect="dark"
        style="width: 100%">
        <el-table-column
          prop="name"
          label="实例类型"
          show-overflow-tooltip>
        </el-table-column>
        <el-table-column
          prop="place"
          label="地域"
          width="120">
        </el-table-column>
        <el-table-column
          prop="time"
          label="启动时间"
          width="120">
        </el-table-column>
        <el-table-column
          prop="cpuRate"
          label="CPU性能"
          width="120">
        </el-table-column>
        <el-table-column
          prop="storageRate"
          label="内存性能"
          width="120">
        </el-table-column>
        <el-table-column
          prop="driveRate"
          label="磁盘性能"
          width="120">
        </el-table-column>
        <el-table-column
          prop="netRate"
          label="网络性能"
          width="120">
        </el-table-column>
        <el-table-column
          prop="Rate"
          label="总评分"
          >
        </el-table-column>
      </el-table> -->
      <List
        ref="multipleTable"
        :columns="serviceRankColumns"
        getListAction="serviceRecommendation.scoreTotal"
        highlight-current-row
        :query="{ mark: 'page' }"
        tooltip-effect="dark"
      ></List>
    </div>
    <h2>最佳性价比实例推荐</h2>
    <div class="manufacturer">
      <List
        ref="multipleTable"
        :columns="serviceRankColumns"
        getListAction="serviceRecommendation.scoreCost"
        highlight-current-row
        :query="{ mark: 'page' }"
        tooltip-effect="dark"
      ></List>
    </div>
  </div>
</template>
<script>
import List from "@/components/list.vue";
import { serviceRankColumns } from "@/scripts/map/serviceRank";

export default {
  name: "ServiceRecom",
  data() {
    return {
      serviceRankColumns,
      map: {
        index: "id",
        instance: "实例类型",
        vendor: "地址",
        benchTime: "启动时间",
        totalcpu: "cpu性能",
        totalmem: "内存性能",
        totalio: "磁盘性能",
        totalnet: "网络性能",
        total: "总分数"
      }
    };
  },
  components: {
    List
  },
  methods: {}
};
</script>
