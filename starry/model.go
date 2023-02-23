package starry

type Info struct {
	Title string
	Time  float64
	Round int
	Teams []Team
}

type unityData struct {
	Type string
	Data interface{}
}

type Team struct {
	Id    int
	Name  string
	Rank  int
	Img   string // 队伍logo 的url
	Score int
}

type msg struct {
	Type string
	Data interface{}
}

type status struct {
	Id     uint
	Status string
}

type attack struct {
	From uint
	To   uint
}

type rank struct {
	Teams []Team
}

type round struct {
	Round uint
}

type restTime struct {
	Time float64
}
type clearStatus struct {
	Id uint
}
