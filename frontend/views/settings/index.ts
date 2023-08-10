import App from "../../system";
import {Interface, GPIO} from "./types";
import Reset from "./reset";
import GPIOTable from "./gpio";
import {Tab, Tabs} from "../../components/tabs";
import Datetime from "./datetime";
import Network from "./network";
import {$html} from "../index";
import UART from "./uart";
import SPI from "./spi";
import I2C from "./i2c";
import {NTP} from "../../types/datetime";
import {WiFiAp, WAN, DNS} from "../../types/network";

export default class Settings {
	private datetime: Datetime;
	private network: Network;
	private gpio: GPIOTable;

	private UARTs: Map<string, UART>;
	private SPIs: Map<string, SPI>;
	private I2Cs: Map<string, I2C>;
	private reset: Reset;

	constructor(private readonly app: App) {
		this.UARTs = new Map<string, UART>();
		this.SPIs = new Map<string, SPI>();
		this.I2Cs = new Map<string, I2C>();
		this.reset = new Reset(app);
	}

	public async render(selector: string) {
		$html(selector,
			`<div class="row">
				<div class="col-12">
					<div class="card card-primary card-outline card-outline-tabs">
						<div class="card-header p-0 border-bottom-0" id="tabs-container"></div>
						<div class="card-body" id="content-container"></div>
						<div class="card-footer">
							<div id="reset-form"></div>
						</div>
					</div>
				</div>
			</div>`
		);

		await this.launch();
	}

	private async launch(): Promise<void> {
		try {
			await this.loadAll();
			await this.handleDynamicContent();
		} catch (err) {
			this.app.error(err);
		}
	}

	private async loadAll(): Promise<void> {
		const interfaces = await this.app.call("system.getInterfaces") as any as Interface[];

		interfaces.forEach(iface => {
			const id = iface.type + "-" + iface.number;

			switch (iface.type) {
				case "UART": return this.UARTs.set(id, new UART(this.app, iface.number, iface.name));
				case "SPI":  return this.SPIs.set(id, new SPI(this.app, iface.number, iface.name));
				case "I2C":  return this.I2Cs.set(id, new I2C(this.app, iface.number, iface.name));
			}
		})

		this.datetime = new Datetime(
			this.app,
			await this.app.call("datetime.get") as any as string,
			await this.app.call("datetime.getNTP") as any as NTP,
		);

		this.network = new Network(
			this.app,
			await this.app.call("network.getWiFiAp") as any as WiFiAp,
			await this.app.call("network.getWAN") as any as WAN,
			await this.app.call("network.getDNS") as any as DNS,
		);

		this.gpio = new GPIOTable(
			this.app,
			await this.app.call("gpio.get") as any as GPIO[]
		);
	}

	private async handleDynamicContent() {
		await this.renderDynamicContent();
		this.visit();
		this.handleEvents();
	}

	private async renderDynamicContent() {
		const staticTabs: Tab[]  = [
			new Tab("datetime", "Дата и Время", this.datetime),
			new Tab("network", "Сеть", this.network),
			new Tab("gpio-table", "Входы и Выходы", this.gpio),
		];

		const ifaceTabs: Tab[] = [];
		this.UARTs.forEach((uart, id) => ifaceTabs.push(new Tab(id, uart.name, uart)));
		this.SPIs.forEach((spi, id) => ifaceTabs.push(new Tab(id, spi.name, spi)));
		this.I2Cs.forEach((i2c, id) => ifaceTabs.push(new Tab(id, i2c.name, i2c)));

		const tabs = new Tabs("settings", [...staticTabs, ...ifaceTabs]);
		await tabs.render("#tabs-container", "#content-container");

		this.reset.render("#reset-form");
	}

	private visit(): void {
		this.datetime.visit();
		this.network.visit();
		this.gpio.visit();
		this.UARTs.forEach(uart => uart.visit());
		this.SPIs.forEach(spi => spi.visit());
		this.I2Cs.forEach(i2c => i2c.visit());
	}

	private handleEvents(): void {
		this.datetime.handleEvents();
		this.network.handleEvents();
		this.gpio.handleEvents();
		this.UARTs.forEach(uart => uart.handleEvents());
		this.SPIs.forEach(spi => spi.handleEvents());
		this.I2Cs.forEach(i2c => i2c.handleEvents());
	}

	public leave(): void {
		this.datetime?.leave();
		this.network?.leave();
		this.gpio?.leave();
		this.UARTs.forEach(uart => uart?.leave());
		this.SPIs.forEach(spi => spi?.leave());
		this.I2Cs.forEach(i2c => i2c?.leave());
	}
}
