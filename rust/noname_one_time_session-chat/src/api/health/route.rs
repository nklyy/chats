use rocket::serde::json::{json, Value};

#[get("/check")]
pub fn check() -> Value {
    json!({"status": "OK"})
}
