export default {
  formatBytes(bytes, decimals = 2) {
    if (bytes === 0) return "0 Bytes";

    const k = 1024;
    const dm = decimals < 0 ? 0 : decimals;
    const sizes = ["Bytes", "KB", "MB", "GB", "TB", "PB", "EB", "ZB", "YB"];

    const i = Math.floor(Math.log(bytes) / Math.log(k));

    return `${parseFloat((bytes / k ** i).toFixed(dm))} ${sizes[i]}`;
  },
  sleep(time) {
    return new Promise(resolve => setTimeout(resolve, time));
  },
  underline(data) {
    if (typeof data === "string") {
      return data.replace(/([A-Z])/g, (p, m) => `_${m.toLowerCase()}`);
    }
    if (Array.isArray(data)) {
      return data.map(item => this.underline(item));
    }
    return undefined;
  },
  underlineUpperCase(data) {
    if (typeof data === "string") {
      return this.underline(data).toUpperCase();
    }
    if (Array.isArray(data)) {
      return data.map(item => this.underline(item));
    }
    return undefined;
  }
};
