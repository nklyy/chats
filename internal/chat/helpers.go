package chat

import (
	"encoding/base64"
	"encoding/json"
	"golang.org/x/crypto/scrypt"
	"noname-realtime-support-chat/internal/chat/room"
)

func (s *service) encodeMessage(msg room.MessageResponse) ([]byte, error) {
	encMsg, err := json.Marshal(msg)
	if err != nil {
		s.logger.Errorf("failed to encode message %v", err)
		return nil, err
	}
	return encMsg, err
}

func (s *service) createHash(data string) (string, error) {
	hash, err := scrypt.Key([]byte(data), []byte(s.salt), 16384, 8, 1, 32)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(hash), nil
}

func (s *service) cleanupOldConnections(userFingerprint string) {
	for client := range s.clients {
		if client.Fingerprint == userFingerprint {
			client.Connection.Close()
			delete(s.clients, client)
		}
	}
}
