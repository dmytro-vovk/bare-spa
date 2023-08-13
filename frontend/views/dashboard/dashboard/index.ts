import App from '../../../system'
import {$html} from '../../index'

export default class DashboardPanel {
	constructor(private readonly app: App) {
	}

	public render(selector: string): void {
		$html(selector, `
			<div class="row">
				<div class="col-12">
					<div class="card">
						<div class="card-header">
							<h3 class="card-title">Якась інфа</h3>
						</div>
						<div class="card-body">
							Щось написано.
						</div>
					</div>
				</div>
			</div>
			`)
	}
}
