package logic

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"core-rpc/internal/model/category"
	"core-rpc/internal/svc"
	"core-rpc/pb/core"

	"github.com/go-sql-driver/mysql"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type CreateCategoryLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateCategoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateCategoryLogic {
	return &CreateCategoryLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// Category 部分
func (l *CreateCategoryLogic) CreateCategory(in *core.CreateCategoryReq) (*core.CreateCategoryResp, error) {
	if in == nil || in.Name == "" {
		return nil, errors.New("category name is required")
	}

	// prepare slug: if provided in request name contains spaces, create slug
	slug := strings.ToLower(strings.ReplaceAll(in.Name, " ", "-"))

	var err error

	// check by name or slug
	_, err = l.svcCtx.CategoryModel.FindOneByName(l.ctx, in.Name)
	if err == nil {
		return nil, errors.New("category already exists")
	}
	if err != nil && !errors.Is(err, sqlx.ErrNotFound) {
		return nil, err
	}
	_, err = l.svcCtx.CategoryModel.FindOneBySlug(l.ctx, slug)
	if err == nil {
		return nil, errors.New("category slug already exists")
	}
	if err != nil && !errors.Is(err, sqlx.ErrNotFound) {
		return nil, err
	}

	newCat := &category.Category{
		Name:        in.Name,
		Slug:        slug,
		Description: sql.NullString{String: "", Valid: false},
	}

	_, err = l.svcCtx.CategoryModel.Insert(l.ctx, newCat)
	if err != nil {
		// handle duplicate key (concurrent insert) - detect MySQL duplicate error 1062
		var me *mysql.MySQLError
		if errors.As(err, &me) && me.Number == 1062 {
			return nil, errors.New("category already exists")
		}
		return nil, err
	}

	return &core.CreateCategoryResp{}, nil
}
