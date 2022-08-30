use actix_web::{get, web, HttpResponse};
use serde::Serialize;

use crate::databases::Databases;

#[derive(Serialize)]
struct CheckResponse {
    status: String,
}

#[get("/check")]
pub async fn check(_dbs: web::Data<Databases>) -> HttpResponse {
    let resp = CheckResponse {
        status: "OK".to_string(),
    };

    HttpResponse::Ok().json(resp)
}
