package main

import (
    "DexDump/cmd"
    "fmt"
    "log"
    "os"
)

func main() {

    // 初始化日志记录器
    log.SetOutput(os.Stdout)

    args := cmd.ParseArgs()

    if err := cmd.Dump_dex_file(args); err != nil {
        fmt.Println("[!] Error:", err)
    }
}
