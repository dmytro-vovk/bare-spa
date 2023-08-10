export const btoa = (cond: boolean): string => cond ? 'true' : 'false';

export const markAs = (cond: boolean, as: string, ...other: string[]): string => cond ? ` ${[as, ...other].join(' ')}` : '';

export const markChecked = (cond: boolean): string => markAs(cond, ' checked');

export const markDisabled = (cond: boolean): string => markAs(cond, ' disabled');

export const markSelected = (cond: boolean): string => markAs(cond, ' selected');