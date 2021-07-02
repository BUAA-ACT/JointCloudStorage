<template>
  <div>
    <el-table :data="clouds" style="width: 100%">
      <el-table-column type="expand">
        <template slot-scope="props">
          <el-form label-position="left" inline class="cloud-detail-info-table-expand" label-width="90px">
            <el-form-item label="云际地址">
              <span>https://{{ props.row.Address }}</span>
            </el-form-item>
            <el-form-item label="Access Key">
              <span>{{ props.row.AccessKey }}</span>
            </el-form-item>
            <el-form-item label="Secret Key">
              <span>{{ props.row.SecretKey }}</span>
            </el-form-item>
            <el-form-item label="访问端点">
              <span>https://{{ props.row.Endpoint }}</span>
            </el-form-item>
            <el-form-item label="存储价格">
              <span>{{ props.row.StoragePrice }}</span>
            </el-form-item>
            <el-form-item label="流量价格">
              <span>{{ props.row.TrafficPrice }}</span>
            </el-form-item>
            <el-form-item label="地理位置">
              <span>{{ props.row.Location }}{{ cloudCity(props.row.Location) }}</span>
            </el-form-item>
            <el-form-item label="可用性">
              <span>{{ props.row.Availability * 100 }}%</span>
            </el-form-item>
          </el-form>
        </template>
      </el-table-column>
      <el-table-column prop="CloudID" label="云节点ID" width="180"> </el-table-column>
      <el-table-column prop="CloudName" label="云节点名称" width="180"> </el-table-column>
      <el-table-column prop="ProviderName" label="供应商" width="60"> </el-table-column>
      <el-table-column prop="Status" label="状态" width="80">
        <template slot-scope="status">
          <el-tag :type="status2Type[status.row.Status]">{{ status.row.Status }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="操作">
        <template v-slot="op">
          <slot :cloud="op.row" name="operations"></slot>
        </template>
      </el-table-column>
    </el-table>
  </div>
</template>

<script>
/**
 * Cloud List
 * @usage: {@link @/components/adminCenter/manageClouds.vue}
 * <cloud-list :clouds="[]" >
 *   <template #operations="prop">
 *     The Cloud is {{ prop.cloud }}.
 *   </template>
 * </cloud-list>
 * @see: manageClouds.vue
 */
import Location from "@/utils/location";

export default {
  name: "cloudList",
  props: {
    clouds: {
      type: Array,
      required: true
    }
  },
  data() {
    return {
      status2Type: {
        UP: "success",
        DOWN: "danger"
      },
      cloud2Cities: {}
    };
  },
  methods: {
    cloudCity(cloudLocation) {
      const ret = this.cloud2Cities[cloudLocation];
      if (ret) {
        return ret;
      }
      this.convertLocation(cloudLocation);
      return "";
    },
    /**
     * Convert location string to city name
     * @param {string} location
     */
    async convertLocation(location) {
      const ret = await Location.convertToCityName(location);
      this.$log(ret);
      return ret;
    }
  }
};
</script>

<style scoped lang="scss">
.cloud-detail-info-table-expand {
  font-size: 0;
  label {
    width: 90px;
    color: #99a9bf;
  }
  .el-form-item {
    margin-right: 0;
    margin-bottom: 0;
    width: 33%;
  }
}
</style>
