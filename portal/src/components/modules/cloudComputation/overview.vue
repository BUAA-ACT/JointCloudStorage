<template>
  <div>
    <div class="mBanner mBanner-m3">
      <p class="bannerTitle">
        跨云计算
        <br />
        <span>HCloud是一套跨云计算框架</span>
        <br />
        <span>提供了基本的Serverless计算服务 <br />同时能够在不同的云平台，以及不同云服务之间进行快速动态迁移。</span>
      </p>
    </div>
    <h2>跨云计算</h2>
    <div class="manufacturer">
      <div class="processImg">
        <img src="@/assets/modules/cloudComputation/process3.png" alt="" />
      </div>
      <el-row class="overviewInfo" :gutter="30">
        <el-col :span="12">
          <p>HCloud是一套跨云计算框架</p>
          <span>
            -HCloud将Serverless计算任务分配给位于不同云平台上的工作进程。<br />
            -能够动态地在多个平台之间对计算任务进行快速迁移。<br />
            -扩容和迁移过程对用户透明。
          </span>
        </el-col>
        <el-col :span="12">
          <p>充分利用各个云平台的优势</p>
          <span>
            -HCloud将Serverless任务调度至时延最低，吞吐最大的云平台。<br />
            -能够选择成本最低的云平台和运行方式。
          </span>
        </el-col>
      </el-row>
    </div>
  </div>
</template>
<script>
export default {
  name: "overview",
  data() {
    return {};
  }
};
</script>
<style lang="scss">
.overviewInfo {
  margin-top: 10px;
  margin-bottom: 50px;
  p {
    font-weight: bold;
  }
  span {
    line-height: 2.5rem;
    font-size: 14px;
  }
}
</style>
