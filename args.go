package main

func ParseArgs(args []string) (Command, error) {
	if len(args) == 0 {
		return nil, nil
	}

	switch args[0] {
	default:
		return nil, nil
	}
}
