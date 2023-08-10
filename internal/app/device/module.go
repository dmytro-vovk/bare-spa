package device

import (
	"fmt"

	"github.com/Sergii-Kirichok/pr/internal/app/device/module"
)

func (kind Type) Modules() map[module.Type]struct{} {
	return map[Type]map[module.Type]struct{}{
		TypeEthIO2: {
			module.TypeTriac2: {},
			module.TypeTriac4: {},
			module.TypeTriac6: {},
			module.TypeRelay2: {},
			module.TypeRelay4: {},
			module.TypeRelay6: {},
		},
		TypeEthIO4: {
			module.TypeTriac2: {},
			module.TypeTriac4: {},
			module.TypeTriac6: {},
			module.TypeRelay2: {},
			module.TypeRelay4: {},
			module.TypeRelay6: {},
		},
		TypeIncubator1: {},
		TypeIncubator2: {
			module.TypeTriac2: {},
			module.TypeTriac4: {},
			module.TypeTriac6: {},
			module.TypeRelay2: {},
			module.TypeRelay4: {},
			module.TypeRelay6: {},
		},
		TypeOmega1: {
			module.TypeTriac2: {},
			module.TypeTriac4: {},
			module.TypeTriac6: {},
			module.TypeRelay2: {},
			module.TypeRelay4: {},
			module.TypeRelay6: {},
		},
		TypeUPS12: {},
		TypeEmbedded: {
			module.TypeTriac2: {},
			module.TypeRelay2: {},
		},
	}[kind]
}

type ErrModuleNotSupported struct {
	DeviceType Type
	ModuleType module.Type
}

func (e *ErrModuleNotSupported) Error() string {
	return fmt.Sprintf("module %q not supported for device type %q", e.ModuleType, e.DeviceType)
}

func (kind Type) IsModuleSupported(mod module.Type) bool {
	_, ok := kind.Modules()[mod]
	return ok
}
