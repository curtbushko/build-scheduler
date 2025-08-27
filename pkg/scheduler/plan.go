package scheduler

type Plan struct {
	name  string
	Slots int
}

var HobbyPlan = Plan{name: "hobby", Slots: 2}
var ProPlan = Plan{name: "pro", Slots: 16}
var EnterprisePlan = Plan{name: "enterprise", Slots: 32}
