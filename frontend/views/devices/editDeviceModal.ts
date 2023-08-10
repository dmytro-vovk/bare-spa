import App from "../../system";
import {Device, DeviceConfig, DeviceInterface, DeviceInterfaceValue} from "./types";
import ModalWithForm from "../../core/modal";
import SelectElement from "../../core/select";

export default class EditDeviceModal extends ModalWithForm {

	private device: Device
	private interfaces: Array<DeviceInterfaceValue>
	private addresses: Array<number>

	constructor(private readonly app: App) {
		super({
			modal: "#edit-device-modal",
			form:  "#edit-device-form",
		});

		this.init();
	}

	private init(): void {
		this.handleEvents();
	}

	private async load(id: number): Promise<void> {
		this.device = await this.app.call("devices.read", { id: id }) as unknown as Device;
		this.interfaces = await this.app.call("devices.supportedInterfaces", { type: this.device.type }) as unknown as Array<DeviceInterfaceValue>;
		this.addresses = await this.app.call("devices.availableAddresses") as unknown as Array<number>;
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

		const deviceNameEl = this.getElement("name") as HTMLInputElement;
		deviceNameEl.value = this.device.name;

		const interfacesEl = this.getElement("interface") as SelectElement;
		interfacesEl.usePrompt(this.device.interface, true);
		interfacesEl.setOptions(this.interfaces);

		const addressesEl = this.getElement("address") as SelectElement;
		addressesEl.usePrompt(this.device.address, true);
		addressesEl.setOptions(this.addresses);
	}

	private handleModalClose(): void {
		this.clear();
		this.device = null;
		this.interfaces = null;
		this.addresses = null;
	}

	private handleSubmit(event: Event): void {
		event.preventDefault(); // don't remove!
		this.editDevice({
			id: this.device.id,
			config: {
				name:      this.value("name"),
				interface: DeviceInterface[this.value("interface")],
				address:   Number(this.value("address")),
			}
		});
	}

	private editDevice(req: { id: number, config: DeviceConfig }): void {
		this.app.call("devices.update", req)
			.catch(err => this.app.error(err))
			.finally(() => this.close());
	}
}
