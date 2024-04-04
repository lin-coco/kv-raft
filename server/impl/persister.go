package impl

import (
	"fmt"
	"io"
	"os"
)

type FilePersister struct {
	LogFile   *os.File
	StateFile *os.File
}

func NewFilePersister(logFile, stateFile *os.File) FilePersister {
	return FilePersister{LogFile: logFile, StateFile: stateFile}
}

func (f FilePersister) Persist(data []string) error {
	for _, d := range data {
		_, err := f.LogFile.Write([]byte(d + "\n"))
		if err != nil {
			return fmt.Errorf("f.File.Write err: %v", err)
		}
	}
	_ = f.LogFile.Sync()
	return nil
}

func (f FilePersister) ReplaceSnapshot(data []byte) error {
	if err := f.LogFile.Truncate(0); err != nil {
		return fmt.Errorf("f.File.Truncate err: %v", err)
	}
	_, err := f.LogFile.Write(data)
	if err != nil {
		return fmt.Errorf("f.File.Write err: %v", err)
	}
	_ = f.LogFile.Sync()
	return nil
}

func (f FilePersister) Snapshot() ([]byte, error) {
	all, err := io.ReadAll(f.LogFile)
	if err != nil {
		return nil, fmt.Errorf("io.ReadAll err: %v", err)
	}
	return all, nil
}

func (f FilePersister) SaveState(state []byte) error {
	if err := f.StateFile.Truncate(0); err != nil {
		return fmt.Errorf("f.File.Truncate err: %v", err)
	}
	// 将文件指针位置设置为文件开头
	if _, err := f.StateFile.Seek(0, 0); err != nil {
		return fmt.Errorf("f.File.Seek err: %v", err)
	}
	_, err := f.StateFile.Write(state)
	if err != nil {
		return fmt.Errorf("f.StateFile.Write err: %v", err)
	}
	return nil
}

func (f FilePersister) ReadState() ([]byte, error) {
	all, err := io.ReadAll(f.StateFile)
	if err != nil {
		return nil, fmt.Errorf("io.ReadAll err: %v", err)
	}
	return all, nil
}
