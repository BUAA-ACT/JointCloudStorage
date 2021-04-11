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
                <p class="title">
                  {{ item.functionName }}
                  <el-button type="text" class="insideFloatRight" @click="goDetail(item)">details</el-button>
                </p>
                <el-divider></el-divider>
                <p>env: {{ item.runEnv }}</p>
                <p>memory: {{ item.memorySize }}</p>
                <p>timeout: {{ item.timeout }}</p>
                <el-popconfirm
                  confirm-button-text="好的"
                  cancel-button-text="不用了"
                  icon="el-icon-info"
                  icon-color="red"
                  title="确认删除此函数么？"
                  @confirm="deleteFunc(item.functionId)"
                >
                  <el-button slot="reference" type="danger" size="small">删除</el-button>
                </el-popconfirm>
              </el-card>
            </el-col>
          </el-row>
          <el-pagenation
            background
            :hide-on-single-page="true"
            @current-change="currentChange"
            :current-page="page"
            layout="prev, pager, next"
            :total="total"
          >
          </el-pagenation>
        </el-tab-pane>
        <el-tab-pane label="创建函数" name="1">
          <div class="bg-gray" v-if="activeName === '1'">
            <el-row :gutter="30">
              <el-col :span="12">
                <el-form v-model="form" ref="form" label-width="150px">
                  <el-form-item label="Function Name:" prop="functionName">
                    <el-input v-model="form.functionName" :maxlength="10"></el-input>
                  </el-form-item>
                  <el-form-item label="Env:">
                    <el-select v-model="form.runEnv" placeholder="请选择">
                      <el-option label="java8" value="java8"></el-option>
                      <el-option label="python3" value="python3"></el-option>
                    </el-select>
                  </el-form-item>
                  <el-form-item label="TimeOut(1-300s):">
                    <el-input v-model="form.timeout" :maxlength="3"></el-input>
                  </el-form-item>
                  <el-form-item label="Memory Size(128-1024MB):">
                    <el-input v-model="form.memorySize" :maxlength="4"></el-input>
                  </el-form-item>
                </el-form>
              </el-col>
              <el-col :span="12">
                <span class="el-form-item__label">Entrypoint:</span>
                <div class="codemirror" v-if="form.runEnv === 'python3'">
                  <codemirror v-model="code" :options="cmOption" />
                </div>
                <div v-if="form.runEnv === 'java8'">
                  <el-upload class="upload-demo" action="#" :http-request="httpRequest" :limit="1" :file-list="fileList">
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
                <td :colspan="3">{{ form.functionName }}</td>
                <td>Created Time:</td>
                <td>{{ form.createTime }}</td>
              </tr>
              <tr>
                <td>Env:</td>
                <td>{{ form.runEnv }}</td>
                <td>Memory Size:</td>
                <td>{{ form.memorySize }}</td>
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
                <el-input type="textarea" rows="6" v-model="jsonObject"></el-input>
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
        <el-tab-pane label="过程监管" name="3">
          <iframe id="iframe" src="http://117.50.133.184:3000/d/MtxJZLoMk/hcloud-backend?orgId=1&refresh=5s" frameborder="0"></iframe>
          <!-- <div class="bg-gray">
            <el-progress class="progress" :text-inside="true"
            :stroke-width="26" :percentage="percentage"></el-progress>
            <div id="myEchart" style="width:1000px; height:600px" ref="myEchart">
            </div>
          </div> -->
        </el-tab-pane>
      </el-tabs>
    </div>
  </div>
</template>
<script>
import dedent from "dedent";
import { codemirror } from "vue-codemirror";

// language
import "codemirror/mode/python/python";

import "codemirror/lib/codemirror.css";

// theme css
import "codemirror/theme/base16-light.css";

// require active-line.js
import "codemirror/addon/selection/active-line";

import echarts from "echarts";
import pagenation from "@/mixin/pagenation";

export default {
  name: "computeExample",
  components: {
    codemirror
  },
  data() {
    return {
      interval: "",
      activeName: "0",
      form: {
        runEnv: "python3"
      },
      jsonObject: "{}",
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
      functionList: [], // 函数列表
      fileList: []
    };
  },
  mixins: [pagenation],
  methods: {
    currentChange(e) {
      console.log(e);
      this.page = e;
      this.getList();
    },
    httpRequest(data) {
      const isJar = data.file.name.indexOf(".jar") === data.file.name.length - 4;

      if (!isJar) {
        this.$message.error("上传文件只能是 jar 格式!");
      } else {
        // 转base64
        this.getBase64(data.file).then(resBase64 => {
          this.fileList = [{ name: data.file.name, code: resBase64.split(",")[1] }]; // 直接拿到base64信息
          this.$message.success("文件上传成功");
        });
      }
    },
    getBase64(file) {
      return new Promise((resolve, reject) => {
        const reader = new FileReader();
        let fileResult = "";
        reader.readAsDataURL(file);
        reader.onload = () => {
          fileResult = reader.result;
        };
        reader.onerror = error => {
          reject(error);
        };
        reader.onloadend = () => {
          resolve(fileResult);
        };
      });
    },
    getList() {
      const query = {
        page: this.page,
        size: this.size
      };
      this.$Api.cloudComputation.listFunctions(query).then(res => {
        this.functionList = res.rows;
        this.total = res.total;
      });
    },
    deleteFunc(id) {
      this.$Api.cloudComputation.deleteFunction(id).then(res => {
        if (res) {
          this.$message.success("操作成功");
          this.getList();
        }
      });
    },
    handleClick() {
      if (this.activeName === "1") {
        this.form = { env: "python" };
      }
    },
    goDetail(val) {
      this.form = val;
      this.activeName = "2";
    },
    createFunc() {
      // 新建函数
      const funcForm = { code: { zipFile: this.form.runEnv === "java8" ? this.fileList[0].code : this.transBase64(this.code) }, ...this.form };
      this.$Api.cloudComputation.createFunction(funcForm).then(res => {
        if (res) {
          this.$message.success("创建函数成功");
          this.goDetail(this.form);
        }
      });
    },
    transBase64(code) {
      return Buffer.from(code).toString("base64");
    },
    invoke() {
      this.$Api.cloudComputation.invokeFunction({ functionName: this.form.functionName, jsonObject: this.jsonObject }).then(res => {
        if (res) {
          console.log(res);
          this.result = "hello world";
        }
      });
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
