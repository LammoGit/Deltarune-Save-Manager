package saves

//go:generate stringer -type=Phone -trimprefix=Phone
type Phone int

const PhoneEmpty Phone = 0

const (
	PhoneHome Phone = iota + 201
	PhoneSans
)
