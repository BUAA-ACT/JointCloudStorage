<template>
  <div class="manufacturer">
    <el-row>
      <el-col :span="3" :offset="3">
        <span style="line-height:40px">评测时间</span>
      </el-col>
      <el-col :span="4">
        <el-select v-model="evaluationTime">
          <el-option label="2020/08/24" value="20200824"></el-option>
          <el-option label="2020/08/17" value="20200817"></el-option>
        </el-select>
      </el-col>
    </el-row>

    <div id="radarChart" ref="radarChart"></div>
    <div id="lineChart" ref="lineChart"></div>
  </div>
</template>

<script>
import echarts from "echarts";

export default {
  name: "compare",
  data() {
    return {
      evaluationTime: ""
    };
  },
  created() {
    this.$nextTick(() => {
      this.initCharts();
      this.drawLinechart();
    });
  },
  methods: {
    initCharts() {
      const chart = echarts.init(this.$refs.radarChart);

      // 把配置和数据放这里
      chart.setOption({
        title: {
          text: "实例各类指标性能",
          left: "36%"
        },
        tooltip: {},
        legend: {
          // data: names,
          top: "30px"
        },
        radar: {
          // shape: 'circle',
          name: {
            textStyle: {
              color: "#fff",
              backgroundColor: "#999",
              borderRadius: 3,
              padding: [3, 5]
            }
          },
          indicator: [
            { name: "总分数", max: 600 },
            { name: "启动时间（ms）", max: 15 },
            { name: "CPU性能分数", max: 150 },
            { name: "内存性能分数", max: 150 },
            { name: "磁盘性能分数", max: 150 },
            { name: "网络性能分数", max: 150 }
          ],
          radius: 200
        },
        series: [
          {
            name: "实例性能",
            type: "radar",
            // areaStyle: {normal: {}},
            data: [
              {
                value: [355, 13, 100, 100, 100, 100],
                name: "阿里云通用型ecs.g6.large"
              },
              {
                value: [400, 10, 110, 110, 110, 110],
                name: "阿里云通用型ecs.hfg6.large"
              }
            ]
          }
        ]
      });
      window.onresize = () => {
        console.log("resize");
        chart.resize();
      };
    },
    drawLinechart() {
      const lineChart = echarts.init(this.$refs.lineChart);
      const data = [
        { name: "2020/11/11", value: ["2020/11/11", 800] },
        { name: "2020/11/12", value: ["2020/11/12", 810] },
        { name: "2020/11/13", value: ["2020/11/13", 830] },
        { name: "2020/11/14", value: ["2020/11/14", 850] },
        { name: "2020/11/15", value: ["2020/11/15", 860] },
        { name: "2020/11/16", value: ["2020/11/16", 880] }
      ];
      // let date = [];
      lineChart.setOption({
        title: {
          text: "实例性能随时间变化图",
          left: "33%"
        },
        tooltip: {
          trigger: "axis"
        },
        legend: {
          // data: names,
          top: "30px"
        },
        xAxis: {
          type: "time"
          // data: date
        },
        yAxis: {
          type: "value",
          boundaryGap: [0, "100%"],
          splitLine: {
            show: false
          }
        },
        series: [
          {
            name: "阿里云通用型ecs.g6.large",
            type: "line",
            data
          },
          {
            name: "阿里云通用型ecs.hfg6.large",
            type: "line",
            data
          }
        ]
      });
    }
  }
};
</script>

<style lang="scss">
.manufacturer {
  #radarChart,
  #lineChart {
    width: 800px;
    height: 600px;
    margin: 30px auto;
  }
}
</style>
