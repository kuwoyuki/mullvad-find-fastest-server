## mullvad-find-fastest-server

> Find the fastest Wireguard server for Mullvad

```sh
# setcap the binary to allow it to bind to raw sockets
go build && doas setcap cap_net_raw=+ep ./mullvad-find-fastest-server &&
./mullvad-find-fastest-server

# or just run as root
doas go run main.go
```

## License

MIT 2018 - Victor Bjelkholm
