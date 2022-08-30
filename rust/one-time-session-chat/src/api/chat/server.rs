use std::collections::{HashMap, HashSet};

use actix::{Actor, Context, Handler, Message, Recipient};
use serde::{Deserialize, Serialize};
use uuid::Uuid;

/// Chat server sends this messages to session
#[derive(Message, Serialize, Deserialize)]
#[rtype(result = "()")]
pub struct SessionMessage {
    action: String,
    message: Option<String>,
    session_id: Option<String>,
    from: Option<String>,
}

#[derive(Message)]
#[rtype(result = "()")]
pub struct Publish {
    pub message: String,
    pub session_id: String,
}

#[derive(Message)]
#[rtype(String)]
pub struct Connect {
    pub addr: Recipient<SessionMessage>,
}

/// Session is disconnected
#[derive(Message)]
#[rtype(result = "()")]
pub struct Disconnect {
    pub session_id: String,
}

#[derive(Debug, Clone)]
pub struct ChatServer {
    pub sessions: HashMap<String, Recipient<SessionMessage>>,
    pub rooms: HashMap<String, HashSet<String>>,
}

impl ChatServer {
    pub fn new() -> ChatServer {
        ChatServer {
            sessions: HashMap::new(),
            rooms: HashMap::new(),
        }
    }
}

impl ChatServer {
    fn find_free_user(&mut self, current_user: String) -> Option<(&String, &mut HashSet<String>)> {
        for (room, room_users) in &mut self.rooms {
            let get_current_user = room_users.get(&current_user);
            if room_users.len() == 1 && get_current_user == None {
                return Some((room, room_users));
            }
        }
        None
    }
}

impl ChatServer {
    /// Send message to all users in the room
    fn send_disconnect_message(&self, room: &str) {
        if let Some(sessions) = self.rooms.get(room) {
            for id in sessions {
                if let Some(addr) = self.sessions.get(id) {
                    addr.do_send(SessionMessage {
                        action: "disconnected".to_string(),
                        message: None,
                        session_id: None,
                        from: None,
                    });
                }
            }
        }
    }

    fn send_connect_message(&self, room: &str) {
        if let Some(sessions) = self.rooms.get(room) {
            for id in sessions {
                if let Some(addr) = self.sessions.get(id) {
                    addr.do_send(SessionMessage {
                        action: "connected".to_string(),
                        message: None,
                        session_id: Some(id.to_string()),
                        from: None,
                    });
                }
            }
        }
    }

    fn send_publish_message(&self, room: &str, message: &str, from: &str) {
        if let Some(sessions) = self.rooms.get(room) {
            for id in sessions {
                if let Some(addr) = self.sessions.get(id) {
                    addr.do_send(SessionMessage {
                        action: "publish".to_string(),
                        message: Some(message.to_string()),
                        session_id: None,
                        from: Some(from.to_string()),
                    });
                }
            }
        }
    }
}

/// Make actor from `ChatServer`
impl Actor for ChatServer {
    /// We are going to use simple Context, we just need ability to communicate
    /// with other actors.
    type Context = Context<Self>;
}

/// Join room, send disconnect message to old room
/// send join message to new room
impl Handler<Publish> for ChatServer {
    type Result = ();

    fn handle(&mut self, msg: Publish, _: &mut Context<Self>) {
        for (name, sessions) in &self.rooms {
            let found_user = sessions.get(&msg.session_id);
            if found_user != None {
                self.send_publish_message(&name, &msg.message, &msg.session_id)
            }
        }
    }
}

/// Handler for Connect message.
/// Register new session and assign unique id to this session
impl Handler<Connect> for ChatServer {
    type Result = String;

    fn handle(&mut self, msg: Connect, _: &mut Context<Self>) -> Self::Result {
        println!("Someone joined");

        // register session with random id
        let id = Uuid::new_v4();
        self.sessions.insert(id.to_string(), msg.addr);

        match self.find_free_user(id.to_string().to_owned()) {
            Some((room_name, room_sessions)) => {
                room_sessions.insert(id.to_string().to_owned());
                // println!("{}", room_name);
                let room_name = room_name.clone();
                self.send_connect_message(&room_name);
            }
            None => {
                self.rooms
                    .entry(Uuid::new_v4().to_string())
                    .or_insert_with(HashSet::new)
                    .insert(id.to_string());
            }
        }

        id.to_string()
    }
}

/// Handler for Disconnect message.
impl Handler<Disconnect> for ChatServer {
    type Result = ();

    fn handle(&mut self, msg: Disconnect, _: &mut Context<Self>) {
        println!("Someone disconnected");

        let mut rooms: Vec<String> = Vec::new();

        // remove address
        if self.sessions.remove(&msg.session_id).is_some() {
            // remove session from all rooms
            for (name, sessions) in &mut self.rooms {
                if sessions.remove(&msg.session_id) {
                    rooms.push(name.to_owned());
                }
            }
        }

        for room in rooms {
            self.send_disconnect_message(&room);
            self.rooms.remove(&room);
        }
    }
}
