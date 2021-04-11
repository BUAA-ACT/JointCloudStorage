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
                  <td>配置名称</td>
                  <td>CPU</td>
                  <td>内存</td>
                  <td>网络资源</td>
                </tr>
                <tr v-for="(li, ind) in formData[k]" :key="ind">
                  <td>{{ li.key1 }}</td>
                  <td>{{ li.key2 }}</td>
                  <td>{{ li.key3 }}</td>
                  <td>{{ li.key4 }}</td>
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
      <el-button @click="$router.go(-1)">重新选择</el-button>
      <el-button type="primary" @click="buy">立即购买</el-button>
    </el-card>
    <el-dialog title="购买成功" :visible.sync="dialogVisible" width="500px">
      <i class="el-icon-success" style="display:block;text-align: center;color:green;font-size:3rem;margin: 20px auto 50px auto"></i>
      <p>尊敬的admin用户，您的【kubeX】集群管理服务已购买成功。</p>
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
          key2: "地域及可用区"
        },
        {
          key3: "集群名称",
          key4: "集群服务商",
          key5: "集群规模",
          key6: "控制节点数",
          key7: "Master配置选择",
          key8: "操作系统",
          key9: "集群网络",
          key10: "服务网段",
          key11: "服务网络"
        },
        {
          key12: "节点名称",
          key13: "Worker配置选择",
          key14: "所在子网",
          key15: "容器网络",
          key16: "Worker节点数"
        }
      ],
      formData: {
        key1: "带宽计费",
        key2: "北方一区-北京",
        key3: "vpc-5zq94h1j",
        key4: "华为",
        key5: "50节点",
        key6: "1",
        key7: [
          {
            key1: "s1.large.1",
            key2: "2核",
            key3: "4G",
            key4: "20Mbps"
          }
        ],
        key8: "Ubantu Server 18.04",
        key9: "外网云集群",
        key10: "使用默认网段",
        key11: "172.168.0.0",
        key12: "vpc-5zq94h1j",
        key13: [
          {
            key1: "s1.large.1",
            key2: "2核",
            key3: "4G",
            key4: "20Mbps"
          }
        ],
        key14: "外网云集群",
        key15: "172.168.0.0/14，172.168.0.0/15，172.168.0.0/16",
        key16: "3"
      },
      dialogVisible: false
    };
  },
  methods: {
    buy() {
      this.dialogVisible = true;
      setTimeout(() => {
        this.dialogVisible = false;
        this.$router.push("deploymentComplete");
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
