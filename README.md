### Desafio técnico: rate-limiter
Desenvolver um rate limiter em Go que possa ser configurado para limitar o número máximo de requisições por segundo com base em um endereço IP específico ou em um token de acesso. 

## Pré-requisitos

- Docker
- Docker Compose

## Executando o projeto

1. Clone o repositório (https): https://github.com/saraqueirozs/challenge-goexpert-rate-limiter.git 

2. Inicie os serviços usando Docker Compose:
   ```
   docker-compose up
   ```
## Para testar local 

O serviço deverá retornar a informação de sucesso (200), caso exceda o limite, deverá bloquear.
```
curl --request GET \
--url http://localhost:8080/ \
--header 'API_KEY: qwe321'

```