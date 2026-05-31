package logic

import (
	"context"
	"core-rpc/internal/utils"
	"strings"

	"core-rpc/internal/model/entity"
	"core-rpc/internal/svc"
	"core-rpc/pb/core"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetArticleCommentsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetArticleCommentsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetArticleCommentsLogic {
	return &GetArticleCommentsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetArticleCommentsLogic) GetArticleComments(in *core.GetArticleCommentsReq) (*core.GetArticleCommentsResp, error) {
	page := utils.NormalizePage(in.Page)
	size := utils.NormalizeSize(in.Size, 10)
	off, limit := utils.OffsetLimit(page, size)

	q := l.svcCtx.Db.Model(&entity.Comment{}).
		Where("article_id = ? AND parent_id = 0", in.ArticleId)
	if strings.EqualFold(in.OrderBy, "hot") {
		q = q.Order("like_count DESC, created_at DESC")
	} else {
		q = q.Order("created_at DESC")
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, err
	}

	var comments []entity.Comment
	if err := q.Offset(off).Limit(limit).Find(&comments).Error; err != nil {
		return nil, err
	}

	userMap, err := fetchUserMap(l.svcCtx.Db, collectUserIDsFromComments(comments))
	if err != nil {
		return nil, err
	}

	previewMap, err := l.loadPreviewReplies(comments, userMap)
	if err != nil {
		return nil, err
	}

	return &core.GetArticleCommentsResp{
		Comments: commentsToProtoList(comments, userMap, previewMap),
		Page:     int32(page),
		Size:     int32(size),
		Total:    int32(total),
	}, nil
}

func (l *GetArticleCommentsLogic) loadPreviewReplies(topComments []entity.Comment, userMap map[uint64]entity.User) (map[uint64][]*core.CommentInfo, error) {
	if len(topComments) == 0 {
		return nil, nil
	}
	rootIDs := make([]uint64, len(topComments))
	for i, c := range topComments {
		rootIDs[i] = c.ID
	}

	var replies []entity.Comment
	if err := l.svcCtx.Db.Where("root_id IN ? AND parent_id > 0", rootIDs).
		Order("created_at ASC").Find(&replies).Error; err != nil {
		return nil, err
	}

	extraIDs := collectUserIDsFromComments(replies)
	for _, id := range extraIDs {
		if _, ok := userMap[id]; !ok {
			u, err := fetchUserMap(l.svcCtx.Db, []uint64{id})
			if err != nil {
				return nil, err
			}
			for k, v := range u {
				userMap[k] = v
			}
		}
	}

	const previewLimit = 3
	grouped := make(map[uint64][]entity.Comment)
	for _, r := range replies {
		list := grouped[r.RootID]
		if len(list) < previewLimit {
			grouped[r.RootID] = append(list, r)
		}
	}

	out := make(map[uint64][]*core.CommentInfo, len(grouped))
	for rootID, list := range grouped {
		items := make([]*core.CommentInfo, 0, len(list))
		for i := range list {
			name, avatar := utils.UserDisplay(userMap, list[i].UserID)
			items = append(items, utils.CommentToProto(&list[i], name, avatar, nil))
		}
		out[rootID] = items
	}
	return out, nil
}
