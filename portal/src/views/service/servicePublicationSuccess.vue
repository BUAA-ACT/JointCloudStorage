<template>
  <div class="dashboard-container">
    <HeadNav />
    <ServiceSuccess />
    <!-- <Footer/> -->
  </div>
</template>

<script>
import HeadNav from "@/components/headNav.vue";
import ServiceSuccess from "@/components/service/servicePublicationSuccess.vue";
// import Footer from "@/components/footer.vue";

export default {
  name: "RequirementReleaseSuccess",
  components: {
    HeadNav,
    ServiceSuccess
    // Footer
  },
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
