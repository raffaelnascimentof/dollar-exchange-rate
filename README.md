# dollar-exchange-rate

Desafio Go Expert - Full Cycle

Criar um server.go para consumir uma API de cotação do Dólar com base no Real e um client.go para realizar uma request HTTP no server.go solicitando a cotação do Dólar.

Criterios de aceite:
- [x] Server deve ter timeout de 200ms para consumir a API de cotação https://economia.awesomeapi.com.br/json/last/USD-BRL
- [x] Server deve salvar a cotação em um banco de dados com timeout de 10ms
- [x] Client deve receber apenas o valor do câmbio do Dólar
- [x] Client deve escrever em um arquivo (cotacao.txt) a cotação solicitada


# Como rodar ?

- Clonar repositorio
- Baixar dependencias (go mod tidy)
- Rodar arquivo do docker-compose para criar banco de dados (docker-compose up -d)
- Rodar server.go (go run server.go)
- Rodar client.go (go run client.go)
