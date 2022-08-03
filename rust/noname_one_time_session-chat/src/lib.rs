mod api;
mod config;
mod databases;

use config::Config;
use databases::Databases;
use rocket::serde::json::{json, Value};

#[macro_use]
extern crate rocket;

#[catch(404)]
fn not_found() -> Value {
    json!({
        "status": "error",
        "reason": "Resource was not found."
    })
}

#[launch]
pub fn rocket() -> _ {
    // let local_cfg = Config::new().unwrap_or_else(|err| {
    //     eprintln!("[ERROR] Problem parsing env arguments: {}", err);
    //     process::exit(1);
    // });

    let local_cfg = Config::new();

    rocket::custom(Config::from_env())
        .attach(Databases::init(local_cfg.mongo_uri, local_cfg.redis_uri))
        .mount("/api", routes![api::health::route::check])
        .register("/", catchers![not_found])
}
