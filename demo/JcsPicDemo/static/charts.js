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
            let drawData = [[], [], []]
            let fails = [[], [], []]
            let start;
            for (let i = 0; i < 3; i++) {
                now = new Date()
                if (resp.data.node_states[i].state !== "ERROR") {
                    data[i].push([now, resp.data.node_states[i].finish_num, true]);
                } else {
                    data[i].push([now, resp.data.node_states[i].finish_num, false]);
                }

                if (data[i].length > 30) {
                    data[i].shift();
                }

                start = -1
                for (let j = 0; j < data[i].length; j++) {
                    d = data[i][j]
                    drawData[i].push(d.slice(0, 2))
                    if (d[2] === false && start === -1) {
                        start = j
                    }
                    if (d[2] === true && start !== -1) {
                        fails[i].push([data[i][start][0], data[i][j][0]])
                        start = -1
                    }
                }
                if (start !== -1) {
                    fails[i].push([data[i][start][0], data[data[i].length][0]])
                }
            }
            chart1.setOption({
                color: '#5470c6',
                xAxis: {
                    //data: date
                },
                series: [{
                    name: '成功',
                    data: drawData[0]
                }]
            });
            chart2.setOption({
                color: '#5470c6',
                xAxis: {
                    //data: date
                },
                series: [{
                    name: '成功',
                    data: drawData[1]
                }]
            });
            if (resp.data.node_states[2].state !== "ERROR") {
                chart3.setOption({
                    color: '#5470c6',
                    series: [{
                        name: '成功',
                        data: drawData[2]
                    }],
                });
            } else {
                chart3.setOption({
                    color: 'red',
                    series: [{
                        name: '成功',
                        data: drawData[2]
                    }],
                });
            }
        }
    ).catch(e =>
        console.log(e)
    )
}, 500);

option && chart1.setOption(option);
option && chart2.setOption(option);
option && chart3.setOption(option);
