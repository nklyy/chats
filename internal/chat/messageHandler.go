package chat

import (
	"context"
	"encoding/json"
	"noname-realtime-support-chat/internal/chat/room"
	"noname-realtime-support-chat/internal/chat/user"
	"time"
)

func (s *service) messageHandler(jsonMessage []byte) {
	var message room.Message
	if err := json.Unmarshal(jsonMessage, &message); err != nil {
		s.logger.Errorf("Error on unmarshal JSON message %s", err)
		return
	}

	switch message.Action {
	case "publish-room":
		dbUser, err := s.userSvc.GetUserByFingerprint(context.Background(), message.Fingerprint)
		if err != nil {
			s.logger.Errorf("failed to get user %v", err)
		}

		dbRoom, err := s.roomSvc.GetRoomByName(context.Background(), *dbUser.RoomName)
		if err != nil {
			s.logger.Errorf("failed to get room %v", err)
		}

		for r := range s.rooms {
			if r.Name == *dbUser.RoomName {
				var msg []*room.RoomMessage

				if dbRoom.Messages == nil {
					msg = append(msg, &room.RoomMessage{
						Id:      dbUser.ID,
						Time:    time.Now(),
						Message: message.Message,
					})
				} else {
					msg = append(*dbRoom.Messages, &room.RoomMessage{
						Id:      dbUser.ID,
						Time:    time.Now(),
						Message: message.Message,
					})
				}
				dbRoom.Messages = &msg

				////////
				err := s.roomSvc.UpdateRoom(context.Background(), dbRoom)
				if err != nil {

					r.Broadcast <- &room.BroadcastMessage{
						Action: message.Action,
						Message: room.MessageResponse{
							Action:  "",
							Message: nil,
							From:    "",
							Error:   "failed update room",
						},
						RoomName: *dbUser.RoomName,
					}
				}

				r.Broadcast <- &room.BroadcastMessage{
					Action: message.Action,
					Message: room.MessageResponse{
						Action:  message.Action,
						Message: &message.Message,
						From:    dbUser.ID,
						Error:   nil,
					},
					RoomName: *dbUser.RoomName,
				}
			}
		}
	case "disconnect":
		dbUser, err := s.userSvc.GetUserByFingerprint(context.Background(), message.Fingerprint)
		if err != nil {
			s.logger.Errorf("failed to get user %v", err)
		}

		for r := range s.rooms {
			if r.Name == *dbUser.RoomName {
				for client := range r.Clients {
					rUser, err := s.userSvc.GetUserById(context.Background(), client.Id)
					if err != nil {
						s.logger.Errorf("failed to get user %v", err)
					}

					userEntity, _ := user.MapToEntity(rUser)
					userEntity.SetRoom(nil)
					err = s.userSvc.UpdateUser(context.Background(), user.MapToDTO(userEntity))
					if err != nil {
						msg, _ := s.encodeMessage(room.MessageResponse{
							Action:  "",
							Message: nil,
							From:    "",
							Error:   "failed update user",
						})
						client.Send <- msg
						return
					}

					close(client.Send)
					err = client.Connection.Close()
					if err != nil {
						s.logger.Errorf("failed close connection %v", err)
					}

					for serverClient := range s.clients {
						if serverClient.Id == client.Id {
							delete(s.clients, client)
						}
					}
				}

				err := s.roomSvc.DeleteRoom(context.Background(), r.Name)
				if err != nil {
					s.logger.Errorf("failed ddelete room %v", err)
				}
				delete(s.rooms, r)
			}
		}
	}
}
