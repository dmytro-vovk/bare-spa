export function $$(selector: string): HTMLElement | null {
	return document.querySelector(selector)
}

export function $$$(selector: string): NodeListOf<Element> | null {
	return document.querySelectorAll(selector)
}

export function $html(element: string | HTMLElement, content: string | string[], glue = ''): void {
	if (typeof element === 'string') {
		const el = $$(element)
		if (el === null) {
			return
		}
		element = el
	}
	if (Array.isArray(content)) {
		content = content.join(glue)
	}
	element.innerHTML = content
}

export function $replace(element: string | HTMLElement, content: string | string[], glue = ''): void {
	if (typeof element === 'string') {
		const el = $$(element)
		if (el === null) {
			return
		}
		element = el
	}
	if (Array.isArray(content)) {
		content = content.join(glue)
	}

	const template = document.createElement('template')
	content = content.trim()
	template.innerHTML = content

	element.parentNode?.replaceChild(template.content, element)
}

export function $text(element: string | HTMLElement, content: number | string): void {
	if (typeof content === 'number') {
		content = content.toString()
	}

	if (typeof element === 'string') {
		const el = $$(element)
		if (el === null) {
			return
		}
		el.innerText = content
	} else {
		element.innerText = content
	}
}

export function $onClick<K extends keyof HTMLElementEventMap>(element: HTMLElement | string, callback: (this: HTMLSelectElement, ev: HTMLElementEventMap[K]) => void): void {
	if (typeof element === 'string') {
		const el = $$$(element)
		if (el === null) {
			return
		}
		el.forEach(e => {
			e.addEventListener('click', callback)
		})
	} else {
		element.addEventListener('click', callback)
	}
}

export function $onChange<K extends keyof HTMLElementEventMap>(element: string | HTMLElement, callback: (this: HTMLSelectElement, ev: HTMLElementEventMap[K]) => void): void {
	if (typeof element === 'string') {
		const el = $$(element)
		if (el === null) {
			return
		}
		el.addEventListener('change', callback)
	} else {
		element.addEventListener('change', callback)
	}
}

export function $newElement(tagName: string, options = {}): HTMLElement {
	const e = document.createElement(tagName)
	for (const prop of Object.keys(options)) {
		e[prop] = options[prop]
	}
	return e
}

export function $buttonBusy(selector: string | HTMLButtonElement): void {
	if (typeof selector === 'string') {
		selector = $$(selector) as HTMLButtonElement
	}
	selector.innerHTML = '<i class="fa-solid fa-spinner fa-spin"></i> ' + selector.innerHTML
	selector.disabled = true
}

export function $buttonReady(selector: string | HTMLButtonElement): void {
	if (typeof selector === 'string') {
		selector = $$(selector) as HTMLButtonElement
	}
	const i = selector.querySelector('i:first-child') as HTMLElement
	selector.removeChild(i)
	selector.disabled = false
}

export function $money(amount: number): string {
	return Number(amount / 100).toLocaleString('uk', {minimumFractionDigits: 2})
}

export function $date(value: number | string | Date): string {
	const options: Intl.DateTimeFormatOptions = {
		year: 'numeric',
		month: 'long',
		day: 'numeric'
	}
	return new Date(value).toLocaleDateString('uk', options)
}

export function $dateTime(value: number | string | Date): string {
	const options: Intl.DateTimeFormatOptions = {
		year: 'numeric',
		month: 'long',
		day: 'numeric',
		hour: 'numeric',
		minute: 'numeric',
		second: 'numeric'
	}
	return new Date(value).toLocaleDateString('uk', options)
}

export function $modalShow(selector: string): void {
	(jQuery(selector) as JQuery).modal('show')
}

export function $modalHide(selector: string): void {
	(jQuery(selector) as JQuery).modal('hide')
}

export function $bytes(n: number): string {
	const t = 1024
	let suffix = ''
	if (n > t) {
		n /= t
		suffix = 'KB'
	}

	if (n > t) {
		n /= t
		suffix = 'MB'
	}

	if (n > t) {
		n /= t
		suffix = 'GB'
	}

	if (n > t) {
		n /= t
		suffix = 'TB'
	}

	return `${n.toFixed(2)}&thinsp;${suffix}`
}
