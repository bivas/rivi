package util

type StringSet struct {
	items       map[string]bool
	Transformer func(string) string
}

func (s *StringSet) Add(add string) bool {
	item := add
	if s.Transformer != nil {
		item = s.Transformer(add)
	}
	if s.items == nil {
		s.items = make(map[string]bool)
	}
	_, exists := s.items[item]
	if exists {
		return false
	}
	s.items[item] = true
	return true
}

func (s *StringSet) AddAll(add []string) *StringSet {
	for _, item := range add {
		s.Add(item)
	}
	return s
}

func (s *StringSet) Values() []string {
	result := make([]string, 0)
	for item := range s.items {
		result = append(result, item)
	}
	return result
}

func StringSliceContains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
