import {$$, $html, $onClick, $setSelect} from "../../index";
import {WiFiAp, WiFiCl, WAN, DNS} from "../../../types/network";
import App from "../../../system";
import {clearForm, ErrorCode, setInvalidFeedback} from "../../../components/errors";

// todo: разбить на пакеты!
export default class Network {
	private wifiApForm: {
		enabled: HTMLInputElement
		ssid: HTMLInputElement,
		password: HTMLInputElement,
		channel: HTMLSelectElement,
		ttl: HTMLSelectElement,
		encryption: HTMLSelectElement,
		ip: HTMLInputElement,
		mask: HTMLInputElement,
		dhcp: HTMLInputElement,
		apply: HTMLButtonElement,
		feedbacks: HTMLDivElement[],
	};

	private wifiClForm: {
		enabled: HTMLInputElement
		ssid: HTMLInputElement,
		password: HTMLInputElement,
		ip: HTMLInputElement,
		mask: HTMLInputElement,
		gateway: HTMLInputElement,
		dhcp: HTMLInputElement,
		apply: HTMLButtonElement,
		feedbacks: HTMLDivElement[],
	};

	private wanForm: {
		dhcp: HTMLInputElement,
		ip: HTMLInputElement,
		mask: HTMLInputElement,
		gateway: HTMLInputElement,
		apply: HTMLButtonElement,
		feedbacks: HTMLDivElement[],
	};

	private dnsForm: {
		dns1: HTMLInputElement,
		dns2: HTMLInputElement,
		apply: HTMLButtonElement,
		feedbacks: HTMLDivElement[],
	};

	private subscriptions = new Map<string, (data) => void>([
		['network.setWiFiAp', this.loggedSetWiFiAp.bind(this)],
		['network.setWiFiCl', this.loggedSetWiFiCl.bind(this)],
		['network.setWAN',    this.loggedSetWAN.bind(this)],
		['network.setDNS',    this.loggedSetDNS.bind(this)],
	]);

	constructor(
		private readonly app: App,
		private wifiAp: WiFiAp,
		// private wifiCl: WiFiCl,
		private wan: WAN,
		private dns: DNS,
	) {}

	public render(selector: string): void {
		$html(selector, `
			<div class="row">
				${Network.renderWiFiAp()}
				${Network.renderWiFiCl()}
				${Network.renderWAN()}
				${Network.renderDNS()}
			</div>`
		);

		this.wifiApForm = {
			enabled:    $$('#wifi-ap-enabled') as HTMLInputElement,
			ssid:       $$('#wifi-ap-ssid') as HTMLInputElement,
			password:   $$('#wifi-ap-password') as HTMLInputElement,
			channel:    $$('#wifi-ap-channel') as HTMLSelectElement,
			ttl:        $$('#wifi-ap-ttl') as HTMLSelectElement,
			encryption: $$('#wifi-ap-enc-type') as HTMLSelectElement,
			ip:         $$('#wifi-ap-ip') as HTMLInputElement,
			mask:       $$('#wifi-ap-mask') as HTMLInputElement,
			dhcp:       $$('#wifi-ap-dhcp') as HTMLInputElement,
			apply:      $$('#wifi-ap-apply') as HTMLButtonElement,
			feedbacks:  [],
		};

		this.wifiClForm = {
			enabled:    $$('#wifi-cl-enabled') as HTMLInputElement,
			ssid:       $$('#wifi-cl-ssid') as HTMLInputElement,
			password:   $$('#wifi-cl-passwd') as HTMLInputElement,
			ip:         $$('#wifi-cl-ip') as HTMLInputElement,
			mask:       $$('#wifi-cl-mask') as HTMLInputElement,
			gateway:    $$('#wifi-cl-gateway') as HTMLInputElement,
			dhcp:       $$('#wifi-cl-dhcp') as HTMLInputElement,
			apply:      $$('#wifi-cl-apply') as HTMLButtonElement,
			feedbacks:  [],
		};

		this.wanForm = {
			dhcp:       $$('#wan-dhcp') as HTMLInputElement,
			ip:         $$('#wan-ip') as HTMLInputElement,
			mask:       $$('#wan-mask') as HTMLInputElement,
			gateway:    $$('#wan-gateway') as HTMLInputElement,
			apply:      $$('#wan-apply') as HTMLButtonElement,
			feedbacks:  [],
		};

		this.dnsForm = {
			dns1:      $$('#dns-1') as HTMLInputElement,
			dns2:      $$('#dns-2') as HTMLInputElement,
			apply:     $$('#dns-apply') as HTMLButtonElement,
			feedbacks: [],
		};

		this.setWiFiAp(this.wifiAp);
		// this.setWiFiCl(this.wifiCl);
		this.setWAN(this.wan);
		this.setDNS(this.dns);
	}

