import request from "@/utils/request";
import { Message } from "element-ui";

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
  },
  changeUserInfo(field, newVal, AccessToken) {
    let err = "";
    if (AccessToken === undefined) {
      err = "获取AccessToken失败！";
      Message.error(err);
      return Promise.reject(err);
    }
    let Field = "";
    switch (field) {
      case "username":
      case "password":
        Field = field.charAt(0).toUpperCase() + field.slice(1);
        break;
      default:
        err = `未知字段：${field}`;
        Message.error(err);
        return Promise.reject(err);
    }
    return request({
      url: `/user/changeUser${Field}`,
      method: "post",
      data: {
        AccessToken,
        [`New${Field}`]: newVal
      }
    });
  },
  changePassword(OriginPassword, NewPassword, AccessToken) {
    return request({
      url: `/user/changePassword`,
      method: "post",
      data: {
        AccessToken,
        OriginPassword,
        NewPassword
      }
    });
  }
};
