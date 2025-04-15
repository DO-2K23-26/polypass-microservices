package health

const (
	Ok    = "Ok"
	Failure = "Failure"
)

type HealthResponse map[string]string
