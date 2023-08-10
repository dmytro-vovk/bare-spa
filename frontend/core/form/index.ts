import SelectElement from "../select";

export declare interface FormConfig {
	form: string | HTMLFormElement
}

type FormElement = HTMLInputElement | SelectElement

export default class Form {

	private readonly form: HTMLFormElement
	private readonly elements: Array<FormElement>
	private readonly submit: HTMLButtonElement

	constructor(config: FormConfig) {
		if (typeof config.form === "string") {
			this.form = document.querySelector(config.form) as HTMLFormElement;
		}

		const elements = Array.from(this.form.elements) as Array<HTMLInputElement | HTMLSelectElement>;
		this.elements = elements.map(element => {
			if ( element instanceof HTMLSelectElement ) {
				return new SelectElement({ select: element as HTMLSelectElement })
			}

			return element as HTMLInputElement
		});

		this.submit = this.form.querySelector('button[type="submit"]');
	}

	protected getElement(name: string): FormElement {
		return this.elements.find(el => el.name === name)
	}

	public value(name: string): string {
		return this.getElement(name)?.value || ""
	}

	public clear(): void {
		this.elements.forEach(element => {
			if (element instanceof SelectElement ) {
				element.restore();
				return
			}

			element.value = "";
		});
	}

	public hideField(name: string): void {
		this.getField(name)?.classList.add("d-none");
	}

	public isFieldShown(name: string): boolean {
		const field = this.getField(name);
		return field ? !field.classList.contains("d-none") : false
	}

	public showField(name: string): void {
		this.getField(name)?.classList.remove("d-none");
	}

	private getField(name: string): HTMLDivElement {
		return this.getElement(name)?.closest(".form-group")
	}

	public onSubmit(cb: (event: Event) => void): void {
		this.form.addEventListener("submit", cb);
	}

	public disableForm(): void {
		this.disableSubmit();
		this.elements.forEach(el => el.setAttribute("disabled", ""));
	}

	public disableSubmit(): void {
		this.submit?.setAttribute("disabled", "");
	}

	public disableField(name: string): void {
		this.getElement(name)?.setAttribute("disabled", "");
	}

	public enableForm(): void {
		this.enableSubmit();
		this.elements.forEach(el => el.removeAttribute("disabled"));
	}

	public enableSubmit(): void {
		this.submit?.removeAttribute("disabled");
	}

	public enableField(name: string): void {
		this.getElement(name)?.removeAttribute("disabled");
	}
}
