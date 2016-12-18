
## populate the redis data

```
HMSET album:1 title "Electric Ladyland" artist "Jimi Hendrix" price 4.95 likes 8
HMSET album:2 title "Back in Black" artist "AC/DC" price 5.95 likes 3
HMSET album:3 title "Rumours" artist "Fleetwood Mac" price 7.95 likes 12
HMSET album:4 title "Nevermind" artist "Nirvana" price 5.95 likes 8
ZADD likes 8 1 3 2 12 3 8 4
```

## Refs
* <https://xiequan.info/working-redis-go/>
