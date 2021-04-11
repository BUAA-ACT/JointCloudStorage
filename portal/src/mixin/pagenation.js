export default {
  data() {
    return {
      page: 1,
      size: 10,
      total: 0
    };
  },
  mounted() {
    this.getList();
  }
};
