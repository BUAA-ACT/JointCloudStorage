<template>
  <div>
    <el-card shadow="always" class="requirementForm settingConfirm">
      <template v-for="(item, index) in formDatakey">
        <div :key="index">
          <div v-for="(i, k) in item" :key="k" class="el-form-item el-form-item--small">
            <label class="el-form-item__label">{{ i }}</label>
            <div class="el-form-item__content">
              {{ formData[k] instanceof Array ? "" : formData[k] }}
              <table v-if="formData[k] instanceof Array">
                <tr>
                  <td>CPU</td>
                  <td>内存</td>
                  <td>架构</td>
                </tr>
                <tr v-for="(li, ind) in formData[k]" :key="ind">
                  <td>{{ li.key1 }}</td>
                  <td>{{ li.key2 }}</td>
                  <td>{{ li.key3 }}</td>
                </tr>
              </table>
            </div>
          </div>
          <el-divider></el-divider>
        </div>
      </template>
      <div class="el-form-item el-form-item--small">
        <label class="el-form-item__label">总计费用</label>
        <div class="el-form-item__content">
          <p><span class="orange">0.11</span>元/小时（配置费用）</p>
          <p><span class="orange">5.80</span>元/GB（网络费用-按使用流量）</p>
        </div>
      </div>
      <el-divider></el-divider>
      <el-form ref="form" :model="form" label-width="100px" label-position="left">
        <el-form-item label="服务器数量"> <el-input-number :min="1" :max="128" v-model="form.num"></el-input-number>台 </el-form-item>
        <el-form-item label="购买时长"> <el-input-number :min="1" :max="5000" v-model="form.duration"></el-input-number>月 </el-form-item>
        <el-form-item label="总计费用">
          <span class="orange">{{ form.num * form.duration * 2000 }}</span
          >元/小时 (详细费用2000/台)
        </el-form-item>
      </el-form>
      <el-button @click="$router.go(-1)">重新选择</el-button>
      <el-button type="primary" @click="buy">立即购买</el-button>
    </el-card>
    <el-dialog title="购买成功" :visible.sync="dialogVisible" width="500px">
      <i class="el-icon-success" style="display:block;text-align: center;color:green;font-size:3rem;margin: 20px auto 50px auto"></i>
      <p>尊敬的admin用户，您的服务已购买成功。</p>
      <span slot="footer" class="dialog-footer"> </span>
    </el-dialog>
  </div>
</template>

<script>
export default {
  data() {
    return {
      formDatakey: [
        {
          key1: "计费模式",
          key2: "带宽值",
          key3: "付费模式",
          key4: "地域及可用区"
        },
        {
          key5: "硬件规格配置",
          key6: "存储系统盘",
          key7: "存储数据盘",
          key8: "公共镜像",
          key9: "版本",
          key10: "网络",
          key11: "网络IP"
        }
      ],
      formData: {
        key1: "按固定带宽",
        key2: "20Mbps",
        key3: "包年包月",
        key4: "北方一区-北京",
        key5: [
          {
            key1: "Intel酷睿i5 10400F 14核 2.9GHz",
            key2: "16G",
            key3: "X86"
          }
        ],
        key6: "ESSD云盘 40GIB 2880 IOPS 性能级别",
        key7: "ESSD云盘 512GIB 2880 IOPS 性能级别",
        key8: "Alibaba Cloud Linux",
        key9: "Alibaba Cloud Linux2",
        key10: "默认网络",
        key11: "分配公网IPV4地址"
      },
      form: {
        num: 1,
        duration: 1
      },
      dialogVisible: false
    };
  },
  methods: {
    buy() {
      this.dialogVisible = true;
      setTimeout(() => {
        this.dialogVisible = false;
        this.$router.push("servicePurchase/deploymentComplete");
      }, 2000);
    }
  }
};
</script>

<style lang="scss">
.settingConfirm {
  padding: 20px;
  table {
    width: 100%;
    text-align: center;
    border-spacing: 0;
    border: 1px solid #dddddd;
    tr {
      border-bottom: 1px solid #dddddd;
    }
    tr:first-child {
      background-color: #f3f3f3;
      border: 0;
    }
  }
  .el-form-item {
    .el-form-item__label {
      width: 100px;
      text-align: left;
    }
    .el-form-item__content {
      margin-left: 100px;
    }
  }
}
</style>
