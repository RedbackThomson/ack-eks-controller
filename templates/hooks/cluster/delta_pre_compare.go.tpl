    if a.ko.Spec.Tags == nil && b.ko.Spec.Tags != nil {
		a.ko.Spec.Tags = make(map[string]*string, 0)
	}