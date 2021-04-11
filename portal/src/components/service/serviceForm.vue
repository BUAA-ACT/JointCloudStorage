<template>
  <div class="requirement">
    <div class="sBanner sBanner-s2">
      <p class="bannerTitle">
        服务发布
        <br />
        <span
          >以客户为中心，建立360度客户视图。支持个性化和个人 <br />
          自助服务，各渠道与客户信息平台集成</span
        >
      </p>
    </div>
    <div class="bg-gray">
      <el-tabs class="dataCenter" v-model="activeName">
        <el-tab-pane label="云服务计算发布" name="first"><SourceRelease /></el-tab-pane>
        <el-tab-pane label="云服务数据发布" name="second"><DataRelease /></el-tab-pane>
      </el-tabs>
    </div>
    <!-- <div class="bg-gray">
      <div class="manufacturer">
        <el-card shadow="always" class="requirementForm">
          <p class="indexTitle" >云服务资源发布</p>
          <el-form ref="form" :model="form" label-width="100px">
            <el-form-item label="服务名称">
              <el-input v-model="form.desc"></el-input>
            </el-form-item>
            <el-form-item label="资源类型">
              <el-select v-model="form.region" placeholder="请选择资源类型">
                <el-option label="需求类型" value="shanghai"></el-option>
                <el-option label="需求类型2" value="beijing"></el-option>
              </el-select>
            </el-form-item>
            <el-form-item label="发布资源相关描述">
              <el-input type="textarea" v-model="form.desc"></el-input>
            </el-form-item>
            <el-form-item label="行业类型">
              <el-radio-group v-model="form.resource">
                <el-radio label="地产"></el-radio>
                <el-radio label="电子商务"></el-radio>
              </el-radio-group>
            </el-form-item>
            <el-form-item label="联系人">
              <el-input v-model="form.desc"></el-input>
            </el-form-item>
            <el-form-item label="联系方式">
              <el-input v-model="form.desc"></el-input>
            </el-form-item>
          </el-form>
          <el-button type="primary" class="submitBtn" @click="onSubmit">确认提交</el-button>
        </el-card>
      </div>

    </div> -->
  </div>
</template>
<script>
import SourceRelease from "./sourceRelease.vue";
import DataRelease from "./dataRelease.vue";

export default {
  name: "RequirementForm",
  data() {
    return {
      form: {},
      activeName: "first",
      marks: {
        0: "0Mbps",
        30: "30Mbps",
        120: "120Mbps",
        200: "200Mbps"
      },
      fileList: [],
      dialogImageUrl: "",
      dialogVisible: false,
      disabled: false
    };
  },
  components: { SourceRelease, DataRelease },
  methods: {
    onSubmit() {}
  }
};
</script>
<style lang="scss">
.requirement .dataCenter {
  width: 86%;
  max-width: 1200px;
  margin: 0 auto;
  margin-top: 20px;
  .el-tabs__header {
    padding: 0 20px;
    margin: 0;
  }
}
</style>
<style lang="scss" scoped>
.bg-gray {
  padding-bottom: 50px;
}
</style>
