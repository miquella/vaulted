package main

type Shell struct {
	VaultName string
	Command   []string
}

func (s *Shell) Run(steward Steward) error {
	_, env, err := steward.GetEnvironment(s.VaultName, nil)
	if err != nil {
		return err
	}

	code, err := env.Spawn(s.Command, nil)
	if err != nil {
		return ErrorWithExitCode{err, 2}
	} else if *code != 0 {
		return ErrorWithExitCode{ErrNoError, *code}
	}

	return nil
}
