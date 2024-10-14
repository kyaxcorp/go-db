package filterV2

type OrderBy struct {
	FieldName string  `json:"field_name"`
	Direction *string `json:"direction"`
}
