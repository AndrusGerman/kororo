package types

type Backend string

func (ctx Backend) String() string {
	return string(ctx)
}
