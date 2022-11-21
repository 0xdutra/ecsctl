all:
	go build

install:
	sudo install -o root -g root -m 0755 ecsctl /usr/local/bin/ecsctl

clean:
	rm -f ecsctl
