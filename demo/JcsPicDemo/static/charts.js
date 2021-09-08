var chart1 = echarts.init(document.getElementById('graph1'));
var chart2 = echarts.init(document.getElementById('graph2'));
var chart3 = echarts.init(document.getElementById('graph3'));

// 指定图表的配置项和数据
var now = new Date();
var data = [[[now, 0]], [[now, 0]], [[now, 0]]];

// for (var i = 1; i < 30; i++) {
//     addData();
// }

option = {
  xAxis: {
    type: 'time',
    boundaryGap: false,
    minInterval: 2,
    min: 'dataMin',
    axisLabel: {
      interval: 2
    }
    //data: date
  },
  yAxis: {
    boundaryGap: [0, '50%'],
    type: 'value',
    min: 'dataMin',
  },
  series: [
    {
      name: '成功',
      type: 'line',
      smooth: false,
      symbol: 'none',
      stack: 'a',
      areaStyle: {
        normal: {}
      },
      data: data,
    }
  ],
  grid: {
    x: 30,
    y: 30,
    x2: 30,
    y2: 30
  }
};

setInterval(function () {
  axios.get("/info").then(
    resp => {
      for (let i = 0; i < 3; i++) {
        now = new Date()
        data[i].push([now, resp.data.node_states[i].finish_num]);
        if (data[i].length > 30) {
          data[i].shift();
        }
      }
      console.log(data[2])
      chart1.setOption({
        xAxis: {
          //data: date
        },
        series: [{
          name: '成功',
          data: data[0]
        }]
      });
      chart2.setOption({
        xAxis: {
          //data: date
        },
        series: [{
          name: '成功',
          data: data[1]
        }]
      });
      chart3.setOption({
        xAxis: {
          //data: date
        },
        series: [{
          name: '成功',
          data: data[2]
        }]
      });
    }
  ).catch(e =>
    console.log(e)
  )
}, 500);

option && chart1.setOption(option);
option && chart2.setOption(option);
option && chart3.setOption(option);
