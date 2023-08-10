type Option = {
	text: string,
	value: string,
	selected?: boolean,
}

type Options = Array<string> | Array<number> | Array<Option>

export declare interface SelectConfig {
	select: HTMLSelectElement
	options?: Options
	prompt?: string | number
}

export default class SelectElement {

	private readonly select: HTMLSelectElement
	private options: Array<Option> = []
	private default: Option
	private prompt?: Option

	constructor(config: SelectConfig) {
		this.select = config.select;
		this.usePrompt(config.prompt);
		this.setOptions(config.options);
	}

	public get name(): string {
		return this.select.name
	}

	public closest<K extends keyof HTMLElementTagNameMap>(selector: K): HTMLElementTagNameMap[K] | null;
	public closest<K extends keyof SVGElementTagNameMap>(selector: K): SVGElementTagNameMap[K] | null;
	public closest<E extends Element = Element>(selectors: string): E | null;
	public closest(selector: string): Element {
		return this.select.closest(selector)
	}

	public setAttribute(qualifiedName: string, value: string): void {
		this.select.setAttribute(qualifiedName, value);
	}

	public removeAttribute(qualifiedName: string): void {
		this.select.removeAttribute(qualifiedName);
	}

	public usePrompt(prompt: string | number, asValue?: boolean): void {
		if ( prompt != null ) {
			const message = String(prompt);
			this.prompt = {
				text: message,
				value: asValue ? message : "",
			};
		}
	}

	public setOptions(options: Options): void {
		this.discard();
		this.addOptions(options);
	}

	public addOptions(options: Options): void {
		options = this.asOptionsArray(options);

		if ( !this.default ) {
			const [first] = options;
			this.default = options.find(opt => opt.selected) ?? first;
		}

		if ( this.prompt ) {
			this.default = this.prompt;
			options.push(this.prompt);
		}

		this.renderOptions(options);
	}

	private asOptionsArray(options: Options): Array<Option> {
		return (options?.map(opt => {
			if ( typeof opt === "string" || typeof opt === "number" ) {
				return { text: String(opt), value: String(opt) }
			}

			return opt
		}) ?? []) as Array<Option>;
	}

	private renderOptions(options: Array<Option>): void {
		options.forEach(option => {
			if ( !this.hasOption(option) ) {
				this.options.push(option);
				this.select.insertAdjacentHTML("beforeend", this.renderOption(option));
			}
		});
	}

	private hasOption(option: Option): boolean {
		return !!this.options.find(opt => opt.value === option.value)
	}

	private renderOption(option: Option): string {
		const attr = (cond: boolean, name: string) => cond ? name : "";

		const isDefault = option.value === this.default.value;
		const isPrompt = option.text === this.prompt?.text;

		const selected = attr(isDefault, "selected");
		const disabled = attr(isPrompt, "disabled");
		const hidden = attr(isPrompt, "hidden");

		const attributes = `${selected} ${disabled} ${hidden}`.replace(/\s+/g, " ");

		return `<option value="${option.value}" ${attributes}>${option.text}</option>`
	}

	private discard(): void {
		this.select.innerHTML = "";
		this.options = [];
		this.default = null;
	}

	public get value(): string {
		const getDefault = () => {
			this.restore();
			return this.default?.value ?? ""
		}

		return this.options.find(opt => opt.value === this.select.value)?.value ?? getDefault();
	}

	public restore(): void {
		const options = this.options;
		this.discard();
		this.addOptions(options);
	}

	public addEventListener<K extends keyof HTMLElementEventMap>(type: K, listener: (this: HTMLElement, ev: HTMLElementEventMap[K]) => any, options?: boolean | AddEventListenerOptions): void {
		this.select.addEventListener(type, listener, options);
	}
}
