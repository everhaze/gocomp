# gocomp
QoL tool for compiling Go code

## Compiling
```sh
wget https://github.com/everhaze/gocomp/archive/refs/heads/main.tar.gz
tar -xzf gocomp-main
cd gocomp-main
go mod tidy
go build -o gocomp
```

## Usage
There are 2 ways to use it.
### 1. CLI
```sh
#      mode  cgo  name
gocomp qs/cc 0/1 myapp

# qs = quick start
# cc = clean cache

# This is the bare minimum. It is the equivalent to 'CGO_ENABLED=0 go build -ldflags="-s -w" -trimpath'
gocomp qs

# Same example but with CGO_ENABLED=1
gocomp qs 1

# Same example but with a custom name
gocomp qs 0 myapp

# Note: If cgo isnt set to 1, it is always 0.
```

### 2. TUI
```sh
# For the TUI, you just run gocomp without passing any args.
gocomp
```
