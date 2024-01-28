import { GithubOutlined } from '@ant-design/icons';
import { DefaultFooter } from '@ant-design/pro-components';
const Footer: React.FC = () => {
  const defaultMessage = '永世传唱，不休独舞';
  const currentYear = new Date().getFullYear();
  return (
    <DefaultFooter
      copyright={`${currentYear} ${defaultMessage}`}
      links={[
        {
          key: 'Ant Design Pro',
          title: 'gyu_user_center',
          href: '',
          blankTarget: false,
        },
        {
          key: 'github',
          title: <GithubOutlined />,
          href: 'https://github.com/GyuXiao',
          blankTarget: true,
        },
        {
          key: 'gyustudio.site',
          title: '粤ICP备2024172144号',
          href: 'https://beian.miit.gov.cn',
          blankTarget: true,
        },
      ]}
    />
  );
};
export default Footer;
