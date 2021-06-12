<template>
  <el-card class="box-card" shadow="hover" v-on:click="this.jump">
    <div slot="header" class="clearfix">
      <span>{{ name?name:"云际节点" }}</span>
      <el-button style="float: right; padding: 3px 0" type="text" v-on:click="this.jump"> 点击前往 </el-button>
    </div>
    <div class="text item">
      云际标识： {{ id ? id : "null" }}
    </div>
    <div class="text item">
      服务状态：
      <el-tag v-if="status === 'UP' " type="success" size="mini">可用</el-tag>
      <el-tag v-if="status !== 'UP' " type="danger" size="mini">故障</el-tag>
    </div>
    <div class="text item">
      访问延迟： {{ latency }} ms
    </div>
  </el-card>
</template>

<script>
const speedtest = require('speedtest-url');
export default {
  name: "PanCard",
  data() {
    return {
      latency: "检测中"
    }
  },
  props: {
    name: String,
    id: String,
    url: String,
    status: String,
  },
  mounted() {
    this.$nextTick(function (){
      console.log("start test")
      this.latencyTest()
    })
  },
  methods: {
    jump: function () {
      window.location = "http://" + this.url
    },
    latencyTest: async function () {
      var type = [
        "line",
        "bar",
        "radar",
        "horizontalbar",
        "pie",
        "doughnut",
        "polararea",
        "scatter"
      ]
      var data = {
        "type": type[Math.floor(Math.random() * type.length)],
        "title": "website",
        "url_array": [
          "https://" + this.url,
        ]
      }
      var l = []
      for (var i =0; i < 5; i++) {
        speedtest.speedtest_url(data).then(hasil => {
          console.log(hasil)
          l.push(parseInt(hasil.message[0].ping))
          if (l.length === 5) {
            var sum = 0
            for (var d of l) {
              sum += d
            }
            this.latency = parseInt(sum/5)
            this.$emit('getLatency', {'id': this.id, 'latency': this.latency, 'name': this.name, 'address': this.url})
          }
        })
      }
      speedtest.chart_speedtest(data).then(hasil => console.log(hasil))
    }
  }
}
</script>

<style scoped>
.text {
  font-size: 14px;
}


</style>