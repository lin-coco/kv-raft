package command

import "errors"

type ImportCMD struct {
	Filepaths []string
}

func importCheck(split []string) error {
	if len(split) != 1 {
		return errors.New("command is incorrent")
	}
	return nil
}

func importUnmarshal(split []string) ImportCMD {
	return ImportCMD{
		Filepaths: split,
	}
}

func (g ImportCMD) Marshal() error {
	return nil
}
func (g ImportCMD) GetKeys() []string {
	return nil
}
func (g ImportCMD) ExecCMD() string {
	return ""
}
func (g ImportCMD) ReadOnly() bool {
	return false
}