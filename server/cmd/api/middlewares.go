package main

import (
	"context"
	"net/http"

	"connectrpc.com/authn"
	connectcors "connectrpc.com/cors"
	"github.com/rs/cors"
	"nopresh.apetrovic.com/gen/proto/auth/v1/authv1connect"
	"nopresh.apetrovic.com/internal/utils/auth"
)

type Middlewares struct {
	jwt *auth.JWT
}

func (m *Middlewares) WithCors(h http.Handler) http.Handler {
	middleware := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   connectcors.AllowedMethods(),
		AllowedHeaders:   connectcors.AllowedHeaders(),
		ExposedHeaders:   connectcors.ExposedHeaders(),
		AllowCredentials: true,
	})
	return middleware.Handler(h)
}

func (m *Middlewares) WithTokenRefresh(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jwtCookie, jwtErr := r.Cookie("jwt")
		refreshCookie, refreshErr := r.Cookie("refresh")

		if jwtErr != nil || refreshErr != nil {
			next.ServeHTTP(w, r)
			return
		}

		if _, err := m.jwt.VerifyToken(jwtCookie.Value); err == nil {
			next.ServeHTTP(w, r)
			return
		}

		refreshClaims, err := m.jwt.VerifyToken(refreshCookie.Value)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		newToken, err := m.jwt.CreateToken(refreshClaims.ID, refreshClaims.Name, refreshClaims.Email)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "jwt",
			Value:    newToken,
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
			Path:     "/",
		})

		// Update the request so the authn middleware sees the new token
		r.Header.Set("Cookie", "jwt="+newToken+"; refresh="+refreshCookie.Value)

		next.ServeHTTP(w, r)
	})
}

func (m *Middlewares) Authenticate(_ context.Context, req *http.Request) (any, error) {
	allowList := map[string]struct{}{
		authv1connect.AuthServiceRegisterProcedure:                       {},
		authv1connect.AuthServiceLoginProcedure:                          {},
		"/grpc.reflection.v1.ServerReflection/ServerReflectionInfo":      {},
		"/grpc.reflection.v1alpha.ServerReflection/ServerReflectionInfo": {},
	}

	procedure, _ := authn.InferProcedure(req.URL)

	if _, ok := allowList[procedure]; ok {
		return nil, nil
	}

	token, refreshToken, err := m.jwt.ExtractTokens(req)

	if err != nil {
		return nil, authn.Errorf("invalid authorization")
	}

	jwtClaims, refreshClaims, err := m.jwt.ExtractClaims(token, refreshToken)

	if err != nil {
		return nil, authn.Errorf("invalid authorization")
	}

	return &auth.AuthInfo{
		RefreshClaims: refreshClaims,
		JwtClaims:     jwtClaims,
		JWTToken:      token,
		RefreshToken:  refreshToken,
	}, nil
}
