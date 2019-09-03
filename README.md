# Paredão

## Descrição
* Projeto para ouvir filas de uma mensageria (ActiveMQ), processar e salvar em um banco de dados. (Mongo)

## Requisitos Mínimos
* [Go 1.12+](https://golang.org/)
* [Docker](https://www.docker.com/)

## Tecnologias utilizadas
* [Go 1.12+](https://golang.org/)
* [MongoDB](https://www.mongodb.com/)
* [ActiveMQ](https://activemq.apache.org/)
* [Docker](https://www.docker.com/)
* [Docker-compose](https://docs.docker.com/compose/)

## Instalação
1. Baixe o repositório como arquivo zip ou faça um clone;
2. Descompacte os arquivos em seu computador;
3. Abra a pasta decompactada
4. Execute `make up`
5. Aguarde até a stack inteira estar deployada.
6. Acesse o ActiveMQ (www.localhost:8161) para enviar mensagens para a fila e simular um projeto real.


## Exemplo de mensagem

Usuário:

```javascript
{
   "_id":"123",
   "email":"emailteste",
   "username":"usernameTeste",
   "fullName":"fullnameTeste",
   "gender":"genderTeste",
   "status":"statusTeste",
   "birthDate":"birthdateTeste",
   "phones":{
      "phone":"phoneTeste",
      "cellphone":"cellphoneTeste",
      "ddd_cellphone":"21",
      "mobile_phone_confirmed":true
   },
   "clientId":"clientTeste"
}
```

## Arquitetura de Solução
TODO

## Dúvidas?
`Caso tenha dúvidas ou precise de suporte, mande um email para rafacoaraujo@gmail.com`
