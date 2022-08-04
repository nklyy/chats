use mongodb::Client as MongoClient;
use redis::Client as RedisClient;
use rocket::fairing::AdHoc;

#[derive(Debug)]
pub struct Databases {
    pub mongo_client: MongoClient,
    pub redis_client: RedisClient,
}

impl Databases {
    async fn new_clients(mongo_uri: String, redis_uri: String) -> Databases {
        let mongo_client = match MongoClient::with_uri_str(mongo_uri).await {
            Ok(client) => client,
            Err(err) => panic!("failed to create MongoDB client: {}", err),
        };

        let redis_client = match RedisClient::open(redis_uri) {
            Ok(client) => client,
            Err(err) => panic!("failed to create redis client: {}", err),
        };

        Databases {
            mongo_client,
            redis_client,
        }
    }

    pub fn init(mongo_uri: String, redis_uri: String) -> AdHoc {
        AdHoc::on_ignite("Connecting to MongoDB", |rocket| async move {
            let dbs = Self::new_clients(mongo_uri, redis_uri).await;
            rocket.manage(dbs)
        })
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[tokio::test]
    async fn databases_connections() {
        let _dbs = Databases::new_clients(
            "mongodb://localhost:27017".to_string(),
            "redis://localhost".to_string(),
        )
        .await;
    }

    #[tokio::test]
    #[should_panic]
    async fn incorrect_mongo_uri() {
        let _dbs = Databases::new_clients("".to_string(), "redis://localhost".to_string()).await;
    }

    #[tokio::test]
    #[should_panic]
    async fn incorrect_redis_uri() {
        let _dbs =
            Databases::new_clients("mongodb://localhost:27017".to_string(), "".to_string()).await;
    }
}
