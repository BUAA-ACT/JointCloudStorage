<template>
  <div>
    <el-row type="flex" justify="center">
      <h2>云际网盘 JcsPan</h2>
    </el-row>
    <el-row type="flex" justify="center" class="center">
      <el-col :span="20" class="info">
        <p>当前设备 ip 地址：{{ ip }}, 推荐您访问 <el-link type="primary" v-bind:href="'http://'+recommend_address">  {{ recommend_name }}  </el-link>节点，延迟 {{ recommend_latency }} ms
        </p>
        <div class="go-button">
          <el-button type="success" plain v-on:click="jump">点击进入</el-button>
        </div>
        <el-divider></el-divider>
      </el-col>
    </el-row >
    <el-row type="flex" justify="center">
      <el-col :span="18">
        <h3 class="pan-card">云际网盘节点列表</h3>
      </el-col>
    </el-row>
    <el-row v-for="i in rowNum" :key="i" type="flex" justify="center">
      <el-col v-for="j in 3" :key="j" :md="6" :sm="20">
        <pan-card v-if="(i-1)*3+j <= clouds.length" class="pan-card" @getLatency="recordLatency"
          v-bind:name="clouds[(i-1)*3+j-1]['cloud_name']" v-bind:id="clouds[(i-1)*3+j-1]['cloud_id']"
                  v-bind:status="clouds[(i-1)*3+j-1]['status']"
                  v-bind:url="clouds[(i-1)*3+j-1]['address']"
        ></pan-card>
      </el-col>
    </el-row>
  </div>
</template>

<script>
import PanCard from "@/components/PanCard";
const publicIp = require('public-ip');
export default {
  name: "JcsPanIndex",
  components: {
    PanCard,
  },
  data: function () {
    return {
      clouds:     [
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
          "address": "182.61.24.215",
        }
      ],
      ip: "",
      recommend_name: "",
      recommend_address: "",
      recommend_latency: null,
    };
  },
  computed: {
    rowNum: function () {
      //return 4/3
      return Math.ceil(this.clouds.length / 3);
    }
  },
  created() {
    (async () => {
      this.ip = await publicIp.v4( {
        fallbackUrls: [
          'https://ifconfig.me/ip'
        ]
      });
    })();
  },
  methods: {
    jump: function () {
      console.log(this.recommend_address)
      window.location = "http://" + this.recommend_address
    },
    recordLatency(info) {
      if (!this.recommend_latency || info.latency<this.recommend_latency){
        this.recommend_latency = info.latency
        this.recommend_name = info.name
        this.recommend_address = info.address
      }
    }
  }
}
</script>


<style scoped>
.pan-card {
  margin-top: 20px;
  margin-left: 20px;
}
.center {
  text-align: center;
}
.go-button{
  padding-top: 20px;
}
.info {
  font-size: 14px;
}
</style>