package filterV2

type HigherOrEqual struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type HigherThan struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type Contains struct {
	Name            string `json:"name"`
	Value           string `json:"value"`
	CaseInsensitive *bool  `json:"case_insensitive"`
}

type NotContains struct {
	Name            string `json:"name"`
	Value           string `json:"value"`
	CaseInsensitive *bool  `json:"case_insensitive"`
}

type BeginsWith struct {
	Name            string `json:"name"`
	Value           string `json:"value"`
	CaseInsensitive *bool  `json:"case_insensitive"`
}

type NotBeginsWith struct {
	Name            string `json:"name"`
	Value           string `json:"value"`
	CaseInsensitive *bool  `json:"case_insensitive"`
}

type EndsWith struct {
	Name            string `json:"name"`
	Value           string `json:"value"`
	CaseInsensitive *bool  `json:"case_insensitive"`
}

type NotEndsWith struct {
	Name            string `json:"name"`
	Value           string `json:"value"`
	CaseInsensitive *bool  `json:"case_insensitive"`
}

type Between struct {
	Name  string  `json:"name"`
	Start *string `json:"start"`
	End   *string `json:"end"`
}

type BetweenUnixTimestamp struct {
	Name  string `json:"name"`
	Start *int64 `json:"start"`
	End   *int64 `json:"end"`
}

type IsEmpty struct {
	Name string `json:"name"`
}

type Equal struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type In struct {
	Name  string    `json:"name"`
	Value []*string `json:"value"`
}

type LowerOrEqual struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type LowerThan struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type NotBetween struct {
	Name  string  `json:"name"`
	Start *string `json:"start"`
	End   *string `json:"end"`
}

type NotBetweenUnixTimestamp struct {
	Name  string `json:"name"`
	Start *int64 `json:"start"`
	End   *int64 `json:"end"`
}

type IsNotEmpty struct {
	Name string `json:"name"`
}

type NotEqual struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type NotIn struct {
	Name  string    `json:"name"`
	Value []*string `json:"value"`
}

type IsNotNull struct {
	Name string `json:"name"`
}

type IsNull struct {
	Name string `json:"name"`
}

type IsTrue struct {
	Name string `json:"name"`
}

type IsFalse struct {
	Name string `json:"name"`
}
