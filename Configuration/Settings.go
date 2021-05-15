package Configuration

type Settings struct {
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
}
