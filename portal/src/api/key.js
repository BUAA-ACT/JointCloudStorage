import request from "@/utils/request";
import store from "@/store";

export default {
  addKey(Comment) {
    return request({
      url: `/user/newKey`,
      method: "post",
      data: {
        AccessToken: store.getters.token,
        Comment
      }
    });
  },
  getAllKeys() {
    return request({
      url: `/user/getUserKeys`,
      method: "post",
      data: {
        AccessToken: store.getters.token
      }
    });
  },
  deleteKey(AccessKey) {
    return request({
      url: `/user/deleteKey`,
      method: "post",
      data: {
        AccessToken: store.getters.token,
        AccessKey
      }
    });
  },
  changeKeyComment(AccessKey, Comment) {
    return request({
      url: `/user/changeKeyComment`,
      method: "post",
      data: {
        AccessToken: store.getters.token,
        AccessKey,
        Comment
      }
    });
  },
  changeKeyStatus(AccessKey, Status) {
    return request({
      url: `/user/changeKeyStatus`,
      method: "post",
      data: {
        AccessToken: store.getters.token,
        AccessKey,
        Status
      }
    });
  },
  resetKey(AccessKey) {
    return request({
      url: `/user/remakeKey`,
      method: "post",
      data: {
        AccessToken: store.getters.token,
        AccessKey
      }
    });
  }
};
