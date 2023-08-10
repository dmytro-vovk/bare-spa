import SimpleRPC from 'simple-jsonrpc-js';
import ReconnectingWebSocket from 'reconnecting-websocket';

type request = {
	id?: number
	jsonrpc: string
	method: string
	params?: any // todo: specified obj type or array with params (second declared in specification)
}

// type response = {
// 	id: number
// 	jsonrpc: string
// 	result?: any
// 	error?: error
// }
//
// type error = {
// 	code: number
// 	message: string
// 	data?: any
// }

export default class RPC {
	private ws: ReconnectingWebSocket;
	private rpc: SimpleRPC;
	private readonly handlers: Record<string, ((data: any) => void)[]>;

	constructor(address: string) {
		this.ws = new ReconnectingWebSocket(address);
		this.rpc = new SimpleRPC();
		this.rpc.toStream = message => this.ws.send(message);
		this.handlers = {};

		this.ws.onerror = error => {
			console.error(error);
			const onlineLabel = document.querySelector('#online-status') as HTMLSpanElement;
			if (onlineLabel) {
				onlineLabel.classList.add('bg-danger');
				onlineLabel.classList.remove('bg-success');
				onlineLabel.innerHTML = '<i class="fas fa-sync-alt fa-spin"></i> Connecting&hellip;';
			}
		};

		this.ws.onopen = () => {
			const onlineLabel = document.querySelector('#online-status') as HTMLSpanElement;
			if (onlineLabel) {
				onlineLabel.classList.add('bg-success');
				onlineLabel.classList.remove('bg-danger');
				onlineLabel.innerText = 'Online';
			}
			for (const method in this.handlers) {
				this.notify('subscribe', method)
			}
		};

		this.ws.onmessage = (event: MessageEvent) => {
			const msg = JSON.parse(event.data as any as string) as request;
			if (msg.id) {
				this.rpc.messageHandler(event.data);
			} else if (this.handlers[msg.method]) {
				this.handlers[msg.method].forEach(h => h(msg.params));
			}
		};
	}

	public async call(method: string, data?: any): Promise<void> {
		return this.rpc.call(method, data);
	}

	public notify(method: string, data: any): void {
		this.rpc.notification(method, data);
	}

	public subscribe(method: string, handler: (data) => void): void {
		if (this.handlers[method]) {
			this.handlers[method].push(handler);
		} else {
			this.handlers[method] = [handler];
		}

		this.notify('subscribe', method)
	}

	public unsubscribe(method: string, handler: (data) => void): void {
		if (this.handlers[method]) {
			const index = this.handlers[method].indexOf(handler);
			if (index !== -1) this.handlers[method].splice(index, 1);
		}

		this.notify('unsubscribe', method)
	}
}
