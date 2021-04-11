/* eslint-disable import/no-extraneous-dependencies */
<template>
  <div>
    <div class="mBanner mBanner-m3">
      <p class="bannerTitle">
        跨云计算
        <br />
        <span>HCloud是一套跨云计算框架</span>
        <br />
        <span>提供了基本的Serverless计算服务 <br />同时能够在不同的云平台，以及不同云服务之间进行快速动态迁移。</span>
      </p>
    </div>
    <h2>跨云计算</h2>
    <div class="manufacturer cloudComputation">
      <el-tabs v-model="activeName" @tab-click="handleClick">
        <el-tab-pane label="函数面板" name="0">
          <el-button style="margin-top:30px" @click="activeName = '1'">创建函数</el-button>
          <el-row :gutter="30">
            <el-col :span="6" v-for="(item, index) in functionList" :key="index">
              <el-card class="functionCard">
                <p class="title">{{ item.name }} <el-button type="text" class="insideFloatRight" @click="goDetail(item)">details</el-button></p>
                <el-divider></el-divider>
                <p>env:{{ item.env }}</p>
                <p>memory:{{ item.memory }}</p>
                <p>timeout:{{ item.timeout }}</p>
                <el-popconfirm confirm-button-text="好的" cancel-button-text="不用了" icon="el-icon-info" icon-color="red" title="确认删除此函数么？">
                  <el-button slot="reference" type="danger" size="small">删除</el-button>
                </el-popconfirm>
              </el-card>
            </el-col>
          </el-row>
        </el-tab-pane>
        <el-tab-pane label="创建函数" name="1">
          <div class="bg-gray" v-if="activeName === '1'">
            <el-row :gutter="30">
              <el-col :span="12">
                <el-form v-model="form" ref="form" label-width="150px">
                  <el-form-item label="Function Name:">
                    <el-input v-model="form.name"></el-input>
                  </el-form-item>
                  <el-form-item label="Env:">
                    <el-select v-model="form.env" placeholder="请选择">
                      <el-option label="java" value="java"></el-option>
                      <el-option label="python" value="python"></el-option>
                    </el-select>
                  </el-form-item>
                  <el-form-item label="TimeOut(1-300s):">
                    <el-input v-model="form.timeout"></el-input>
                  </el-form-item>
                  <el-form-item label="Memory Size(128-1024MB):">
                    <el-input v-model="form.memory"></el-input>
                  </el-form-item>
                </el-form>
              </el-col>
              <el-col :span="12">
                <span class="el-form-item__label">Entrypoint:</span>
                <div class="codemirror" v-if="form.env === 'python'">
                  <codemirror v-model="code" :options="cmOption" />
                </div>
                <div v-if="form.env === 'java'">
                  <el-upload
                    class="upload-demo"
                    action="https://jsonplaceholder.typicode.com/posts/"
                    :on-preview="handlePreview"
                    :on-remove="handleRemove"
                    :before-remove="beforeRemove"
                    :limit="1"
                    :on-exceed="handleExceed"
                    :file-list="fileList"
                  >
                    <el-button size="small" type="primary">上传jar包</el-button>
                    <div slot="tip" class="el-upload__tip">只能上传jar文件</div>
                  </el-upload>
                </div>
              </el-col>
            </el-row>
            <el-button @click="createFunc">下一步</el-button>
          </div>
        </el-tab-pane>
        <el-tab-pane label="提交运行" name="2">
          <div class="bg-gray">
            <table class="el-table funcTable">
              <tr>
                <td>Function Name:</td>
                <td :colspan="3">{{ form.name }}</td>
                <td>Created Time:</td>
                <td>{{ form.createTime }}</td>
              </tr>
              <tr>
                <td>Env:</td>
                <td>{{ form.env }}</td>
                <td>Memory Size:</td>
                <td>{{ form.memory }}</td>
                <td>Timeout:</td>
                <td>{{ form.timeout }}</td>
              </tr>
              <tr>
                <td>Code Size:</td>
                <td>{{ form.codeSize }}</td>
                <td>Code Checksum:</td>
                <td :colspan="3">{{ form.codeChecksum }}</td>
              </tr>
              <tr>
                <td>Description</td>
                <td :colspan="5">{{ form.description }}</td>
              </tr>
              <tr>
                <td>EnviromentVariable:</td>
                <td :colspan="2">{{ form.envVar }}</td>
                <td :colspan="3">
                  <!-- Enable Native Serverless:
                  <el-switch
                    v-model="form.enableNS"
                    size="small">
                  </el-switch> -->
                </td>
              </tr>
            </table>

            <el-row style="margin-bottom:30px" :gutter="30">
              <el-col :span="9">
                <p>Args：</p>
                <el-input type="textarea" rows="6" v-model="args"></el-input>
              </el-col>
              <el-col :span="2" :offset="1">
                <el-button class="invokeBtn" @click="invoke">invoke</el-button>
              </el-col>
              <el-col :offset="1" :span="9">
                <p>函数返回值：</p>
                <el-input rows="6" type="textarea" readonly v-model="result"></el-input>
              </el-col>
            </el-row>
            <el-button @click="activeName = '1'">上一步</el-button>
            <el-button @click="getResult">查看过程</el-button>
          </div>
        </el-tab-pane>
      </el-tabs>
    </div>
  </div>
