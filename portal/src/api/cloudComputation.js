import request from "@/utils/request";

// 云计算接口
export default {
  createFunction(query) {
    return request({
      url: "/function/createFunction",
      method: "post",
      data: query
    });
  },
  deleteFunctionByName() {
    return request({
      url: "/function/delete/functionByName",
      method: "delete"
    });
  },
  deleteFunction(functionId) {
    return request({
      url: `/function/deleteFunction/${functionId}`,
      method: "delete"
    });
  },
  getFunctionByFunctionName() {
    return request({
      url: "/function/delete/functionByName",
      method: "get"
    });
  },
  invokeFunction(query) {
    return request({
      url: "/function/invokeFunction",
      method: "post",
      data: query
    });
  },
  listFunctions(query) {
    return request({
      url: "/function/lists",
      method: "get",
      params: query || { page: 1, size: 10 }
    });
  },
  updateFunction() {
    return request({
      url: "/function/updateFunction",
      method: "put"
    });
  }
};
