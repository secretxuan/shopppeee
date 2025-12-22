import { useState, useEffect } from 'react';
import { List, Card, Tag, Button, Space, message, Empty, Tabs } from 'antd';
import { useNavigate } from 'react-router-dom';
import { getOrderList, cancelOrder, confirmReceipt } from '../api/order';
import './Orders.css';

const { TabPane } = Tabs;

const Orders = () => {
  const [orders, setOrders] = useState([]);
  const [loading, setLoading] = useState(false);
  const [activeTab, setActiveTab] = useState('all');
  const navigate = useNavigate();

  const statusMap: Record<string, { text: string; color: string }> = {
    pending: { text: '待支付', color: 'orange' },
    paid: { text: '已支付', color: 'blue' },
    shipped: { text: '已发货', color: 'cyan' },
    completed: { text: '已完成', color: 'green' },
    cancelled: { text: '已取消', color: 'red' },
  };

  const loadOrders = async (status?: string) => {
    setLoading(true);
    try {
      const params = status === 'all' ? {} : { status };
      const res = await getOrderList(params);
      setOrders(res.data.data || []);
    } catch (error) {
      message.error('加载订单失败');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadOrders(activeTab);
  }, [activeTab]);

  const handleCancel = async (id: number) => {
    try {
      await cancelOrder(id);
      message.success('订单已取消');
      loadOrders(activeTab);
    } catch (error) {
      message.error('取消订单失败');
    }
  };

  const handleConfirm = async (id: number) => {
    try {
      await confirmReceipt(id);
      message.success('确认收货成功');
      loadOrders(activeTab);
    } catch (error) {
      message.error('确认收货失败');
    }
  };

  const renderActions = (order: any) => {
    const actions = [];

    if (order.status === 'pending') {
      actions.push(
        <Button key="pay" type="primary" size="small">
          去支付
        </Button>
      );
      actions.push(
        <Button key="cancel" size="small" onClick={() => handleCancel(order.id)}>
          取消订单
        </Button>
      );
    }

    if (order.status === 'shipped') {
      actions.push(
        <Button
          key="confirm"
          type="primary"
          size="small"
          onClick={() => handleConfirm(order.id)}
        >
          确认收货
        </Button>
      );
    }

    if (order.status === 'completed') {
      actions.push(
        <Button key="review" size="small">
          评价
        </Button>
      );
    }

    actions.push(
      <Button key="detail" size="small">
        查看详情
      </Button>
    );

    return actions;
  };

  return (
    <div className="orders-page">
      <h1>我的订单</h1>

      <Tabs activeKey={activeTab} onChange={setActiveTab}>
        <TabPane tab="全部订单" key="all" />
        <TabPane tab="待支付" key="pending" />
        <TabPane tab="已支付" key="paid" />
        <TabPane tab="已发货" key="shipped" />
        <TabPane tab="已完成" key="completed" />
      </Tabs>

      <List
        loading={loading}
        dataSource={orders}
        locale={{
          emptyText: <Empty description="暂无订单" />,
        }}
        renderItem={(order: any) => (
          <Card className="order-card" key={order.id}>
            <div className="order-header">
              <div className="order-info">
                <span className="order-no">订单号：{order.order_no}</span>
                <span className="order-time">
                  {new Date(order.created_at).toLocaleString()}
                </span>
              </div>
              <Tag color={statusMap[order.status]?.color}>
                {statusMap[order.status]?.text}
              </Tag>
            </div>

            <div className="order-items">
              {order.order_items?.map((item: any) => (
                <div key={item.id} className="order-item">
                  <img
                    src={item.product_image || 'https://via.placeholder.com/80'}
                    alt={item.product_name}
                  />
                  <div className="item-info">
                    <div className="item-name">{item.product_name}</div>
                    <div className="item-sku">SKU: {item.product_sku}</div>
                  </div>
                  <div className="item-price">¥{item.price.toFixed(2)}</div>
                  <div className="item-quantity">x{item.quantity}</div>
                  <div className="item-subtotal">¥{item.sub_total.toFixed(2)}</div>
                </div>
              ))}
            </div>

            <div className="order-footer">
              <div className="order-total">
                合计：<span className="total-amount">¥{order.total_amount.toFixed(2)}</span>
              </div>
              <Space className="order-actions">{renderActions(order)}</Space>
            </div>
          </Card>
        )}
      />
    </div>
  );
};

export default Orders;
