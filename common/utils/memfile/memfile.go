package memfile

import (
	"errors"
	"io"
	"io/fs"
	"sync"
	"time"
)

var errInvalid = errors.New("invalid argument")

// File is an in-memory emulation of the I/O operations of os.File.
// The zero value for File is an empty file ready to use.
type File struct {
	m sync.Mutex
	b []byte
	i int

	closed bool

	filename string
}

type fileInfo struct {
	size int64
	name string
}

func (f *fileInfo) Name() string {
	if f.name == "" {
		return "memfile"
	}
	return f.name
}

func (f *fileInfo) Size() int64 {
	return f.size
}

func (f *fileInfo) Mode() fs.FileMode {
	return 0666
}

func (f *fileInfo) ModTime() time.Time {
	return time.Now()
}

func (f *fileInfo) IsDir() bool {
	return false // 由于是文件，返回 false
}

func (f *fileInfo) Sys() any {
	return nil
}

var _ fs.FileInfo = (*fileInfo)(nil)

// Stat returns the FileInfo structure describing file.
func (fb *File) Stat() (fs.FileInfo, error) {
	return &fileInfo{name: fb.filename, size: int64(len(fb.b))}, nil
}

// New creates and initializes a new File using b as its initial contents.
// The new File takes ownership of b.
func New(b []byte) *File {
	return &File{b: b}
}

func NewWithName(i string, b []byte) *File {
	return &File{b: b, filename: i}
}

// Name returns the name of the file.
func (fb *File) Name() string {
	return fb.filename
}

// Read reads up to len(b) bytes from the File.
// It returns the number of bytes read and any error encountered.
// At end of file, Read returns (0, io.EOF).
func (fb *File) Read(b []byte) (int, error) {
	fb.m.Lock()
	defer fb.m.Unlock()

	n, err := fb.readAt(b, int64(fb.i))
	fb.i += n
	return n, err
}

// SetName sets the name of the file.
func (fb *File) SetName(name string) {
	fb.filename = name
}

// ReadAt reads len(b) bytes from the File starting at byte offset.
// It returns the number of bytes read and the error, if any.
// At end of file, that error is io.EOF.
func (fb *File) ReadAt(b []byte, offset int64) (int, error) {
	fb.m.Lock()
	defer fb.m.Unlock()
	return fb.readAt(b, offset)
}

var errClosed = errors.New("memfile: read/write on closed file")

func (fb *File) readAt(b []byte, off int64) (int, error) {
	if fb.closed {
		return 0, errClosed
	}

	if off < 0 || int64(int(off)) < off {
		return 0, errInvalid
	}
	if off > int64(len(fb.b)) {
		return 0, io.EOF
	}
	n := copy(b, fb.b[off:])
	if n < len(b) {
		return n, io.EOF
	}
	return n, nil
}

// Write writes len(b) bytes to the File.
// It returns the number of bytes written and an error, if any.
// If the current file offset is past the io.EOF, then the space in-between are
// implicitly filled with zero bytes.
func (fb *File) Write(b []byte) (int, error) {
	fb.m.Lock()
	defer fb.m.Unlock()

	n, err := fb.writeAt(b, int64(fb.i))
	fb.i += n
	return n, err
}

// WriteAt writes len(b) bytes to the File starting at byte offset.
// It returns the number of bytes written and an error, if any.
// If offset lies past io.EOF, then the space in-between are implicitly filled
// with zero bytes.
func (fb *File) WriteAt(b []byte, offset int64) (int, error) {
	fb.m.Lock()
	defer fb.m.Unlock()
	return fb.writeAt(b, offset)
}
func (fb *File) writeAt(b []byte, off int64) (int, error) {
	if fb.closed {
		return 0, errClosed
	}

	if off < 0 || int64(int(off)) < off {
		return 0, errInvalid
	}
	if off > int64(len(fb.b)) {
		fb.truncate(off)
	}
	n := copy(fb.b[off:], b)
	fb.b = append(fb.b, b[n:]...)
	return len(b), nil
}

// Seek sets the offset for the next Read or Write on file with offset,
// interpreted according to whence: 0 means relative to the origin of the file,
// 1 means relative to the current offset, and 2 means relative to the end.
func (fb *File) Seek(offset int64, whence int) (int64, error) {
	fb.m.Lock()
	defer fb.m.Unlock()

	if fb.closed {
		return 0, errClosed
	}

	var abs int64
	switch whence {
	case io.SeekStart:
		abs = offset
	case io.SeekCurrent:
		abs = int64(fb.i) + offset
	case io.SeekEnd:
		abs = int64(len(fb.b)) + offset
	default:
		return 0, errInvalid
	}
	if abs < 0 {
		return 0, errInvalid
	}
	fb.i = int(abs)
	return abs, nil
}

// Truncate changes the size of the file. It does not change the I/O offset.
func (fb *File) Truncate(n int64) error {
	fb.m.Lock()
	defer fb.m.Unlock()
	return fb.truncate(n)
}
func (fb *File) truncate(n int64) error {
	if fb.closed {
		return errClosed
	}

	switch {
	case n < 0 || int64(int(n)) < n:
		return errInvalid
	case n <= int64(len(fb.b)):
		fb.b = fb.b[:n]
		return nil
	default:
		fb.b = append(fb.b, make([]byte, int(n)-len(fb.b))...)
		return nil
	}
}

// Bytes returns the full contents of the File.
// The result in only valid until the next Write, WriteAt, or Truncate call.
func (fb *File) Bytes() []byte {
	fb.m.Lock()
	defer fb.m.Unlock()
	return fb.b
}

func (fb *File) Close() error {
	fb.m.Lock()
	defer fb.m.Unlock()
	fb.closed = true
	return nil
}
