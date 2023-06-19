# lambda-iam-autentication

POC Lambda for technical purposes

Lambda call another endpoint with IAM autorizar enabled for test a signV4 autentication

Diagrama Flow

    Lambda ( get SecretKey (access-key-id + access-key) )==>
    ApiGW (AIM authorizer) == Lambda SignIn Mock JWT 

## Compile

    GOOD=linux GOARCH=amd64 go build -o ../build/main main.go

    zip -jrm ../build/main.zip ../build/main

    aws lambda update-function-code \
    --function-name lambda-iam-autentication \
    --zip-file fileb:///mnt/c/Eliezer/workspace/github.com/lambda-iam-autentication/build/main.zip \
    --publish

## Endpoint

Worker (no endpoint)


