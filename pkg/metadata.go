package pkg

import (
	"bytes"
	"encoding/binary"
	"io/fs"
	"os"
	"syscall"
	"time"
)

type Metadata struct {
	Ctime time.Time
	Mtime time.Time
	Mode  uint32
	Size  int64
}

func getUnixMetadata(fileInfo os.FileInfo) (Metadata, error) {
	stat, ok := fileInfo.Sys().(*syscall.Stat_t)
	if !ok {
		return Metadata{}, syscall.EINVAL
	}

	return Metadata{
		Ctime: time.Unix(int64(stat.Ctim.Sec), int64(stat.Ctim.Nsec)),
		Mtime: time.Unix(int64(stat.Mtim.Sec), int64(stat.Mtim.Nsec)),
		Mode:  uint32(stat.Mode),
		Size:  stat.Size,
	}, nil
}

func GetFileMetadata(fileInfo fs.FileInfo) (Metadata, error) {

	var metadata Metadata

	// switch runtime.GOOS {

	// default:
	metadata, err := getUnixMetadata(fileInfo)
	Check(err)
	// }

	return metadata, nil
}

func MetadataToBytes(metadata Metadata) []byte {
	var buf bytes.Buffer

	ctimeBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(ctimeBytes, uint64(metadata.Ctime.UnixNano()))
	buf.Write(ctimeBytes)

	mtimeBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(mtimeBytes, uint64(metadata.Mtime.UnixNano()))
	buf.Write(mtimeBytes)

	modeBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(modeBytes, metadata.Mode)
	buf.Write(modeBytes)

	sizeBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(sizeBytes, uint64(metadata.Size))
	buf.Write(sizeBytes)

	return buf.Bytes()
}

func BytesToMetadata(data []byte) Metadata {
	var ctime int64
	var mtime int64
	var mode uint32
	var size int64

	buf := bytes.NewReader(data)

	err := binary.Read(buf, binary.LittleEndian, &ctime)
	Check(err)
	err = binary.Read(buf, binary.LittleEndian, &mtime)
	Check(err)
	err = binary.Read(buf, binary.LittleEndian, &mode)
	Check(err)
	err = binary.Read(buf, binary.LittleEndian, &size)
	Check(err)

	return Metadata{
		Ctime: time.Unix(0, ctime),
		Mtime: time.Unix(0, mtime),
		Mode:  mode,
		Size:  size,
	}
}
