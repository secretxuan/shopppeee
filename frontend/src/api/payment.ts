import request from './axios';

// 创建支付
export const createPayment = (data: {
  order_id: number;
  payment_method: string;
}) => request.post('/payments', data);

// 获取支付详情
export const getPayment = (id: number) => request.get(`/payments/${id}`);

// 模拟支付回调
export const mockPaymentCallback = (data: {
  payment_no: string;
  third_party_no: string;
  status: string;
}) => request.post('/payment-callback', data);
