import RPC from './rpc'
import * as toastr from 'toastr'
import Navigo, {Match} from 'navigo'
import {$onClick, $replace} from '../views'

export interface Renderable {
	render: (selector: string, match?: Match) => void
	leave?: () => void
}

export default class App {
	private readonly rpc: RPC
	private readonly rpcPrefix: string
	private readonly contentSelector
	private menuHTML
	private menuSelector
	private router: Navigo
	private readonly msg: typeof toastr

	constructor(rpc: RPC, rpcNS: string, contentSelector: string) {
		this.rpc = rpc
		this.rpcPrefix = rpcNS + '.'
		this.contentSelector = contentSelector
		this.msg = toastr
		this.msg.options.positionClass = 'toast-bottom-right'
		this.msg.options.progressBar = true
		this.router = new Navigo('/')
		this.router.notFound(() => {
			window.location.href = '/'
		})
	}

	public setMenu(selector: string, html: string): App {
		this.menuSelector = selector
		this.menuHTML = html
		return this
	}

	// Call this to go to another page properly
	public navigate(path: string): void {
		this.router.navigate(path)
	}

	// Make RPC call
	public async call(method: string, data?: unknown): Promise<void> {
		return this.rpc.call(this.rpcPrefix + method, data)
	}

	public subscribe(method: string, handler: (data) => void): void {
		this.rpc.subscribe(method, handler)
	}

	public unsubscribe(method: string): void {
		this.rpc.unsubscribe(method)
	}

	// Show error message
	public error(err: unknown): void {
		console.error(err)
		this.msg.error(err)
	}

	// Show information message
	public info(message: string): void {
		this.msg.info(message)
	}

	// Show success message
	public success(message: string): void {
		this.msg.success(message)
	}

	// Update sidebar links display
	public navToggle(href: string): void {
		document.querySelectorAll('.navbar-nav a.nav-link')
			.forEach(e => {
				e.classList.remove('active')
			})

		const link = document.querySelector(`.navbar-nav a.nav-link[href="${href}"]`)
		if (link != null) {
			link.classList.add('active')
		}
	}

	public route(path: string, handler: Renderable): App {
		this.router.on(path, (match: Match) => {
			this.renderMenu(path)
			handler.render(this.contentSelector, match)
		}, {
			leave: (done: () => void) => {
				handler.leave?.()
				done()
			}
		})
		this.router.resolve()
		return this
	}

	private renderMenu(activeLink: string): void {
		$replace(this.menuSelector, this.menuHTML)
		this.navToggle(activeLink)
		$onClick('#logout', evt => {
			evt.preventDefault()
			this.call('logout')
				.then(() => {
					window.location.reload()
				}, err => {
					this.error(err)
				})
		})
		$onClick('#sync-all', () => {
			this.call('admin.syncSystem').then(() => {
				this.info('Завершено')
			}, err => {
				this.error(err)
			})
		})
	}
}
