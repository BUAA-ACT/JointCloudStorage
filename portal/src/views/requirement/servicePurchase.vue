<template>
  <div class="dashboard-container">
    <HeadNav />
    <BuyService />
    <Footer />
  </div>
</template>

<script>
import HeadNav from "@/components/headNav.vue";
import BuyService from "@/components/requirement/servicePurchase.vue";
import Footer from "@/components/footer.vue";

export default {
  components: { HeadNav, BuyService, Footer }
};
</script>

<style lang="scss" scoped></style>
