import {$html} from "../index";

export interface App {
	render(selector: string): void;
}

export default class Page {
	constructor(
		private readonly title: string,
		private readonly app: App
	) {}

	public render(selector: string): void {
		$html(selector,
			`<section class="content-header">
				<div class="container-fluid">
					<div class="row mb-2">
                    	<div class="col-sm-12">
                        	<h1>${this.title}</h1>
                        </div>
                    </div>
                </div>
            </section>
            <section class="content mb-5">
                <div class="container-fluid">
                    <div class="row">
                        <div class="col-12" id="content"></div>
                    </div>
                </div>
            </section>`
		);

		this.app.render('#content');
	}
}
