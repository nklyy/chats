use actix::{Actor, ActorContext, StreamHandler};
use actix_web::{get, web, Error, HttpRequest, HttpResponse};
use actix_web_actors::ws::{self, CloseReason};
use serde::Deserialize;

#[derive(Deserialize)]
struct ChatMessage {
    action: String,
    message: EncryptedMessage,
    fingerprint: String,
}

#[derive(Deserialize)]
struct EncryptedMessage {
    data: String,
    salt: String,
    iv: String,
}

struct Chat;

impl Actor for Chat {
    type Context = ws::WebsocketContext<Self>;
}

/// Handler for ws::Message message
impl StreamHandler<Result<ws::Message, ws::ProtocolError>> for Chat {
    fn handle(&mut self, msg: Result<ws::Message, ws::ProtocolError>, ctx: &mut Self::Context) {
        match msg {
            Ok(ws::Message::Ping(msg)) => ctx.pong(&msg),
            Ok(ws::Message::Text(text)) => {
                let msg: ChatMessage = match serde_json::from_str(&text) {
                    Ok(r) => r,
                    Err(_) => {
                        ctx.close(Some(CloseReason {
                            code: ws::CloseCode::Invalid,
                            description: Some("invalid request json".to_string()),
                        }));
                        return ctx.stop();
                    }
                };

                // println!("{:?}", msg.message.data);
                ctx.text(text)
            }
            // Ok(ws::Message::Binary(bin)) => ctx.binary(bin),
            _ => (),
        }
    }
}

#[get("/chat")]
async fn chat(req: HttpRequest, stream: web::Payload) -> Result<HttpResponse, Error> {
    let resp = ws::start(Chat {}, &req, stream);
    println!("{:?}", resp);
    resp
}
