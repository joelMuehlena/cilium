// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Cilium

package k8s

import (
	"context"
	"log/slog"

	"github.com/cilium/hive/cell"

	cilium_v2alpha1 "github.com/cilium/cilium/pkg/k8s/apis/cilium.io/v2alpha1"
	"github.com/cilium/cilium/pkg/k8s/client"
	"github.com/cilium/cilium/pkg/k8s/resource"
	"github.com/cilium/cilium/pkg/k8s/synced"
	"github.com/cilium/cilium/pkg/option"
)

const (
	k8sAPIGroupCiliumNetworkPolicyV2 = "cilium/v2alpha1::CiliumQoSPolicy"
)

// FIXME:
var Cell = cell.Module(
	"qos-policy-watcher",
	"Watches Cilium QoS policies",

	cell.Invoke(startCiliumQoSPolicyWatcher),
)

type PolicyWatcherParams struct {
	cell.In

	Lifecycle cell.Lifecycle

	ClientSet client.Clientset
	Config    *option.DaemonConfig
	Logger    *slog.Logger

	K8sResourceSynced *synced.Resources
	K8sAPIGroups      *synced.APIGroups

	CiliumQoSPolicies resource.Resource[*cilium_v2alpha1.CiliumQoSPolicy]
}

func startCiliumQoSPolicyWatcher(params PolicyWatcherParams) {
	if !params.ClientSet.IsEnabled() {
		return // skip watcher if K8s is not enabled
	}

	// We want to subscribe before the start hook is invoked in order to not miss
	// any events
	ctx, cancel := context.WithCancel(context.Background())

	p := &qosPolicyWatcher{
		log:               params.Logger,
		config:            params.Config,
		k8sResourceSynced: params.K8sResourceSynced,
		k8sAPIGroups:      params.K8sAPIGroups,
		ciliumQoSPolicies: params.CiliumQoSPolicies,
	}

	params.Lifecycle.Append(cell.Hook{
		OnStart: func(startCtx cell.HookContext) error {
			p.watchResources(ctx)
			return nil
		},
		OnStop: func(cell.HookContext) error {
			if cancel != nil {
				cancel()
			}
			return nil
		},
	})

	if params.Config.EnableCiliumQoSPolicy {
		p.registerResourceWithSyncFn(ctx, k8sAPIGroupCiliumNetworkPolicyV2, func() bool {
			return true
		})
	}
}
