package model

type Dialect struct {
	Name  *string `json:"name,omitempty"`  // The name of the dialect
	Label *string `json:"label,omitempty"` // The human-readable label of the connection
}

type Connection struct {
	Name        *string  `json:"name,omitempty"` // Name of the connection. Also used as the unique identifier
	Dialect     *Dialect `json:"dialect,omitempty"`
	Host        *string  `json:"host,omitempty"`         // Host name/address of server
	Port        *string  `json:"port,omitempty"`         // Port number on server
	Database    *string  `json:"database,omitempty"`     // Database name
	Schema      *string  `json:"schema,omitempty"`       // Scheme name
	DialectName *string  `json:"dialect_name,omitempty"` // (Read/Write) SQL Dialect name
	Managed     *bool    `json:"managed,omitempty"`      // Is this connection created and managed by Looker
}
