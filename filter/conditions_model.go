package filter

type He struct {
	Name  string `json:"Name"`
	Value string `json:"Value"`
}

type Ht struct {
	Name  string `json:"Name"`
	Value string `json:"Value"`
}

type Contains struct {
	Name            string `json:"Name"`
	Value           string `json:"Value"`
	CaseInsensitive *bool  `json:"CaseInsensitive"`
}

type NotContains struct {
	Name            string `json:"Name"`
	Value           string `json:"Value"`
	CaseInsensitive *bool  `json:"CaseInsensitive"`
}

type BeginsWith struct {
	Name            string `json:"Name"`
	Value           string `json:"Value"`
	CaseInsensitive *bool  `json:"CaseInsensitive"`
}

type NotBeginsWith struct {
	Name            string `json:"Name"`
	Value           string `json:"Value"`
	CaseInsensitive *bool  `json:"CaseInsensitive"`
}

type EndsWith struct {
	Name            string `json:"Name"`
	Value           string `json:"Value"`
	CaseInsensitive *bool  `json:"CaseInsensitive"`
}

type NotEndsWith struct {
	Name            string `json:"Name"`
	Value           string `json:"Value"`
	CaseInsensitive *bool  `json:"CaseInsensitive"`
}

type Between struct {
	Name   string  `json:"Name"`
	Value1 *string `json:"Start"`
	Value2 *string `json:"End"`
}

type Empty struct {
	Name string `json:"Name"`
}

type Eq struct {
	Name  string `json:"Name"`
	Value string `json:"Value"`
}

type In struct {
	Name  string    `json:"Name"`
	Value []*string `json:"Value"`
}

type Le struct {
	Name  string `json:"Name"`
	Value string `json:"Value"`
}

type Lt struct {
	Name  string `json:"Name"`
	Value string `json:"Value"`
}

type NotBetween struct {
	Name   string  `json:"Name"`
	Value1 *string `json:"Start"`
	Value2 *string `json:"End"`
}

type NotEmpty struct {
	Name string `json:"Name"`
}

type NotEq struct {
	Name  string `json:"Name"`
	Value string `json:"Value"`
}

type NotIn struct {
	Name  string    `json:"Name"`
	Value []*string `json:"Value"`
}

type NotNull struct {
	Name string `json:"Name"`
}

type IsNull struct {
	Name string `json:"Name"`
}

type IsTrue struct {
	Name string `json:"Name"`
}

type IsFalse struct {
	Name string `json:"Name"`
}
