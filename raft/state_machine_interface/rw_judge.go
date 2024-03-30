package state_machine_interface

type RWJudge interface {
	ReadOnly(command string) bool
}
