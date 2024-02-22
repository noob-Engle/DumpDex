package cmd

import (
    "fmt"
    "io/ioutil"
    "regexp"
    "unsafe"
)

var patternDexMagic = regexp.MustCompile(`dex\n\d{3}\x00`)

type Header struct {
    Magic         [8]byte
    Checksum      uint32
    Signature     [20]byte
    FileSize      uint32
    HeaderSize    uint32
    EndianTag     uint32
    LinkSize      uint32
    LinkOff       uint32
    MapOff        uint32
    StringIdsSize uint32
    StringIdsOff  uint32
    TypeIdsSize   uint32
    TypeIdsOff    uint32
    ProtoIdsSize  uint32
    ProtoIdsOff   uint32
    FieldIdsSize  uint32
    FieldIdsOff   uint32
    MethodIdsSize uint32
    MethodIdsOff  uint32
    ClassIdsSize  uint32
    ClassIdsOff   uint32
    DataSize      uint32
    DataOff       uint32
}

func NewHeader(buffer []byte) Header {
    var header Header
    headerSize := unsafe.Sizeof(header)
    headerBytes := buffer[:headerSize]

    headerPtr := (*Header)(unsafe.Pointer(&headerBytes[0]))
    header = *headerPtr

    return header
}

type MemoryDex struct {
    memory *RemoteMemory
    mmap   *MemoryMap
    header Header
}

func NewMemoryDex(memory *RemoteMemory, mmap *MemoryMap) (*MemoryDex, error) {
    headerSize := unsafe.Sizeof(Header{})
    buffer := make([]byte, headerSize)

    err := memory.ReadMemory(mmap, buffer)
    if err != nil {
        return nil, err
    }

    header := NewHeader(buffer)

    return &MemoryDex{
        memory: memory,
        mmap:   mmap,
        header: header,
    }, nil
}

func (m *MemoryDex) IsValid() bool {
    if !patternDexMagic.Match(m.header.Magic[:]) {
        fmt.Printf("Invalid header magic: %s\n", string(m.header.Magic[:]))
        return false
    }
    if m.mmap.Size() < uint64(m.header.FileSize) {
        fmt.Println("File size doesn't match")
        return false
    }
    // https://source.android.com/docs/core/runtime/dex-format?hl=zh-cn#endian-constant
    if m.header.EndianTag != 0x12345678 && m.header.EndianTag != 0x78563412 {
        fmt.Println("Invalid endian tag")
        return false
    }
    if m.header.TypeIdsSize > 65535 || m.header.ProtoIdsSize > 65535 {
        fmt.Println("Too many method ids or proto ids")
        return false
    }
    return true
}

func (m *MemoryDex) Size() uint32 {
    return m.header.FileSize
}

func (m *MemoryDex) Dump(outputFile string) error {
    buffer := make([]byte, m.Size())

    err := m.memory.ReadMemory(m.mmap, buffer)
    if err != nil {
        return err
    }

    err = ioutil.WriteFile(outputFile, buffer, 0644)
    if err != nil {
        return err
    }

    return nil
}
