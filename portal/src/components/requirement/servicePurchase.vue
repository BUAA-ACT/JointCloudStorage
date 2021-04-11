<template>
  <div class="requirement">
    <div class="bg-gray">
      <div class="manufacturer disasterRecovery">
        <table style="width:100%">
          <tr>
            <td style="width:80px;padding:0;vertical-align: super;">
              <p style="margin:5px">服务购买</p>
            </td>
            <td>
              <el-steps :active="activeNum" style="margin-bottom:30px">
                <el-step title="确认购买"></el-step>
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
