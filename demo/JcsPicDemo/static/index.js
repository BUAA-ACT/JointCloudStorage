Vue.createApp({
  setup() {
    let info = Vue.reactive({data: {}});
    // setInterval(() =>
    //     axios.get("/info").then(
    //       resp => {
    //         info.data = resp.data;
    //       }
    //     ).catch(e =>
    //       console.log(e)
    //     )
    //   , 100)
    info = {
      data: {
        "py/object": "app.Info",
        "runner_state": "running",
        "node_states": [{
          "py/object": "node.node.NodeState",
          "finish_num": 165,
          "fail_num": 0,
          "endpoint_name": "呼和浩特",
          "endpoint_address": "http://jsi-aliyun-hohhot.jointcloudstorage.cn/",
          "file_processing": "picture/6.jpg"
        }, {
          "py/object": "node.node.NodeState",
          "finish_num": 165,
          "fail_num": 0,
          "endpoint_name": "青岛",
          "endpoint_address": "http://jsi-aliyun-qingdao.jointcloudstorage.cn/",
          "file_processing": "pic2021-09-07-20-05-30.jpg",
          "state": "OK"
        }, {
          "py/object": "node.node.NodeState",
          "finish_num": 163,
          "fail_num": 3,
          "endpoint_name": "成都",
          "endpoint_address": "http://jsi-txyun-chengdu.jointcloudstorage.cn/",
          "file_processing": "pic2021-09-07-20-05-22.jpg",
          "state": "OK"
        }]
      }
    }
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
  <div class="col-8" v-else>
    <div class="card" v-for="(state, index) in info.data.node_states" :key="state.endpoint_name">
      <div class="card-body">
        <div class="card-title" style="height: 50px">
          <span class="row justify-content-between">
            <span class="col-3">节点{{ index + 1 }}：{{ state.endpoint_name }}</span>
            <span class="col-6 grey-link">{{ state.endpoint_address }}</span>
            <span :class="'bg-' + (state.state === 'OK' ? 'success' : 'danger') + ' col'"
                  style="border: 1px solid #ccc; border-radius: 4px;text-align: center; color: white; max-width: 130px">
                当前状态：<i :class="'bi-' + (state.state === 'OK' ? 'check' : 'x')"></i>
            </span>
          </span>
        </div>
        <div class="card-text">
          <table class="table detail-table">
            <tbody>
              <th class="success">成功</th>
              <th class="fail">失败</th>
              <tr>
                <td class="success">{{ state.finish_num }}</td>
                <td class="fail">{{ state.fail_num }}</td>
              </tr>
            </tbody>
          </table>
          <span class="text-info">正在处理：{{ state.file_processing }}</span>
        </div>
      </div>
      
    </div>
  </div>
</div>
        `
})
  .mount("#app")