	private static renderWiFiAp(): string {
		const ttls: string[] = ["5", "10", "15"];
		const channels: string[] = ["auto", "1", "2", "3", "4", "5", "6", "7", "8", "9", "10"];
		const encryptions = { psk: "WPA (PSK)", psk2: "WPA2 (PSK)" };

		return (
			`<div class="col-sm-12 col-md-12 col-lg-6">
				<div class="card card-outline">
					<div class="card-header">
						<h3 class="card-title">Wi-Fi Сервер</h3>
						<div class="card-tools m-0">
							<div class="custom-control custom-switch">
							  <input type="checkbox" class="custom-control-input" id="wifi-ap-enabled">
							  <label class="custom-control-label" for="wifi-ap-enabled"></label>
							</div>
						</div>
					</div>
					<div class="card-body">
						<div class="form-group row">
							<label for="wifi-ap-ssid" class="col-sm-4 col-form-label">SSID</label>
							<div class="col-sm-8">
								<input type="text" id="wifi-ap-ssid" class="form-control">
							</div>
						</div>
						<div class="form-group row">
							<label for="wifi-ap-password" class="col-sm-4 col-form-label">Пароль</label>
							<div class="col-sm-8">
								<input type="password" id="wifi-ap-password" class="form-control">
							</div>
						</div>
						<div class="form-group row">
							<label for="wifi-ap-channel" class="col-sm-4 col-form-label">Канал</label>
							<div class="col-sm-8">
								<select id="wifi-ap-channel" class="custom-select">
									${channels.map(ch => `<option>${ch}</option>`).join('')}
								</select>
							</div>
						</div>
						<div class="form-group row">
							<label for="wifi-ap-ttl" class="col-sm-4 col-form-label">TTL (min)</label>
							<div class="col-sm-8">
								<select id="wifi-ap-ttl" class="custom-select">
									${ttls.map(ttl => `<option>${ttl}</option>`).join('')}
								</select>
							</div>
						</div>
						<div class="form-group row">
							<label for="wifi-ap-enc-type" class="col-sm-4 col-form-label">Тип шифрования</label>
							<div class="col-sm-8">
								<select id="wifi-ap-enc-type" class="custom-select">
									${Object.keys(encryptions).map(key => `<option value="${key}">${encryptions[key]}</option>`).join('')}
								</select>
							</div>
						</div>
						<div class="form-group row">
							<label for="wifi-ap-ip" class="col-sm-4 col-form-label">IP Адрес</label>
							<div class="col-sm-8">
								<input type="text" id="wifi-ap-ip" class="form-control">
							</div>
						</div>
						<div class="form-group row">
							<label for="wifi-ap-mask" class="col-sm-4 col-form-label">Маска сети</label>
							<div class="col-sm-8">
								<input type="text" id="wifi-ap-mask" class="form-control">
							</div>
						</div>
						<div class="form-group row">
							<label for="date" class="col-sm-4 col-form-label"></label>
							<div class="col-sm-8">
								<div class="custom-control custom-switch">
								  <input type="checkbox" class="custom-control-input" id="wifi-ap-dhcp">
								  <label class="custom-control-label" for="wifi-ap-dhcp">DHCP Сервер</label>
								</div>
							</div>
						</div>
					</div>
					<div class="card-footer">
						<button type="button" id="wifi-ap-apply" class="btn btn-success float-right">Применить</button>
					</div>
				</div>
			</div>`
		);
	}

