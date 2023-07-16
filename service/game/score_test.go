package game

import (
	"Evo/config"
	"Evo/db"
	"testing"
)

func BenchmarkCalcScore(b *testing.B) {
	config.InitConfig()
	db.InitDB()
	for n := 0; n < b.N; n++ {
		CalcScore(1)
	}
}

//func TestInsertAttack(t *testing.T) {
//	config.InitConfig()
//	db.InitDB()
//	for i := 1; i < 5; i++ {
//		attack := &model.Attack{
//			Attacker:    4,
//			TeamID:      1,
//			Round:       1,
//			ChallengeId: uint(i),
//			GameBoxId:   uint(i),
//		}
//		db.DB.Save(attack)
//	}
//}
