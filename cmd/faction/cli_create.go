package main

type cliCreateCmd struct {
	optCliConn

	Object struct {
		Name string `positional-arg-name:"object" description:"Object to create"`
	} `positional-args:"true" required:"true"`
}

func (c *cliCreateCmd) Execute(args []string) error {
	return nil
}
