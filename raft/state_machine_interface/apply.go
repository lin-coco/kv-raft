package state_machine_interface

type Apply interface {
	ApplyCommand(command string) string
}
