<template>
  <div>
    <div class="mBanner mBanner-m1">
      <p class="bannerTitle">
        云服务评分
        <br />
        <span>对阿里云、腾讯云、华为云、uCloud上多种常用实例类型进行评测。</span>
      </p>
    </div>
    <h2>云服务测评</h2>
    <div class="manufacturer">
      <div class="processImg">
        <img src="@/assets/modules/serviceRecommendation/process1.png" alt="" />
        <!-- <el-row>
          <el-col :span="8">
            <p>VMs</p>
          </el-col>
          <el-col :span="8">
            <p>Run Benchmark</p>
          </el-col>
          <el-col :span="8">
            <p>Make Scores</p>
          </el-col>
        </el-row> -->
      </div>
      <p class="lineHeight15">
        使用Sysbench评测程序分别对CPU、内存和IO进行评测并打分。对三个指标的评测都分为单线程的评测和多线程的评测，多线程评测的线程数量等于实例的CPU数量。
        对内存的评测又分为读和写，对IO的评测分为顺序读、顺序写、随机读、随机写和随机读写。
      </p>
      <p>
        使用Netperf对内网性能进行评测，评测分为TCP和UDP两部分。
      </p>
      <br />
      <List
        ref="multipleTable"
        :columns="columns"
        :filterMap="filterMap"
        getListAction="serviceRecommendation.score"
        highlight-current-row
        :query="{ mark: 'page' }"
        tooltip-effect="dark"
        @selection-change="handleSelectionChange"
      >
        <template v-slot:filterBtns>
          <el-button type="primary" @click="compare">
            比较
          </el-button>
        </template>
      </List>
    </div>
  </div>
</template>
<script>
import List from "@/components/list.vue";
import { selectedRankColumns } from "@/scripts/map/serviceRank";

export default {
  name: "ServiceMark",
  data() {
    return {
      selectedRankColumns,
      evaluationTime: "20200824",
      selectedList: [],
      checkAllVal: false,
      filterMap: {
        evaluationTime: {
          label: "评测时间",
          type: "select",
          dataSource: [
            { label: "2020/08/24", value: "2020/08/24" },
            { label: "2020/08/17", value: "2020/08/17" }
          ]
        }
      },
      columns: [
        { type: "selection", label: "排名" },
        { prop: "instance", label: "实例类型" },
        { prop: "vendor", label: "地址" },
        { prop: "benchTime", label: "启动时间", formatter: (row, column, cellValue) => `${cellValue || 0}s` },
        { prop: "totalcpu", label: "cpu性能", formatter: (row, column, cellValue) => Number(cellValue).toFixed(2) },
        { prop: "totalmem", label: "内存性能", formatter: (row, column, cellValue) => Number(cellValue).toFixed(2) },
        { prop: "totalio", label: "磁盘性能", formatter: (row, column, cellValue) => Number(cellValue).toFixed(2) },
        { prop: "totalnet", label: "网络性能", formatter: (row, column, cellValue) => Number(cellValue).toFixed(2) },
        { prop: "total", label: "总分数", formatter: (row, column, cellValue) => Number(cellValue).toFixed(2) },
        {
          label: "操作",
          formatter: row => (
            <div>
              <el-button onClick={() => this.suspend(row)} type="text" size="small">
                编辑
              </el-button>
            </div>
          )
        }
      ]
    };
  },
  mounted() {},
  watch: {
    evaluationTime() {
      // console.log(e);
    }
  },
  components: {
    List
  },
  methods: {
    suspend(row) {
      this.$message.warning(row.instance);
    },
    handleSelectionChange(e) {
      // console.log(e);
      this.selectedList = [...e];
    },
    onCheckAllClick() {},
    getChecked() {},
    onCheckboxClick() {},
    compare() {
      if (this.selectedList.length < 2) {
        this.$message({
          message: "请选择至少两条进行对比",
          type: "warning"
        });
      } else {
        // console.log('compare');
        this.$router.push("/serviceRecommendation/compare");
      }
    }
  }
};
</script>
<style lang="scss" scoped>
.filterPart {
  margin-top: 40px;
  * {
    margin-right: 30px;
  }
}
</style>
