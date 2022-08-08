mod api;
mod config;
mod databases;

use actix_web::{middleware::Logger, web, App, HttpServer};
use config::Config;
use databases::Databases;

pub async fn execute() -> Result<(), std::io::Error> {
    // let local_cfg = Config::new().unwrap_or_else(|err| {
    //     eprintln!("[ERROR] Problem parsing env arguments: {}", err);
    //     process::exit(1);
    // });
    std::env::set_var("RUST_LOG", "info,actix_web=info");
    env_logger::init();

    let local_cfg = Config::init();
    let dbs = Databases::init(local_cfg.mongo_uri, local_cfg.redis_uri).await;

    HttpServer::new(move || {
        App::new()
            .wrap(Logger::default())
            .app_data(web::Data::new(dbs.clone()))
            .service(api::chat::service::chat)
            .service(
                web::scope("/api").service(api::health::service::check), // ...so this handles requests for `GET /app/index.html`
            )
            .default_service(web::to(api::not_found))
    })
    .bind(("127.0.0.1", local_cfg.port))?
    .run()
    .await
}
