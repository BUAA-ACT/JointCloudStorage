<template>
  <div class="dashboard-container">
    <HeadNav />
    <SolutionCase />
    <Footer />
  </div>
</template>

<script>
import HeadNav from "../components/headNav.vue";
import SolutionCase from "../components/solutionCase.vue";
import Footer from "../components/footer.vue";

export default {
  components: { HeadNav, SolutionCase, Footer },
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
