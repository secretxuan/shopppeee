import request from './axios';

// 创建评价
export const createReview = (data: {
  product_id: number;
  order_id?: number;
  rating: number;
  content?: string;
  images?: string;
}) => request.post('/reviews', data);

// 获取商品评价列表
export const getProductReviews = (
  productId: number,
  params?: { page?: number; page_size?: number }
) => request.get(`/reviews/products/${productId}`, { params });

// 获取我的评价
export const getMyReviews = (params?: { page?: number; page_size?: number }) =>
  request.get('/reviews/my', { params });

// 删除评价
export const deleteReview = (id: number) => request.delete(`/reviews/${id}`);

// 管理员回复评价
export const adminReplyReview = (id: number, reply: string) =>
  request.post(`/reviews/${id}/reply`, { reply });
