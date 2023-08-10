import App from "../../system";
import ModalWithForm from "../../core/modal";
import {$$} from "../index";
import {Device} from "./types";

export default class DeleteDeviceModal extends ModalWithForm {

	private device: Device

	constructor(private readonly app: App) {
		super({
			modal: "#delete-device-modal",
			form:  "#delete-device-form",
		});

		this.init();
	}

	private init(): void {
		this.handleEvents();
	}

	private async load(id: number): Promise<void> {
		this.device = await this.app.call("devices.read", { id: id }) as unknown as Device;
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

		$$("#delete-device-name").innerText = this.device.name;
	}

	private handleModalClose(): void {
		this.clear();
		this.device = null;
	}

	private handleSubmit(event: Event): void {
		event.preventDefault(); // don't remove!
		this.deleteDevice({ id: this.device.id });
	}

	private deleteDevice(req: { id: number }): void {
		this.app.call("devices.delete", req)
			.catch(err => this.app.error(err))
			.finally(() => this.close());
	}
}
