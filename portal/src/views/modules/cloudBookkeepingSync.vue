<template>
  <!-- 国防科大 江枫 云际同步记账模块-->
  <div>
    <Header :navList="navList" title="云际同步记账" />
    <router-view />
    <Footer type="mini" />
  </div>
</template>

<script>
import Header from "@/components/modules/header.vue";
import Footer from "@/components/footer.vue";

export default {
  name: "CloudBookkeepingSync",
  components: { Header, Footer },
  data() {
    return {
      navList: [
        {
          name: "概述",
          url: "overview"
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
