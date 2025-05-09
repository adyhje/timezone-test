Curl Befehle zum testen:

Einfügen
#####
```
curl -X POST http://localhost:8080/times \
  -H "Content-Type: application/json" \
  -d '{"timestamp":"2025-05-09T12:00:00Z"}'

curl -X POST http://localhost:8080/times \
  -H "Content-Type: application/json" \
  -d '{"timestamp":"2025-05-09T14:00:00+02:00"}'
```

Selektieren
#####
```
curl http://localhost:8080/times
```
