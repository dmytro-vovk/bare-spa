import App from "../../system";
import ModalWithForm from "../../core/modal";
import {Device, ModuleConfig, ModuleType, ModuleTypeValue} from "./types";
import SelectElement from "../../core/select";

export default class CreateDeviceModuleModal extends ModalWithForm {

	private device: Device
	private modules: Array<ModuleTypeValue>

	constructor(private readonly app: App) {
		super({
			modal: "#create-device-module-modal",
			form:  "#create-device-module-form",
		});

		this.init();
	}

	public async hasModules(req: { type: string }): Promise<boolean> {
		const modules = await this.getSupportedModules(req);
		return modules.length !== 0;
	}

	private init(): void {
		this.handleEvents();
	}

	private async load(id: number): Promise<void> {
		this.device = await this.app.call("devices.read", { id: id }) as unknown as Device;
		this.modules = await this.getSupportedModules({ type: this.device.type });
	}

	private handleEvents(): void {
		this.onModalShow(this.handleModalOpen.bind(this));
		this.onModalHidden(this.handleModalClose.bind(this));
		this.onSubmit(this.handleSubmit.bind(this));
	}

	private async handleModalOpen(event: MouseEvent): Promise<void> {
		const button = event.relatedTarget as HTMLButtonElement;
		const tableRow = button.closest("tr[data-device-id]") as HTMLTableRowElement;
		await this.load(+tableRow.dataset.deviceId);

		const modulesEl = this.getElement("type") as SelectElement;
		modulesEl.setOptions(this.modules);
	}

	private handleModalClose(): void {
		this.clear();
		this.device = null;
		this.modules = null;
	}

	private handleSubmit(event: Event): void {
		event.preventDefault(); // don't remove!
		this.createModule({
			deviceID: this.device.id,
			type: this.value("type") as ModuleType,
			config: {
				name: this.value("name"),
			}
		});
	}

	private createModule(req: { deviceID: number, type: ModuleType, config: ModuleConfig }): void {
		this.app.call("devices.createModule", req)
			.catch(err => this.app.error(err))
			.finally(() => this.close());
	}

	private async getSupportedModules(req: { type: string }): Promise<Array<ModuleTypeValue>> {
		return await this.app.call("devices.supportedModules", req) as unknown as Array<ModuleTypeValue>;
	}
}
