import request from "@/utils/request";

// const BACKEND_URL = "http://localhost:8081";
export default {
  // create(query) {
  //   return request({
  //     url: "/api/register",
  //     method: "post",
  //     params: query
  //   });
  // },
  login(query) {
    return request({
      url: `/user/login`,
      method: "post",
      data: query
    });
  },
  register(query) {
    return request({
      url: `/user/register`,
      method: "post",
      data: query
    });
  },
  checkToken(AccessToken) {
    return request({
      url: `/user/checkValidity`,
      method: "post",
      data: { AccessToken }
    });
  },
  getUserInfo(AccessToken) {
    return request({
      url: `/user/getUserInfo`,
      method: "post",
      data: { AccessToken }
    });
  },
  logout(AccessToken) {
    return request({
      url: `/user/logout`,
      method: "post",
      data: { AccessToken }
    });
  }
};
