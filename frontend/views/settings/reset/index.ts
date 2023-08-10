import App from "../../../system";
import {$$, $html, $onClick} from "../../index";

export default class Reset {
	private resetButton: HTMLButtonElement;

	constructor(private readonly app: App) {
	}

	public render(selector: string): void {
		$html(selector, `
			<button type="button" class="btn btn-danger" data-toggle="modal" data-target="#reset-modal">Сброс настроек системы</button>
			<div class="modal fade" id="reset-modal" aria-modal="true" role="dialog">
				<div class="modal-dialog">
					<div class="modal-content bg-danger">
						<div class="modal-header">
							<h4 class="modal-title">Сброс настроек</h4>
							<button type="button" class="close" data-dismiss="modal" aria-label="закрыть">
								<span aria-hidden="true">×</span>
							</button>
						</div>
						<div class="modal-body">
							<p>Сбросить все настройки системной платы?</p>
						</div>
						<div class="modal-footer justify-content-between">
							<button type="button" class="btn btn-outline-light" data-dismiss="modal">Отмена</button>
							<button type="button" class="btn btn-outline-light" id="system-reset">Сбросить</button>
						</div>
					</div>
				</div>
			</div>`);

		this.resetButton = $$('#system-reset') as HTMLButtonElement;

		$onClick(this.resetButton, () => {
			this.app.systemReset().then(
				() => ($$('[data-dismiss=modal]') as HTMLButtonElement).click(),
				error => {
					this.app.error(error);
					($$('[data-dismiss=modal]') as HTMLButtonElement).click();
				},
			);
		});
	}
}
