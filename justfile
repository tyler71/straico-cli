build:
 go build -ldflags "-s -w" -o straico-cli
debug:
 go build -gcflags="all=-N -l" -o straico-cli
 echo 0 | sudo tee /proc/sys/kernel/yama/ptrace_scope
 echo "Run -> Attach to Process"

