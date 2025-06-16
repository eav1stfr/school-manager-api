package utils

func AuthorizeUser(role string, allowedRoles ...string) (bool, error) {
	for _, allowedRole := range allowedRoles {
		if role == allowedRole {
			return true, nil
		}
	}
	return false, UserNotAuthorizedError
}
