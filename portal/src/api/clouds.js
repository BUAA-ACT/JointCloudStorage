import request from "@/utils/request";
import store from "@/store";
// 云厂商接口
export default {
  /**
   * @deprecated
   * @see getAllClouds()
   */
  getAllCloud() {
    return {
      clouds: [
        {
          cloud_name: "阿里云-杭州",
          cloud_id: "aliyun-hangzhou",
          storage_price: 0.12,
          traffic_price: 0.5,
          availability: 0.99995,
          status: "UP",
          endpoint: "oss-cn-hangzhou.aliyuncs.com",
          access_key: "<ak-aliyun>",
          secret_key: "<sk-aliyun>",
          location: "120.188938,30.301958",
          bucket: "jcspan-hangzhou",
          address: "aliyun-hangzhou.jointcloudstorage.cn"
        },
        {
          cloud_name: "阿里云-呼和浩特",
          cloud_id: "aliyun-hohhot",
          storage_price: 0.1,
          traffic_price: 0.6,
          availability: 0.99995,
          status: "UP",
          endpoint: "oss-cn-huhehaote.aliyuncs.com",
          access_key: "<ak-aliyun>",
          secret_key: "<sk-aliyun>",
          location: "111.744578,40.873867",
          bucket: "jcspan-huhehaote",
          address: "aliyun-hohhot.jointcloudstorage.cn"
        },
        {
          cloud_name: "阿里云-青岛",
          cloud_id: "aliyun-qingdao",
          storage_price: 0.15,
          traffic_price: 0.4,
          availability: 0.99995,
          status: "UP",
          endpoint: "oss-cn-qingdao.aliyuncs.com",
          access_key: "<ak-aliyun>",
          secret_key: "<sk-aliyun>",
          location: "120.382109,36.075311",
          bucket: "jcspan-qingdao",
          address: "aliyun-qingdao.jointcloudstorage.cn"
        },
        {
          cloud_name: "腾讯云-成都",
          cloud_id: "txyun-chengdu",
          storage_price: 0.099,
          traffic_price: 0.5,
          availability: 0.9999,
          status: "UP",
          endpoint: "cos.ap-chengdu.myqcloud.com",
          access_key: "<ak-txyun>",
          secret_key: "<sk-txyun>",
          location: "104.072745,30.664271",
          bucket: "jcspan-chengdu-1259241028",
          address: "txyun-chengdu.jointcloudstorage.cn"
        },
        {
          cloud_name: "百度云-广州",
          cloud_id: "bdyun-guangzhou",
          storage_price: 0.119,
          traffic_price: 0.49,
          availability: 0.9995,
          status: "UP",
          endpoint: "s3.gz.bcebos.com",
          access_key: "<ak-bdyun>",
          secret_key: "<sk-bdyun>",
          location: "113.260506,23.132943",
          bucket: "jcspan-guangzhou",
          address: "bdyun-guangzhou.jointcloudstorage.cn"
        }
      ],
      ip: "",
      recommend_name: "",
      recommend_address: "",
      recommend_latency: null
    };
  },
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
  },
  /**
   * `新`获取所有云的接口
   * @see http://gitlab.act.buaa.edu.cn/jointcloudstorage/jcspan/-/wikis/httpserver-%E8%AE%BE%E8%AE%A1%E6%96%87%E6%A1%A3
   */
  getAllClouds() {
    return request({
      url: `/cloud/getAllClouds`,
      method: "post",
      data: {
        AccessToken: store.getters.token
      }
    });
  },
  addNewCloud(Cloud) {
    return request({
      url: `/cloud/newCloud`,
      method: "post",
      data: {
        AccessToken: store.getters.token,
        Cloud
      }
    });
  },
  /**
   * @description DO NOT CHANGE CLOUD ID!
   * @param Cloud
   * @returns Promise
   */
  changeCloudInfo(Cloud) {
    return request({
      url: `/cloud/changeCloudInfo`,
      method: "post",
      data: {
        AccessToken: store.getters.token,
        Cloud
      }
    });
  },
  getVoteRequests() {
    return request({
      url: `/cloud/getVoteRequests`,
      method: "post",
      data: {
        AccessToken: store.getters.token
      }
    });
  },
  /**
   * @param {String} CloudID
   * @param {Boolean} VoteResult
   * @returns {Promise}
   */
  voteForCloud(CloudID, VoteResult) {
    return request({
      url: `/cloud/vote`,
      method: "post",
      data: {
        AccessToken: store.getters.token,
        CloudID,
        VoteResult
      }
    });
  }
};
