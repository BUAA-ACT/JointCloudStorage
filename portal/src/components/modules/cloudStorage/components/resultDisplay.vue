<template>
  <div>
    <el-card>
      <span>迁移源概况</span>
      <el-row class="numMark" :gutter="20">
        <el-col :span="4" v-for="(item, index) in originBrief" :key="index">
          <div :class="index === 0 ? 'blue' : item.type === 'error' && item.num > 0 ? 'red' : 'green'">
            <i class="bigNum">{{ item.num }}</i
            >{{ item.info }}
          </div>
        </el-col>
      </el-row>
      <span>迁移任务概况</span>
      <div>
        <el-row class="numMark" :gutter="20">
          <el-col :span="4" v-for="(item, index) in taskBrief" :key="index">
            <div :class="index === 0 ? 'blue' : item.type === 'error' && item.num > 0 ? 'red' : 'green'">
              <i class="bigNum">{{ item.num }}</i
              >{{ item.info }}
            </div>
          </el-col>
        </el-row>
      </div>
    </el-card>
    <el-card>
      <span>迁移任务列表</span>
      <el-table ref="multipleTable" :data="tableData" highlight-current-row tooltip-effect="dark" style="width: 100%">
        <el-table-column prop="name" label="任务名称" show-overflow-tooltip> </el-table-column>
        <el-table-column prop="origin" width="150" label="迁移源"> </el-table-column>
        <el-table-column prop="target" width="150" label="迁移目标"> </el-table-column>
        <el-table-column prop="startTime" label="开始时间" width="180"> </el-table-column>
        <el-table-column prop="progress" label="当前完成进度" width="200">
          <template slot-scope="scope">
            <el-progress :percentage="scope.row.progress"></el-progress>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="120">
          <template slot-scope="scope">
            <el-button @click.native.prevent="suspend(scope.$index, tableData)" type="text" size="small">
              暂停
            </el-button>
            <el-button @click.native.prevent="deleteRow(scope.$index, tableData)" type="text" size="small">
              取消
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
    <el-card>
      <span>已完成迁徙任务</span>
      <el-table ref="multipleTable" :data="tableData" highlight-current-row tooltip-effect="dark" style="width: 100%">
        <el-table-column prop="name" label="任务名称" show-overflow-tooltip> </el-table-column>
        <el-table-column prop="origin" width="150" label="迁移源"> </el-table-column>
        <el-table-column prop="target" width="150" label="迁移目标"> </el-table-column>
        <el-table-column prop="startTime" label="开始时间" width="180"> </el-table-column>
        <el-table-column prop="endTime" label="结束时间" width="180"> </el-table-column>
        <el-table-column label="操作" width="120">
          <template slot-scope="scope">
            <el-button @click.native.prevent="deleteRow(scope.$index, tableData)" type="text" size="small">
              查看
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<script>
export default {
  name: "resultDisplay",
  data() {
    return {
      tableData: [
        {
          name: "北京服务器: e-a13ha86weg153we ",
          origin: "阿里云-北京",
          target: "百度云-上海",
          startTime: "2020/11/18 07:22:21",
          endTime: "2020/11/19 19:11:02",
          progress: 60
        },
        {
          name: "东京服务器: e-haurawwroi3rg6ar4ge ",
          origin: "阿里云-北京",
          target: "百度云-上海",
          startTime: "2020/10/18 08:20:51",
          endTime: "2020/11/20 14:33:07",
          progress: 60
        }
        //  {
        //   name: '多伦多服务器: o-itfxci562y3rge ',
        //   origin: '阿里云-多伦多',
        //   target: '百度云-上海',
        //   startTime: '2020/10/23 08:27:20',
        //   endTime: '2020/10/23 09:10:07',
        //   progress: 60,
        // }, {
        //   name: '巴黎服务器: b-irtui12e3ra9i3wrt ',
        //   origin: '阿里云-巴黎',
        //   target: '百度云-上海',
        //   startTime: '2020/11/14 09:23:45',
        //   endTime: '2020/11/29 22:41:22',
        //   progress: 60,
        // }, {
        //   name: '挪威服务器: N-sje8ecd4rg7tt5eas ',
        //   origin: '阿里云-挪威',
        //   target: '百度云-上海',
        //   startTime: '2020/11/18 17:17:35',
        //   endTime: '2020/11/19 18:51:07',
        //   progress: 60,
        // }, {
        //   name: '伦敦服务器: L-eus8ef2eqq4dojhv ',
        //   origin: '阿里云-伦敦',
        //   target: '百度云-上海',
        //   startTime: '2020/10/18 09:21:31',
        //   endTime: '2020/10/22 17:01:58',
        //   progress: 60,
        // }
      ],
      originBrief: [
        {
          num: 4,
          info: "个总数"
        },
        {
          num: 2,
          info: "个在线"
        },
        {
          num: 0,
          info: "个离线"
        },
        {
          num: 0,
          info: "个异常",
          type: "error"
        },
        {
          num: 2,
          info: "个迁移中"
        },
        {
          num: 0,
          info: "个已完成"
        }
      ],
      taskBrief: [
        {
          num: 2,
          info: "个总数"
        },
        {
          num: 0,
          info: "个已暂停"
        },
        {
          num: 0,
          info: "个未开始"
        },
        {
          num: 0,
          info: "个出错",
          type: "error"
        },
        {
          num: 2,
          info: "个迁移中"
        },
        {
          num: 0,
          info: "个已完成"
        }
      ]
    };
  },
  methods: {
    deleteRow(index, rows) {
      rows.splice(index, 1);
    },
    suspend() {}
  }
};
</script>

<style lang="scss" scoped>
.el-card .el-table {
  margin: 20px auto 30px auto;
}
.numMark {
  margin: 20px 0;
  .green {
    background-color: #00d9a611;
    border-top: 2px solid #00d9a6;
  }
  .blue {
    background-color: #419eff11;
    border-top: 2px solid #419eff;
  }
  .red {
    background-color: #ff454511;
    border-top: 2px solid #ff4545;
  }
  .el-col div {
    border-radius: 5px;
    font-size: 14px;
    height: 40px;
    text-align: center;
    .bigNum {
      font-size: 26px;
      margin: 0 5px 0 0;
      line-height: 40px;
    }
  }
}
</style>