	private static renderWiFiCl(): string {
		return (
			`<div class="col-sm-12 col-md-12 col-lg-6">
				<div class="card">
					<div class="card-header">
						<h3 class="card-title">Wi-Fi Клиент</h3>
						<div class="card-tools m-0">
							<div class="custom-control custom-switch">
								<input type="checkbox" class="custom-control-input" id="wifi-cl-enabled" checked>
								<label class="custom-control-label" for="wifi-cl-enabled"></label>
							</div>
						</div>
					</div>
					<div class="card-body">
						<div class="form-group row">
							<label for="wifi-cl-ssid" class="col-sm-4 col-form-label">SSID</label>
							<div class="col-sm-8">
								<input type="text" id="wifi-cl-ssid" class="form-control">
							</div>
						</div>
						<div class="form-group row">
							<label for="wifi-cl-passwd" class="col-sm-4 col-form-label">Пароль</label>
							<div class="col-sm-8">
								<input type="password" id="wifi-cl-passwd" class="form-control">
							</div>
						</div>
						<div class="form-group row">
							<label for="date" class="col-sm-4 col-form-label"></label>
							<div class="col-sm-8">
								<div class="custom-control custom-switch">
								  <input type="checkbox" class="custom-control-input" id="wifi-cl-dhcp">
								  <label class="custom-control-label" for="wifi-cl-dhcp">DHCP</label>
								</div>
							</div>
						</div>
						<div class="form-group row">
							<label for="wifi-cl-ip" class="col-sm-4 col-form-label">IP Адрес</label>
							<div class="col-sm-8">
								<input type="text" id="wifi-cl-ip" class="form-control">
							</div>
						</div>
						<div class="form-group row">
							<label for="wifi-cl-mask" class="col-sm-4 col-form-label">Маска сети</label>
							<div class="col-sm-8">
								<input type="text" id="wifi-cl-mask" class="form-control">
							</div>
						</div>
						<div class="form-group row">
							<label for="wifi-cl-gateway" class="col-sm-4 col-form-label">Шлюз по-умолчанию</label>
							<div class="col-sm-8">
								<input type="text" id="wifi-cl-gateway" class="form-control">
							</div>
						</div>
					</div>
					<div class="card-footer">
						<button type="button" id="wifi-cl-apply" class="btn btn-success float-right">Применить</button>
					</div>
				</div>
			</div>`
		);
	}

	private static renderWAN(): string {
		return (
			`<div class="col-sm-12 col-md-12 col-lg-6">
				<div class="card">
					<div class="card-header">
						<h3 class="card-title">Конфигурация WAN (Ethernet)</h3>
					</div>
					<div class="card-body">
						<div class="form-group row">
							<label for="date" class="col-sm-4 col-form-label"></label>
							<div class="col-sm-8">
								<div class="custom-control custom-switch">
								  <input type="checkbox" class="custom-control-input" id="wan-dhcp">
								  <label class="custom-control-label" for="wan-dhcp">DHCP</label>
								</div>
							</div>
						</div>
						<div class="form-group row">
							<label for="wan-ip" class="col-sm-4 col-form-label">IP Адрес</label>
							<div class="col-sm-8">
								<input type="text" id="wan-ip" class="form-control">
							</div>
						</div>
						<div class="form-group row">
							<label for="wan-mask" class="col-sm-4 col-form-label">Маска сети</label>
							<div class="col-sm-8">
								<input type="text" id="wan-mask" class="form-control">
							</div>
						</div>
						<div class="form-group row">
							<label for="wan-gateway" class="col-sm-4 col-form-label">Шлюз по-умолчанию</label>
							<div class="col-sm-8">
								<input type="text" id="wan-gateway" class="form-control">
							</div>
						</div>
					</div>
					<div class="card-footer">
						<button type="button" id="wan-apply" class="btn btn-success float-right">Применить</button>
					</div>
				</div>
			</div>`
		);
	}

	private static renderDNS(): string {
		return (
			`<div class="col-sm-12 col-md-12 col-lg-6">
				<div class="card">
					<div class="card-header">
						<h3 class="card-title">Конфигурация DNS</h3>
					</div>
					<div class="card-body">
						<div class="form-group row">
							<label for="dns-1" class="col-sm-4 col-form-label">DNS 1</label>
							<div class="col-sm-8">
								<input type="text" id="dns-1" class="form-control">
							</div>
						</div>
						<div class="form-group row">
							<label for="dns-2" class="col-sm-4 col-form-label">DNS 2</label>
							<div class="col-sm-8">
								<input type="text" id="dns-2" class="form-control">
							</div>
						</div>
					</div>
					<div class="card-footer">
						<button type="button" id="dns-apply" class="btn btn-success float-right">Применить</button>
					</div>
				</div>
			</div>`
		);
	}

	public visit(): void {
		for (const [method, handler] of this.subscriptions) {
			this.app.subscribe(method, handler);
		}
	}

	public leave(): void {
		for (const [method, handler] of this.subscriptions) {
			this.app.unsubscribe(method, handler);
		}
	}

