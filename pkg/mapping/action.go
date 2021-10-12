package mapping

type CompositeMappingAction struct {
	Actions []MappingAction
}

var _ MappingAction = &CompositeMappingAction{}

func ActionComposedOf(actions ...MappingAction) *CompositeMappingAction {
	return &CompositeMappingAction{
		Actions: actions,
	}
}

func (a *CompositeMappingAction) Apply(source interface{}, resultSet ResultSet) error {
	for _, action := range a.Actions {
		if err := action.Apply(source, resultSet); err != nil {
			return err
		}
	}

	return nil
}
