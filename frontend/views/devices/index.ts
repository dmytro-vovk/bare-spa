import App from "../../system";
import {$$, $html} from "../index";
import {Device, DeviceType} from "./types";
import CreateDeviceModal from "./createDeviceModal";
import DeleteDeviceModal from "./deleteDeviceModal";
import EditDeviceModal from "./editDeviceModal";
import CreateDeviceModuleModal from "./createDeviceModuleModal";
import Pagination, {defaultCurrentPage} from "../../components/pagination";
import SelectElement from "../../core/select/index";
import EditDeviceModuleModal from "./editDeviceModuleModal";
import DeleteDeviceModuleModal from "./deleteDeviceModuleModal";

export default class DevicesTable {

	private createDeviceModal: CreateDeviceModal
	private editDeviceModal: EditDeviceModal
	private deleteDeviceModal: DeleteDeviceModal
	private createDeviceModuleModal: CreateDeviceModuleModal
	private editDeviceModuleModal: EditDeviceModuleModal
	private deleteDeviceModuleModal: DeleteDeviceModuleModal

	private showLimitEl: SelectElement
	private tableBodyEl: HTMLTableSectionElement
	private pagination: Pagination

	private devices: Array<Device>

	private subscriptions = new Map<string, (data) => void>([
		["devices.create", async () => {
			await this.load(this.currentPage);
			this.app.info("Добавлено новое устройство");
		}],
		["devices.update", async () => {
			await this.load(this.currentPage);
			this.app.info("Настройки устройства изменены");
		}],
		["devices.delete", async () => {
			await this.load(this.currentPage);
			this.app.info("Устройство успешно удалено")
		}],
		["devices.createModule", async () => {
			await this.load(this.currentPage);
			this.app.info("Добавлен новый модуль");
		}],
		["devices.updateModule", async () => {
			await this.load(this.currentPage);
			this.app.info("Настройки модуля изменены");
		}],
		["devices.deleteModule", async () => {
			await this.load(this.currentPage);
			this.app.info("Модуль успешно удален")
		}],
	]);

	constructor(private readonly app: App) {}

	public async render(selector: string): Promise<void> {
		$html(selector,
			`<div class="card">
				<div class="card-header">
					<div class="row">
						<div class="col d-flex">
							<label class="d-flex align-items-center mb-0" for="devices-show-limit">
								Показывать <select id="devices-show-limit" class="custom-select ml-2 mr-2"></select> устройств
							</label>
						</div>
						<div class="col d-flex justify-content-end">
							<button type="button" class="btn btn-sm btn-default" data-toggle="modal" data-target="#create-device-modal">
								<i class="fas fa-plus"></i> Добавить устройство
							</button>
						</div>
					</div>
				</div>
				<div class="card-body">
					<div class="row p-0">
						<table id="devices-table" class="table table-last-underlined">
							<thead class="thead-light">
								<tr>
									<th rowspan="1" colspan="1" style="width: 5%;"></th>
									<th rowspan="1" colspan="1" style="width: 25%;">Название</th>
									<th rowspan="1" colspan="1" style="width: 20%;">Тип</th>
									<th rowspan="1" colspan="1" style="width: 20%;">Интерфейс</th>
									<th rowspan="1" colspan="1" style="width: 15%;">Адрес</th>
									<th rowspan="1" colspan="1" style="width: 15%;">Действия</th>
								</tr>
							</thead>
							<tbody></tbody>
						</table>
					</div>
				</div>
				<div class="card-footer">
					<div class="row align-items-center">
						<div class="col">
							<p id="shown-devices-info" class="mb-0"></p>
						</div>
						<div class="col">
							<nav id="devices-pagination" aria-label="Available devices navigation"></nav>
						</div>
					</div>
				</div>
			</div>
			<div id="create-device-modal" class="modal fade" tabindex="0"></div>
			<div id="edit-device-modal" class="modal fade" tabindex="0"></div>
			<div id="delete-device-modal" class="modal fade" tabindex="0"></div>
			<div id="create-device-module-modal" class="modal fade" tabindex="0"></div>
			<div id="edit-device-module-modal" class="modal fade" tabindex="0"></div>
			<div id="delete-device-module-modal" class="modal fade" tabindex="0"></div>`
		);

		DevicesTable.renderCreateDeviceModal("#create-device-modal");
		DevicesTable.renderEditDeviceModal("#edit-device-modal");
		DevicesTable.renderDeleteDeviceModal("#delete-device-modal");
		DevicesTable.renderCreateDeviceModuleModal("#create-device-module-modal");
		DevicesTable.renderEditDeviceModuleModal("#edit-device-module-modal");
		DevicesTable.renderDeleteDeviceModuleModal("#delete-device-module-modal");

		this.init();
		await this.handleDynamicContent();
	}

