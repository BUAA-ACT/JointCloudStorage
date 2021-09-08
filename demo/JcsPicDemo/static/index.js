Vue.createApp({
  setup() {
    let noData = Vue.ref(true);
    let info = Vue.reactive({data: {}});
    Vue.watch(
      () => info.data,
      (newVal) => {
        const curStates = newVal.node_states;
        if (curStates && curStates.length > 0 && noData.value) {
          const chartScript = document.createElement('script');
          chartScript.type = "text/javascript";
          chartScript.src = "/static/charts.js"
          document.body.appendChild(chartScript)
          noData.value = false;
        }
      },
    )
    setInterval(() =>
        axios.get("/info").then(
          resp => {
            info.data = resp.data;
          }
        ).catch(e =>
          console.log(e)
        )
      , 500)
    return {
      info
    }
  },
  template: `
<nav style="text-align: center" class="navbar navbar-dark bg-dark">
    <a class="navbar-brand" style="padding-left: 10px">基于云际网盘的媒体资源流水线处理引擎</a>
</nav>
<div class="row justify-content-center">
  <div v-if="!info.data.node_states || info.data.node_states.length === 0" style="text-align: center">
    <a href="/start"><button class="btn btn-primary shadow-lg" style="margin-top: 15vh; font-size: 36px">启动流水线</button></a>
  </div>
  <div v-else>
    <div class="col-12 row" v-for="(state, index) in info.data.node_states" :key="state.endpoint_name">
      <div class="card col-7">
        <div class="card-body">
          <div class="card-title" style="">
            <span class="row justify-content-between">
              <span class="col-10">节点{{ index + 1 }}：{{ state.endpoint_name }}</span>
              <span :class="'bg-' + (state.state === 'OK' ? 'success' : 'danger') + ' col'"
                    style="border: 1px solid #ccc; border-radius: 4px;text-align: center; color: white; max-width: 100px; margin: 0 -10px">
                  {{ state.state }}
              </span>
              <span class="col-12 grey-link">{{ state.endpoint_address }}</span>
            </span>
          </div>
          <div class="card-text">
            <table class="table detail-table">
              <tbody>
                <th class="success">成功</th>
                <th :class="state.fail_num === 0 ? 'ignore' : 'fail'">失败</th>
<!--                <th class="ignore">失败</th>-->
                <tr>
                  <td class="success">{{ state.finish_num }}</td>
                  <td :class="state.fail_num === 0 ? 'ignore' : 'fail'">{{ state.fail_num }}</td>
<!--                  <td class="ignore">{{ state.fail_num }}</td>-->
                </tr>
              </tbody>
            </table>
            <span style="color: #555">正在处理：{{ state.file_processing }}</span>
          </div>
        </div>
      </div>
      <div :id="'graph' + (index + 1)" style="height:280px" class="col-5"></div>
    </div>
  </div>
</div>
        `
})
  .mount("#app")