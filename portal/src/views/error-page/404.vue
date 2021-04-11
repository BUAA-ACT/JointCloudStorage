<template>
  <div>
    <HeadNav />
    <div class="manufacturer">
      <img class="img404" src="@/assets/developing.png" alt="" />
    </div>
    <!-- <a href="/">back home</a> -->
  </div>
</template>
<script>
import HeadNav from "@/components/headNav.vue";

export default {
  name: "Page404",
  components: {
    HeadNav
  }
};
</script>
<style lang="scss">
.img404 {
  width: 50%;
  display: block;
  margin: 50px auto;
}
</style>
