
# GoPad

## Uses Fyne for GUI toolkit
https://fyne.io/


## Fedora 42 Dependencies
```
# Linux deps (X11)
sudo dnf install libX11-devel libXrandr-devel libXcursor-devel libXinerama-devel libXi-devel libXxf86vm-devel libX11-devel
sudo dnf install mesa-libGL-devel

# Windows deps
sudo dnf install mingw64-gcc mingw64-crt mingw64-winpthreads mingw64-windows-default-manifest
```


## Compilation
```
# for linux
go build main.go

# for Windows
CGO_ENABLED=1 GOOS=windows GOARCH=amd64 CC=x86_64-w64-mingw32-gcc go build -o gopad.exe
```

## Packaging
```
go install fyne.io/tools/cmd/fyne@latest
export PATH=$PATH:$HOME/go/bin
GOOS=windows CC=x86_64-w64-mingw32-gcc fyne package -os windows -name gopad.exe
```
