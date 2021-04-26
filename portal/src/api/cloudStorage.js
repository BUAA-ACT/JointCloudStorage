import request from "@/utils/request";
import store from "@/store";

// const BACKEND_URL = "http://localhost:8081";
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
  updatePreference(form) {
    return request({
      url: `/user/changeUserPreference`,
      method: "post",
      data: {
        AccessToken: store.getters.token,
        ...form
      }
    });
  },

  getFiles(path) {
    return request({
      url: `/file/getFiles`,
      method: "post",
      data: {
        AccessToken: store.getters.token,
        FilePath: path
      }
    });
  },

  getUploadAddress(path) {
    return request({
      url: `/file/preUploadFile`,
      method: "post",
      data: {
        AccessToken: store.getters.token,
        FilePath: path
      }
    });
  },

  upload(item, token, url) {
    const formData = new FormData();
    formData.append("file", item.file, item.filename);
    formData.append("token", token);
    return request({
      url,
      method: "post",
      headers: {
        "Content-Type": "multipart/form-data"
      },
      timeout: 0,
      data: formData,
      onUploadProgress: progressEvent => {
        const complete = (progressEvent.loaded / progressEvent.total) * 100 || 0;
        item.onProgress({ percent: complete });
      }
    });
  },

  getDownloadAddress(filename) {
    return request({
      url: `/file/downloadFile`,
      method: "post",
      data: {
        AccessToken: store.getters.token,
        FilePath: filename
      }
    });
  },

  download(url) {
    return request({
      url,
      method: "get"
    });
  },

  deleteFile(filename) {
    return request({
      url: `/file/deleteFile`,
      method: "post",
      data: {
        AccessToken: store.getters.token,
        FilePath: filename
      }
    });
  }
};
