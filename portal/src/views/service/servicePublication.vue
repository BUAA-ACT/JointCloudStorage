<template>
  <div class="dashboard-container">
    <HeadNav />
    <ServiceForm />
    <Footer />
  </div>
</template>

<script>
import HeadNav from "@/components/headNav.vue";
import ServiceForm from "@/components/service/serviceForm.vue";
import Footer from "@/components/footer.vue";

export default {
  name: "ServicePublication",
  components: { HeadNav, ServiceForm, Footer },
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
