package mysql

import (
	"database/sql"
	"go.uber.org/zap"
	"go_code/ginStudy/gin_b2/bluebell/models"
)

func GetCommunityList() (communityList []*models.Community, err error) {
	sqlStr := "select community_id, community_name from community"
	if db.Select(&communityList, sqlStr); err != nil {
		if err == sql.ErrNoRows {
			zap.L().Warn("no community list")
			err = nil
		}
	}
	return
}

// GetCommunityDetailByID 根据id查询社区详情
func GetCommunityDetailByID(id int64) (community *models.CommunityDetail, err error) {
	community = new(models.CommunityDetail)
	sqlStr := `select community_id, community_name , introduction, create_time 
				from community 
				where community_id = ?
				`
	if err := db.Get(community, sqlStr, id); err != nil {
		if err == sql.ErrNoRows {
			err = ErrorInvalidID
		}
	}
	return community, nil
}
