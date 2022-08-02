use dotenv::dotenv;
use std::env;

pub struct Config {
    pub port: String,
    pub environment: String,

    pub mongo_uri: String,
    pub redis_uri: String,
}

impl Config {
    pub fn new() -> Result<Config, String> {
        dotenv().ok();

        let port = match env::var("PORT") {
            Ok(port) => port,
            Err(_) => return Err("incorrect port".to_string()),
        };

        let environment = match env::var("APP_ENV") {
            Ok(environment) => environment,
            Err(_) => return Err("incorrect app_env".to_string()),
        };

        let mongo_uri = match env::var("MONGO_URI") {
            Ok(environment) => environment,
            Err(_) => return Err("incorrect mongo_uri".to_string()),
        };

        let redis_uri = match env::var("REDIS_URI") {
            Ok(environment) => environment,
            Err(_) => return Err("incorrect redis_uri".to_string()),
        };

        Ok(Config {
            port,
            environment,
            mongo_uri,
            redis_uri,
        })
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn create_config() {
        let c = Config::new().unwrap();
        assert_eq!(c.port.chars().count() > 0, true);
        assert_eq!(c.environment.chars().count() > 0, true);
        assert_eq!(c.mongo_uri.chars().count() > 0, true);
        assert_eq!(c.redis_uri.chars().count() > 0, true);
    }
}
