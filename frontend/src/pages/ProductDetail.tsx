import { useEffect, useState } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import {
  Row,
  Col,
  Image,
  Typography,
  Button,
  InputNumber,
  Tag,
  Spin,
  message,
  Breadcrumb,
} from 'antd'
import { ShoppingCartOutlined, HomeOutlined } from '@ant-design/icons'
import { productAPI, Product } from '@/api/product'
import { useCartStore } from '@/store/useCartStore'
import './ProductDetail.css'

const { Title, Paragraph } = Typography

const ProductDetail = () => {
  const { id } = useParams<{ id: string }>()
  const navigate = useNavigate()
  const { addItem } = useCartStore()
  const [product, setProduct] = useState<Product | null>(null)
  const [loading, setLoading] = useState(true)
  const [quantity, setQuantity] = useState(1)

  useEffect(() => {
    if (id) {
      loadProduct(parseInt(id))
    }
  }, [id])

  const loadProduct = async (productId: number) => {
    try {
      setLoading(true)
      const data = await productAPI.getProductById(productId)
      setProduct(data)
    } catch (error) {
      console.error('加载商品详情失败:', error)
      message.error('商品不存在')
      navigate('/products')
    } finally {
      setLoading(false)
    }
  }

  const handleAddToCart = () => {
    if (product) {
      addItem(product, quantity)
      message.success('已添加到购物车')
    }
  }

  if (loading) {
    return (
      <div className="loading-container">
        <Spin size="large" tip="加载中..." />
      </div>
    )
  }

  if (!product) {
    return null
  }

  return (
    <div className="product-detail-page">
      <Breadcrumb className="breadcrumb">
        <Breadcrumb.Item href="/">
          <HomeOutlined /> 首页
        </Breadcrumb.Item>
        <Breadcrumb.Item href="/products">商品列表</Breadcrumb.Item>
        <Breadcrumb.Item>{product.name}</Breadcrumb.Item>
      </Breadcrumb>

      <div className="product-detail-content">
        <Row gutter={[32, 32]}>
          <Col xs={24} md={12}>
            <div className="product-images">
              <Image
                src={product.images?.[0] || 'https://via.placeholder.com/600x600?text=Product'}
                alt={product.name}
                className="main-image"
              />
            </div>
          </Col>

          <Col xs={24} md={12}>
            <div className="product-info">
              <Title level={2}>{product.name}</Title>

              <div className="product-meta">
                <Tag color="blue">SKU: {product.sku}</Tag>
                <Tag color={product.stock > 0 ? 'success' : 'error'}>
                  库存: {product.stock}
                </Tag>
                <Tag color={product.status === 'active' ? 'success' : 'default'}>
                  {product.status === 'active' ? '在售' : '下架'}
                </Tag>
              </div>

              <div className="product-price">
                <span className="price-label">价格：</span>
                <span className="price-value">¥{product.price.toFixed(2)}</span>
              </div>

              <div className="product-description">
                <Title level={4}>商品描述</Title>
                <Paragraph>{product.description}</Paragraph>
              </div>

              <div className="product-actions">
                <div className="quantity-selector">
                  <span className="quantity-label">数量：</span>
                  <InputNumber
                    min={1}
                    max={product.stock}
                    value={quantity}
                    onChange={(value) => setQuantity(value || 1)}
                    disabled={product.stock === 0}
                  />
                  <span className="stock-info">
                    （库存 {product.stock} 件）
                  </span>
                </div>

                <div className="action-buttons">
                  <Button
                    type="primary"
                    size="large"
                    icon={<ShoppingCartOutlined />}
                    onClick={handleAddToCart}
                    disabled={product.stock === 0}
                    block
                  >
                    加入购物车
                  </Button>
                </div>
              </div>
            </div>
          </Col>
        </Row>
      </div>
    </div>
  )
}

export default ProductDetail
