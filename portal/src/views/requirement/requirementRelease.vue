<template>
  <div class="dashboard-container">
    <HeadNav />
    <RequirementForm />
    <!-- <Footer/> -->
  </div>
</template>

<script>
import HeadNav from "@/components/headNav.vue";
import RequirementForm from "@/components/requirement/requirementForm.vue";
// import Footer from "@/components/footer.vue";

export default {
  name: "RequirementRelease",
  components: {
    HeadNav,
    RequirementForm
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
