docker run --name mongodb -p 27017:27017 -d mongodb/mongodb-community-server:latest
docker run --name mongo-express -p 8081:8081 --link mongodb:mongo -e ME_CONFIG_MONGODB_SERVER=mongo -d mongo-express
docker run -p 9000:9000 -p 9001:9001 -d -e "MINIO_ROOT_USER=AKIAIOSFODNN7EXAMPLE" -e "MINIO_ROOT_PASSWORD=wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY" quay.io/minio/minio server /data --console-address ":9001"

docker run -p 8000:8000 casbin/casdoor-all-in-one         



# Setup with Docker
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

docker run \
  -e driverName=mysql \
  -e dataSourceName='myuser:mypassword@tcp(host.docker.internal:3307)/' \
  -p 8007:8000 \
  casbin/casdoor:latest




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