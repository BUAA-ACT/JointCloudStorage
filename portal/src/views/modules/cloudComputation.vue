<template>
  <!-- 上海交大 跨云计算模块-->
  <div>
    <Header :navList="navList" title="跨云计算" />
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
          name: "计算示例",
          url: "computeExample"
        },
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
