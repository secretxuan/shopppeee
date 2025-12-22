import { useNavigate } from 'react-router-dom'
import {
  Table,
  Button,
  InputNumber,
  Empty,
  Typography,
  Space,
  Popconfirm,
  Image,
  Card,
} from 'antd'
import { DeleteOutlined, ShoppingOutlined } from '@ant-design/icons'
import { useCartStore, CartItem } from '@/store/useCartStore'
import type { ColumnsType } from 'antd/es/table'
import './Cart.css'

const { Title, Text } = Typography

const Cart = () => {
  const navigate = useNavigate()
  const { items, updateQuantity, removeItem, clearCart, getTotalPrice } = useCartStore()

  const columns: ColumnsType<CartItem> = [
    {
      title: '商品',
      dataIndex: 'name',
      key: 'name',
      render: (_, record) => (
        <Space>
          <Image
            width={80}
            src={record.images?.[0] || 'https://via.placeholder.com/80x80?text=Product'}
            alt={record.name}
            preview={false}
          />
          <div>
            <div className="product-name">{record.name}</div>
            <Text type="secondary">SKU: {record.sku}</Text>
          </div>
        </Space>
      ),
    },
    {
      title: '单价',
      dataIndex: 'price',
      key: 'price',
      render: (price: number) => <span>¥{price.toFixed(2)}</span>,
    },
    {
      title: '数量',
      dataIndex: 'quantity',
      key: 'quantity',
      render: (quantity: number, record) => (
        <InputNumber
          min={1}
          max={record.stock}
          value={quantity}
          onChange={(value) => updateQuantity(record.id, value || 1)}
        />
      ),
    },
    {
      title: '小计',
      key: 'subtotal',
      render: (_, record) => (
        <span className="subtotal">
          ¥{(record.price * record.quantity).toFixed(2)}
        </span>
      ),
    },
    {
      title: '操作',
      key: 'action',
      render: (_, record) => (
        <Popconfirm
          title="确定要删除这个商品吗？"
          onConfirm={() => removeItem(record.id)}
          okText="确定"
          cancelText="取消"
        >
          <Button type="text" danger icon={<DeleteOutlined />}>
            删除
          </Button>
        </Popconfirm>
      ),
    },
  ]

  if (items.length === 0) {
    return (
      <div className="cart-page">
        <Card>
          <Empty
            description="购物车空空如也"
            image={Empty.PRESENTED_IMAGE_SIMPLE}
          >
            <Button
              type="primary"
              icon={<ShoppingOutlined />}
              onClick={() => navigate('/products')}
            >
              去逛逛
            </Button>
          </Empty>
        </Card>
      </div>
    )
  }

  return (
    <div className="cart-page">
      <Card>
        <Title level={3}>购物车</Title>

        <Table
          columns={columns}
          dataSource={items}
          rowKey="id"
          pagination={false}
          scroll={{ x: 800 }}
        />

        <div className="cart-footer">
          <div className="cart-actions">
            <Popconfirm
              title="确定要清空购物车吗？"
              onConfirm={clearCart}
              okText="确定"
              cancelText="取消"
            >
              <Button danger>清空购物车</Button>
            </Popconfirm>
            <Button onClick={() => navigate('/products')}>
              继续购物
            </Button>
          </div>

          <div className="cart-summary">
            <div className="summary-item">
              <Text>商品总数：</Text>
              <Text strong>{items.reduce((sum, item) => sum + item.quantity, 0)} 件</Text>
            </div>
            <div className="summary-item total">
              <Text>总计：</Text>
              <Text strong className="total-price">
                ¥{getTotalPrice().toFixed(2)}
              </Text>
            </div>
            <Button type="primary" size="large" block>
              去结算
            </Button>
          </div>
        </div>
      </Card>
    </div>
  )
}

export default Cart
