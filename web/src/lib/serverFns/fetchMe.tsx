import { me } from '#/gen/proto/auth/v1/auth-AuthService_connectquery';
import { callUnaryMethod } from '@connectrpc/connect-query';
import { createConnectTransport } from '@connectrpc/connect-web';
import { createServerFn } from '@tanstack/react-start';

export const checkMe = createServerFn().handler(async () => {
  const { getCookie } = await import('@tanstack/react-start/server');
  const jwtValue = getCookie('jwt') ?? '';
  const refreshValue = getCookie('refresh') ?? '';
  const t = createConnectTransport({
    baseUrl: 'http://localhost:5000/api',
    interceptors: [(next) => async (req) => {
      req.header.set('Cookie', `jwt=${jwtValue}; refresh=${refreshValue}`);
      return next(req);
    }],
  });

  return callUnaryMethod(t, me, {});
});
