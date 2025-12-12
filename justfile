app_name := "straico-cli"
module_name := "github.com/tyler71/straico-cli/m/v0"

build:
    go build -ldflags "-s -w" -o straico

build-all:
    ./build.sh

debug:
    go build -gcflags="all=-N -l" -o straico-cli
    echo 0 | sudo tee /proc/sys/kernel/yama/ptrace_scope
    echo "Run -> Attach to Process"

test:
    go test -v ./...

    echo "Running tests with coverage..."
    go test -coverprofile=.coverage.out ./...
    go tool cover -func=.coverage.out

demo:
    vhs demo.tape