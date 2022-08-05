use actix_web::HttpResponse;
use serde::Serialize;

pub mod health;

#[derive(Serialize)]
struct NotFoundResponse {
    error: String,
}

pub async fn not_found() -> HttpResponse {
    let resp = NotFoundResponse {
        error: "Page NotFound".to_string(),
    };

    HttpResponse::NotFound().json(resp)
}
