import { Card, Button, Tag, message } from 'antd'
import { ShoppingCartOutlined, EyeOutlined } from '@ant-design/icons'
import { useNavigate } from 'react-router-dom'
import { Product } from '@/api/product'
import { useCartStore } from '@/store/useCartStore'
import './ProductCard.css'

interface ProductCardProps {
  product: Product
}

const ProductCard: React.FC<ProductCardProps> = ({ product }) => {
  const navigate = useNavigate()
  const { addItem } = useCartStore()

  const handleAddToCart = (e: React.MouseEvent) => {
    e.stopPropagation()
    addItem(product, 1)
    message.success('已添加到购物车')
  }

  const handleViewDetails = () => {
    navigate(`/products/${product.id}`)
  }

  return (
    <Card
      hoverable
      className="product-card"
      cover={
        <div className="product-image-wrapper">
          <img
            alt={product.name}
            src={product.images?.[0] || 'https://via.placeholder.com/300x300?text=Product'}
            className="product-image"
          />
          {product.stock === 0 && (
            <div className="out-of-stock-overlay">
              <Tag color="red">缺货</Tag>
            </div>
          )}
        </div>
      }
      onClick={handleViewDetails}
    >
      <div className="product-info">
        <h3 className="product-name">{product.name}</h3>
        <p className="product-description">{product.description}</p>
        <div className="product-footer">
          <div className="product-price">
            <span className="price-symbol">¥</span>
            <span className="price-value">{product.price.toFixed(2)}</span>
          </div>
          <div className="product-actions">
            <Button
              type="text"
              icon={<EyeOutlined />}
              onClick={handleViewDetails}
            >
              查看
            </Button>
            <Button
              type="primary"
              icon={<ShoppingCartOutlined />}
              onClick={handleAddToCart}
              disabled={product.stock === 0}
            >
              加入购物车
            </Button>
          </div>
        </div>
        <div className="product-meta">
          <Tag color={product.stock > 0 ? 'success' : 'error'}>
            库存: {product.stock}
          </Tag>
          <Tag color="blue">SKU: {product.sku}</Tag>
        </div>
      </div>
    </Card>
  )
}

export default ProductCard
