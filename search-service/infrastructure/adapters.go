package infrastructure

// The "infrastructure" package contains adapters for various infrastructure repositories, such as database repositories.
// Below are defined common interfaces for such adapters.

type IDbAdapter interface {
	CheckHealth() bool
	Migrate() error
}

type ISqlAdapter interface {
	Get()
	Set(...)
}

type IGormAdapter interface {
	IDbAdapter
	ISqlAdapter
}
