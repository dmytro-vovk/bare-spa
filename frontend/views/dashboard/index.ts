import App from '../../system';
import {$html} from "../index";

export default class Dashboard {

	private subscriptions = new Map<string, (data) => void>([
		["system.setTemperatures", this.setTemperatures.bind(this)],
	]);

	constructor(private readonly app: App) {}

	public async render(selector: string) {
		$html(selector, `<div class="row" id="thermometers"></div>`);

		await this.launch();
	}

	private async launch(): Promise<void> {
		try {
			await this.handleDynamicContent();
		} catch (err) {
			this.app.error(err);
		}
	}

	private async handleDynamicContent() {
		await this.renderDynamicContent();
		this.visit();
	}

	private async renderDynamicContent() {
		const thermometers = await this.app.call("system.getTemperatures") as any as Record<string, number>;

		$html("#thermometers", Object.keys(thermometers).map(name => this.renderThermometer(name)).join(""));

		this.setTemperatures(thermometers);
	}

	public visit(): void {
		for (const [method, handler] of this.subscriptions) {
			this.app.subscribe(method, handler);
		}
	}

	private renderThermometer(name: string): string {
		return `
			<div class="col-6">
				<div class="info-box">
					<span class="info-box-icon bg-gradient-green elevation-1">
						<i class="fas fa-thermometer-half"></i>
					</span>
					<div class="info-box-content">
						<span class="info-box-text">${name}</span>
						<span class="info-box-number js-thermometer" data-name="${name}"></span>
					</div>
				</div>
			</div>`
	}

	private setTemperatures(temperatures: Record<string, number>): void {
		const thermometers = document.querySelectorAll(".js-thermometer") as NodeListOf<HTMLSpanElement>;

		thermometers.forEach(therm => {
			const temp = Object.keys(temperatures).find(temp => temp === therm.dataset.name);
			therm.innerText = (String(temperatures[temp]) ?? "-") + " °";
		});
	}

	public leave(): void {
		for (const [method, handler] of this.subscriptions) {
			this.app.unsubscribe(method, handler);
		}
	}
}
