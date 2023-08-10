import {$$, $html} from '../../views';
import {btoa, markAs} from '../index';
import {App} from "../../views/page";

export class Tab {
	constructor(
		private id: string,
		private title: string,
		private app?: App
	) {}

	public getID(): string {
		return this.id;
	}

	public getTitle(): string {
		return this.title;
	}

	public async render() {
		if (this.app != null) {
			await this.app.render(`#${this.id}`)
		}
	}
}

export class Tabs {
	private tabsHeader: HTMLUListElement;
	private tabsBody: HTMLDivElement;

	private isInitialTab = true;

	constructor(
		private id: string,
		private tabs: Tab[]
	) {}

	public async render(headerSelector: string, bodySelector: string) {
		$html(headerSelector, `<ul class="nav nav-tabs" id="${this.id}-tabs" role="tablist"></ul>`);
		$html(bodySelector, `<div class="tab-content" id="${this.id}-tabs-content"></div>`);

		this.tabsHeader = $$(`#${this.id}-tabs`) as HTMLUListElement;
		this.tabsBody = $$(`#${this.id}-tabs-content`) as HTMLDivElement;

		this.tabs.forEach(tab => this.renderTabContainers(tab.getID(), tab.getTitle()));
		this.tabs.forEach(tab => tab.render()); // we iterate twice to have correct event listeners pointers
	}

	private renderTabContainers(id: string, title: string): void {
		this.renderTabHeader(id, title);
		this.renderTabBody(id);
		this.isInitialTab = false;
	}

	private renderTabHeader(id: string, title: string): void {
		this.tabsHeader.innerHTML += (
			`<li class="nav-item" role="presentation">
				<a class="nav-link${markAs(this.isInitialTab, 'active')}" 
				   id="${id}-tab" data-toggle="tab" href="#${id}" role="tab"
				   aria-controls="${id}" aria-selected="${btoa(this.isInitialTab)}">${title}</a>
  			</li>`
		);
	}

	private renderTabBody(id: string): void {
		this.tabsBody.innerHTML += (
			`<div class="tab-pane fade${markAs(this.isInitialTab, 'active', 'show')}"
 				  id="${id}" role="tabpanel" aria-labelledby="${id}-tab"></div>`
		);
	}
}
