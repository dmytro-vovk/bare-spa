import App from '../../system'
import {$html} from '../index'
import DashboardPanel from './dashboard'

export default class Dashboard {
	private readonly dash: DashboardPanel

	constructor(private readonly app: App) {
		this.dash = new DashboardPanel(app)
	}

	public render(selector: string): void {
		$html(selector, `
			<div class="row">
				<div class="col-12" id="dashboard"></div>
			</div>`)

		this.dash.render('#dashboard')
	}
}
