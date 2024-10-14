package filter

type GroupCondition struct {
	Or         *bool             `json:"Or"`
	Conditions []*Condition      `json:"Conditions"`
	Groups     []*GroupCondition `json:"Groups"`
}
