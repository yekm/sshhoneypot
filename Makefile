default: run

id_rsa:
	ssh-keygen -P "" -t rsa -f id_rsa

build:
	go build

run: build id_rsa
	./sshhoneypot id_rsa 0.0.0.0:13222 |& tee -a log

u:
	cat log | cut -f2 -d- | cut -f2 -d' ' | sort -u
