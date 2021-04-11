<template>
  <!-- 国防科大 李琳慧 云际异步记账模块-->
  <div>
    <Header :navList="navList" title="云际异步记账" />
    <router-view />
    <Footer type="mini" />
  </div>
</template>

<script>
import Header from "@/components/modules/header.vue";
import Footer from "@/components/footer.vue";

export default {
  name: "CloudBookkeepingAsync",
  components: { Header, Footer },
  data() {
    return {
      navList: [
        {
          name: "概述",
          url: "overview"
        },
        {
          name: "云际记账示例",
          url: "example"
        },
        {
          name: "技术亮点",
          url: "technicalSpot"
        }
      ]
    };
  }
};
</script>
