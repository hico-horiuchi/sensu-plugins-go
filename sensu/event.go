package sensu

type EventStruct struct {
	Client      clientStruct
	Check       checkStruct
	Occurrences int
	Action      string
}
