// interfaces/credentials/embed.go
package avro    // nom de package au choix

import "embed"

//go:embed *.avsc 
var FS embed.FS