package enums

// ENUM(query, fragment, form_post)
type ResponseMode string

func (r *ResponseMode) WithDefault() {
	if *r == "" {
		*r = ResponseModeQuery
	}
}
