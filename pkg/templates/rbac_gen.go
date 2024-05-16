// Code generated by go generate; DO NOT EDIT.

package main

//+kubebuilder:rbac:groups="",resources=configmaps,verbs=create;get;list;watch;delete;update
//+kubebuilder:rbac:groups="",resources=configmaps,verbs=get;list;patch;update;watch;create
//+kubebuilder:rbac:groups="",resources=configmaps,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups="",resources=configmaps,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups="",resources=configmaps,verbs=get;update
//+kubebuilder:rbac:groups="",resources=configmaps;endpoints;events;secrets;serviceaccounts;services;services/proxy,verbs=get;list;watch;create;update;patch;delete;deletecollection
//+kubebuilder:rbac:groups="",resources=configmaps;events,verbs=get;list;watch;create;update;delete;deletecollection;patch
//+kubebuilder:rbac:groups="",resources=configmaps;jobs;namespaces;pods;secrets,verbs=list;watch
//+kubebuilder:rbac:groups="",resources=configmaps;secrets;serviceaccounts;services;persistentvolumeclaims;pods;endpoints,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups="",resources=events,verbs=create;delete;get;list;patch;update;watch
//+kubebuilder:rbac:groups="",resources=events,verbs=create;get;list;patch;update;watch
//+kubebuilder:rbac:groups="",resources=events,verbs=create;patch
//+kubebuilder:rbac:groups="",resources=events,verbs=create;patch
//+kubebuilder:rbac:groups="",resources=events,verbs=create;patch
//+kubebuilder:rbac:groups="",resources=events;secrets;configmaps;serviceaccounts;services,verbs=create;delete;get;list;patch;update;watch
//+kubebuilder:rbac:groups="",resources=groups;secrets;serviceaccounts;services;users,verbs=create;delete;get;impersonate;list;patch;update;watch
//+kubebuilder:rbac:groups="",resources=namespaces;serviceaccounts,verbs=create;get;list;watch;update;patch;delete
//+kubebuilder:rbac:groups="",resources=pods,verbs=get
//+kubebuilder:rbac:groups="",resources=pods,verbs=get
//+kubebuilder:rbac:groups="",resources=pods,verbs=get;list
//+kubebuilder:rbac:groups="",resources=pods,verbs=get;list;watch
//+kubebuilder:rbac:groups="",resources=pods;services;endpoints,verbs=get;list;watch
//+kubebuilder:rbac:groups="",resources=pods;services;services/finalizers;endpoints;persistentvolumeclaims;events;configmaps;secrets,verbs=create;delete;get;list;patch;update;watch
//+kubebuilder:rbac:groups="",resources=pods;services;services/finalizers;endpoints;persistentvolumeclaims;events;configmaps;secrets;serviceaccounts;namespaces;nodes,verbs=create;delete;get;list;patch;update;watch
//+kubebuilder:rbac:groups="",resources=secrets,verbs=create
//+kubebuilder:rbac:groups="",resources=secrets,verbs=create
//+kubebuilder:rbac:groups="",resources=secrets,verbs=create;get;list;watch;delete
//+kubebuilder:rbac:groups="",resources=secrets,verbs=delete;get;list;patch;update;watch
//+kubebuilder:rbac:groups="",resources=secrets,verbs=delete;patch;update;watch
//+kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch
//+kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch
//+kubebuilder:rbac:groups="",resources=serviceaccounts,verbs=create
//+kubebuilder:rbac:groups="",resources=serviceaccounts,verbs=delete;get;patch;update
//+kubebuilder:rbac:groups="",resources=services,verbs=list
//+kubebuilder:rbac:groups="";events.k8s.io,resources=events,verbs=create;patch;update
//+kubebuilder:rbac:groups=*,resources=*,verbs=*
//+kubebuilder:rbac:groups=*,resources=*,verbs=*
//+kubebuilder:rbac:groups=*,resources=*,verbs=get;list;watch
//+kubebuilder:rbac:groups=*,resources=*,verbs=get;list;watch
//+kubebuilder:rbac:groups=addon.open-cluster-management.io,resources=addondeploymentconfigs,verbs=get;list;watch
//+kubebuilder:rbac:groups=addon.open-cluster-management.io,resources=addondeploymentconfigs,verbs=get;list;watch
//+kubebuilder:rbac:groups=addon.open-cluster-management.io,resources=addondeploymentconfigs,verbs=get;list;watch
//+kubebuilder:rbac:groups=addon.open-cluster-management.io,resources=addondeploymentconfigs,verbs=get;list;watch
//+kubebuilder:rbac:groups=addon.open-cluster-management.io,resources=addondeploymentconfigs;clustermanagementaddons;managedclusteraddons,verbs=create;delete;get;list;update
//+kubebuilder:rbac:groups=addon.open-cluster-management.io,resources=clustermanagementaddons,verbs=create;get;list;watch;update
//+kubebuilder:rbac:groups=addon.open-cluster-management.io,resources=clustermanagementaddons,verbs=get;list;watch
//+kubebuilder:rbac:groups=addon.open-cluster-management.io,resources=clustermanagementaddons,verbs=get;list;watch
//+kubebuilder:rbac:groups=addon.open-cluster-management.io,resources=clustermanagementaddons,verbs=get;list;watch;patch
//+kubebuilder:rbac:groups=addon.open-cluster-management.io,resources=clustermanagementaddons/finalizers,verbs=update
//+kubebuilder:rbac:groups=addon.open-cluster-management.io,resources=clustermanagementaddons/finalizers,verbs=update
//+kubebuilder:rbac:groups=addon.open-cluster-management.io,resources=clustermanagementaddons/finalizers,verbs=update
//+kubebuilder:rbac:groups=addon.open-cluster-management.io,resources=clustermanagementaddons/finalizers;managedclusteraddons/finalizers,verbs=update
//+kubebuilder:rbac:groups=addon.open-cluster-management.io,resources=clustermanagementaddons/status,verbs=patch;update
//+kubebuilder:rbac:groups=addon.open-cluster-management.io,resources=clustermanagementaddons/status,verbs=patch;update
//+kubebuilder:rbac:groups=addon.open-cluster-management.io,resources=clustermanagementaddons/status,verbs=update;patch
//+kubebuilder:rbac:groups=addon.open-cluster-management.io,resources=clustermanagementaddons/status;managedclusteraddons/status,verbs=get;patch;update
//+kubebuilder:rbac:groups=addon.open-cluster-management.io,resources=clustermanagementaddons;clustermanagementaddons/finalizers,verbs=create;update;get;delete;list;watch
//+kubebuilder:rbac:groups=addon.open-cluster-management.io,resources=clustermanagementaddons;managedclusteraddons,verbs=list;watch
//+kubebuilder:rbac:groups=addon.open-cluster-management.io,resources=managedclusteraddons,verbs=create;get;list;update;watch
//+kubebuilder:rbac:groups=addon.open-cluster-management.io,resources=managedclusteraddons,verbs=delete
//+kubebuilder:rbac:groups=addon.open-cluster-management.io,resources=managedclusteraddons,verbs=get;list;watch
//+kubebuilder:rbac:groups=addon.open-cluster-management.io,resources=managedclusteraddons,verbs=get;list;watch;create;update;delete
//+kubebuilder:rbac:groups=addon.open-cluster-management.io,resources=managedclusteraddons,verbs=get;list;watch;patch;update;delete
//+kubebuilder:rbac:groups=addon.open-cluster-management.io,resources=managedclusteraddons/finalizers,verbs=update
//+kubebuilder:rbac:groups=addon.open-cluster-management.io,resources=managedclusteraddons/finalizers,verbs=update
//+kubebuilder:rbac:groups=addon.open-cluster-management.io,resources=managedclusteraddons/finalizers,verbs=update
//+kubebuilder:rbac:groups=addon.open-cluster-management.io,resources=managedclusteraddons/status,verbs=patch;update
//+kubebuilder:rbac:groups=addon.open-cluster-management.io,resources=managedclusteraddons/status,verbs=patch;update
//+kubebuilder:rbac:groups=addon.open-cluster-management.io,resources=managedclusteraddons/status,verbs=update;patch
//+kubebuilder:rbac:groups=addon.open-cluster-management.io,resources=managedclusteraddons;managedclusteraddons/status;clustermanagementaddons,verbs=create;delete;get;list;patch;update;watch
//+kubebuilder:rbac:groups=addon.open-cluster-management.io,resources=managedclusteraddons;managedclusteraddons/status;managedclusteraddons/finalizers,verbs=watch;create;update;delete;get;list;patch
//+kubebuilder:rbac:groups=addon.open-cluster-management.io;agent.open-cluster-management.io;apps.open-cluster-management.io;cluster.open-cluster-management.io;operator.open-cluster-management.io;work.open-cluster-management.io;view.open-cluster-management.io;authentication.open-cluster-management.io;policy.open-cluster-management.io,resources=channels;channels/status;channels/finalizers;deployables;deployables/status;gitopsclusters;gitopsclusters/status;helmreleases;helmreleases/status;klusterletaddonconfigs;manifestworks;manifestworks/status;managedclusters;managedclusterviews;managedclusterviews/status;managedclusteraddons;managedserviceaccounts;multiclusterhubs;placements;placements/status;placement/finalizers;placementbindings;placementbindings/finalizers;placementdecisions;placementdecisions/status;placementdecisions/finalizers;placementrules;placementrules/status;placementrules/finalizers;subscriptions;subscriptions/finalizers;subscriptions/status;subscriptionstatuses;subscriptionreports;multiclusterapplicationsetreports;multiclusterapplicationsetreports/status;policies,verbs=get;list;watch;update;patch;create;delete;deletecollection
//+kubebuilder:rbac:groups=admissionregistration.k8s.io,resources=mutatingwebhookconfigurations;validatingwebhookconfigurations,verbs=create;delete;get;list;patch;update;watch
//+kubebuilder:rbac:groups=admissionregistration.k8s.io;certificates.k8s.io;coordination.k8s.io;apiextensions.k8s.io,resources=certificatesigningrequests;customresourcedefinitions;leases;mutatingwebhookconfigurations;validatingwebhookconfigurations,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=agent-install.openshift.io,resources=agents;infraenvs;nmstateconfigs;agentserviceconfigs,verbs=list;watch
//+kubebuilder:rbac:groups=agent.open-cluster-management.io,resources=klusterletaddonconfigs;klusterletaddonconfigs/finalizers;klusterletaddonconfigs/status,verbs=create;delete;get;list;patch;update;watch
//+kubebuilder:rbac:groups=apiextensions.k8s.io,resources=customresourcedefinitions,verbs=*
//+kubebuilder:rbac:groups=apiextensions.k8s.io,resources=customresourcedefinitions,verbs=list;watch
//+kubebuilder:rbac:groups=apiextensions.k8s.io,resources=customresourcedefinitions;customresourcedefinitions/finalizers,verbs=create;get;list;watch;update;patch;delete
//+kubebuilder:rbac:groups=app.k8s.io,resources=applications,verbs=list;watch
//+kubebuilder:rbac:groups=app.k8s.io;argoproj.io,resources=applications;applications/status;applicationsets;applicationsets/status,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=create;delete;get;list;patch;update;watch
//+kubebuilder:rbac:groups=apps,resources=deployments/finalizers,verbs=update
//+kubebuilder:rbac:groups=apps,resources=deployments/finalizers,verbs=update
//+kubebuilder:rbac:groups=apps,resources=deployments;daemonsets;replicasets;statefulsets,verbs=create;delete;get;list;patch;update;watch
//+kubebuilder:rbac:groups=apps,resources=deployments;daemonsets;replicasets;statefulsets,verbs=create;delete;get;list;patch;update;watch
//+kubebuilder:rbac:groups=apps,resources=deployments;deployments/finalizers,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=apps,resources=deployments;replicasets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=apps,resources=replicasets,verbs=get
//+kubebuilder:rbac:groups=apps,resources=replicasets,verbs=get;list
//+kubebuilder:rbac:groups=apps,resources=replicasets;deployments,verbs=get
//+kubebuilder:rbac:groups=apps,resources=replicasets;deployments,verbs=get
//+kubebuilder:rbac:groups=apps.open-cluster-management.io,resources=*,verbs=update
//+kubebuilder:rbac:groups=apps.open-cluster-management.io,resources=channels;gitopsclusters;helmreleases;placementrules;subscriptions;subscriptionreports;multiclusterapplicationsetreports,verbs=list;watch
//+kubebuilder:rbac:groups=apps.open-cluster-management.io,resources=placementrules,verbs=get;list;watch
//+kubebuilder:rbac:groups=apps.open-cluster-management.io,resources=subscriptions,verbs=get;list;watch
//+kubebuilder:rbac:groups=argoproj.io,resources=applications;applicationsets;argocds,verbs=list;watch
//+kubebuilder:rbac:groups=authentication.k8s.io,resources=tokenreviews,verbs=create
//+kubebuilder:rbac:groups=authentication.k8s.io,resources=tokenreviews,verbs=create
//+kubebuilder:rbac:groups=authentication.k8s.io,resources=tokenreviews,verbs=create
//+kubebuilder:rbac:groups=authentication.k8s.io,resources=tokenreviews,verbs=create
//+kubebuilder:rbac:groups=authentication.k8s.io,resources=tokenreviews,verbs=create
//+kubebuilder:rbac:groups=authentication.k8s.io;authorization.k8s.io,resources=uids;userextras/authentication.kubernetes.io/pod-name;userextras/authentication.kubernetes.io/pod-uid,verbs=impersonate
//+kubebuilder:rbac:groups=authentication.open-cluster-management.io,resources=managedserviceaccounts,verbs=create;delete
//+kubebuilder:rbac:groups=authentication.open-cluster-management.io,resources=managedserviceaccounts,verbs=get;list;watch
//+kubebuilder:rbac:groups=authorization.k8s.io,resources=subjectaccessreviews,verbs=create
//+kubebuilder:rbac:groups=authorization.k8s.io,resources=subjectaccessreviews,verbs=create
//+kubebuilder:rbac:groups=authorization.k8s.io,resources=subjectaccessreviews,verbs=create
//+kubebuilder:rbac:groups=authorization.k8s.io,resources=subjectaccessreviews,verbs=create;get
//+kubebuilder:rbac:groups=authorization.k8s.io,resources=subjectaccessreviews,verbs=create;get
//+kubebuilder:rbac:groups=authorization.k8s.io,resources=subjectaccessreviews,verbs=get;create
//+kubebuilder:rbac:groups=capi-provider.agent-install.openshift.io,resources=agentmachines,verbs=list;watch
//+kubebuilder:rbac:groups=certificates.k8s.io,resources=certificatesigningrequests,verbs=get;list;watch
//+kubebuilder:rbac:groups=certificates.k8s.io,resources=certificatesigningrequests,verbs=get;list;watch
//+kubebuilder:rbac:groups=certificates.k8s.io,resources=certificatesigningrequests,verbs=get;list;watch
//+kubebuilder:rbac:groups=certificates.k8s.io,resources=certificatesigningrequests,verbs=list;watch
//+kubebuilder:rbac:groups=certificates.k8s.io,resources=certificatesigningrequests/approval,verbs=update;patch
//+kubebuilder:rbac:groups=certificates.k8s.io,resources=certificatesigningrequests/status,verbs=update
//+kubebuilder:rbac:groups=certificates.k8s.io,resources=certificatesigningrequests/status;certificatesigningrequests/approval,verbs=update
//+kubebuilder:rbac:groups=certificates.k8s.io,resources=certificatesigningrequests;certificatesigningrequests/approval,verbs=create;get;list;update;watch
//+kubebuilder:rbac:groups=certificates.k8s.io,resources=certificatesigningrequests;certificatesigningrequests/approval,verbs=create;get;list;update;watch
//+kubebuilder:rbac:groups=certificates.k8s.io,resources=certificatesigningrequests;certificatesigningrequests/approval,verbs=get;list;watch;create;update
//+kubebuilder:rbac:groups=certificates.k8s.io,resources=signers,verbs=approve
//+kubebuilder:rbac:groups=certificates.k8s.io,resources=signers,verbs=approve
//+kubebuilder:rbac:groups=certificates.k8s.io,resources=signers,verbs=approve
//+kubebuilder:rbac:groups=certificates.k8s.io,resources=signers,verbs=approve
//+kubebuilder:rbac:groups=certificates.k8s.io,resources=signers,verbs=approve
//+kubebuilder:rbac:groups=certificates.k8s.io,resources=signers,verbs=sign
//+kubebuilder:rbac:groups=certmanager.k8s.io,resources=*,verbs=*
//+kubebuilder:rbac:groups=cluster.open-cluster-management.io,resources=clusterclaims,verbs=get
//+kubebuilder:rbac:groups=cluster.open-cluster-management.io,resources=manageclusters,verbs=get;list;watch
//+kubebuilder:rbac:groups=cluster.open-cluster-management.io,resources=managedclusters,verbs=get;list
//+kubebuilder:rbac:groups=cluster.open-cluster-management.io,resources=managedclusters,verbs=get;list;watch
//+kubebuilder:rbac:groups=cluster.open-cluster-management.io,resources=managedclusters,verbs=get;list;watch
//+kubebuilder:rbac:groups=cluster.open-cluster-management.io,resources=managedclusters,verbs=get;list;watch
//+kubebuilder:rbac:groups=cluster.open-cluster-management.io,resources=managedclusters,verbs=list;get;watch
//+kubebuilder:rbac:groups=cluster.open-cluster-management.io,resources=managedclusters,verbs=watch;get;list
//+kubebuilder:rbac:groups=cluster.open-cluster-management.io,resources=managedclusters;managedclusters/finalizers,verbs=create;get;list;patch;update;watch
//+kubebuilder:rbac:groups=cluster.open-cluster-management.io,resources=managedclusters;managedclustersets,verbs=get;list;watch;update;patch
//+kubebuilder:rbac:groups=cluster.open-cluster-management.io,resources=managedclusters;managedclustersets;managedclustersetbindings;clustercurators;placements;placementdecisions,verbs=list;watch
//+kubebuilder:rbac:groups=cluster.open-cluster-management.io,resources=managedclusters;placementdecisions;placements,verbs=get;list;watch
//+kubebuilder:rbac:groups=config.openshift.io,resources=*;infrastructures,verbs=*
//+kubebuilder:rbac:groups=config.openshift.io,resources=apiservers;infrastructures;infrastructures/status,verbs=get
//+kubebuilder:rbac:groups=config.openshift.io,resources=clusterversions,verbs=list;get;watch
//+kubebuilder:rbac:groups=config.openshift.io,resources=infrastructures,verbs=get;list;watch
//+kubebuilder:rbac:groups=config.openshift.io,resources=infrastructures,verbs=get;list;watch
//+kubebuilder:rbac:groups=config.openshift.io;console.openshift.io;project.openshift.io;tower.ansible.com,resources=infrastructures;consolelinks;projects;featuregates;ansiblejobs;clusterversions,verbs=list;get;watch
//+kubebuilder:rbac:groups=console.open-cluster-management.io,resources=userpreferences,verbs=create;get;list;patch;watch
//+kubebuilder:rbac:groups=coordination.k8s.io,resources=leases,verbs=create
//+kubebuilder:rbac:groups=coordination.k8s.io,resources=leases,verbs=create;get;list;patch;update;watch
//+kubebuilder:rbac:groups=coordination.k8s.io,resources=leases,verbs=create;get;list;patch;update;watch;delete
//+kubebuilder:rbac:groups=coordination.k8s.io,resources=leases,verbs=delete;get;list;watch;create;update;patch
//+kubebuilder:rbac:groups=coordination.k8s.io,resources=leases,verbs=get;list;patch;update;watch
//+kubebuilder:rbac:groups=coordination.k8s.io,resources=leases,verbs=get;list;watch;create;update;delete
//+kubebuilder:rbac:groups=coordination.k8s.io,resources=leases,verbs=get;list;watch;create;update;patch
//+kubebuilder:rbac:groups=coordination.k8s.io,resources=leases,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=coordination.k8s.io,resources=leases,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=coordination.k8s.io,resources=leases,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=coordination.k8s.io,resources=leases,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core.observatorium.io,resources=observatoria,verbs=create;delete;get;list;patch;update;watch
//+kubebuilder:rbac:groups=discovery.open-cluster-management.io,resources=discoveryconfigs;discoveredclusters,verbs=list;watch
//+kubebuilder:rbac:groups=extensions.hive.openshift.io,resources=agentclusterinstalls,verbs=list;watch
//+kubebuilder:rbac:groups=hive.openshift.io,resources=clusterclaims;clusterdeployments;clusterpools;clusterimagesets;clusterprovisions;clusterdeprovisions;machinepools,verbs=list;watch
//+kubebuilder:rbac:groups=hive.openshift.io;multicluster.openshift.io,resources=clusterimagesets;multiclusterengines,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=hypershift.openshift.io,resources=hostedclusters;nodepools,verbs=list;watch
//+kubebuilder:rbac:groups=imageregistry.open-cluster-management.io,resources=managedclusterimageregistries,verbs=get;list;watch
//+kubebuilder:rbac:groups=imageregistry.open-cluster-management.io,resources=managedclusterimageregistries;managedclusterimageregistries/status,verbs=get;list;watch
//+kubebuilder:rbac:groups=integreatly.org,resources=grafanas;grafanas/status;grafanas/finalizers;grafanadashboards;grafanadashboards/status;grafanadatasources;grafanadatasources/status,verbs=get;list;create;update;delete;deletecollection;watch
//+kubebuilder:rbac:groups=internal.open-cluster-management.io,resources=managedclusterinfos,verbs=list;watch
//+kubebuilder:rbac:groups=metal3.io,resources=baremetalhosts;provisionings,verbs=list;watch
//+kubebuilder:rbac:groups=migration.k8s.io,resources=storageversionmigrations,verbs=create;delete;get;list;update;watch
//+kubebuilder:rbac:groups=monitor.open-cluster-management.io,resources=*,verbs=create;delete;get;list;patch;update;watch
//+kubebuilder:rbac:groups=monitoring.coreos.com,resources=prometheusrules,verbs=create;delete;get;list;patch;update;watch
//+kubebuilder:rbac:groups=monitoring.coreos.com,resources=servicemonitors,verbs=create;delete;get;list
//+kubebuilder:rbac:groups=monitoring.coreos.com,resources=servicemonitors,verbs=create;delete;get;list;patch;update;watch
//+kubebuilder:rbac:groups=monitoring.coreos.com,resources=servicemonitors,verbs=get;create
//+kubebuilder:rbac:groups=multicluster.openshift.io,resources=multiclusterengines,verbs=get;list
//+kubebuilder:rbac:groups=multicluster.openshift.io,resources=multiclusterengines,verbs=list;watch
//+kubebuilder:rbac:groups=networking.k8s.io,resources=ingresses,verbs=get;list;create;update;delete;deletecollection;watch
//+kubebuilder:rbac:groups=oauth.openshift.io,resources=oauthclients,verbs=create;delete;get;list;patch;update;watch
//+kubebuilder:rbac:groups=observability.open-cluster-management.io,resources=*;multiclusterobservabilities;endpointmonitorings,verbs=create;delete;get;list;patch;update;watch
//+kubebuilder:rbac:groups=operator.open-cluster-management.io,resources=multiclusterglobalhubs,verbs=get;list
//+kubebuilder:rbac:groups=operator.open-cluster-management.io,resources=multiclusterhubs,verbs=get
//+kubebuilder:rbac:groups=operator.open-cluster-management.io,resources=multiclusterhubs,verbs=get;list
//+kubebuilder:rbac:groups=operator.open-cluster-management.io,resources=multiclusterhubs,verbs=get;list;watch
//+kubebuilder:rbac:groups=operator.open-cluster-management.io,resources=multiclusterhubs,verbs=watch;get;list
//+kubebuilder:rbac:groups=operator.openshift.io,resources=ingresscontrollers,verbs=get;list;watch
//+kubebuilder:rbac:groups=operators.coreos.com,resources=subscriptions,verbs=get;list;watch
//+kubebuilder:rbac:groups=policy.open-cluster-management.io,resources=*,verbs=create;delete;get;list;patch;update;watch
//+kubebuilder:rbac:groups=policy.open-cluster-management.io,resources=placementbindings;policies;policyautomations;policysets,verbs=list;watch
//+kubebuilder:rbac:groups=policy.open-cluster-management.io,resources=policies,verbs=create;delete;get;list;patch;update;watch
//+kubebuilder:rbac:groups=policy.open-cluster-management.io,resources=policies,verbs=list;get;watch
//+kubebuilder:rbac:groups=policy.open-cluster-management.io,resources=policies/finalizers,verbs=update
//+kubebuilder:rbac:groups=policy.open-cluster-management.io,resources=policies/status,verbs=get;patch;update
//+kubebuilder:rbac:groups=policy.open-cluster-management.io,resources=policysets,verbs=create;delete;get;list;patch;update;watch
//+kubebuilder:rbac:groups=policy.open-cluster-management.io,resources=policysets/finalizers,verbs=update
//+kubebuilder:rbac:groups=policy.open-cluster-management.io,resources=policysets/status,verbs=get;patch;update
//+kubebuilder:rbac:groups=proxy.open-cluster-management.io,resources=clusterstatuses/aggregator,verbs=create
//+kubebuilder:rbac:groups=rbac.authorization.k8s.io,resources=*,verbs=*
//+kubebuilder:rbac:groups=rbac.authorization.k8s.io,resources=clusterrolebindings;clusterroles;rolebindings;roles,verbs=create;delete;get;list
//+kubebuilder:rbac:groups=rbac.authorization.k8s.io,resources=clusterroles,verbs=create
//+kubebuilder:rbac:groups=rbac.authorization.k8s.io,resources=clusterroles,verbs=create;get;delete
//+kubebuilder:rbac:groups=rbac.authorization.k8s.io,resources=clusterroles,verbs=delete;get;patch;update
//+kubebuilder:rbac:groups=rbac.authorization.k8s.io,resources=rolebindings,verbs=create
//+kubebuilder:rbac:groups=rbac.authorization.k8s.io,resources=rolebindings,verbs=create;get;list;watch;update;patch;delete
//+kubebuilder:rbac:groups=rbac.authorization.k8s.io,resources=rolebindings,verbs=delete;get;patch;update
//+kubebuilder:rbac:groups=rbac.authorization.k8s.io,resources=rolebindings;roles;clusterrolebindings;clusterroles,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=rbac.authorization.k8s.io,resources=roles,verbs=create;get;list;watch;update;patch;delete;bind;escalate
//+kubebuilder:rbac:groups=rbac.authorization.k8s.io,resources=roles;rolebindings,verbs=create;delete;get;list;patch;update;watch
//+kubebuilder:rbac:groups=rbac.authorization.k8s.io,resources=roles;rolebindings,verbs=get;list;watch;create;update;delete
//+kubebuilder:rbac:groups=rbac.open-cluster-management.io,resources=clusterpermissions,verbs=create;delete;get;list;patch;update;watch
//+kubebuilder:rbac:groups=rbac.open-cluster-management.io,resources=clusterpermissions/finalizers,verbs=update
//+kubebuilder:rbac:groups=rbac.open-cluster-management.io,resources=clusterpermissions/status,verbs=get;patch;update
//+kubebuilder:rbac:groups=route.openshift.io,resources=routes,verbs=create
//+kubebuilder:rbac:groups=route.openshift.io,resources=routes,verbs=delete;get;list;update;watch
//+kubebuilder:rbac:groups=route.openshift.io,resources=routes;routes/custom-host;routes/status,verbs=get;list;create;update;delete;deletecollection;watch;create
//+kubebuilder:rbac:groups=search.open-cluster-management.io,resources=searches,verbs=get;list;patch;update;watch
//+kubebuilder:rbac:groups=search.open-cluster-management.io,resources=searches/finalizers,verbs=update
//+kubebuilder:rbac:groups=search.open-cluster-management.io,resources=searches/status,verbs=get;patch;update
//+kubebuilder:rbac:groups=storage.k8s.io,resources=storageclasses,verbs=list;watch
//+kubebuilder:rbac:groups=storage.k8s.io,resources=storageclasses,verbs=watch;get;list
//+kubebuilder:rbac:groups=submariner.io,resources=brokers,verbs=create;get;update;delete
//+kubebuilder:rbac:groups=submarineraddon.open-cluster-management.io,resources=submarinerconfigs,verbs=get;list;watch;update;patch
//+kubebuilder:rbac:groups=submarineraddon.open-cluster-management.io,resources=submarinerconfigs,verbs=list;watch
//+kubebuilder:rbac:groups=submarineraddon.open-cluster-management.io,resources=submarinerconfigs/status,verbs=update;patch
//+kubebuilder:rbac:groups=tower.ansible.com,resources=ansiblejobs,verbs=create;delete;get;list;patch;update;watch
//+kubebuilder:rbac:groups=wgpolicyk8s.io,resources=policyreports,verbs=create;delete;get;list;patch;update;watch
//+kubebuilder:rbac:groups=wgpolicyk8s.io,resources=policyreports,verbs=list;watch
//+kubebuilder:rbac:groups=work.open-cluster-management.io,resources=manifestworks,verbs=*
//+kubebuilder:rbac:groups=work.open-cluster-management.io,resources=manifestworks,verbs=create;delete;get;list
//+kubebuilder:rbac:groups=work.open-cluster-management.io,resources=manifestworks,verbs=create;delete;get;list;patch;update;watch
//+kubebuilder:rbac:groups=work.open-cluster-management.io,resources=manifestworks,verbs=create;delete;get;list;patch;update;watch
//+kubebuilder:rbac:groups=work.open-cluster-management.io,resources=manifestworks,verbs=create;get;list;watch;update;patch;delete
//+kubebuilder:rbac:groups=work.open-cluster-management.io,resources=manifestworks,verbs=create;update;get;list;watch;delete;deletecollection;patch
//+kubebuilder:rbac:groups=work.open-cluster-management.io,resources=manifestworks;manifestworks/finalizers,verbs=create;delete;get;list;patch;update;watch
