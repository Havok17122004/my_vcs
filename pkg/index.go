package pkg

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"unsafe"
)

type Header struct {
	// signature     [4]byte  needed??
	// version       [4]byte  needed??
	numberEntries uint32
}

type Entry struct {
	Metadata Metadata
	id       string
	pathsize uint16
}

type Index struct {
	Header  Header
	Entries map[string]Entry
}

func ParseIndex() *Index {
	var idx Index
	file, err := os.OpenFile(filepath.Join(VCSDirPath, "index"), os.O_CREATE|os.O_RDWR, 0777)
	Check(err)

	fileInfo, _ := file.Stat()
	fileSize := fileInfo.Size()
	if fileSize == 0 {
		return &Index{
			Header:  Header{numberEntries: 0},
			Entries: make(map[string]Entry),
		}
	}

	_, err = file.Seek(-20, 2)
	Check(err)
	checksum := make([]byte, 20)
	_, err = file.Read(checksum)
	Check(err)

	bufferSize := fileSize - 20
	if bufferSize < 0 {
		fmt.Println("buffersize less than 0")
		return &Index{}
	}
	buffer := make([]byte, bufferSize)
	_, err = file.Seek(0, 0)
	Check(err)
	_, err = file.Read(buffer)
	Check(err)

	checksumStr := fmt.Sprintf("%x", checksum)
	sha := fmt.Sprintf("%x", GetSHAofText(string(buffer)))
	if checksumStr != sha {
		fmt.Println("sha of indexfile not matching")
		return &Index{}
	}

	_, err = file.Seek(0, 0)
	Check(err)

	var h Header
	// fmt.Println(unsafe.Sizeof(h))
	hSlice := make([]byte, unsafe.Sizeof(h))
	_, err = file.Read(hSlice)
	// fmt.Println(hSlice)
	Check(err)

	hBuf := bytes.NewBuffer(hSlice[:])
	err = binary.Read(hBuf, binary.BigEndian, &h.numberEntries)
	// fmt.Println(h.numberEntries)
	Check(err)

	idx.Header = h
	idx.Entries = make(map[string]Entry)
	for i := uint32(0); i < h.numberEntries; i++ {
		var e Entry
		metadataSlice := make([]byte, 28)
		entryIdSlice := make([]byte, 40)
		sizeSlice := make([]byte, 2)

		_, err = file.Read(metadataSlice)
		// fmt.Println(string(metadataSlice))
		Check(err)
		_, err = file.Read(entryIdSlice)
		// fmt.Printf("%x\n", entryIdSlice)
		// fmt.Println(string(entryIdSlice))
		Check(err)
		_, err = file.Read(sizeSlice)
		// fmt.Println(sizeSlice)
		Check(err)

		meta := BytesToMetadata(metadataSlice)
		// fmt.Println(meta)
		e.Metadata = meta
		e.id = string(entryIdSlice)

		sizeBuf := bytes.NewBuffer(sizeSlice)
		err = binary.Read(sizeBuf, binary.BigEndian, &e.pathsize)
		// fmt.Println(e.pathsize)
		Check(err)

		nameSlice := make([]byte, e.pathsize)
		_, err = file.Read(nameSlice)
		Check(err)
		Name := string(nameSlice)
		idx.Entries[Name] = e
	}
	// fmt.Println(idx.Header)
	// fmt.Println(idx.Entries)
	fmt.Println("Parsed index file as ", idx)
	return &idx

}

func (index *Index) ModifyIndex(path string, info fs.FileInfo, Oid string) {
	if index.Entries == nil {
		index.Entries = make(map[string]Entry)
	}
	meta, err := GetFileMetadata(info)
	// fmt.Println(meta)
	Check(err)

	// fmt.Println(Oid)
	Check(err)
	// fmt.Println(len(path))
	index.Entries[path] = Entry{Metadata: meta, id: Oid, pathsize: uint16(len(path))}
	// fmt.Printf("%s\n\n", path)
	index.Header.numberEntries = uint32(len(index.Entries))
	// fmt.Println(index)
	fmt.Println("Modified index and added ", index.Entries[path], " for the path ", path)
}

func (index *Index) SaveIndex() {

	file, err := os.OpenFile(filepath.Join(VCSDirPath, "index"), os.O_RDWR|os.O_CREATE, 0777)
	Check(err)
	defer file.Close()
	var buffer bytes.Buffer

	var h Header
	headerBytes := make([]byte, unsafe.Sizeof(h))
	binary.BigEndian.PutUint32(headerBytes, index.Header.numberEntries) //--------------------------------------------------------------------------------
	// fmt.Println(headerBytes)
	_, err = buffer.Write(headerBytes)
	Check(err)
	_, err = file.Write(headerBytes)
	Check(err)

	keys := make([]string, 0, len(index.Entries))
	for key := range index.Entries {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, key := range keys {
		entry := index.Entries[key]

		metadataBytes := MetadataToBytes(entry.Metadata)
		// fmt.Println(metadataBytes)
		_, err = buffer.Write(metadataBytes)
		Check(err)
		_, err = file.Write(metadataBytes)
		Check(err)

		_, err = buffer.Write([]byte(entry.id))
		Check(err)
		_, err = file.Write([]byte(entry.id))
		Check(err)

		sizeBytes := make([]byte, 2)
		binary.BigEndian.PutUint16(sizeBytes, uint16(entry.pathsize))
		// fmt.Println(sizeBytes)
		_, err = buffer.Write(sizeBytes)
		Check(err)
		_, err = file.Write(sizeBytes)
		Check(err)

		_, err = buffer.Write([]byte(key))
		Check(err)
		_, err = file.Write([]byte(key))
		Check(err)
	}

	checksum := GetSHAofText(buffer.String())
	// fmt.Println(buffer.String())
	_, err = file.Write(checksum[:])
	Check(err)
	fmt.Println("Index saved as ", index)
}
