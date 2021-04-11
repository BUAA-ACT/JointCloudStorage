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
            <el-button type="primary" @click="$router.push('requirementRelease')">需求发布</el-button>
          </div>
          <div class="insideFloatRight">
            <el-input class="search" placeholder="请输入您需要的服务" v-model="search">
              <el-button type="primary" slot="append" icon="el-icon-search"></el-button>
            </el-input>
          </div>
        </div>
        <el-divider></el-divider>
        <div class="content-card" v-for="(item, index) in requirementList" :key="index">
          <div class="content-card-head">
            <el-form label-width="100px" label-position="left">
              <el-form-item label="需求名称">
                {{ item.requireName }}
                <el-button type="success" plain class="insideFloatRight">已成交</el-button>
                <el-button type="success" class="insideFloatRight">{{ item.type || "通用型" }}</el-button>
              </el-form-item>
            </el-form>
          </div>
          <el-divider></el-divider>
          <div class="content-card-body">
            <el-form v-model="form" label-width="100px" label-position="left">
              <el-form-item label="基础配置">
                <table>
                  <tr>
                    <td><span>付费模式：</span>{{ item.payType }}</td>
                    <td><span>地域及可用区：</span>{{ item.area }}</td>
                    <td><span>实例：</span>{{ item.requireName }}</td>
                  </tr>
                  <tr>
                    <td><span>购买数量：</span>{{ item.count }}</td>
                    <td><span>镜像：</span>{{ item.mirror }}</td>
                    <td><span>系统盘：</span>{{ item.systemDisk }}</td>
                  </tr>
                  <tr>
                    <td><span>购买时长：</span>{{ item.duration }}</td>
                    <td><span>时间：</span>{{ item.time }}</td>
                    <td><span>用户：</span>{{ item.userName }}</td>
                  </tr>
                </table>
              </el-form-item>
              <el-form-item label="网络和安全组">
                <table>
                  <tr>
                    <td><span>网络 ：</span>{{ item.network }}</td>
                    <td><span>VPC ：</span>{{ item.VPC }}</td>
                    <td><span>交换机 ：</span>{{ item.exchange }}</td>
                  </tr>
                  <tr>
                    <td><span>公网带宽 ：</span>{{ item.bandwidth }}</td>
                    <td><span>安全组 ：</span>{{ item.safety }}</td>
                    <td></td>
                  </tr>
                </table>
              </el-form-item>
            </el-form>
          </div>
        </div>
        <el-pagination background style="text-align:right;margin-right: 10px" layout="prev, pager, next" :total="2"> </el-pagination>
      </el-card>
    </div>
  </div>
</template>
<script>
// import SingleRequirement from "./singleRequirement.vue";
// import ComprehensiveRequirement from "./comprehensiveRequirement.vue";

export default {
  name: "RequirementForm",
  data() {
    return {
      search: "",
      form: {
        resource: "全部"
      },
      requirementList: [
        {
          requireName: "高主频通用型 hfg7 /ecs.hfg5.xlarge（4vCPU 16GiB）",
          type: "通用型",
          payType: "包年包月",
          area: "华东 1 （杭州） / 随机分配",
          count: "1 台",
          mirror: "Alibaba Cloud Linux 2.1903 LTS 64位（安全加固）",
          systemDisk: "ESSD云盘 40GiB ，随实例释放，PL0（单盘IOPS性能上限1万）",
          duration: "5 个月 ",
          time: "2020-11-20 09:11:55",
          userName: "麒麟软件有限公司（企业用户）",
          network: "专有网络",
          VPC: "默认专有网络",
          exchange: "默认交换机",
          bandwidth: "按固定带宽 1Mbps",
          safety: "1). 默认安全组（自定义端口）"
        },
        {
          requireName: "大数据存储型  d2s /ecs.d2s.20xlarge(	80vCPU	352 GiB)",
          type: "存储型",
          payType: "按量付费",
          area: "华南 1（深圳） / 随机分配",
          count: "1 台",
          mirror: "Alibaba Cloud Linux 2.1903 LTS 64位（安全加固）",
          systemDisk: "ESSD云盘 40GiB ，随实例释放，PL0（单盘IOPS性能上限1万）",
          duration: "1个月 ",
          time: "2021-01-20 13:21:25",
          userName: "中国电子科技集团有限公司（企业用户）",
          network: "专有网络",
          VPC: "默认专有网络",
          exchange: "默认交换机",
          bandwidth: "按使用流量 5Mbps",
          safety: "1). 默认安全组（自定义端口）"
        },
        {
          requireName: "高主频通用型 hfc7 / ecs.hfc7.xlarge（4vCPU 8GiB）",
          payType: "包年包月",
          area: "华东 1（上海） / 随机分配",
          count: "1 台",
          mirror: "Windows Server 2019 数据中心版 64位英文版（安全加固）",
          systemDisk: "ESSD云盘 40GiB ，随实例释放，PL0（单盘IOPS性能上限1万）",
          duration: "14 个月",
          time: "2020-12-01 08:21:23",
          userName: "中国软件与技术服务股份有限公司（企业用户）",
          network: "专有网络",
          VPC: "默认专有网络",
          exchange: "默认交换机",
          bandwidth: "按使用流量 25Mbps",
          safety: "1). 默认安全组（自定义端口）"
        },
        {
          requireName: "通用型 g6 / ecs.g6.large（2vCPU 8GiB）",
          payType: "包年包月",
          area: "华北 1 （青岛） / 随机分配",
          count: "1 台",
          mirror: "Ubuntu 18.04 64位（安全加固）",
          systemDisk: "高效云盘 40GiB ，随实例释放",
          duration: "24个月 ",
          time: "2020-12-11 09:31:33",
          userName: "国家超级计算无锡中心（企业用户）",
          network: "专有网络",
          VPC: "默认专有网络",
          exchange: "默认交换机",
          bandwidth: "按固定带宽 1Mbps",
          safety: "1). 默认安全组（自定义端口）"
        },
        {
          requireName: "内存型 r6 / ecs.r6.3xlarge（12vCPU 96GiB）",
          type: "内存型",
          payType: "按量付费",
          area: "美国（硅谷） / 随机分配",
          count: "1 台",
          mirror: "Ubuntu 18.04 64位（安全加固）",
          systemDisk: "高效云盘 40GiB ，随实例释放",
          duration: "6个月",
          time: "2021-01-01 10:51:01",
          userName: "商汤上海人工智能超算中心（企业用户）",
          network: "专有网络",
          VPC: "默认专有网络",
          exchange: "默认交换机",
          bandwidth: "按使用流量 5Mbps",
          safety: "1). 默认安全组（自定义端口）"
        },
        {
          requireName: "计算网络增强型 sn1ne / ecs.sn1ne.xlarge（4vCPU 8GiB）",
          type: "增强型",
          payType: "按量付费",
          area: "中东东部1 可用区A",
          count: "1 台",
          mirror: "Ubuntu 20.04 64位（安全加固）",
          systemDisk: "高效云盘 20GiB ，随实例释放",
          duration: "10个月",
          time: "2021-01-22 15:11:31",
          userName: "John Wick （个人用户）",
          network: "专有网络",
          VPC: "默认专有网络",
          exchange: "默认交换机",
          bandwidth: "按使用流量 5Mbps",
          safety: "1). 默认安全组（自定义端口）"
        }
      ]
    };
  },
  // components: { SingleRequirement, ComprehensiveRequirement },
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
