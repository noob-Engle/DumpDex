package cmd

import (
    "fmt"
    "io/ioutil"
    "os"
    "regexp"
    "strconv"
    "strings"
)

var PATTERN_MEMORY_MAP = regexp.MustCompile(`(\S+)-(\S+) (\S+) \S+ \S+ \S+\s+(\S.*)?`)

type MemoryMap struct {
    Address  [2]uint64
    Perms    string
    Filepath string
}

func (m *MemoryMap) StartAddress() uint64 {
    return m.Address[0]
}

func (m *MemoryMap) Size() uint64 {
    return m.Address[1] - m.Address[0]
}

func (m *MemoryMap) Name() string {
    if m.Filepath != "" {
        return m.Filepath
    } else {
        return "[anonymous memory]"
    }
}

type RemoteMemory struct {
    pid int
}

func NewRemoteMemory(pid int) *RemoteMemory {
    return &RemoteMemory{pid: pid}
}

func (r *RemoteMemory) ReadMaps() ([]MemoryMap, error) {
    mapsBytes, err := ioutil.ReadFile(fmt.Sprintf("/proc/%d/maps", r.pid))
    if err != nil {
        return nil, err
    }

    mapsStr := string(mapsBytes)
    lines := strings.Split(mapsStr, "\n")

    var memoryMaps []MemoryMap

    for _, line := range lines {
        if line == "" {
            continue
        }

        captures := PATTERN_MEMORY_MAP.FindStringSubmatch(line)[1:]
        startAddress, _ := strconv.ParseUint(captures[0], 16, 64)
        endAddress, _ := strconv.ParseUint(captures[1], 16, 64)
        perms := captures[2]
        filepath := captures[3]

        memoryMap := MemoryMap{
            Address:  [2]uint64{startAddress, endAddress},
            Perms:    perms,
            Filepath: filepath,
        }

        memoryMaps = append(memoryMaps, memoryMap)
    }

    return memoryMaps, nil
}

func (r *RemoteMemory) ReadMemory(memoryMap *MemoryMap, buffer []byte) error {
    memFile, err := os.Open(fmt.Sprintf("/proc/%d/mem", r.pid))
    if err != nil {
        return err
    }
    defer memFile.Close()

    _, err = memFile.Seek(int64(memoryMap.StartAddress()), 0)
    if err != nil {
        return err
    }

    _, err = memFile.Read(buffer)
    if err != nil {
        return err
    }

    return nil
}
