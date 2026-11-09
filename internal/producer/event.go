// Package producer struct for events
package producer

type (
	Events struct {
		EventType string `json:"event_type"`
		EventID   string `json:"event_id"`
		Payload   any    `json:"payload"`
	}

	 APIResponse struct {
		Status  int    `json:"status"`
		Message string `json:"mssage"`
		Data    any    `json:"data"`
	}
) 
