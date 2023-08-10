import RPC from './rpc';
import * as toastr from 'toastr';
import Navigo from 'navigo';

const CallSystemLog = 'system.logRead';
const CallSystemReset = 'system.reset';

export default class App {
	private readonly rpc: RPC;
	private router: Navigo;
	private readonly msg: typeof toastr;

	constructor(rpc: RPC) {
		this.rpc = rpc;
		this.msg = toastr;
		this.msg.options.positionClass = 'toast-bottom-right';
	}

	public setRouter(r: Navigo): void {
		this.router = r;
		this.router.resolve();
	}

	// Call this to go to another page properly
	public navigate(path: string): void {
		this.router.navigate(path);
	}

	// Call this when page done updating the layout
	public updatePage(): void {
		this.router.updatePageLinks();
	}

	// Make RPC call
	public async call(method: string, data?: any): Promise<void> {
		return this.rpc.call(method, data);
	}

	public async systemReset(): Promise<void> {
		return this.rpc.call(CallSystemReset);
	}

	public async systemLog(): Promise<void> {
		return this.rpc.call(CallSystemLog);
	}

	public subscribe(method: string, handler: (data) => void): void {
		this.rpc.subscribe(method, handler);
	}

	public unsubscribe(method: string, handler: (data) => void): void {
		this.rpc.unsubscribe(method, handler);
	}

	// Show error message
	public error(err: any): void {
		console.error(err);
		this.msg.error(err.message);
	}

	// Show information message
	public info(message: string): void {
		this.msg.info(message);
	}

	// Update sidebar links display
	public sideBarToggle(href: string): void {
		document.querySelectorAll('.nav-sidebar a.nav-link').forEach(e => e.classList.remove('active'));

		const link = document.querySelector(`.nav-sidebar a.nav-link[href="${href}"]`);
		if (link) {
			link.classList.add('active');
		}
	}
}
