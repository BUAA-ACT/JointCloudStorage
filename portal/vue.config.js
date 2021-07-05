/* eslint-disable */
// const url = process.env.API_BASE_URL;
// const isProduction = process.env.NODE_ENV === 'production';
// compression-webpack-plugin 打包的时候开启gzip可以很大程度减少包的大小，非常适合于上线部署。
// const UglifyJsPlugin = require('uglifyjs-webpack-plugin');
const TerserPlugin = require("terser-webpack-plugin");
const FileManagerPlugin = require("filemanager-webpack-plugin");

const path = require("path");

function resolve(dir) {
  return path.join(__dirname, dir);
}

module.exports = {
  publicPath: "./",
  lintOnSave: true,
  productionSourceMap: false,
  chainWebpack: config => {
    // hmr fix
    config.resolve.symlinks(true);
    // 添加别名
    config.resolve.alias.set("@", resolve("src"));
    const cdn = {
      // 访问https://unpkg.com/element-ui/lib/theme-chalk/index.css获取最新版本
      css: ["https://unpkg.com/element-ui@2.15.1/lib/theme-chalk/index.css"],
      js: [
        "https://cdn.bootcdn.net/ajax/libs/vue/2.6.11/vue.min.js",
        "https://cdn.bootcdn.net/ajax/libs/echarts/4.9.0-rc.1/echarts.min.js",
        "https://cdn.bootcdn.net/ajax/libs/element-ui/2.15.1/index.js",
        "https://cdn.bootcdn.net/ajax/libs/vuex/3.5.1/vuex.min.js",
        "https://cdn.bootcdn.net/ajax/libs/vue-router/3.4.4/vue-router.min.js",
        "https://cdn.bootcdn.net/ajax/libs/axios/0.20.0-0/axios.min.js",
        "https://cdn.jsdelivr.net/npm/echarts@4.9.0/dist/echarts.js"
      ]
    };
    // 如果使用多页面打包，使用vue inspect --plugins查看html是否在结果数组中
    config.plugin("html").tap(args => {
      // html中添加cdn
      const webpackArg = [...args];
      webpackArg[0].cdn = process.env.NODE_ENV === "production" ? cdn : [];
      return webpackArg;
    });
  },
  configureWebpack: config => {
    if (process.env.NODE_ENV === "production") {
      config.externals = {
        vue: "Vue",
        "element-ui": "ELEMENT",
        "vue-router": "VueRouter",
        vuex: "Vuex",
        axios: "axios",
        echarts: "echarts"
      };
      // 去掉注释
      config.optimization.minimizer.push(
        new TerserPlugin({
          // 使用 TerserPlugin 替代 UglifyJsPlugin
          terserOptions: {
            format: {
              comments: false // 去掉注释
            },
            compress: {
              // warnings: false,
              drop_console: true,
              drop_debugger: false,
              pure_funcs: ["console.log", "this.$log"] // 移除console
            }
          }
        })
      );
      // config.plugins.push(
      //   new FileManagerPlugin({
      //     events: {
      //       onEnd: {
      //         delete: ["./jointCloud.zip"],
      //         archive: [
      //           {
      //             source: "./dist",
      //             destination: "./jointCloud.zip"
      //           }
      //         ]
      //       }
      //     }
      //   })
      // );
    }
  },
  // 配置转发代理
  devServer: {
    proxy: {
      "/file": {
        target: `http://127.0.0.1:8081`
      },
      "/plan": {
        target: `http://127.0.0.1:8081`
      },
      "/user": {
        target: `http://127.0.0.1:8081`
      },
      "/cloud": {
        target: `http://127.0.0.1:8081`
      },
      "/upload": {
        target: `http://192.168.105.13:8083`
      }
    }
  }
};
