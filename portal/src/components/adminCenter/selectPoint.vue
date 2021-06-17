<template>
  <el-form>
    <el-row :gutter="24">
      <el-col :span="12">
        <el-col :span="12">
          <el-form-item label="位置经度" prop="lng">
            <el-input v-model="model.lng" type="number" class="input_number" @mousewheel.native.prevent />
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item label="位置纬度" prop="lat">
            <el-input v-model="model.lat" type="number" class="input_number" @mousewheel.native.prevent />
          </el-form-item>
        </el-col>
        <el-col :span="24">
          <baidu-map class="bm-view" ak="QFgFQorrB84maOZh0pPGC8kUiP0mGIhx" :center="center" :zoom="zoom" :scroll-wheel-zoom="true" @ready="createMap">
            <!--            搜索-->
            <bm-local-search :keyword="model.address" :auto-viewport="true" style="display: none"></bm-local-search>
            <!--            标记-->
            <bm-marker :position="{ lng: model.lng, lat: model.lat }" />
          </baidu-map>
          <el-input
            v-model="model.address"
            placeholder="搜索地点"
            style="margin-left: 10px;width: 200px;position: absolute;top: 25%;opacity: 0.9"
            prefix-icon="el-icon-search"
          ></el-input>
        </el-col>
      </el-col>
    </el-row>
  </el-form>
</template>

<script>
import BaiduMap from "vue-baidu-map/components/map/Map.vue";
import BmLocalSearch from "vue-baidu-map/components/search/LocalSearch.vue";
import BmMarker from "vue-baidu-map/components/overlays/Marker.vue";

export default {
  components: {
    BaiduMap,
    BmLocalSearch,
    BmMarker,
  },
  data() {
    return {
      center: { lng: 0, lat: 0 },
      zoom: 10,
      model: {
        lng: "",
        lat: "",
        address: ""
      }
    };
  },
  methods: {
    createMap({ BMap, map }) {
      // 百度地图API功能
      this.center.lng = 116.413387;
      this.center.lat = 39.926861;
      const viewthis = this;
      map.addEventListener("click", function(e) {
        viewthis.model.lng = e.point.lng;
        viewthis.model.lat = e.point.lat;
        viewthis.$emit("getPoint", { lng: e.point.lng, lat: e.point.lat });
      });

      // 设置区域图
      function getBoundary(data, bdary) {
        // eslint-disable-next-line no-param-reassign
        data = data.split("-");
        bdary.get(data[0], function(rs) {
          // 获取行政区域
          const count = rs.boundaries.length; // 行政区域的点有多少个
          // var pointArray = []
          for (let i = 0; i < count; i+=1) {
            const ply = new BMap.Polygon(rs.boundaries[i], {
              strokeWeight: 2,
              strokeColor: "#ff0000",
              fillOpacity: 0.1,
              fillColor: data[1]
            }); // 建立多边形覆盖物
            map.addOverlay(ply); // 添加覆盖物
          }
        });
      }
      // 区域图
      const datas = new Array("徐州市-#665599");
      const bdary = new BMap.Boundary();
      for (let i = 0; i < datas.length; i+=1) {
        getBoundary(datas[i], bdary);
      }


    }
  }
};
</script>

<style scoped>
.bm-view {
  width: 100%;
  height: 300px;
}
</style>
