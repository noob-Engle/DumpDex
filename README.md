## 简易脱壳小工具
参考frida-dexdump

```bash
adb push DumpDex /data/local/tmp
adb shell chmod +x /data/local/tmp/DumpDex
adb shell
su 
cd /data/local/tmp
./DumpDex -pid '$(pidof com.example.app)' -o <output_dir>
```
