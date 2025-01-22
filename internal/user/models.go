package user

type NewUserRequest struct {
	Name                   string   `json:"name"`
	Email                  string   `json:"email"`
	PreferredNotifications []string `json:"preferredNotifications"`
}

type User struct {
	ID                     int      `json:"id"`
	Name                   string   `json:"name"`
	Email                  string   `json:"email"`
	PreferredNotifications []string `json:"preferredNotifications"`
}

type UpdatePreferencesRequest struct {
	UserID                 int      `json:"userId"`
	PreferredNotifications []string `json:"preferredNotifications"`
}

type GetUsersByPreferencesRequest struct {
	PreferredNotifications []string `json:"preferredNotifications"`
}
