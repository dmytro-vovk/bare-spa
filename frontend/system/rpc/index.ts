import SimpleRPC from 'simple-jsonrpc-js'
import ReconnectingWebSocket from 'reconnecting-websocket'

interface request {
	id?: number
	jsonrpc: string
	method: string
	params?: unknown
}

export default class RPC {
	private readonly ws: ReconnectingWebSocket
	private readonly rpc: SimpleRPC
	private readonly handlers: Record<string, Array<(data: unknown) => void>>
	private reconnected = false

	constructor(address: string) {
		this.ws = new ReconnectingWebSocket(address)
		this.rpc = new SimpleRPC()
		this.rpc.toStream = message => {
			this.ws.send(message)
		}
		this.handlers = {}
		this.ws.onerror = () => {
			this.reconnected = true
			const onlineLabel = document.querySelector('#online-status')
			if (onlineLabel != null) {
				onlineLabel.classList.add('bg-danger')
				onlineLabel.classList.remove('bg-success')
				onlineLabel.innerHTML = '<i class="fas fa-sync-alt fa-spin"></i> Підключення&hellip;'
			}
		}
		this.ws.onopen = () => {
			const onlineLabel = document.querySelector('#online-status') as HTMLSpanElement
			if (onlineLabel != null) {
				onlineLabel.classList.add('bg-success')
				onlineLabel.classList.remove('bg-danger')
				onlineLabel.innerText = 'На зв\'язку'
			}
			if (this.reconnected) {
				for (const method in this.handlers) {
					this.notify('subscribe', method)
				}
			}
		}
		this.ws.onmessage = (event: MessageEvent) => {
			const msg = JSON.parse(event.data as string) as request
			if (msg.id) {
				this.rpc.messageHandler(event.data)
			} else if (msg.method in this.handlers) {
				this.handlers[msg.method].forEach(h => {
					h(msg.params)
				})
			}
		}
	}

	public async call(method: string, data?: unknown): Promise<void> {
		return this.rpc.call(method, data)
	}

	public notify(method: string, data: unknown): void {
		this.rpc.notification(method, data)
	}

	public subscribe(method: string, handler: (data) => void): void {
		if (method in this.handlers) {
			this.handlers[method].push(handler)
		} else {
			this.handlers[method] = [handler]
		}
		this.notify('subscribe', method)
	}

	public unsubscribe(method: string): void {
		delete this.handlers[method]
		this.notify('unsubscribe', method)
	}
}
