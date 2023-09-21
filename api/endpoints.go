package api

// endpoints.go contains the endpoint definitions, including per-method request
// and response structs. Endpoints are the binding between the service and
// transport.

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"

	"github.com/go-kit/kit/endpoint"
	"user/db"
	"user/users"
)

// Endpoints collects the endpoints that comprise the Service.
type Endpoints struct {
	LoginEndpoint       endpoint.Endpoint
	RegisterEndpoint    endpoint.Endpoint
	UserGetEndpoint     endpoint.Endpoint
	UserPostEndpoint    endpoint.Endpoint
	AddressGetEndpoint  endpoint.Endpoint
	AddressPostEndpoint endpoint.Endpoint
	CardGetEndpoint     endpoint.Endpoint
	CardPostEndpoint    endpoint.Endpoint
	DeleteEndpoint      endpoint.Endpoint
	HealthEndpoint      endpoint.Endpoint
}

// MakeEndpoints returns an Endpoints structure, where each endpoint is
// backed by the given service.
func MakeEndpoints(s Service) Endpoints {
	return Endpoints{
		LoginEndpoint:       MakeLoginEndpoint(s),
		RegisterEndpoint:    MakeRegisterEndpoint(s),
		HealthEndpoint:      MakeHealthEndpoint(s),
		UserGetEndpoint:     MakeUserGetEndpoint(s),
		UserPostEndpoint:    MakeUserPostEndpoint(s),
		AddressGetEndpoint:  MakeAddressGetEndpoint(s),
		AddressPostEndpoint: MakeAddressPostEndpoint(s),
		CardGetEndpoint:     MakeCardGetEndpoint(s),
		DeleteEndpoint:      MakeDeleteEndpoint(s),
		CardPostEndpoint:    MakeCardPostEndpoint(s),
	}
}

// MakeLoginEndpoint returns an endpoint via the given service.
func MakeLoginEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		tr := otel.Tracer("Login")
		_, span := tr.Start(ctx, "Login")
		span.SetAttributes(attribute.Key("service").String("user"))
		defer span.End()
		req := request.(loginRequest)
		u, err := s.Login(req.Username, req.Password)
		return userResponse{User: u}, err
	}
}

// MakeRegisterEndpoint returns an endpoint via the given service.
func MakeRegisterEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		tr := otel.Tracer("Register")
		_, span := tr.Start(ctx, "register")
		span.SetAttributes(attribute.Key("service").String("user"))
		defer span.End()
		req := request.(registerRequest)
		id, err := s.Register(req.Username, req.Password, req.Email, req.FirstName, req.LastName)
		return postResponse{ID: id}, err
	}
}

// MakeUserGetEndpoint returns an endpoint via the given service.
func MakeUserGetEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		tr := otel.Tracer("Get Users")
		ctx, span := tr.Start(ctx, "Get Users")
		span.SetAttributes(attribute.Key("service").String("user"))
		defer span.End()

		req := request.(GetRequest)

		ctx, userspan := tr.Start(ctx, "users from db")
		usrs, err := s.GetUsers(req.ID)
		userspan.End()
		if req.ID == "" {
			return EmbedStruct{usersResponse{Users: usrs}}, err
		}
		if len(usrs) == 0 {
			if req.Attr == "addresses" {
				return EmbedStruct{addressesResponse{Addresses: make([]users.Address, 0)}}, err
			}
			if req.Attr == "cards" {
				return EmbedStruct{cardsResponse{Cards: make([]users.Card, 0)}}, err
			}
			return users.User{}, err
		}
		user := usrs[0]
		ctx, attributespan := tr.Start(ctx, "attributes from db")
		db.GetUserAttributes(&user)
		attributespan.End()
		if req.Attr == "addresses" {
			return EmbedStruct{addressesResponse{Addresses: user.Addresses}}, err
		}
		if req.Attr == "cards" {
			return EmbedStruct{cardsResponse{Cards: user.Cards}}, err
		}
		return user, err
	}
}

// MakeUserPostEndpoint returns an endpoint via the given service.
func MakeUserPostEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		tr := otel.Tracer("Post User")
		ctx, span := tr.Start(ctx, "Post User")
		span.SetAttributes(attribute.Key("service").String("user"))
		defer span.End()
		req := request.(users.User)
		id, err := s.PostUser(req)
		return postResponse{ID: id}, err
	}
}

