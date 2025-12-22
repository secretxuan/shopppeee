import { useEffect, useState } from 'react'
import { Row, Col, Input, Select, Spin, Empty, Pagination, Card } from 'antd'
import { SearchOutlined } from '@ant-design/icons'
import { productAPI, Product } from '@/api/product'
import ProductCard from '@/components/ProductCard'
import './Products.css'

const { Search } = Input
const { Option } = Select

const Products = () => {
  const [products, setProducts] = useState<Product[]>([])
  const [loading, setLoading] = useState(true)
  const [total, setTotal] = useState(0)
  const [page, setPage] = useState(1)
  const [pageSize, setPageSize] = useState(12)
  const [keyword, setKeyword] = useState('')
  const [sortBy, setSortBy] = useState<string>('')

  useEffect(() => {
    loadProducts()
  }, [page, pageSize, sortBy])

  const loadProducts = async () => {
    try {
      setLoading(true)
      const params: any = {
        page,
        page_size: pageSize,
      }
      if (sortBy) params.sort = sortBy
      
      const data = keyword
        ? await productAPI.searchProducts(keyword, page, pageSize)
        : await productAPI.getProductList(params)
      
      setProducts(data.list || [])
      setTotal(data.total || 0)
    } catch (error) {
      console.error('加载商品失败:', error)
    } finally {
      setLoading(false)
    }
  }

  const handleSearch = (value: string) => {
    setKeyword(value)
    setPage(1)
    loadProducts()
  }

  const handleSortChange = (value: string) => {
    setSortBy(value)
    setPage(1)
  }

  const handlePageChange = (newPage: number, newPageSize: number) => {
    setPage(newPage)
    setPageSize(newPageSize)
  }

  return (
    <div className="products-page">
      <Card className="filter-card">
        <Row gutter={16} align="middle">
          <Col xs={24} md={12}>
            <Search
              placeholder="搜索商品..."
              allowClear
              enterButton={<SearchOutlined />}
              size="large"
              onSearch={handleSearch}
            />
          </Col>
          <Col xs={24} md={12}>
            <Select
              placeholder="排序方式"
              size="large"
              style={{ width: '100%' }}
              onChange={handleSortChange}
              value={sortBy}
            >
              <Option value="">默认排序</Option>
              <Option value="price_asc">价格从低到高</Option>
              <Option value="price_desc">价格从高到低</Option>
              <Option value="created_desc">最新上架</Option>
            </Select>
          </Col>
        </Row>
      </Card>

      <div className="products-content">
        {loading ? (
          <div className="loading-container">
            <Spin size="large" tip="加载中..." />
          </div>
        ) : products.length > 0 ? (
          <>
            <div className="products-count">
              共找到 <strong>{total}</strong> 件商品
            </div>
            <Row gutter={[16, 16]}>
              {products.map((product) => (
                <Col xs={24} sm={12} md={8} lg={6} key={product.id}>
                  <ProductCard product={product} />
                </Col>
              ))}
            </Row>
            <div className="pagination-container">
              <Pagination
                current={page}
                pageSize={pageSize}
                total={total}
                onChange={handlePageChange}
                showSizeChanger
                showQuickJumper
                showTotal={(total) => `共 ${total} 件商品`}
                pageSizeOptions={['12', '24', '36', '48']}
              />
            </div>
          </>
        ) : (
          <Empty description="暂无商品" />
        )}
      </div>
    </div>
  )
}

export default Products
