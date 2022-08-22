use std::time::{Duration, Instant};

use actix::{
    fut, Actor, ActorContext, ActorFutureExt, Addr, AsyncContext, ContextFutureSpawner, Handler,
    Running, StreamHandler, WrapFuture,
};
use actix_web::{get, web, Error, HttpRequest, HttpResponse};
use actix_web_actors::ws;
use actix_web_actors::ws::CloseReason;
use log::{error, info};
use serde::Deserialize;

use crate::api::chat::server;

/// How often heartbeat pings are sent
const HEARTBEAT_INTERVAL: Duration = Duration::from_secs(5);

/// How long before lack of client response causes a timeout
const CLIENT_TIMEOUT: Duration = Duration::from_secs(10);

#[derive(Deserialize)]
struct ChatMessage {
    action: String,
    message: String,
}

#[derive(Debug, Clone)]
pub struct Session {
    /// Client must send ping at least once per 10 seconds (CLIENT_TIMEOUT),
    /// otherwise we drop connection.
    pub hb: Instant,
    // pub session_id: Option<String>,
    pub session_id: String,
    pub addr: Addr<server::ChatServer>,
}

impl Session {
    /// helper method that sends ping to client every 5 seconds (HEARTBEAT_INTERVAL).
    /// also this method checks heartbeats from client
    fn hb(&self, ctx: &mut ws::WebsocketContext<Self>) {
        ctx.run_interval(HEARTBEAT_INTERVAL, |act, ctx| {
            // check client heartbeats
            if Instant::now().duration_since(act.hb) > CLIENT_TIMEOUT {
                // heartbeat timed out
                println!("Websocket Client heartbeat failed, disconnecting!");

                // notify chat server
                act.addr.do_send(server::Disconnect {
                    session_id: act.session_id.to_owned(),
                });

                // stop actor
                ctx.stop();

                // don't try to send a ping
                return;
            }

            // println!("ping");
            ctx.ping(b"ping");
        });
    }
}

impl Actor for Session {
    type Context = ws::WebsocketContext<Self>;

    /// Method is called on actor start.
    /// We register ws session with ChatServer
    fn started(&mut self, ctx: &mut Self::Context) {
        // we'll start heartbeat process on session start.
        self.hb(ctx);

        // register self in chat server. `AsyncContext::wait` register
        // future within context, but context waits until this future resolves
        // before processing any other events.
        // HttpContext::state() is instance of WsChatSessionState, state is shared
        // across all routes within application
        let addr = ctx.address();

        self.addr
            .send(server::Connect {
                addr: addr.recipient(),
            })
            .into_actor(self)
            .then(|res, act, ctx| {
                match res {
                    Ok(res) => act.session_id = res,

                    // something is wrong with chat server
                    _ => ctx.stop(),
                }

                fut::ready(())
            })
            .wait(ctx)
    }

    // delete all rooms with current clients
    fn stopping(&mut self, _: &mut Self::Context) -> Running {
        // notify chat server
        self.addr.do_send(server::Disconnect {
            session_id: self.session_id.to_owned(),
        });
        Running::Stop
    }
}

/// Handle messages from chat server, we simply send it to peer websocket
impl Handler<server::SessionMessage> for Session {
    type Result = ();

    fn handle(&mut self, msg: server::SessionMessage, ctx: &mut Self::Context) {
        // let msg = match serde_json::to_string(&msg) {
        //     Ok(r) => r,
        //     Err(_) => {
        //         ctx.close(Some(CloseReason {
        //             code: ws::CloseCode::Invalid,
        //             description: Some("message json".to_string()),
        //         }));
        //         return ctx.stop();
        //     }
        // };
        let msg = serde_json::to_string(&msg).unwrap();
        ctx.text(msg);
    }
}

/// Handler for ws::Message message
impl StreamHandler<Result<ws::Message, ws::ProtocolError>> for Session {
    fn handle(&mut self, msg: Result<ws::Message, ws::ProtocolError>, ctx: &mut Self::Context) {
        match msg {
            Ok(ws::Message::Ping(msg)) => {
                self.hb = Instant::now();
                ctx.pong(&msg);
            }
            Ok(ws::Message::Pong(_)) => {
                self.hb = Instant::now();
            }
            Ok(ws::Message::Text(text)) => {
                let msg: ChatMessage = match serde_json::from_str(&text) {
                    Ok(r) => r,
                    Err(err) => {
                        error!("invalid request json: {}", err);
                        ctx.close(Some(CloseReason {
                            code: ws::CloseCode::Invalid,
                            description: Some("invalid request json".to_string()),
                        }));
                        return ctx.stop();
                    }
                };

                if msg.action == "publish-room" {
                    self.addr.do_send(server::Publish {
                        message: msg.message.to_owned(),
                        session_id: self.session_id.to_owned(),
                    });
                }

                if msg.action == "disconnect" {
                    return ctx.stop();
                }

                if msg.action != "publish-room" || msg.action != "disconnect" {
                    info!("unsupported action: {}", msg.action);
                    ctx.close(Some(CloseReason {
                        code: ws::CloseCode::Unsupported,
                        description: Some("unsupported action".to_string()),
                    }));
                    ctx.stop()
                }

                // ctx.text(text)
            }
            // Ok(ws::Message::Binary(bin)) => ctx.binary(bin),
            Ok(ws::Message::Close(reason)) => {
                ctx.close(reason);
                ctx.stop();
            }
            Ok(ws::Message::Continuation(_)) => {
                ctx.stop();
            }
            Ok(ws::Message::Nop) => (),
            _ => (),
        }
    }
}

#[get("/chat")]
async fn chat(
    req: HttpRequest,
    stream: web::Payload,
    srv: web::Data<Addr<server::ChatServer>>,
) -> Result<HttpResponse, Error> {
    let client = Session {
        hb: Instant::now(),
        session_id: "".to_string(),
        addr: srv.get_ref().clone(),
    };

    // let query_params = match web::Query::<ChatRequest>::from_query(req.query_string()) {
    //     Ok(params) => params,
    //     Err(err) => return Ok(HttpResponse::BadRequest().body("invalid query param")),
    // };

    // println!("{}", query_params.fingerprint);

    // let resp = ws::start(client, &req, stream);
    // println!("{:?}", resp);
    // resp

    ws::start(client, &req, stream)
}
