import Filter from "./Filter";
import NostrEvent from "./NostrEvent";

interface RelayEventMap {
    EVENT: CustomEvent<RelayResponseEvent>
    OK: CustomEvent<RelayResponseOk>
    EOSE: CustomEvent<RelayResponseEose>
    CLOSED: CustomEvent<RelayResponseClosed>
    NOTICE: CustomEvent<RelayResponseNotice>
}

interface RelayEventTarget extends EventTarget {
    addEventListener<K extends keyof RelayEventMap>(type: K, listener: (this: Relay, ev: RelayEventMap[K]) => any, options?: boolean | AddEventListenerOptions): void;
    addEventListener(type: string, listener: EventListenerOrEventListenerObject | null, options?: AddEventListenerOptions | boolean): void;
}

const RelayEventTarget = EventTarget as {
    new(): RelayEventTarget
    prototype: RelayEventTarget
}

export default class Relay extends RelayEventTarget {
    private websocket: WebSocket

    constructor(url: string) {
        super()
        this.websocket = new WebSocket(url)
        this.websocket.onmessage = (event) => {
            const msg_array: Array<any> = JSON.parse(event.data)
            switch (msg_array[0]) {
                case "EVENT":
                    this.dispatchEvent(new CustomEvent("EVENT", { detail: <RelayResponseEvent>{ subscription_id: msg_array[1], event: new NostrEvent(msg_array[2]) } }))
                    break;
                case "OK":
                    this.dispatchEvent(new CustomEvent("OK", { detail: <RelayResponseOk>{ event_id: msg_array[1], success: msg_array[2], message: msg_array[3] } }))
                    break;
                case "EOSE":
                    this.dispatchEvent(new CustomEvent("EOSE", { detail: <RelayResponseEose>{ subscription_id: msg_array[1] }}))
                    break;
                case "CLOSED":
                    this.dispatchEvent(new CustomEvent("CLOSED", { detail: <RelayResponseClosed>{ subscription_id: msg_array[1], message: msg_array[2] }}))
                    break;
                case "NOTICE":
                    this.dispatchEvent(new CustomEvent("NOTICE", { detail: <RelayResponseNotice>{ message: msg_array[1] }}))
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

export interface RelayResponseEvent {
    subscription_id: string
    event: NostrEvent
}

export interface RelayResponseOk {
    event_id: string
    success: boolean
    message: string
}

export interface RelayResponseEose {
    subscription_id: string
}

export interface RelayResponseClosed {
    subscription_id: string
    message: string
}

export interface RelayResponseNotice {
    message: string
}
