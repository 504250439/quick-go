package service

import (
	"quick-go/app/entity"
	"quick-go/app/repo"
	"quick-go/global/consts"
	"quick-go/utils/quickErrors"
	"time"

	"github.com/gin-gonic/gin"
)

type SpuService struct {
	ctx     *gin.Context
	SpuRepo repo.ISpuRepo
}

func SpuServiceNew(ctx *gin.Context) *SpuService {
	svc := SpuService{ctx: ctx, SpuRepo: repo.NewMysqlSpuRepository()}
	return &svc
}

// getSpuInfo 获取商品的信息
func (s *SpuService) GetSpuInfo(req *entity.GetSpuInfoReq) (res *entity.GetSpuInfoRes, err error) {
	// 直接查数据然后返回
	spuInfo, err := s.SpuRepo.GetSpuDetail(s.ctx, req.AppID, req.SpuID)
	if err != nil {
		return nil, quickErrors.New(consts.CurdSelectErrorCode, consts.CurdSelectErrorMsg, err.Error())
	}

	if spuInfo.AppID == "" {
		return nil, quickErrors.New(consts.CurdSelectErrorCode, consts.CurdSelectErrorMsg, "spu not found")
	}

	res = &entity.GetSpuInfoRes{
		AppID:        spuInfo.AppID,
		SpuID:        spuInfo.SpuID,
		SpuType:      spuInfo.SpuType,
		ResourceID:   spuInfo.ResourceID,
		ResourceType: spuInfo.ResourceType,
		GoodsName:    spuInfo.GoodsName,
		GoodsImg:     spuInfo.GoodsImg,
	}
	return res, nil
}

// getSpuSaleTimeStamp 获取商品的销售时间戳
func getSpuSaleTimeStamp(saleAtString string) (saleAtStamp int64, err error) {
	saleAt, err := time.Parse("2006-01-02 15:04:05", saleAtString)
	if err != nil {
		return 0, err
	}
	return saleAt.Unix(), nil
}
