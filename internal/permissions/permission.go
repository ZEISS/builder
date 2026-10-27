package permissions

import (
	"github.com/google/uuid"
)

// PermissionRequest represents a request for a permission.
type PermissionRequest struct {
	// ID is the unique identifier for the permission request.
	ID uuid.UUID `json:"id"`
	// SessionID is the unique identifier for the session.
	SessionID uuid.UUID `json:"session_id"`
	// ToolCallID is the unique identifier for the tool call.
	ToolCallID string `json:"tool_call_id"`
	// ToolName is the name of the tool.
	ToolName string `json:"tool_name"`
	// Description is the description of the permission request.
	Description string `json:"description"`
	// Action is the action to be performed.
	Action string `json:"action"`
	// Params are the parameters for the action.
	Params any `json:"params"`
	// Path is the path for the action.
	Path string `json:"path"`
}
