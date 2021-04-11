<template>
  <div class="dashboard-container">
    <HeadNav :mini="true" />
    <ConsoleDashboard />
    <Footer />
  </div>
</template>

<script>
import HeadNav from "@/components/headNav.vue";
import ConsoleDashboard from "@/components/console.vue";
import Footer from "@/components/footer.vue";

export default {
  name: "Index",
  components: { HeadNav, ConsoleDashboard, Footer },
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
