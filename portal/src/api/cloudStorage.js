import request from "@/utils/request";
import store from "@/store";

const BACKEND_URL = "http://localhost:8081";
// 云存储接口
export default {
  makeBucket(query) {
    return request({
      url: "/storage/makeBucket",
      method: "post",
      params: query
    });
  },
  listBuckets() {
    return request({
      url: "/storage/listBuckets",
      method: "get"
    });
  },
  listCloudFiles() {
    return request({
      url: "/storage/listCloudFiles",
      method: "get"
    });
  },
  uploadFile(query) {
    return request({
      url: "/storage/uploadFile",
      method: "post",
      params: query
    });
  },
  downloadFile() {
    return request({
      url: "/storage/downloadFile",
      method: "get"
    });
  },
  deleteFile(query) {
    return request({
      url: "/storage/deleteFile",
      method: "post",
      params: query
    });
  },
  updatePreference(form) {
    return request({
      url: `${BACKEND_URL}/user/changeUserPreference`,
      method: "post",
      data: {
        AccessToken: store.getters.token,
        ...form
      }
    });
  },
  getStoragePlans() {
    return request({
      url: `${BACKEND_URL}/plan/getAllStoragePlan`,
      method: "post",
      data: {
        AccessToken: store.getters.token
      }
    });
  },
  changeStoragePlan(StoragePlan) {
    return request({
      url: `${BACKEND_URL}/plan/chooseStoragePlan`,
      method: "post",
      data: {
        AccessToken: store.getters.token,
        StoragePlan
      }
    });
  },

  getFiles(path) {
    return request({
      url: `${BACKEND_URL}/file/getFiles`,
      method: 'post',
      data: {
        AccessToken: store.getters.token,
        FilePath: path
      }
    })
  },

  getUploadAddress(path) {
    return request({
      url: `${BACKEND_URL}/file/preUploadFile`,
      method: 'post',
      data: {
        AccessToken: store.getters.token,
        FilePath: path
      }
    })
  },
  
  upload(item, token, url) {
    var form_data = new FormData()
    form_data.append('file', item.file, item.filename)
    form_data.append("token", token)
    return request({
      url: url,
      method: 'post',
      headers: {
        'Content-Type': 'multipart/form-data'
      },
      timeout: 0,
      data: form_data,
      onUploadProgress: progressEvent => {
        const complete = (progressEvent.loaded / progressEvent.total * 100 | 0)
        item.onProgress({ percent: complete })
      }
    })
  },

  getDownloadAddress(filename) {
    return request({
      url: `${BACKEND_URL}/file/downloadFile`,
      method: 'post',
      data: {
        AccessToken: store.getters.token,
        FilePath: filename
      }
    })
  },

  download(url) {
    return request({
      url: url,
      method: 'get',
    })
  },
  
  deleteFile(filename) {
    return request({
      url: `${BACKEND_URL}/file/deleteFile`,
      method: 'post',
      data: {
        AccessToken: store.getters.token,
        FilePath: path
      }
    })
  },
  
  genDownloadLink(filename) {
    // TODO: insecure, should use temporary token
    return process.env.VUE_APP_BASE_API + '/storage/download?filename=' + filename + '&t=' + getToken()
  }
};
