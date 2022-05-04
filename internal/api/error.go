package api

type Error struct {
	Cause string
}

func (this Error) Error() string {
	return this.Cause
}
