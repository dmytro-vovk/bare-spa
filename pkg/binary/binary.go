package binary

type Bits uint32

func GetBits(n uint) Bits          { return 1 << n }
func SetBits(b, flag Bits) Bits    { return b | flag }
func ClearBits(b, flag Bits) Bits  { return b &^ flag }
func ToggleBits(b, flag Bits) Bits { return b ^ flag }
func HasBits(b, flag Bits) bool    { return b&flag == flag }

func ClearFlag(bitPos, numBits uint) (b Bits) {
	for i := uint(0); i < numBits; i++ {
		b |= GetBits(bitPos + i)
	}

	return
}

func SetFromTo(from, to uint, b, flag Bits) Bits {
	// todo: include additional checks for border situations
	//if from > to {
	//	return 0, errors.New(`"from" can't be greater than "to"`)
	//}

	cf := ClearFlag(from, to-from+1) // include left value
	return SetBits(ClearBits(b, cf), flag<<from)
}
