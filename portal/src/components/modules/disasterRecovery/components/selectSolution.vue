<template>
  <div>
    <el-form ref="form" class="selectSolution" :model="form" :key="$i18n.locale" :label-width="$i18n.locale !== 'en-US' ? '100px' : '200px'">
      <el-form-item :label="$t('disasterReco.formLabel1')">
        <el-input-number :min="6" :max="128" v-model="form.num"></el-input-number>
        {{ $i18n.locale === "zh-CN" ? "台" : "" }}
        <span class="tips">{{ $t("disasterReco.inputTip1") }}</span>
      </el-form-item>
      <el-form-item :label="$t('disasterReco.formLabel2')">
        <el-select v-model="form.centerNum">
          <el-option label="2" value="2"></el-option>
        </el-select>
        <el-input class="shortInput" :placeholder="'>=3' + $t('disasterReco.must')" v-model="form.centerNum"></el-input>
        <span class="tips">{{ $t("disasterReco.inputTip2") }}</span>
      </el-form-item>
      <el-form-item :label="$t('disasterReco.formLabel3')">
        <el-select v-model="form.centerNum"> <el-option label="2" value="2"></el-option> </el-select>GB
        <el-input class="shortInput" :placeholder="'>=2' + $t('disasterReco.must')" v-model="form.centerNum"></el-input>
        <span class="tips">{{ $t("disasterReco.inputTip3") }}</span>
      </el-form-item>
      <el-form-item :label="$t('disasterReco.formLabel4')">
        <el-select v-model="form.centerNum"> <el-option label="2" value="2"></el-option> </el-select>GB
        <el-input class="shortInput" :placeholder="'>=2' + $t('disasterReco.must')" v-model="form.centerNum"></el-input>
        <span class="tips">{{ $t("disasterReco.inputTip4") }}</span>
      </el-form-item>
      <el-form-item :label="$t('disasterReco.formLabel5')">
        <el-select v-model="form.centerNum"> <el-option label="2" value="2"></el-option> </el-select>GB
        <el-input class="shortInput" :placeholder="'>=2' + $t('disasterReco.must')" v-model="form.centerNum"></el-input>
        <span class="tips">{{ $t("disasterReco.inputTip5") }}</span>
      </el-form-item>
      <el-form-item :label="$t('disasterReco.formLabel6')">
        <el-select v-model="form.centerNum"> <el-option label="2" value="2"></el-option> </el-select>GB
        <el-input class="shortInput" :placeholder="'>=2' + $t('disasterReco.must')" v-model="form.centerNum"></el-input>
        <span class="tips">{{ $t("disasterReco.inputTip6") }}</span>
      </el-form-item>
      <el-form-item :label="$t('disasterReco.formLabel7')">
        <el-select v-model="form.centerNum"> <el-option label="2" value="2"></el-option> </el-select>Mbps
        <el-input class="shortInput" :placeholder="'>=2' + $t('disasterReco.must')" v-model="form.centerNum"></el-input>
        <span class="tips">{{ $t("disasterReco.inputTip7") }}</span>
      </el-form-item>
    </el-form>
    <el-button @click="back">{{ $t("disasterReco.prev") }}</el-button>
    <el-button type="primary" class="submitBtn" @click="onSubmit"> {{ $t("disasterReco.next") }}</el-button>
  </div>
</template>

<script>
export default {
  data() {
    return {
      form: {}
    };
  },
  methods: {
    back() {
      this.$router.go(-1);
    },
    onSubmit() {
      this.$router.push({ path: "caseConfirm" });
    }
  }
};
</script>
<style lang="scss">
.shortInput {
  width: 100px;
}
.selectSolution {
  .el-input-number,
  .el-select {
    margin-right: 10px;
  }
  .shortInput {
    margin-left: 10px;
  }
}
</style>
