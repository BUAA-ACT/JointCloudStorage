<template>
  <div>
    <el-card shadow="always" class="requirementForm kuberx">
      <el-form ref="form" :model="form" label-width="100px" label-position="left">
        <el-form-item label="计费模式">
          <el-radio-group v-model="form.resource">
            <el-radio-button label="带宽计费"></el-radio-button>
            <el-radio-button label="按量计费"></el-radio-button>
          </el-radio-group>
        </el-form-item>
        <el-divider></el-divider>
        <el-form-item label="地域及可用区">
          <table>
            <tr>
              <td>
                <el-select>
                  <el-option label="华北2（北京）" value="1"></el-option>
                </el-select>
              </td>
              <td>
                <el-radio-group v-model="form.resource">
                  <el-radio-button label="随机分配"></el-radio-button>
                  <el-radio-button label="可用区H"></el-radio-button>
                  <el-radio-button label="可用区G"></el-radio-button>
                  <el-radio-button label="可用区F"></el-radio-button>
                  <el-radio-button label="可用区C"></el-radio-button>
                  <el-radio-button label="可用区E"></el-radio-button>
                  <el-radio-button label="可用区D"></el-radio-button>
                </el-radio-group>
              </td>
            </tr>
          </table>

          <span class="tips">不同地域的实例之间内网互不相通；选择靠近您客户的地域，可降低网络时延，提高您客户的访问速度。</span>
        </el-form-item>
        <el-divider></el-divider>
        <el-form-item label="集群名称">
          <el-input v-model="form.resource"></el-input>
        </el-form-item>
        <el-form-item label="集群服务商">
          <el-radio-group v-model="form.resource">
            <el-radio-button label="华为"></el-radio-button>
            <el-radio-button label="腾讯"></el-radio-button>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="集群规模">
          <el-radio-group v-model="form.resource">
            <el-radio-button label="50节点"></el-radio-button>
            <el-radio-button label="200节点"></el-radio-button>
            <el-radio-button label="1000节点"></el-radio-button>
            <el-radio-button label="2000节点"></el-radio-button>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="控制节点数">
          <el-radio-group v-model="form.resource">
            <el-radio-button label="3"></el-radio-button>
            <el-radio-button label="1"></el-radio-button>
          </el-radio-group>
        </el-form-item>

        <el-form-item label="Master配置选择">
          <el-table ref="multipleTable" :data="tableData" tooltip-effect="dark" style="width: 100%" @selection-change="handleSelectionChange">
            <el-table-column type="selection" width="55"> </el-table-column>
            <el-table-column label="配置名称" prop="name"> </el-table-column>
            <el-table-column prop="cpu" label="CPU"> </el-table-column>
            <el-table-column prop="storage" label="内存"> </el-table-column>
            <el-table-column prop="net" label="网络资源"> </el-table-column>
          </el-table>
        </el-form-item>

        <el-form-item label="操作系统">
          <el-radio-group v-model="form.resource">
            <el-radio-button label="Ubantu Server 18.04"></el-radio-button>
            <el-radio-button label="CentOs 7.6 bit"></el-radio-button>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="集群网络">
          <table>
            <tr>
              <td><el-select></el-select></td>
              <td><el-button @click="netDialogVisible = true">新建网络</el-button></td>
            </tr>
          </table>
        </el-form-item>
        <el-form-item label="服务网段">
          <el-radio-group v-model="form.resource">
            <el-radio-button label="使用默认网段"></el-radio-button>
            <el-radio-button label="手动设置网段"></el-radio-button>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="服务网络">
          <table>
            <tr>
              <td>
                <el-input></el-input>
              </td>
              <td>
                <el-input></el-input>
              </td>
              <td>
                .0.
              </td>
              <td>
                <el-input></el-input>
              </td>
            </tr>
            <tr>
              <td>
                <el-radio label="自动选择"></el-radio>
              </td>
            </tr>
          </table>
        </el-form-item>

        <el-divider></el-divider>
        <el-form-item label="节点名称">
          <el-input></el-input>
        </el-form-item>
        <el-form-item label="Worker配置选择">
          <el-table ref="multipleTable" :data="tableData" tooltip-effect="dark" style="width: 100%" @selection-change="handleSelectionChange">
            <el-table-column type="selection" width="55"> </el-table-column>
            <el-table-column label="配置名称" prop="name"> </el-table-column>
            <el-table-column prop="cpu" label="CPU"> </el-table-column>
            <el-table-column prop="storage" label="内存"> </el-table-column>
            <el-table-column prop="net" label="网络资源"> </el-table-column>
          </el-table>
        </el-form-item>
        <el-form-item label="Worker配置选择"> <el-input-number :min="6" :max="128" v-model="form.num"></el-input-number>台 </el-form-item>
        <el-form-item label="登录">
          <table>
            <tr>
              <td>用户名</td>
              <td><el-input></el-input></td>
              <td>密码</td>
              <td><el-input></el-input></td>
              <td>确认密码</td>
              <td><el-input></el-input></td>
            </tr>
          </table>
        </el-form-item>
        <el-form-item label="所在子网">
          <el-select></el-select>
        </el-form-item>
        <el-form-item label="容器网络">
          <table>
            <tr>
              <td>
                <el-input></el-input>
              </td>
              <td>
                <el-input></el-input>
              </td>
              <td>
                .0.0./
              </td>
              <td>
                <el-input></el-input>
              </td>
            </tr>
            <tr>
              <td>
                <el-radio label="自动选择"></el-radio>
              </td>
            </tr>
          </table>
        </el-form-item>
      </el-form>
      <el-button type="primary" class="submitBtn" @click="onSubmit">确认提交</el-button>
    </el-card>
    <el-dialog title="新建VPC" :visible.sync="netDialogVisible" width="800px">
      <el-form ref="form" :model="form" label-width="100px" label-position="left">
        <p>私有网络信息</p>
        <el-form-item label="所属地域">
          华南地区（广州）
        </el-form-item>
        <el-form-item label="名称">
          <el-input placeholder="请输入名称"></el-input>
        </el-form-item>
        <el-form-item label="IPv4 CIDR">
          <table class="formTable">
            <tr>
              <td>
                <el-input></el-input>
              </td>
              <td>
                <el-input></el-input>
              </td>
              <td>
                .0.0./
              </td>
              <td>
                <el-input></el-input>
              </td>
              <td>
                <span class="orange">创建后不可修改</span>
              </td>
            </tr>
          </table>

          <span class="tips">为了您可以更好的使用私用网络服务，建议您提前做好网络规划。</span>
        </el-form-item>
        <p>初始子网信息</p>
        <el-form-item label="子网名称">
          <el-input placeholder="请输入名称"></el-input>
        </el-form-item>
        <el-form-item label="IPv4 CIDR">
          <table class="formTable">
            <tr>
              <td>
                10.0
              </td>
              <td>
                <el-input></el-input>
              </td>
              <td>
                .0/
              </td>
              <td>
                <el-input></el-input>
              </td>
            </tr>
          </table>

          <span class="tips">IP地址剩余253个</span>
        </el-form-item>
        <el-form-item label="可用区">
          <el-radio-group>
            <el-radio-button label="随机可用区"></el-radio-button>
            <el-radio-button label="南京一区"></el-radio-button>
            <el-radio-button label="南京二区"></el-radio-button>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="关联路由器">
          默认
        </el-form-item>
      </el-form>
      <span slot="footer" class="dialog-footer">
        <el-button type="primary" @click="netDialogVisible = false">确认提交</el-button>
      </span>
    </el-dialog>
  </div>
</template>

<script>
export default {
  data() {
    return {
      form: {},
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
      this.$router.push({ path: "clusterService/settingConfirm" });
    }
  }
};
</script>

<style lang="scss" scoped>
.kuberx {
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
.formTable {
  white-space: nowrap;
}
</style>
