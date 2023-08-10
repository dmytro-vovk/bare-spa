import App from "../../system";
import {$$, $html} from "../index";
import SelectElement from "../../core/select";
import Pagination from "../../components/pagination";
import {Alias} from "./types";

export default class AliasesTable {

	private showDeviceAliases: SelectElement
	private tableBodyEl: HTMLTableSectionElement
	private pagination: Pagination

	private aliases: Array<Alias>

	private subscriptions = new Map<string, (data) => void>([]);

	constructor(private readonly app: App) {}

	public async render(selector: string): Promise<void> {
		$html(selector,
			`<div class="card">
				<div class="card-header">
					<div class="row">
						<div class="col d-flex">
							<label class="d-flex align-items-center mb-0" for="show-device-aliases">
								Привязки устройства <select id="show-device-aliases" class="custom-select ml-2 mr-2">
									<option>Микроконтроллер</option>
								</select>
							</label>
						</div>
						<div class="col d-flex justify-content-end">
							<button type="button" class="btn btn-sm btn-default" data-toggle="modal" data-target="#create-alias-modal">
								<i class="fas fa-plus"></i> Добавить привязку
							</button>
						</div>
					</div>
				</div>
				<div class="card-body">
					<div class="row p-0">
						<table id="aliases-table" class="table table-last-underlined">
							<thead class="thead-light">
								<tr>
									<th rowspan="1" colspan="1">#</th>
									<th rowspan="1" colspan="1">Название</th>
									<th rowspan="1" colspan="1">Состояние</th>
									<th rowspan="1" colspan="1">Путь</th>
									<th rowspan="1" colspan="1">Действия</th>
								</tr>
							</thead>
							<tbody></tbody>
						</table>
					</div>
				</div>
				<div class="card-footer">
					<div class="row align-items-center">
						<div class="col">
							<p id="shown-aliases-info" class="mb-0"></p>
						</div>
						<div class="col">
							<nav id="aliases-pagination" aria-label="Available aliases navigation"></nav>
						</div>
					</div>
				</div>
			</div>`
		);

		this.init();
		await this.handleDynamicContent();
	}

	private init(): void {
		this.showDeviceAliases = new SelectElement({
			select: $$("#show-device-aliases") as HTMLSelectElement,
			options: ["Микроконтроллер"],
		});
		this.tableBodyEl = $$("#aliases-table tbody") as HTMLTableSectionElement;
		this.pagination = new Pagination("#aliases-pagination", this.load.bind(this));
	}

	private async handleDynamicContent(): Promise<void> {
		await this.load(this.currentPage);
		this.handleEvents();
		this.visit();
	}

	private async load(page: number): Promise<void> {
		const limit = 3;
		const total = await this.app.call("aliases.count") as unknown as number;

		try {
			await this.renderAliases(page, limit);
		} catch (e) {
			console.log(`exception when load ${page} page:`, e);
			if (page - 1 >= 1) {
				return await this.load(page - 1)
			}

			console.log("immediate stopping");
			return
		}

		this.renderPagination(page, limit, total);
	}

	private async renderAliases(page: number, limit: number): Promise<void> {
		this.aliases = await this.getAliases({
			page:  page,
			limit: limit,
		});

		this.tableBodyEl.innerHTML = "";
		for (const alias of this.aliases) {
			this.renderAlias(alias);
		}
	}

	private async getAliases(req: { page: number, limit: number }): Promise<Array<Alias>> {
		return await this.app.call("aliases.list", req) as unknown as Array<Alias>;
	}

	private renderAlias(alias: Alias): void {
		this.tableBodyEl.insertAdjacentHTML("beforeend",
			`<tr data-alias-id="${alias.id}">
				<td>${alias.id}</td>
				<td>${alias.name}</td>
				<td>${alias.state}</td>
				<td>${alias.path}</td>
				<td class="text-left">
					<span class="action-icon edit-icon" data-toggle="modal" data-target="#edit-alias-modal">
						<i class="fas fa-edit"></i>
					</span>
					<span class="action-icon delete-icon" data-toggle="modal" data-target="#delete-alias-modal">
						<i class="fas fa-trash"></i>
					</span>
				</td>
			</tr>`
		);
	}

	private get currentPage(): number {
		return this.pagination.getCurrentPage();
	}

	private renderPagination(page: number, limit: number, total: number): void {
		this.pagination.render(page, limit, total);
		// const info = this.pagination.render(page, limit, total);
		// const suffix = info.from && info.to ? `: показаны ${info.from}-${info.to}`: "";
		// $html("#shown-aliases-info", `Всего привязок ${total}${suffix}`);
		$html("#shown-aliases-info", `Всего привязок 3: показаны 1-3`);
	}

	private handleEvents(): void {
		console.log("handle events");
	}

	public visit(): void {
		this.subscriptions.forEach((handler, method) => this.app.subscribe(method, handler));
	}

	public leave(): void {
		this.subscriptions.forEach((handler, method) => this.app.unsubscribe(method, handler));
	}
}
