package mapping

func MapObject(actionResolver MappingActionResolver, resultSet ResultSet, source interface{}) error {
	actions, err := actionResolver.ResolveMappingActions(source)
	if err != nil {
		return err
	}

	action := ActionComposedOf(actions...)
	if err := action.Apply(source, resultSet); err != nil {
		return err
	}

	return nil
}
