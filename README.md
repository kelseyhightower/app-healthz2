# app-healthz2

app-healthz2 is a sample app that demonstrates how to leverage the health endpoint pattern.

## Create Docker Image

Build the go binary

```
GOOS=linux bash build
```

```
docker build -t kelseyhightower/app-healthz:2.0.0
```

```
docker push kelseyhightower/app-healthz:2.0.0
```
