import App from "../../../system";
import {$$, $html, $onClick, $setSelect, renderOptions} from "../../index";
import {clearForm, ErrorCode, setInvalidFeedback} from "../../../components/errors";

const baudRates = [
	9_600,
	14_400,
	19_200,
	31_250,
	38_400,
	56_000,
	57_600,
	76_800,
	115_200,
	128_000,
	230_400,
	250_000,
	256_000,
];

const dataBits = [
	5,
	6,
	7,
	8,
];

const parity = {
	"None":  "N",
	"Odd":   "O",
	"Even":  "E",
	"Space": "S",
	"Mark":  "M",
};

const stopBits = [
	1,
	2,
];

type Config = {
	enabled:  boolean;
	baudRate: number;
	dataBits: number;
	parity:   string;
	stopBits: number;
}

export default class UART {
	private form: {
		enabled:   HTMLInputElement,
		baudRate:  HTMLSelectElement,
		dataBits:  HTMLSelectElement,
		parity:    HTMLSelectElement,
		stopBits:  HTMLSelectElement,
		apply:     HTMLButtonElement,
		feedbacks: HTMLDivElement[],
	};

	private subscriptions = new Map<string, (data) => void>([
		[`uart-${this.number}.setConfig`, (req: Config) => {
			this.setConfig(req);
			this.app.info(`Настройки UART-${this.number} применены`)
		}],
	]);

	constructor(
		private readonly app: App,
		readonly number: number,
		readonly name: string,
	) {}

	public async render(selector: string) {
		$html(selector,
			`<div class="row">
				<div class="col-sm-12 col-md-12 col-lg-6">
					<div class="card">
						<div class="card-header">
							<h3 class="card-title">Настройки подключения</h3>
							<div class="card-tools m-0">			
								<div class="col-sm-8">
									<div class="custom-control custom-switch">
									  <input id="uart-${this.number}-enabled" class="custom-control-input" type="checkbox">
									  <label for="uart-${this.number}-enabled" class="custom-control-label"></label>
									</div>
								</div>
                    		</div>	
						</div>
						<div class="card-body">
							<div class="form-group row">
								<label for="uart-${this.number}-baud-rate" class="col-sm-4 col-form-label">Baud Rate</label>
								<div class="col-sm-8">
									<select id="uart-${this.number}-baud-rate" class="custom-select">
										${renderOptions(baudRates)}
									</select>
								</div>
							</div>
							<div class="form-group row">
								<label for="uart-${this.number}-data-bits" class="col-sm-4 col-form-label">Data Bits</label>
								<div class="col-sm-8">
									<select id="uart-${this.number}-data-bits" class="custom-select">
										${renderOptions(dataBits)}
									</select>
								</div>
							</div>
							<div class="form-group row">
								<label for="uart-${this.number}-parity" class="col-sm-4 col-form-label">Parity</label>
								<div class="col-sm-8">
									<select id="uart-${this.number}-parity" class="custom-select">
										${renderOptions(parity)}
									</select>
								</div>
							</div>
							<div class="form-group row">
								<label for="uart-${this.number}-stop-bits" class="col-sm-4 col-form-label">Stop Bits</label>
								<div class="col-sm-8">
									<select id="uart-${this.number}-stop-bits" class="custom-select">
										${renderOptions(stopBits)}
									</select>
								</div>
							</div>							
						</div>
						<div class="card-footer">							
							<button id="uart-${this.number}-apply" class="btn btn-success float-right" type="button">Установить</button>
						</div>
					</div>
				</div>				
			</div>`
		);

		this.form = {
			enabled:   $$(`#uart-${this.number}-enabled`) as HTMLInputElement,
			baudRate:  $$(`#uart-${this.number}-baud-rate`) as HTMLSelectElement,
			dataBits:  $$(`#uart-${this.number}-data-bits`) as HTMLSelectElement,
			parity:    $$(`#uart-${this.number}-parity`) as HTMLSelectElement,
			stopBits:  $$(`#uart-${this.number}-stop-bits`) as HTMLSelectElement,
			apply:     $$(`#uart-${this.number}-apply`) as HTMLButtonElement,
			feedbacks: [],
		};

		const config = await this.app.call("uart.getConfig", { id: this.number }) as any as Config;
		this.setConfig(config);
	}

	public visit(): void {
		for (const [method, handler] of this.subscriptions) {
			this.app.subscribe(method, handler);
		}
	}

	public handleEvents(): void {
		$onClick(this.form.apply, () => {
			const config: Config = {
				enabled:  this.form.enabled.checked,
				baudRate: Number(this.form.baudRate.value),
				dataBits: Number(this.form.dataBits.value),
				parity:   this.form.parity.value,
				stopBits: Number(this.form.stopBits.value),
			};

			this.app.call("uart.setConfig", {
				id:     this.number,
				config: config,
			}).catch(err => {
				if (err.code === ErrorCode.ValidationErr && err.data) {
					setInvalidFeedback(this.form, err.data);
					return
				}

				this.app.error(err);
			})
		});
	}

	public setConfig(req: Config): void {
		clearForm(this.form);

		this.form.enabled.checked = req.enabled;
		$setSelect(this.form.baudRate, req.baudRate);
		$setSelect(this.form.dataBits, req.dataBits);
		$setSelect(this.form.parity, req.parity);
		$setSelect(this.form.stopBits, req.stopBits);
	}

	public leave(): void {
		for (const [method, handler] of this.subscriptions) {
			this.app.unsubscribe(method, handler);
		}
	}
}
