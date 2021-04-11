<template>
  <div class="dashboard-container">
    <HeadNav />
    <RequirementSuccess />
    <!-- <Footer/> -->
  </div>
</template>

<script>
import HeadNav from "@/components/headNav.vue";
import RequirementSuccess from "@/components/requirement/requirementReleaseSuccess.vue";
// import Footer from "@/components/footer.vue";

export default {
  name: "RequirementReleaseSuccess",
  components: {
    HeadNav,
    RequirementSuccess
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
