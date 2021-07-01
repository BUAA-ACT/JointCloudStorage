<template>
  <div>
    <div v-if="newPlanAvailable" class="migration-operations">
      <el-button type="primary" @click="accept" :loading="submitLoading">迁移</el-button><br />
      <!--      <el-button type="primary" @click="cancel" :loading="cancelLoading" :disabled="!plansLoaded">我再想想</el-button><br />-->
      <el-button type="info" @click="cancel" :loading="submitLoading">取消</el-button><br />
    </div>
    <div v-if="newPlanAvailable">
      <div class="plans-viewer-container">
        <el-card class="box-card">
          <div slot="header" class="clearfix">
            旧存储方案
          </div>
          <div class="text item">
            存储模式： {{ Advices.StoragePlanOld.StorageMode }}<br />
            存储价格： {{ Advices.StoragePlanOld.StoragePrice }}<br />
            流量价格： {{ Advices.StoragePlanOld.TrafficPrice }}<br />
            可用性：{{ Advices.StoragePlanOld.Availability }}
          </div>
        </el-card>
        <el-card class="box-card">
          <div slot="header" class="clearfix">
            新存储方案
          </div>
          <div class="text item">
            存储模式： {{ Advices.StoragePlanNew.StorageMode }}<br />
            存储价格： {{ Advices.StoragePlanNew.StoragePrice }}<br />
            流量价格： {{ Advices.StoragePlanNew.TrafficPrice }}<br />
            可用性：{{ Advices.StoragePlanNew.Availability }}
          </div>
        </el-card>
      </div>
      <location-viewer :clouds="Advices.CloudsOld" :new-clouds="Advices.CloudsNew" :dynamic="migrating" class="location-viewer" ref="viewer" />
    </div>
    <div class="no-new-plan" v-else>
      <i class="el-icon-success tip-icon"></i><br />
      暂时没有为您找到更优的存储方案！
    </div>
  </div>
</template>

<script>
import Plan from "@/api/plan";
import locationViewer from "@/components/viewer/locationViewer.vue";

export default {
  name: "planAdvice",
  components: {
    locationViewer
  },
  data() {
    return {
      Advices: {},
      submitLoading: false,
      newPlanAvailable: false,
      migrating: false
    };
  },
  methods: {
    async getNewAdvice() {
      this.newPlanAvailable = false;
      Plan.getNewAdvice().then(resp => {
        if (resp) {
          if (resp.Advices && resp.Advices.length > 0) {
            [this.Advices] = resp.Advices;
            this.newPlanAvailable = true;
          } else {
            this.Advices = {};
          }
        }
      });
    },
    async accept() {
      this.submitLoading = true;
      this.migrating = true;
      await Plan.acceptAdvice();
      // TODO: 轮询
      this.submitLoading = false;
    },
    async cancel() {
      this.submitLoading = true;
      await Plan.abandonAdvice();
      this.submitLoading = false;
    }
  },
  mounted() {
    this.getNewAdvice();
  }
};
</script>

<style scoped lang="scss">
.no-new-plan {
  .tip-icon {
    font-size: 100px;
    color: #67c23a;
  }
}
.migration-operations {
  display: flex;
  width: 200px;
  justify-content: space-between;
}
.plans-viewer-container {
  display: -webkit-box;
  display: -webkit-flex;
  display: -ms-flexbox;
  display: flex;
  margin-top: 30px;
  -webkit-box-pack: justify;
  -webkit-justify-content: space-between;
  -ms-flex-pack: justify;
  justify-content: space-between;
  width: 800px;
}
.location-viewer {
  width: 800px;
  height: 400px;
}
</style>
