package autodoc

// 兼容低版本 Golang

// A stringBuilder is used to efficiently build a string using Write methods.
// It minimizes memory copying. The zero value is ready to use.
// Do not copy a non-zero Builder.
type stringBuilder struct {
	buf []byte
}

// WriteString appends the contents of s to b's buffer.
func (sb *stringBuilder) WriteString(s string) {
	sb.buf = append(sb.buf, s...)
}

// String returns the accumulated string.
func (sb *stringBuilder) String() string {
	return string(sb.buf)
	//return *(*string)(unsafe.Pointer(&sb.buf))
}
