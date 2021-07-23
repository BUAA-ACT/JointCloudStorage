import Other from "../utils/other";

/**
 * My Implement of WebSocket
 * @copyright Bear 2021.7.21
 * @licence MIT
 */
export default class MyWS {
  /**
   * WebSocket Pool
   * @type {Object}
   */
  static wsPool = {};

  /**
   * Class Instance of WebSocket
   * @type {WebSocket}
   */
  ws = null;

  /**
   * Class Instance of current WebSocket full path
   * @type {string}
   */
  fullPath = "";

  /**
   * Constructor of class MyWS
   * Use WebSocket Pool to manage a number of WebSockets
   * @param {string} path relative path of websocket server to host
   * @param {string} host host of Websocket Server
   * @param {string | number} port port of WebSocket Server
   */
  constructor(path, host = window.location.hostname, port = window.location.port) {
    let protocol;
    if (window.location.protocol === "https:") {
      protocol = "wss:";
    } else {
      protocol = "ws:";
    }
    this.fullPath = `${protocol}//${host}:${port}${path}`;
    let ws = MyWS.wsPool[this.fullPath];
    if (!(ws && ws.readyState === WebSocket.OPEN)) {
      ws = new WebSocket(this.fullPath);
      MyWS.wsPool[this.fullPath] = ws;
    }
    this.ws = ws;
  }

  /**
   * Send Message via current WebSocket
   * @param {JSON} message
   * @returns {Promise<void|*>}
   */
  async send(message) {
    if (!this.ws) {
      return Promise.reject(new Error("WebSocket 连接不存在！"));
    }
    switch (this.ws.readyState) {
      case WebSocket.OPEN:
        this.ws.send(JSON.stringify(message));
        return Promise.resolve();
      case WebSocket.CONNECTING:
        await Other.sleep(50);
        return this.send(message);
      case WebSocket.CLOSING:
      case WebSocket.CLOSED:
        return Promise.reject(new Error("连接已关闭"));
      default:
        return Promise.reject(new Error("未知错误"));
    }
  }

  /**
   * Close Current WebSocket Connection
   * And remove from WebSocket Pool
   */
  close() {
    console.log("CLOSING WS");
    this.ws.close();
    console.log(`CLOSED WS: ${this.ws.readyState}`);
    this.ws.onmessage = null;
    this.ws.onerror = null;
    this.ws.onopen = null;
    this.ws = null;
    MyWS.wsPool[this.fullPath] = null;
    this.fullPath = null;
  }

  /**
   * addAction
   * @param { Function } handler
   * @param {"message" | "error" | "open" | "close"} type
   * @returns {boolean|undefined}
   */
  addAction(handler, type) {
    if (!this.ws) {
      return undefined;
    }
    this.ws.addEventListener(type, handler);
    return true;
  }

  /**
   * removeAction
   * @param {Function} handler
   * @param {"message" | "error" | "open" | "close"} type
   * @returns {boolean|undefined}
   */
  removeAction(handler, type) {
    if (!this.ws) {
      return undefined;
    }
    this.ws.removeEventListener(type, handler);
    return true;
  }
}
