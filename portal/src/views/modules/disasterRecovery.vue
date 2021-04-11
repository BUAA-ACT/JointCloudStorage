<template>
  <!-- 国防科大 baohan 云际容灾模块-->
  <div>
    <Header :navList="navList" title="云际容灾" />
    <router-view />
    <Footer type="mini" />
  </div>
</template>

<script>
import Header from "@/components/modules/header.vue";
import Footer from "@/components/footer.vue";

export default {
  name: "DisasterRecovery",
  components: { Header, Footer },
  data() {
    return {
      navList: [
        {
          name: "概述",
          url: "overview"
        },
        {
          name: "云际容灾示例",
          url: "example"
        },
        {
          name: "技术亮点",
          url: "technicalSpot"
        },
        {
          name: "文档资料",
          url: "documents"
        }
      ]
    };
  }
};
</script>
