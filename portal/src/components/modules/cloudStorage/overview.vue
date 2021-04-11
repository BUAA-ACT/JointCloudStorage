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
      <el-row class="overviewInfo" :gutter="30">
        <el-col :span="12">
          <div class="text">
            <p>
              本模块支撑软件定义的云际计算运行与管理的关键，针对大数据应用、计算框架等对资源动态个性化的需求，
              以云际计算环境跨地理位置分布的存储资源为主体，重点研究云际存储资源抽象聚合与评估、云际分布存储与数据管理、
              多源异构数据感知的应用部署与协同调度等，实现全局的效能优化与最小化多副本开销，降低云际不确定性对服务性能带来的影响，
              实现软件定义的云际大数据高效能存储管理，形成一系列云际计算环境下大数据资源分布式管理与处理的高效运行支撑理论、机制和关键技术等。
            </p>
            <p>
              针对现有公有云存储服务存在的供应商锁定、单云故障影响上层业务等问题，云际存储系统提出了多云互联、
              共同向用户提供存储能力的解决方案。该系统可根据用户自定义的存储策略，将用户数据分布在多个公有云厂商的多个数据中心内，
              并提供纠删码和多副本两种冗余方案。当单云发生故障时，系统仍能够利用分布在其他云的冗余数据块恢复用户数据，不影响上层应用的正常运行。
            </p>
            <p>
              云际存储系统适合对数据可靠性要求很高，或希望将数据按需分布到多个数据中心的用户。
            </p>
          </div>
        </el-col>
        <el-col :span="12">
          <img src="@/assets/modules/cloudStorage/img1.png" alt="" />
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
  .text p {
    font-weight: normal;
    line-height: 2rem;
    font-size: 14px;
    color: #000000;
    text-indent: 2rem;
  }
  img {
    width: 100%;
  }
}
</style>
