package p9p
import "io"
type Msg interface {
	ReadBinary(r io.Reader) (err error)
	WriteBinary(w io.Writer) (err error)
}
