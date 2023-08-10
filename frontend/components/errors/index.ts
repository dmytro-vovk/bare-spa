import {$insertAfter, $newElement} from "../../views";

export const enum ErrorCode {
	NoType,
	ParsingErr,
	BadRequest,
	NotFound,
	InvalidParams,
	ValidationErr,
	CommandErr,
}

// todo: what about jsonrpc errors?

export function setInvalidFeedback(form = {}, data = {}): void {
	clearForm(form);
	form["feedbacks"] = Object.keys(data).map(field => {
		const input = form[field] as HTMLInputElement | HTMLSelectElement;
		const message = data[field] as string;
		if (message && input) {
			input.classList.add("is-invalid");
			return $insertAfter(
				input,
				$newElement("div", {
					id: `${input.id}-feedback`,
					className: "invalid-feedback",
					innerText: message,
				}),
			) as HTMLDivElement;
		}
	});
}

export function clearForm(form = {}): void {
	(form["feedbacks"] as HTMLDivElement[]).forEach(feedback => {
		document.getElementById(feedback.id.replace("-feedback", "")).classList.remove("is-invalid");
		feedback.remove();
	});
}