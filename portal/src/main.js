/* eslint-disable */
import Vue from "vue";
import ElementUI from "element-ui";
import "element-ui/lib/theme-chalk/index.css";
import enLocale from "element-ui/lib/locale/lang/en";
import zhLocale from "element-ui/lib/locale/lang/zh-CN";
import VueRouter from "vue-router";
import VueI18n from "vue-i18n";
import appApi from "@/api/actions";
import App from "./App.vue";
import store from "./store";
import router from "./router";
import "@/styles/index.scss"; // global css

// Vue.use(ElementUI);
Vue.use(VueRouter);
Vue.use(VueI18n);
Vue.config.productionTip = false;

Vue.prototype.$Api = appApi;

const i18n = new VueI18n({
  locale: "zh-CN", // 语言标识
  // this.$i18n.locale // 通过切换locale的值来实现语言切换
  messages: {
    "zh-CN": Object.assign(require("@/utils/lang/zh"), zhLocale), // 中文语言包
    "en-US": Object.assign(require("@/utils/lang/en"), enLocale) // 中文语言包
    // 'en-US': {
    //   require('@/utils/lang/en'),
    //   ...enLocale,
    // }, // 英文语言包
  }
});
Vue.use(ElementUI, {
  i18n: (key, value) => i18n.t(key, value),
  size: "small",
  zIndex: 3000
});

new Vue({
  router,
  store,
  i18n,
  render: h => h(App)
}).$mount("#app");
