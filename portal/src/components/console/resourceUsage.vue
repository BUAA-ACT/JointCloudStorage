<template>
  <el-card class="chartCard" style="height: 450px">
    <el-row>
      <el-col :span="24">
        <el-col :span="6" class="center">
          <el-progress type="circle" :percentage="80" style="width:130px;margin-top: 20px"></el-progress>
        </el-col>
        <el-col :span="18" class="center">
          <el-card style="height: 75px;margin-top: 10px;padding:10px">
            <el-col :span="6">
              <el-progress type="circle" :percentage="12" :width="75"></el-progress>
            </el-col>
            <el-col :span="18">
              <el-col :span="8">
                <p>12%</p>
                <p>CPU</p>
              </el-col>
              <el-col :span="8">
                <p>1.49 core</p>
                <p>已使用</p>
              </el-col>
              <el-col :span="8">
                <p>12 core</p>
                <p>总计</p>
              </el-col>
            </el-col>
          </el-card>
          <el-card style="height: 75px;margin-top: 10px;padding:10px">
            <el-col :span="6">
              <el-progress type="circle" :percentage="39" :width="75"></el-progress>
            </el-col>
            <el-col :span="18">
              <el-col :span="8">
                <p>39%</p>
                <p>内存</p>
              </el-col>
              <el-col :span="8">
                <p>9.18 Gi</p>
                <p>已使用</p>
              </el-col>
              <el-col :span="8">
                <p>23.39 Gi</p>
                <p>总计</p>
              </el-col>
            </el-col>
          </el-card>
          <el-card style="height: 75px;margin-top: 10px;padding:10px">
            <el-col :span="6">
              <el-progress type="circle" :percentage="31" :width="75"></el-progress>
            </el-col>
            <el-col :span="18">
              <el-col :span="8">
                <p>31%</p>
                <p>Pods</p>
              </el-col>
              <el-col :span="8">
                <p>102</p>
                <p>已使用</p>
              </el-col>
              <el-col :span="8">
                <p>330</p>
                <p>总计</p>
              </el-col>
            </el-col>
          </el-card>
          <el-card style="height: 75px;margin-top: 10px;padding:10px">
            <el-col :span="6">
              <el-progress type="circle" :percentage="86" :width="75"></el-progress>
            </el-col>
            <el-col :span="18">
              <el-col :span="8">
                <p>86%</p>
                <p>本地存储</p>
              </el-col>
              <el-col :span="8">
                <p>54.12 GB</p>
                <p>已使用</p>
              </el-col>
              <el-col :span="8">
                <p>63 GB</p>
                <p>总计</p>
              </el-col>
            </el-col>
          </el-card>
        </el-col>
      </el-col>
    </el-row>
  </el-card>
</template>

<script>
export default {
  name: "ResourceUsage"
};
</script>

<style lang="scss"></style>
