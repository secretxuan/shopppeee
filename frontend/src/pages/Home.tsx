import { useEffect, useState } from 'react'
import { Row, Col, Carousel, Typography, Spin, Empty } from 'antd'
import { ShoppingOutlined, SafetyOutlined, CustomerServiceOutlined, RocketOutlined } from '@ant-design/icons'
import { productAPI, Product } from '@/api/product'
import ProductCard from '@/components/ProductCard'
import './Home.css'

const { Title, Paragraph } = Typography

const Home = () => {
  const [products, setProducts] = useState<Product[]>([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    loadProducts()
  }, [])

  const loadProducts = async () => {
    try {
      setLoading(true)
      const data = await productAPI.getProductList({ page: 1, page_size: 8 })
      setProducts(data.list || [])
    } catch (error) {
      console.error('加载商品失败:', error)
    } finally {
      setLoading(false)
    }
  }

  const banners = [
    {
      title: '欢迎来到 Shoppee',
      subtitle: '优质商品，优质服务',
      color: '#336699',
    },
    {
      title: '新品上架',
      subtitle: '发现更多精选好物',
      color: '#52c41a',
    },
    {
      title: '限时优惠',
      subtitle: '全场商品特价促销',
      color: '#ff4d4f',
    },
  ]

  const features = [
    {
      icon: <ShoppingOutlined />,
      title: '精选商品',
      description: '严选优质商品，品质保证',
    },
    {
      icon: <SafetyOutlined />,
      title: '安全支付',
      description: '多重支付保障，购物更安心',
    },
    {
      icon: <RocketOutlined />,
      title: '快速配送',
      description: '专业物流配送，送货上门',
    },
    {
      icon: <CustomerServiceOutlined />,
      title: '贴心服务',
      description: '7x24小时客服，随时为您服务',
    },
  ]

  return (
    <div className="home-page">
      {/* 轮播图 */}
      <Carousel autoplay className="home-carousel">
        {banners.map((banner, index) => (
          <div key={index}>
            <div className="carousel-item" style={{ background: banner.color }}>
              <div className="carousel-content">
                <Title level={1} style={{ color: '#fff', marginBottom: 16 }}>
                  {banner.title}
                </Title>
                <Paragraph style={{ color: '#fff', fontSize: 18 }}>
                  {banner.subtitle}
                </Paragraph>
              </div>
            </div>
          </div>
        ))}
      </Carousel>

      {/* 特色服务 */}
      <div className="features-section">
        <Row gutter={[24, 24]}>
          {features.map((feature, index) => (
            <Col xs={24} sm={12} md={6} key={index}>
              <div className="feature-card">
                <div className="feature-icon">{feature.icon}</div>
                <Title level={4}>{feature.title}</Title>
                <Paragraph type="secondary">{feature.description}</Paragraph>
              </div>
            </Col>
          ))}
        </Row>
      </div>

      {/* 热门商品 */}
      <div className="products-section">
        <div className="section-header">
          <Title level={2}>热门商品</Title>
          <Paragraph type="secondary">精选优质商品，满足您的需求</Paragraph>
        </div>

        {loading ? (
          <div className="loading-container">
            <Spin size="large" tip="加载中..." />
          </div>
        ) : products.length > 0 ? (
          <Row gutter={[16, 16]}>
            {products.map((product) => (
              <Col xs={24} sm={12} md={8} lg={6} key={product.id}>
                <ProductCard product={product} />
              </Col>
            ))}
          </Row>
        ) : (
          <Empty description="暂无商品" />
        )}
      </div>
    </div>
  )
}

export default Home
