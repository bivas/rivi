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

func (s *StringSet) Remove(remove string) bool {
	if s.items == nil {
		return false
	}
	item := remove
	if s.Transformer != nil {
		item = s.Transformer(remove)
	}
	_, exists := s.items[item]
	if exists {
		delete(s.items, item)
		return true
	}
	return false
}

func (s *StringSet) Values() []string {
	result := make([]string, len(s.items))
	i := 0
	for item := range s.items {
		result[i] = item
		i++
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
