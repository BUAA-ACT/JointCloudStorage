<template>
  <div class="dashboard-container">
    <HeadNav />
    <BuyKubeX />
    <Footer />
  </div>
</template>

<script>
import HeadNav from "@/components/headNav.vue";
import BuyKubeX from "@/components/clusterService.vue";
import Footer from "@/components/footer.vue";

export default {
  components: { HeadNav, BuyKubeX, Footer }
};
</script>

<style lang="scss" scoped></style>
