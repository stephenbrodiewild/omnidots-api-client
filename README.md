# Omnidots API Client
The Omnidots API docs are [here](https://honeycomb.omnidots.com/api/docs).

Fed those docs to ChatGPT, told it to make an OpenAPI spec, then used a code [generator](https://github.com/deepmap/oapi-codegen) to get the Go client code

## Run the example 
```bash
# set your token 
export OMNIDOTS_TOKEN=...
go run cmd/example/main.go
```