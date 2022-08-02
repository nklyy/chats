mod config;
mod databases;

use config::Config;
use databases::Databases;
use std::process;

pub async fn execute() {
    println!("Hello, world!");

    let local_cfg = Config::new().unwrap_or_else(|err| {
        eprintln!("[ERROR] Problem parsing env arguments: {}", err);

        process::exit(1);
    });

    let _dbs = Databases::new(local_cfg.mongo_uri, local_cfg.redis_uri)
        .await
        .unwrap_or_else(|err| {
            eprintln!("[ERROR] Failed connect to databases: {}", err);

            process::exit(1);
        });

    println!("{}", local_cfg.port);
}
