<template>
  <div>
    <div class="mBanner mBanner-m3">
      <p class="bannerTitle">
        云际容灾
        <br />
        <span>xxxx</span>
      </p>
    </div>

    <div class="manufacturer">
      <h2>技术亮点</h2>
      <el-divider></el-divider>
      <div class="text">
        <p>
          现有的基于纠删码的存储系统均是面向单云数据中心环境设计的。在跨云数据中心环境下，它们写入、读取、修复数据时需要跨云数据中心（通过公网）传输大量数据，由于公网带宽远低于内网带宽，所以它们写入、读取、修复数据的速度较低。本系统包含了一套专为跨云数据中心纠删码设计的数据写入方法、数据读取方法和数据修复方法，可能通过最小化跨云数据中心传输量来提高写入、读取、修复速度。
        </p>
        <h4>
          1、跨云数据中心纠删码写入方法
        </h4>
        <p>
          本方法通过对比接收端云数据中心和发送端云数据中心中的编码块数来选择编码节点，使跨云传输量始终等于接收端云数据中心和发送端云数据中心中编码块数的最小值，从而有效降低了写入数据时的跨云数据中心传输量，从而提高了写入速度。此外，由于本方法将数据的写入过程组织为两级流水线（第一级：客户端写入数据块；第二级：各存储节点协同完成校验块的编码），因而能够进一步提高批量数据的写入速度。
        </p>
        <h4>1. 降低启动成本</h4>
        <p>
          Serverless服务使用可以免去自己搭建版本管理服务器、数据库服务器等基础设施，仅需要编写业务逻辑代码即可，对初创企业更加友好。部署、运维等其他服务均由Serverless服务商提供
        </p>
        <h4>2、跨云数据中心纠删码读取方法</h4>
        <p>
          本方法通过读取距离客户端最近的若干编码块、并对其进行解码来获取原始数据。这种读取方式以增加计算开销为代价达到降低跨云数据中心传输量的目的，由于跨云数据中心纠删码的传输耗时远大于计算耗时，因此这种方式能够提高读取速度。
        </p>
        <h4>3、跨云数据中心纠删码单云服务器修复方法</h4>
        <p>
          本方法在修复单个故障云服务器时采用树-星混合拓扑组织数据传输和数据合并：先在各个云数据中心内部采用星型传输拓扑对辅助块进行合并，然后构建最小网络距离生成树来组织各个合并后的辅助块的传输。本方法使得修复云服务器中的每个编码块时的跨云数据中心传输块数不大于云数据中心数减一，因此能够有效提高修复速度。
        </p>
        <h4>4、跨云数据中心纠删码云数据中心修复方法</h4>
        <p>
          本方法通过对比接收端云数据中心中的待修复块数和发送端云数据中心中的辅助块数来选择对辅助修复块进行合并的节点，使跨云数据中心传输量始终等于接收端云数据中心中待修复块数和发送端云数据中心中辅助块数的最小值，从而有效降低了云数据中心修复时的跨云传输量，进而提高了云数据中心的修复速度。
        </p>
      </div>
    </div>
  </div>
</template>
<script>
export default {
  name: "technicalSpot",
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
}
</style>