	public handleEvents(): void {
		$onClick(this.wifiApForm.apply, () => this.app.call('network.setWiFiAp', {
			enabled:    this.wifiApForm.enabled.checked,
			ssid:       this.wifiApForm.ssid.value,
			password:   this.wifiApForm.password.value,
			channel:    this.wifiApForm.channel.value,
			ttl:        Number(this.wifiApForm.ttl.value),
			encryption: this.wifiApForm.encryption.value,
			ip:         this.wifiApForm.ip.value,
			mask:       this.wifiApForm.mask.value,
			dhcp:       this.wifiApForm.dhcp.checked,
		}).catch(err => {
			if (err.code === ErrorCode.ValidationErr && err.data) {
				setInvalidFeedback(this.wifiApForm, err.data);
				return
			}

			this.app.error(err);
		}));

		$onClick(this.wifiClForm.apply, () => this.app.call('network.setWiFiCl', {
			enabled:  this.wifiClForm.enabled.checked,
			ssid:     this.wifiClForm.ssid.value,
			password: this.wifiClForm.password.value,
			ip:       this.wifiClForm.ip.value,
			mask:     this.wifiClForm.mask.value,
			gateway:  this.wifiClForm.gateway.value,
			dhcp:     this.wifiClForm.dhcp.checked,
		}).catch(err => {
			if (err.code === ErrorCode.ValidationErr && err.data) {
				setInvalidFeedback(this.wifiClForm, err.data);
				return
			}

			this.app.error(err);
		}));

		$onClick(this.wanForm.apply, () => this.app.call('network.setWAN', {
			dhcp:    this.wanForm.dhcp.checked,
			ip:      this.wanForm.ip.value,
			mask:    this.wanForm.mask.value,
			gateway: this.wanForm.gateway.value,
		}).catch(err => {
			if (err.code === ErrorCode.ValidationErr && err.data) {
				setInvalidFeedback(this.wanForm, err.data);
				return
			}

			this.app.error(err);
		}));

		$onClick(this.dnsForm.apply, () => this.app.call('network.setDNS', {
			dns1: this.dnsForm.dns1.value,
			dns2: this.dnsForm.dns2.value,
		}).catch(err => {
			if (err.code === ErrorCode.ValidationErr && err.data) {
				setInvalidFeedback(this.dnsForm, err.data);
				return
			}

			this.app.error(err);
		}));
	}

	private loggedSetWiFiAp(req: WiFiAp): void {
		this.setWiFiAp(req);
		this.app.info('Настройки точки доступа Wi-Fi применены');
	}

	private setWiFiAp(req: WiFiAp): void {
		clearForm(this.wifiApForm);

		this.wifiApForm.enabled.checked = req.enabled;
		this.wifiApForm.ssid.value = req.ssid;
		this.wifiApForm.password.value = req.password;
		$setSelect(this.wifiApForm.channel, req.channel);
		$setSelect(this.wifiApForm.ttl, req.ttl);
		$setSelect(this.wifiApForm.encryption, req.encryption);
		this.wifiApForm.ip.value = req.ip;
		this.wifiApForm.mask.value = req.mask;
		this.wifiApForm.dhcp.checked = req.dhcp;
	}

	private loggedSetWiFiCl(req: WiFiCl): void {
		this.setWiFiCl(req);
		this.app.info('Настройки Wi-Fi клиента применены');
	}

	private setWiFiCl(req: WiFiCl): void {
		clearForm(this.wifiClForm);

		this.wifiClForm.enabled.checked = req.enabled;
		this.wifiClForm.ssid.value = req.ssid;
		this.wifiClForm.password.value = req.password;
		this.wifiClForm.dhcp.checked = req.dhcp;
		this.wifiClForm.ip.value = req.ip;
		this.wifiClForm.mask.value = req.mask;
		this.wifiClForm.gateway.value = req.gateway;

		const disabled = req.dhcp;
		this.wifiClForm.ip.disabled = disabled;
		this.wifiClForm.mask.disabled = disabled;
		this.wifiClForm.gateway.disabled = disabled;
		this.dnsForm.dns1.disabled = disabled;
		this.dnsForm.dns2.disabled = disabled;
		this.dnsForm.apply.disabled = disabled;
	}

	private loggedSetWAN(req: WAN): void {
		this.setWAN(req);
		this.app.info('Настройки WAN применены');
	}

	private setWAN(req: WAN): void {
		clearForm(this.wanForm);

		this.wanForm.ip.value = req.ip;
		this.wanForm.mask.value = req.mask;
		this.wanForm.gateway.value = req.gateway;
		this.wanForm.dhcp.checked = req.dhcp;

		const disabled = req.dhcp;
		this.wanForm.ip.disabled = disabled;
		this.wanForm.mask.disabled = disabled;
		this.wanForm.gateway.disabled = disabled;
		this.dnsForm.dns1.disabled = disabled;
		this.dnsForm.dns2.disabled = disabled;
		this.dnsForm.apply.disabled = disabled;
	}

	private loggedSetDNS(req: DNS): void {
		this.setDNS(req);
		this.app.info('Настройки DNS применены');
	}

	private setDNS(req: DNS): void {
		clearForm(this.dnsForm);

		this.dnsForm.dns1.value = req.dns1;
		this.dnsForm.dns2.value = req.dns2 ?? "";
	}
}
