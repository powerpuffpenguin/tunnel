Target="tunnel"
Docker="king011/tunnel"
Dir=$(cd "$(dirname $BASH_SOURCE)/.." && pwd)
Platforms=(
    darwin/amd64
    windows/amd64
    linux/arm
    linux/amd64
)