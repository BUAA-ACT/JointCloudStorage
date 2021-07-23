<template>
  <el-card>
    <el-steps :active="curStep" finish-status="success">
      <el-step v-for="index in finishStep + 1" :title="`步骤 ${index}`" :key="String(index)"></el-step>
    </el-steps>
    <el-form v-if="curStep === 0" label-position="top">
      <el-form-item label="请选择存储方式" size="large">
        <el-radio v-model="curPlan.StorageMode" label="Replica">多副本</el-radio>
        <el-radio v-model="curPlan.StorageMode" label="EC">纠删码</el-radio>
      </el-form-item>
    </el-form>
    <!--    多副本模式-->
    <h3 v-if="curPlan.StorageMode === 'Replica'">您当前选择的是多副本模式</h3>
    <h3 v-if="curPlan.StorageMode === 'EC'">您当前选择的是纠删码模式</h3>
    <el-form label-position="top">
      <el-form-item label="请选择存储云" size="large" v-if="curStep === 1">
        <el-transfer v-model="ReplicaClouds" :data="formattedClouds" :titles="['当前可用云', '已选择']">
          <!--                        <el-tooltip slot-scope="{ option }">-->
          <!--                          <div slot="content">-->
          <!--                            存储价格：{{ allClouds[option.index].StoragePrice }}元/GB/月<br />-->
          <!--                            流量价格：{{ allClouds[option.index].TrafficPrice }}元/GB<br />-->
          <!--                            可用性：{{ allClouds[option.index].Availability }}<br />-->
          <!--                          </div>-->
          <!--                          <span>{{ option.label }}</span>-->
          <!--                        </el-tooltip>-->
        </el-transfer>
      </el-form-item>
      <el-form-item label="请选择存储云" size="large" v-if="curStep === 2 && curPlan.StorageMode === 'EC'">
        <el-transfer v-model="ECNKClouds" :data="formattedECClouds" :titles="['作为数据分块的云', '作为校验分块的云']"> </el-transfer>
      </el-form-item>
    </el-form>

    <el-button v-if="curStep > 0" @click="prevStep">上一步</el-button>
    <el-button type="primary" @click="nextStep" :disabled="!curStepFinish">{{ curStep === finishStep ? "完成" : "下一步" }}</el-button>
  </el-card>
</template>

<script>
import Clouds from "@/api/clouds";
import Plan from "@/api/plan";

export default {
  name: "customizeStoragePlan",
  inject: ["formatPrice"],
  data() {
    return {
      curStep: 0,
      curPlan: {
        StorageMode: "Replica"
      },
      allClouds: [],
      formattedClouds: [],
      ReplicaClouds: [],
      ECNKClouds: []
    };
  },
  methods: {
    prevStep() {
      this.curStep -= 1;
    },
    nextStep() {
      if (this.curStep === this.finishStep) {
        if (this.curPlan.StorageMode === "Replica") {
          Plan.changeStoragePlan({
            ...this.curPlan,
            N: this.ReplicaClouds.length,
            K: 1,
            Clouds: this.ReplicaClouds.map(val => this.allClouds[val])
          })
            .then(resp => {
              if (resp) {
                this.$message.success("自定义存储方案成功！");
                this.curStep = 0;
                this.$store.dispatch("updateInfo", "StoragePlan");
                this.$emit("success");
              }
            })
            .catch(() => {
              this.$emit("failed");
            });
        } else if (this.curPlan.StorageMode === "EC") {
          Plan.changeStoragePlan({
            ...this.curPlan,
            N: this.ReplicaClouds.length,
            K: this.ReplicaClouds.length - this.ECNKClouds.length,
            Clouds: this.ReplicaClouds.filter(value => !this.ECNKClouds.includes(value))
              .map(val => this.allClouds[val])
              .concat(this.ECNKClouds.map(value => this.allClouds[value]))
          })
            .then(resp => {
              if (resp) {
                this.$message.success("自定义存储方案成功！");
                this.curStep = 0;
                this.$store.dispatch("updateInfo", "StoragePlan");
                this.$emit("success");
              }
            })
            .catch(() => {
              this.$emit("failed");
            });
        } else {
          this.$message.error("存储方式不存在！");
          this.$emit("error");
        }
      } else {
        this.curStep += 1;
        switch (this.curStep) {
          case 1:
            this.ReplicaClouds = [];
            break;
          case 2:
            this.ECNKClouds = [];
            break;
          default:
            break;
        }
      }
    },
    async getAllClouds() {
      Clouds.getAllClouds()
        .then(resp => {
          if (resp && resp.Clouds) {
            this.allClouds = resp.Clouds;
          }
        })
        .catch(() => {})
        .then(() => {
          this.allClouds = this.allClouds || [];
          this.formattedClouds = this.allClouds.map((value, index) => {
            return {
              ...value,
              key: index,
              label: value.CloudName
            };
          });
        });
    },
    renderCloud(h, option) {
      return (
        <el-tooltip slot-scope="{ option }">
          <div slot="content">
            存储价格：{this.formatPrice(this.allClouds[option.index].StoragePrice)}元/GB/月
            <br />
            流量价格：{this.formatPrice(this.allClouds[option.index].TrafficPrice)}元/GB
            <br />
            可用性：{this.allClouds[option.index].Availability}
            <br />
          </div>
          <span>{option.label}</span>
        </el-tooltip>
      );
    }
  },
  computed: {
    finishStep() {
      if (this.curPlan.StorageMode === "Replica") {
        return 2;
      }
      if (this.curPlan.StorageMode === "EC") {
        return 3;
      }
      return -1;
    },
    ReplicaCloudInPlan() {
      return this.ReplicaClouds.map(val => {
        return { ...this.allClouds[val] };
      });
    },
    formattedECClouds() {
      return this.ReplicaClouds.map(value => {
        return {
          key: value,
          label: this.allClouds[value].CloudName
        };
      });
    },
    curStepFinish() {
      switch (this.curStep) {
        case 0:
          return !!this.curPlan.StorageMode;
        case 1:
          if (this.curPlan.StorageMode === "Replica") {
            return this.ReplicaClouds.length > 0;
          }
          if (this.curPlan.StorageMode === "EC") {
            return this.ReplicaClouds.length >= 2;
          }
          return false;
        case 2:
          if (this.curPlan.StorageMode === "Replica") {
            return true;
          }
          if (this.curPlan.StorageMode === "EC") {
            return this.ECNKClouds.length > 0 && this.ECNKClouds.length <= this.ReplicaClouds.length / 2;
          }
          return false;
        default:
          return this.curStep === this.finishStep;
      }
    }
  },
  beforeMount() {
    this.getAllClouds();
  },
  watch: {
    curStep(newVal) {
      if (newVal === 0) {
        this.ReplicaClouds = [];
        this.ECNKClouds = [];
      }
    }
  }
};
</script>

<style scoped></style>
