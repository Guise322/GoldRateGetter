package bot

type Command struct {
	Name        string
	Description string
	Action      func(msg Message) string
}
