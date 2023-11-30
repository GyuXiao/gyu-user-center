import {extend} from 'umi-request';
import {message} from "antd";
import {history} from "umi";
import {stringify} from "querystring";

/**
 * 通过设置全局请求和响应的拦截器，统一的进行错误处理；
 * 很好的尝试，但暂时不用
 */

/**
 * 配置request请求时的默认参数
 */
const request = extend({
  credentials: 'include', // 默认请求是否带上cookie
  // requestType: 'form',
});

/**
 * 所以请求拦截器
 */
request.interceptors.request.use((url, options): any => {
  console.log(`do request url = ${url}`)
  return {
    url,
    options: {
      ...options,
      headers: {},
    },
  };
});

/**
 * 所有响应拦截器
 */
request.interceptors.response.use(async (response, options): Promise<any> => {
  const res = await response.clone().json();
  if (res.code === 20000000) {
    return res.data
  }
  if (res.code === 20010009) {
    message.error('请先登录');
    history.replace({
      pathname: `/user/login`,
      search: stringify({
        redirect: location.pathname,
      }),
    });
  }

  message.error(res.msg)
  return res;
});

export default request;
