<template>
  <div>
    <div class="manufacturer">
      <el-card shadow="always" class="requirementForm">
        <p class="indexTitle">综合需求发布 > 云服务存储需求</p>
        <!-- <el-image
          style="width: 100%;"
          src="@/assets/wmap.png"
          fit="cover">
        </el-image> -->
        <div class="comprehensiveCase" id="chartDiv">
          <!-- <img src="@/assets/wmap.png"> -->
          <div id="myEchart" ref="myEchart"></div>
          <baidu-map class="bd-map" :ak="ak" center="北京">
            <bm-navigation anchor="BMAP_ANCHOR_TOP_RIGHT"></bm-navigation>
            <bm-overview-map anchor="BMAP_ANCHOR_BOTTOM_RIGHT" :isOpen="true"></bm-overview-map>
          </baidu-map>
          <p>综合服务推荐列表</p>
          <el-row :gutter="30">
            <el-col :span="8" v-for="(item, index) in comprehensiveCase" :key="index">
              <el-card shadow="always">
                <p class="title">{{ item.name }}</p>
                <p>{{ item.description }}</p>
                <a href="">查看详情 ></a>
              </el-card>
            </el-col>
          </el-row>
        </div>

        <!-- <el-button type="primary" class="submitBtn" @click="onSubmit">确认提交</el-button> -->
      </el-card>
    </div>
    <div class="manufacturer">
      <RecomList />
    </div>
  </div>
</template>
<script>
import echarts from "echarts";
//   import '../../node_modules/echarts/map/js/world.js'
import "echarts/map/js/world";
import RecomList from "@/components/requirement/recomList.vue";
import BaiduMap from "vue-baidu-map/components/map/Map.vue";

