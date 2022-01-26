package api

const (
	ActionAddItem = iota + 1
	ActionDelItem
	ActionGetItem
	ActionGetAll
)

type ItemKey string
type ItemValue string

type Item struct {
	Key   ItemKey `json:"key"`
	Value ItemValue `json:"value"`
}

type Msg struct {
	Action int  `json:"action"`
	Item   Item `json:"item"`
}
