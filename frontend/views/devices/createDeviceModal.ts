import App from "../../system";
import ModalWithForm from "../../core/modal";
import {DeviceConfig, DeviceInterface, DeviceInterfaceValue, DeviceTypeValue} from "./types";
import SelectElement from "../../core/select";

export default class CreateDeviceModal extends ModalWithForm {

	constructor(private readonly app: App) {
		super({
			modal: "#create-device-modal",
			form:  "#create-device-form",
		});

		this.init();
	}

	private init(): void {
		(this.getElement("type") as SelectElement).usePrompt("Выберите тип устройства");

		this.handleEvents();
	}

	private handleEvents(): void {
		this.onModalShow(this.handleModalOpen.bind(this));
		this.onTypeChange(this.handleTypeChange.bind(this));
		this.onModalHidden(this.handleModalClose.bind(this));
		this.onSubmit(this.handleSubmit.bind(this));
	}

	private async handleModalOpen(): Promise<void> {
		const types = await this.app.call("devices.supportedTypes") as unknown as Array<DeviceTypeValue>;
		(this.getElement("type") as SelectElement).setOptions(types);

		this.hideField("interface");
		this.hideField("address");
		this.disableSubmit();
	}

	private onTypeChange(cb: (event: Event) => void): void {
		this.getElement("type").addEventListener("change", cb);
	}

	private async handleTypeChange(): Promise<void> {
		const interfaces = await this.app.call("devices.supportedInterfaces", { type: this.value("type") }) as unknown as Array<DeviceInterfaceValue>;
		(this.getElement("interface") as SelectElement).setOptions(interfaces);
		this.showField("interface");

		if ( !this.isFieldShown("address") ) {
			const addresses = await this.app.call("devices.availableAddresses") as unknown as Array<number>;
			(this.getElement("address") as SelectElement).setOptions(addresses);
			this.showField("address");
		}

		this.enableSubmit();
	}

	private handleModalClose(): void {
		this.clear();
	}

	private handleSubmit(event: Event): void {
		event.preventDefault(); // don't remove!
		this.createDevice({
			type: this.value("type"),
			config: {
				name: this.value("name"),
				interface: this.value("interface") as DeviceInterface,
				address: Number(this.value("address")),
			},
		});
	}

	private createDevice(req: { type: string, config: DeviceConfig }): void {
		this.app.call("devices.create", req)
			.catch(err => this.app.error(err))
			.finally(() => this.close());
	}
}
