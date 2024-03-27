import '../css/style.css'
import { MessageInputElement, SendMessageButtomElement } from './Element'
import NostrEvent from './NostrEvent'
import Relay, { RelayResponseOk } from './Relay'

const relay = new Relay("")

relay.addEventListener("OK", (ev: CustomEvent<RelayResponseOk>) => {
    // Do Something
})

SendMessageButtomElement.onclick = function () {
    const message = MessageInputElement.value
    relay.SendEvent(new NostrEvent({ kind: 1, content: message, tags: [] }))
}
