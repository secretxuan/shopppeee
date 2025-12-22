import { Layout as AntLayout, Menu, Badge, Avatar, Dropdown, Space } from 'antd'
import {
  HomeOutlined,
  ShoppingOutlined,
  ShoppingCartOutlined,
  UserOutlined,
  LoginOutlined,
  LogoutOutlined,
} from '@ant-design/icons'
import { Outlet, useNavigate, useLocation } from 'react-router-dom'
import { useAuthStore } from '@/store/useAuthStore'
import { useCartStore } from '@/store/useCartStore'
import type { MenuProps } from 'antd'
import './Layout.css'

const { Header, Content, Footer } = AntLayout

const Layout = () => {
  const navigate = useNavigate()
  const location = useLocation()
  const { isAuthenticated, user, logout } = useAuthStore()
  const { getTotalItems } = useCartStore()

  const menuItems: MenuProps['items'] = [
    {
      key: '/',
      icon: <HomeOutlined />,
      label: '首页',
      onClick: () => navigate('/'),
    },
    {
      key: '/products',
      icon: <ShoppingOutlined />,
      label: '商品',
      onClick: () => navigate('/products'),
    },
    {
      key: '/cart',
      icon: (
        <Badge count={getTotalItems()} offset={[10, 0]}>
          <ShoppingCartOutlined />
        </Badge>
      ),
      label: '购物车',
      onClick: () => navigate('/cart'),
    },
  ]

  const userMenuItems: MenuProps['items'] = isAuthenticated
    ? [
        {
          key: 'profile',
          icon: <UserOutlined />,
          label: '个人中心',
          onClick: () => navigate('/profile'),
        },
        {
          type: 'divider',
        },
        {
          key: 'logout',
          icon: <LogoutOutlined />,
          label: '退出登录',
          onClick: () => {
            logout()
            navigate('/')
          },
        },
      ]
    : [
        {
          key: 'login',
          icon: <LoginOutlined />,
          label: '登录',
          onClick: () => navigate('/login'),
        },
      ]

  return (
    <AntLayout className="layout">
      <Header className="header">
        <div className="header-content">
          <div className="logo" onClick={() => navigate('/')}>
            <ShoppingOutlined style={{ fontSize: 28, marginRight: 8 }} />
            <span>Shoppee</span>
          </div>
          <Menu
            theme="dark"
            mode="horizontal"
            selectedKeys={[location.pathname]}
            items={menuItems}
            className="menu"
          />
          <div className="user-section">
            <Dropdown menu={{ items: userMenuItems }} placement="bottomRight">
              <Space className="user-dropdown">
                <Avatar icon={<UserOutlined />} src={user?.avatar} />
                <span className="username">
                  {isAuthenticated ? user?.username : '未登录'}
                </span>
              </Space>
            </Dropdown>
          </div>
        </div>
      </Header>
      <Content className="content">
        <div className="content-wrapper">
          <Outlet />
        </div>
      </Content>
      <Footer className="footer">
        <div className="footer-content">
          <p>Shoppee 电商平台 © 2025 - 优质商品，优质服务</p>
          <p>基于 Go + React + PostgreSQL 构建</p>
        </div>
      </Footer>
    </AntLayout>
  )
}

export default Layout
