<template>
  <div class="block">
    <el-carousel :interval="0" arrow="always" :initial-index="4">
      <el-carousel-item v-for="item in 5" :key="item">
        <div style="height:100%" :class="'imgIcon-banner-' + item" @click="goto(item)"></div>
      </el-carousel-item>
    </el-carousel>
    <!--    <div class="manufacturer">-->
    <!--      <el-row :gutter="30">-->
    <!--        <el-col :span="12">-->
    <!--          <el-card class="card" shadow="always">-->
    <!--            <div class="imgIcon imgIcon-xuqiu"></div>-->
    <!--            <p class="fabu">需求发布</p>-->
    <!--            <p class="hoverSee title">-->
    <!--              需求发布-->
    <!--            </p>-->
    <!--            <p class="hoverSee">-->
    <!--              需求发布者通过发布需求，填写提交需求信息，将需求添加至需求大厅展示。同时云际平台智能计算出服务与需求适配度，生成服务推荐方案列表，为需求发布方提供智能购买方案。-->
    <!--            </p>-->
    <!--            <el-button class="fabuBtn" type="text" @click="go('/requirementRelease')">发布</el-button>-->
    <!--          </el-card>-->
    <!--        </el-col>-->
    <!--        <el-col :span="12">-->
    <!--          <el-card class="card" shadow="always">-->
    <!--            <div class="imgIcon imgIcon-fuwu"></div>-->
    <!--            <p class="fabu">服务发布</p>-->
    <!--            <p class="hoverSee title">-->
    <!--              服务发布-->
    <!--            </p>-->
    <!--            <p class="hoverSee">-->
    <!--              服务发布者根据自身资源与数据业务支撑能力，通过填写提交相关服务信息，即可将服务信息在服务大厅中展示。云际平台支持智能磋商、自主选择等方式促成需求与服务达成一致。-->
    <!--            </p>-->
    <!--            <el-button class="fabuBtn" type="text" @click="go('/servicePublication')">发布</el-button>-->
    <!--          </el-card>-->
    <!--        </el-col>-->
    <!--      </el-row>-->
    <!--    </div>-->
    <div class="bg-gray">
      <div class="manufacturer">
        <p class="indexTitle">云际计算服务平台</p>
        <el-row :gutter="2">
          <el-col class="cloudCard" :span="12" v-for="(item, index) in functionList" :key="index">
            <div class="imgIcon" :class="'imgIcon-' + item.icon"></div>
            <p class="title">{{ item.name }}</p>
            <p>{{ item.info }}</p>
            <el-button type="text" class="button" @click="go(item.url)">了解详情</el-button>
          </el-col>
        </el-row>
      </div>
    </div>
    <div class="bg-bubble">
      <div class="manufacturer">
        <p class="indexTitle">解决方案</p>
        <el-row :gutter="30">
          <el-col v-for="(item, index) in solutionList" :span="item.span" :key="'solution' + index">
            <el-card :class="'solutionCard solutionCard-' + item.class" shadow="hover">
              <p class="title">{{ item.title }}</p>
              <p v-html="item.desc"></p>
            </el-card>
          </el-col>
        </el-row>
      </div>
      <div class="manufacturer">
        <p class="indexTitle">合作伙伴</p>
        <el-row :gutter="30">
          <el-col v-for="item in 12" :span="6" :key="'partner' + item">
            <el-card :class="'solutionCard co' + item" shadow="hover"> </el-card>
          </el-col>
        </el-row>
      </div>
    </div>
    <div class="bg-earth">
      <div class="manufacturer">
        <el-row :gutter="60">
          <el-col :span="10">
            <p class="title">万物互联 智能生态 <br />自主协作 合作共赢</p>
            <p>以支撑万物互联的智能计算场景为牵引，支持资源纵横汇聚实现智能匹配，构建安全可信的智能计算新基建。</p>
            <el-button type="text">
              了解全球服务中心 >>
            </el-button>
          </el-col>
          <el-col :span="14">
            <el-tabs class="dataCenter" v-model="activeName">
              <el-tab-pane label="国内服务中心" name="first">北京1</el-tab-pane>
              <el-tab-pane label="全球服务中心" name="second">中国1</el-tab-pane>
            </el-tabs>
          </el-col>
        </el-row>
      </div>
    </div>
    <!-- <div>
      <p style="text-align: center;line-height: 32px;font-size: 32px">厂商信息</p>
      <p style="text-align: center;line-height: 40px;font-size: 18px">Cloud computing services</p>
      <el-row :gutter="12">
        <el-col :span="8" v-for="(item, index) in cardInfo" :key="index">
          <el-card shadow="always">
            <div class="imgIcon" :class="'imgIcon-'+item.icon"></div>
            <p>{{item.info}}</p>
            <el-button type="text" class="button">立即选配</el-button>
          </el-card>
        </el-col>
      </el-row>
    </div> -->
  </div>
</template>

