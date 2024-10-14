package filter

type Condition struct {
	Or *bool `json:"Or"`
	// Group       *GroupCondition `json:"Group"`
	Eq    *Eq    `json:"Eq"`
	NotEq *NotEq `json:"NotEqual"`

	Ht *Ht `json:"HigherThan"`
	He *He `json:"HigherOrEqual"`

	Lt *Lt `json:"LowerThan"`
	Le *Le `json:"LowerOrEqual"`

	Contains    *Contains    `json:"Contains"`
	NotContains *NotContains `json:"NotContains"`

	BeginsWith    *BeginsWith    `json:"BeginsWith"`
	NotBeginsWith *NotBeginsWith `json:"NotBeginsWith"`

	EndsWith    *EndsWith    `json:"EndsWith"`
	NotEndsWith *NotEndsWith `json:"NotEndsWith"`

	// TODO: add JSON functions: in array, has key etc...

	In    *In    `json:"In"`
	NotIn *NotIn `json:"NotIn"`

	IsNull  *IsNull  `json:"IsNull"`
	NotNull *NotNull `json:"IsNotNull"`

	IsTrue  *IsTrue  `json:"IsTrue"`
	IsFalse *IsFalse `json:"IsFalse"`

	Empty    *Empty    `json:"IsEmpty"`
	NotEmpty *NotEmpty `json:"IsNotEmpty"`

	Between    *Between    `json:"Between"`
	NotBetween *NotBetween `json:"NotBetween"`
}
