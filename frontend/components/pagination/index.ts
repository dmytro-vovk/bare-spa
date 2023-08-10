export const defaultCurrentPage = 1;
const defaultEntriesLimit = 1;
const defaultTotalEntries = 0;
const defaultTotalPages = 1;
const defaultPagesAround = 1;

type info = {
	from: number,
	to: number,
}

export default class Pagination {

	private currentPage = defaultCurrentPage
	private entriesLimit = defaultEntriesLimit
	private totalPages = defaultTotalPages
	private totalEntries = defaultTotalEntries
	private pagesAround = defaultPagesAround

	private readonly list: HTMLUListElement;

	constructor(
		selector: string,
		private readonly load: (page: number) => void,
	) {
		const element = document.querySelector(selector);
		element.innerHTML = (
			`<ul class="pagination justify-content-end mb-0"></ul>`
		);

		this.list = element.firstElementChild as HTMLUListElement;
		this.list.addEventListener("click", this.transition.bind(this));
	}

	public render(currentPage: number, entriesLimit: number, totalEntries: number, pagesAround?: number): info {
		this.init(currentPage, entriesLimit, totalEntries, pagesAround);
		this.list.innerHTML = this.renderItems();
		return this.info
	}

	private init(currentPage: number, entriesLimit: number, totalEntries: number, pagesAround?: number): void {
		this.setEntriesLimit(entriesLimit);
		this.setTotalEntries(totalEntries);
		this.setTotalPages(totalEntries);
		this.setCurrentPage(currentPage);
		this.setPagesAround(pagesAround);
	}

	public getCurrentPage(): number {
		return this.currentPage
	}

	private get info(): info {
		if ( this.totalEntries === 0 ) {
			return {from: 0, to: 0};
		}

		const maxOnPage = this.currentPage * this.entriesLimit;
		return {
			from:  (this.currentPage - 1) * this.entriesLimit + 1,
			to:    maxOnPage <= this.totalEntries ? maxOnPage : this.totalEntries,
		};
	}

	private setEntriesLimit(limit: number): void {
		this.entriesLimit = limit >= defaultEntriesLimit ? limit : defaultEntriesLimit;
	}

	private setTotalEntries(entries: number): void {
		this.totalEntries = entries >= defaultTotalEntries ? entries : defaultTotalEntries;
	}

	private setTotalPages(entries: number): void {
		this.totalPages = entries >= defaultTotalPages ? Math.ceil(entries / this.entriesLimit) : defaultTotalPages;
	}

	private setCurrentPage(page: number): void {
		this.currentPage = page >= defaultCurrentPage && page <= this.totalPages ? page : defaultCurrentPage;
	}

	private setPagesAround(around: number): void {
		this.pagesAround = around >= defaultPagesAround ? around : defaultPagesAround;
	}

	private renderLeftArrow(): string {
		const isLast = this.currentPage === 1;
		const icon = `<i class="fas fa-angle-left"></i>`;
		const link = isLast ? Pagination.renderPageLink(icon) : Pagination.renderPageLink(icon, this.currentPage - 1);
		const state = isLast ? "disabled" : "";
		return Pagination.renderPageItem(link, state)
	}

	private static renderPageItem(link: string, state: "active" | "disabled" | "" = "") {
		return `<li class="page-item ${state}">${link}</li>`
	}

	private static renderPageLink(html: string, page?: number, isActive?: boolean) {
		if ( !isActive && page ) {
			return `<a role="button" class="page-link" data-page="${page}">${html}</a>`
		}

		return `<span class="page-link">${html}</span>`
	}

	private async transition(event: Event): Promise<void> {
		const link = (event.target as HTMLElement).closest("[data-page]") as HTMLAnchorElement;

		const page = Number(link?.dataset.page ?? 0);
		if ( page === this.currentPage || page < 1 || page > this.totalPages ) {
			return
		}

		await this.load(page);
	}

	private renderItems(): string {
		let hasEllipsis: boolean;
		let items = this.renderLeftArrow();

		for (let page = 1; page <= this.totalPages; page++) {
			if ( this.isVisible(page) ) {
				items += this.renderItem(page);
				hasEllipsis = false;
			} else if ( !hasEllipsis ) {
				items += Pagination.renderEllipsis();
				hasEllipsis = true
			}
		}

		items += this.renderRightArrow();
		return items
	}

	private isVisible(page: number): boolean {
		const start = this.currentPage - this.pagesAround;
		const end = this.currentPage + this.pagesAround;
		const around = [...Array(end - start + 1).keys()]
			.map(p => p + start)
			.filter(p => p > 1 && p < this.totalPages);

		return [1, ...around, this.totalPages].includes(page)
	}

	private renderItem(page: number): string {
		const isActive = this.currentPage === page;
		const link = Pagination.renderPageLink(String(page), page, isActive);
		const state = isActive ? "active" : "";
		return Pagination.renderPageItem(link, state)
	}

	private static renderEllipsis(): string {
		const icon = `<i class="fas fa-ellipsis-h"></i>`;
		return Pagination.renderPageItem(Pagination.renderPageLink(icon), "disabled")
	}

	private renderRightArrow(): string {
		const isLast = this.currentPage === this.totalPages;
		const icon = `<i class="fas fa-angle-right"></i>`;
		const link = isLast ? Pagination.renderPageLink(icon) : Pagination.renderPageLink(icon, this.currentPage + 1);
		const state = isLast ? "disabled" : "";
		return Pagination.renderPageItem(link, state)
	}
}
