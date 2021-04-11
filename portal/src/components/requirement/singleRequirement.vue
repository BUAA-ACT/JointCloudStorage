<template>
  <div>
    <div class="manufacturer">
      <el-card shadow="always" class="requirementForm">
        <p class="indexTitle">单一需求发布 > 云服务计算需求</p>
        <el-form ref="form" :model="form" label-width="100px" label-position="left">
          <el-form-item label="需求基本信息">
            <table>
              <tr>
                <td>类型</td>
                <td>
                  <el-radio-group v-model="form.resource">
                    <el-radio-button label="通用型"></el-radio-button>
                    <el-radio-button label="标准型"></el-radio-button>
                    <el-radio-button label="高IO型"></el-radio-button>
                    <el-radio-button label="计算型"></el-radio-button>
                    <el-radio-button label="内存型"></el-radio-button>
                    <el-radio-button label="大数据型"></el-radio-button>
                    <el-radio-button label="共享型"></el-radio-button>
                    <el-radio-button label="增强型"></el-radio-button>
                  </el-radio-group>
                </td>
              </tr>
              <tr>
                <td><p></p></td>
              </tr>
              <tr>
                <td>名称</td>
                <td>
                  <el-input placeholder="请输入需求名称"></el-input>
                </td>
              </tr>
            </table>
          </el-form-item>
          <el-divider></el-divider>
          <el-form-item label="硬件规格配置">
            <table border="0">
              <tr>
                <td>cpu</td>
                <td>
                  <el-select v-model="form.region" placeholder="选择vCPU">
                    <el-option v-for="(item, index) in cpuList" :key="index" :label="item + ' vCPU'" :value="index"></el-option>
                  </el-select>
                </td>
                <td>内存</td>
                <td>
                  <el-select v-model="form.region" placeholder="选择内存">
                    <el-option v-for="(item, index) in storageList" :key="index" :label="item + ' GiB'" :value="index"></el-option>
                  </el-select>
                </td>
                <td>
                  I/O优化示例
                  <el-tooltip content="">
                    <i class="el-icon-question" effect="light"></i>
                  </el-tooltip>
                </td>
                <td>
                  <el-select v-model="form.region" placeholder="是否支持IPv6">
                    <el-option label="全部" value="all"></el-option>
                    <el-option label="IPv6" value="ipv6"></el-option>
                  </el-select>
                </td>
              </tr>
              <tr>
                <td><p></p></td>
              </tr>
              <tr>
                <td>架构</td>
                <td :colspan="4">
                  <el-button>
                    x86计算
                  </el-button>
                  <el-button>
                    异构计算GPU/FPGA/NPU
                  </el-button>
                  <el-button>
                    弹性裸金属服务器（神龙）
                  </el-button>
                </td>
              </tr>
            </table>
          </el-form-item>
          <el-divider></el-divider>
          <el-form-item label="存储">
            <table>
              <tr>
                <td>系统盘</td>
              </tr>
              <tr>
                <td>
                  <el-select value="essd">
                    <el-option label="ESSD云盘" value="essd"></el-option>
                  </el-select>
                </td>
                <td>
                  <el-tooltip content="容量范围20~500" placement="top" effect="light">
                    <el-input-number v-model="form.gb" :step="1"></el-input-number>
                  </el-tooltip>
                </td>
                <td>
                  {{ "2280" + " IOPS" }}
                </td>
                <td>
                  性能级别
                  <el-tooltip placement="top" effect="light">
                    <div slot="content" style="width:200px">
                      ESSD 云盘容量越大，可供选择的 性能级别 越高（460 GiB 以上可选 PL2，1260 GiB 以上可选
                      PL3），性能级别越高相应的费用也越高，如何合理选择 ESSD 云盘性能级别，
                      <a href="">查看详情></a>
                    </div>
                    <i class="el-icon-question"></i>
                  </el-tooltip>
                  ：
                </td>
                <td>
                  <el-select v-model="form.region">
                    <el-option label="PL0(单盘IOPS性能上限1万)" value="pl0"></el-option>
                    <el-option label="PL1（单盘IOPS性能上限5万）" value="pl1"></el-option>
                  </el-select>
                </td>
                <td>
                  <el-checkbox v-model="form.checked">随实例释放</el-checkbox>
                </td>
              </tr>
              <tr>
                <td>数据盘</td>
              </tr>
              <tr>
                <td colspan="3">
                  <span class="tips"
                    >您已选择<i class="orange"> {{ dataDisk.length }} </i>块盘，还可以选择 <i class="orange">{{ 16 - dataDisk.length }}</i> 块盘</span
                  >
                </td>
              </tr>
              <tr v-for="(item, index) in dataDisk" :key="index">
                <td>
                  <i class="el-icon-minus" style="margin-right:5px" @click="deleteDataDisk"></i>
                  <el-select value="essd">
                    <el-option label="ESSD云盘" value="essd"></el-option>
                  </el-select>
                </td>
                <td>
                  <el-tooltip content="容量范围20~500" placement="top" effect="light">
                    <el-input-number v-model="item.gb" :step="1"></el-input-number>
                  </el-tooltip>
                </td>
                <td>
                  {{ "2280" + " IOPS" }}
                </td>
                <td>
                  性能级别：
                  <el-tooltip placement="top" effect="light">
                    <div slot="content" style="width:200px">
                      ESSD 云盘容量越大，可供选择的 性能级别 越高（460 GiB 以上可选 PL2，1260 GiB 以上可选
                      PL3），性能级别越高相应的费用也越高，如何合理选择 ESSD 云盘性能级别，
                      <a href="">查看详情></a>
                    </div>
                    <i class="el-icon-question"></i>
                  </el-tooltip>
                </td>
                <td>
                  <el-select v-model="item.level">
                    <el-option label="PL0(单盘IOPS性能上限1万)" value="pl0"></el-option>
                    <el-option label="PL1（单盘IOPS性能上限5万）" value="pl1"></el-option>
                  </el-select>
                </td>
                <td>数量：<el-input-number v-model="form.gb" :step="1"></el-input-number></td>
              </tr>
              <tr>
                <td>
                  <el-button @click="addDataDisk">
                    <i class="el-icon-plus"></i>
                    新增一块数据盘
                  </el-button>
                </td>
              </tr>
            </table>
          </el-form-item>
          <el-divider></el-divider>
          <el-form-item label="镜像">
            <table>
              <tr colspan="5">
                <td>
                  <el-button>公共镜像</el-button>
                  <el-button>自定义镜像</el-button>
                  <el-button>共享镜像</el-button>
                  <el-button>镜像市场</el-button>
                  <el-tooltip content="">
                    <i class="el-icon-question" effect="light"></i>
                  </el-tooltip>
                </td>
              </tr>
              <tr>
                <td><p></p></td>
              </tr>
              <tr>
                <td>
                  <el-select>
                    <el-option v-for="(item, index) in mirrorList" :key="index" :label="item" :value="index"></el-option>
                  </el-select>
                </td>
                <td colspan="3">
                  <el-select placeholder="请选择版本"></el-select>
                </td>
                <td>
                  <el-checkbox v-model="form.checked">安全加固</el-checkbox>
                  <el-tooltip
                    effect="light"
                    content="云服务器加载基础安全组件，提供网站漏洞检查、云产品安全配置检查、主机登录异常告警等安全功能，并可以通过云安全中心统一管理。"
                  >
                    <i class="el-icon-question"></i>
                  </el-tooltip>
                </td>
              </tr>
            </table>
          </el-form-item>
          <el-divider></el-divider>
          <el-form-item label="网络">
            <table>
              <tr>
                <td>
                  <el-button>专有网络</el-button>
                  <el-tooltip effect="light" content="">
                    <i class="el-icon-question"></i>
                  </el-tooltip>
                </td>
              </tr>
              <tr>
                <td><p></p></td>
              </tr>
              <tr>
                <td>
                  <el-select placeholder="默认专有网络"> </el-select>
                </td>
              </tr>
            </table>
          </el-form-item>
          <el-divider></el-divider>
          <el-form-item label="公网IP">
            <el-checkbox>分配公网IPv4地址</el-checkbox>
            <span class="tips">系统会分配公网IP，也可采用更加灵活的弹性公网IP方案</span>
          </el-form-item>
          <el-form-item label="带宽计费模式">
            <el-button>按固定带宽</el-button>
            <el-button>按使用流量</el-button>
            <span class="tips">贷款费用合并在ECS实例中收取</span>
          </el-form-item>
          <el-form-item label="带宽值">
            <el-slider v-model="form.progress" :marks="marks" :step="10" show-input show-stops :max="200"> </el-slider>
            <span class="tips">注意：购买云服务器后如需使用公网负载均衡，则需为公网负载均衡单独购买带宽</span>
          </el-form-item>
          <el-divider></el-divider>
          <el-form-item label="付费模式">
            <el-button>包年包月</el-button>
            <el-button>按量付费</el-button>
            <el-button>抢占式实例</el-button>
          </el-form-item>
          <el-divider></el-divider>
          <el-form-item label="地域及可用区">
            <table>
              <tr>
                <td>
                  <el-select>
                    <el-option label="华北2（北京）" value="1"></el-option>
                  </el-select>
                </td>
                <td>
                  <el-radio-group v-model="form.resource">
                    <el-radio-button label="随机分配"></el-radio-button>
                    <el-radio-button label="可用区H"></el-radio-button>
                    <el-radio-button label="可用区G"></el-radio-button>
                    <el-radio-button label="可用区F"></el-radio-button>
                    <el-radio-button label="可用区C"></el-radio-button>
                    <el-radio-button label="可用区E"></el-radio-button>
                    <el-radio-button label="可用区D"></el-radio-button>
                  </el-radio-group>
                </td>
              </tr>
            </table>

            <span class="tips">不同地域的实例之间内网互不相通；选择靠近您客户的地域，可降低网络时延，提高您客户的访问速度。</span>
          </el-form-item>
        </el-form>
        <el-button type="primary" class="submitBtn" @click="onSubmit">确认提交</el-button>
      </el-card>
    </div>
  </div>
