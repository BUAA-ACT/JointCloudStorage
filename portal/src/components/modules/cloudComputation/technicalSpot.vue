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

    <div class="manufacturer">
      <h2>HCloud 技术亮点</h2>
      <el-divider></el-divider>
      <div class="text">
        <p>
          HCloud是一套跨云Serverless计算框架，提供了基本的Serverless计算服务，同时能够让Serverless任务在不同的云平台（如AWS和阿里云），以及不同云服务（平台原生Serverless服务以及虚拟机）之间进行快速动态迁移，在保证用户业务逻辑稳定的前提条件下，降低用户成本，简化用户运维工作。
        </p>
        <h3>
          一·提供Serverless计算服务
        </h3>
        <p>
          Serverless计算，又称Function-as-a-Service（FaaS）、函数计算、无服务计算等，是一种新型的编程模型。用户以函数的形式将自定义的业务逻辑代码上传到Serverless平台，并通过事件驱动，HTTP调用等触发器对业务逻辑进行调用。Serverless具有如下优势：
        </p>
        <h4>1. 降低启动成本</h4>
        <p>
          Serverless服务使用可以免去自己搭建版本管理服务器、数据库服务器等基础设施，仅需要编写业务逻辑代码即可，对初创企业更加友好。部署、运维等其他服务均由Serverless服务商提供
        </p>
        <h4>2. 细粒度计费机制</h4>
        <p>
          Serverless按照函数调用的具体次数和每个函数执行的具体时间（精确到100ms）进行细粒度计费。对比在虚拟机上运行服务，由于虚拟机的粗粒度长期购买机制，即使没有请求处理和计算资源的使用，虚拟机仍将持续计费
        </p>
        <h4>3. 自动扩容</h4>
        <p>
          Serverless得益于其细粒度部署的特点，易于实现自动扩容，当请求吞吐量激增时，能够快速创建出大量的实例进行计算，保证每个实例的计算效率以及整体业务的稳定性。自动扩容的实现由Serverless平台完成，对用户透明无感知
        </p>
        <h3>二· 提供跨平台的Serverless服务</h3>
        <p>
          云平台通常提供原生的Serverless服务，这些服务往往具有不同的接口，不同的写法要求和运行方式。因此，当云用户希望将自己的函数代码部署在多云从而灵活利用各平台优势时，不得不针对每个平台进行代码的适配和分别部署。
        </p>
        <p>
          HCloud向用户提供统一的接口，能够将用户的同一份代码运行在不同的云平台。当用户具有Serverless任务迁移的需求时，能够将任务动态地从一个云迁移到另一个云上，整个过程对服务调用方透明。
        </p>
        <h3>三· 提供跨服务的Serverless架构</h3>
        <p>
          Serverless按量计费的特点，使得在请求吞吐量小时，比起虚拟机而言，能够避免资源空转带来的成本浪费。然而，当请求吞吐量非常大时，Serverless的服务成本较高，而此时使用相同性能的虚拟机处理请求，成本相对使用Serverless更加低廉。
        </p>
        <p>
          因此，当请求吞吐量低时，HCloud能够将Serverless任务运行在原生的Serverless平台（如AWS
          Lambda），利用Serverless按量付费的优势；当请求吞吐量高时，HCloud可以将任务迁移到虚拟机上运行，从而降低成本p>
        </p>
        <p>
          通常情况下，单个Function的调用请求量较少，往往不需要占用所有的虚拟机资源。
          HCloud可以将多个Function的实例一同运行在虚拟机中，从而提升虚拟机整体的资源利用率，最终达到降低每个用户价格的目标。该原理类似于日常生活中的“拼单”。
        </p>
        <p>
          云服务厂商通常提供Spot
          Instance，即可以临时使用的虚拟机，这样的虚拟机通常价格非常便宜，但随时可能被云服务商回收。HCloud提供的Serverless计算模型天然利于在Spot
          Instance上运行，当某台虚拟机被回收时，仅需将任务迁移到其他虚拟机上即可，该过程对用户无感知。（注：目前功能完善中）
        </p>
      </div>

      <!-- 集成北大、北航的现有工作（注：目前功能完善中）

       北大（云际磋商）：当HCloud有扩容需求（如请求量突增，需要将任务从原生Serverless服  务迁移至虚拟机，
       或者需要更多的虚拟机来支持计算）时，能够使用北大的资源推荐API，从 而获得一个合适的虚拟机。

       北航（云际存储）：HCloud能够支持云际存储，给用户使用。当用户的自定义业务逻辑代码  中有使用存储的需求时，可以通过云际存储的API来使用这样的存储服务。 -->
    </div>
  </div>
</template>
<script>
export default {
  name: "ServiceMark",
  data() {
    return {
      evaluationTime: "20200824",
      selectedList: [],
      tableData: [
        {
          date: "2016-05-03",
          name: "王小虎",
          address: "上海市普陀区金沙江路 1518 弄"
        },
        {
          date: "2016-05-03",
          name: "王小虎",
          address: "上海市普陀区金沙江路 1518 弄"
        }
      ]
    };
  },
  watch: {
    evaluationTime() {
      // console.log(e);
    }
  },
  methods: {
    handleSelectionChange(e) {
      // console.log(e);
      this.selectedList = [...e];
    },
    compare() {
      if (this.selectedList.length < 2) {
        this.$message({
          message: "请选择至少两条进行对比",
          type: "warning"
        });
      } else {
        // console.log('compare');
      }
    }
  }
};
</script>
<style lang="scss">
.text {
  p {
    text-indent: 2em;
    line-height: 2em;
    -webkit-margin-before: 0em;
    margin-block-start: 0em;
    -webkit-margin-after: 0em;
    margin-block-end: 0em;
  }
}
</style>
