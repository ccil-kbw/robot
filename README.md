# Development
First you need to create a discord bot and generate a token for it.

Join a test server of your own and get your server's id

In the context of this code these are called Token and Guild. When you got those run the following command (replace those numbers with the guild id and token you got previously):

``` shell
go run ./cmd/discord_bot/main.go -guild 1234567890 -token 1234567890
```
