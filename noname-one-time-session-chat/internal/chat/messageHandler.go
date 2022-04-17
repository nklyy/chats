package chat

import (
	"context"
	"encoding/json"
	"noname-realtime-support-chat/internal/chat/room"
)

func (s *service) messageHandler(jsonMessage []byte) {
	var message room.Message
	if err := json.Unmarshal(jsonMessage, &message); err != nil {
		s.logger.Errorf("Error on unmarshal JSON message %s", err)
		return
	}

	switch message.Action {
	case "publish-room":
		var serverUser *room.Client
		for client := range s.clients {
			if client.Fingerprint == message.Fingerprint {
				serverUser = client
			}
		}

		for r := range s.rooms {
			if r.Name == serverUser.Room.Name {
				r.Broadcast <- &room.BroadcastMessage{
					Action: message.Action,
					Message: room.MessageResponse{
						Action:  message.Action,
						Message: &message.Message,
						From:    message.Fingerprint,
						Error:   nil,
					},
					RoomName: serverUser.Room.Name,
				}
			}
		}
	case "disconnect":
		var serverUser *room.Client
		for client := range s.clients {
			if client.Fingerprint == message.Fingerprint {
				serverUser = client
			}
		}

		for r := range s.rooms {
			if r.Name == serverUser.Room.Name {
				for roomClient := range r.Clients {
					//roomClient.Room = nil
					//close(roomClient.Send)

					err := roomClient.PubSub.Unsubscribe(context.Background(), roomClient.Room.Name)
					if err != nil {
						s.logger.Errorf("failed unsubscrube from channel %v", err)
					}

					//roomClient.PubSub.Close()
					roomClient.Connection.Close()

					for serverClient := range s.clients {
						if serverClient.Fingerprint == roomClient.Fingerprint {
							delete(s.clients, roomClient)
						}
					}
				}
				delete(s.rooms, r)
			}
		}
	}
}
