<template>
  <!-- 北航 云际存储模块-->
  <div>
    <Header :navList="navList" title="云际存储" />
    <router-view />
    <Footer type="mini" />
  </div>
</template>

<script>
import Header from "@/components/modules/header.vue";
import Footer from "@/components/footer.vue";

export default {
  name: "CloudComputation",
  components: { Header, Footer },
  data() {
    return {
      navList: [
        {
          name: "概述",
          url: "overview"
        },
        {
          name: "存储构建",
          url: "storageBuild"
        },
        {
          name: "存储迁移",
          url: "storageMove"
        },
        // {
        //   name: '云际容灾',
        //   url: 'disasterRecovery',
        // },
        {
          name: "技术亮点",
          url: "technicalSpot"
        }
        // {
        //   name: '文档资料',
        //   url: 'documents',
        // },
      ]
    };
  }
};
</script>
