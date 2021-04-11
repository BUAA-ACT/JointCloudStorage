<template>
  <div>
    <el-card class="chartCard">
      <h4>
        存储卷信息
      </h4>
      <span class="tips" style="margin-bottom:10px">存储卷供用户创建的工作负载使用，是将工作负载数据持久化的一种资源对象。</span>
      <el-collapse v-model="activeName" accordion>
        <el-collapse-item title="什么是存储卷类型?" name="1">
          <div>存储卷类型 (StorageClass) 是由集群管理员配置存储服务端的参数，并按类型提供存储给集群用户使用。</div>
        </el-collapse-item>
        <el-collapse-item title="什么是本地存储卷 (Local Volume) ?" name="2">
          <div>本地存储卷表示挂载的本地存储设备，如磁盘、分区或目录。</div>
        </el-collapse-item>
      </el-collapse>
    </el-card>
    <el-card class="chartCard" style="height: 700px; padding:20px 0">
      <el-table :data="nodeStatus">
        <el-table-column prop="name" label="名称" width="300"> </el-table-column>
        <el-table-column prop="status" label="状态" width="140"> </el-table-column>
        <el-table-column prop="mode" label="访问模式" width="200"> </el-table-column>
        <el-table-column prop="mount" label="挂载" width="140"> </el-table-column>
        <el-table-column prop="createtime" label="创建时间" width="140"> </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<script>
export default {
  name: "StorageForm",
  data() {
    return {
      activeName: "1",
      nodeStatus: [
        {
          name: "minio-87dz96",
          status: "准备就绪",
          mode: "ReadWriteOnce",
          mount: "已挂载",
          createtime: "2021-01-07 13:49"
        },
        {
          name: "data-postgre-ux292w-postgresql-0",
          status: "准备就绪",
          mode: "ReadWriteOnce",
          mount: "已挂载",
          createtime: "2021-01-05 12:09"
        },
        {
          name: "db-data-mongodb-5gupv8-0",
          status: "准备就绪",
          mode: "ReadWriteOnce",
          mount: "已挂载",
          createtime: "2021-01-04 09:59"
        },
        {
          name: "ks-jenkins",
          status: "准备就绪",
          mode: "ReadWriteOnce",
          mount: "未挂载",
          createtime: "2020-12-09 20:23"
        },
        {
          name: "harbor-hrypxk-harbor-registry",
          status: "准备就绪",
          mode: "ReadWriteOnce",
          mount: "已挂载",
          createtime: "2020-11-22 03:00"
        },
        {
          name: "harbor-hrypxk-harbor-jobservice",
          status: "准备就绪",
          mode: "ReadWriteOnce",
          mount: "已挂载",
          createtime: "2020-11-22 03:00"
        },
        {
          name: "harbor-hrypxk-harbor-chartmuseum",
          status: "准备就绪",
          mode: "ReadWriteOnce",
          mount: "已挂载",
          createtime: "2020-11-22 03:00"
        }
      ]
    };
  }
};
</script>

<style lang="scss">
.chartCard {
  text-align: left;
}
</style>
