import axios from "axios";
import { Message } from "element-ui";
// import store from "@/store";
// import { getToken } from '@/utils/auth'
// import qs from "qs";

const successCode = 200;

// 创建axios 实例
const request = axios.create({
  // baseURL: "/api",
  timeout: 10000 // 请求超时时间
});
// request 拦截器
request.interceptors.request.use(
  config => {
    // if (store.getters.token) {
    //   // let each request carry token
    //   // ['X-Token'] is a custom headers key
    //   // please modify it according to the actual situation
    //   config.headers['X-Token'] = getToken()
    // }
    const newConfig = config;
    if (newConfig.method === "post") {
      newConfig.headers["Content-Type"] = "application/json";
      newConfig.data = config.data;
    }
    return newConfig;
  },
  error => {
    // do something with request error
    console.log(error); // for debug
    // Promise.reject(error);
  }
);
request.interceptors.response.use(
  response => {
    // if (store.getters.token) {
    //   // let each request carry token
    //   // ['X-Token'] is a custom headers key
    //   // please modify it according to the actual situation
    //   config.headers['X-Token'] = getToken()
    // }
    if (response.status === 200) {
      const res = response.data;
      console.log(res);
      if (res.code === successCode) {
        return res.data || true;
      }
      Message.error(res.msg);
      return false;
    }
    return false;
  },
  error => {
    // do something with request error
    console.log(error.response); // for debug
    Message.error(`Status: ${error.response.status} , Message: ${error.response.data.error}`);
  }
);
export default request;
