# Splat - a small Slack app library
Currently Splat just uses the slash commands API for Slack

### Usage
This example registers a command `/test` and listens for requests on port `:3000` at the `/slackcommands` endpoint
```go
app := splat.New("YourSigningSecret")

app.RegisterCommand("/test", func(p *splat.Payload) *splat.Response {
    return &splat.Response{Text: "Hello, World!"}
})

log.Fatal(app.Open(":3000", "/slackcommands"))
```