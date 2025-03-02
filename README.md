docker run --name mongodb -p 27017:27017 -d mongodb/mongodb-community-server:latest
docker run --name mongo-express -p 8081:8081 --link mongodb:mongo -e ME_CONFIG_MONGODB_SERVER=mongo -d mongo-express

docker run -p 8000:8000 casbin/casdoor-all-in-one         





# Next Steps

* Add Create-option for Recipes
* Integrate Minio for the Pictures
* Fix Header beeing Responsive
* Add Dockerfile
* Add Readme (Build & Containers & usw.)
* Fix Makefile
* Add Configuration
* Add Multilanguage support