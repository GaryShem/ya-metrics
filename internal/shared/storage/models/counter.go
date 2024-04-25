package models

type Counter struct {
	Type  string `json:"type"`
	Name  string `json:"name"`
	Value int64  `json:"value"`
}

func NewCounter(name string, value int64) *Counter {
	return &Counter{
		Type:  "counter",
		Name:  name,
		Value: value,
	}
}

func CopyCounter(original Counter) *Counter {
	return &Counter{
		Type:  "counter",
		Name:  original.Name,
		Value: original.Value,
	}
}

func (cm *Counter) Update(value int64) {
	cm.Value += value
}

func (cm *Counter) Reset() {
	cm.Value = 0
}
