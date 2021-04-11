<template>
  <div class="app-container">
    <el-upload
      action="foobar"
      :http-request="handleUpload"
    >
      <el-button type="primary"><i class="el-icon-upload2"></i> 上传文件</el-button>
    </el-upload>
    <el-table
      v-loading="listLoading"
      :data="files.filter(data => !search || data.FileInfo.FileName.toLowerCase().includes(search.toLowerCase()))"
      fit
    >
      <el-table-column
        label="文件名"
        prop="FileInfo.FileName"
        min-width="200"
        :formatter="filenameFormatter"
      />
      <el-table-column
        label="大小"
        prop="FileInfo.Size"
        :formatter="sizeFormatter"
      />
      <el-table-column
        label="修改时间"
        prop="FileInfo.LastModified"
        :formatter="dateFormatter"
      />
      <el-table-column
        align="right"
      >
        <template slot="header" slot-scope="scope">
          <el-input
            v-model="search"
            size="mini"
            placeholder="搜索文件"
          />
        </template>
        <template slot-scope="scope">
          <el-button
            size="mini"
            type="primary"
            @click="handleDownload(scope.row.FileInfo.FileName)"
          >下载</el-button>
          <el-button
            size="mini"
            type="danger"
            @click="handleDelete(scope.row.FileInfo.FileName)"
          >删除</el-button>
        </template>
      </el-table-column>
    </el-table>
  </div>
</template>

<script>
import cloudStorage from "@/api/cloudStorage";
import { Message } from "element-ui";

export default {
  data() {
    return {
      files: [],
      search: '',
      listLoading: false
    }
  },
  created() {
    this.fetchData()
  },
  methods: {
    fetchData() {
      this.listLoading = true
      cloudStorage.getFiles("/").then(response => {
        if (response.Files != null) {
          this.files = response.Files
        } else {
          this.files = []
        }
        this.listLoading = false
      })
    },
    handleUpload(item) {
      var self = this
      cloudStorage.getUploadAddress(item.file.name).then(response => {
        var token = response.Token
        cloudStorage.upload(item, token, "http://localhost:8083/upload").then(() => {
          self.fetchData()
        })
      })
    },
    handleDownload(filename) {
      cloudStorage.getDownloadAddress(filename).then(response => {
        var type = response.Type
        var url = response.Result
        if (type == "URL") {
          var link = document.createElement('a')
          link.href = url
          link.setAttribute('download', filename)
          document.body.appendChild(link)
          link.click()
          URL.revokeObjectURL(link)
          document.body.removeChild(link)
        } else if(type == "TID") {
          Message.info("正在重建文件")
        }
      })
    },
    handleDelete(filename) {
      var self = this
      deleteFile(filename).then(() => {
        self.fetchData()
      })
    },
    sizeFormatter(row, column, bytes, index) {
      const si = false
      var thresh = si ? 1000 : 1024
      if (Math.abs(bytes) < thresh) {
        return bytes + ' B'
      }
      var units = !si
        ? ['KB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB']
        : ['KiB', 'MiB', 'GiB', 'TiB', 'PiB', 'EiB', 'ZiB', 'YiB']
      var u = -1
      do {
        bytes /= thresh
        ++u
      } while (Math.abs(bytes) >= thresh && u < units.length - 1)
      return bytes.toFixed(1) + ' ' + units[u]
    },
    dateFormatter(row, column, timestamp, index) {
      var date = new Date(timestamp)
      return date.toLocaleString('zh-CN')
    },
    filenameFormatter(row, column, name, index) {
      while(name.charAt(0)=='/') {
        name = name.substring(1);
      }
      while(name.charAt(name.length-1)=='/') {
        name = name.substring(0,name.length-1);
      }
      return name
    }
  }
}
</script>

<style>
.el-tag+.el-tag{
  margin-left: 5px;
}
</style>
