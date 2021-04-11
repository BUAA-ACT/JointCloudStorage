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
      :data="files.filter(data => !search || data.filename.toLowerCase().includes(search.toLowerCase()))"
      fit
    >
      <el-table-column
        label="文件名"
        prop="filename"
        min-width="200"
      />
      <el-table-column
        label="大小"
        prop="size"
        :formatter="sizeFormatter"
      />
      <el-table-column
        label="修改时间"
        prop="last_modified"
        :formatter="dateFormatter"
      />
      <el-table-column class-name="status-col" label="位置">
        <template slot-scope="scope">
          <el-tag v-for="loc in scope.row.sites" :key="loc" type="info">{{ loc }}</el-tag>
        </template>
      </el-table-column>
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
            @click="handleDownload(scope.row.filename)"
          >下载</el-button>
          <el-button
            size="mini"
            type="danger"
            @click="handleDelete(scope.row.filename)"
          >删除</el-button>
        </template>
      </el-table-column>
    </el-table>
  </div>
</template>

<script>
import cloudStorage from "@/api/cloudStorage";

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
        console.log("token=", token)
        cloudStorage.upload(item, token, "http://localhost:8083/upload").then(() => {
          self.fetchData()
        })
      })
    },
    handleDownload(filename) {
      var url = genDownloadLink(filename)
      var link = document.createElement('a')
      link.href = url
      link.setAttribute('download', filename)
      document.body.appendChild(link)
      link.click()
      URL.revokeObjectURL(link)
      document.body.removeChild(link)
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
      var date = new Date(timestamp * 1000)
      return date.toLocaleString('zh-CN')
    }
  }
}
</script>

<style>
.el-tag+.el-tag{
  margin-left: 5px;
}
</style>
