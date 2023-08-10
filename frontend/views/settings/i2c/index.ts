import App from "../../../system";
import {$$, $html, $onClick, $setSelect, renderOptions} from "../../index";
import {clearForm, ErrorCode, setInvalidFeedback} from "../../../components/errors";

const baudRates = {
	"SM 10 kBit/s":   10_000,
	"SM 100 kBit/s":  100_000,
	"FM 400 kBit/s": 400_000,
	"FM 1 mBit/s":   1_000_000,
	"FM 3.4 mBit/s": 3_400_000,
};

const addressing = [
	7,
	16,
];

type Config = {
	enabled:    boolean;
	baudRate:   number;
	addressing: number;
}

export default class I2C {
	private form: {
		enabled:    HTMLInputElement,
		baudRate:   HTMLSelectElement,
		addressing: HTMLSelectElement,
		apply:      HTMLButtonElement,
		feedbacks:  HTMLDivElement[],
	};

	private subscriptions = new Map<string, (data) => void>([
		[`i2c-${this.number}.setConfig`, (req: Config) => {
			this.setConfig(req);
			this.app.info(`Настройки I2C-${this.number} применены`)
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
									  <input id="i2c-${this.number}-enabled" class="custom-control-input" type="checkbox">
									  <label for="i2c-${this.number}-enabled" class="custom-control-label"></label>
									</div>
								</div>
                    		</div>							
						</div>
						<div class="card-body">							
							<div class="form-group row">
								<label for="i2c-${this.number}-baud-rate" class="col-sm-4 col-form-label">Baud Rate</label>
								<div class="col-sm-8">
									<select id="i2c-${this.number}-baud-rate" class="custom-select">
										${renderOptions(baudRates)}
									</select>
								</div>
							</div>
							<div class="form-group row">
								<label for="i2c-${this.number}-addressing" class="col-sm-4 col-form-label">Адресация (bit)</label>
								<div class="col-sm-8">
									<select id="i2c-${this.number}-addressing" class="custom-select">
										${renderOptions(addressing)}
									</select>
								</div>
							</div>														
						</div>
						<div class="card-footer">							
							<button id="i2c-${this.number}-apply" class="btn btn-success float-right" type="button">Установить</button>
						</div>
					</div>
				</div>	
			</div>`
		);

		this.form = {
			enabled:    $$(`#i2c-${this.number}-enabled`) as HTMLInputElement,
			baudRate:   $$(`#i2c-${this.number}-baud-rate`) as HTMLSelectElement,
			addressing: $$(`#i2c-${this.number}-addressing`) as HTMLSelectElement,
			apply:      $$(`#i2c-${this.number}-apply`) as HTMLButtonElement,
			feedbacks:  [],
		};

		const config = await this.app.call("i2c.getConfig", { id: this.number }) as any as Config;
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
				enabled:    this.form.enabled.checked,
				baudRate:   Number(this.form.baudRate.value),
				addressing: Number(this.form.addressing.value),
			};

			this.app.call("i2c.setConfig", {
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
		$setSelect(this.form.addressing, req.addressing);
	}

	public leave(): void {
		for (const [method, handler] of this.subscriptions) {
			this.app.unsubscribe(method, handler);
		}
	}
}
