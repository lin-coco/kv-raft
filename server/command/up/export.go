package command

import "errors"

type ExportCMD struct {
	Fuzzys []string
}

func exportCheck(split []string) error {
	if len(split) != 1 {
		return errors.New("command is incorrent")
	}
	return nil
}

func exportUnmarshal(split []string) ExportCMD {
	return ExportCMD{
		Fuzzys: split,
	}
}

func (g ExportCMD) Marshal() error {
	return nil
}
func (g ExportCMD) GetKeys() []string {
	return nil
}
func (g ExportCMD) ExecCMD() string {
	return ""
}
func (g ExportCMD) ReadOnly() bool {
	return false
}