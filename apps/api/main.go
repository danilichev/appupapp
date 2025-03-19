package main

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/danielgtaylor/huma/v2/humacli"
	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"

	_ "github.com/danielgtaylor/huma/v2/formats/cbor"
)

type Options struct {
	Port int `help:"Port to listen on" short:"p" default:"8888"`
}

type GreetingOutput struct {
	Body struct {
		Message string `json:"message" example:"Hello, world!" doc:"Greeting message"`
	}
}

type ReviewInput struct {
	Body struct {
		Author  string `json:"author" maxLength:"10" doc:"Author of the review"`
		Rating  int    `json:"rating" minimum:"1" maximum:"5" doc:"Rating from 1 to 5"`
		Message string `json:"message,omitempty" maxLength:"100" doc:"Review message"`
	}
}

func addRoutes(api huma.API) {
	grp := huma.NewGroup(api, "/v1")
	grp1 := huma.NewGroup(grp, "")

	huma.Get(grp, "/greeting_/{name}", func(ctx context.Context, input *struct {
		Name string `path:"name" maxLength:"30" example:"world" doc:"Name to greet"`
	}) (*GreetingOutput, error) {
		resp := &GreetingOutput{}
		resp.Body.Message = fmt.Sprintf("Hello, %s!", input.Name)
		return resp, nil
	})

	huma.Register(grp1, huma.Operation{
		OperationID:   "post-review",
		Method:        http.MethodPost,
		Path:          "/reviews",
		Summary:       "Post a review",
		Tags:          []string{"Reviews"},
		DefaultStatus: http.StatusCreated,
	}, func(ctx context.Context, i *ReviewInput) (*struct{}, error) {
		return nil, nil
	})

	usersGrp := huma.NewGroup(grp, "/users")
	usersGrp.UseMiddleware(JWTMiddleware(api, []byte("secret")))
	usersGrp.UseSimpleModifier(func(op *huma.Operation) {
		op.Summary = "A summary for all operations in this group"
		op.Tags = []string{"my-tag"}
	})
}

type ctxKey string

const userCtxKey ctxKey = "user"

func JWTMiddleware(api huma.API, secret []byte) func(ctx huma.Context, next func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		authHeader := ctx.Header("Authorization")

		if authHeader == "" {
			huma.WriteErr(api, ctx, http.StatusUnauthorized, "Missing Authorization header")
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			huma.WriteErr(api, ctx, http.StatusUnauthorized, "Invalid Authorization header format (requires 'Bearer <token>')")
			return
		}
		tokenStr := parts[1]

		token, err := jwt.ParseWithClaims(tokenStr, &jwt.RegisteredClaims{}, func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
			return secret, nil
		})

		if err != nil {
			huma.WriteErr(api, ctx, http.StatusUnauthorized, fmt.Sprintf("Failed to parse token: %v", err))
			return
		}

		claims, ok := token.Claims.(*jwt.RegisteredClaims)
		if !ok || !token.Valid {
			huma.WriteErr(api, ctx, http.StatusUnauthorized, "Invalid or expired token")
			return
		}

		ctx = huma.WithValue(ctx, userCtxKey, claims)

		next(ctx)
	}
}

func main() {
	cli := humacli.New(func(hooks humacli.Hooks, options *Options) {
		router := chi.NewMux()
		api := humachi.New(router, huma.DefaultConfig("My API", "1.0.0"))

		addRoutes(api)

		router.Get("/docs", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/html")
			w.Write([]byte(`<!doctype html>
		<html>
		  <head>
			<title>API Reference</title>
			<meta charset="utf-8" />
			<meta
			  name="viewport"
			  content="width=device-width, initial-scale=1" />
		  </head>
		  <body>
			<script
			  id="api-reference"
			  data-url="/openapi.json"></script>
			<script src="https://cdn.jsdelivr.net/npm/@scalar/api-reference"></script>
		  </body>
		</html>`))
		})

		hooks.OnStart(func() {
			fmt.Printf("Starting server on port %d...\n", options.Port)
			http.ListenAndServe(fmt.Sprintf(":%d", options.Port), router)
		})
	})

	cli.Run()
}
