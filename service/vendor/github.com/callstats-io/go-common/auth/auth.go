package auth

// SigningMethodFromString returns the signing method corresponding to a constant string
func SigningMethodFromString(method string) SigningMethod {
	switch method {
	case "HS256":
		return SigningMethodHS256
	case "HS384":
		return SigningMethodHS384
	case "HS512":
		return SigningMethodHS512
	case "ES256":
		return SigningMethodES256
	case "ES384":
		return SigningMethodES384
	case "ES512":
		return SigningMethodES512
	default:
		panic("Unknown signing method: " + method)
	}
}
