// simple reconnecting websocket client for receiving global events
type Msg = { type: string; data: any };

class WSClient {
  private url: string;
  private socket: WebSocket | null = null;
  private listeners: ((m: Msg) => void)[] = [];
  private reconnectDelay = 1000;
  private maxDelay = 5000;

  constructor(url: string) {
    this.url = url;
    this.connect();
  }

  connect() {
    try {
      this.socket = new WebSocket(this.url);
      this.socket.onopen = () => {
        this.reconnectDelay = 1000;
        // console.log("ws open");
      };
      this.socket.onmessage = (ev) => {
        try {
          const msg = JSON.parse(ev.data);
          this.listeners.forEach((l) => l(msg));
        } catch {}
      };
      this.socket.onclose = () => {
        setTimeout(() => {
          this.reconnectDelay = Math.min(
            this.maxDelay,
            this.reconnectDelay * 2
          );
          this.connect();
        }, this.reconnectDelay);
      };
      this.socket.onerror = () => {
        // ignore, close triggers reconnect
      };
    } catch (err) {
      setTimeout(() => this.connect(), this.reconnectDelay);
    }
  }

  addListener(fn: (m: Msg) => void) {
    this.listeners.push(fn);
  }
  removeListener(fn: (m: Msg) => void) {
    this.listeners = this.listeners.filter((l) => l !== fn);
  }
}

const host = location.hostname;
const wsUrl = `ws://${host}:8080/ws`; // update if backend uses other host/port
export const wsClient = new WSClient(wsUrl);
