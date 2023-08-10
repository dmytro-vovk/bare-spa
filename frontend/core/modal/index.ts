import Form, {FormConfig} from "../form";

export declare interface ModalWithFormConfig extends FormConfig {
	modal: string | HTMLDivElement
}

export default class ModalWithForm extends Form {
	protected readonly modal: HTMLDivElement
	private readonly cancel: HTMLButtonElement

	constructor(config: ModalWithFormConfig) {
		super(config);

		if (typeof config.modal === "string") {
			this.modal = document.querySelector(config.modal) as HTMLDivElement;
		}

		this.cancel = this.modal.querySelector('button[data-dismiss="modal"]');
	}

	public onModalShow(cb: (event: Event) => void): void {
		$(this.modal).on("show.bs.modal", cb);
	}

	public onModalShown(cb: (event: Event) => void): void {
		$(this.modal).on("shown.bs.modal", cb);
	}

	public onModalHide(cb: (event: Event) => void): void {
		$(this.modal).on("hide.bs.modal", cb);
	}

	public onModalHidden(cb: (event: Event) => void): void {
		$(this.modal).on("hidden.bs.modal", cb);
	}

	public close(): void {
		this.cancel.click();
	}
}
