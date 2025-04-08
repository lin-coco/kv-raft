package storage

import "strings"


type Gomap map[string]string

func NewGomap(initCap int) *Gomap {
	gomap := Gomap(make(map[string]string, initCap))
	return &gomap
}

func (m *Gomap) Put(key, value string) {
	(*m)[key] = value
}

func (m *Gomap) Get(key string) (string,bool) {
	v,b := (*m)[key]
	return v,b
}

func (m *Gomap) Del(key string) {
	delete(*m, key)
}

func (m *Gomap) Prefix(prefix string) []string {
	result := make([]string, 0)
	for k := range *m {
		if k == "" {
			continue
		}
		if strings.HasPrefix(k, prefix) {
			result = append(result, k)
		}
	}
	return result
}
func (m *Gomap)	Suffix(suffix string) []string {
	result := make([]string, 0)
	for k := range *m {		
		if k == "" {
			continue
		}
		if strings.HasSuffix(k, suffix) {
			result = append(result, k)
		}
	}
	return result

}
func (m *Gomap)	Contains(sub string) []string {
	result := make([]string, 0)
	for k := range *m {		
		if k == "" {
			continue
		}
		if strings.Contains(k, sub) {
			result = append(result, k)
		}
	}
	return result
}