package controller

import (
	// _itemShopRepository "github.com/TxChar/isekai-shop-api/pkg/itemShop/repository"
	_itemShopService "github.com/TxChar/isekai-shop-api/pkg/itemShop/service"
)

type itemShopControllerImpl struct{
	itemShopService _itemShopService.ItemShopService
}

func NewItemShopServiceImpl(
	itemShopService _itemShopService.ItemShopService,
) ItemShopController {
	return &itemShopControllerImpl{itemShopService}
}