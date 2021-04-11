<template>
  <div>
    <div class="mBanner mBanner-m6">
      <p class="bannerTitle">
        云际异步记账
        <br />
        <span>安全、灵活、高吞吐量和低延迟相结合</span>
      </p>
    </div>
    <h2>云际异步记账</h2>
    <div class="manufacturer">
      <p class="text">
        万物互联背景下，复杂的异步交互的动态性和不确定性直接决定了数据流通和共识的难度。DAG结构的区块链技术具有更公平，高效和异步的拜占庭容错能力，可以解决当前面临的挑战，建立异步云际记账。实现在没有可靠网络，不保证消息延迟的完全异步环境下，不影响多类型，多来源的主体之间的协作进程(事务推进)，使得在同一时间序列上的事务不需要达到最终确认状态，即可进行下一个事务的处理。随着时间推移，最终达成共识。
        异步云际记账由三个模块组成，包括初始化，异步协作和查询。
        初始化阶段，确定由多个数据服务提供者DSP和数据服务消费者DSC组成协作委员会，DSP上传数据。 异步协作过程，引入新的“Gossip about
        gossip”协议进行随机交流。一番杂乱无章的通信之后，数据不断同步。而后通过投票和提交两阶段协议达成共识。一轮共识后，事件将被排序，提交给整个网络并存储在本地账本中。
        查询相同的分布式分类帐，如果已达成共识，则返回共识结果，共识时间戳和顺序。
      </p>
      <img src="@/assets/modules/cloudBookkeepingAsync/img.gif" alt="" />
    </div>
  </div>
</template>

<script>
export default {
  name: "overView",
  data() {
    return {};
  },
  methods: {}
};
</script>

<style lang="scss" scoped>
.text {
  font-weight: normal;
  line-height: 2rem;
  font-size: 14px;
  color: #000000;
  text-indent: 2rem;
}
.manufacturer img {
  width: 100%;
  margin: auto;
}
</style>
