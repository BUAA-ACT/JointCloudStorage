<template>
  <el-card class="chartCard">
    <p>数据存储总量<i class="bigNum">45.9GB</i></p>
    <el-row :gutter="30">
      <el-col :span="8" class="center">
        <el-progress type="circle" :percentage="80"></el-progress>
        <p>阿里云-北京</p>
      </el-col>
      <el-col :span="8" class="center">
        <el-progress type="circle" color="#30c0b7" :percentage="80"></el-progress>
        <p>百度云-上海</p>
      </el-col>
      <el-col :span="8" class="center">
        <el-progress type="circle" color="#feb23a" :percentage="75"></el-progress>
        <p>金山云-广州</p>
      </el-col>
    </el-row>
  </el-card>
</template>

<script>
export default {
  name: "ServerPercentage"
};
</script>

<style lang="scss"></style>
