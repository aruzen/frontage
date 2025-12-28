package data

type ActionSummary struct {
	ActionTag string // 効果名
	Type      string
	Data      map[string]interface{}
}

type ActionResult struct {
	ActionTag  string
	Context    map[string]interface{}
	State      map[string]interface{}
	SummaryIdx int
}
