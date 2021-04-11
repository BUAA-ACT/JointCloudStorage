<template>
  <div class="dashboard-container">
    <HeadNav />
    <ServiceCatalogue />
    <Footer />
  </div>
</template>

<script>
import HeadNav from "@/components/headNav.vue";
import ServiceCatalogue from "@/components/service/serviceList.vue";
import Footer from "@/components/footer.vue";

export default {
  name: "ServiceList",
  components: { HeadNav, ServiceCatalogue, Footer },
  data() {
    return {};
  },
  computed: {},
  created() {
    // if (!this.roles.includes('admin')) {
    //   this.currentRole = 'editorDashboard'
    // }
  }
};
</script>