	private init(): void {
		this.createDeviceModal = new CreateDeviceModal(this.app);
		this.editDeviceModal = new EditDeviceModal(this.app);
		this.deleteDeviceModal = new DeleteDeviceModal(this.app);
		this.createDeviceModuleModal = new CreateDeviceModuleModal(this.app);
		this.editDeviceModuleModal = new EditDeviceModuleModal(this.app);
		this.deleteDeviceModuleModal = new DeleteDeviceModuleModal(this.app);

		this.showLimitEl = new SelectElement({
			select: $$("#devices-show-limit") as HTMLSelectElement,
			options: ["1", "2", "3"],
		});
		this.tableBodyEl = $$("#devices-table tbody") as HTMLTableSectionElement;
		this.pagination = new Pagination("#devices-pagination", this.load.bind(this));
	}

	private static renderCreateDeviceModal(selector: string): void {
		$html(selector,
			`<div class="modal-dialog modal-lg">
				<div class="modal-content">
					<div class="modal-header">
						<h4 class="modal-title">Добавить устройство</h4>
						<button type="button" class="close" data-dismiss="modal" aria-label="Close">&times;</button>
					</div>
					<form id="create-device-form">
						<div class="modal-body">
							<div class="form-group">
								<label for="create-device-name">Название</label>
								<input id="create-device-name" class="form-control" name="name" placeholder="Теплица 1">
							</div>
							<div class="form-group">
								<label for="create-device-type">Тип</label>
								<select id="create-device-type" class="custom-select" name="type"></select>
							</div>
							<div class="form-group">
								<label for="create-device-interface">Интерфейс</label>
								<select id="create-device-interface" class="custom-select" name="interface"></select>
							</div>
							<div class="form-group">
								<label for="create-device-address">Адрес</label>
								<select id="create-device-address" class="custom-select" name="address"></select>
							</div>
						</div>
						<div class="modal-footer justify-content-between">
							<button type="button" class="btn btn-default js-dismiss-modal" data-dismiss="modal">Отмена</button>
							<button type="submit" class="btn btn-primary">Сохранить</button>
						</div>
					</form>
				</div>
			</div>`
		);
	}

	private async renderDevice(dev: Device, isActive: boolean): Promise<void> {
		const hasModules = await this.createDeviceModuleModal.hasModules({type: dev.type});
		const deployable = hasModules && dev.modules.length;

		this.tableBodyEl.insertAdjacentHTML("beforeend",
			`<tr data-device-id="${dev.id}">
				<td>
					<span class="js-show-modules">
						${ deployable ? DevicesTable.renderPlusIcon() : DevicesTable.renderMinusIcon()}
					</span>
				</td>
				<td>${dev.name}</td>
				<td>${dev.type}</td>
				<td>${dev.interface}</td>
				<td>${dev.address}</td>
				<td class="text-right">${DevicesTable.renderDeviceActions(dev.type, hasModules)}</td>
			</tr>`
		);

		if (isActive) this.showModules(dev)
	}

	private static renderDeviceActions(deviceType: string, hasModules: boolean): string {
		const actions = DevicesTable.renderCreateDeviceModuleActionButton(hasModules);

		if (deviceType == DeviceType.Embedded) return actions

		return actions +
		`<span class="action-icon edit-icon" data-toggle="modal" data-target="#edit-device-modal">
			<i class="fas fa-edit"></i>
		</span>
		<span class="action-icon delete-icon" data-toggle="modal" data-target="#delete-device-modal">
			<i class="fas fa-trash"></i>
		</span>`
	}

	private static renderCreateDeviceModuleActionButton(hasModules: boolean): string {
		if (!hasModules) return ""

		return (
			`<span class="action-icon text-muted" data-toggle="modal" data-target="#create-device-module-modal">
				M+
			</span>`
		)
	}

	private static renderEditDeviceModal(selector: string): void {
		$html(selector,
			`<div class="modal-dialog modal-lg">
				<div class="modal-content">
					<div class="modal-header">
						<h4 class="modal-title">Изменить устройство</h4>
						<button type="button" class="close" data-dismiss="modal" aria-label="Close">&times;</button>
					</div>
					<form id="edit-device-form">
						<div class="modal-body">
							<div class="form-group">
								<label for="edit-device-name">Название</label>
								<input id="edit-device-name" class="form-control" name="name" placeholder="Теплица 1">
							</div>
							<div class="form-group">
								<label for="edit-device-interface">Интерфейс</label>
								<select id="edit-device-interface" class="custom-select" name="interface"></select>
							</div>
							<div class="form-group">
								<label for="edit-device-address">Адрес</label>
								<select id="edit-device-address" class="custom-select" name="address"></select>
							</div>
						</div>
						<div class="modal-footer justify-content-between">
							<button type="button" class="btn btn-default" data-dismiss="modal">Отмена</button>
							<button type="submit" class="btn btn-primary">Сохранить</button>
						</div>
					</form>
				</div>
			</div>`
		);
	}

