use mongodb::Client as MongoClient;
use redis::Client as RedisClient;

#[derive(Debug)]
pub struct Databases {
    pub mongo_client: MongoClient,
    pub redis_client: RedisClient,
}

impl Databases {
    pub async fn new(mongo_uri: String, redis_uri: String) -> Result<Databases, String> {
        let mongo_client = match MongoClient::with_uri_str(mongo_uri).await {
            Ok(client) => client,
            Err(_) => return Err("failed to create MongoDB client".to_string()),
        };

        let redis_client = match RedisClient::open(redis_uri) {
            Ok(client) => client,
            Err(_) => return Err("failed to create redis client".to_string()),
        };

        Ok(Databases {
            mongo_client,
            redis_client,
        })
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[tokio::test]
    async fn databases_connections() {
        let _dbs = Databases::new(
            "mongodb://localhost:27017".to_string(),
            "redis://localhost".to_string(),
        )
        .await
        .unwrap();
    }

    #[tokio::test]
    async fn incorrect_mongo_uri() {
        let dbs = Databases::new("".to_string(), "redis://localhost".to_string())
            .await
            .unwrap_err();

        assert_eq!(dbs, "failed to create MongoDB client".to_string())
    }

    #[tokio::test]
    async fn incorrect_redis_uri() {
        let dbs = Databases::new("mongodb://localhost:27017".to_string(), "".to_string())
            .await
            .unwrap_err();

        assert_eq!(dbs, "failed to create redis client".to_string())
    }
}
