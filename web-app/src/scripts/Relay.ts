import Filter from "./Filter";
import NostrEvent from "./NostrEvent";

export default class Relay extends EventTarget {
    private websocket: WebSocket

    constructor(url: string) {
        super()
        this.websocket = new WebSocket(url)
        this.websocket.onmessage = (event) => {
            const msg_array: Array<any> = JSON.parse(event.data)
            switch (msg_array[0]) {
                case "EVENT":
                    this.dispatchEvent(new CustomEvent("EVENT", { detail: { data: new NostrEvent(msg_array[1]) }}))
                    break;
                case "OK":
                    this.dispatchEvent(new CustomEvent("OK", { detail: { data: { event_id: msg_array[1], success: msg_array[2], message: msg_array[3] } }}))
                    break;
                case "EOSE":
                    this.dispatchEvent(new CustomEvent("EOSE", { detail: { subscription_id: msg_array[1] }}))
                    break;
                case "CLOSED":
                    this.dispatchEvent(new CustomEvent("CLOSED", { detail: { data: { subscription_id: msg_array[1], message: msg_array[2] } }}))
                    break;
                case "NOTICE":
                    this.dispatchEvent(new CustomEvent("NOTICE", { detail: { message: msg_array[1] }}))
                    break;
            }
        }
    }

    SendEvent(event: NostrEvent) {
        this.websocket.send(JSON.stringify(["EVENT", event]))
    }

    NewSubscription(subscriptionID: string, ...filters: Filter[]) {
        this.websocket.send(JSON.stringify(["REQ", subscriptionID, ...filters]))
    }

    CloseSubscription(subscriptionID: string) {
        this.websocket.send(JSON.stringify(["CLOSE", subscriptionID]))
    }
}
