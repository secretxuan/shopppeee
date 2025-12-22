import { useState, useEffect } from 'react';
import {
  Table,
  Button,
  Space,
  Modal,
  Form,
  Input,
  InputNumber,
  Select,
  message,
  Popconfirm,
  Tag,
} from 'antd';
import {
  PlusOutlined,
  EditOutlined,
  DeleteOutlined,
  EyeOutlined,
} from '@ant-design/icons';
import {
  getProductList,
  createProduct,
  updateProduct,
  deleteProduct,
  updateProductStatus,
} from '../../api/product';
import { getCategoryList } from '../../api/category';
import './ProductManage.css';

const ProductManage = () => {
  const [products, setProducts] = useState([]);
  const [categories, setCategories] = useState([]);
  const [loading, setLoading] = useState(false);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [pageSize] = useState(10);
  const [modalVisible, setModalVisible] = useState(false);
  const [editingProduct, setEditingProduct] = useState<any>(null);
  const [form] = Form.useForm();

  // 加载商品列表
  const loadProducts = async () => {
    setLoading(true);
    try {
      const res = await getProductList({ page, page_size: pageSize });
      setProducts(res.list || []);
      setTotal(res.total || 0);
    } catch (error) {
      message.error('加载失败');
    } finally {
      setLoading(false);
    }
  };

  // 加载分类列表
  const loadCategories = async () => {
    try {
      const res = await getCategoryList();
      setCategories(res.list || []);
    } catch (error) {
      console.error('加载分类失败');
    }
  };

  useEffect(() => {
    loadProducts();
    loadCategories();
  }, [page]);

  // 打开创建/编辑弹窗
  const showModal = (product?: any) => {
    setEditingProduct(product);
    if (product) {
      form.setFieldsValue(product);
    } else {
      form.resetFields();
    }
    setModalVisible(true);
  };

  // 处理保存
  const handleSave = async () => {
    try {
      const values = await form.validateFields();
      if (editingProduct) {
        await updateProduct(editingProduct.id, values);
        message.success('更新成功');
      } else {
        await createProduct(values);
        message.success('创建成功');
      }
      setModalVisible(false);
      loadProducts();
    } catch (error) {
      message.error('操作失败');
    }
  };

  // 删除商品
  const handleDelete = async (id: number) => {
    try {
      await deleteProduct(id);
      message.success('删除成功');
      loadProducts();
    } catch (error) {
      message.error('删除失败');
    }
  };

  // 更新状态
  const handleStatusChange = async (id: number, status: string) => {
    try {
      await updateProductStatus(id, status);
      message.success('状态更新成功');
      loadProducts();
    } catch (error) {
      message.error('状态更新失败');
    }
  };

  const columns = [
    {
      title: 'ID',
      dataIndex: 'id',
      width: 80,
    },
    {
      title: '商品名称',
      dataIndex: 'name',
      width: 200,
    },
    {
      title: '分类',
      dataIndex: ['category', 'name'],
      width: 120,
    },
    {
      title: 'SKU',
      dataIndex: 'sku',
      width: 120,
    },
    {
      title: '价格',
      dataIndex: 'price',
      width: 100,
      render: (price: number) => `¥${price.toFixed(2)}`,
    },
    {
      title: '库存',
      dataIndex: 'stock',
      width: 100,
    },
    {
      title: '销量',
      dataIndex: 'sale_count',
      width: 100,
    },
    {
      title: '状态',
      dataIndex: 'status',
      width: 120,
      render: (status: string, record: any) => (
        <Select
          value={status}
          style={{ width: 100 }}
          onChange={(value) => handleStatusChange(record.id, value)}
        >
          <Select.Option value="active">
            <Tag color="green">上架</Tag>
          </Select.Option>
          <Select.Option value="inactive">
            <Tag color="red">下架</Tag>
          </Select.Option>
        </Select>
      ),
    },
    {
      title: '操作',
      key: 'action',
      width: 200,
      render: (_: any, record: any) => (
        <Space>
          <Button
            type="link"
            size="small"
            icon={<EyeOutlined />}
            onClick={() => window.open(`/products/${record.id}`, '_blank')}
          >
            查看
          </Button>
          <Button
            type="link"
            size="small"
            icon={<EditOutlined />}
            onClick={() => showModal(record)}
          >
            编辑
          </Button>
          <Popconfirm
            title="确定删除该商品吗？"
            onConfirm={() => handleDelete(record.id)}
            okText="确定"
            cancelText="取消"
          >
            <Button type="link" size="small" danger icon={<DeleteOutlined />}>
              删除
            </Button>
          </Popconfirm>
        </Space>
      ),
    },
  ];

  return (
    <div className="product-manage">
      <div className="page-header">
        <h2>商品管理</h2>
        <Button type="primary" icon={<PlusOutlined />} onClick={() => showModal()}>
          添加商品
        </Button>
      </div>

      <Table
        columns={columns}
        dataSource={products}
        rowKey="id"
        loading={loading}
        pagination={{
          current: page,
          pageSize,
          total,
          onChange: (p) => setPage(p),
          showSizeChanger: false,
          showTotal: (total) => `共 ${total} 条`,
        }}
      />

      <Modal
        title={editingProduct ? '编辑商品' : '添加商品'}
        open={modalVisible}
        onOk={handleSave}
        onCancel={() => setModalVisible(false)}
        width={600}
        okText="保存"
        cancelText="取消"
      >
        <Form form={form} layout="vertical">
          <Form.Item
            name="name"
            label="商品名称"
            rules={[{ required: true, message: '请输入商品名称' }]}
          >
            <Input placeholder="请输入商品名称" />
          </Form.Item>

          <Form.Item name="description" label="商品描述">
            <Input.TextArea rows={3} placeholder="请输入商品描述" />
          </Form.Item>

          <Form.Item
            name="category_id"
            label="分类"
            rules={[{ required: true, message: '请选择分类' }]}
          >
            <Select placeholder="请选择分类">
              {categories.map((cat: any) => (
                <Select.Option key={cat.id} value={cat.id}>
                  {cat.name}
                </Select.Option>
              ))}
            </Select>
          </Form.Item>

          <Form.Item
            name="sku"
            label="SKU"
            rules={[{ required: true, message: '请输入SKU' }]}
          >
            <Input placeholder="请输入SKU" />
          </Form.Item>

          <Form.Item
            name="price"
            label="售价"
            rules={[{ required: true, message: '请输入售价' }]}
          >
            <InputNumber
              min={0}
              step={0.01}
              precision={2}
              style={{ width: '100%' }}
              placeholder="请输入售价"
            />
          </Form.Item>

          <Form.Item name="orig_price" label="原价">
            <InputNumber
              min={0}
              step={0.01}
              precision={2}
              style={{ width: '100%' }}
              placeholder="请输入原价"
            />
          </Form.Item>

          <Form.Item
            name="stock"
            label="库存"
            rules={[{ required: true, message: '请输入库存' }]}
          >
            <InputNumber min={0} style={{ width: '100%' }} placeholder="请输入库存" />
          </Form.Item>

          <Form.Item name="status" label="状态" initialValue="active">
            <Select>
              <Select.Option value="active">上架</Select.Option>
              <Select.Option value="inactive">下架</Select.Option>
            </Select>
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
};

export default ProductManage;
