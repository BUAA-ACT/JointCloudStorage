<template>
  <div class="dashboard-container">
    <HeadNav />
    <IndexPage />
    <Footer />
  </div>
</template>

<script>
import HeadNav from "../components/headNav.vue";
import IndexPage from "../components/indexPage.vue";
import Footer from "../components/footer.vue";

export default {
  name: "Index",
  components: { HeadNav, IndexPage, Footer },
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
