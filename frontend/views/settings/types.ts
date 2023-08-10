export type GPIO = {
	pin: number;
	name: string;
	register: number;
	level: boolean;
	direction: boolean;
	inversion: boolean;
	type: string;
	default_state: number;
}

export type Interface = {
	name: string;
	type: string;
	number: number;
}
