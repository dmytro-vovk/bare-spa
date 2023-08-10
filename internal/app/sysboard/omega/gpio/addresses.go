//go:build sysboard

package gpio

const (
	baseAddr = 0x1000_0000
	baseSize = 0x1000
)

// Nums at the end describe what GPIO{XX} is it
const (
	modeOffset01 = 0x0060
	modeOffset02 = 0x0064
)

const (
	ctrl0Offset = 0x0600
	ctrl1Offset = 0x0604
	ctrl2Offset = 0x0608
)

// For inversion settings(polarity)
// rights: RW
const (
	pol0Offset = 0x0610
	pol1Offset = 0x0614
	pol2Offset = 0x0618
)

// For setting whole register data
// rights: RW
const (
	data0Offset = 0x0620
	data1Offset = 0x0624
	data2Offset = 0x0628
)

// 1 for set the register; 0 has no effect
// rights: WO
const (
	dset0Offset = 0x0630
	dset1Offset = 0x0634
	dset2Offset = 0x0638
)

// 1 for set the register; 0 has no effect
// rights: WO
const (
	dclr0Offset = 0x0640
	dclr1Offset = 0x0644
	dclr2Offset = 0x0648
)