</template>
<script>
export default {
  name: "SingleRequirement",
  data() {
    return {
      form: {},
      cpuList: [
        "1",
        "2",
        "4",
        "8",
        "10",
        "12",
        "16",
        "20",
        "24",
        "28",
        "32",
        "40",
        "48",
        "52",
        "56",
        "64",
        "72",
        "80",
        "82",
        "96",
        "104",
        "160",
        "192",
        "208"
      ],
      storageList: [
        "0.5",
        "1",
        "2",
        "4",
        "8",
        "12",
        "15",
        "16",
        "24",
        "30",
        "31",
        "32",
        "46",
        "48",
        "60",
        "62",
        "64",
        "88",
        "92",
        "93",
        "96",
        "112",
        "120",
        "128",
        "160",
        "176",
        "186",
        "192",
        "224",
        "256",
        "288",
        "336",
        "352",
        "372",
        "384",
        "480",
        "512",
        "704",
        "768",
        "960",
        "1536",
        "1920",
        "3072"
      ],
      dataDisk: [],
      mirrorList: ["Alibaba Cloud Linux", "Windows Server", "CentOS", "Ubuntu", "Debian", "SUSE Linux", "OpenSUSE", "CoreOS", "FreeBSD"],
      marks: {
        0: "0Mbps",
        30: "30Mbps",
        120: "120Mbps",
        200: "200Mbps"
      }
    };
  },
  methods: {
    onSubmit() {
      this.$router.push({ path: "/requirementReleaseSuccess" });
    },
    deleteDataDisk(e) {
      console.log(e);
    },
    addDataDisk() {
      if (this.dataDisk.length < 16) {
        this.dataDisk.push({
          gb: 40,
          level: "pl1"
        });
      } else {
        this.$message.warning("数据盘已达到上限");
      }
    }
  }
};
</script>
<style lang="scss">
.requirement .manufacturer {
  padding-bottom: 20px;
  .el-icon-question {
    color: #dddddd;
    margin: 0 5px;
  }
  table {
    border-spacing: 0;
    td {
      white-space: nowrap;
    }
  }
}
</style>
<style lang="scss" scoped>
.bg-gray {
  padding-bottom: 50px;
}
</style>
