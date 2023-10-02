#### You can run the app directly from this url: https://fib.dre4success.com

##### To run it locally, docker is the easiest way and below are the steps:

- Install Docker on your machine if you haven't already. You can download Docker from the official website: https://www.docker.com/products/docker-desktop

- Clone the repository to your local machine:
```
git clone https://github.com/dre4success/fibonacci
```
- Navigate to the root directory of the project.

- Build the docker image:
```
docker build -t app-name .
```
- Run the docker container:
```
docker run -p 8080:8080 -p 8081:8081 app-name
```