<script>
export default {
  name: "IndexPage",
  data() {
    return {
      activeIndex: "1",
      activeName: "first",
      cardInfo: [
        {
          info:
            "阿里云（www.aliyun.com）创立于2009年，是全球领先的云计算及人工智能科技公司，为200多个国家和地区的企业、开发者和政府机构提供服务。2017年1月阿里云成为奥运会全球指定云服务商。2017年8月阿里巴巴财报数据显示，阿里云付费云计算用户超过100万。阿里云致力于以在线公共服务的方式，提供安全、可靠的...",
          icon: "alibaba"
        },
        {
          info:
            "华为云是华为的云服务品牌，用在线的方式将华为30多年在ICT基础设施领域的技术积累和产品解决方案开放给客户，致力于提供稳定可靠、安全可信、可持续创新的云服务，做智能世界的“黑土地”，推进实现“用得起、用得好、用得放心”的普惠AI。华为云作为底座，为华为全栈全场景AI战略提供强大的算力平台和更易用的开发平台。",
          icon: "huawei"
        },
        {
          info:
            "腾讯云基于QQ、微信、腾讯游戏等海量业务的技术锤炼， 从基础架构到精细化运营，从平台实力到生态能力建设，腾讯云将之整合并面向市场， 使之能够为企业和创业者提供集云计算、云数据、云运营于一体的云端服务体验。",
          icon: "tencent"
        }
      ],
      functionList: [
        // {
        // name: '云际磋商', icon: 'yj_1', info: '根据需求侧特定信息整合供给侧资源, 形成满足特定领域用户需求的云服务资源提供商候选对象列表, 协助用户选择最优化的云服务资源组合。', url: '/serviceRecommendation',
        // },
        {
          name: "虚拟专云",
          icon: "yj_5",
          info:
            "基于虚拟专用云平台管理系统，建立 “工作空间即服务” 的一种新型云际应用服务模式，对传统应用平滑上云，与对云际资源的动态自适应形成有效支撑。",
          url: "/serviceRecommendation"
        },
        {
          name: "云际存储",
          icon: "yj_3",
          info:
            "量化评估任务执行时间异常造成的影响，利用工作集负载感知的高效虚拟机迁移技术，实现虚拟机实时迁移，多维度云际存储优化，为云际应用开发运维提供支撑。",
          url: "/cloudStorage"
        },
        {
          name: "跨云计算",
          icon: "yj_4",
          info:
            "构建面向计算与存储的资源抽象模型，运用应用与平台的解耦方法，实现计算资源的统一管理与聚合，通过对效能的动态评价寻找最满足功能约束的资源分配计算方案。",
          url: "/cloudComputation"
        },
        {
          name: "云际容灾",
          icon: "yj_6",
          info:
            "利用“纠删码”与“多副本”容灾技术，实现数据分散存放，形成有效数据备份；结合数据安全、存储利用率、读写效率等特性，实现低成本/高效率的云际容灾管理。",
          url: "/disasterRecovery"
        },
        {
          name: "云际同步记账",
          icon: "yj_8",
          info:
            "提供分布记账、智能合约等相关支撑机制, 基于区块链技术构建云际分布记账账本，交易完成后向参与方提供包含资源集成、代理服务和资源交易支付等完整记录。",
          url: "/cloudBookkeepingSync"
        },
        {
          name: "云际异步记账",
          icon: "yj_8",
          info:
            "以异步共识等相关机制为支撑，向参与方提供包含数据协作详细结果，顺序和时间戳等完整记录，利用有向无环图(DAG)结构技术构建云际分布记账账本。",
          url: "/cloudBookkeepingAsync"
        },
        {
          name: "云际互联",
          icon: "yj_2",
          info:
            "实现对云际资源与服务之间的互联互通，构建云服务描述与发现标准体系；实现云际服务的快速发现、标准化索引和计量，构建云际服务的可信平台以及价值交换平台。",
          url: "/404"
        },
        {
          name: "云际学习",
          icon: "yj_7",
          info: "结合AI数据集与应用场景，通过云际大数据融合过程实现数据智能学习，实现对问题解决的聚焦、分类、推荐、解决、预测。",
          url: "/404"
        }
      ],
      solutionList: [
        {
          span: 12,
          class: "game",
          title: "区块链",
          desc: "应用于分布式数据存储、点对点传输、共识机制、加密算法等计算机技术的新型应用模式。"
        },
        {
          span: 6,
          class: "finance",
          title: "金融",
          desc: "以灵活多变的云服务框架，满足金融行业的趋势发展。"
        },
        {
          span: 6,
          class: "education",
          title: "教育",
          desc: "助力在线教育行业发展，为高校与教育管理部门提供解决方案。"
        },
        {
          span: 6,
          class: "bigdata",
          title: "大数据",
          desc: "大数据采集共享，实时发布，为用户提供丰富的数据资源。"
        },
        {
          span: 6,
          class: "video",
          title: "数字政府",
          desc: "提供电子政务、互联网+监管、智慧城市、精细化管理等集成服务体系。"
        },
        {
          span: 12,
          class: "development",
          title: "云开发",
          desc: "提供Faas等多种服务及云端支持，免去服务部署，即可实现快速开发与发布。"
        }
      ]
    };
  },
  methods: {
    handleSelect() {},
    handleCommand() {},
    go(url) {
      this.$router.push(url);
    },
    imgUrlChange(icon) {
      return `@/assets/${icon}.png`;
    },
    goto(item) {
      if (item === 5) {
        this.$router.push("/userCenter");
      }
    }
  }
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style lang="scss">
.imgIcon-banner-5 {
  background-image: url(../assets/jcspan.png), url(../assets/banner-bg3.png);
  background-repeat: no-repeat;
  background-size: auto 80%, auto 100%;
  background-position: center;
  background-color: #fff;
  cursor: pointer;
}
</style>
