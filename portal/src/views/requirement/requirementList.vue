<template>
  <div class="dashboard-container">
    <HeadNav />
    <RequirementCatalogue />
    <Footer />
  </div>
</template>

<script>
import HeadNav from "@/components/headNav.vue";
import RequirementCatalogue from "@/components/requirement/requirementList.vue";
import Footer from "@/components/footer.vue";

export default {
  name: "RequirementList",
  components: { HeadNav, RequirementCatalogue, Footer },
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
