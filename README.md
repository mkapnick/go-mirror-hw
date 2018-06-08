## Running locally
- `GOOS=linux go build -o main`
- `docker run -p 9000:9000 go-mirror`

## Deployment
- Deployed on Kubernetes/GCP
- Can be found here: [http://35.232.131.106](http://35.232.131.106)

## DockerHub
- Can be found on docker hub here: [https://hub.docker.com/r/mdotm/go-mirror](https://hub.docker.com/r/mdotm/go-mirror)
- `docker pull mdotm/go-mirror`
- `docker run -p 9000:9000 mdotm/go-mirror`
