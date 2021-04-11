<template>
  <div>
    <div class="filter-field" v-if="filterMap">
      <div v-for="(item, key) in filterMap" :key="key" class="filter-item">
        {{ item.label }}：
        <RenderCell v-if="item.render" :render="item.render" :itemKey="key" :filterData="filterData"> </RenderCell>
        <template v-else>
          <el-select v-if="item.type === 'select'" v-model="filterData[key]">
            <el-option v-for="(it, idx) in item.dataSource" :key="key + idx" :label="it.label" :value="it.value"></el-option>
          </el-select>
          <el-datepicker v-else-if="item.type === 'date'" v-model="filterData[key]" type="date"> </el-datepicker>
          <template v-else-if="item.type === 'datetimeRange'">
            <el-datepicker v-model="filterData[item.key1]" type="date" :clearable="false" :format="item.format"> </el-datepicker>
            -
            <el-datepicker v-model="filterData[item.key2]" type="date" :clearable="false" :format="item.format"> </el-datepicker>
          </template>
          <el-input v-else v-model="filterData[key]" @enter="onSearchClick"> </el-input>
        </template>
      </div>
    </div>
    <slot name="filterBtns"></slot>
    <el-table
      ref="multipleTable"
      :data="tableData"
      v-bind="$attrs"
      v-on="$listeners"
      :highlight-current-row="listColumns.length > 0 && listColumns[0].type !== 'selection'"
      tooltip-effect="dark"
      style="width: 100%"
    >
      <el-table-column
        v-for="(item, index) in listColumns"
        :key="index"
        :prop="item.prop"
        :label="item.label"
        :width="item.width"
        :formatter="item.formatter"
        :sortable="item.sortable"
        :show-overflow-tooltip="item.showOverflowTooltip"
        :type="item.type || ''"
      >
        <!-- <template slot="header" slot-scope="scope">
          <RenderCell
            v-if="typeof item.label === 'function'"
            :render="item.label"
            :index="scope.$index"
            :row="scope.row" />
          <span v-else>
            {{item.label}}
          </span>
        </template>
        <template slot-scope="scope" >
          <RenderCell
            v-if="item.render"
            :render="item.render"
            :index="scope.$index"
            :row="scope.row" />
          <span v-else>
            {{
              item.formatter? item.formatter(scope.row, scope.row[item.prop]): scope.row[item.prop]
            }}
          </span>
        </template> -->
      </el-table-column>
    </el-table>

    <el-pagenation
      v-if="pagenation"
      background
      :hide-on-single-page="true"
      @current-change="currentChange"
      :current-page="page"
      layout="prev, pager, next"
      :total="total"
    >
    </el-pagenation>
  </div>
</template>

<script>
import RenderCell from "./renderCell";

export default {
  name: "List",
  props: {
    getListAction: String, // 数据返回api
    map: Object, // 数据
    columns: Array, // 表头
    pagenation: {
      // 分页
      type: Boolean,
      default: false
    },
    query: Object, // api query
    filterMap: {
      // 列表筛选
      type: Object,
      default: () => {}
    },
    defaultFilterData: {
      type: Object,
      default: () => {}
    }
  },
  components: { RenderCell },
  data() {
    const filterData = { ...this.defaultFilterData };
    return {
      page: 1,
      size: 10,
      total: 0,
      tableData: [],
      listColumns: [],
      filterData
    };
  },
  mounted() {
    this.getList();
    this.listColumns = this.columns || this.mapToColums(this.map);
  },
  methods: {
    mapToColums(map) {
      const column = [];
      Object.keys(map).forEach(e => {
        if (e === "selection" || e === "index") {
          column.push({ type: e, key: map[e] });
        } else {
          column.push({ prop: e, label: map[e] });
        }
      });
      console.log(column);
      return column;
    },
    currentChange(page) {
      this.page = page;
      this.getList();
    },
    onSearchClick() {
      this.getList(); // todo
    },
    getList() {
      const params = {
        page: this.page,
        size: this.size,
        ...this.query
      };
      this.$Api[this.getListAction.split(".")[0]][this.getListAction.split(".")[1]](params).then(res => {
        if (res) {
          this.tableData = res.rows;
          this.total = res.total;
        }
      });
    }
  }
};
</script>

<style lang="scss" scoped>
.filter-field {
  .filter-item {
    margin-right: 20px;
    margin-bottom: 20px;
    .el-input {
      width: 100px;
    }
    float: left;
  }
}
</style>
