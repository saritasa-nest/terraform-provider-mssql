package model

type OptionsList map[string]NullString

func (o OptionsList) Parse(d map[string]interface{}) OptionsList {
	result := make(OptionsList)
	if d != nil {
		for opt := range d {
			value := d[opt]
			if value == nil {
				result[opt] = ""
			} else {
				result[opt] = NullString(value.(string))
			}
		}
	}
	return result
}
