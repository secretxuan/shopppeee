import api from './axios'

export interface Product {
  id: number
  name: string
  description: string
  price: number
  stock: number
  sku: string
  category_id: number
  images?: string[]
  status: string
  created_at: string
  updated_at: string
}

export interface ProductListParams {
  page?: number
  page_size?: number
  category_id?: number
  sort?: string
  keyword?: string
}

export interface ProductListResponse {
  list: Product[]
  total: number
  page: number
  page_size: number
}

export const productAPI = {
  // 获取商品列表
  getProductList: (params?: ProductListParams) => {
    return api.get<any, ProductListResponse>('/products', { params })
  },

  // 获取商品详情
  getProductById: (id: number) => {
    return api.get<any, Product>(`/products/${id}`)
  },

  // 搜索商品
  searchProducts: (keyword: string, page = 1, page_size = 20) => {
    return api.get<any, ProductListResponse>('/products/search', {
      params: { keyword, page, page_size },
    })
  },

  // 创建商品（管理员）
  createProduct: (data: Partial<Product>) => {
    return api.post('/products', data)
  },

  // 更新商品（管理员）
  updateProduct: (id: number, data: Partial<Product>) => {
    return api.put(`/products/${id}`, data)
  },

  // 删除商品（管理员）
  deleteProduct: (id: number) => {
    return api.delete(`/products/${id}`)
  },

  // 更新商品状态（管理员）
  updateProductStatus: (id: number, status: string) => {
    return api.patch(`/products/${id}/status`, { status })
  },

  // 批量更新库存（管理员）
  batchUpdateStock: (updates: Record<number, number>) => {
    return api.post('/products/batch-stock', updates)
  },
}

// 导出独立函数供旧代码使用
export const getProductList = productAPI.getProductList
export const getProductById = productAPI.getProductById
export const searchProducts = productAPI.searchProducts
export const createProduct = productAPI.createProduct
export const updateProduct = productAPI.updateProduct
export const deleteProduct = productAPI.deleteProduct
export const updateProductStatus = productAPI.updateProductStatus
export const batchUpdateStock = productAPI.batchUpdateStock
