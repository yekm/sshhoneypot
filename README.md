# sshhoneypot

2 days:
```
$ make counties |head
make ip | xargs -n1 ./ggeoip /usr/share/GeoIP/GeoIPCity.dat | cut -f 2 | sort | uniq -c | sort -rn
    109 CN
     56 US
     36 RU
     26 KR
     26 BR
     13 FR
     12 JP
     11 ID
     10 SG
```
