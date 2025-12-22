package service

import (
	"errors"

	"github.com/shoppee/ecommerce/internal/database"
	"github.com/shoppee/ecommerce/internal/models"
	"github.com/shoppee/ecommerce/pkg/logger"
	"go.uber.org/zap"
)

// AddressService 地址服务
type AddressService struct{}

// NewAddressService 创建地址服务实例
func NewAddressService() *AddressService {
	return &AddressService{}
}

// GetUserAddresses 获取用户所有地址
func (s *AddressService) GetUserAddresses(userID uint) ([]models.Address, error) {
	var addresses []models.Address
	if err := database.DB.Where("user_id = ?", userID).Order("is_default DESC, created_at DESC").Find(&addresses).Error; err != nil {
		return nil, err
	}
	return addresses, nil
}

// GetAddress 获取地址详情
func (s *AddressService) GetAddress(id, userID uint) (*models.Address, error) {
	var address models.Address
	if err := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&address).Error; err != nil {
		return nil, err
	}
	return &address, nil
}

// CreateAddress 创建地址
func (s *AddressService) CreateAddress(userID uint, req interface{}) (*models.Address, error) {
	reqMap, ok := req.(map[string]interface{})
	if !ok {
		// 尝试结构体类型
		type AddressReq struct {
			Name      string `json:"name"`
			Phone     string `json:"phone"`
			Province  string `json:"province"`
			City      string `json:"city"`
			District  string `json:"district"`
			Detail    string `json:"detail"`
			IsDefault bool   `json:"is_default"`
		}
		reqStruct, ok := req.(*AddressReq)
		if !ok {
			return nil, errors.New("无效的请求数据")
		}
		
		address := &models.Address{
			UserID:    userID,
			Name:      reqStruct.Name,
			Phone:     reqStruct.Phone,
			Province:  reqStruct.Province,
			City:      reqStruct.City,
			District:  reqStruct.District,
			Detail:    reqStruct.Detail,
			IsDefault: reqStruct.IsDefault,
		}
		
		if err := database.DB.Create(address).Error; err != nil {
			return nil, err
		}
		
		logger.Info("创建地址成功", zap.Uint("user_id", userID), zap.Uint("address_id", address.ID))
		return address, nil
	}
	
	address := &models.Address{
		UserID:    userID,
		Name:      reqMap["name"].(string),
		Phone:     reqMap["phone"].(string),
		Province:  reqMap["province"].(string),
		City:      reqMap["city"].(string),
		District:  reqMap["district"].(string),
		Detail:    reqMap["detail"].(string),
	}
	
	if isDefault, ok := reqMap["is_default"].(bool); ok {
		address.IsDefault = isDefault
	}
	
	if err := database.DB.Create(address).Error; err != nil {
		return nil, err
	}
	
	logger.Info("创建地址成功", zap.Uint("user_id", userID), zap.Uint("address_id", address.ID))
	return address, nil
}

// UpdateAddress 更新地址
func (s *AddressService) UpdateAddress(id, userID uint, updates map[string]interface{}) error {
	var address models.Address
	if err := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&address).Error; err != nil {
		return err
	}
	
	if err := database.DB.Model(&address).Updates(updates).Error; err != nil {
		return err
	}
	
	logger.Info("更新地址成功", zap.Uint("address_id", id))
	return nil
}

// DeleteAddress 删除地址
func (s *AddressService) DeleteAddress(id, userID uint) error {
	result := database.DB.Where("id = ? AND user_id = ?", id, userID).Delete(&models.Address{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("地址不存在")
	}
	
	logger.Info("删除地址成功", zap.Uint("address_id", id))
	return nil
}

// SetDefaultAddress 设置默认地址
func (s *AddressService) SetDefaultAddress(id, userID uint) error {
	var address models.Address
	if err := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&address).Error; err != nil {
		return err
	}
	
	// 先将该用户的所有地址设为非默认
	database.DB.Model(&models.Address{}).Where("user_id = ?", userID).Update("is_default", false)
	
	// 设置当前地址为默认
	if err := database.DB.Model(&address).Update("is_default", true).Error; err != nil {
		return err
	}
	
	logger.Info("设置默认地址成功", zap.Uint("address_id", id))
	return nil
}

// GetDefaultAddress 获取默认地址
func (s *AddressService) GetDefaultAddress(userID uint) (*models.Address, error) {
	var address models.Address
	if err := database.DB.Where("user_id = ? AND is_default = ?", userID, true).First(&address).Error; err != nil {
		return nil, err
	}
	return &address, nil
}
