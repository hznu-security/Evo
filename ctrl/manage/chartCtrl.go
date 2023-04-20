package manage

import (
	"Evo/db"
	"Evo/model"
	"Evo/util"
	"github.com/gin-gonic/gin"
	"sort"
)

type ChartLine struct {
	Name string    `json:"name"`
	Data []float64 `json:"data"`
}

type ChartData struct {
	XData  []int       `json:"xData"`
	YData  []ChartLine `json:"yData"`
	Legend []string    `json:"legend"`
}

func GetChart(c *gin.Context) {
	round := 15
	var chartData ChartData
	xData := make([]int, round)
	for i := 0; i < int(round); i++ {
		xData[i] = i
	}
	chartData.XData = xData

	// 纵坐标
	yData := make([]ChartLine, 0)
	teams := make([]model.Team, 0)
	legend := make([]string, 0)
	db.DB.Model(&model.Team{}).Select([]string{"name"}).Find(&teams)
	for _, team := range teams {
		line := ChartLine{}
		line.Name = team.Name
		charts := make([]model.Chart, 0)
		db.DB.Where("team_name = ?", team.Name).Find(&charts)
		// 将该队伍的chart按升序排序
		sort.Slice(charts, func(i, j int) bool {
			return charts[i].Round < charts[j].Round
		})
		data := make([]float64, 0)
		for _, chart := range charts {
			data = append(data, chart.Score)
		}
		line.Data = data
		yData = append(yData, line)
		legend = append(legend, team.Name)
	}
	chartData.YData = yData
	chartData.Legend = legend
	util.Success(c, "success", gin.H{
		"data": chartData,
	})
}