// MakeAddressGetEndpoint returns an endpoint via the given service.
func MakeAddressGetEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		tr := otel.Tracer("Get Users")
		ctx, span := tr.Start(ctx, "Get Users")
		span.SetAttributes(attribute.Key("service").String("user"))
		defer span.End()

		req := request.(GetRequest)

		ctx, addrspan := tr.Start(ctx, "address from db")

		adds, err := s.GetAddresses(req.ID)
		addrspan.End()
		if req.ID == "" {
			return EmbedStruct{addressesResponse{Addresses: adds}}, err
		}
		if len(adds) == 0 {
			return users.Address{}, err
		}
		return adds[0], err
	}
}

// MakeAddressPostEndpoint returns an endpoint via the given service.
func MakeAddressPostEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		tr := otel.Tracer("Post Address")
		ctx, span := tr.Start(ctx, "Post Address")
		span.SetAttributes(attribute.Key("service").String("user"))
		defer span.End()
		req := request.(addressPostRequest)
		id, err := s.PostAddress(req.Address, req.UserID)
		return postResponse{ID: id}, err
	}
}

// MakeUserGetEndpoint returns an endpoint via the given service.
func MakeCardGetEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		tr := otel.Tracer("Get Users")
		ctx, span := tr.Start(ctx, "Get Users")
		span.SetAttributes(attribute.Key("service").String("user"))
		defer span.End()

		req := request.(GetRequest)
		ctx, cardspan := tr.Start(ctx, "card from db")
		cards, err := s.GetCards(req.ID)
		cardspan.End()
		if req.ID == "" {
			return EmbedStruct{cardsResponse{Cards: cards}}, err
		}
		if len(cards) == 0 {
			return users.Card{}, err
		}
		return cards[0], err
	}
}

// MakeCardPostEndpoint returns an endpoint via the given service.
func MakeCardPostEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		tr := otel.Tracer("Post Card")
		ctx, span := tr.Start(ctx, "Post Card")
		span.SetAttributes(attribute.Key("service").String("user"))
		defer span.End()
		req := request.(cardPostRequest)
		id, err := s.PostCard(req.Card, req.UserID)
		return postResponse{ID: id}, err
	}
}

// MakeLoginEndpoint returns an endpoint via the given service.
func MakeDeleteEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		tr := otel.Tracer("Delete Entity")
		ctx, span := tr.Start(ctx, "Delete Entity")
		span.SetAttributes(attribute.Key("service").String("user"))
		defer span.End()
		req := request.(deleteRequest)
		err = s.Delete(req.Entity, req.ID)
		if err == nil {
			return statusResponse{Status: true}, err
		}
		return statusResponse{Status: false}, err
	}
}

// MakeHealthEndpoint returns current health of the given service.
func MakeHealthEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		tr := otel.Tracer("Health Check")
		ctx, span := tr.Start(ctx, "Health Check")
		span.SetAttributes(attribute.Key("service").String("user"))
		defer span.End()
		health := s.Health()
		return healthResponse{Health: health}, nil
	}
}

type GetRequest struct {
	ID   string
	Attr string
}

type loginRequest struct {
	Username string
	Password string
}

type userResponse struct {
	User users.User `json:"user"`
}

type usersResponse struct {
	Users []users.User `json:"customer"`
}

type addressPostRequest struct {
	users.Address
	UserID string `json:"userID"`
}

type addressesResponse struct {
	Addresses []users.Address `json:"address"`
}

type cardPostRequest struct {
	users.Card
	UserID string `json:"userID"`
}

type cardsResponse struct {
	Cards []users.Card `json:"card"`
}

type registerRequest struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type statusResponse struct {
	Status bool `json:"status"`
}

type postResponse struct {
	ID string `json:"id"`
}

type deleteRequest struct {
	Entity string
	ID     string
}

type healthRequest struct {
	//
}

type healthResponse struct {
	Health []Health `json:"health"`
}

type EmbedStruct struct {
	Embed interface{} `json:"_embedded"`
}
