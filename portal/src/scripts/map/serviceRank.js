const rankColumns = [
  { prop: "instance", label: "实例类型" },
  { prop: "vendor", label: "地址" },
  { prop: "benchTime", label: "启动时间", formatter: (row, column, cellValue) => `${cellValue || 0}s` },
  { prop: "totalcpu", label: "cpu性能", formatter: (row, column, cellValue) => Number(cellValue).toFixed(2) },
  { prop: "totalmem", label: "内存性能", formatter: (row, column, cellValue) => Number(cellValue).toFixed(2) },
  { prop: "totalio", label: "磁盘性能", formatter: (row, column, cellValue) => Number(cellValue).toFixed(2) },
  { prop: "totalnet", label: "网络性能", formatter: (row, column, cellValue) => Number(cellValue).toFixed(2) },
  { prop: "total", label: "总分数", formatter: (row, column, cellValue) => Number(cellValue).toFixed(2) }
];

const serviceRankColumns = [{ type: "index", label: "排名" }].concat(rankColumns);
const selectedRankColumns = [{ type: "selection", label: "排名" }].concat(rankColumns);
export { serviceRankColumns, selectedRankColumns };
