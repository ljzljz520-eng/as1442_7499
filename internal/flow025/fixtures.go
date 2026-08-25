package flow025

import "noticeword/internal/model"

func FixtureRecords() []model.Record {
	return []model.Record{
		{ID: "notice-001", Community: "星河社群", Title: "春季活动", Body: "请在周六上午九点集合", Description: "集合地点与报名说明", Status: model.StatusPublished, CharacterCount: 11, Revision: 1},
		{ID: "notice-002", Community: "星河社群", Title: "志愿者招募", Body: "欢迎报名社区志愿服务", Description: "报名截止日期和联系人", Status: model.StatusPublished, CharacterCount: 10, Revision: 1},
		{ID: "notice-003", Community: "星河社群", Title: "读书分享", Body: "本月主题是城市记忆", Description: "分享会流程和准备事项", Status: model.StatusPublished, CharacterCount: 10, Revision: 1},
	}
}

func FixtureQuery(page, pageSize int) model.SearchQuery {
	return model.SearchQuery{Community: "星河", Page: page, PageSize: pageSize}
}
