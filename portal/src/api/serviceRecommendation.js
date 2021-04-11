import request from "@/utils/request";

// 云推荐接口
export default {
  scoreTotal(query) {
    return request({
      url: "/scoreTotal/page",
      method: "get",
      params: query
    });
  },
  scoreCost(query) {
    return request({
      url: "/scoreCost/page",
      method: "get",
      params: query
    });
  },
  score(query) {
    return request({
      url: "/score/page",
      method: "get",
      params: query
    });
  }
};
