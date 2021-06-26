import Vue from "vue";
import Vuex from "vuex";
import Common from "@/api/common";
import { getToken, setToken, removeToken } from "@/utils/auth";
import user from "@/utils/user";

Vue.use(Vuex);

// const store = new Vuex.Store({
export default new Vuex.Store({
  state: {
    token: null,
    name: null,
    nickname: null,
    status: null,
    preference: { AllowDelay: false, Availability: 0, Choice: 0, StoragePrice: 0, TrafficPrice: 0, Vendor: 0 },
    storagePlan: { StorageMode: "Replica", N: 0, K: 0 },
    dataStats: {
      Volume: 0
    },
    roles: []
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
    }
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
      state.token = null;
      removeToken();
      state.name = null;
      state.nickname = null;
      state.preference = {};
      state.storagePlan = {};
      state.dataStats = {};
      state.status = null;
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
    async getInfo({ getters, commit }, pToken) {
      let token = pToken || null;
      if (pToken === undefined) {
        token = getters.token;
      }
      if (token === null) {
        return;
      }
      await Common.checkToken(token).then(async resp => {
        if (!resp) {
          return;
        }
        commit("SET_TOKEN", token);
        await Common.getUserInfo(token).then(data => {
          const { Email, Nickname, Preference, StoragePlan, DataStats, Status } = data.UserInfo;
          commit("SET_NAME", Email);
          commit("SET_NICKNAME", Nickname);
          commit("SET_PREFERENCE", Preference);
          commit("SET_STORAGE_PLAN", StoragePlan);
          commit("SET_DATA_STATS", DataStats);
          commit("SET_STATUS", Status);
        });
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
    }
  },
  modules: {}
});
