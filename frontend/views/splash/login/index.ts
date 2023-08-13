import App from '../../../system'
import {$$, $buttonBusy, $buttonReady, $html, $onClick} from '../../index'

export default class Login {
	private submit: HTMLButtonElement
	private name: HTMLInputElement
	private password: HTMLInputElement

	constructor(private readonly app: App) {
	}

	public render(selector: string): void {
		$html(selector, `
			<div class="card card-primary">
				<div class="card-header">
					<h3 class="card-title">Вхід в систему</h3>
				</div>
				<form>
					<div class="card-body">
						<div class="form-group">
							<label for="loginEmail">Ім'я користувача</label>
							<input type="email" class="form-control" id="loginName"
								   placeholder="Введіть своє ім'я користувача">
						</div>
						<div class="form-group">
							<label for="loginPassword">Пароль</label>
							<input type="password" class="form-control" id="loginPassword"
								   placeholder="Введіть свій пароль">
						</div>
					</div>
					<div class="card-footer">
						<button type="submit" class="btn btn-primary" id="loginSubmit">Увійти</button>
					</div>
				</form>
			</div>
			`)

		this.submit = $$('#loginSubmit') as HTMLButtonElement
		this.name = $$('#loginName') as HTMLInputElement
		this.password = $$('#loginPassword') as HTMLInputElement

		$onClick(this.submit, evt => {
			$buttonBusy(this.submit)
			const name = this.name.value
			const password = this.password.value

			let valid = false

			if (name === '') {
				this.name.classList.add('is-invalid')
			} else {
				this.name.classList.remove('is-invalid')
				valid = true
			}

			if (password === '') {
				this.password.classList.add('is-invalid')
				valid = false
			} else {
				this.password.classList.remove('is-invalid')
				valid = valid && true
			}

			if (valid) {
				this.app.call('login', {
					name,
					password
				})
					.then(success => {
						if (success as unknown as boolean) {
							window.setTimeout(() => {location.reload()}, 500)
						}
					})
					.catch(err => {
						this.app.error(err)
						$buttonReady(this.submit)
					})
			} else $buttonReady(this.submit)
			evt.preventDefault()
		})
	}
}
