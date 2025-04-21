docker run --name mongodb -p 27017:27017 -d mongodb/mongodb-community-server:latest
docker run --name mongo-express -p 8081:8081 --link mongodb:mongo -e ME_CONFIG_MONGODB_SERVER=mongo -d mongo-express
docker run -p 9000:9000 -p 9001:9001 -d -e "MINIO_ROOT_USER=AKIAIOSFODNN7EXAMPLE" -e "MINIO_ROOT_PASSWORD=wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY" quay.io/minio/minio server /data --console-address ":9001"

docker run -p 8000:8000 casbin/casdoor-all-in-one   



docker run --name miniocontainer --network minio-network -p 9000:9000 -p 9001:9001 -d -e "MINIO_ROOT_USER=AKIAIOSFODNN7EXAMPLE" -e "MINIO_ROOT_PASSWORD=wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY" quay.io/minio/minio server /data --console-address ":9001"


docker run --network minio-network -p 3000:8000 -e MINIO_ENDPOINT=host.docker.internal:9000 your-image-name



# Setup with Docker
## Setup Casdoor

After the container for casdoor is running the application has to be registered.

mysql for casdoor
```
docker run --name my-mysql-container \
  -e MYSQL_ROOT_PASSWORD=rootpassword \
  -e MYSQL_USER=myuser \
  -e MYSQL_PASSWORD=mypassword \
  -e MYSQL_DATABASE=casdoor \
  -p 3307:3306 \
  -d mysql:latest
```

casdoor
```
docker run \
  -e driverName=mysql \
  -e dataSourceName='myuser:mypassword@tcp(host.docker.internal:3307)/' \
  -p 8007:8000 \
  casbin/casdoor:latest
```

## docker container for the app

create network:
```docker
docker network create recipe-app-network
```

evtl. change the network.
Minio needs to have a client and key to start
```
docker run --name miniocontainer --network recipe-app-network -p 9000:9000 -p 9001:9001 -d \
  -e "MINIO_ROOT_USER=uS4yfmTlVP1rBakV1TzW" \
  -e "MINIO_ROOT_PASSWORD=peOlGbNFAAP2Rcr4DvfPijgB1FKEkBMynURpOZAG" \
  quay.io/minio/minio server /data --console-address ":9001"
```

Run Redis
docker run -d --name redis --network recipe-app-network -p 6379:6379 -p 8001:8001 redis/redis-stack:latest

run mongodb
docker run --name mongodb --network recipe-app-network -p 27017:27017 -d mongodb/mongodb-community-server:latest

// docker run --name mongo-express --network minio-network -p 8081:8081 --link mongodb:mongo -e ME_CONFIG_MONGODB_SERVER=mongo -d mongo-express



# Next Steps

* ~~Add Create-option for Recipes~~ 👍
* ~~Integrate Minio for the Pictures~~ 👍
* ~~Keywords for Creating~~ 👍
* ~~Correct Timeformat~~ 👍
* ~~Import from other recipes~~ 👍
* ~~Correct User name display~~ 👍
* ~~Edit mode~~
* change from recipe.Nutrition.ServingSize to recipeYield
* Fix Header beeing Responsive
* Add Dockerfile
* Add Readme (Build & Containers & usw.)
* Fix Makefile
* Add Configuration
* Add Multilanguage support
* Favorite-Support