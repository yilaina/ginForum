package logic

import (
	"go_code/ginStudy/gin_b2/bluebell/dao/mysql"
	"go_code/ginStudy/gin_b2/bluebell/models"
)

func GetCommunityList() ([]*models.Community, error) {
	// 查数据库 查找到所有的community 并返回
	return mysql.GetCommunityList()
}
func GetCommunityDetail(id int64) (*models.CommunityDetail, error) {
	return mysql.GetCommunityDetailByID(id)
}
