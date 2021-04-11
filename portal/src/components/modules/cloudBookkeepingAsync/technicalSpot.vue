<template>
  <div>
    <div class="mBanner mBanner-m6">
      <p class="bannerTitle">
        云际异步记账
        <br />
        <span>安全、灵活、高吞吐量和低延迟相结合</span>
      </p>
    </div>

    <div class="manufacturer">
      <h2>技术亮点</h2>
      <el-divider></el-divider>
      <div class="text">
        <p>
          区块链技术凭借其优势和特性可以在不可信环境中低成本建立信任的新协作模式。但是，区块链目前面临可扩展性的挑战。有向无环图（DAG）是一种不同于区块链的分布式账本技术。将区块的链式存储结构转变为网状拓扑结构。使得事务上链的操作可以并发执行，并且达成共识的过程无需以高能耗为前提的挖矿。
        </p>
        <p>
          Hashgraph是一种DAG结构的异步不确定的算法。异步指除了网络消息的最终交付外，对网络消息延迟没有保证，流程的协同完全由消息传递事件驱动。不确定性意味着指在
          r 轮 (r有可能增长至无穷大)
          后，非错误过程的未确定的概率趋近于零。与其他共识算法相比，Hashgraph共识的达成是在本地进行，不需要额外的通信开销。在可扩展性、安全性、效率和共识达成灵活程度上，Hashgraph都有很大突破。
        </p>
        <p>
          我们在Hedera共识服务上构建异步云际记账。异步记账流程如图所示。首先，Alice生成数据并创建一个包含消息和主题的事务。TopicID被附加到消息并发送到测试网。主题允许具有相同topicID的消息分组在一起。Alice可以删除主题，也可以指定协作成员以向该主题提交消息。该消息包括有关交易的详细信息。使用topicID，可以设置数据共享范围，以便只有知道topicID的成员才能看到数据。这提高了数据的隐私性和安全性。Testnet将验证topicID是否有效，并向Alice返回预检查的结果。如果有效，则将其打包为事件。测试网会将有关该事件的消息八卦传播给网络中的其他人。否则，该消息将被丢弃。Bob是Alice的协作者，将其交易提交到相同的topicID。他还可以提交到其他主题。假设topicID相同，事件经历投票和提交的两阶段，最终确定达成共识。
        </p>
        <p>
          •投票。将协作过程分为多个轮次。当新事件传播到下一个成员时，该成员添加见证消息。见证消息用于管理创建轮次和接收轮次。某一轮次创建的第一个事件x，称为见证人，该轮次记为r。例如事件x是否是知名见证人，需要由后面轮次的见证人投票决定。第r
          + 1轮的见证人构成投票委员会的成员。可见意味着则投YES，否则为NO。
        </p>
        <p>
          •提交。如果第r + 2轮的见证人能够强可见第r +
          1轮的见证人，则收集该见证人的选票。当票数超过2/3，则确定该事件是知名见证人。没有收集到足够的选票，会继续进行。确定事件是否是知名见证人不会永远持续下去。第r
          + n轮（n通常为10）还未确定，则引入硬币轮次。当所有或绝大多数知名见证人可见事件时，事件将达到最终确认状态。
          然后共识被持久保存到镜像节点以显示给其他节点。此时会生成一条记录，其中包括消息，主题，序列号和共识时间戳。该时间戳一定是最终的时间戳，通常在几秒钟内即可达成共识。镜像节点会侦听给定主题的记录。由于镜像节点接收来自Testnet的所有信息，因此了解事务及其一致性顺序，共识时间戳。它还可以构造状态证明，向第三方证明某一主题接收到的消息的确切列表，以及以何种顺序和时间戳接收的消息。镜像节点看到新消息后发布给主题订阅者Alice和Bob。
        </p>
        <img src="@/assets/modules/cloudBookkeepingAsync/img1.png" alt="" />
        <p style="text-align:center">异步云际记账流程</p>
      </div>
    </div>
  </div>
</template>
<script>
export default {
  name: "technicalSpot",
  data() {
    return {};
  },
  watch: {
    evaluationTime() {
      // console.log(e);
    }
  },
  methods: {}
};
</script>
<style lang="scss">
.text {
  p {
    -webkit-margin-before: 0em;
    margin-block-start: 0em;
    -webkit-margin-after: 0em;
    margin-block-end: 0em;
    font-size: 0.8rem;
    font-weight: normal;
    line-height: 2rem;
    font-size: 14px;
    color: #000000;
    text-indent: 2rem;
  }
  img {
    width: 85%;
    display: block;
    margin: 10px auto;
  }
}
</style>
