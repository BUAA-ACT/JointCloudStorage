<template>
  <div class="requirement bg-gray">
    <div class="dataCenter">
      <el-card class="filter-field">
        <div class="filter-item">
          需求类型：
          <el-radio-group v-model="form.resource">
            <el-radio-button label="全部"></el-radio-button>
            <el-radio-button label="通用型"></el-radio-button>
            <el-radio-button label="标准型"></el-radio-button>
            <el-radio-button label="高IO型"></el-radio-button>
            <el-radio-button label="计算型"></el-radio-button>
            <el-radio-button label="内存型"></el-radio-button>
            <el-radio-button label="大数据型"></el-radio-button>
            <el-radio-button label="共享型"></el-radio-button>
            <el-radio-button label="增强型"></el-radio-button>
          </el-radio-group>
        </div>
        <el-divider></el-divider>
        <div class="filter-item">
          可用区域：
          <el-radio-group v-model="form.resource">
            <el-radio-button label="全部"></el-radio-button>
            <el-radio-button label="北京"></el-radio-button>
            <el-radio-button label="上海"></el-radio-button>
            <el-radio-button label="杭州"></el-radio-button>
            <el-radio-button label="广州"></el-radio-button>
            <el-radio-button label="长沙"></el-radio-button>
            <el-radio-button label="武汉"></el-radio-button>
            <el-radio-button label="天津"></el-radio-button>
          </el-radio-group>
          <el-divider></el-divider>
        </div>
        <div class="filter-item">
          成交状态：
          <el-radio-group v-model="form.resource">
            <el-radio-button label="全部"></el-radio-button>
            <el-radio-button label="已成交"></el-radio-button>
            <el-radio-button label="进行中"></el-radio-button>
            <el-radio-button label="未成交"></el-radio-button>
          </el-radio-group>
        </div>
      </el-card>
      <el-card class="requireCard">
        <div class="filter-head">
          <div class="filter-item">
            综合排序
          </div>
          <div class="filter-item">排序 <i class="el-icon-d-caret"></i></div>
          <div class="insideFloatRight">
            <el-button type="primary" @click="$router.push('servicePublication')">服务发布</el-button>
          </div>
          <div class="insideFloatRight">
            <el-input class="search" placeholder="请输入您需要的服务" v-model="search">
              <el-button type="primary" slot="append" icon="el-icon-search"></el-button>
            </el-input>
          </div>
        </div>
        <el-divider></el-divider>
        <div class="content-card" v-for="(item, index) in serviceList" :key="index">
          <div class="content-card-head">
            {{ item.serviceName }}
            <el-button type="success" class="insideFloatRight">{{ item.type }}</el-button>
          </div>
          <el-divider></el-divider>
          <div class="content-card-body">
            <el-form v-model="form" label-width="100px" label-position="left">
              <el-form-item label="计算服务">
                <table>
                  <tr>
                    <td><span>服务名称：</span>{{ item.countName }}</td>
                  </tr>
                  <tr>
                    <td><span>描述：</span>{{ item.countDisc }}</td>
                  </tr>
                </table>
              </el-form-item>
              <el-form-item label="数据服务">
                <table>
                  <tr>
                    <td><span>服务名称：</span>{{ item.dataName }}</td>
                  </tr>
                  <tr>
                    <td><span>描述：</span>{{ item.dataDisc }}</td>
                  </tr>
                </table>
              </el-form-item>
            </el-form>
          </div>
        </div>
        <el-pagination background style="text-align:right;margin-right: 10px" layout="prev, pager, next" :total="6"> </el-pagination>
      </el-card>
    </div>
  </div>
