package logic

import (
	"context"
	"core-rpc/internal/model/entity"
	"core-rpc/internal/svc"
	"core-rpc/internal/utils"
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

	total, comments, err := l.svcCtx.CommentRepo.ListTopComments(in.ArticleId, in.OrderBy, off, limit)
	if err != nil {
		return nil, err
	}

	// collect user ids from comments
	userIDs := make([]uint64, 0, len(comments))
	seen := make(map[uint64]struct{})
	for _, c := range comments {
		if _, ok := seen[c.UserID]; !ok {
			seen[c.UserID] = struct{}{}
			userIDs = append(userIDs, c.UserID)
		}
	}
	userMap, err := l.svcCtx.UserRepo.FindByIDs(userIDs)
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

	// load preview replies via repo
	const previewLimit = 3
	grouped, err := l.svcCtx.CommentRepo.LoadPreviewReplies(rootIDs, previewLimit)
	if err != nil {
		return nil, err
	}

	// ensure userMap contains all users referenced by replies
	extraIDs := make([]uint64, 0)
	for _, list := range grouped {
		for _, r := range list {
			if _, ok := userMap[r.UserID]; !ok {
				extraIDs = append(extraIDs, r.UserID)
			}
		}
	}
	if len(extraIDs) > 0 {
		extraMap, err := l.svcCtx.UserRepo.FindByIDs(extraIDs)
		if err != nil {
			return nil, err
		}
		for k, v := range extraMap {
			userMap[k] = v
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
