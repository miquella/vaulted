package vaulted

type Operation int

const (
	OpenOperation Operation = iota
	SealOperation
)

type Steward interface {
	GetPassword(operation Operation, name string) (string, error)
}

type StewardMaxTries interface {
	GetMaxOpenTries() int
}

type StaticSteward struct {
	Password string
}

func NewStaticSteward(password string) *StaticSteward {
	return &StaticSteward{
		Password: password,
	}
}

func (s *StaticSteward) GetPassword(operation Operation, name string) (string, error) {
	return s.Password, nil
}
