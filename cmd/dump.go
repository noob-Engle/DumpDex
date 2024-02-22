package cmd

import (
    "fmt"
    "io/ioutil"
    "log"
    "os"
    "os/exec"
    "path/filepath"
    "regexp"
    "strconv"
    "strings"
)

var PATTERN_TOP_ACTIVITY = regexp.MustCompile(`ACTIVITY (\S+) .*pid=(\d+)\n.*\n.*?mResumed=true`)

func isDirEmpty(dir string) (bool, error) {
    f, err := os.Open(dir)
    if err != nil {
        return false, err
    }
    defer f.Close()

    _, err = f.Readdirnames(1)
    if err == nil {
        return false, nil
    }
    if len(err.Error()) > 0 && err.Error()[0] == 'n' {
        return true, nil
    }
    return false, err
}

func Dump_dex_file(args Args) error {
    if args.OutputDir != "" {
        err := os.MkdirAll(args.OutputDir, os.ModePerm)
        if err != nil {
            return err
        }
        empty, err := isDirEmpty(args.OutputDir)
        if err != nil {
            return err
        }
        if !empty {
            fmt.Println("[!] Target directory is not empty, abort!")
            return nil
        }
    }
    //memory
    var memory *RemoteMemory
    if args.Pid != 0 {
        _ = NewRemoteMemory(args.Pid)
    } else {
        output, err := exec.Command("/system/bin/dumpsys", "activity", "top").Output()
        if err != nil {
            log.Fatal(err)
        }
        captures := PATTERN_TOP_ACTIVITY.FindSubmatch(output)
        if len(captures) != 3 {
            log.Fatal("failed to get pid of the top activity")
        }

        topActivity := string(captures[1])
        topActivityPid, err := strconv.Atoi(string(captures[2]))
        if err != nil {
            log.Fatal(err)
        }

        fmt.Printf("[*] Top activity: %s (%d)\n", topActivity, topActivityPid)

        _ = NewRemoteMemory(topActivityPid)
    }
    mappingsslice := make([]string, 0) // 创建一个空的字符串切片
    mappings, err := memory.ReadMaps()
    if err != nil {
        // 处理错误
    }
    for _, memoryMap := range mappings {
        dex, err := NewMemoryDex(memory, &memoryMap)
        if err != nil {
            if args.Verbose {
                fmt.Printf("[*] Skipped: %s\n", memoryMap.Name())
            }
            continue
        }

        if !dex.IsValid() {
            if args.Verbose {
                fmt.Printf("[*] Skipped: %s\n", memoryMap.Name())
            }
            continue
        }

        mappingsslice = append(mappingsslice, memoryMap.Name())

        fmt.Printf("[*] Processed: %s\n", memoryMap.Name())
    }

    // 遍历映射
    for _, mapping := range mappingsslice {
        fmt.Println(mapping)
    }
    if args.OutputDir != "" {
        outputFile := filepath.Join(args.OutputDir, "mappings.txt")
        content := strings.Join(mappingsslice, "\n")
        err := ioutil.WriteFile(outputFile, []byte(content), 0644)
        if err != nil {
            return err
        }

        fmt.Printf("[*] Dumped %d dex file(s) to %s\n", len(mappings), args.OutputDir)
    }
    return nil
}
