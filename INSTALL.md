# Install Go on Ubuntu Linux #

## With brew (recommended)
> brew is one of the most popular package managers for macOS and is also available on Linux, and it works better than
> other package managers, like snap or flatpak.

Just: `brew install go` installs the latest version of Go.

## Ubuntu APT
Installs the 1.24.2 version of Go.

- `sudo add-apt-repository ppa:longsleep/golang-backports && apt update`
- `sudo apt install -y golang-go`

## Check the installation

`go version`
