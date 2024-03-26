import { NostrEvent } from "./NostrEvent";

export { };

declare global {
    interface Window {
        nostr: {
            getPublicKey: () => string
            signEvent: (event: { created_at: number, kind: number, tags: string[][], content: string }) => NostrEvent
            getRelays?: () => { [url: string]: {read: boolean, write: boolean} }
            nip04?: {
                encrypt: (pubkey, plaintext) => string
                decrypt: (pubkey, ciphertext) => string
            }
            nip44?: {
                encrypt: (pubkey, plaintext) => string
                decrypt: (pubkey, ciphertext) => string
            }
        }
    }
}