<template>
  <div class="app-container">
    <div class="refresh-btn">
      <el-button type="primary" @click="fetchData"><i class="el-icon-refresh"></i></el-button>
    </div>
    <el-upload action="foobar" :http-request="handleUpload">
      <el-button type="primary"><i class="el-icon-upload2"></i> 上传文件</el-button>
    </el-upload>
    <el-table
      v-loading="listLoading"
      :data="files.filter(data => !search || data.FileInfo.FileName.toLowerCase().includes(search.toLowerCase()))"
      fit
    >
      <el-table-column type="expand">
        <template slot-scope="props">
          <el-form label-position="left" inline class="demo-table-expand">
            <el-form-item label="同步状态">
              <span>{{ props.row.FileInfo.SyncStatus }}</span>
            </el-form-item>
            <el-form-item label="重建状态">
              <span>{{ reconstructStatusFormatter(props.row.FileInfo.ReconstructStatus) }}</span>
            </el-form-item>
            <el-form-item label="重建时间">
              <span>{{ reconstructDateFormatter(props.row.FileInfo.ReconstructStatus, props.row.FileInfo.LastReconstructed) }}</span>
            </el-form-item>
          </el-form>
        </template>
      </el-table-column>
      <el-table-column label="文件名" prop="FileInfo.FileName" :formatter="filenameFormatter" />
      <el-table-column label="大小" prop="FileInfo.Size" :formatter="sizeFormatter" />
      <el-table-column label="修改时间" prop="FileInfo.LastModified" :formatter="dateFormatter" />
      <el-table-column align="right">
        <template slot="header" slot-scope="scope">
          <el-input v-model="search" size="mini" placeholder="搜索文件" />
        </template>
        <template slot-scope="scope">
          <el-button size="mini" type="primary" @click="handleDownload(scope.row.FileInfo.FileName)">下载</el-button>
          <el-button size="mini" type="danger" @click="handleDelete(scope.row.FileInfo.FileName)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>
  </div>
</template>

<script>
import cloudStorage from "@/api/cloudStorage";
import { Message } from "element-ui";
import other from "@/utils/other";

export default {
  data() {
    return {
      files: [],
      search: "",
      listLoading: false
    };
  },
  created() {
    this.fetchData();
  },
  methods: {
    fetchData() {
      this.listLoading = true;
      cloudStorage.getFiles("/").then(response => {
        setTimeout(() => {
          if (response.Files != null) {
            this.files = response.Files;
          } else {
            this.files = [];
          }
          this.listLoading = false;
        }, 300);
      });
    },

    handleUpload(item) {
      const self = this;
      cloudStorage.getUploadAddress(item.file.name).then(response => {
        const token = response.Token;
        const addr = `${window.location.protocol}//${window.location.hostname}/upload`;
        cloudStorage.upload(item, token, addr).then(() => {
          self.fetchData();
        });
      });
    },

    handleDownload(filename) {
      cloudStorage.getDownloadAddress(filename).then(response => {
        const type = response.Type;
        let url = response.Result;
        if (type === "url") {
          if (!url.startsWith("http")) {
            url = `${window.location.protocol}//${window.location.hostname}${url}`;
          }
          const link = document.createElement("a");
          link.href = url;
          link.setAttribute("download", filename);
          document.body.appendChild(link);
          link.click();
          URL.revokeObjectURL(link);
          document.body.removeChild(link);
        } else if (type === "tid") {
          Message.info("正在重建文件");
        }
      });
    },

    handleDelete(filename) {
      const self = this;
      cloudStorage.deleteFile(filename).then(() => {
        self.fetchData();
      });
    },

    sizeFormatter(row, column, bytes) {
      return other.formatBytes(bytes);
    },

    dateFormatter(row, column, timestamp) {
      const date = new Date(timestamp);
      return date.toLocaleString("zh-CN");
    },

    filenameFormatter(row, column, name) {
      let newName = name;
      while (newName.charAt(0) === "/") {
        newName = newName.substring(1);
      }
      while (newName.charAt(newName.length - 1) === "/") {
        newName = newName.substring(0, name.length - 1);
      }
      return newName;
    },

    reconstructStatusFormatter(status) {
      if (status.length === 0) {
        return "未执行";
      }
      return status;
    },

    reconstructDateFormatter(status, timestamp) {
      if (status.length === 0) {
        return "-";
      }
      const date = new Date(timestamp);
      return date.toLocaleString("zh-CN");
    }
  }
};
</script>

<style>
.el-tag + .el-tag {
  margin-left: 5px;
}

.demo-table-expand {
  font-size: 0;
}
.demo-table-expand label {
  width: 90px;
  color: #99a9bf;
}
.demo-table-expand .el-form-item {
  margin-right: 0;
  margin-bottom: 0;
  width: 100%;
}
.app-container {
  text-align: left;
}
.refresh-btn {
  float: left;
  margin-right: 5px;
}
</style>