</template>
<script>
export default {
  name: "RequirementForm",
  data() {
    return {
      search: "",
      form: {},
      serviceList: [
        {
          serviceName: "乌镇之光（桐乡）超算中心",
          type: "通用型",
          countName: "互联网、大数据、人工智能和实体经济存储计算",
          countDisc: "提供阿里云与百度云等多种云际存储，提供VM虚拟机部署K8S服务部署等多种服务管理。",
          dataName: "2019年至今浙江省医疗保健消费数据",
          dataDisc: "提供海量医疗保健消费数据知识图谱关联，连接医疗服务网点各大数据网点确保数据实时更新。"
        },
        {
          serviceName: "腾讯云计算（北京）有限责任公司",
          type: "通用型",
          countName: "腾讯云大数据存储计算",
          countDisc: "提供腾讯云云际存储，提供VM虚拟机部署K8S服务部署等多种服务管理。",
          dataName: "服务器代运维数据",
          dataDisc: "运维团队对多云上数据迁移、解决方案、架构设计、数据库设置、网站程序设置、远程监测数据获取"
        },
        {
          serviceName: "华为技术有限公司（华为云）",
          type: "通用型",
          countName: "华为云大数据存储计算",
          countDisc: "提供华为云等多种云际存储，提供VM虚拟机部署K8S服务部署等多种服务管理。",
          dataName: "华为云部分上云数据",
          dataDisc: "信息安全等级测评，等保测评，等保二级，等保三级数据。"
        },
        {
          serviceName: "优刻得科技股份有限公司（UCloud)",
          type: "通用型",
          countName: "UCloud-基于虚拟化技术的主机服务",
          countDisc: "提供一系列专业解决方案，包括计算资源、存储资源和网络资源等企业必须的基础IT架构服务，满足互联网研发团队在不同场景下的各类需求。",
          dataName: "标准化网站环境配置",
          dataDisc: "高效快速环境搭建，限时折扣助力轻松上云。"
        },
        {
          serviceName: "阿里云计算中心",
          type: "通用型",
          countName: "阿里云大数据存储计算",
          countDisc: "提供阿里云与百度云等多种云际存储，提供VM虚拟机部署K8S服务部署等多种服务管理。",
          dataName: "ECS故障排查服务数据统计 ",
          dataDisc: "IIS、apache,nginx、tomcat等多种服务器＋数据库问题排查数据统计。"
        },
        {
          serviceName: "之江实验室大数据中心",
          type: "通用型",
          countName: "大数据融合",
          countDisc: "依托阿里云＋华为云＋腾讯云等存储资源，针对不同企业应用场景，支持跨云融合。",
          dataName: "融合处理＋边缘 ",
          dataDisc: "2020年浙江省智慧交通对外公布统计数据。"
        },
        {
          serviceName: "北京金山云网络技术有限公司（金山云）",
          type: "通用型",
          countName: "金山云存储服务",
          countDisc: "构建了完备的云计算基础架构和运营体系，通过与人工智能、大数据、物联网、区块链、边缘计算、AR/VR等先进技术有机结合。",
          dataName: "DataCloud-大数据云平台",
          dataDisc: "提供数据汇聚、数据开发、智能调度、数据管理、数据分析与服务等一站式大数据Web化处理环境，帮助客户快速发现数据价值。"
        },
        {
          serviceName: "国家超级计算无锡中心",
          type: "通用型",
          countName: "跨云存储＋融合计算",
          countDisc: "可提供数十PB的高性能的存储系统。高水平的技术人才队伍为使用计算和存储资源提供全面的技术支持和服务。",
          dataName: "云计算与大数据处理",
          dataDisc: "运用云计算技术，由云模式下的几百万个处理器提供强大的计算能力，将对大数据产业提供强大的支撑，并对其整体发展起到极大的促进作用。"
        },
        {
          serviceName: "商汤上海人工智能超算中心",
          type: "通用型",
          countName: "高性能异构计算",
          countDisc:
            "面向当代AI技术发展，为计算机视觉和机器学习提供了丰富的高性能计算体系。该体系支持包括X86、GPU、ARM、FPGA, 和各种新型AI芯片在内的多种计算架构。",
          dataName: "AI 高性能存储",
          dataDisc: "对数据传输，读写流程，数据缓存做了专门的优化，支持Al场景下大规模数据海量存储以及高速读取"
        }
      ]
    };
  },
  methods: {
    onSubmit() {}
  }
};
</script>
<style lang="scss">
.requirement .dataCenter {
  width: 86%;
  max-width: 1200px;
  margin: 0 auto;
  margin-top: 20px;
  .search {
    margin-bottom: 20px;
  }
  .insideFloatRight {
    margin-left: 20px;
  }
  .filter-field {
    margin-bottom: 20px;
    text-align: left;
    font-size: 0.8rem;
    .el-divider--horizontal {
      margin: 10px 0;
    }
  }
  .filter-head {
    padding: 0 20px;
    height: 48px;
    line-height: 32px;
    .filter-item {
      float: left;
      display: block;
      font-size: 0.9rem;
      color: gray;
      margin-right: 20px;
    }
    .el-form-item {
      margin-bottom: 0 !important;
    }
  }
}
.requireCard {
  margin-bottom: 20px;
  //

  .el-divider--horizontal {
    margin: 5px 0;
  }
  .el-card__body {
    padding: 20px 0;
    width: 100%;
    height: 100%;
    position: relative;
    p {
      margin-top: 0;
    }
    .el-table th {
      background: #ffffff;
    }
    .leftTitle {
      width: 25px;
      height: 100%;
      background: #409eff;
      color: #fff;
      position: absolute;
      padding-top: 10%;
      font-size: 0.8em;
      line-height: 1.9rem;
      // padding-left: 5px;
      top: 0;
      left: 0;
    }
    .deal {
      background: #0dbeae;
    }
    .content-card {
      margin: 20px;
      padding: 10px 20px;
      // margin-bottom: 20px;
      background: #f9f9f9;
      .content-card-head {
        text-align: left;
        height: 40px;
        line-height: 40px;
        .el-form-item__label,
        .el-form-item__content {
          font-size: 1rem;
          color: black;
          font-weight: 600;
        }
      }
      > p {
        margin: 0;
        text-align: left;
        padding: 10px 20px;
        border-bottom: 1px dashed #dddddd;
      }
      > div {
        .el-form {
          // padding-bottom: 10px;
        }
        .el-form-item {
          border-bottom: 1px dashed #dddddd;
          text-align: left;
          margin-bottom: 5px;
          .tips {
            font-size: 0.8em;
          }
        }
        .el-form-item:last-child {
          border: 0;
          margin-bottom: 10px;
        }
        .el-radio-button__orig-radio:checked + .el-radio-button__inner {
          background-color: #9ac1ea7d;
          color: #666666;
        }
        table {
          text-align: left;
          font-size: 12px;
          td {
            padding-left: 10px;
            vertical-align: baseline;
            line-height: 25px;
            span {
              font-weight: bold;
            }
          }
        }
      }
    }
  }
}
</style>
