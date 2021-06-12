<template>
  <div class="hello">
    <img alt="Vue logo" src="../assets/logo.png">
    <h1>{{ msg }}</h1>
    <h2> Joint Cloud Storage </h2>
    <el-row type="flex" justify="center">
      <el-col :span="20" >
        <p>
          Joint Cloud Storage 是一个分布式云际存储系统。多个云服务访问方式统一，实现互操作性。
          <br>在此基础上，通过统一定义的接口互相协作，形成云际存储系统。
        </p>
      </el-col>
    </el-row>
  </div>
</template>

<script>
export default {
  name: 'HelloWorld',
  props: {
    msg: String
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
.hello {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  margin-top: 60px;
}

h3 {
  margin: 40px 0 0;
}
ul {
  list-style-type: none;
  padding: 0;
}
li {
  display: inline-block;
  margin: 0 10px;
}
a {
  color: #42b983;
}
</style>
