package data

type EventSummary struct {
	ActionTag string // 効果名
	Data      map[string]interface{}
}

type ActionResult struct {
	ActionTag string
	Data      map[string]interface{}
}
