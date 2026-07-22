package saves

//go:generate stringer -type=Phone -trimprefix=Phone
type Phone int

const (
	PhoneEmpty Phone = iota
	PhoneHome        = iota + 200
	PhoneSans
)
