use dotenv::dotenv;
use std::env;

pub struct Config {
    pub port: u16,
    pub environment: String,

    pub mongo_uri: String,
    pub redis_uri: String,
}

impl Config {
    pub fn init() -> Config {
        dotenv().ok();

        let port = match env::var("PORT") {
            Ok(port) => port,
            Err(_) => panic!("incorrect port"),
        };

        let port: u16 = port.parse().expect("Can't parse port into number");

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
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn create_config() {
        let c = Config::init();
        assert_eq!(c.port > 0, true);
        assert_eq!(c.environment.chars().count() > 0, true);
        assert_eq!(c.mongo_uri.chars().count() > 0, true);
        assert_eq!(c.redis_uri.chars().count() > 0, true);
    }
}
