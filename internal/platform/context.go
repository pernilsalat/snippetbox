package platform

type contextKey string

const IsAuthenticatedContextKey = contextKey("isAuthenticated")

const UserContextKey = contextKey("user")
