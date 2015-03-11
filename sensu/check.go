package sensu

type checkStruct struct {
	Name        string
	Issued      int
	Output      string
	Status      int
	Command     string
	Subscribers []string
	Handler     string
	History     []string
	Flapping    bool
}
