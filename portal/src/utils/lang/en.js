const base = {
  lang: "中"
};
const disasterReco = {
  title: "Cloud Disaster Tolerance",
  intro:
    "Recently, cloud data center (CDC) failures are commonplace. To prevent users' data from being inaccessible due to CDC failures, we design a JointCloud disaster tolerance system (JCDT). JCDT encodes data to several coded blocks with erasure code and distributes these coded blocks to cloud hosts of multiple CDCs. When some of these CDCs fail, users' original data can be retrieved by coded blocks stored in available CDCs. JCDT designs a data writing method, a data read method, and a data repair method for JointCloud erasure codes to achieve high disaster tolerance, low redundancy, and high efficiency.",
  bannerText: "Disaster Tolerance service with high performance, low cost, verifiable and easy operation and maintenance",
  tab1: "Erasure Code Selection",
  tab2: "Settings",
  tab3: "Setting Confirm",
  tab4: "File Viewer",
  select: "Select",
  type: "Type of erasure code",
  explanation: "Explanation",
  next: "Next",
  prev: "Prev",
  option1:
    "The JCDT system with this erasure code requires 13 cloud hosts (12 storage hosts and 1 manager host) distributed in 3 areas. It can tolerate the failure of arbitrary 3 storage hosts or 1 area. The redundancy is 1.  The repair efficiency is relatively high.",
  option2:
    "The JCDT system with this erasure code requires 9 cloud hosts (8 storage hosts and 1 manager host) distributed in 2 areas. It can tolerate the failure of arbitrary 3 storage hosts or 1 area. The redundancy is 5/3.  The repair efficiency is relatively high.",
  option3:
    "The JCDT system with this erasure code requires 7 cloud hosts (6 storage hosts and 1 manager host) distributed in 2 areas. It can tolerate the failure of arbitrary 3 storage hosts or 1 area. The redundancy is 1. ",
  option4:
    "The JCDT system with this erasure code requires 7 cloud hosts (6 storage hosts and 1 manager host) distributed in 3 areas. It can tolerate the failure of arbitrary 2 storage hosts or 1 area. The redundancy is 1/2.",
  option5:
    "The JCDT system with this erasure code requires 9 cloud hosts (8 storage hosts and 1 manager host) distributed in 4 areas. It can tolerate the failure of arbitrary 2 storage hosts or 1 area. The redundancy is 1/3.",
  option6:
    "The JCDT system with this erasure code requires 10 cloud hosts (9 storage hosts and 1 manager host) distributed in 3 areas. It can tolerate the failure of arbitrary 3 storage hosts or 1 area. The redundancy is 1/2.",
  option7:
    "The JCDT system with this erasure code requires 13 cloud hosts (12 storage hosts and 1 manager host) distributed in 2 areas. It can tolerate the failure of arbitrary 6 storage hosts or 1 area. The redundancy is 1.",
  option8:
    "The JCDT system with this erasure code requires 13 cloud hosts (12 storage hosts and 1 manager host) distributed in 3 areas. It can tolerate the failure of arbitrary 4 storage hosts or 1 area. The redundancy is 1/2.",
  formLabel1: "Number of cloud hosts",
  inputTip1: "Explanation: the number of cloud hosts cannot be smaller than the number of the coded blocks in each stripe of the erasure code.",
  formLabel2: "Number of Areas",
  inputTip2: "Explanation: the number of areas cannot be smaller than 1/disaster-tolerance.",
  formLabel3: "Number of CPUs",
  inputTip3:
    "Explanation: JCDT system needs to run HDFS, MySql, Redis, and OpenEC on each cloud host. Hence, it is suggested that the number of CPUs in each cloud host is not smaller than 2.",
  formLabel4: "System Disk",
  inputTip4:
    "Explanation: JCDT system needs to run HDFS, MySql, Redis, and OpenEC on each cloud host. Hence, it is suggested that the system disk in each cloud host is not smaller than 2 GB.",
  formLabel5: "Data Disk",
  inputTip5: "Explanation: the maximum storage space * (1+redundancy)= the number of cloud hosts * the size of the data disk of each cloud host.",
  formLabel6: "Memory",
  inputTip6:
    "Explanation: JCDT system needs to run HDFS, MySql, Redis, and OpenEC on each cloud host. Hence, it is suggested that the memory in each cloud host is not smaller than 2 GB.",
  formLabel7: "Bandwidth",
  inputTip7:
    "Explanation: in the JCDT system, the efficiency of data writing, data read and data repair depends on the bandwidth. Hence it is suggested that the bandwidth is not less than 2 Mbps.",
  must: "Must",
  caseTab1: "Intelligent recommendation scheme",
  caseTab2: "Custom scheme",
  sourceList: "Storage resource List",
  technicalSpot: "Highlights",
  technicalSpotHtml: `
  <p>The existing erasure-coded storage systems are mainly designed for a single cloud data center (CDC) environment. In a cross-CDC, they need to transfer a large amount of data across CDCs (via the public network) to write, read, and repair data. Since the bandwidth of the public network is much lower than the bandwidth of the internal network, their speed of writing, reading, and repairing data is low. This JCDT system includes a data writing method, a data read method, and a data repair method designed for cross-CDC erasure codes, which can improve the speed of data writing, data read, and data repair by minimizing the cross-CDC traffic.
  </p>
  <h4>1. Data writing method of cross-CDC erasure codes</h4>
  <p>This method selects encoding nodes by comparing the number of coded blocks in the receiving end CDC and the sending end CDC so that the cross-CDC traffic is always equal to the minimum number of coded blocks in the receiving end CDC and the sending end CDC. Hence it can effectively reduce the cross-CDC traffic of data writing, thereby increasing the writing speed. In addition, because this method organizes the data writing process into a two-stage pipeline (the first stage: the client writes the data blocks; the second stage: the storage nodes coordinate to complete the encoding of parity blocks), it can further improve the writing speed of batched data.
  </p>
  <h4>2. Data read method of cross-CDC erasure codes</h4>
  <p>This method retrieves the original data by reading and decoding several coded blocks closest to the client. This reading method achieves the purpose of reducing the cross-CDC traffic at the cost of increasing computing overhead. Since the transmission time of cross-CDC erasure codes is much longer than the computing time, this method can increase the read speed.
  </p>
  <h4>3. Single cloud host repair method of cross-CDC erasure codes</h4>
  <p>This method uses a tree-star hybrid topology to organize data transmission and data merging when repairing a single cloud host: first, it uses the star transmission topology to merge the auxiliary blocks in each CDC. Then, it builds a minimum network distance spanning tree to organize the transmission of auxiliary blocks merged in each CDC. This method makes the number of blocks transferred across CDC to repair each block in the failed cloud host is not greater than the number of CDCs minus one, so the repair speed can be effectively improved.
  </p>
  <h4>4. Single CDC repair method for cross-CDC erasure codes</h4>
  <p>This method selects the node that merges the auxiliary blocks by comparing the number of blocks to be repaired in the receiving end CDC and the number of auxiliary blocks in the sending end CDC so that the cross-CDC traffic is always equal to the minimum value of the number of blocks to be repaired in the receiving end CDC and the number of auxiliary blocks in the sending end CDC. Therefore, it can effectively reduce the cross-CDC traffic of repairing a CDC, thereby increasing the repair speed of the CDC.
  </p>`,
  documentTitle1: "Erasure code selection and the requirement of cloud hosts",
  documentTitle2: "Network topology",
  documentTitle3: "Manual deployment process (completed when purchasing resources)",
  documentTitle4: "How to use",
  disasterRecoRequire: `<h1 class="page-title">Erasure code selection and the requirement of cloud hosts</h1>
  <p>DFC-12,6,3,3	The JCDT system with this erasure code requires 13 cloud hosts (12 storage hosts and 1 manager host) distributed in 3 areas. It can tolerate the failure of arbitrary 3 storage hosts or 1 area. The redundancy is 1.  The repair efficiency is relatively high.
</p><p>DFC-8,3,3,2	The JCDT system with this erasure code requires 9 cloud hosts (8 storage hosts and 1 manager host) distributed in 2 areas. It can tolerate the failure of arbitrary 3 storage hosts or 1 area. The redundancy is 5/3.  The repair efficiency is relatively high.
</p><p>RS-6,3,3,2	The JCDT system with this erasure code requires 7 cloud hosts (6 storage hosts and 1 manager host) distributed in 2 areas. It can tolerate the failure of arbitrary 3 storage hosts or 1 area. The redundancy is 1.
</p><p>RS-6,4,2,3	The JCDT system with this erasure code requires 7 cloud hosts (6 storage hosts and 1 manager host) distributed in 3 areas. It can tolerate the failure of arbitrary 2 storage hosts or 1 area. The redundancy is 1/2.
</p><p>RS-8,6,2,4	The JCDT system with this erasure code requires 9 cloud hosts (8 storage hosts and 1 manager host) distributed in 4 areas. It can tolerate the failure of arbitrary 2 storage hosts or 1 area. The redundancy is 1/3.
</p><p>RS-9,6,3,3	The JCDT system with this erasure code requires 10 cloud hosts (9 storage hosts and 1 manager host) distributed in 3 areas. It can tolerate the failure of arbitrary 3 storage hosts or 1 area. The redundancy is 1/2.
</p><p>RS-12,6,6,2	The JCDT system with this erasure code requires 13 cloud hosts (12 storage hosts and 1 manager host) distributed in 2 areas. It can tolerate the failure of arbitrary 6 storage hosts or 1 area. The redundancy is 1.
</p><p>RS-12,8,4,3	The JCDT system with this erasure code requires 13 cloud hosts (12 storage hosts and 1 manager host) distributed in 3 areas. It can tolerate the failure of arbitrary 4 storage hosts or 1 area. The redundancy is 1/2.
</p>`
};

export { base, disasterReco };
