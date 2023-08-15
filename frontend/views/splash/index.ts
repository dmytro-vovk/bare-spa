import App from '../../system'
import {$html} from '../index'
import Login from './login'

export default class Splash {
	private login: Login

	constructor(private readonly app: App) {
		this.login = new Login(app)
	}

	public render(selector: string): void {
		$html(selector, `
			<div class="row">
				<div class="col-lg-6" id="login-container"></div>
			</div>`)

		this.login.render('#login-container')
	}
}
