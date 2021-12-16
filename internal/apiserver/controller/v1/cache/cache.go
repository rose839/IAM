// Package cache defines a cache service which can return all secrets and policies.
package cache

import (
	"context"
	"fmt"
	"sync"

	metav1 "github.com/rose839/IAM/api/meta/v1"
	pb "github.com/rose839/IAM/api/proto/apiserver/v1"
	"github.com/rose839/IAM/internal/apiserver/store"
	"github.com/rose839/IAM/internal/pkg/code"
	"github.com/rose839/IAM/pkg/errors"
)

type Cache struct {
	store store.Factory
}

var (
	cacheServer *Cache
	once        sync.Once
)

// GetCacheInsOr return cache server instance with given factory.
func GetCacheInsOr(store store.Factory) (*Cache, error) {
	if store != nil {
		once.Do(func() {
			cacheServer = &Cache{store}
		})
	}

	if cacheServer == nil {
		return nil, fmt.Errorf("got nil cache server")
	}

	return cacheServer, nil
}

// ListSecrets returns all secrets.
func (c *Cache) ListSecrets(ctx context.Context, r *pb.ListSecretsRequest) (*pb.ListSecretsResponse, error) {
	opts := metav1.ListOptions{
		Offset: r.Offset,
		Limit:  r.Limit,
	}

	secrets, err := c.store.Secrets().List(ctx, "", opts)
	if err != nil {
		return nil, errors.WithCode(code.ErrDatabase, err.Error())
	}

	items := make([]*pb.SecretInfo, 0)
	for _, secret := range secrets.Items {
		items = append(items, &pb.SecretInfo{
			SecretId:    secret.SecretID,
			Username:    secret.Username,
			SecretKey:   secret.SecretKey,
			Expires:     secret.Expires,
			Description: secret.Description,
			CreatedAt:   secret.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:   secret.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return &pb.ListSecretsResponse{
		TotalCount: secrets.TotalCount,
		Items:      items,
	}, nil
}

// ListPolicies returns all policies.
func (c *Cache) ListPolicies(ctx context.Context, r *pb.ListPoliciesRequest) (*pb.ListPoliciesResponse, error) {
	opts := metav1.ListOptions{
		Offset: r.Offset,
		Limit:  r.Limit,
	}

	policies, err := c.store.Policies().List(ctx, "", opts)
	if err != nil {
		return nil, errors.WithCode(code.ErrDatabase, err.Error())
	}

	items := make([]*pb.PolicyInfo, 0)
	for _, pol := range policies.Items {
		items = append(items, &pb.PolicyInfo{
			Name:         pol.Name,
			Username:     pol.Username,
			PolicyShadow: pol.PolicyShadow,
			CreatedAt:    pol.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return &pb.ListPoliciesResponse{
		TotalCount: policies.TotalCount,
		Items:      items,
	}, nil
}
