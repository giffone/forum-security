package config

const (
	Code200 = 200 // http.StatusOK (GET)
	Code201 = 201 // http.StatusCreated (POST)
	Code204 = 204 // http.StatusNoContent (PUT,PATCH,DELETE)
	Code301 = 301 // http.StatusMovedPermanently
	Code302 = 302 // http.StatusFound
	Code307 = 307 // http.StatusTemporaryRedirect
	Code400 = 400 // http.StatusBadRequest
	Code401 = 401 // http.StatusUnauthorized
	Code403 = 403 // http.StatusForbidden
	Code404 = 404 // http.StatusNotFound
	Code405 = 405 // http.StatusMethodNotAllowed
	Code422 = 422 // http.StatusUnprocessableEntity
	Code500 = 500 // http.StatusInternalServerError
	Code503 = 503 // http.StatusServiceUnavailable

	/*------------------------------------------------------*/

	StatusOK          = "successfully: %s"
	StatusCreated     = "created: %s"
	AlreadyExist      = "can not create: %s already have"
	InvalidCharacters = "invalid: the %s contains invalid characters"
	TooShort          = "too short: %s must be at least %s characters"
	NotMatch          = "no match: the entered %s does not match"
	WrongEnter        = "wrong: the entered %s is wrong"
	InvalidEnter      = "invalid: the entered %s is incorrect, please use valid"
	InvalidState      = "invalid: oauth state \"%s\" does not match with ours, url redirect for safety stopped"
	InternalError     = "internal error: \"%v\""
	AccessDenied      = "access denied: you not authorized or session expired"
	NotWorking        = "error: sorry, %s not working for now"
	ImageNotAllowed   = "image: file type %s can not use, allowed types (jpeg, png, gif)"
	FileSizeToBig     = "size: file size too big, accepted size up to %s"
)
