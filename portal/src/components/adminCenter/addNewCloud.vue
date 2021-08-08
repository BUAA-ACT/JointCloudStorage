<template>
  <div>
    <el-row>
      <el-col>
        <el-card shadow="always" class="requirementForm kuberx">
          <el-form ref="form" :model="form" label-width="100px" label-position="left">
            <el-form-item label="存储服务商">
              <el-radio-group v-model="form.ProviderName">
                <el-radio-button label="阿里"></el-radio-button>
                <el-radio-button label="华为"></el-radio-button>
                <el-radio-button label="腾讯"></el-radio-button>
                <el-radio-button label="百度"></el-radio-button>
                <el-radio-button label="金山"></el-radio-button>
              </el-radio-group>
            </el-form-item>
            <el-form-item label="云存储名称">
              <el-input v-model="form.CloudName" class="input"></el-input>
            </el-form-item>
            <el-form-item label="云际 ID">
              <el-input v-model="form.CloudID" class="input"></el-input>
            </el-form-item>
            <el-divider></el-divider>
            <el-form-item label="存储价格">
              <el-input-number v-model="form.StoragePrice" class="input" controls-position="right" :precision="precision" :step="step" />
              元/GB
            </el-form-item>
            <el-form-item label="流量价格">
              <el-input-number v-model="form.TrafficPrice" class="input" controls-position="right" :precision="precision" :step="step" />
              元/GB
            </el-form-item>
            <el-form-item label="可用性">
              <el-radio-group v-model="form.Availability">
                <el-radio-button v-for="val in availability" :label="val" :key="String(val)" />
              </el-radio-group>
            </el-form-item>
            <el-divider></el-divider>
            <el-form-item label="存储接入点">
              https://
              <el-input v-model="form.Endpoint" class="input"></el-input>
            </el-form-item>
            <el-form-item label="accessKey">
              <el-input v-model="form.AccessKey" class="input"></el-input>
            </el-form-item>
            <el-form-item label="secretKey">
              <el-input v-model="form.SecretKey" class="input"></el-input>
            </el-form-item>
            <el-form-item label="bucket 名称">
              <el-input v-model="form.Bucket" class="input"></el-input>
            </el-form-item>
            <el-divider></el-divider>
            <el-form-item label="云际地址">
              https://
              <el-input v-model="form.Address" class="input"></el-input>
            </el-form-item>
            <el-form-item label="地理位置">
              <select-point @getPoint="getPoint"></select-point>
            </el-form-item>
            <el-divider></el-divider>
          </el-form>
          <el-button type="primary" class="submitBtn" @click="onSubmit">确认提交</el-button>
        </el-card>
      </el-col>
    </el-row>
    <el-dialog title="请确认云信息" :visible.sync="confirmVisible" top="20px" append-to-body>
      <add-new-cloud-confirm ref="confirmDiag" :cloud="form" :modify="modify" @cancel="cancelConfirm" @success="successConfirm" />
    </el-dialog>
  </div>
</template>

<script>
import selectPoint from "@/components/adminCenter/selectPoint.vue";
import addNewCloudConfirm from "@/components/adminCenter/addNewCloudConfirm.vue";

export default {
  components: {
    selectPoint,
    addNewCloudConfirm
  },
  props: {
    cloud: {
      type: Object,
      required: false
    },
    modify: {
      type: Boolean,
      default: false
    }
  },
  data() {
    return {
      form: {
        CloudID: "ksyun-beijing",
        Endpoint: "ks3-cn-beijing.ksyun.com",
        AccessKey: "",
        SecretKey: "",
        StoragePrice: 0.12,
        TrafficPrice: 0.4,
        Availability: 0.9995,
        Status: "UP",
        Location: "116.411087,39.912695",
        Address: "ksyun-beijing.jointcloudstorage.cn",
        CloudName: "金山云-北京",
        ProviderName: "金山",
        Bucket: "jcspan-beijing"
      },
      confirmVisible: false,
      precision: 2,
      step: 0.1,
      availability: [0.99995, 0.9995, 0.995, 0.99, 0.95]
    };
  },
  methods: {
    onSubmit() {
      this.confirmVisible = true;
    },
    getPoint(point) {
      this.form.Location = `${point.lng},${point.lat}`;
      this.$log(this.form.Location);
    },
    successConfirm() {
      this.closeConfirm();
      this.$emit("success");
    },
    cancelConfirm() {
      this.closeConfirm();
      this.$emit("cancel-confirm");
    },
    closeConfirm() {
      this.confirmVisible = false;
    }
  },
  beforeMount() {
    if (this.cloud) {
      this.form = this.cloud;
    }
  },
  watch: {
    cloud(newVal) {
      if (newVal) {
        this.form = newVal;
      }
    }
  }
};
</script>

<style lang="scss" scoped>
.kuberx {
  text-align: left;
  padding-top: 20px;

  .el-icon-question {
    color: #dddddd;
    margin: 0 5px;
  }

  table {
    border-spacing: 0;

    td {
      white-space: nowrap;
    }
  }
}

.input {
  width: 200px;
}

.formTable {
  white-space: nowrap;
}
</style>
