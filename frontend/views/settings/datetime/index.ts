import {$$, $html, $onClick, range} from "../../index";
import App from "../../../system";
import {NTP} from "../../../types/datetime";

export default class Datetime {
	private systemTimeForm: {
		date:  HTMLInputElement,
		time:  HTMLInputElement,
		label: HTMLLabelElement,
		apply: HTMLButtonElement,
	};

	private ntpServersForm: {
		enabled: HTMLInputElement,
		servers: HTMLInputElement[],
		apply:   HTMLButtonElement,
	};

	private subscriptions = new Map<string, (data) => void>([
		['datetime.set', this.setDatetime.bind(this)],
		['datetime.setNTP', this.loggedSetNTP.bind(this)],
	]);

	constructor(
		private readonly app: App,
		private timestamp: string,
		private ntp: NTP,
	) {}

	public render(selector: string): void {
		$html(selector,
			`<div class="row">
				<div class="col-sm-12 col-md-12 col-lg-6" id="system-time-card"></div>
				<div class="col-sm-12 col-md-12 col-lg-6" id="ntp-servers-card"></div>
			</div>`
		);

		Datetime.renderSystemTime("#system-time-card");
		Datetime.renderNTPServers("#ntp-servers-card");

		this.systemTimeForm = {
			date:  $$('#date-input') as HTMLInputElement,
			time:  $$('#time-input') as HTMLInputElement,
			label: $$('#current-datetime') as HTMLLabelElement,
			apply: $$('#apply-datetime') as HTMLButtonElement,
		};

		this.ntpServersForm = {
			enabled: $$('#is-ntp-enabled') as HTMLInputElement,
			servers: [...document.querySelectorAll(".js-ntp-server") as NodeListOf<HTMLInputElement>],
			apply:   $$('#apply-ntp') as HTMLButtonElement,
		};

		this.setDatetime({ timestamp: this.timestamp });
		this.setNTP({ ntp: this.ntp });
	}

	private static renderSystemTime(selector: string): void {
		$html(selector,
			`<div class="card">
				<div class="card-header">
					<h3 class="card-title">Системное время (UTC)</h3>
				</div>
				<div class="card-body">
					<div class="form-group row">
						<label for="date" class="col-sm-4 col-form-label">Дата</label>
						<div class="col-sm-8">
							<input type="date" id="date-input" class="form-control">
						</div>
					</div>
					<div class="form-group row">
						<label for="time" class="col-sm-4 col-form-label">Время</label>
						<div class="col-sm-8">
							<input type="time" id="time-input" class="form-control">
						</div>
					</div>
				</div>
				<div class="card-footer">
					<label class="text-muted col-sm-6 col-form-label" id="current-datetime"></label>
					<button type="button" id="apply-datetime" class="btn btn-success float-right">Установить</button>
				</div>
			</div>`
		);
	}

	private static renderNTPServers(selector: string): void {
		const count = 3;
		$html(selector,
			`<div class="card">
				<div class="card-header">
					<h3 class="card-title">Сервера точного времени</h3>
				</div>
				<div class="card-body">
					<div class="form-group row">
						<div class="offset-sm-4 col-sm-4">
							<div class="form-check">
								<input type="checkbox" class="form-check-input" id="is-ntp-enabled">
								<label class="form-check-label" for="is-ntp-enabled">Использовать</label>
							</div>
						</div>
					</div>
					${range(count).map(idx => `
					<div class="form-group row">
						<label for="ntp-server-${idx}" class="col-sm-4 col-form-label">Сервер ${idx}</label>
						<div class="col-sm-8">
							<input type="text" class="form-control js-ntp-server" id="ntp-server-${idx}">
						</div>
					</div>`).join('')}
				</div>
				<div class="card-footer">
					<button type="button" id="apply-ntp" class="btn btn-success float-right">Применить</button>
				</div>
			</div>`
		);
	}

	public visit(): void {
		for (const [method, handler] of this.subscriptions) {
			this.app.subscribe(method, handler);
		}
	}

	public leave(): void {
		for (const [method, handler] of this.subscriptions) {
			this.app.unsubscribe(method, handler);
		}
	}

	public handleEvents(): void {
		$onClick(this.systemTimeForm.apply, () => {
			this.callSetDatetime({
				timestamp: this.systemTimeForm.date.value + " " + this.systemTimeForm.time.value + ":00",
			})
				.then(() => this.app.info('Системное время обновлено'))
				.catch(err => this.app.error(err));
		});

		$onClick(this.ntpServersForm.apply, () => {
			this.callSetNTP({
				ntp: {
					enabled: this.ntpServersForm.enabled.checked,
					servers: this.ntpServersForm.servers.map(inp => inp.value).filter(inp => inp),
				},
			}).catch(err => this.app.error(err));
		});
	}

	private async callSetDatetime(req: { timestamp: string }): Promise<void> {
		return this.app.call('datetime.set', req)
	}

	private setDatetime(req: { timestamp: string }): void {
		this.timestamp = req.timestamp;
		this.systemTimeForm.label.innerText = req.timestamp;
	}

	private async callSetNTP(req: { ntp: NTP }): Promise<void> {
		return this.app.call('datetime.setNTP', req);
	}

	private loggedSetNTP(req: { ntp: NTP }): void {
		this.setNTP(req)
		this.app.info('NTP сервера обновлены');
	}

	private setNTP(req: { ntp: NTP }) {
		this.ntp = req.ntp;
		this.ntpServersForm.servers.forEach((inp, idx) => inp.value = req.ntp.servers[idx] ?? "");

		const disabled = req.ntp.enabled;
		this.ntpServersForm.enabled.checked = disabled;
		this.systemTimeForm.date.disabled = disabled;
		this.systemTimeForm.time.disabled = disabled;
		this.systemTimeForm.apply.disabled = disabled;
	}
}
