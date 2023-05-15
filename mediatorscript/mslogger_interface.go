package mediatorscript

type MsLogger interface {
	Infof(format string, a ...any)
	Warningf(format string, a ...any)
}
