<template>
  <div class="storageMove">
    <div class="mBanner mBanner-m4">
      <p class="bannerTitle">
        云际存储
        <br />
        <span>海量，安全，低成本，高可靠的云存储服务</span>
      </p>
    </div>
    <div class="manufacturer">
      <el-tabs v-model="activeName" @tab-click="handleClick">
        <el-tab-pane label="服务状态监控" name="1">
          <SelectPoint v-if="activeName === '1'" @selectedPointSumbit="onSelectedOrigin" />
        </el-tab-pane>
        <!-- <el-tab-pane label="迁移目标选择" name="2">
          <SelectServer v-if="activeName==='2'" title="资源需求分布图"
          @selectedServerSumbit="onSelectServer"/>
        </el-tab-pane> -->
        <el-tab-pane label="存储迁移" name="2">
          <div v-if="activeName === '2'">
            <el-card class="requirementForm">
              <el-row :gutter="30">
                <el-col :span="12">
                  <p>迁移源</p>
                  <div class="radioGroup">
                    <el-radio v-model="radio" label="1" border> 阿里云-北京</el-radio>
                    <el-radio v-model="radio" label="2" border> 百度云-上海</el-radio>
                    <el-radio v-model="radio" label="3" border> 金山云-广州</el-radio>
                  </div>
                </el-col>
                <el-col :span="12">
                  <p>迁移目标</p>
                  <div class="radioGroup">
                    <el-radio v-model="target" label="1" border> 阿里云-北京</el-radio>
                    <el-radio v-model="target" label="2" border> 百度云-上海</el-radio>
                    <el-radio v-model="target" label="3" border> 金山云-广州</el-radio>
                  </div>
                </el-col>
              </el-row>
              <el-button type="primary" class="submitBtn" @click="onSubmit">开始迁移</el-button>
            </el-card>
          </div>
        </el-tab-pane>
        <el-tab-pane label="迁移结果" name="3">
          <ResultDisplay v-if="activeName === '3'" />
        </el-tab-pane>
      </el-tabs>
    </div>
  </div>
</template>
<script>
import ResultDisplay from "./components/resultDisplay.vue";
import SelectPoint from "./components/selectPoint.vue";
// import SelectServer from './components/selectServer.vue';

export default {
  name: "storageMove",
  components: {
    // SelectServer,
    SelectPoint,
    ResultDisplay
  },
  data() {
    return {
      activeName: "1",
      form: {},
      radio: "",
      target: ""
    };
  },

  methods: {
    handleClick() {},
    onSubmit() {
      // loading
      this.activeName = "3";
    },
    onSelectServer(form) {
      console.log(form);
      this.activeName = "3";
    },
    onSelectedOrigin() {
      this.activeName = "2";
    },
    onSelectedTarget() {
      this.activeName = "4";
    }
  }
};
</script>
<style lang="scss">
.storageMove {
  .el-card {
    margin-bottom: 30px;
  }
  .requirementForm {
    p {
      text-indent: 20px;
    }
    .radioGroup {
      // padding: 0 20px;
      margin: 20px;
      .el-radio {
        // display: block;
        width: 100%;
        margin-bottom: 10px;
      }
      .el-radio.is-bordered + .el-radio.is-bordered {
        margin-left: 0;
      }
    }
  }
}
</style>