	private static renderDeleteDeviceModal(selector: string): void {
		$html(selector,
			`<div class="modal-dialog modal-lg">
				<div class="modal-content">
					<div class="modal-header">
						<h4 class="modal-title">Удалить устройство</h4>
						<button type="button" class="close" data-dismiss="modal" aria-label="Close">&times;</button>
					</div>
					<form id="delete-device-form">
						<div class="modal-body">
							<p>Вы уверенны, что хотите удалить устройство <strong id="delete-device-name"></strong>?</p>
						</div>
						<div class="modal-footer justify-content-between">
							<button type="button" class="btn btn-default" data-dismiss="modal">Отмена</button>
							<button type="submit" class="btn btn-danger">Удалить</button>
						</div>
					</form>
				</div>
			</div>`
		);
	}

	private static renderCreateDeviceModuleModal(selector: string): void {
		$html(selector,
			`<div class="modal-dialog modal-lg">
				<div class="modal-content">
					<div class="modal-header">
						<h4 class="modal-title">Добавить модуль устройства</h4>
						<button type="button" class="close" data-dismiss="modal" aria-label="Close">&times;</button>
					</div>
					<form id="create-device-module-form">
						<div class="modal-body">
							<div class="form-group">
								<label for="create-device-module-name">Название</label>
								<input id="create-device-module-name" class="form-control" name="name" placeholder="Модуль расширитель">
							</div>
							<div class="form-group">
								<label for="create-device-module-type">Тип</label>
								<select id="create-device-module-type" class="custom-select" name="type"></select>
							</div>
						</div>
						<div class="modal-footer justify-content-between">
							<button type="button" class="btn btn-default" data-dismiss="modal">Отмена</button>
							<button type="submit" class="btn btn-primary">Сохранить</button>
						</div>
					</form>
				</div>
			</div>`
		);
	}

	private static renderEditDeviceModuleModal(selector: string): void {
		$html(selector,
			`<div class="modal-dialog modal-lg">
				<div class="modal-content">
					<form id="edit-device-module-form">
						<div class="modal-header">
							<h4 class="modal-title">Изменить модуль</h4>
							<button type="button" class="close" data-dismiss="modal" aria-label="Close">&times;</button>
						</div>
						<div class="modal-body">
							<div class="form-group">
								<label for="edit-device-module-name">Название</label>
								<input id="edit-device-module-name" class="form-control" name="name" placeholder="Реле">
							</div>
						</div>
						<div class="modal-footer justify-content-between">
							<button type="button" class="btn btn-default" data-dismiss="modal">Отмена</button>
							<button type="submit" class="btn btn-primary">Сохранить</button>
						</div>
					</form>
				</div>
			</div>`
		);
	}

	private static renderDeleteDeviceModuleModal(selector: string): void {
		$html(selector,
			`<div class="modal-dialog modal-lg">
				<div class="modal-content">
					<form id="delete-device-module-form">
						<div class="modal-header">
							<h4 class="modal-title">Удалить модуль</h4>
							<button type="button" class="close" data-dismiss="modal" aria-label="Close">&times;</button>
						</div>
						<div class="modal-body">
							<p>Вы уверенны, что хотите удалить модуль <strong id="delete-device-module-name"></strong>?</p>
						</div>
						<div class="modal-footer justify-content-between">
							<button type="button" class="btn btn-default" data-dismiss="modal">Отмена</button>
							<button type="submit" class="btn btn-danger">Удалить</button>
						</div>
					</form>
				</div>
			</div>`
		);
	}

	private async handleDynamicContent(): Promise<void> {
		await this.load(this.currentPage);
		this.handleEvents();
		this.visit();
	}

	private async load(page: number): Promise<void> {
		const limit = Number(this.showLimitEl.value);
		const total = await this.app.call("devices.count") as unknown as number;

		try {
			await this.renderDevices(page, limit);
		} catch (e) {
			console.log(`exception when load ${page} page:`, e);
			if (page - 1 >= 1) {
				return await this.load(page - 1)
			}

			console.log("immediate stopping");
			return
		}

		this.renderPagination(page, limit, total);
	}

