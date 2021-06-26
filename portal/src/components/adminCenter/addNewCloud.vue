<template>
  <div>
    <el-row>
      <el-col>
        <el-card shadow="always" class="requirementForm kuberx">
          <el-form ref="form" :model="form" label-width="100px" label-position="left">
            <el-form-item label="存储服务商">
              <el-radio-group v-model="form.resource">
                <el-radio-button label="阿里"></el-radio-button>
                <el-radio-button label="华为"></el-radio-button>
                <el-radio-button label="腾讯"></el-radio-button>
                <el-radio-button label="百度"></el-radio-button>
                <el-radio-button label="金山"></el-radio-button>
              </el-radio-group>
            </el-form-item>
            <el-form-item label="云存储名称">
              <el-input v-model="form.cloudName" class="input"></el-input>
            </el-form-item>
            <el-form-item label="云际 id">
              <el-input v-model="form.cloudId" class="input"></el-input>
            </el-form-item>
            <el-divider></el-divider>
            <el-form-item label="存储价格">
              <el-input v-model="form.storagePrice" class="input"></el-input>
              元/GB
            </el-form-item>
            <el-form-item label="流量价格">
              <el-input v-model="form.trafficPrice" class="input"></el-input>
              元/GB
            </el-form-item>
            <el-form-item label="可用性">
              <el-radio-group v-model="form.availability">
                <el-radio-button label="0.99995"></el-radio-button>
                <el-radio-button label="0.9995"></el-radio-button>
                <el-radio-button label="0.995"></el-radio-button>
              </el-radio-group>
            </el-form-item>
            <el-divider></el-divider>
            <el-form-item label="存储接入点">
              https://
              <el-input v-model="form.endpoint" class="input"></el-input>
            </el-form-item>
            <el-form-item label="accessKey">
              <el-input v-model="form.accessKey" class="input"></el-input>
            </el-form-item>
            <el-form-item label="secretKey">
              <el-input v-model="form.secretKey" class="input"></el-input>
            </el-form-item>
            <el-form-item label="bucket 名称">
              <el-input v-model="form.bucket" class="input"></el-input>
            </el-form-item>
            <el-divider></el-divider>
            <el-form-item label="云际地址">
              https://
              <el-input v-model="form.address" class="input"></el-input>
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
  </div>
</template>

<script>
import selectPoint from "@/components/adminCenter/selectPoint.vue";

export default {
  components: {
    selectPoint
  },
  data() {
    return {
      form: {
        storagePrice: "0.01",
        cloudName: "",
        cloudId: "",
        location: "",
        resource: "阿里",
        trafficPrice: "0.01",
        availability: "0.9995",
        endPoint: "",
        accessKey: "",
        secretKey: "",
        bucket: "",
        address: ""
      },
      cpuList: [],
      storageList: [],
      dataDisk: [],
      tableData: [
        {
          name: "S1.large.1",
          cpu: "2核",
          storage: "4G",
          net: "20Mbps"
        },
        {
          name: "S1.large.1",
          cpu: "2核",
          storage: "4G",
          net: "20Mbps"
        },
        {
          name: "S1.large.1",
          cpu: "2核",
          storage: "4G",
          net: "20Mbps"
        }
      ],
      netDialogVisible: false
    };
  },
  methods: {
    onSubmit() {
      this.$router.push({ name: "addNewCloudConfirm", params: { formData: this.form } });
    },
    getPoint(point) {
      this.form.location = `${point.lng},${point.lat}`;
      this.$log(this.form.location);
    }
  },
  mounted() {
    if (this.$route.params.formData) {
      this.$log("ok");
      this.form = this.$route.params.formData;
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

.el-input__inner {
  width: 50px;
}

.formTable {
  white-space: nowrap;
}
</style>
