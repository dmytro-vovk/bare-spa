import App from "../../system";
import ModalWithForm from "../../core/modal";
import {Module, ModuleConfig} from "./types";

export default class EditDeviceModuleModal extends ModalWithForm {

	private deviceID: number
	private module: Module

	constructor(private readonly app: App) {
		super({
			modal: "#edit-device-module-modal",
			form:  "#edit-device-module-form",
		});

		this.init();
	}

	private init(): void {
		this.handleEvents();
	}

	private async load(deviceId: number, moduleId: number): Promise<void> {
		this.deviceID = deviceId;
		this.module = await this.app.call("devices.readModule", { deviceId: deviceId, moduleId: moduleId }) as unknown as Module;
	}

	private handleEvents(): void {
		this.onModalShow(this.handleModalOpen.bind(this));
		this.onModalHidden(this.handleModalClose.bind(this));
		this.onSubmit(this.handleSubmit.bind(this));
	}

	private async handleModalOpen(event: MouseEvent): Promise<void> {
		const button = event.relatedTarget as HTMLButtonElement;
		const tableRow = button.closest("tr[data-device-module]") as HTMLTableRowElement;
		await this.load(+tableRow.dataset.deviceId, +tableRow.dataset.moduleId);

		const moduleNameEl = this.getElement("name") as HTMLInputElement;
		moduleNameEl.value = this.module.name;
	}

	private handleModalClose(): void {
		this.clear();
		this.deviceID = null;
		this.module = null;
	}

	private handleSubmit(event: Event): void {
		event.preventDefault(); // don't remove!
		this.editModule({
			deviceID: this.deviceID,
			moduleID: this.module.id,
			config: {
				name: this.value("name"),
			}
		});
	}

	private editModule(req: { deviceID: number, moduleID: number, config: ModuleConfig }): void {
		this.app.call("devices.updateModule", req)
			.catch(err => this.app.error(err))
			.finally(() => this.close());
	}
}
