package logic

import (
	"go.uber.org/zap"
	"shop-backend/dao/mysql"
	"shop-backend/models/vo"
	"sync"
	"time"
)

// 存放协程函数信息的通道
var errorChannel = make(chan error, 2)

// GetCategories 获取sku分类信息
func GetCategories(revChan chan<- []*vo.CategoryVO, cid1, cid2 int64) {
	categories, err := mysql.SelectSpuCategoryByCID(cid1, cid2)
	if err != nil {
		errorChannel <- err
	}
	// 封装categoriesVO
	categoriesVO := make([]*vo.CategoryVO, 0)
	for _, cate := range categories {
		categoriesVO = append(categoriesVO, &vo.CategoryVO{
			ID:   cate.ID,
			Name: cate.Name,
		})
	}
	// 发送到类型为[]*vo.CategoryVO的通道中
	revChan <- categoriesVO
	// 此通道为一发送一接收的对应关系，所以尽量在发送方关闭通道(即使关闭了通道，也能获取通道里数据)
	defer close(revChan)
}

// GetSkuList 获取spu下属的sku集合
func GetSkuList(revChan chan<- []*vo.SkuVO, spuID int64) {
	// 获取skuList
	skuList, err := mysql.SelectSkuListBySpuID(spuID)
	if err != nil {
		errorChannel <- err
	}
	// 封装skuListVO
	skuListVO := make([]*vo.SkuVO, 0)
	for _, sku := range skuList {
		// 获取sku的所有商品图片
		skuPicList, _ := mysql.SelectSkuPicBySkuID(sku.ID)
		// 封装skuPicListVO
		skuPicListVO := make([]*vo.SkuPicVO, 0)
		for _, pic := range skuPicList {
			skuPicListVO = append(skuPicListVO, &vo.SkuPicVO{
				ID:        pic.ID,
				SkuID:     pic.SkuID,
				PicUrl:    pic.PicUrl,
				IsDefault: pic.IsDefault,
			})
		}
		skuListVO = append(skuListVO, &vo.SkuVO{
			ID:                      sku.ID,
			SpuID:                   sku.SpuID,
			Title:                   sku.Title,
			Price:                   sku.Price,
			Unit:                    sku.Unit,
			Stock:                   sku.Stock,
			ProductSkuSpecification: sku.ProductSkuSpecification,
			SkuPicList:              skuPicListVO,
		})
	}
	// 发送到类型为[]*vo.SkuVO的通道中
	revChan <- skuListVO
	// 此通道为一发送一接收的对应关系，所以尽量在发送方关闭通道(即使关闭了通道，也能获取通道里数据)
	defer close(revChan)
}

// GetProductDetailWithConcurrent 多协程获取商品详情
func GetProductDetailWithConcurrent(skuID int64) (*vo.ProductDetailVO, error) {
	// start := time.Now()
	detail := new(vo.ProductDetailVO)
	// 获取spu
	spu, err := mysql.SelectSpuBySkuID(skuID)
	if err != nil {
		return nil, err
	}
	// 封装spuVO
	spuVO := &vo.SpuVO{
		ID:                   spu.ID,
		Sale:                 spu.Sale,
		SubTitle:             spu.SubTitle,
		ProductSpecification: spu.ProductSpecification,
	}
	detail.Spu = spuVO

	// 创建两个类型不同的通道，分别用于获取category和skuList
	revChan1 := make(chan []*vo.CategoryVO, 1)
	revChan2 := make(chan []*vo.SkuVO, 1)

	// 开启两个协程，并将对应的通道传递到协程函数中，由协程函数负责往通道里装入数据
	go GetCategories(revChan1, spu.CID1, spu.CID2)
	go GetSkuList(revChan2, spu.ID)

	// case1：如果协程函数未执行完，从通道中接收数据的操作会在通道底层的接收队列阻塞等待。所以程序会阻塞在下面两行
	// case2：如果协程函数已经执行完，那么直接从通道中获取数据，完成detail的赋值操作
	detail.Categories = <-revChan1
	detail.SkuList = <-revChan2
	// 处理异常
	if len(errorChannel) > 0 {
		err := <-errorChannel
		zap.L().Error("多协程获取商品详情失败", zap.Error(err))
		return nil, err
	}
	// cost := time.Since(start)
	// zap.L().Info(" GetProductDetailWithConcurrent 商品详情接口耗时:", zap.Duration("cost", cost))
	return detail, nil
}

// GetCategories2 获取sku分类信息
func GetCategories2(channel chan *vo.ProductDetailVO, cid1, cid2 int64) {
	categories, _ := mysql.SelectSpuCategoryByCID(cid1, cid2)
	// 封装categoriesVO
	categoriesVO := make([]*vo.CategoryVO, 0)
	for _, cate := range categories {
		categoriesVO = append(categoriesVO, &vo.CategoryVO{
			ID:   cate.ID,
			Name: cate.Name,
		})
	}

	// 通道中只有一个detail对象，如果两个协程函数都在方法第一行获取通道中的对象，那么肯定有一个协程函数会进入通道底层的接收队列中进行阻塞等待
	// 与其直接在里面等待，不如先查询出要添加到detail的数据，在成功获取到通道中的detail时，直接赋值
	// 从通道里获取detail
	detail := <-channel
	// 更新detail
	detail.Categories = categoriesVO
	// 重新放入到通道
	channel <- detail
	// 任务完成，waitGroup计数-1
	wg.Done()
}

