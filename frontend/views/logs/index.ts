import App from "../../system";
import {$$, $html} from "../index";

export default class Logs {
	private logText: HTMLElement;

	constructor(private app: App) {
	}

	public render(selector: string): void {
		$html(selector, `
			<div class="card">
				<div class="card-header">
					<h3 class="card-title">Системный лог</h3>
				</div>
				<div class="card-body">
					<pre id="system-log" class="mb-0"></pre>
				</div>
			</div>`)

		this.logText = $$('#system-log');

		this.load();
	}

	private load(): void {
		this.app.systemLog().then(
			data => {
				const lines = (data as unknown as string).split("\n");
				for (let i = 0; i < lines.length; i++) {
					lines[i] = lines[i].replace(/([a-z0-9-_]+.go:\d+:)/, '<span class="location">$1</span>');
					if (lines[i].match(' Error ')) {
						lines[i] = lines[i].replace(
							/^(.*)\sdaemon\.err\swebserver\[\d+]:\s(.*)/,
							'<span class="meta">$1</span> <span class="error"><i class="fas fa-exclamation-triangle"></i> $2</span>',
						);
					} else if (lines[i].match(/shutting down|Shutdown/)) {
						lines[i] = lines[i].replace(
							/^(.*)\sdaemon\.err\swebserver\[\d+]:\s(.*)/,
							'<span class="meta">$1</span> <span class="notice"><i class="fas fa-info-circle"></i> $2</span>',
						);
					} else {
						lines[i] = lines[i].replace(
							/^(.*)\sdaemon\.err\swebserver\[\d+]:\s(.*)/,
							'<span class="meta">$1</span> $2',
						);
					}
				}
				$html(this.logText, lines, "\n");
			},
			error => {
				this.app.error(error)
			}
		);
	}
}
