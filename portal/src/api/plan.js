import request from "@/utils/request";
import store from "@/store";

export default {
  getStoragePlans() {
    return request({
      url: `/plan/getAllStoragePlan`,
      method: "post",
      data: {
        AccessToken: store.getters.token
      }
    });
  },
  changeStoragePlan(StoragePlan) {
    return request({
      url: `/plan/chooseStoragePlan`,
      method: "post",
      data: {
        AccessToken: store.getters.token,
        StoragePlan
      }
    });
  },
  getNewAdvice() {
    return request({
      url: "/plan/getNewAdvice",
      method: "post",
      data: {
        AccessToken: store.getters.token
      }
    });
  },
  acceptAdvice() {
    return request({
      url: "/plan/acceptAdvice",
      method: "post",
      data: {
        AccessToken: store.getters.token
      }
    });
  },
  abandonAdvice() {
    return request({
      url: "/plan/abandonAdvice",
      method: "post",
      data: {
        AccessToken: store.getters.token
      }
    });
  }
};