// GetSkuList2 获取spu下属的sku集合
func GetSkuList2(channel chan *vo.ProductDetailVO, spuID int64) {
	// 获取skuList
	skuList, _ := mysql.SelectSkuListBySpuID(spuID)
	// 封装skuListVO
	skuListVO := make([]*vo.SkuVO, 0)
	for _, sku := range skuList {
		// 获取sku的所有商品图片
		skuPicList, _ := mysql.SelectSkuPicBySkuID(sku.ID)
		// 封装skuPicListVO
		skuPicListVO := make([]*vo.SkuPicVO, 0)
		for _, pic := range skuPicList {
			skuPicListVO = append(skuPicListVO, &vo.SkuPicVO{
				ID:        pic.ID,
				SkuID:     pic.SkuID,
				PicUrl:    pic.PicUrl,
				IsDefault: pic.IsDefault,
			})
		}
		skuListVO = append(skuListVO, &vo.SkuVO{
			ID:                      sku.ID,
			SpuID:                   sku.SpuID,
			Title:                   sku.Title,
			Price:                   sku.Price,
			Unit:                    sku.Unit,
			Stock:                   sku.Stock,
			ProductSkuSpecification: sku.ProductSkuSpecification,
			SkuPicList:              skuPicListVO,
		})
	}

	// 通道中只有一个detail对象，如果两个协程函数都在方法第一行获取通道中的对象，那么肯定有一个协程函数会进入通道底层的接收队列中进行阻塞等待
	// 与其直接在里面等待，不如先查询出要添加到detail的数据，在成功获取到通道中的detail时，直接赋值
	// 从通道里获取detail
	detail := <-channel
	// 更新detail
	detail.SkuList = skuListVO
	channel <- detail
	// 任务完成，waitGroup计数-1
	wg.Done()
}

var wg sync.WaitGroup

// GetProductDetailWithConcurrent2 多协程获取商品详情
func GetProductDetailWithConcurrent2(skuID int64) (*vo.ProductDetailVO, error) {
	start := time.Now()
	detail := new(vo.ProductDetailVO)
	// 获取spu
	spu, _ := mysql.SelectSpuBySkuID(skuID)
	// 封装spuVO
	spuVO := &vo.SpuVO{
		ID:                   spu.ID,
		Sale:                 spu.Sale,
		SubTitle:             spu.SubTitle,
		ProductSpecification: spu.ProductSpecification,
	}
	detail.Spu = spuVO

	// 创建一个类型为*vo.ProductDetailVO类型的通道
	// 这里必须指定通道的缓存区，如果不执行的话，下一行添加数据到通道，就会放入到通道底层的发送队列中，等待被消费。
	// 代码会在下一行阻塞，执行不了协程函数消费数据，从而造成程序永久阻塞
	channel := make(chan *vo.ProductDetailVO, 1)
	// 将detail添加到通道
	channel <- detail
	// 此时通道的对应关系为一发送对多接收，所以在发送方关闭通道比较合适
	defer close(channel)

	// 下面的两个协程函数，如果未在函数返回前执行，那么最后返回的detail对象，skuList和categories为空，因为协程函数还没有执行完。
	// 所以使用WaitGroup。
	wg.Add(2)
	go GetCategories2(channel, spu.CID1, spu.CID2)
	go GetSkuList2(channel, spu.ID)

	cost := time.Since(start)
	zap.L().Info("GetProductDetailWithConcurrent2 商品详情接口耗时:", zap.Duration("cost", cost))
	wg.Wait()
	// 因为通道中的数据类型为指针类型，所以最后返回的detail包含spu、categories、skuList
	return <-channel, nil
}

// GetProductDetail 获取商品详情
func GetProductDetail(skuID int64) (*vo.ProductDetailVO, error) {
	start := time.Now()
	detail := new(vo.ProductDetailVO)
	// 获取spu
	spu, err := mysql.SelectSpuBySkuID(skuID)
	if err != nil {
		return nil, err
	}
	// 封装spuVO
	spuVO := &vo.SpuVO{
		ID:                   spu.ID,
		Sale:                 spu.Sale,
		SubTitle:             spu.SubTitle,
		ProductSpecification: spu.ProductSpecification,
	}
	detail.Spu = spuVO

	// 获取categories
	categories, err := mysql.SelectSpuCategoryByCID(spu.CID1, spu.CID2)
	if err != nil {
		return nil, err
	}
	// 封装categoriesVO
	categoriesVO := make([]*vo.CategoryVO, 0)
	for _, cate := range categories {
		categoriesVO = append(categoriesVO, &vo.CategoryVO{
			ID:   cate.ID,
			Name: cate.Name,
		})
	}
	detail.Categories = categoriesVO

	// 获取skuList
	skuList, err := mysql.SelectSkuListBySpuID(spu.ID)
	if err != nil {
		return nil, err
	}
	// 封装skuListVO
	skuListVO := make([]*vo.SkuVO, 0)
	for _, sku := range skuList {
		// 获取sku的所有商品图片
		skuPicList, err := mysql.SelectSkuPicBySkuID(sku.ID)
		if err != nil {
			return nil, err
		}
		// 封装skuPicListVO
		skuPicListVO := make([]*vo.SkuPicVO, 0)
		for _, pic := range skuPicList {
			skuPicListVO = append(skuPicListVO, &vo.SkuPicVO{
				ID:        pic.ID,
				SkuID:     pic.SkuID,
				PicUrl:    pic.PicUrl,
				IsDefault: pic.IsDefault,
			})
		}
		skuListVO = append(skuListVO, &vo.SkuVO{
			ID:                      sku.ID,
			SpuID:                   sku.SpuID,
			Title:                   sku.Title,
			Price:                   sku.Price,
			Unit:                    sku.Unit,
			Stock:                   sku.Stock,
			ProductSkuSpecification: sku.ProductSkuSpecification,
			SkuPicList:              skuPicListVO,
		})
	}
	detail.SkuList = skuListVO
	cost := time.Since(start)
	zap.L().Info("商品详情接口耗时:", zap.Duration("cost", cost))
	return detail, nil
}
