import request from './axios';

// 获取购物车
export const getCart = () => request.get('/cart');

// 添加到购物车
export const addToCart = (data: { product_id: number; quantity: number }) =>
  request.post('/cart/items', data);

// 更新购物车项
export const updateCartItem = (id: number, quantity: number) =>
  request.put(`/cart/items/${id}`, { quantity });

// 删除购物车项
export const removeFromCart = (id: number) => request.delete(`/cart/items/${id}`);

// 清空购物车
export const clearCart = () => request.delete('/cart/clear');

// 选中/取消选中购物车项
export const selectCartItem = (id: number, selected: boolean) =>
  request.patch(`/cart/items/${id}/select`, { selected });
