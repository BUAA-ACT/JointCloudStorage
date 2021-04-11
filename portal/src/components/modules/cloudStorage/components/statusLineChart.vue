<template>
  <el-card class="chartCard">
    <p>服务运行状态</p>
    <div id="lineChart" ref="lineChart"></div>
  </el-card>
</template>

<script>
import echarts from "echarts";

export default {
  name: "StatusLineChart",
  data() {
    return {
      lineChartOption: {
        xAxis: {
          type: "category",
          boundaryGap: false,
          axisLine: { lineStyle: { color: "#999999" } },
          data: ["Z1", "Z2", "Z3", "Z4", "Z5", "Z6"]
        },
        yAxis: {
          type: "value",
          axisLine: { lineStyle: { color: "#999999" } }
        },
        series: [
          {
            data: [51, 49, 76, 48, 45, 55],
            type: "line",
            itemStyle: {
              color: "rgba(0, 55, 255, 1)"
            },
            markPoint: {
              itemStyle: {
                color: "#409eff"
              },
              data: [{ type: "max", name: "最大值" }]
            },
            smooth: true,
            areaStyle: {
              color: {
                type: "linear",
                x: 0,
                y: 0,
                x2: 0,
                y2: 1,
                colorStops: [
                  {
                    offset: 0,
                    color: "#409eff" // 0% 处的颜色
                  },
                  {
                    offset: 1,
                    color: "white" // 100% 处的颜色
                  }
                ],
                global: false // 缺省为 false
              }
            }
          }
        ]
      }
    };
  },
  mounted() {
    this.$nextTick(() => {
      this.drawLinechart();
    });
  },
  methods: {
    drawLinechart() {
      const lineChart = echarts.init(this.$refs.lineChart);
      lineChart.setOption(this.lineChartOption);
    }
  }
};
</script>

<style lang="scss">
#lineChart {
  width: 100%;
  height: 200px;
}
</style>
