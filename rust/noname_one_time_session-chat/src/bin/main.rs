use noname_one_time_session_chat;

#[rocket::main]
async fn main() -> Result<(), rocket::Error> {
    let _rocket = noname_one_time_session_chat::rocket().launch().await?;
    Ok(())
}
