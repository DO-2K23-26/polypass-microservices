package search

type SearchRepository interface {
	SearchCredentials(query SearchCredentialQuery) (*SearchCredentialResult, error)
}
