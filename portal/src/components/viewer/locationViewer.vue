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
    }
  },
  methods: {
    initCharts() {
      console.log(this.clouds);
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
            name: "Top 5",
            type: "effectScatter",
            coordinateSystem: "geo",
            data: [...this.clouds],
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
                console.log(params);
                return params.data.value[2];
              }
            }
          }
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
