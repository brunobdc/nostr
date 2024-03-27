import { v4 as uuid } from 'uuid'
import '../css/style.css'
import { MessageInputElement, SendMessageButtomElement } from './Element'
import NostrEvent from './NostrEvent'
import Relay, { RelayResponseEvent, RelayResponseOk } from './Relay'

const relay = new Relay("")
const subscription_id = uuid()

relay.addEventListener("EVENT", (ev: CustomEvent<RelayResponseEvent>) => {
    if (ev.detail.subscription_id == subscription_id) {
        // Do something
    }
    // Do something
})

relay.NewSubscription(subscription_id)

relay.addEventListener("OK", (ev: CustomEvent<RelayResponseOk>) => {
    // Do Something
})

SendMessageButtomElement.onclick = function () {
    const message = MessageInputElement.value
    relay.SendEvent(new NostrEvent({ kind: 1, content: message, tags: [] }))
}
