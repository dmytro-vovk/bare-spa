import {$html} from '../index'
import {Match} from 'navigo'

export interface Renderable {
	render: (selector: string, match?: Match) => void
	leave?: () => void
}

export default class Page {
	constructor(
		private readonly title: string,
		private readonly content: Renderable
	) {
	}

	public render(selector: string, match: Match): void {
		$html(selector, `
			<div class="content-header">
				<div class="container">
					<div class="row mb-2">
						<div class="col-sm-6">
							<h1 class="m-0">${this.title}</h1>
						</div>
					</div>
				</div>
			</div>
			<div class="content">
				<div class="container" id="content"></div>
			</div>`)
		this.content.render('#content', match)
	}

	public leave() {
		this.content.leave?.()
	}
}
