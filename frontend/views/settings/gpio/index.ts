import App from "../../../system";
import {GPIO} from "../types";
import {$$, $html, $onChange} from "../../index";
import {markChecked, markSelected} from "../../../components";

const input = 'input';
const output = 'output';

const high = 'Высокий';
const low = 'Низкий';

export default class GPIOTable {
	private gpioTableBody: HTMLTableSectionElement;
	private subscriptions = new Map<string, (data) => void>([
		["gpio.setDirection", this.setDirection.bind(this)],
		["gpio.setInversion", this.setInversion.bind(this)],
		["gpio.setLevel", this.setLevel.bind(this)],
	]);

	constructor(
		private readonly app: App,
		private gpio: Record<number, GPIO>
	) {}

	public render(selector: string): void {
		$html(selector,
			`<table class="table table-bordered table-striped" style="width: auto">
				<thead>
					<tr>
						<th>Пин</th>
						<th>Название</th>
						<th>Уровень</th>
						<th>Направление</th>
						<th>Исходный уровень</th>
						<th>Инверсия</th>
					</tr>
				</thead>
				<tbody id="gpio-table-body"></tbody>
			</table>`
		);

		this.gpioTableBody = $$('#gpio-table tbody') as HTMLTableSectionElement;
		this.renderTableBody();
	}

	private renderTableBody(): void {
		$html(this.gpioTableBody, Object
			.keys(this.gpio)
			.map(pin => GPIOTable.renderRow(this.gpio[pin]))
		);
	}

	private static renderRow(io: GPIO): string {
		return (
			`<tr>
				<td class="text-right">${io.pin}</td>
				<td>${io.name}</td>
				<td>
					<input type="checkbox" id="level-${io.pin}" ${markChecked(io.level)}>
					<label for="level-${io.pin}">${io.level ? high : low}</label>
				</td>
				<td>
					<select class="custom-select" id="direction-${io.pin}">
						<option value="${output}" ${markSelected(io.direction)}>Выход</option>
						<option value="${input}" ${markSelected(!io.direction)}>Вход</option>
					</select>
				</td>
				<td class="text-center"> ${io.default_state
				? `<i class="fas fa-level-up-alt   mr-1"></i>${high}`
				: `<i class="fas fa-level-down-alt mr-1"></i>${low}`}
				</td>
				<td>
					<input type="checkbox" id="inversion-${io.pin}" ${markChecked(io.inversion)}>
				</td>
			</tr>`
		);
	}

	public visit(): void {
		for (const [method, handler] of this.subscriptions) {
			this.app.subscribe(method, handler);
		}
	}

	public handleEvents(): void {
		$onChange(this.gpioTableBody, (event) => {
			const target = event.target as HTMLElement;
			const [typ, pin] = target.id.split('-');

			switch (typ) {
				case 'level':
					return this.changeLevel(+pin, target as HTMLInputElement);
				case 'direction':
					return this.changeDirection(+pin, target as HTMLSelectElement);
				case 'inversion':
					return this.changeInversion(+pin, target as HTMLInputElement);
			}
		});
	}

	private changeDirection(pin: number, el: HTMLSelectElement): void {
		this.app.call('gpio.setDirection', {
			pin: pin,
			direction: el.value === 'output'
		}).catch(error => this.app.error(error));
	}

	private setDirection(req: { pin: number, direction: boolean }): void {
		const direction = this.gpioTableBody.querySelector(`#direction-${req.pin}`) as HTMLSelectElement;

		// use it as written in documentation (see it on backend)
		const setDirection = req.direction ? output : input;
		direction.selectedIndex = Array.from(direction.options).findIndex(opt => opt.value === setDirection);
	}

	private changeInversion(pin: number, el: HTMLInputElement): void {
		this.app.call('gpio.setInversion', {
			pin: pin,
			inversion: el.checked
		}).catch(error => this.app.error(error));
	}

	private setInversion(req: { pin: number, inversion: boolean }): void {
		const inversion = this.gpioTableBody.querySelector(`#inversion-${req.pin}`) as HTMLInputElement;

		inversion.checked = req.inversion;
	}

	private changeLevel(pin: number, el: HTMLInputElement): void {
		this.app.call('gpio.setLevel', {
			pin: pin,
			level: el.checked
		}).catch(error => this.app.error(error));
	}

	private setLevel(req: { pin: number, level: boolean }): void {
		const level = this.gpioTableBody.querySelector(`#level-${req.pin}`) as HTMLInputElement;
		const label = this.gpioTableBody.querySelector(`[for=${level.id}]`) as HTMLLabelElement;

		level.checked = req.level;
		label.textContent = level.checked ? high : low;
	}

	public leave(): void {
		for (const [method, handler] of this.subscriptions) {
			this.app.unsubscribe(method, handler);
		}
	}
}
