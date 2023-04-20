package game

//// 插入新一轮的数据
//func insertChart() {
//	teams := make([]model.Team, 0)
//	db.DB.Select("score").Find(&teams)
//	charts := make([]model.Chart, 0)
//	for _, team := range teams {
//		charts = append(charts, model.Chart{
//			TeamName: team.Name,
//			Score:    team.Score,
//			Round:    config.ROUND_NOW - 1,
//		})
//	}
//	db.DB.Create(&charts)
//}
//
//func InitChart() {
//	reNewChart()
//}

//
//func reNewCache() {
//	// 横坐标
//	round := config.ROUND_NOW
//	var chartData ChartData
//	xData := make([]int, round)
//	for i := 0; i < int(round); i++ {
//		xData[i] = i
//	}
//	chartData.XData = xData
//
//	// 纵坐标
//	yData := make([]ChartLine, 0)
//	teams := make([]model.Team, 0)
//	db.DB.Model(&model.Team{}).Select([]string{"name"}).Find(&teams)
//	for _, team := range teams {
//		line := ChartLine{}
//		line.Name = team.Name
//		charts := make([]model.Chart, 0)
//		db.DB.Where("team_name").Find(&charts)
//		// 将该队伍的chart按升序排序
//		sort.Slice(charts, func(i, j int) bool {
//			return charts[i].Round < charts[j].Round
//		})
//		data := make([]float64, 0)
//		for _, chart := range charts {
//			data = append(data, chart.Score)
//		}
//		line.Data = data
//		yData = append(yData, line)
//	}
//	chartData.YData = yData
//	db.Set("chart", chartData)
//}
//
//func reNewChart() {
//	//insertChart()
//	//reNewCache()
//}
