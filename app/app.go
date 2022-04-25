package app

var (
	relayer *relayer.Relayer
)

func Relayer() relayer.Relayer {
	if relayer == nil {
		// do something
	}
	return relayer
}
