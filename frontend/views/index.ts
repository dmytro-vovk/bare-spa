export function $$(selector: string): HTMLElement {
	return document.querySelector(selector);
}

export function $html(element: string | HTMLElement, content: string | string[], glue = ''): void {
	if (typeof element === 'string') {
		element = $$(element);
	}
	if (Array.isArray(content)) {
		content = content.join(glue);
	}
	element.innerHTML = content;
}

export function $text(element: string | HTMLElement, content: number | string): void {
	if (typeof element === 'string') {
		element = $$(element);
	}
	if (typeof content === 'number') {
		content = content.toString();
	}
	element.innerText = content;
}

export function $onClick<K extends keyof HTMLElementEventMap>(element: HTMLElement | string, callback: (this: HTMLSelectElement, ev: HTMLElementEventMap[K]) => any): void {
	if (typeof element === 'string') {
		element = document.querySelector(element) as HTMLElement;
	}
	element.addEventListener('click', callback);
}

export function $onChange<K extends keyof HTMLElementEventMap>(element: string | HTMLElement, callback: (this: HTMLSelectElement, ev: HTMLElementEventMap[K]) => any): void {
	if (typeof element === 'string') {
		element = $$(element) as HTMLElement;
	}
	element.addEventListener('change', callback);
}

export function $newElement(tagName: string, options = {}): HTMLElement {
	const e = document.createElement(tagName);
	for (const prop of Object.keys(options)) {
		e[prop] = options[prop];
	}
	return e
}

type stringer = number|string

export function renderOptions(obj: Array<stringer> | Record<string, stringer>): string {
	const option = (k: stringer, v: stringer) => `<option value="${v}">${k}</option>`;

	if (Array.isArray(obj)) {
		return obj.map(k => option(k, k)).join('');
	}

	return Object.keys(obj).filter(k => isNaN(Number(k))).map(k => option(k, obj[k])).join("")
}

export function $insertAfter(referenceNode: string | HTMLElement, newNode: HTMLElement): HTMLElement {
	if (typeof referenceNode === 'string') {
		referenceNode = $$(referenceNode);
	}

	referenceNode.parentNode.insertBefore(newNode, referenceNode.nextSibling);
	return newNode;
}

// todo: write func $setValue (for different inputs)
export function $setInput(element: string | HTMLInputElement, value: number | string): void {
	if (typeof element === 'string') {
		element = $$(element) as HTMLInputElement;
	}

	if (typeof value === 'number') {
		value = value.toString();
	}

	if (element.type === "password" || element.type === "text") {
		element.value = value;
		return
	}

	console.error(`selected element isn't "password" or "text", it is "${element.type}"`)
}

export function $setCheckbox(element: string | HTMLInputElement, checked: boolean): void {
	if (typeof element === 'string') {
		element = $$(element) as HTMLInputElement;
	}

	if (element.type !== "checkbox") {
		console.error(`selected element isn't "checkbox", it is "${element.type}"`)
		return
	}

	element.checked = checked;
}

export function $setSelect(element: string | HTMLSelectElement, value: number | string): void {
	if (typeof element === 'string') {
		element = $$(element) as HTMLSelectElement;
	}

	element.selectedIndex = Array.from(element.options).findIndex(opt => opt.value == value);
}

export const range = n => [...Array(n).keys()];
