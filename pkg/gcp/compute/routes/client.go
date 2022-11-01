package routes

import (
	"context"
	"fmt"

	compute "cloud.google.com/go/compute/apiv1"
	"github.com/go-logr/logr"
	"sigs.k8s.io/controller-runtime/pkg/log"

	"google.golang.org/api/iterator"
	computepb "google.golang.org/genproto/googleapis/cloud/compute/v1"
	capg "sigs.k8s.io/cluster-api-provider-gcp/api/v1beta1"
)

type Client struct {
	routesClient *compute.RoutesClient
}

func NewClient(client *compute.RoutesClient) *Client {
	return &Client{
		routesClient: client,
	}
}

func (c *Client) DeleteRoutes(ctx context.Context, cluster *capg.GCPCluster) error {
	logger := c.getLogger(ctx)
	project := cluster.Spec.Project
	filter := fmt.Sprintf(`name : "%s*"`, cluster.Name)

	req := &computepb.ListRoutesRequest{
		Filter:               &filter,
		Project:              project,
	}

	it := c.routesClient.List(ctx, req)
	for {
		route, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			logger.Error(err, "failed to list routes")
			return err
		}

		message := fmt.Sprintf("Deleting route - %s for cluster %v", *route.Name, cluster.Name)
		logger.Info(message)

		err = c.deleteRoute(ctx, route, project)
		if err != nil {
			return err
		}

		logger.Info("Deleted!")
	}

	return nil
}

func (c *Client) deleteRoute(ctx context.Context, route *computepb.Route, project string) error {
	req := &computepb.DeleteRouteRequest{
		Project:   project,
		Route:     *route.Name,
	}

	op, err := c.routesClient.Delete(ctx, req)
	if err != nil {
		return err
	}

	err = op.Wait(ctx)
	return err
}

func (c *Client) getLogger(ctx context.Context) logr.Logger {
	logger := log.FromContext(ctx)
	return logger.WithName("gcp-compute-routes-client")
}
