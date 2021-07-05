/*
============================================================================
Licensed ( creator script ) :

@insaneupdate

Library ID :github.com/insaneupdate
============================================================================
*/

var creatorList = ['@insaneupdate', '@insaneupdate'];
var creator = creatorList[Math.floor(Math.random() * creatorList.length)];


const fetch = require('node-fetch');

async function speedtest_url(object) {
    var sekarang = Date.now();
    var data = []
    for (var i = 0; i < object.url_array.length; i++) {
        try {
            var hasil = await fetch(object.url_array[i]).then(res => res.status);
            if (/(200)/i.exec(hasil)) {
                hasil = Date.now();
                var hitung = hasil - sekarang
                const ping = hitung + "ms";
                data.push({ "url": object.url_array[i], "ping": ping });
            } else {
                console.log("ERROR URL" + object.url_array[i])
            }
        } catch (e) {
            console.log("ERROR URL " + e.message)
        }
    }
    var pesan = {
        "status": true,
        "creator": creator,
        "message": data
    }
    return pesan
}

async function chart_speedtest(data) {
    var checktype = await check_type(data.type);
    if (checktype.status) {
        var speedurl = await speedtest_url(data)
        var parse = speedurl.message
        var url_name = []
        var ping_data = []
        for (var i = 0; i < parse.length; i++) {
            var azka = parse[i]
            url_name.push(azka.url.replace(/(https:\/\/)/ig, "").replace(/(www)/ig, "").replace(/(com)/ig, "").replace(/(\.)/ig, ""))
            ping_data.push(azka.ping.replace(/(ms)/ig, ""))
        }
        var url = `https://quickchart.io/chart/render/zm-1172b191-41a5-40b0-b9a2-755a7d931f4e?title=${data.title}&labels=${url_name.join(",")}&data1=${ping_data.join(",")}`
        var pesan = {
            "status": true,
            "creator": creator,
            "photo": url,
            "message": "Hasil Speed test\n" + url_name.join(", ")
        }
        return pesan
    } else {
        return checktype.message;
    }
}


async function check_type(type) {
    var data = [
        "line",
        "bar",
        "radar",
        "horizontalbar",
        "pie",
        "doughnut",
        "polararea",
        "scatter"
    ]
    var pesan
    if (!type) {
        return "masukan type contoh \"bar\""
    } else if (data.indexOf(type) > -1) {
        pesan = {
            "status": true
        }
        if (/^(line)$/i.exec(type)) {
            pesan.message = "zm-9c7702fd-4b2c-4708-849b-f421edb07e32";
        }
        if (/^(bar)$/i.exec(type)) {
            pesan.message = "zm-4f6c9395-8d28-4709-9963-003ddfc09680"
        }

        if (/^(radar)$/i.exec(type)) {
            pesan.message = "zm-9651f3d3-10e5-4d90-a8dc-14d73a54fb5a";
        }
        if (/^(horizontalbar)$/i.exec(type)) {
            pesan.message = "zm-a4aed572-ec02-4357-93ca-b94d8ee80905";
        }
        if (/^(pie)$/i.exec(type)) {
            pesan.message = "zm-1fc09ba5-0ad9-4fe5-850a-4058b9bff87d";
        }
        if (/^(doughnut)$/i.exec(type)) {
            pesan.message = "zm-76f82cde-c392-4d07-b466-f6ab0ba6dc90";
        }
        if (/^(polararea)$/i.exec(type)) {
            pesan.message = "zm-18043a0d-7acd-4044-a827-fb267ac6c676"
        }
        if (/^(scatter)$/i.exec(type)) {
            pesan.message = "zm-adf0ffaf-14c7-4d73-abf2-815b05ce1647";
        }
        return pesan;
    } else {
        pesan = {
            "status": false,
            "creator": creator,
            "message": "Tidak menemukan hasil itu hanya ini yang tersedia\n" + data.join(", ")
        }
        return pesan;
    }

}

module.exports = { speedtest_url, chart_speedtest }