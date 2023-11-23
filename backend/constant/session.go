package constant

const (
	AuthMethodSession = "session"
	SessionName       = "psession"

	AuthMethodJWT = "jwt"
	JWTHeaderName = "PanelAuthorization"
	JWTBufferTime = 3600
	JWTIssuer     = "PanelX"

	PasswordExpiredName = "expired"
)
