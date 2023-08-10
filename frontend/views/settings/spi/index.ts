import App from "../../../system";
import {$$, $html, $onClick, $setSelect, renderOptions} from "../../index";
import {clearForm, ErrorCode, setInvalidFeedback} from "../../../components/errors";

const baudRates = {
	"31.25 kBit/s": 31_250,
	"62.5 kBit/s":  62_500,
	"128 kBit/s":   128_000,
	"256 kBit/s":   256_000,
	"512 kBit/s":   512_000,
	"2 mBit/s":     2_000_000,
	"4 mBit/s":     4_000_000,
};

const clockPolarities = [
	"Low",
	"High",
];

const clockPhases = [
	"1 Edge",
	"2 Edge",
];

const dataBits = [
	8,
	16,
];

const firstBit = [
	"MSB",
	"LSB",
];

type Config = {
	enabled:       boolean;
	baudRate:      number;
	clockPolarity: string;
	clockPhase:    string;
	dataBits:      number;
	firstBit:      string;
}

export default class SPI {
	private form: {
		enabled:       HTMLInputElement,
		baudRate:      HTMLSelectElement,
		clockPolarity: HTMLSelectElement,
		clockPhase:    HTMLSelectElement,
		dataBits:      HTMLSelectElement,
		firstBit:      HTMLSelectElement,
		apply:         HTMLButtonElement,
		feedbacks:     HTMLDivElement[],
	};

	private subscriptions = new Map<string, (data) => void>([
		[`spi-${this.number}.setConfig`, (req: Config) => {
			this.setConfig(req);
			this.app.info(`Настройки SPI-${this.number} применены`)
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
									  <input id="spi-${this.number}-enabled" class="custom-control-input" type="checkbox">
									  <label for="spi-${this.number}-enabled" class="custom-control-label"></label>
									</div>
								</div>
                    		</div>							
						</div>
						<div class="card-body">							
							<div class="form-group row">
								<label for="spi-${this.number}-baud-rate" class="col-sm-4 col-form-label">Baud Rate</label>
								<div class="col-sm-8">
									<select id="spi-${this.number}-baud-rate" class="custom-select">
										${renderOptions(baudRates)}
									</select>
								</div>
							</div>
							<div class="form-group row">
								<label for="spi-${this.number}-clock-polarity" class="col-sm-4 col-form-label">Clock Polarity (CPOL)</label>
								<div class="col-sm-8">
									<select id="spi-${this.number}-clock-polarity" class="custom-select">
										${renderOptions(clockPolarities)}
									</select>
								</div>
							</div>
							<div class="form-group row">
								<label for="spi-${this.number}-clock-phase" class="col-sm-4 col-form-label">Clock Phase (CPHA)</label>
								<div class="col-sm-8">
									<select id="spi-${this.number}-clock-phase" class="custom-select">
										${renderOptions(clockPhases)}
									</select>
								</div>
							</div>
							<div class="form-group row">
								<label for="spi-${this.number}-data-bits" class="col-sm-4 col-form-label">Data Bits (DFF)</label>
								<div class="col-sm-8">
									<select id="spi-${this.number}-data-bits" class="custom-select">
										${renderOptions(dataBits)}
									</select>
								</div>
							</div>
							<div class="form-group row">
								<label for="spi-${this.number}-first-bit" class="col-sm-4 col-form-label">First Bit</label>
								<div class="col-sm-8">
									<select id="spi-${this.number}-first-bit" class="custom-select">
										${renderOptions(firstBit)}
									</select>
								</div>
							</div>
						</div>
						<div class="card-footer">							
							<button id="spi-${this.number}-apply" class="btn btn-success float-right" type="button">Установить</button>
						</div>
					</div>
				</div>				
			</div>`
		);

		this.form = {
			enabled:       $$(`#spi-${this.number}-enabled`) as HTMLInputElement,
			baudRate:      $$(`#spi-${this.number}-baud-rate`) as HTMLSelectElement,
			clockPolarity: $$(`#spi-${this.number}-clock-polarity`) as HTMLSelectElement,
			clockPhase:    $$(`#spi-${this.number}-clock-phase`) as HTMLSelectElement,
			dataBits:      $$(`#spi-${this.number}-data-bits`) as HTMLSelectElement,
			firstBit:      $$(`#spi-${this.number}-first-bit`) as HTMLSelectElement,
			apply:         $$(`#spi-${this.number}-apply`) as HTMLButtonElement,
			feedbacks:     [],
		};

		const config = await this.app.call("spi.getConfig", { id: this.number }) as any as Config;
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
				enabled:       this.form.enabled.checked,
				baudRate:      Number(this.form.baudRate.value),
				clockPolarity: this.form.clockPolarity.value,
				clockPhase:    this.form.clockPhase.value,
				dataBits:      Number(this.form.dataBits.value),
				firstBit:      this.form.firstBit.value,
			};

			this.app.call("spi.setConfig", {
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
		$setSelect(this.form.clockPolarity, req.clockPolarity);
		$setSelect(this.form.clockPhase, req.clockPhase);
		$setSelect(this.form.dataBits, req.dataBits);
		$setSelect(this.form.firstBit, req.firstBit);
	}

	public leave(): void {
		for (const [method, handler] of this.subscriptions) {
			this.app.unsubscribe(method, handler);
		}
	}
}
