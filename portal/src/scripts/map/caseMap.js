// CPS Area VCPU Memory System Disk Type Data Disk Type Bandwidth Number Price
const caseMap = [
  { prop: "instance", label: "厂商" },
  { prop: "vendor", label: "可用区" },
  { prop: "vendor", label: "VCPU" },
  { prop: "vendor", label: "内存" },
  { prop: "vendor", label: "系统盘" },
  { prop: "vendor", label: "类型" },
  { prop: "benchTime", label: "数据盘", formatter: (row, column, cellValue) => `${cellValue || 0}GB` },
  { prop: "benchTime", label: "公网带宽", formatter: (row, column, cellValue) => `${cellValue || 0}Mbps` },
  { prop: "vendor", label: "数量" },
  { prop: "vendor", label: "主机单价" }
];

const selectedCaseMap = [{ type: "selection", label: "选择" }].concat(caseMap);

const fileViewMap = [
  { type: "index", label: "序号" },
  { prop: "vender", label: "文件名称", formatter: () => "/t1efne231" },
  { prop: "vendor", label: "修改时间", formatter: () => "2020/10/16 11:22:33" },
  { prop: "vendor", label: "编码类型", formatter: () => "edu_3_2" },
  { prop: "vendor", label: "条带数", formatter: () => "4" }
];
export { caseMap, selectedCaseMap, fileViewMap };
