<template>
  <div class="manufacturer">
    <ul>
      <li>
        <el-link href="#/documents/disasterRecoRequire" target="_blank">{{ $t("disasterReco.documentTitle1") }}</el-link>
      </li>
      <li>
        <el-link href="#/documents/disasterRecoNetwork" target="_blank">{{ $t("disasterReco.documentTitle2") }}</el-link>
      </li>
      <li>
        <el-link href="#/documents/disasterRecoProcess" target="_blank">{{ $t("disasterReco.documentTitle3") }}</el-link>
      </li>
      <li>
        <el-link href="#/documents/disasterRecoManual" target="_blank">{{ $t("disasterReco.documentTitle4") }}</el-link>
      </li>
    </ul>
  </div>
</template>
<script>
export default {
  name: "disasterRecoDocument"
};
</script>
<style lang="scss" scoped>
li {
  line-height: 2rem;
}
</style>
