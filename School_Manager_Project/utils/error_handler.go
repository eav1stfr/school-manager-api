package utils

import "net/http"

type AppErrors struct {
	errMessage string
	statusCode int
}

func (e *AppErrors) Error() string {
	return e.errMessage
}

func (e *AppErrors) GetStatusCode() int {
	return e.statusCode
}

func (e *AppErrors) SetErrorMessage(errMsg string) {
	e.errMessage = errMsg
}

func (e *AppErrors) SetErrStatusCode(code int) {
	e.statusCode = code
}

var (
	InvalidSortParameterError = &AppErrors{
		errMessage: "invalid sort filter parameter",
		statusCode: http.StatusBadRequest}

	ConnectingToDatabaseError = &AppErrors{
		errMessage: "error connecting to database",
		statusCode: http.StatusInternalServerError}

	DatabaseQueryError = &AppErrors{
		errMessage: "error database query",
		statusCode: http.StatusInternalServerError}

	UnitNotFoundError = &AppErrors{
		errMessage: "unit not found",
		statusCode: http.StatusNotFound}

	UnableToStartTransactionError = &AppErrors{
		errMessage: "unable to start transaction",
		statusCode: http.StatusInternalServerError}

	InvalidIdError = &AppErrors{
		errMessage: "invalid ID",
		statusCode: http.StatusBadRequest}

	InvalidUpdateParametersError = &AppErrors{
		errMessage: "invalid update parameters",
		statusCode: http.StatusBadRequest}

	ErrorCommitingTransaction = &AppErrors{
		errMessage: "error commiting the transaction",
		statusCode: http.StatusInternalServerError}

	MissingFieldsError = &AppErrors{
		errMessage: "invalid request body - all fields are required",
		statusCode: http.StatusBadRequest}

	DuplicateEmailError = &AppErrors{
		errMessage: "duplicate email - email must be unique",
		statusCode: http.StatusBadRequest}

	ClassTeacherNotFound = &AppErrors{
		errMessage: "class / class teacher not found",
		statusCode: http.StatusBadRequest}

	ErrorEncodingData = &AppErrors{
		errMessage: "error encoding data",
		statusCode: http.StatusInternalServerError}

	ErrorGeneratingSaltForHashing = &AppErrors{
		errMessage: "error hashing password",
		statusCode: http.StatusInternalServerError}

	InvalidRequestBodyError = &AppErrors{
		errMessage: "invalid request body",
		statusCode: http.StatusBadRequest}

	AccountInactiveError = &AppErrors{
		errMessage: "account is inactive",
		statusCode: http.StatusForbidden}

	InvalidEncodedHashFormat = &AppErrors{
		errMessage: "invalid encoded hash format",
		statusCode: http.StatusForbidden}

	FailedToDecodeSalt = &AppErrors{
		errMessage: "failed to decode the salt",
		statusCode: http.StatusForbidden}

	FailedToDecodeHashError = &AppErrors{
		errMessage: "failed to decode the hashed password",
		statusCode: http.StatusForbidden}

	IncorrectPasswordError = &AppErrors{
		errMessage: "incorrect password",
		statusCode: http.StatusForbidden}

	ErrorGeneratingJwtToken = &AppErrors{
		errMessage: "error generating jwt token",
		statusCode: http.StatusInternalServerError}

	UnknownInternalServerError = &AppErrors{
		errMessage: "unknown internal server error",
		statusCode: http.StatusInternalServerError}

	TokenExpiredError = &AppErrors{
		errMessage: "token is expired",
		statusCode: http.StatusUnauthorized}

	InvalidLoginTokenError = &AppErrors{
		errMessage: "invalid login token",
		statusCode: http.StatusUnauthorized}

	UnexpectedSigningMethodError = &AppErrors{
		errMessage: "unexpected signing method",
		statusCode: http.StatusUnauthorized}

	UserNotAuthorizedError = &AppErrors{
		errMessage: "user not authorized",
		statusCode: http.StatusUnauthorized}
)
