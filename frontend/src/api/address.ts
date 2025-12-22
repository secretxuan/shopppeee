import request from './axios';

// 获取地址列表
export const getAddressList = () => request.get('/addresses');

// 获取默认地址
export const getDefaultAddress = () => request.get('/addresses/default');

// 获取地址详情
export const getAddress = (id: number) => request.get(`/addresses/${id}`);

// 创建地址
export const createAddress = (data: {
  name: string;
  phone: string;
  province: string;
  city: string;
  district: string;
  detail: string;
  is_default?: boolean;
}) => request.post('/addresses', data);

// 更新地址
export const updateAddress = (id: number, data: any) =>
  request.put(`/addresses/${id}`, data);

// 删除地址
export const deleteAddress = (id: number) => request.delete(`/addresses/${id}`);

// 设置默认地址
export const setDefaultAddress = (id: number) =>
  request.patch(`/addresses/${id}/default`);
