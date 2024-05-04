package v2

type Iqama interface {
	GetTodayTimes() (*IqamaDailyTimes, error)
	GetTomorrowTimes() (*IqamaDailyTimes, error)
	GetDiscordPrettified() string
	GetShellPrettified() string
}
