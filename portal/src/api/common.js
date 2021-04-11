import request from "@/utils/request";

const BACKEND_URL = "http://localhost:8081";
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
      url: `${BACKEND_URL}/user/login`,
      method: "post",
      data: query
    });
  },
  register(query) {
    return request({
      url: `${BACKEND_URL}/user/register`,
      method: "post",
      data: query
    });
  },
  checkToken(AccessToken) {
    return request({
      url: `${BACKEND_URL}/user/checkValidity`,
      method: "post",
      data: { AccessToken }
    });
  },
  getUserInfo(AccessToken) {
    return request({
      url: `${BACKEND_URL}/user/getUserInfo`,
      method: "post",
      data: { AccessToken }
    });
  },
  logout(AccessToken) {
    return request({
      url: `${BACKEND_URL}/user/logout`,
      method: "post",
      data: { AccessToken }
    });
  }
};
