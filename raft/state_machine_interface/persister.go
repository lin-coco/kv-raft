package state_machine_interface

type Persister interface {
	Persist(data []string) error       // 数据应该是递增的存储
	ReplaceSnapshot(data []byte) error // 替换整个快照
	Snapshot() ([]byte, error)         // 获取整个快照
	SaveState(state []byte) error      // 设置状态
	ReadState() ([]byte, error)        // 阅读状态
}

//type FilePersister struct {
//	filePath         string
//	metaDataFilePath string
//}
//
//func NewFilePersister(filePath string) *FilePersister {
//	return &FilePersister{
//		filePath:         filePath,
//		metaDataFilePath: filePath + ".meta",
//	}
//}
//
//func (p *FilePersister) SaveState(meta []byte) error {
//	if err := os.WriteFile(p.metaDataFilePath, meta, 0644); err != nil {
//		return fmt.Errorf("os.WriteFile err: %v", err)
//	}
//	return nil
//}
//
//func (p *FilePersister) ReadState() ([]byte, error) {
//	return os.ReadFile(p.metaDataFilePath)
//}
//
//func (p *FilePersister) Persist(data [][]byte) error {
//	// 打开文件，以追加模式写入
//	file, err := os.OpenFile(p.filePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
//	if err != nil {
//		return fmt.Errorf("os.OpenFile err: %v", err)
//	}
//	defer func() {
//		_ = file.Close()
//	}()
//	var allBytes []byte
//	for i := 0; i < len(data); i++ {
//		allBytes = append(allBytes, data[i]...)
//		allBytes = append(allBytes, '\n') // 添加换行符
//	}
//	// 写入数据到文件
//	_, err = file.Write(allBytes)
//	if err != nil {
//		return fmt.Errorf("file.Write err: %v", err)
//	}
//	return nil
//}
//
//func (p *FilePersister) ReadPersist() ([][]byte, error) {
//	var lines [][]byte
//
//	// 打开文件进行读取
//	file, err := os.Open(p.filePath)
//	if err != nil {
//		return nil, fmt.Errorf("os.Open err: %v", err)
//	}
//	defer func() {
//		_ = file.Close()
//	}()
//
//	// 逐行读取文件内容
//	scanner := bufio.NewScanner(file)
//	for scanner.Scan() {
//		lines = append(lines, scanner.Bytes())
//	}
//	if err = scanner.Err(); err != nil {
//		return nil, fmt.Errorf("bufio.NewScannerr err: %v", err)
//	}
//	return lines, nil
//}
