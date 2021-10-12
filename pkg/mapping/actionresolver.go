package mapping

type CompositeActionResolver struct {
	Resolvers []MappingActionResolver
}

var _ MappingActionResolver = &CompositeActionResolver{}

func ActionResolverComposedOf(resolvers ...MappingActionResolver) *CompositeActionResolver {
	return &CompositeActionResolver{
		Resolvers: resolvers,
	}
}

func (r *CompositeActionResolver) ResolveMappingActions(source interface{}) ([]MappingAction, error) {
	actions := []MappingAction{}

	for _, resolver := range r.Resolvers {
		resolvedActions, err := resolver.ResolveMappingActions(source)
		if err != nil {
			return nil, err
		}

		actions = append(actions, resolvedActions...)
	}

	return actions, nil
}
