import request from './axios';

// 获取分类列表
export const getCategoryList = (params?: { parent_id?: number }) =>
  request.get('/categories', { params });

// 获取分类详情
export const getCategory = (id: number) => request.get(`/categories/${id}`);

// 创建分类（管理员）
export const createCategory = (data: {
  name: string;
  description?: string;
  icon?: string;
  parent_id?: number;
  sort?: number;
}) => request.post('/categories', data);

// 更新分类（管理员）
export const updateCategory = (id: number, data: any) =>
  request.put(`/categories/${id}`, data);

// 删除分类（管理员）
export const deleteCategory = (id: number) => request.delete(`/categories/${id}`);
