<template>
  <div class="manufacturer">
    <el-card shadow="always" class="requirementForm">
      <p class="indexTitle">数据发布信息</p>
      <el-form ref="form" :model="form" label-width="100px">
        <el-form-item label="数据名称">
          <el-input v-model="form.desc"></el-input>
        </el-form-item>
        <el-form-item label="数据标题">
          <el-input v-model="form.desc"></el-input>
        </el-form-item>
        <el-form-item label="数据图标">
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
            <div slot="tip" class="el-upload__tip">只能上传jpg/png文件，且不超过500kb</div>
          </el-upload>
        </el-form-item>
        <el-form-item label="行业类型">
          <el-upload action="#" list-type="picture-card" :auto-upload="false">
            <i slot="default" class="el-icon-plus"></i>
            <div slot="file" slot-scope="{ file }">
              <img class="el-upload-list__item-thumbnail" :src="file.url" alt="" />
              <span class="el-upload-list__item-actions">
                <span class="el-upload-list__item-preview" @click="handlePictureCardPreview(file)">
                  <i class="el-icon-zoom-in"></i>
                </span>
                <span v-if="!disabled" class="el-upload-list__item-delete" @click="handleDownload(file)">
                  <i class="el-icon-download"></i>
                </span>
                <span v-if="!disabled" class="el-upload-list__item-delete" @click="handleRemove(file)">
                  <i class="el-icon-delete"></i>
                </span>
              </span>
            </div>
          </el-upload>
        </el-form-item>
        <el-form-item label="数据用途">
          <el-input v-model="form.desc"></el-input>
        </el-form-item>
        <el-form-item label="定价策略">
          <el-radio-group v-model="form.resource">
            <el-radio-button label="免费"></el-radio-button>
            <el-radio-button label="120元/次 按次计费"></el-radio-button>
            <el-radio-button label="270元/年 包年付费"></el-radio-button>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="描述">
          <el-input type="textarea" v-model="form.desc"></el-input>
        </el-form-item>
        <el-form-item label="上传条款">
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
        <el-form-item label="样本数据">
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
      <el-button class="submitBtn" type="primary" @click="onSubmit">确认提交</el-button>
    </el-card>
  </div>
</template>
<script>
export default {
  name: "dataRelease",
  data() {
    return {
      form: {},
      fileList: []
    };
  },
  methods: {
    onSubmit() {
      this.$router.push({ path: "/servicePublicationSuccess" });
    },
    handleRemove(file, fileList) {
      console.log(file, fileList);
    },
    handlePreview(file) {
      console.log(file);
    },
    handleExceed(files, fileList) {
      this.$message.warning(`当前限制选择 1 个文件，本次选择了 ${files.length} 个文件，共选择了 ${files.length + fileList.length} 个文件`);
    },
    beforeRemove(file) {
      return this.$confirm(`确定移除 ${file.name}？`);
    },
    handlePictureCardPreview(file) {
      this.dialogImageUrl = file.url;
      this.dialogVisible = true;
    },
    handleDownload(file) {
      console.log(file);
    }
  }
};
</script>
<style lang="scss">
@keyframes hover-blue {
  0% {
    -webkit-filter: brightness(90%);
    background-color: #419eff;
    opacity: 0.1;
  }
  50% {
    -webkit-filter: brightness(90%);
    background-color: #419eff;
    opacity: 0.3;
  }
  100% {
    /*css3滤镜亮度百分比*/
    -webkit-filter: brightness(100%);
    background-color: #419eff;
    opacity: 1;
  }
}
.comprehensiveCase {
  padding: 0 20px;
  margin-bottom: 20px;
  > p {
    text-align: center;
    font-weight: 400;
    color: #000000;
  }
  img {
    width: 100%;
  }
  .el-card {
    min-height: 200px;
    .el-card__body {
      padding: 20px !important;
      overflow: auto;
    }
    p,
    a {
      font-size: 14px;
      color: #666666;
      line-height: 1.8em;
    }
    .title {
      font-size: 16px;
      font-weight: 400;
    }
    .title,
    a {
      color: #333333;
    }
    a {
      float: right;
    }
  }
  .el-card:hover {
    -webkit-animation-name: hover-blue;
    animation-name: hover-blue;
    -webkit-animation-duration: 0.5s;
    animation-duration: 0.5s;
    -webkit-animation-timing-function: ease-in;
    animation-timing-function: ease-in;
    -webkit-animation-iteration-count: 1;
    animation-iteration-count: 1;
    background-color: #419eff;
    p,
    a {
      color: #ffffff;
    }
  }
}
#myEchart {
  width: 800px;
  height: 500px;
  margin: 0 auto;
}
</style>
