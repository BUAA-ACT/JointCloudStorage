<template>
  <div>
    <el-table
      ref="multipleTable"
      :data="selectList"
      highlight-current-row
      @current-change="handleCurrentChange"
      tooltip-effect="dark"
      style="width: 100%"
    >
      <el-table-column prop="name" width="50" :label="$t('disasterReco.select')"> </el-table-column>
      <el-table-column prop="type" width="120" :label="$t('disasterReco.type')"> </el-table-column>
      <el-table-column prop="note" :label="$t('disasterReco.explanation')" show-overflow-tooltip> </el-table-column>
    </el-table>
    <el-button type="primary" class="submitBtn" @click="onSubmit">
      {{ $t("disasterReco.next") }}
    </el-button>
  </div>
</template>

<script>
export default {
  data() {
    return {};
  },
  computed: {
    selectList() {
      if (this.$i18n.locale === "en-US") {
        return [
          {
            type: "DFC-12,6,3,3",
            note: this.$t("disasterReco.option1")
          },
          {
            type: "DFC-8,3,3,2",
            note: this.$t("disasterReco.option2")
          },
          {
            type: "RS-6,3,3,2",
            note: this.$t("disasterReco.option3")
          },
          {
            type: "RS-6,4,2,3",
            note: this.$t("disasterReco.option4")
          },
          {
            type: "RS-8,6,2,4",
            note: this.$t("disasterReco.option5")
          },
          {
            type: "RS-9,6,3,3",
            note: this.$t("disasterReco.option6")
          },
          {
            type: "RS-12,6,6,2",
            note: this.$t("disasterReco.option7")
          },
          {
            type: "RS-12,8,4,3",
            note: this.$t("disasterReco.option8")
          }
        ];
      }
      return [
        {
          type: "DFC-12,6,3,3",
          note: this.$t("disasterReco.option1")
        },
        {
          type: "DFC-8,3,3,2",
          note: this.$t("disasterReco.option2")
        },
        {
          type: "RS-6,3,3,2",
          note: this.$t("disasterReco.option3")
        },
        {
          type: "RS-6,4,2,3",
          note: this.$t("disasterReco.option4")
        },
        {
          type: "RS-8,6,2,4",
          note: this.$t("disasterReco.option5")
        },
        {
          type: "RS-9,6,3,3",
          note: this.$t("disasterReco.option6")
        },
        {
          type: "RS-12,6,6,2",
          note: this.$t("disasterReco.option7")
        },
        {
          type: "RS-12,8,4,3",
          note: this.$t("disasterReco.option8")
        }
      ];
    }
  },
  methods: {
    handleCurrentChange() {},
    onSubmit() {
      this.$router.push({ path: "example/selectSolution" });
    }
  }
};
</script>
