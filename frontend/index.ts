import App from "./system";
import Dashboard from "./views/dashboard";
import Navigo from "navigo";
import Page from "./views/page";
import RPC from "./system/rpc";
import Settings from "./views/settings";
import Logs from "./views/logs";
import DevicesTable from "./views/devices";
import AliasesTable from "./views/aliases";

const contentSelector = "#content-wrapper";
const wsURL = window.location.protocol.replace("http", "ws") + "//" + window.location.host + "/ws";
const app = new App(new RPC(wsURL));
const dashboard = new Dashboard(app);
const homePage = new Page("Dashboard", dashboard);
const logPage = new Page("Log", new Logs(app));
const settings = new Settings(app)
const settingsPage = new Page("Settings", settings);

const devices = new DevicesTable(app);
const devicesPage = new Page(`<i class="fal fa-clipboard-list"></i> Список подключенных устройств`, devices);

const aliases = new AliasesTable(app);
const aliasesPage = new Page(`<i class="fal fa-clipboard-list"></i> Список привязок`, aliases);

app.setRouter(new Navigo("/")
	.on("/", () => {
		app.sideBarToggle("/");
		homePage.render(contentSelector);
	}, {
		leave(done) {
			dashboard.leave();
			done();
		}
	})
	.on("/log", () => {
		app.sideBarToggle("/log")
		logPage.render(contentSelector);
	})
	.on("/settings", () => {
		app.sideBarToggle("/settings")
		settingsPage.render(contentSelector);
	}, {
		leave(done) {
			settings.leave();
			done();
		}
	})
	.on("/devices", () => {
		app.sideBarToggle("/devices")
		devicesPage.render(contentSelector);
	}, {
		leave(done) {
			devices.leave();
			done();
		}
	})
	.on("/aliases", () => {
		app.sideBarToggle("/aliases")
		aliasesPage.render(contentSelector);
	}, {
		leave(done) {
			aliases.leave();
			done();
		}
	}),
);