// 引入中国地图数据
export default {
  name: "ComprehensiveRequirement",
  data() {
    return {
      comprehensiveCase: [
        {
          name: "方案一",
          description: "建议用户可在北京地区增购华为云2核4GB云服务器1台；将闲置服务器sa1.small2迁移至阿里云1核2GB服务器，此方案性价比最高。"
        },
        {
          name: "方案二",
          description:
            "建议用户可在北京地区将负载过重服务器ecs.s2.xlarge迁移至阿里云4核16GB服务器；将闲置服务器sa1.small2迁移至腾讯云1核1GB服务器，此方案价格最便宜。"
        },
        {
          name: "方案三",
          description:
            "建议用户可在北京地区将负载过重服务器ecs.s2.xlarge迁移至阿里云4核16GB服务器，同时增购华为云2核4GB云服务器1台；将闲置服务器sa1.small2迁移至腾讯云1核1GB服务器，此方案性能最稳定。"
        }
      ]
    };
  },
  components: { RecomList, BaiduMap },
  mounted() {
    this.initCharts();
  },
  methods: {
    onSubmit() {},
    initCharts() {
      // const width = document.getElementById('chartDiv').offsetWidth;
      // const height = width * 0.6;
      // debugger;

      const chart = echarts.init(this.$refs.myEchart);

      // 把配置和数据放这里
      chart.setOption({
        backgroundColor: "#ffffff",
        title: {
          // text: '设备分布',
          left: "40%",
          top: "0px",
          textStyle: {
            color: "#fff",
            opacity: 0.7
          }
        },
        // dataRange: {
        //   show: false,
        //   min: 0,
        //   max: 1000000,
        //   text: ['High', 'Low'],
        //   realtime: true,
        //   calculable: true,
        //   color: ['orangered', 'yellow', 'lightskyblue'],
        // },

        tooltip: {
          trigger: "item"
        },
        geo: {
          map: "world",
          label: {
            emphasis: {
              show: false
            }
          },
          roam: false,
          silent: true,
          zoom: 1.2,
          itemStyle: {
            normal: {
              areaColor: "#dddddd",
              borderColor: "#dddddd"
            },
            emphasis: {
              borderColor: "#fff",
              areaColor: "#5b9bd5",
              borderWidth: 1
            }
          }
        },
        series: [
          // {
          //   type: 'map',
          //   mapType: 'world',
          //   // zoom: 1,
          //   mapLocation: {
          //     y: 100,
          //   },
          //   data: [{
          //     name: 'Afghanistan',
          //     value: 186,
          //   },
          //   {
          //     name: 'Spain',
          //     value: 38,
          //   },
          //   {
          //     name: 'France',
          //     value: 866,
          //   },
          //   {
          //     name: 'United Kingdom',
          //     value: 35,
          //   },
          //   {
          //     name: 'Greece',
          //     value: 99,
          //   },
          //   {
          //     name: 'Greenland',
          //     value: 56,
          //   },
          //   {
          //     name: 'India',
          //     value: 120,
          //   },
          //   {
          //     name: 'Ireland',
          //     value: 41,
          //   },
          //   {
          //     name: 'Japan',
          //     value: 33,
          //   },
          //   {
          //     name: 'South Korea',
          //     value: 52,
          //   },
          //   {
          //     name: 'Mexico',
          //     value: 4,
          //   },
          //   {
          //     name: 'Malaysia',
          //     value: 35,
          //   },
          //   {
          //     name: 'Vietnam',
          //     value: 397,
          //   },
          //   {
          //     name: 'China',
          //     value: 1993,
          //   },
          //   ],
          //   symbolSize: 12,
          //   label: {
          //     normal: {
          //       show: false,
          //     },
          //     emphasis: {
          //       show: false,
          //     },
          //   },
          //   itemStyle: {
          //     normal: {
          //       areaColor: '#dddddd',
          //       borderColor: '#dddddd',
          //     },
          //     emphasis: {
          //       borderColor: '#fff',
          //       areaColor: '#5b9bd5',
          //       borderWidth: 1,
          //     },
          //   },
          // },
          {
            name: "服务器数量",
            type: "scatter",
            coordinateSystem: "geo",
            data: [
              { name: "Brazil", value: [-51.92528, -14.235004, 35] },
              { name: "Canada", value: [-106.346771, 56.130366, 24] },
              { name: "China", value: [104.195397, 35.86166, 550] },
              { name: "Germany", value: [10.451526, 51.165691, 45] },
              { name: "Finland", value: [25.748151, 61.92411, 23] },
              { name: "Hong Kong", value: [114.109497, 22.396428, 350] },
              { name: "Japan", value: [138.252924, 36.204824, 100] },
              { name: "South Korea", value: [127.766922, 35.907757, 50] },
              { name: "Malaysia", value: [101.975766, 4.210484, 430] }
            ],
            symbolSize(val) {
              return val[2] / 10;
            },
            encode: {
              value: 2
            },
            label: {
              formatter: "{b}",
              position: "right",
              show: false
            },
            itemStyle: {
              color: "#5b9bd5"
            },
            emphasis: {
              label: {
                show: true
              }
            }
          },
          {
            name: "Top 5",
            type: "effectScatter",
            coordinateSystem: "geo",
            // data: convertData(data.sort((a, b) => b.value - a.value).slice(0, 6)),
            data: [
              { name: "China", value: [104.195397, 35.86166, 550] },
              { name: "Hong Kong", value: [114.109497, 22.396428, 350] },
              { name: "Japan", value: [138.252924, 36.204824, 100] },
              { name: "Malaysia", value: [101.975766, 4.210484, 430] }
            ],
            symbolSize(val) {
              return val[2] / 10;
            },
            encode: {
              value: 2
            },
            showEffectOn: "render",
            rippleEffect: {
              brushType: "stroke"
            },
            hoverAnimation: true,
            label: {
              formatter: "{b}",
              position: "right",
              show: true
            },
            itemStyle: {
              color: "#5b9bd5",
              shadowBlur: 10,
              shadowColor: "#fff"
            },
            zlevel: 1
          }
        ]
      });
      // this.$nextTick(()=>{
      window.onresize = () => {
        console.log("resize");
        chart.resize();
      };
      // })
    }
    // chinaConfigure() {
    //   const myChart = echarts.init(this.$refs.myEchart); // 这里是为了获得容器所在位置
    //   window.onresize = myChart.resize;
    //   myChart.setOption({ // 进行相关配置
    //     backgroundColor: '#ffffff',
    //     tooltip: {}, // 鼠标移到图里面的浮动提示框
    //     dataRange: {
    //       show: false,
    //       min: 0,
    //       max: 1000,
    //       text: ['High', 'Low'],
    //       realtime: true,
    //       calculable: true,
    //       color: ['orangered', 'yellow', 'lightskyblue'],
    //     },
    //     geo: { // 这个是重点配置区
    //       map: 'china', // 表示中国地图
    //       roam: true,
    //       label: {
    //         normal: {
    //           show: false, // 是否显示对应地名
    //           textStyle: {
    //             color: 'rgba(0,0,0,0.4)',
    //           },
    //         },
    //       },
    //       itemStyle: {
    //         normal: {
    //           borderColor: 'rgba(0, 0, 0, 0.2)',
    //         },
    //         emphasis: {
    //           areaColor: null,
    //           shadowOffsetX: 0,
    //           shadowOffsetY: 0,
    //           shadowBlur: 20,
    //           borderWidth: 0,
    //           shadowColor: 'rgba(0, 0, 0, 0.5)',
    //         },
    //       },
    //     },
    //     series: [{
    //       type: 'scatter',
    //       coordinateSystem: 'geo', // 对应上方配置
    //     },
    //     {
    //       name: '服务器数量', // 浮动框的标题
    //       type: 'map',
    //       geoIndex: 0,
    //       data: [{
    //         name: '北京',
    //         value: 599,
    //       }, {
    //         name: '上海',
    //         value: 142,
    //       }, {
    //         name: '黑龙江',
    //         value: 44,
    //       }, {
    //         name: '深圳',
    //         value: 92,
    //       }, {
    //         name: '湖北',
    //         value: 810,
    //       }, {
    //         name: '四川',
    //         value: 453,
    //       }],
    //     }],
    //   });
    // },
  },
  computed: {
    ak() {
      return process.env.VUE_APP_BDMAP_AK;
    }
  }
};
</script>
<style lang="scss">
@keyframes hover-blue {
  0% {
    -webkit-filter: brightness(90%);
    background-color: #419eff;
    opacity: 0.1;
  }
  50% {
    -webkit-filter: brightness(90%);
    background-color: #419eff;
    opacity: 0.3;
  }
  100% {
    /*css3滤镜亮度百分比*/
    -webkit-filter: brightness(100%);
    background-color: #419eff;
    opacity: 1;
  }
}
.comprehensiveCase {
  padding: 0 20px;
  margin-bottom: 20px;
  > p {
    text-align: center;
    font-weight: 400;
    color: #000000;
  }
  img {
    width: 100%;
  }
  .el-card {
    min-height: 200px;
    .el-card__body {
      padding: 20px !important;
      overflow: auto;
    }
    p,
    a {
      font-size: 14px;
      color: #666666;
      line-height: 1.8em;
    }
    .title {
      font-size: 16px;
      font-weight: 400;
    }
    .title,
    a {
      color: #333333;
    }
    a {
      float: right;
    }
  }
  .el-card:hover {
    -webkit-animation-name: hover-blue;
    animation-name: hover-blue;
    -webkit-animation-duration: 0.5s;
    animation-duration: 0.5s;
    -webkit-animation-timing-function: ease-in;
    animation-timing-function: ease-in;
    -webkit-animation-iteration-count: 1;
    animation-iteration-count: 1;
    background-color: #419eff;
    p,
    a {
      color: #ffffff;
    }
  }
}
#myEchart {
  width: 800px;
  height: 500px;
  margin: 0 auto;
}
.bd-map {
  height: 60vh;
}
</style>
