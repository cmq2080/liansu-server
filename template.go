package main

type Template struct {
	data map[string]string
}

func (t *Template) GetKeys() []string {
	keys := []string{}
	for key := range t.data {
		keys = append(keys, key)
	}

	return keys
}
