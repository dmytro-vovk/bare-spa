export type WiFiAp = {
	enabled: boolean,
	ssid: string,
	password: string,
	channel: string,
	ttl: number,
	encryption: string,
	ip: string,
	mask: string,
	dhcp: boolean,
};

export type WiFiCl = {
	enabled: boolean,
	ssid: string,
	password: string,
	ip: string,
	mask: string,
	gateway: string,
	dhcp: boolean,
};

export type WAN = {
	dhcp: boolean,
	ip: string,
	mask: string,
	gateway: string,
};

export type DNS = {
	dns1: string,
	dns2?: string,
};
