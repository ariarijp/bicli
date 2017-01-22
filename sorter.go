package main

func (s ShortURLs) Len() int {
	return len(s)
}

func (s ShortURLs) Less(i, j int) bool {
	return s[i].LineNum < s[j].LineNum
}

func (s ShortURLs) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
