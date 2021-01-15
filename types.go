package coinbase

type userData struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	Username        string `json:"username"`
	ProfileLocation string `json:"profile_location"`
	ProfileBio      string `json:"profile_bio"`
	ProfileURL      string `json:"profile_url"`
	AvatarURL       string `json:"avatar_url"`
	Resource        string `json:"resource"`
	ResourcePath    string `json:"resource_path"`
	Email           string `json:"email"`
}

type user struct {
	Data userData `json:"data"`
}
