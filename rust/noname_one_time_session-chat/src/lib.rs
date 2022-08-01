mod config;

use config::Config;
use std::process;

pub fn execute() {
    println!("Hello, world!");

    let local_cfg = Config::new().unwrap_or_else(|err| {
        eprintln!("[ERROR] Problem parsing env arguments: {}", err);

        process::exit(1);
    });

    println!("{}", local_cfg.port);
}
