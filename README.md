#### Pre-req
###### make a 'cert' director and keep the .crt and .key file inside it
###### make a .env file and hydrate the variables
##### For local development make sure temporal and terraform are installed

#### Temporal
##### RUN temporal server start-dev

#### API Server
##### RUN go run ./api/main.go

#### Worker
##### RUN go run ./worker/main.go
