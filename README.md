# Omnidots API Client
The Omnidots API docs are [here](https://honeycomb.omnidots.com/api/docs).

Fed those docs to ChatGPT, told it to make an OpenAPI spec, then used a code [generator](https://github.com/deepmap/oapi-codegen) to get the Go client code

## Using the CLI  
```bash
go build -o omnidots-cli ./cmd/cli/main.go
./omnidots-cli --token {TOKEN HERE} --command list-sensors
```