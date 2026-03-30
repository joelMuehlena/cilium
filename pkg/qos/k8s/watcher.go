// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Cilium

package k8s

import (
	"context"
	"log/slog"

	cilium_v2alpha1 "github.com/cilium/cilium/pkg/k8s/apis/cilium.io/v2alpha1"
	"github.com/cilium/cilium/pkg/k8s/resource"
	k8sSynced "github.com/cilium/cilium/pkg/k8s/synced"
	"github.com/cilium/cilium/pkg/option"
)

type qosPolicyWatcher struct {
	log    *slog.Logger
	config *option.DaemonConfig

	k8sResourceSynced *k8sSynced.Resources
	k8sAPIGroups      *k8sSynced.APIGroups

	ciliumQoSPolicies resource.Resource[*cilium_v2alpha1.CiliumQoSPolicy]
}

func (p *qosPolicyWatcher) watchResources(ctx context.Context) {
	go func() {
		var (
			qosPolicyEvents <-chan resource.Event[*cilium_v2alpha1.CiliumQoSPolicy]
		)

		if p.config.EnableCiliumQoSPolicy {
			qosPolicyEvents = p.ciliumQoSPolicies.Events(ctx)
		}

		for event := range qosPolicyEvents {
			switch event.Kind {
			case resource.Sync:
				p.log.Info("Sync policy qos")
			// Ensure the policy is set in the datapath?????
			case resource.Upsert:
				p.log.Info("Upsert policy qos")
			// Ensure the policy is set in the datapath
			case resource.Delete:
				// Ensure the policy is removed in the datapath
				p.log.Info("Delete policy qos")
			}

			event.Done(err)
		}
	}()
}

func (p *qosPolicyWatcher) registerResourceWithSyncFn(ctx context.Context, resource string, syncFn func() bool) {
	p.k8sResourceSynced.BlockWaitGroupToSyncResources(ctx.Done(), nil, syncFn, resource)
	p.k8sAPIGroups.AddAPI(resource)
}
