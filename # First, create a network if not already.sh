# First, create a network if not already done
docker network create minio-network



# Run your application container with the correct endpoint
docker run --network minio-network -p 3000:3000 \
  -e MINIO_ENDPOINT=miniocontainer:9000 \
  -e MINIO_ACCESS_KEY_ID=uS4yfmTlVP1rBakV1TzW \
  -e MINIO_SECRET_ACCESS_KEY=peOlGbNFAAP2Rcr4DvfPijgB1FKEkBMynURpOZAG \
  -e MONGO_ENDPOINT=mongodb://mongodb:27017 \
  -e REDIS_ENDPOINT=redis:6379 \
  -e CALLBACK_ADDRESS=http://localhost:3000/callback \
  -e CASDOOR_ENDPOINT=https://auth.ferber.io \

  your-image-name


// docker run -d --name redis --network minio-network -p 6379:6379 -p 8001:8001 redis:latest


# Run MinIO container
docker run --name miniocontainer --network minio-network -p 9000:9000 -p 9001:9001 -d \
  -e "MINIO_ROOT_USER=uS4yfmTlVP1rBakV1TzW" \
  -e "MINIO_ROOT_PASSWORD=peOlGbNFAAP2Rcr4DvfPijgB1FKEkBMynURpOZAG" \
  quay.io/minio/minio server /data --console-address ":9001"

  
docker run -d --name redis --network minio-network -p 6379:6379 -p 8001:8001 redis/redis-stack:latest
  docker run --name mongodb --network minio-network -p 27017:27017 -d mongodb/mongodb-community-server:latest
  docker run --name mongo-express --network minio-network -p 8081:8081 --link mongodb:mongo -e ME_CONFIG_MONGODB_SERVER=mongo -d mongo-express