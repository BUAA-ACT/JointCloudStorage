<template>
  <div :class="{ fixedFooter: fixed, footer: !fixed }">
    <div class="footer-links">
      <el-row v-if="type !== 'mini'">
        <el-col>
          <el-link href="avascript:" target="_blank">关于我们</el-link>
          <el-link href="avascript:" target="_blank">联系我们</el-link>
          <el-link href="avascript:" target="_blank">加入我们</el-link>
          <el-link href="avascript:" target="_blank">服务体系</el-link>
          <el-link href="avascript:" target="_blank">交易保障</el-link>
          <el-link href="avascript:" target="_blank">网站协议</el-link>
          <el-link href="avascript:" target="_blank">支付方式</el-link>
          <el-link href="avascript:" target="_blank">平台广告</el-link>
          <el-link href="avascript:" target="_blank">致谢信息</el-link>
        </el-col>
        <!-- <div class="layui-col-md3">
        <div class="text-center" style="display: inline-block;float: right;text-align: center;">
          <a href="javascript:">
            <i style="font-size: 26px" class="layui-icon layui-icon-login-weibo"></i></a>
          <p style="font-size: 12px;margin-top: 8px">官方微博</p>
        </div>
        <div class="text-center"
          style="display: inline-block;float: right;margin-right: 10px;text-align: center;">
          <a href="javascript:">
            <i style="font-size: 26px" class="layui-icon layui-icon-login-wechat"></i></a>
          <p style="font-size: 12px;margin-top: 8px">官方公众号</p>
        </div>
      </div> -->
      </el-row>
      <p>
        <!-- 建设单位：国防科技大学计算机院并行与分布式处理国防科技重点实验室 -->
        国家重点研发计划——云际计算项目组（仅用于技术验证）
        <!-- 交流合作：（+86）137-0731-2497 -->
      </p>
      <a href="http://beian.miit.gov.cn/">
        ICP备案号：京ICP备2021011189号-1
      </a>
    </div>
  </div>
</template>

<script>
export default {
  name: "Footer",
  props: {
    type: String
  },
  data() {
    return {
      fixed: this.$route.path.toString().indexOf("documents") > -1
    };
  },
  methods: {
    changeFooter() {
      const contentHeight = document.body.scrollHeight;
      const winHeight = window.innerHeight;
      // 当网页正文高度小于可视窗口高度时，为footer添加类fixed-bottom
      const flag = contentHeight <= winHeight;
      this.fixed = flag;
      return flag;
    }
  },
  watch: {
    $route: {
      handler() {
        this.changeFooter();
      },
      immediate: true
    }
  },
  mounted() {
    this.changeFooter();
  }
};
</script>

<style lang="scss">
.footer {
  background-color: black;
  color: #dddddd;
  padding: 30px 0;
  overflow: auto;
}
.fixedFooter {
  background-color: black;
  color: #dddddd;
  padding: 30px 0;
  overflow: auto;
  width: 100%;
  position: absolute;
  bottom: 0;
}
.footer-links {
  width: 86%;
  max-width: 1200px;
  margin: 0 auto;
  .el-link {
    margin: 18px 20px;
    padding: 10px 0;
    color: #dddddd;
    font-size: 18px;
    font-weight: normal;
    line-height: 1em;
  }
  p {
    margin: 0;
  }
}
</style>
