package mapping

type MappedResultSet struct {
	results map[string]interface{}
}

var _ ResultSet = &MappedResultSet{}

func NewMappedResultSet() *MappedResultSet {
	return &MappedResultSet{
		results: make(map[string]interface{}),
	}
}

func (m *MappedResultSet) AddResult(key string, result interface{}) {
	m.results[key] = result
}

func (m *MappedResultSet) GetResult(key string) (interface{}, bool) {
	result, ok := m.results[key]
	return result, ok
}

func (m *MappedResultSet) GetAllResults() map[string]interface{} {
	var results = make(map[string]interface{})

	for key, result := range m.results {
		results[key] = result
	}

	return results
}
