package token

type Maker interface {
	// CreateToken creates a new token for a specific username and duration
	CreateToken(tokenParams *TokenParams) (string, *Payload, error)

	//  VerifyToken checks if the input token is valid or not
	VerifyToken(token string) (*Payload, error)
}
