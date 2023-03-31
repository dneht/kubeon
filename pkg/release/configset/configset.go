/*
Copyright (c) 2020, Dash

Licensed under the LGPL, Version 3.0 (the "License");
you may not use this file except in compliance with the License.
*/

package configset

func GetKubeletByVersion(version string) string {
	switch version {
	default:
		return kubeletYaml
	}
}

func GetKubeadmByVersion(version string) string {
	switch version {
	default:
		return kubeadmYaml
	}
}

func GetKubeadmInitByVersion(version string) string {
	switch version {
	default:
		return kubeadmInitYaml
	}
}

func GetKubeadmJoinByVersion(version string) string {
	switch version {
	default:
		return kubeadmJoinYaml
	}
}

func GetHealthzReaderByVersion(version string) string {
	switch version {
	default:
		return healthzReaderYaml
	}
}

func GetHaproxyStaticTemplate() string {
	return haproxyStaticYaml
}

func GetApiserverStartupService() string {
	return apiserverStartupService
}

func GetApiserverStartupScript() string {
	return apiserverStartupBash
}

func GetApiserverUpdaterTemplate() string {
	return apiserverUpdaterYaml
}
