package Configuration

type Settings struct {
	// I recognize that the Username / Authtoken should not be
	// in plain text in a public git repo for everyone to see, but this
	// way for demo purposes, it is easiest so that anyone can pull it
	// down and have it immediately work
	UserName                       string
	AuthToken                      string
	Channel                        string
	CardCommand                    string
	DisconnectCommand              string
	ChangeChannelCommand           string
	TwitchServer                   string
	TwitchPort                     string
	TwitchRateLimit                int
	TwitchRateLimitDurationSeconds int
	MagicEndpoint                  string
	MagicRateLimit                 int
	MagicRateLimitDurationSeconds  int
	VerboseMode                    bool
}
