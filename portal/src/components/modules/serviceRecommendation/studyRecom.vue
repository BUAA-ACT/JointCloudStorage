<template>
  <div>
    <div class="mBanner mBanner-m2">
      <p class="bannerTitle">
        机器学习
        <br />
        <span
          >通过得到的数据集进行训练得到一个模型，<br />
          模拟或实现人类的学习行为，以获取云服务预测的技能。
        </span>
      </p>
    </div>
    <h2>机器学习推荐</h2>

    <div class="manufacturer">
      <p class="lineHeight20">
        我们在阿里云、华为云、腾讯云、UCloud四家厂商中常用的实例上运行Hadoop和Spark应用程序来获取大量的性能数据。
        Hadoop和Spark工作负载提取自HiBench和spark-perf。<br />
        具体情况及数据集见：<el-link type="primary" href="https://github.com/jind5182/lumos-dataset" target="_blank"
          >https://github.com/jind5182/lumos-dataset</el-link
        ><br />
        我们通过得到的数据集进行训练得到一个模型。
        您可以在我们支持的实例集合中选取几个实例，使用您自己的程序在这几个实例上运行，收集运行过程中的性能数据。并将收集到的数据提交给我们，我们将反馈对于你的程序最适合的实例。
      </p>
      <div class="bg-gray">
        <p>提交您的性能数据</p>
        <el-form>
          <el-form-item label="性能数据">
            <el-upload
              class="upload-demo"
              action="https://jsonplaceholder.typicode.com/posts/"
              :on-preview="handlePreview"
              :on-remove="handleRemove"
              :before-remove="beforeRemove"
              :limit="1"
              :on-exceed="handleExceed"
              :file-list="fileList"
            >
              <el-button size="small" type="primary">点击上传</el-button>
            </el-upload>
          </el-form-item>
        </el-form>
        <span>
          注意：数据格式具体见提供的github
        </span>
      </div>
    </div>
  </div>
</template>
<script>
export default {
  name: "studyRecom",
  data() {
    return {
      fileList: []
    };
  },
  methods: {
    handlePreview() {},
    handleRemove() {},
    beforeRemove() {},
    handleExceed() {}
  }
};
</script>
<style lang="scss" scoped>
.bg-gray {
  text-align: center;
  padding: 20px;
  p {
    font-weight: bold;
  }
  .el-form {
    width: 300px;
    margin: 0 auto;
  }
  span {
    font-size: 14px;
    color: #999999;
  }
}
</style>
