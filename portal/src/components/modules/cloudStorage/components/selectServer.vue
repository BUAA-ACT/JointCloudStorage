<template>
  <el-card shadow="always" class="requirementForm">
    <p class="indexTitle">{{ title || "跨云存储构建" }}</p>
    <div id="myEchart" ref="myEchart"></div>
    <el-form ref="form" :model="form" label-width="100px">
      <el-form-item label="数据中心">
        <el-checkbox-group v-model="form.city">
          <el-checkbox v-for="city in cities" :label="city.name" :key="city.key"
            >{{ city.name }}
            <el-radio-group v-model="form.server[city.key]">
              <el-radio-button label="阿里云"></el-radio-button>
              <el-radio-button label="百度云"></el-radio-button>
              <el-radio-button label="UCloud"></el-radio-button>
              <el-radio-button label="金山云"></el-radio-button>
            </el-radio-group>
          </el-checkbox>
        </el-checkbox-group>
      </el-form-item>
      <el-form-item label="存储模式">
        <el-radio-group v-model="form.resource">
          <el-radio-button label="多副本存储"></el-radio-button>
          <el-radio-button label="纠删码存储"></el-radio-button>
        </el-radio-group>
      </el-form-item>
    </el-form>
    <el-button type="primary" class="submitBtn" @click="onSubmit">确认提交</el-button>
  </el-card>
</template>

<script>
import echarts from "echarts";
import "echarts/map/js/world";

export default {
  name: "SelectServer",
  props: {
    title: String,
    submitFunc: {}
  },
  data() {
    return {
      form: {
        city: [],
        server: []
      },
      cities: [
        {
          name: "北京",
          key: 1
        },
        {
          name: "上海",
          key: 2
        },
        {
          name: "广州",
          key: 3
        },
        {
          name: "成都",
          key: 4
        },
        {
          name: "青岛",
          key: 5
        }
      ]
    };
  },
  mounted() {
    this.$nextTick(() => {
      this.initCharts();
    });
  },
  methods: {
    onSubmit() {
      this.$emit("selectedServerSumbit", this.form);
    },
    initCharts() {
      // const width = document.getElementById('chartDiv').offsetWidth;
      // const height = width * 0.6;
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
        // console.log('resize');
        chart.resize();
      };
      // })
    }
  }
};
</script>

<style lang="scss" scoped>
#myEchart {
  width: 800px;
  height: 500px;
  margin: 0 auto;
}
</style>
<style lang="scss">
.el-checkbox .el-radio-group {
  margin-left: 30px;
}
</style>
