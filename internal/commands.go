package internal

type Сommands struct {
	Commands map[string]func(*State, Command) error
}

func (c *Сommands) Register(name string, f func(*State, Command) error) {
	c.Commands[name] = f
}

func (c *Сommands) Run(s *State, cmd Command) error {
	err := c.Commands[cmd.Name](s, cmd)
	if err != nil {
		return err
	}
	return nil
}
