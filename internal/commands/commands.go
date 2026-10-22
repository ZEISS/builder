package commands

// Argument represents a command argument with its metadata.
type Argument struct {
	// ID is the unique identifier for the argument
	ID string
	// Title is the human-readable name of the argument
	Title string
	// Description is a brief description of the argument
	Description string
	// Required indicates whether the argument is required
	Required bool
}

// CustomCommand represents a user-defined custom command loaded from markdown files.
type CustomCommand struct {
	// ID is the unique identifier for the command
	ID string
	// Name is the human-readable name of the command
	Name string
	// Content is the command content as a string
	Content string
	// Arguments is a list of arguments for the command
	Arguments []Argument
}

type commandSource struct {
	path   string
	prefix string
}
