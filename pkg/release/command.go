/*
Copyright 2020 Dasheng.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package release

import (
	"fmt"
	"github.com/dneht/kubeon/pkg/define"
	"strconv"
)

func BuildIstioctlArgs(input *IstioTemplate, install, local bool) []string {
	imageHub := input.MirrorUrl + define.IstioMirrorImagePrefix
	if local {
		imageHub = define.DockerImageRepo + define.IstioImagePrefix
	}
	istioArgs := make([]string, 0, 80)
	if install {
		istioArgs = append(istioArgs, "install")
	} else {
		istioArgs = append(istioArgs, "manifest", "generate")
	}
	autoInject := define.IstioProxyAutoInjectEnable
	if !input.EnableAutoInject {
		autoInject = define.IstioProxyAutoInjectDisable
	}
	clusterDomain := input.ProxyClusterDomain
	if "" == clusterDomain {
		clusterDomain = define.DefaultClusterDNSDomain
	}
	istioArgs = append(istioArgs, "--set", "profile=default", "--set", "values.global.hub="+imageHub,
		"--set", "values.global.imagePullPolicy="+define.ImagePullPolicyNotPresent,
		"--set", "values.global.proxy.autoInject="+autoInject, "--set", "values.global.proxy.clusterDomain="+clusterDomain)
	if !local {
		istioArgs = append(istioArgs, "--set", "values.global.proxy.image="+define.IstioProxyImage,
			"--set", "values.global.proxy_init.image="+define.IstioProxyImage,
			"--set", "values.pilot.image="+define.IstioPilotImage, "--set", "values.cni.image="+define.IstioCNIImage)
	}
	if input.EnableNetworkPlugin {
		istioArgs = append(istioArgs, "--set", "components.cni.enabled=true",
			"--set", "values.cni.excludeNamespaces[0]=kube-system", "--set", "values.cni.excludeNamespaces[1]=istio-system")
	}
	istioArgs = append(istioArgs, "--set", "meshConfig.enableAutoMtls="+strconv.FormatBool(input.EnableAutoMTLS))
	if input.EnableHttp2AutoUpgrade {
		istioArgs = append(istioArgs, "--set", "meshConfig.h2UpgradePolicy="+define.IstioHttp2AutoUpgrade)
	} else {
		istioArgs = append(istioArgs, "--set", "meshConfig.h2UpgradePolicy="+define.IstioHttp2DontAutoUpgrade)
	}
	if input.EnableIngressGateway {
		istioArgs = append(istioArgs, "--set", "values.gateways.istio-ingressgateway.enabled=true")
		if "" != input.IngressGatewayType {
			istioArgs = append(istioArgs, "--set", "values.gateways.istio-ingressgateway.type="+input.IngressGatewayType)
		}
	} else {
		istioArgs = append(istioArgs, "--set", "values.gateways.istio-ingressgateway.enabled=false")
	}
	if input.EnableEgressGateway {
		istioArgs = append(istioArgs, "--set", "values.gateways.istio-egressgateway.enabled=true")
		if "" != input.EgressGatewayType {
			istioArgs = append(istioArgs, "--set", "values.gateways.istio-egressgateway.type="+input.EgressGatewayType)
		}
	} else {
		istioArgs = append(istioArgs, "--set", "values.gateways.istio-egressgateway.enabled=false")
	}
	if len(input.ServiceEntryExportTo) > 0 {
		for idx, exportTo := range input.ServiceEntryExportTo {
			istioArgs = append(istioArgs, "--set", fmt.Sprintf("meshConfig.defaultServiceExportTo[%d]=%s", idx, exportTo))
		}
	}
	if len(input.VirtualServiceExportTo) > 0 {
		for idx, exportTo := range input.VirtualServiceExportTo {
			istioArgs = append(istioArgs, "--set", fmt.Sprintf("meshConfig.defaultVirtualServiceExportTo[%d]=%s", idx, exportTo))
		}
	}
	if len(input.DestinationRuleExportTo) > 0 {
		for idx, exportTo := range input.DestinationRuleExportTo {
			istioArgs = append(istioArgs, "--set", fmt.Sprintf("meshConfig.defaultDestinationRuleExportTo[%d]=%s", idx, exportTo))
		}
	}
	if input.EnableSkywalking || input.EnableSkywalkingAll {
		istioArgs = append(istioArgs, "--set", "meshConfig.enableTracing=true",
			"--set", "meshConfig.extensionProviders[0].name=tracer.skywalking",
			"--set", fmt.Sprintf("meshConfig.extensionProviders[0].skywalking.service=%s", input.SkywalkingService),
			"--set", fmt.Sprintf("meshConfig.extensionProviders[0].skywalking.port=%d", input.SkywalkingPort))
		if "" != input.SkywalkingAccessToken {
			istioArgs = append(istioArgs, "--set", fmt.Sprintf("meshConfig.extensionProviders[0].skywalking.accessToken=%s", input.SkywalkingAccessToken))
		}
		istioArgs = append(istioArgs, "--set", "meshConfig.defaultProviders.tracing[0]=tracer.skywalking")
		if input.EnableSkywalkingAll {
			input.MetricsServiceAddress = fmt.Sprintf("%s:%d", input.SkywalkingService, input.SkywalkingPort)
		}
	} else if input.EnableZipkin {
		istioArgs = append(istioArgs, "--set", "meshConfig.enableTracing=true",
			"--set", "meshConfig.extensionProviders[0].name=tracer.zipkin",
			"--set", fmt.Sprintf("meshConfig.extensionProviders[0].zipkin.service=%s", input.ZipkinService),
			"--set", fmt.Sprintf("meshConfig.extensionProviders[0].zipkin.port=%d", input.ZipkinPort))
		istioArgs = append(istioArgs, "--set", "meshConfig.defaultProviders.tracing[0]=tracer.zipkin",
			"--set", fmt.Sprintf("meshConfig.defaultConfig.tracing.zipkin.address=%s:%d", input.ZipkinService, input.ZipkinPort))
	}
	if "" != input.MetricsServiceAddress {
		istioArgs = append(istioArgs,
			"--set", "meshConfig.defaultConfig.envoyMetricsService.address="+input.MetricsServiceAddress)
	}
	if "" != input.AccessLogServiceAddress {
		istioArgs = append(istioArgs, "--set", "meshConfig.enableEnvoyAccessLogService=true",
			"--set", "meshConfig.defaultConfig.envoyAccessLogService.address="+input.AccessLogServiceAddress)
	}
	istioArgs = append(istioArgs, "--set", "meshConfig.defaultConfig.proxyStatsMatcher.inclusionRegexps[0]=.*membership_healthy.*",
		"--set", "meshConfig.defaultConfig.proxyStatsMatcher.inclusionRegexps[1]=.*upstream_cx_active.*",
		"--set", "meshConfig.defaultConfig.proxyStatsMatcher.inclusionRegexps[2]=.*upstream_cx_total.*",
		"--set", "meshConfig.defaultConfig.proxyStatsMatcher.inclusionRegexps[3]=.*upstream_rq_active.*",
		"--set", "meshConfig.defaultConfig.proxyStatsMatcher.inclusionRegexps[4]=.*upstream_rq_total.*",
		"--set", "meshConfig.defaultConfig.proxyStatsMatcher.inclusionRegexps[5]=.*upstream_rq_pending_active.*",
		"--set", "meshConfig.defaultConfig.proxyStatsMatcher.inclusionRegexps[6]=.*lb_healthy_panic.*",
		"--set", "meshConfig.defaultConfig.proxyStatsMatcher.inclusionRegexps[7]=.*upstream_cx_none_healthy.*")
	return istioArgs
}
