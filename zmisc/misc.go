package zmisc

func FirstNonEmpty(ss ...string) string {
	var s string
	for _, s = range ss {
		if s != "" {
			break
		}
	}
	return s
}
