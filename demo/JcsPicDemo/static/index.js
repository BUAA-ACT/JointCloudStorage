Vue.createApp({
  setup() {
    let info = Vue.reactive({data: {}});
    setInterval(() => axios.get("/info").then(
      resp => {
        info.data = resp.data;
      }
    ).catch(e =>
      console.log(e)
    ), 100)
    return {
      info
    }
  },
  template: `
<nav style="text-align: center" class="navbar navbar-inverse navbar-fixed-top shadow-sm">
    <a class="navbar-brand">基于云际网盘的媒体资源流水线处理引擎</a>
</nav>
<div class="row" style="margin-top: 70px">
  <div v-if="info.data.node_states.length === 0" style="text-align: center">
    <a href="/start"><button class="btn btn-primary shadow-lg" style="margin-top: 15vh; font-size: 36px">启动流水线</button></a>
  </div>
  <div class="col-md-8 col-md-offset-2" v-else>
    <div class="panel panel-default" v-for="(state, index) in info.data.node_states" :key="state.endpoint_name">
      <div class="panel-heading" style="height: 50px">
        <span class="row" style>
        <span class="col-md-3">节点{{ index + 1 }}：{{ state.endpoint_name }}</span>
        <span :class="'text-' + (state.state === 'OK' ? 'success' : 'danger') + ' col-md-3 col-md-offset-1'">
            当前状态：<i :class="'glyphicon glyphicon-' + (state.state === 'OK' ? 'ok' : 'remove')" />
        </span>
        <a class="col-md-2 col-md-offset-3" :href="state.endpoint_address">接入点</a>
        </span>
      </div>
      <div class="panel-body">
        <table style="text-align: center" class="table">
          <tbody>
            <th class="success">成功</th>
            <th class="warning">失败</th>
            <tr>
              <td>{{ state.finish_num }}</td>
              <td>{{ state.fail_num }}</td>
            </tr>
          </tbody>
        </table>
        <span class="text-info">正在处理：{{ state.file_processing }}</span>
      </div>
    </div>
  </div>
</div>
        `
})
  .mount("#app")