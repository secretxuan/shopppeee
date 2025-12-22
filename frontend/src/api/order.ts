import request from './axios';

// 创建订单
export const createOrder = (data: {
  address_id: number;
  cart_item_ids: number[];
  payment_method: string;
  remark?: string;
}) => request.post('/orders', data);

// 获取订单列表
export const getOrderList = (params: {
  page?: number;
  page_size?: number;
  status?: string;
}) => request.get('/orders', { params });

// 获取订单详情
export const getOrderDetail = (id: number) => request.get(`/orders/${id}`);

// 取消订单
export const cancelOrder = (id: number) => request.post(`/orders/${id}/cancel`);

// 确认收货
export const confirmReceipt = (id: number) => request.post(`/orders/${id}/confirm`);

// 管理员获取订单列表
export const adminGetOrderList = (params: {
  page?: number;
  page_size?: number;
  status?: string;
}) => request.get('/orders/admin', { params });

// 管理员更新订单状态
export const adminUpdateOrderStatus = (id: number, status: string) =>
  request.patch(`/orders/admin/${id}/status`, { status });
