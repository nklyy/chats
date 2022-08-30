package chat

import (
	"context"
	"encoding/json"
	"one-time-session-chat/internal/chat/room"
)

func (s *service) encodeMessage(msg room.MessageResponse) ([]byte, error) {
	encMsg, err := json.Marshal(msg)
	if err != nil {
		s.logger.Errorf("failed to encode message %v", err)
		return nil, err
	}
	return encMsg, err
}

func (s *service) cleanupOldConnections(userFingerprint string) {
	for client := range s.clients {
		if client.Fingerprint == userFingerprint {
			if client.Room != nil {
				for rm := range s.rooms {
					if rm.Name == client.Room.Name {
						for roomClient := range rm.Clients {
							err := roomClient.PubSub.Unsubscribe(context.Background(), roomClient.Room.Name)
							if err != nil {
								s.logger.Errorf("failed unsubscrube from channel %v", err)
							}

							//roomClient.PubSub.Close()
							err = roomClient.Connection.Close()
							if err != nil {
								s.logger.Errorf("failed to close room client connection %v", err)
							}
							delete(s.clients, roomClient)
						}
						delete(s.rooms, rm)
					}
				}
			}

			err := client.Connection.Close()
			if err != nil {
				s.logger.Errorf("failed to close client connection %v", err)
			}
			delete(s.clients, client)
		}
	}
}
