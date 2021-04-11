import { Notification, Message } from "element-ui";
import Vue from "vue";

export default {
  firstLoginNotification() {
    Notification.warning({
      title: "警告",
      message: "您还没有设置存储方案，请点击本通知前往设置存储偏好",
      onClick: () => {
        Vue.$router.push({ path: "/cloudStorage/userPreference" });
      }
    });
  },
  logout() {
    Message.success("您已成功退出！");
  }
};
