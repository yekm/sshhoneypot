default: run

id_rsa:
	ssh-keygen -P "" -t rsa -f id_rsa

sshhoneypot: sshhoneypot.go
	go build

run: sshhoneypot id_rsa
	./sshhoneypot id_rsa 0.0.0.0:13222 |& tee -a log

u:
	cat log | cut -f2 -d- | cut -f2 -d' ' | sort -u

ggeoip: ggeoip.c
	$(CC) -lGeoIP $< -o $@

ip:
	cat log | cut -f3 -d' ' | cut -f1 -d':' | sort -u

counties: ggeoip
	make ip | xargs -n1 ./ggeoip /usr/share/GeoIP/GeoIPCity.dat | cut -f 2 | sort | uniq -c | sort -rn
