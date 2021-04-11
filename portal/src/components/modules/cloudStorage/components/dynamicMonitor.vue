<template>
  <div>
    <h4>用户服务实时监控 <a class="insideFloatRight" @click="back">返回</a></h4>
    <el-row :gutter="30">
      <el-col :span="8">
        <el-card class="chartCard">
          <p>数据完成率</p>
          <div id="gaugeChart" ref="gaugeChart"></div>
        </el-card>
      </el-col>
      <el-col :span="16">
        <ServerPercentage />
      </el-col>
      <el-col :span="8">
        <el-card class="chartCard">
          <p>数据资源生产量</p>
          <div id="barChart" ref="barChart"></div>
        </el-card>
      </el-col>
      <el-col :span="16">
        <el-card class="chartCard">
          <p>数据资源运行产生趋势</p>
          <div id="lineChart" ref="lineChart"></div>
        </el-card>
      </el-col>
      <el-col :span="8">
        <StatusLineChart />
      </el-col>
      <el-col :span="8">
        <el-card class="chartCard">
          <p>设备工作率</p>
          <div id="ybarChart" ref="ybarChart"></div>
        </el-card>
      </el-col>
      <el-col :span="8">
        <el-card class="chartCard">
          <p>服务设备资源利用率</p>
          <div class="center">
            <el-progress type="dashboard" :percentage="75" :color="colors"></el-progress>
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script>
import echarts from "echarts";
import StatusLineChart from "./statusLineChart.vue";
import ServerPercentage from "./serverPercentage.vue";

export default {
  name: "DynamicMonitor",
  data() {
    return {
      colors: [
        { color: "#f56c6c", percentage: 20 },
        { color: "#feb23a", percentage: 40 },
        { color: "#5cb87a", percentage: 60 },
        { color: "#1989fa", percentage: 80 },
        { color: "#6f7ad3", percentage: 100 }
      ],
      barChartOption: {
        color: {
          type: "linear",
          x: 0,
          y: 0,
          x2: 0,
          y2: 1,
          colorStops: [
            {
              offset: 0,
              color: "#9fc6ef" // 0% 处的颜色
            },
            {
              offset: 1,
              color: "#409eff" // 100% 处的颜色
            }
          ],
          global: false // 缺省为 false
        },
        tooltip: {
          trigger: "axis",
          axisPointer: {
            // 坐标轴指示器，坐标轴触发有效
            type: "shadow" // 默认为直线，可选为：'line' | 'shadow'
          }
        },
        xAxis: {
          // boundaryGap: false,
          axisLine: { lineStyle: { color: "#999999" } },
          data: ["Z1", "Z2", "Z3", "Z4", "Z5", "Z6"]
        },
        yAxis: [
          {
            axisLine: { lineStyle: { color: "#999999" } },
            type: "value"
          }
        ],
        series: [
          {
            name: "数据资源生产量",
            type: "bar",
            barWidth: "50%",
            data: [51, 49, 76, 48, 45, 55]
          }
        ]
      },
      yAxisBarChartOption: {
        color: {
          type: "linear",
          x: 0,
          y: 0,
          x2: 1,
          y2: 0,
          colorStops: [
            {
              offset: 0,
              color: "#409eff" // 0% 处的颜色
            },
            {
              offset: 1,
              color: "#00c8dc" // 100% 处的颜色
            }
          ],
          global: false // 缺省为 false
        },
        yAxis: {
          type: "category",
          axisLine: { lineStyle: { color: "#999999" } },
          data: ["Z1", "Z2", "Z3", "Z4", "Z5", "Z6"]
        },
        xAxis: {
          axisLine: { lineStyle: { color: "#999999" } },
          type: "value"
        },
        series: [
          {
            data: [5, 6, 3, 4, 2, 4],
            type: "bar",
            barWidth: "30%"
          }
        ]
      },
      lineChartOption: {
        xAxis: {
          type: "category",
          boundaryGap: false,
          axisLine: { lineStyle: { color: "#999999" } },
          data: ["Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"]
        },
        yAxis: {
          type: "value",
          axisLine: { lineStyle: { color: "#999999" } }
        },
        series: [
          {
            data: [820, 932, 901, 934, 1090, 850, 1020],
            type: "line",
            smooth: true,
            itemStyle: {
              color: "#feb23a"
            }
          },
          {
            data: [1000, 800, 901, 800, 1190, 800, 900],
            type: "line",
            smooth: true,
            itemStyle: {
              color: "#409eff"
            }
          }
        ]
      },
      gaugeChartOption: {
        tooltip: {
          formatter: "{a} <br/>{b} : {c}%"
        },
        series: [
          {
            name: "完成率",
            type: "gauge",
            detail: {
              formatter: "{value}%",
              offsetCenter: [0, "0"]
            },
            data: [{ value: 60, name: "" }],
            itemStyle: {
              color: "transparent"
              // show: false,
            },
            splitLine: { show: false },
            axisTick: { show: false },
            axisLabel: { show: false },
            axisLine: {
              lineStyle: {
                color: [
                  [0.2, "#dddddd"],
                  [0.8, "#409eff"],
                  [1, "#30c0b7"]
                ]
              }
            }
          }
        ]
      }
    };
  },
  components: {
    StatusLineChart,
    ServerPercentage
  },
  mounted() {
    this.$nextTick(() => {
      this.drawBarchart();
    });
  },
  methods: {
    drawBarchart() {
      const barChart = echarts.init(this.$refs.barChart);
      barChart.setOption(this.barChartOption);
      const ybarChart = echarts.init(this.$refs.ybarChart);
      ybarChart.setOption(this.yAxisBarChartOption);
      const lineChart = echarts.init(this.$refs.lineChart);
      lineChart.setOption(this.lineChartOption);
      const gaugeChart = echarts.init(this.$refs.gaugeChart);
      gaugeChart.setOption(this.gaugeChartOption, true);
    },
    back() {
      this.$emit("back");
    }
  }
};
</script>

<style lang="scss" scoped>
#barChart,
#ybarChart,
#lineChart,
#gaugeChart {
  width: 100%;
  height: 240px;
}
.center {
  margin-top: 50px;
}
.insideFloatRight {
  float: right;
  font-weight: normal;
  color: #999999;
  font-size: 14px;
  margin: 10px;
}
</style>
