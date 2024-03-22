package model

import (
	"fmt"
	"gorm.io/gorm"
	"time"
)

func GetVotes() []Vote {
	ret := make([]Vote, 0)
	err := Conn.Table("vote").Find(&ret).Error
	if err != nil {
		fmt.Printf("err:%s", err.Error())
	}
	return ret
}

func GetVote(id int64) VoteWithOpt {
	var ret Vote
	err := Conn.Table("vote").Where("id = ?", id).First(&ret).Error
	if err != nil {
		fmt.Printf("err:%s", err.Error())
	}
	opt := make([]VoteOpt, 0)
	err = Conn.Table("vote_opt").Where("vote_id = ?", id).Find(&opt).Error
	if err != nil {
		fmt.Printf("err:%s", err.Error())
	}
	return VoteWithOpt{
		Vote: ret,
		Opt:  opt,
	}
}

func PostVote(userId, voteId int64, optIDs []int64) bool {
	var ret Vote
	err := Conn.Table("vote").Where("id = ?", voteId).First(&ret).Error
	if err != nil {
		fmt.Printf("err:%s", err.Error())
		return false
	}
	for _, value := range optIDs {
		err = Conn.Table("vote_opt").Where("id = ?", value).Update("count", gorm.Expr("count + ?", 1)).Error
		if err != nil {
			fmt.Printf("err:%s", err.Error())
			return false
		}

		user := VoteOptUser{
			VoteId:      voteId,
			UserId:      userId,
			VoteOptId:   value,
			CreatedTime: time.Now(),
		}
		err = Conn.Table("vote_opt_user").Create(&user).Error
		if err != nil {
			fmt.Printf("err:%s", err.Error())
			return false
		}
	}

	return true
}
