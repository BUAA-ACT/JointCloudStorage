<template>
  <div class="viewer-container">
    <div class="location-viewer" ref="previewMap" id="previewMap"></div>
  </div>
</template>

<script>
/**
 * 云存储可视化工具
 * 格式：clouds 中存放云的位置
 * clouds[0]: [Longitude, Latitude, Caption]
 * 经度， 纬度， 注解
 */
import echarts from "echarts";
import "echarts/map/js/china";

export default {
  name: "locationViewer",
  props: {
    clouds: {
      type: Array,
      default: () => []
      // [Longitude, Latitude, Caption]
    },
    newClouds: {
      type: Array
      // [Longitude, Latitude, Caption]
    },
    formatFunction: {
      type: Function,
      default: Clouds => {
        return Clouds.map(value => {
          return {
            name: value.CloudID,
            value: value.Location.split(",").concat([
              `存储价格：${value.StoragePrice}元/GB/月<br/>
          流量价格：${value.TrafficPrice}元/GB<br/>
          可用性：${value.Availability}<br />`
            ])
          };
        });
      }
    },
    dynamic: {
      type: Boolean,
      default: false
    }
  },
  methods: {
    initCharts() {
      const chart = echarts.init(this.$refs.previewMap);

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
        tooltip: {
          trigger: "item"
        },
        geo: {
          map: "china",
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
          {
            name: "当前服务器",
            type: "effectScatter",
            coordinateSystem: "geo",
            data: [...this.formattedClouds],
            symbolSize: 30,
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
            zlevel: 100,
            tooltip: {
              formatter(params) {
                return params.data.value[2];
              }
            }
          },
          ...this.cloudMigration
        ]
      });
      window.onresize = () => {
        console.log("resize");
        chart.resize();
      };
    }
  },
  mounted() {
    this.initCharts();
  },
  watch: {
    clouds() {
      this.initCharts();
    },
    dynamic() {
      this.initCharts();
    }
  },
  computed: {
    cloudMigration() {
      if (!(this.newClouds && this.newClouds.length > 0)) {
        return [];
      }
      const newCloudsObj = {
        name: "新服务器",
        type: "effectScatter",
        coordinateSystem: "geo",
        data: [...this.formattedNewClouds],
        symbolSize: 30,
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
          color: "#4fb648",
          shadowBlur: 10,
          shadowColor: "#ffffff"
        },
        zlevel: 101,
        tooltip: {
          formatter(params) {
            return params.data.value[2];
          }
        }
      };
      const migrationLines = [];
      // [{
      //   coord: [116.4551,40.2539]
      // }, {
      //   coord: [121.4648,31.2891]
      // }]
      for (let i = 0; i < this.clouds.length; i += 1) {
        migrationLines.push([
          {
            coord: this.formattedClouds[i].value.slice(0, 2)
          },
          {
            coord: this.formattedNewClouds[i].value.slice(0, 2)
          }
        ]);
      }
      const migrationObj = {
        name: "",
        type: "lines",
        zlevel: 102,
        effect: {
          show: true,
          symbolSize: 10
        },
        lineStyle: {
          normal: { type: "dotted", color: "#0077ff", width: 2, curveness: 0.2 }
        },
        data: migrationLines
      };
      if (!this.dynamic) {
        return [newCloudsObj];
      }
      return [newCloudsObj, migrationObj];
    },
    formattedClouds() {
      return this.formatFunction(this.clouds);
    },
    formattedNewClouds() {
      return this.formatFunction(this.newClouds);
    }
  }
};
</script>

<style lang="scss" scoped>
.viewer-container {
  .location-viewer {
    width: 100%;
    height: 100%;
  }
}
</style>
