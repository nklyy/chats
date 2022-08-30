use one_time_session_chat;

#[actix_web::main]
async fn main() -> std::io::Result<()> {
    one_time_session_chat::execute().await
}
