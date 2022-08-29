package forms

type errors map[string][]string

func (error errors) Add(field, message string) {
	error[field] = append(error[field], message)
}

func (error errors) Get(field string) string {
	errorMessage := error[field]
	if len(errorMessage) == 0 {
		return ""
	}
	return errorMessage[0]
}
