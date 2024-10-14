package filterV2

type Condition struct {
	Or *bool `json:"or"`
	// Group       *GroupCondition `json:"Group"`
	Equal    *Equal    `json:"equal"`
	NotEqual *NotEqual `json:"not_equal"`

	HigherThan    *HigherThan    `json:"higher_than"`
	HigherOrEqual *HigherOrEqual `json:"higher_or_equal"`

	LowerThan    *LowerThan    `json:"lower_than"`
	LowerOrEqual *LowerOrEqual `json:"lower_or_equal"`

	Contains    *Contains    `json:"contains"`
	NotContains *NotContains `json:"not_contains"`

	BeginsWith    *BeginsWith    `json:"begins_with"`
	NotBeginsWith *NotBeginsWith `json:"not_begins_with"`

	EndsWith    *EndsWith    `json:"ends_with"`
	NotEndsWith *NotEndsWith `json:"not_ends_with"`

	// TODO: add JSON functions: in array, has key etc...

	In    *In    `json:"in"`
	NotIn *NotIn `json:"not_in"`

	IsNull    *IsNull    `json:"is_null"`
	IsNotNull *IsNotNull `json:"is_not_null"`

	IsTrue  *IsTrue  `json:"is_true"`
	IsFalse *IsFalse `json:"is_false"`

	IsEmpty    *IsEmpty    `json:"is_empty"`
	IsNotEmpty *IsNotEmpty `json:"is_not_empty"`

	Between    *Between    `json:"between"`
	NotBetween *NotBetween `json:"not_between"`

	BetweenUnixTimestamp    *BetweenUnixTimestamp    `json:"between_unix_timestamp"`
	NotBetweenUnixTimestamp *NotBetweenUnixTimestamp `json:"not_between_unix_timestamp"`
}
