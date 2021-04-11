<template>
  <div>
    <div class="mBanner mBanner-m4">
      <p class="bannerTitle">
        云际存储
        <br />
        <span>海量，安全，低成本，高可靠的云存储服务</span>
      </p>
    </div>

    <div class="manufacturer">
      <el-row class="text" :gutter="30">
        <el-col :span="12">
          <h4>一· 云际异构存储资源抽象聚合：</h4>
          <p>
            将多云存储资源统一通过对象存储标准协议S3暴露给用户，用户在首次配置之后无需手动管理多个云即可实现数据的自动分发，同时可以在云际存储控制台监控及管理存储资源。
          </p>
        </el-col>
        <el-col :span="12">
          <img src="@/assets/modules/cloudStorage/img2.png" alt="" />
        </el-col>
      </el-row>
      <el-row class="text" :gutter="30">
        <el-col :span="12">
          <h4>二· 云服务商级别容灾：</h4>
          <p>实时监控云服务商的健康状态和价格，支持自动或手动触发数据迁移，迁移过程不影响上层应用。</p>
        </el-col>
        <el-col :span="12">
          <img src="@/assets/modules/cloudStorage/img3.png" alt="" />
        </el-col>
      </el-row>
    </div>
  </div>
</template>
<script>
export default {
  name: "ServiceMark",
  data() {
    return {
      evaluationTime: "20200824",
      selectedList: [],
      tableData: [
        {
          date: "2016-05-03",
          name: "王小虎",
          address: "上海市普陀区金沙江路 1518 弄"
        },
        {
          date: "2016-05-03",
          name: "王小虎",
          address: "上海市普陀区金沙江路 1518 弄"
        }
      ]
    };
  },
  watch: {
    evaluationTime() {
      // console.log(e);
    }
  },
  methods: {
    handleSelectionChange(e) {
      // console.log(e);
      this.selectedList = [...e];
    },
    compare() {
      if (this.selectedList.length < 2) {
        this.$message({
          message: "请选择至少两条进行对比",
          type: "warning"
        });
      } else {
        // console.log('compare');
      }
    }
  }
};
</script>
<style lang="scss">
.text {
  p {
    text-indent: 2em;
    line-height: 2em;
    -webkit-margin-before: 0em;
    margin-block-start: 0em;
    -webkit-margin-after: 0em;
    margin-block-end: 0em;
  }
  img {
    width: 100%;
  }
  .el-col {
    margin-bottom: 30px;
  }
}
</style>
