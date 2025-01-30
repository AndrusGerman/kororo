package types

type Database string

func (ctx Database) String() string {
	return string(ctx)
}
