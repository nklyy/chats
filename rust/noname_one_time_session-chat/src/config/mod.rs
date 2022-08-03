use dotenv::dotenv;
use rocket::{figment::Figment, Config as RocketConfig};
use std::env;

pub struct Config {
    pub port: String,
    pub environment: String,

    pub mongo_uri: String,
    pub redis_uri: String,
}

impl Config {
    pub fn new() -> Config {
        dotenv().ok();

        let port = match env::var("PORT") {
            Ok(port) => port,
            Err(_) => panic!("incorrect port"),
        };

        let environment = match env::var("APP_ENV") {
            Ok(environment) => environment,
            Err(_) => panic!("incorrect app_env"),
        };

        let mongo_uri = match env::var("MONGO_URI") {
            Ok(environment) => environment,
            Err(_) => panic!("incorrect mongo_uri"),
        };

        let redis_uri = match env::var("REDIS_URI") {
            Ok(environment) => environment,
            Err(_) => panic!("incorrect redis_uri"),
        };

        Config {
            port,
            environment,
            mongo_uri,
            redis_uri,
        }
    }

    pub fn from_env() -> Figment {
        let local_cfg = Config::new();

        let port: u16 = local_cfg
            .port
            .parse()
            .expect("PORT environment variable should parse to an integer");

        RocketConfig::figment().merge(("port", port))
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn create_config() {
        let c = Config::new();
        assert_eq!(c.port.chars().count() > 0, true);
        assert_eq!(c.environment.chars().count() > 0, true);
        assert_eq!(c.mongo_uri.chars().count() > 0, true);
        assert_eq!(c.redis_uri.chars().count() > 0, true);
    }
}
