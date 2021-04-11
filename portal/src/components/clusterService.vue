<template>
  <div class="requirement">
    <div class="sBanner sBanner-kubeX"></div>
    <div class="bg-gray">
      <div class="manufacturer disasterRecovery">
        <table style="width:100%">
          <tr>
            <td style="width:80px;padding:0;vertical-align: super;">
              <p style="margin:5px">KUBEX购买</p>
            </td>
            <td>
              <el-steps :active="activeNum" style="margin-bottom:30px">
                <el-step title="服务配置选型"></el-step>
                <el-step title="配置确认"></el-step>
                <el-step title="部署完成"></el-step>
              </el-steps>
            </td>
          </tr>
        </table>

        <router-view></router-view>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  data() {
    return {
      activeNum: 1
    };
  }
};
</script>

<style lang="scss" scoped></style>