	private async renderDevices(page: number, limit: number): Promise<void> {
		this.devices = await this.getDevices({
			page:  page,
			limit: limit,
		});

		const activeDevices = [...document.querySelectorAll(".js-show-modules") as NodeListOf<HTMLSpanElement>]
			.filter(span => span.classList.contains("active"))
			.map(span => span.closest("tr[data-device-id]") as HTMLTableRowElement)
			.map(tableRow => this.devices.find(dev => dev.id === +tableRow.dataset.deviceId))

		const isActive = (id: number): boolean => !!activeDevices.find(dev => dev?.id === id);

		this.tableBodyEl.innerHTML = "";
		for (const device of this.devices) {
			await this.renderDevice(device, isActive(device.id));
		}
	}

	private async getDevices(req: { page: number, limit: number }): Promise<Array<Device>> {
		return await this.app.call("devices.list", req) as unknown as Array<Device>;
	}

	private get currentPage(): number {
		return this.pagination.getCurrentPage();
	}

	private renderPagination(page: number, limit: number, total: number): void {
		const info = this.pagination.render(page, limit, total);
		const suffix = info.from && info.to ? `: показаны ${info.from}-${info.to}`: "";
		$html("#shown-devices-info", `Всего устройств ${total}${suffix}`);
	}

	private handleEvents(): void {
		this.showLimitEl.addEventListener("change", this.onShowLimitChange.bind(this));
		this.tableBodyEl.addEventListener("click", this.onShowDeviceModules.bind(this))
	}

	private async onShowLimitChange(): Promise<void> {
		await this.load(defaultCurrentPage);
	}

	private onShowDeviceModules(event: Event): void {
		const icon = (event.target as HTMLElement).closest(".js-show-modules") as HTMLSpanElement;

		if (!icon) return

		const tableRow = icon.closest("tr[data-device-id]") as HTMLTableRowElement;
		const device = this.devices.find(dev => dev.id === +tableRow.dataset.deviceId);

		if (!device) return

		icon.classList.contains("active") ? this.hideModules(device) : this.showModules(device);
	}

	private showModules(device: Device): void {
		const tableRow = document.querySelector(`tr[data-device-id="${device.id}"]`) as HTMLTableRowElement;
		const icon = tableRow.querySelector(".js-show-modules");

		if (device.modules.length) {
			icon.classList.add("active");
			icon.innerHTML = DevicesTable.renderMinusIcon();
		}

		this.renderModules(tableRow, device);
	}

	private hideModules(device: Device): void {
		const tableRow = document.querySelector(`tr[data-device-id="${device.id}"]`) as HTMLTableRowElement;
		const icon = tableRow.querySelector(".js-show-modules");

		// eslint-disable-next-line @typescript-eslint/no-unused-vars
		device.modules.forEach(_ => tableRow.nextElementSibling?.remove())
		icon.classList.remove("active");
		icon.innerHTML = DevicesTable.renderPlusIcon();
	}

	private renderModules(tableRow: HTMLTableRowElement, device: Device): void {
		tableRow.insertAdjacentHTML("afterend",
			`${device.modules.map(mod => `
			<tr data-device-id="${device.id}" data-module-id="${mod.id}" data-device-module>
				<td></td>
				<td>
					<div class="d-flex align-items-center">
						<span class="text-danger">
							<i class="fas fa-compact-disc"></i>
							<span class="font-weight-bold">Модуль</span>			
						</span>
					</div>
				</td>
				<td>
					<div class="text-dark font-weight-bold">Название</div>
					<div class="text-muted">${mod.name}</div>
				</td>
				<td>
					<div class="text-dark font-weight-bold">Тип</div>
					<div class="text-muted">${mod.type}</div>
				</td>
				<td></td>
				<td>
					<div class="text-right">
						<span class="action-icon edit-icon" data-toggle="modal" data-target="#edit-device-module-modal">
							<i class="fas fa-edit"></i>
						</span>
						<span class="action-icon delete-icon" data-toggle="modal" data-target="#delete-device-module-modal">
							<i class="fas fa-trash"></i>
						</span>
					</div>
				</td>
			</tr>`).join("")}`)
	}

	private static renderPlusIcon(): string {
		return `<i class="bi bi-plus-circle"></i>`
	}

	private static renderMinusIcon(): string {
		return `<i class="bi bi-dash-circle"></i>`
	}

	public visit(): void {
		this.subscriptions.forEach((handler, method) => this.app.subscribe(method, handler));
	}

	public leave(): void {
		this.subscriptions.forEach((handler, method) => this.app.unsubscribe(method, handler));
	}
}
