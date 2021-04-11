<template>
  <div>
    <h1 class="page-title">第一次全面评测报告</h1>
    <div id="header-1">
      <h3 class="page-header"></h3>
      <p>
        我们将通过一部分有代表性的实例（即日常评测时进行评测的12种实例及其他个别实例）来对评测的结果进行分析，
        观察其CPU、内存、IO、网络四个方面的指标。
        着重分析同一云厂商的不同类型云服务器之间的差异和同样的类型规格下不同云服务厂商之间的差异。另外，关于同一类型中不同实例的比较，可以在评分预测推广中查看。
      </p>
    </div>
    <div id="header-2">
      <h3 class="page-header">评测范围</h3>
      <p>
        一共包含了阿里云的23种实例、华为云的21种实例、腾讯云的17种实例和UCloud的23种实例，共计94种不同的实例，
        这些实例主要来源于通用型/标准型、计算优化型、内存优化型三种类型，同时也包含其他一些最新类型。地域都是北京。
      </p>
    </div>
    <div id="header-3">
      <h3 class="page-header">厂商内部</h3>
      <p>这一部分着重分析同一云厂商的不同类型云服务器之间的差异。</p>
      <ul class="secondheader">
        <li>阿里云</li>
        <p>
          从下图中可以看出，阿里云的不同类型的实例各方面的性能都比较稳定，相差不大。
          总的来说，计算型优于内存型，内存型优于通用型。具体的，计算型在IO、内存、网络三个方面都明显比通用型更好，而内存型在内存性能上比通用型更好。
        </p>
        <img style="width: 60%;" src="@/assets/modules/serviceRecommendation/pc10.png" />
        <li>华为云</li>
        <p>
          从下图中可以看出，华为云在CPU、IO、内存上相差不大，内存性能内存型优于通用型，通用型优于计算型。
          但是，在网络性能方面，差距就比较大。通用型的两个实例的网络性能都远远低于其他三个实例，不到一半，
          这说明通用型的实例并不适合对网络性能需求较大的场景，计算型的网络性能高于内存型。
        </p>
        <img style="width: 60%;" src="@/assets/modules/serviceRecommendation/pc11.png" />
        <li>腾讯云</li>
        <p>
          从下图中可以看出，腾讯云仅仅在内存性能上相差不大，在其他三个指标上都有着较大的差距。
          对于CPU性能，计算型的两个实例明显的高于其他三个实例，在IO性能方面，计算型的实例同样也优于其他实例，
          这种差距在网络性能上体现的更为明显。综合来看，腾讯云的计算型实例在各个方面都是最好的，相对于其他的类型，在性能要求上是更好的选择。
        </p>
        <img style="width: 60%;" src="@/assets/modules/serviceRecommendation/pc12.png" />
        <li>UCloud</li>
        <p>从下图中可以看出，UCloud的三个实例除了在网络性能上内存呢最大的实例更优秀，其他几乎没有差别，这恰恰也体现了UCloud并没有区分不同的类型。</p>
        <img style="width: 60%;" src="@/assets/modules/serviceRecommendation/pc13.png" />
      </ul>
    </div>
    <div id="header-4">
      <h3 class="page-header">厂商之间</h3>
      <p>这一部分着重分析同样的类型规格下不同云服务厂商之间的差异。</p>
      <ul class="secondheader">
        <li>通用型 / 标准型</li>
        <p>
          从下图中可以看出，对于CPU内存比为1:4的通用型而言，总体来看阿里云是最优秀的，主要体现在内存性能和网络性能上的差距。
          具体的，CPU性能方面，华为云和UCloud好于阿里云，阿里云好于腾讯云；IO性能方面，华为云最好，腾讯云最差；
          内存方面，华为云和UCloud又比腾讯云和阿里云差了不少；网络性能方面，阿里云优于UCloud，优于腾讯云，优于华为云。
          不同厂商有着各自好的方面，需要根据不同的需求作出选择。
        </p>
        <img style="width: 60%;" src="@/assets/modules/serviceRecommendation/pc14.png" />
        <li>计算优化型</li>
        <p>
          从下图中可以看出，对于CPU内存比为1:2的计算型而言，CPU性能差距不大，IO性能华为云好于阿里云，好于UCloud，
          好于腾讯云，内存方面仍然是阿里云最佳，网络方面阿里云和华为云好于腾讯云和UCloud。 从性能的角度考虑，阿里云和华为云都是不错的选择。
        </p>
        <img style="width: 60%;" src="@/assets/modules/serviceRecommendation/pc15.png" />
        <li>内存优化型</li>
        <p>
          从下图中可以看出，对于CPU内存比为1:8的内存型而言。对于CPU性能，UCloud和华为云好于阿里云和腾讯云；
          对于IO性能，阿里云和华为云好于UCloud和腾讯云；内存上仍然是阿里云最好，其次是腾讯云；
          网络上，UCloud最佳，腾讯云最差。总的来说，腾讯云相对较差，阿里云比较推荐。
        </p>
        <img style="width: 60%;" src="@/assets/modules/serviceRecommendation/pc16.png" />
      </ul>
      <p>
        通过以上的比较可以看出，不同的厂商有不同的特点，需要根据具体的需求进行选择。总的来说，不考虑价格，阿里云和华为云更值得推荐。UCloud也不错，腾讯云稍微逊色一点。
      </p>
    </div>
    <div>
      评测结果数据：
      <a href="datas/bench_results.xlsx" download="datas/bench_results.xlsx">bench_results.xlsx</a>
    </div>
  </div>
</template>
<script>
export default {
  name: "Benchreport1"
};
</script>
