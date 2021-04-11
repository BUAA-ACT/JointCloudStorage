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
    preference: { AllowDelay: false, Availability: 0, Choice: 0, StoragePrice: 0, TrafficPrice: 0, Vendor: 0 },
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
    vendor: state => state.preference.Vendor
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
    RESET_ALL: state => {
      state.token = null;
      removeToken();
      state.name = null;
      state.nickname = null;
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
          const { Email, Nickname, Preference } = data.UserInfo;
          commit("SET_NAME", Email);
          commit("SET_NICKNAME", Nickname);
          commit("SET_PREFERENCE", Preference);
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
