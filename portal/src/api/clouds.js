import request from "@/utils/request";

// 云厂商接口
export default {
  getAllCloud() {
    return {
      clouds: [
        {
          "cloud_name": "阿里云-杭州",
          "cloud_id": "aliyun-hangzhou",
          "storage_price": 0.12,
          "traffic_price": 0.5,
          "availability": 0.99995,
          "status": "UP",
          "endpoint": "oss-cn-hangzhou.aliyuncs.com",
          "access_key": "<ak-aliyun>",
          "secret_key": "<sk-aliyun>",
          "location": "120.188938,30.301958",
          "bucket": "jcspan-hangzhou",
          "address": "aliyun-hangzhou.jointcloudstorage.cn",
        },
        {
          "cloud_name": "阿里云-呼和浩特",
          "cloud_id": "aliyun-hohhot",
          "storage_price": 0.1,
          "traffic_price": 0.6,
          "availability": 0.99995,
          "status": "UP",
          "endpoint": "oss-cn-huhehaote.aliyuncs.com",
          "access_key": "<ak-aliyun>",
          "secret_key": "<sk-aliyun>",
          "location": "111.744578,40.873867",
          "bucket": "jcspan-huhehaote",
          "address": "aliyun-hohhot.jointcloudstorage.cn"
        },
        {
          "cloud_name": "阿里云-青岛",
          "cloud_id": "aliyun-qingdao",
          "storage_price": 0.15,
          "traffic_price": 0.4,
          "availability": 0.99995,
          "status": "UP",
          "endpoint": "oss-cn-qingdao.aliyuncs.com",
          "access_key": "<ak-aliyun>",
          "secret_key": "<sk-aliyun>",
          "location": "120.382109,36.075311",
          "bucket": "jcspan-qingdao",
          "address": "aliyun-qingdao.jointcloudstorage.cn",
        },
        {
          "cloud_name": "腾讯云-成都",
          "cloud_id": "txyun-chengdu",
          "storage_price": 0.099,
          "traffic_price": 0.5,
          "availability": 0.9999,
          "status": "UP",
          "endpoint": "cos.ap-chengdu.myqcloud.com",
          "access_key": "<ak-txyun>",
          "secret_key": "<sk-txyun>",
          "location": "104.072745,30.664271",
          "bucket": "jcspan-chengdu-1259241028",
          "address": "txyun-chengdu.jointcloudstorage.cn",
        },
        {
          "cloud_name": "百度云-广州",
          "cloud_id": "bdyun-guangzhou",
          "storage_price": 0.119,
          "traffic_price": 0.49,
          "availability": 0.9995,
          "status": "UP",
          "endpoint": "s3.gz.bcebos.com",
          "access_key": "<ak-bdyun>",
          "secret_key": "<sk-bdyun>",
          "location": "113.260506,23.132943",
          "bucket": "jcspan-guangzhou",
          "address": "bdyun-guangzhou.jointcloudstorage.cn",
        }
      ],
      ip: "",
      recommend_name: "",
      recommend_address: "",
      recommend_latency: null,
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
  }
};
