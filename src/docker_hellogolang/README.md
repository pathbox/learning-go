http://colobu.com/2015/10/12/create-minimal-golang-docker-images/


1. CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .
2. docker build -t example-scratch .
3. docker run -it -p 9097:9097 example-scratch


push 到你的DockerHub

```
docker tag `IMAGE ID` caryhub/example-scratch

docker images

docker push caryhub/example-scratch
```