For local development.

## Start Docker DynamoDB.

```
docker-compose up dynamo
```

### Start Server

```
make dev-run
```

### Build and run Client

```
make dev-crun
```

### Client runs with

```
./build/client -h
```

## Run DynamoDB GUI client (useful)

https://github.com/Arattian/DynamoDb-GUI-Client

## For AWS depolyment.

```
make prod-build
terraform init
terraform apply -var-file="main.tfvars"
```

## Reference:

https://github.com/nerdguru/go-sls-crudl
https://github.com/developmentseed/tf-lambda-proxy-apigw
https://learn.hashicorp.com/terraform/aws/lambda-api-gateway
https://www.alexedwards.net/blog/serverless-api-with-go-and-aws-lambda
