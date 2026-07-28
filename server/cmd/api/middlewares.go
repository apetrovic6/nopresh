package main

import (
	"context"
	"log/slog"
	"net/http"
	"strings"

	"connectrpc.com/authn"
	connectcors "connectrpc.com/cors"
	"github.com/rs/cors"
	"nopresh.apetrovic.com/gen/proto/auth/v1/authv1connect"
	"nopresh.apetrovic.com/internal/utils/auth"
)

type Middlewares struct {
	jwt     *auth.JWT
	logger  *slog.Logger
	domains []string
}

func NewMiddleware(jwt *auth.JWT, logger *slog.Logger, domains string) Middlewares {
	parsedDomains := strings.Split(domains, ",")

	return Middlewares{
		jwt:     jwt,
		logger:  logger,
		domains: parsedDomains,
	}
}

func (m *Middlewares) WithCors(h http.Handler) http.Handler {
	middleware := cors.New(cors.Options{
		AllowedOrigins:   m.domains,
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

		// No refresh cookie or empty value — nothing we can do.
		if refreshErr != nil || refreshCookie.Value == "" {
			next.ServeHTTP(w, r)
			return
		}

		// JWT present and valid — no refresh needed.
		if jwtErr == nil && jwtCookie.Value != "" {
			if _, err := m.jwt.VerifyToken(jwtCookie.Value); err == nil {
				next.ServeHTTP(w, r)
				return
			}
		}

		// JWT is missing or expired — try to issue a new one from the refresh token.
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

		m.logger.Info("WithTokenRefresh: issued new jwt from refresh token")

		http.SetCookie(w, &http.Cookie{
			Name:     "jwt",
			Value:    newToken,
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
			Path:     "/",
			MaxAge:   900,
		})

		// Apply the fresh token to the in-flight request too, so the downstream
		// Authenticate handler sees it on this same request instead of 401'ing
		// and forcing the client to retry.
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
		m.logger.Error("invalid authorization. Couldn't extract token",
			"error", err.Error(),
		)
		return nil, authn.Errorf("invalid authorization. Couldn't extract token")
	}

	jwtClaims, refreshClaims, err := m.jwt.ExtractClaims(token, refreshToken)

	if err != nil {
		m.logger.Error("invalid authorization. Couldn't extract claims",
			"error", err.Error(),
		)
		return nil, authn.Errorf("invalid authorization. Couldn't extract claims")
	}

	return &auth.AuthInfo{
		RefreshClaims: refreshClaims,
		JwtClaims:     jwtClaims,
		JWTToken:      token,
		RefreshToken:  refreshToken,
	}, nil
}
