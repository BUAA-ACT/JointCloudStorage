/* eslint-disable */
const ak = "QFgFQorrB84maOZh0pPGC8kUiP0mGIhx";

if (!global.BMap) {
  global.BMap = {};
  global.BMap._preloader = new Promise((resolve, reject) => {
    global._initBaiduMap = function() {
      resolve(global.BMap);
      global.document.body.removeChild($script);
      global.BMap._preloader = null;
      global._initBaiduMap = null;
    };
    const $script = document.createElement("script");
    global.document.body.appendChild($script);
    $script.src = `https://api.map.baidu.com/api?v=2.0&ak=${ak}&callback=_initBaiduMap`;
  });
} else if (!global.BMap._preloader) {
  Promise.resolve(global.BMap);
}

export default {
  BMap: global.BMap,
  /**
   * Convert location string
   * @param {string} location `${longitude},${latitude}`
   * @returns {string}
   */
  async convertToCityName(location) {
    const coordinate = location.split(",");
    const [lng, lat] = coordinate;
    const map = new BMap.Map("allmap");
    const point = new BMap.Point(lng, lat);
    const gc = new BMap.Geocoder();
    let ret = "";
    await gc.getLocation(point, rs => {
      const addComp = rs.addressComponents;
      ret = `${addComp.province}, ${addComp.city}`;
    });
    return ret;
  }
};
