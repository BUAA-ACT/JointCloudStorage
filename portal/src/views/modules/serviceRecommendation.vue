<template>
  <!-- 北京大学 资源推荐模块-->
  <div>
    <Header :navList="navList" title="云服务推荐" />
    <router-view />
    <Footer type="mini" />
  </div>
</template>

<script>
import Header from "@/components/modules/header.vue";
import Footer from "@/components/footer.vue";

export default {
  name: "ServiceRecommendation",
  components: { Header, Footer },
  data() {
    return {
      navList: [
        {
          name: "云服务测评",
          url: "serviceMark"
        },
        {
          name: "云服务推荐",
          url: "serviceRecom"
        },
        {
          name: "机器学习",
          url: "studyRecom"
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
