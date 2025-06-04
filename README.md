# API de Clima por CEP

Esta aplicação fornece informações meteorológicas (temperatura em Celsius, Fahrenheit e Kelvin) para um CEP brasileiro.

## Requisitos

- Go 1.21 ou superior
- Docker e Docker Compose (para implantação em contêiner)
- Chave de API WeatherAPI (cadastre-se em https://www.weatherapi.com/)

## Estrutura do Projeto

O projeto segue os princípios da Arquitetura Limpa Simplificada:

- `internal/entity`: Entidades de domínio
- `internal/usecase`: Lógica de negócios
- `internal/gateway`: Definições de interfaces para serviços externos
- `internal/infra/gateway`: Implementações concretas de gateways
- `internal/infra/web/handlers`: Manipuladores HTTP

## Executando Localmente

### Sem Docker

1. Configure a chave WeatherAPI como uma variável de ambiente:
   ```bash
   export WEATHERAPI_KEY=sua_chave_api_aqui
   ```

2. Execute a aplicação:
   ```bash
   go run main.go
   ```

### Com Docker Compose

1. Crie um arquivo `.env` com sua chave WeatherAPI:
   ```
   WEATHERAPI_KEY=sua_chave_api_aqui
   ```

2. Construa e execute o contêiner:
   ```bash
   docker-compose up --build
   ```

## Uso da API

### Obter Temperatura por CEP

```
GET /temperature?zipcode=12345678
```

#### Resposta de Sucesso (200 OK)

```json
{
  "temp_C": 28.5,
  "temp_F": 83.3,
  "temp_K": 301.5
}
```

#### Respostas de Erro

- Formato de CEP inválido (422 Unprocessable Entity):
  ```
  invalid zipcode
  ```

- CEP não encontrado (404 Not Found):
  ```
  can not find zipcode
  ```

## Testes

Execute os testes com:

```bash
go test ./...
```

## Implantação no Google Cloud Run

1. Construa e envie a imagem Docker para o Google Container Registry:
   ```bash
   gcloud builds submit --tag gcr.io/seu-id-projeto/goexpert-clima
   ```

2. Implante no Cloud Run:
   ```bash
   gcloud run deploy goexpert-clima \
     --image gcr.io/seu-id-projeto/goexpert-clima \
     --platform managed \
     --region us-central1 \
     --set-env-vars WEATHERAPI_KEY=sua_chave_api_aqui
   ```

3. O serviço está disponível na seguinte URL fornecida pelo Cloud Run:

- https://cep-weather-902894090648.southamerica-east1.run.app/temperature?zipcode=89240000

## Licença

Este projeto está licenciado sob a Licença MIT.
