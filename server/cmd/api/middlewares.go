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

func WithCors(h http.Handler) http.Handler {
	middleware := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   connectcors.AllowedMethods(),
		AllowedHeaders:   connectcors.AllowedHeaders(),
		ExposedHeaders:   connectcors.ExposedHeaders(),
		AllowCredentials: true,
	})
	return middleware.Handler(h)
}

func (app *app) withTokenRefresh(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jwtCookie, jwtErr := r.Cookie("jwt")
		refreshCookie, refreshErr := r.Cookie("refresh")

		if jwtErr != nil || refreshErr != nil {
			next.ServeHTTP(w, r)
			return
		}

		if _, err := app.jwt.VerifyToken(jwtCookie.Value); err == nil {
			next.ServeHTTP(w, r)
			return
		}

		refreshClaims, err := app.jwt.VerifyToken(refreshCookie.Value)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		newToken, err := app.jwt.CreateToken(refreshClaims.ID, refreshClaims.Name, refreshClaims.Email)
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

func (app *app) authenticate(_ context.Context, req *http.Request) (any, error) {
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

	token, refreshToken, err := app.jwt.ExtractTokens(req)

	if err != nil {
		return nil, authn.Errorf("invalid authorization")
	}

	jwtClaims, refreshClaims, err := app.extractClaims(token, refreshToken)

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
