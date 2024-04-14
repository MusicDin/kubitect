build:
	go build -ldflags "-s -w" -trimpath -o kubitect ./cmd

test:
	go test ./... -v
