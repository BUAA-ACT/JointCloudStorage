<template>
  <div class="dashboard-container">
    <HeadNav :mini="true" />
    <UserCenter />
    <!--    <Footer />-->
  </div>
</template>

<script>
import HeadNav from "../components/headNav.vue";
import UserCenter from "../components/userCenter.vue";

export default {
  name: "Index",
  components: { HeadNav, UserCenter },
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
<style scoped>
.dashboard-container {
  min-height: 100vh;
  height: 100%;
}
</style>
