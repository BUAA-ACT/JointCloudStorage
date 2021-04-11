<template>
  <div>
    <h1 class="page-title">评分预测推广</h1>
    <div id="header-1">
      <h3 class="page-header"></h3>
      <p>在日常评测中，我们只评测了各种类的实例中小规格的实例，需要对其他规格的实例性能评测结果进行预测推广。</p>
    </div>
    <div id="header-2">
      <h3 class="page-header">CPU性能</h3>
      <p>
        通过评测，我们发现不同云厂商的云服务器CPU性能符合相同的规律。这里列举了华为云上计算型c6和通用型s6的实例进行展示说明，阿里云、腾讯云、UCloud的结果类似。
      </p>
      <img style="width: 49%;" src="@/assets/modules/serviceRecommendation/pc1.png" />
      <img style="width: 49%;" src="@/assets/modules/serviceRecommendation/pc2.png" />
      <p>从图中可以看出以下结论：</p>
      <ol>
        <li>
          对于本次CPU性能的测试，CPU性能结果与内存大小几乎没有关系。
          图中的large.2和large.4分别对应着CPU数目与内存大小的比例为1:2和1:4，对应的两条折线几乎重合。
        </li>
        <li>
          对同一种类的实例而言，单个CPU性能是相当的。 图中的medium、large、xlarge、2xlarge、3xlarge分别对应着CPU数目为1、2、4、8、12，
          对于单线程的评测，评测结果在同一种类的实例中是几乎相同的。
        </li>
        <li>
          对同一种类的实例而言，其多线程评测的结果与其CPU数目成比例，是线性关系。
          注意图中的折线是弯曲的的因为x轴对应的CPU数目并不是等距的，换为散点图如下，可以明显的看出线性关系。
        </li>
      </ol>
      <img style="width: 49%;" src="@/assets/modules/serviceRecommendation/pc3.png" />
      <img style="width: 49%;" src="@/assets/modules/serviceRecommendation/pc4.png" />
      <p>
        在日常评测中，我们只评测了各种类的实例中小规格的实例。
        基于以上结论，我们可以在日常的评测中通过小规格实例的评测结果对同种类的其他实例的评测结果进行预测。 我们构建了<span
          style="font-weight: bold; color: darkred;"
          >线性模型</span
        >进行预测，
        由于对同一种类的实例而言，其多线程评测的结果与其CPU数目成比例，是线性关系，所以这样的线性模型是可信的。在评分的展示中会使用不同颜色进行区分。
      </p>
      <p>
        需要注意的是，在日常评测中CPU评分并不符合线性关系，这是因为同时考虑了单线程和多线程的评测结果，
        但是在预测是是分开预测的，所以之前构建的模型与CPU的评分并不矛盾。
      </p>
    </div>
    <div id="header-3">
      <h3 class="page-header">IO性能</h3>
      <p>对于总的IO性能，通过评测发现了以下结论：</p>
      <ul>
        <li>
          1当内存大小大于等于4GB时：对于阿里云和UCloud两家厂商，总IO性能趋于稳定，
          如下图UCloud实例IO评测结果所示，可以发现，无论CPU数目和内存大小如何变化，
          总的IO性能都很稳定的维持在了75分左右。而对于腾讯云和华为云两家厂商，
          总IO性能并不稳定，存在小幅的波动，同时整体趋势随着CPU数目和内存大小的增加，存在少量的提高，这种提高是很细微的，但并不能忽视。
        </li>
        <img style="width: 60%;" src="@/assets/modules/serviceRecommendation/pc7.png" />
        <li>当内存大小为1GB或2GB的时候，对应的总IO性能相对于第一种情况会明显的下降很多。</li>
      </ul>
      <p>
        在日常评测中，我们评测的小规格的内存大小都是大于等于4GB的。利用第一条结论，对于阿里云和UCloud两家厂商，
        <span style="font-weight: bold; color: #8b0000;">直接</span>将测得的IO性能推广到该种类的其他实例中； 对于腾讯云和华为云两家厂商，构建<span
          style="font-weight: bold; color: darkred;"
          >线性模型</span
        >进行预测。
      </p>
    </div>
    <div id="header-4">
      <h3 class="page-header">内存性能</h3>
      <p>内存性能的评测同样分为单线程评测和多线程评测，这里以阿里云为例，其他厂商类似。</p>
      <ul>
        <li>
          在单线程评测中，可以得出结论：内存性能与CPU数目和内存大小没有明显关系，在同一种类的实例中维持稳定。
          如下图所示，在g6和g6e两个实例规格族中，性能得分稳定的维持在92分；在hfg6实例规格族中，性能的分稳定的维持在101分。
        </li>
        <img style="width: 60%;" src="@/assets/modules/serviceRecommendation/pc8.png" />
        <img style="width: 60%;" src="@/assets/modules/serviceRecommendation/pc9.png" />
        <li>
          对于多线程评测，我们发现内存性能即不会维持稳定，也不满足线性关系。可以得出以下两个结论：
          <ol>
            <li>内存性能虽CPU变化，内存大小对内存性能的影响不大；</li>
            <li>
              内存性能虽CPU数目增加而减弱，且CPU数目较小时，内存性能减弱变化较大；
              CPU数目较大时，内存性能减弱变化较小。这个规律在阿里云的实例中体现的非常明显，从图中呢可以看出，在阿里云中，CPU数目的内存性能的关系近似于二次函数，
              然而在腾讯云、华为云、UCloud的实例中这一规律体现的并不明显。
              究其原因，可能是因为只有阿里云在同一类型中的实例是严格遵守同样的CPU内存比，而其他三家厂商比较随意，在同一类型中的实例也有多种CPU内存比。
            </li>
          </ol>
        </li>
      </ul>
      <img style="width: 49%;" src="@/assets/modules/serviceRecommendation/pc5.png" />
      <img style="width: 49%;" src="@/assets/modules/serviceRecommendation/pc6.png" />
      <p>
        总的来说，对于单线程评测而言，<span style="font-weight: bold; color: darkred;">直接</span>
        将小规格的评测结果推广到其他规格的评测结果。对于多线程评测而言，简单的使用
        <span style="font-weight: bold; color: darkred;">二次模型</span>
        进行拟合预测，虽然在某些种类的实例中并不够准确，但能够准确的反映出下降的趋势，仍然具有参考价值。
      </p>
    </div>
    <div id="header-5">
      <h3 class="page-header">网络性能</h3>
      <p>
        网络性能虽着配置的增大的确有增强的趋势，但其与CPU数目及内存大小没有明显的关系。影响网络性能的因素有很多其他的方面。所以，我们不对网络性能进行预测。
      </p>
    </div>
    <div id="header-6">
      <h3 class="page-header">启动时间</h3>
      <p>启动时间与规格不存在明显的关系，所以不进行预测。总的来说，在四家厂商中，启动时间整体上由慢到快的顺序为UCloud > 腾讯云 > 华为云 > 阿里云。</p>
    </div>
  </div>
</template>
<script>
export default {
  name: "Scorepredict"
};
</script>
