import App from '../../system'
import {$html} from '../index'
import Login from './login'
import {Match} from 'navigo'

export default class Splash {
	private login: Login

	constructor(private readonly app: App) {
	}

	public render(selector: string, match?: Match): void {
		$html(selector, `
			<div class="row">
				<div class="col-lg-6" id="login-container"></div>
			</div>`)

		this.login = new Login(this.app)
		this.login.render('#login-container')

		if (match && match.data) {
			this.login.resetPassword(match)
		}
	}
}
