package npvm

func GetAllRunningFiles() []string {
	mu.Lock()
	defer mu.Unlock()
	t := runningFiles
	return t
}
