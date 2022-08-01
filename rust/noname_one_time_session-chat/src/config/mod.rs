use dotenv::dotenv;
use std::env;

pub struct Config {
    pub port: String,
    pub environment: String,
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

        Ok(Config { port, environment })
    }
}