</template>
<script>
import dedent from "dedent";
import * as echarts from "echarts";

export default {
  name: "example",
  components: {},
  data() {
    return {
      interval: "",
      activeName: "0",
      form: {
        env: "python"
      },
      args: "{}",
      result: "",
      percentage: 0,
      second: 0,
      cmOption: {
        autoCloseBrackets: true,
        tabSize: 4,
        styleActiveLine: true,
        lineNumbers: true,
        line: true,
        mode: "text/x-python",
        theme: "base16-light"
        // keyMap: "emacs"
      },
      code: dedent`# Code here
      def handler(event):
          return "hello world"
      `,
      functionList: [
        {
          name: "image_resize",
          env: "python3",
          memory: "128MB",
          timeout: "60s"
        },
        {
          name: "MLtest",
          env: "python3",
          memory: "128MB",
          timeout: "30s"
        },
        {
          name: "MLtest1",
          env: "python3",
          memory: "128MB",
          timeout: "3s"
        },
        {
          name: "hello",
          env: "python3",
          memory: "128MB",
          timeout: "3s"
        }
      ]
    };
  },
  methods: {
    handleClick() {
      if (this.activeName === "1") {
        this.form = { env: "python" };
      }
    },
    goDetail() {
      // this.form = val;
      this.activeName = "2";
    },
    createFunc() {
      this.activeName = "2";
    },
    invoke() {
      this.result = "hello world";
    },
    getResult() {
      this.activeName = "3";
      // this.dynamicData();
    },
    randomData() {
      this.second += 1;
      const value = Math.random() * 21;
      return {
        name: `${this.second}s`,
        value: [`${this.second}s`, Math.round(value)]
      };
    },
    beforeDestroy() {
      clearInterval(this.interval);
    },
    dynamicData() {
      const data = [];
      for (let i = 0; i < 10; i += 1) {
        data.push(this.randomData());
      }

      const option = {
        title: {
          text: "动态数据"
        },
        tooltip: {
          trigger: "axis",
          // formatter: function (params) {
          //     params = params[0];
          //     var date = new Date(params.name);
          //     return date.getDate() + '/' + (date.getMonth() + 1) + '/' + date.getFullYear()
          // + ' : ' + params.value[1];
          // },
          axisPointer: {
            animation: false
          }
        },
        xAxis: {
          type: "category",
          splitLine: {
            show: false
          }
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
            name: "模拟数据",
            type: "line",
            showSymbol: false,
            hoverAnimation: false,
            data
          }
        ]
      };

      // const myChart = echarts.init(this.$refs.myEchart);
      // window.onresize = myChart.resize;
      // myChart.setOption(
      //     option);
      let flag = false;
      const myChart = echarts.init(this.$refs.myEchart);

      this.interval = setInterval(() => {
        for (let i = 0; i < 5; i += 1) {
          // data.shift();
          data.push(this.randomData());
        }
        this.percentage += 10;
        myChart.setOption(option);
        if (flag === false) {
          flag = true;
          setTimeout(() => {
            clearInterval(this.interval);
          }, 9000); // 停止
        }
      }, 1000); // 启动,func不能使用括号
    }
  }
};
</script>
<style lang="scss">
.cloudComputation {
  .el-tabs__header {
    margin-bottom: 0;
  }
  .bg-gray {
    padding: 20px;
  }
  .codemirror {
    width: 100%;
    height: 100%;
    margin: 0;
    overflow: auto;
  }
  .funcTable {
    margin-bottom: 30px;
    border-spacing: 0;
    tr {
      border-top: 1px solid #ebeef5;
    }
    tr:last-child {
      border-bottom: 1px solid #ebeef5;
    }
    td {
      width: 16.6%;
      padding-left: 10px;
      border-left: 1px solid #ebeef5;
    }
    td:nth-last-child {
      border-right: 1px solid #ebeef5;
    }
  }
  .progress {
    width: 300px;
    margin: 30px auto;
  }
  .invokeBtn {
    margin: 40px auto;
  }
}
iframe {
  width: 100%;
  height: 500px;
  margin-top: 30px;
}
.functionCard {
  margin: 30px auto;
  .el-card__body {
    padding: 0;
    padding-bottom: 20px;
    text-align: center;
  }
  .title {
    margin: 20px;
    margin-bottom: 10px;
  }
  p {
    margin: 10px 20px;
    // margin-bottom: 10px;
  }
  .el-button {
    margin: 10px;
  }
  .insideFloatRight {
    margin: 0;
    padding: 0;
    float: right;
  }
}
</style>
