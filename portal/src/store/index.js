import Vue from "vue";
import Vuex from "vuex";
import Common from "@/api/common";
import { getToken, setToken, removeToken } from "@/utils/auth";
import user from "@/utils/user";
import Other from "@/utils/other";

Vue.use(Vuex);

const ROLE_ADMIN = "ADMIN";
const ROLE_GUEST = "GUEST";

// const store = new Vuex.Store({
export default new Vuex.Store({
  state: {
    token: null,
    name: null,
    nickname: null,
    status: null,
    preference: { AllowDelay: false, Availability: 0, Choice: 0, StoragePrice: 0, TrafficPrice: 0, Vendor: 0 },
    storagePlan: { StorageMode: "Replica", N: 0, K: 0, Clouds: [], StoragePrice: 0, TrafficPrice: 0, Availability: 0 },
    dataStats: {
      Volume: 0
    },
    role: [],
    ready: false
  },
  getters: {
    token: state => {
      const localToken = getToken();
      if (localToken) {
        state.token = localToken;
      }
      return state.token;
    },
    name: state => state.name,
    vendor: state => state.preference.Vendor,
    status: state => state.status,
    dataStats: state => {
      const { storagePlan, dataStats } = state;
      if (!dataStats || dataStats === {}) {
        return null;
      }
      let { Volume } = dataStats;
      const { StorageMode, K } = storagePlan;
      if (StorageMode === "EC") {
        Volume /= K;
      }
      const cloudsDetails = storagePlan.Clouds.map(val => {
        return {
          CloudID: val.CloudID,
          Location: val.Location,
          UploadTraffic: dataStats.UploadTraffic[val.CloudID] || 0,
          DownloadTraffic: dataStats.DownloadTraffic[val.CloudID] || 0
        };
      });
      return { Volume, cloudsDetails };
    },
    haveStoragePlan: state => state.storagePlan.N > 0,
    preference: state => state.preference,
    storagePlan: state => state.storagePlan,
    ready: state => state.ready,
    role: state => state.role,
    isAdmin: state => state.role === ROLE_ADMIN,
    isGuest: state => state.role === ROLE_GUEST
  },
  mutations: {
    SET_TOKEN: (state, token) => {
      setToken(token);
      state.token = token;
    },
    SET_NAME: (state, name) => {
      state.name = name;
    },
    SET_NICKNAME: (state, nickname) => {
      state.nickname = nickname;
    },
    SET_PREFERENCE: (state, preference) => {
      state.preference = preference;
    },
    SET_STORAGE_PLAN: (state, storagePlan) => {
      state.storagePlan = storagePlan;
    },
    SET_DATA_STATS: (state, dataStats) => {
      state.dataStats = dataStats;
    },
    SET_STATUS: (state, status) => {
      state.status = status;
    },
    RESET_ALL: state => {
      state.ready = true;
      state.token = null;
      removeToken();
      state.name = null;
      state.nickname = null;
      state.preference = {};
      state.storagePlan = {};
      state.dataStats = {};
      state.status = null;
    },
    SET_READY: (state, ready) => {
      state.ready = ready;
    },
    SET_ROLE: (state, role) => {
      state.role = role;
    }
  },
  actions: {
    login({ dispatch }, loginForm) {
      const { username, password } = loginForm;
      return Common.login({
        Email: username,
        Password: password
      }).then(async resp => {
        await dispatch("getInfo", resp.AccessToken);
      });
    },
    /**
     * 获取用户所有信息
     * `ready` 字段表示当前状态
     * @param getters
     * @param commit
     * @param pToken
     * @returns {Promise<void>}
     */
    async getInfo({ getters, commit }, pToken) {
      let token = pToken || null;
      if (pToken === undefined) {
        token = getters.token;
      }
      if (token === null) {
        return;
      }
      commit("SET_READY", false);
      await Common.checkToken(token)
        .then(async resp => {
          if (!resp) {
            commit("RESET_ALL");
            return;
          }
          commit("SET_TOKEN", token);
          await Common.getUserInfo(token).then(data => {
            const { Email, Nickname, Preference, StoragePlan, DataStats, Status, Role } = data.UserInfo;
            commit("SET_NAME", Email);
            commit("SET_NICKNAME", Nickname);
            commit("SET_PREFERENCE", Preference);
            commit("SET_STORAGE_PLAN", StoragePlan);
            commit("SET_DATA_STATS", DataStats);
            commit("SET_STATUS", Status);
            commit("SET_ROLE", Role);
          });
        })
        .catch()
        .then(() => {
          commit("SET_READY", true);
        });
    },
    async logout({ state, commit }) {
      Common.logout(state.token).then(() => {
        commit("RESET_ALL");
        user.logout();
      });
    },
    async checkPreference({ getters }) {
      if (getters.vendor === 0) {
        user.firstLoginNotification();
      }
    },
    /**
     * 更新用户信息某字段
     * 与 `getInfo` 不同的是，本方法不会影响 `ready` 字段
     * @param getters
     * @param commit
     * @param {"Preference" | "StoragePlan" | "DataStats" | "Status" | "Role"} field
     * @returns {Promise<void>}
     */
    async updateInfo({ getters, commit }, field) {
      let value;
      if (typeof field === "string") {
        await Common.getUserInfo(getters.token).then(data => {
          value = data.UserInfo[field];
          commit(`SET${Other.underlineUpperCase(field)}`, value);
        });
      }
      if (!value) {
        return Promise.reject(new Error("Invalid field name"));
      }
      return Promise.resolve();
    }
  },
  modules: {}
});
