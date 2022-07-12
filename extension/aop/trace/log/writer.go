package log

type Writer interface {
	Write(p []byte)
}
