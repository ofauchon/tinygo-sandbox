//go:build !sx126x

package rfswitch

type CustomSwitch struct {
}

var (
	rfstate int
)

func (s CustomSwitch) InitRFSwitch() {
}

func (s CustomSwitch) SetRfSwitchMode(mode int) error {
	rfstate = mode

	return nil
}
