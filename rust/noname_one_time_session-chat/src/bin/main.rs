use noname_one_time_session_chat;

#[actix_web::main]
async fn main() -> std::io::Result<()> {
    noname_one_time_session_chat::execute().await
}
