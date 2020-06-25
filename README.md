# go-graphql

Run locally

1. rename .env.copy to only .env and set your env variables inside the file
2. execute `go run server.go`

Deploying to Google App Engine

1. Rename app.yaml.copy to only app.yaml and set your env variables inside the file


Deploying

gcloud app deploy

Logs

gcloud app logs tail -s default