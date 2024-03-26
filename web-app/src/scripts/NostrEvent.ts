export interface NostrEventDto {
    id?: string
    pubkey?: string
    created_at?: number
    kind: number
    tags: string[][]
    content: string
    sig?: string
}

export default class NostrEvent {
    id: string
    pubkey: string
    created_at: number
    kind: number
    tags: string[][]
    content: string
    sig: string

    constructor(eventDto: NostrEventDto) {
        if (eventDto.id && eventDto.pubkey && eventDto.created_at && eventDto.sig) {
            this.id = eventDto.id
            this.pubkey = eventDto.pubkey
            this.created_at = eventDto.created_at
            this.kind = eventDto.kind
            this.tags = eventDto.tags
            this.content = eventDto.content
            this.sig = eventDto.sig
            return
        }

        const signedEvent = window.nostr.signEvent({
            created_at: Date.now(),
            kind: eventDto.kind,
            tags: eventDto.tags,
            content: eventDto.content
        })

        this.id = signedEvent.id
        this.pubkey = signedEvent.pubkey
        this.created_at = signedEvent.created_at
        this.kind = signedEvent.kind
        this.tags = signedEvent.tags
        this.content = signedEvent.content
        this.sig = signedEvent.sig
    }
}
