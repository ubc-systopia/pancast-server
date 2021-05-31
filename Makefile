.PHONY: app

all: app

app:
	go build -o app cmd/app/main.go

clean:
	rm app 2> /dev/null || true