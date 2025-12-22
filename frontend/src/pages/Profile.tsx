import { Card, Descriptions, Avatar, Typography, Tag } from 'antd'
import { UserOutlined } from '@ant-design/icons'
import { useAuthStore } from '@/store/useAuthStore'
import dayjs from 'dayjs'
import './Profile.css'

const { Title } = Typography

const Profile = () => {
  const { user } = useAuthStore()

  if (!user) {
    return null
  }

  return (
    <div className="profile-page">
      <Card>
        <div className="profile-header">
          <Avatar size={80} icon={<UserOutlined />} src={user.avatar} />
          <div className="profile-info">
            <Title level={3}>{user.username}</Title>
            <Tag color={user.role === 'admin' ? 'red' : 'blue'}>
              {user.role === 'admin' ? '管理员' : '普通用户'}
            </Tag>
          </div>
        </div>

        <Descriptions title="个人信息" column={1} bordered>
          <Descriptions.Item label="用户ID">{user.id}</Descriptions.Item>
          <Descriptions.Item label="用户名">{user.username}</Descriptions.Item>
          <Descriptions.Item label="邮箱">{user.email}</Descriptions.Item>
          <Descriptions.Item label="手机号">
            {user.phone || '未设置'}
          </Descriptions.Item>
          <Descriptions.Item label="角色">
            {user.role === 'admin' ? '管理员' : '普通用户'}
          </Descriptions.Item>
        </Descriptions>
      </Card>
    </div>
  )
}

export default Profile
