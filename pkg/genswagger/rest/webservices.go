package rest

import (
	"fmt"
	"net/http"
	"reflect"

	restful "github.com/emicklei/go-restful"
	cniv1 "github.com/k8snetworkplumbingwg/network-attachment-definition-client/pkg/apis/k8s.cni.cncf.io/v1"
	"github.com/rancher/wrangler/pkg/slice"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	virtv1 "kubevirt.io/client-go/api/v1"
	cdiv1 "kubevirt.io/containerized-data-importer/pkg/apis/core/v1beta1"
	mime "kubevirt.io/kubevirt/pkg/rest"

	networkv1beta1 "github.com/harvester/harvester-network-controller/pkg/apis/network.harvesterhci.io/v1beta1"
	"github.com/harvester/harvester/pkg/apis/harvesterhci.io/v1beta1"
)

var defaultActions = []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete}

func AggregatedWebServices() []*restful.WebService {
	harvesterAPI := NewGroupWebService(v1beta1.SchemeGroupVersion)
	harvesterv1beta1API := NewGroupVersionWebService(v1beta1.SchemeGroupVersion)
	backupGVR := schema.GroupVersionResource{Group: v1beta1.SchemeGroupVersion.Group, Version: v1beta1.SchemeGroupVersion.Group, Resource: "virtualmachinebackups"}
	imageGVR := schema.GroupVersionResource{Group: v1beta1.SchemeGroupVersion.Group, Version: v1beta1.SchemeGroupVersion.Group, Resource: "virtualmachineimages"}
	keypairGVR := schema.GroupVersionResource{Group: v1beta1.SchemeGroupVersion.Group, Version: v1beta1.SchemeGroupVersion.Group, Resource: "keypairs"}
	restoreGVR := schema.GroupVersionResource{Group: v1beta1.SchemeGroupVersion.Group, Version: v1beta1.SchemeGroupVersion.Group, Resource: "virtualmachinerestores"}
	supportBundleGVR := schema.GroupVersionResource{Group: v1beta1.SchemeGroupVersion.Group, Version: v1beta1.SchemeGroupVersion.Group, Resource: "supportbundles"}
	templateGVR := schema.GroupVersionResource{Group: v1beta1.SchemeGroupVersion.Group, Version: v1beta1.SchemeGroupVersion.Group, Resource: "virtualmachinetemplates"}
	templateVersionGVR := schema.GroupVersionResource{Group: v1beta1.SchemeGroupVersion.Group, Version: v1beta1.SchemeGroupVersion.Group, Resource: "virtualmachinetemplateversions"}
	upgradeGVR := schema.GroupVersionResource{Group: v1beta1.SchemeGroupVersion.Group, Version: v1beta1.SchemeGroupVersion.Group, Resource: "upgrades"}

	AddGenericNamespacedResourceRoutes(harvesterv1beta1API, backupGVR, &v1beta1.VirtualMachineBackup{}, "VirtualMachineBackup", &v1beta1.VirtualMachineBackupList{})
	AddGenericNamespacedResourceRoutes(harvesterv1beta1API, restoreGVR, &v1beta1.VirtualMachineRestore{}, "VirtualMachineRestore", &v1beta1.VirtualMachineRestoreList{})
	AddGenericNamespacedResourceRoutes(harvesterv1beta1API, imageGVR, &v1beta1.VirtualMachineImage{}, "VirtualMachineImage", &v1beta1.VirtualMachineImageList{})
	AddGenericNamespacedResourceRoutes(harvesterv1beta1API, templateGVR, &v1beta1.VirtualMachineTemplate{}, "VirtualMachineTemplate", &v1beta1.VirtualMachineTemplateList{})
	AddGenericNamespacedResourceRoutes(harvesterv1beta1API, templateVersionGVR, &v1beta1.VirtualMachineTemplateVersion{}, "VirtualMachineTemplateVersion", &v1beta1.VirtualMachineTemplateVersionList{})
	AddGenericNamespacedResourceRoutes(harvesterv1beta1API, keypairGVR, &v1beta1.KeyPair{}, "KeyPair", &v1beta1.KeyPairList{})
	AddGenericNamespacedResourceRoutes(harvesterv1beta1API, upgradeGVR, &v1beta1.Upgrade{}, "Upgrade", &v1beta1.UpgradeList{})
	AddGenericNamespacedResourceRoutes(harvesterv1beta1API, supportBundleGVR, &v1beta1.SupportBundle{}, "SupportBundle", &v1beta1.SupportBundleList{})

	harvesterNetworkv1beta1API := NewGroupVersionWebService(networkv1beta1.SchemeGroupVersion)

	clusterNetworkGVR := schema.GroupVersionResource{Group: networkv1beta1.SchemeGroupVersion.Group, Version: networkv1beta1.SchemeGroupVersion.Version, Resource: "clusternetworks"}
	nodeNetworkGVR := schema.GroupVersionResource{Group: networkv1beta1.SchemeGroupVersion.Group, Version: networkv1beta1.SchemeGroupVersion.Version, Resource: "nodenetworks"}

	AddGenericNamespacedResourceRoutes(harvesterNetworkv1beta1API, clusterNetworkGVR, &networkv1beta1.ClusterNetwork{}, "ClusterNetwork", &networkv1beta1.ClusterNetworkList{})
	AddGenericNamespacedResourceRoutes(harvesterNetworkv1beta1API, nodeNetworkGVR, &networkv1beta1.NodeNetwork{}, "NodeNetwork", &networkv1beta1.NodeNetworkList{})

	// CDI
	cdiAPI := NewGroupWebService(cdiv1.SchemeGroupVersion)
	cdiv1beta1API := NewGroupVersionWebService(cdiv1.SchemeGroupVersion)
	dvGVR := schema.GroupVersionResource{Group: cdiv1.SchemeGroupVersion.Group, Version: cdiv1.SchemeGroupVersion.Version, Resource: "datavolumes"}

	AddGenericNamespacedResourceRoutes(cdiv1beta1API, dvGVR, &cdiv1.DataVolume{}, "DataVolume", &cdiv1.DataVolumeList{})

	// kubevirt
	virtAPI := NewGroupWebService(virtv1.SchemeGroupVersion)
	virtv1API := NewGroupVersionWebService(virtv1.SchemeGroupVersion)
	vmGVR := schema.GroupVersionResource{Group: virtv1.GroupVersion.Group, Version: virtv1.GroupVersion.Version, Resource: "virtualmachines"}
	vmiGVR := schema.GroupVersionResource{Group: virtv1.GroupVersion.Group, Version: virtv1.GroupVersion.Version, Resource: "virtualmachineinstances"}
	migrationGVR := schema.GroupVersionResource{Group: virtv1.GroupVersion.Group, Version: virtv1.GroupVersion.Version, Resource: "virtualmachineinstancemigrations"}

	AddGenericNamespacedResourceRoutes(virtv1API, vmiGVR, &virtv1.VirtualMachineInstance{}, "VirtualMachineInstance", &virtv1.VirtualMachineInstanceList{}, http.MethodGet)
	AddGenericNamespacedResourceRoutes(virtv1API, vmGVR, &virtv1.VirtualMachine{}, "VirtualMachine", &virtv1.VirtualMachineList{})
	AddGenericNamespacedResourceRoutes(virtv1API, migrationGVR, &virtv1.VirtualMachineInstanceMigration{}, "VirtualMachineInstanceMigration", &virtv1.VirtualMachineInstanceMigrationList{})

	// multus
	cniv1API := NewGroupVersionWebService(cniv1.SchemeGroupVersion)
	nadGVR := schema.GroupVersionResource{Group: cniv1.SchemeGroupVersion.Group, Version: cniv1.SchemeGroupVersion.Version, Resource: "network-attachment-definitions"}

	//network-attachment-definitions
	AddGenericNamespacedResourceRoutes(cniv1API, nadGVR, &cniv1.NetworkAttachmentDefinition{}, "NetworkAttachmentDefinition", &cniv1.NetworkAttachmentDefinitionList{})

	return []*restful.WebService{harvesterAPI, harvesterv1beta1API, harvesterNetworkv1beta1API, cdiAPI, cdiv1beta1API, virtAPI, virtv1API, cniv1API}
}

func NewGroupVersionWebService(gv schema.GroupVersion) *restful.WebService {
	ws := new(restful.WebService)
	ws.Doc("The Harvester API.")
	ws.Path(GroupVersionBasePath(gv))

	return ws
}

func AddGenericNonNamespacedResourceRoutes(ws *restful.WebService, gvr schema.GroupVersionResource, objPointer runtime.Object, objKind string, objListPointer runtime.Object, actions ...string) {
	AddGenericResourceRoutes(ws, gvr, objPointer, objKind, objListPointer, false, actions...)
}

func AddGenericNamespacedResourceRoutes(ws *restful.WebService, gvr schema.GroupVersionResource, objPointer runtime.Object, objKind string, objListPointer runtime.Object, actions ...string) {
	AddGenericResourceRoutes(ws, gvr, objPointer, objKind, objListPointer, true, actions...)
}

func AddGenericResourceRoutes(ws *restful.WebService, gvr schema.GroupVersionResource, objPointer runtime.Object, objKind string, objListPointer runtime.Object, namespaced bool, actions ...string) {
	if actions == nil {
		actions = defaultActions
	}

	objExample := reflect.ValueOf(objPointer).Elem().Interface()
	listExample := reflect.ValueOf(objListPointer).Elem().Interface()

	if slice.ContainsString(actions, http.MethodGet) {
		ws.Route(addGetParams(
			ws.GET(ResourcePath(gvr, namespaced)).
				Produces(mime.MIME_JSON, mime.MIME_YAML, mime.MIME_JSON_STREAM).
				Operation("readNamespaced"+objKind).
				To(Noop).Writes(objExample).
				Doc("Get a "+objKind+" object.").
				Metadata("kind", objKind).
				Returns(http.StatusOK, "OK", objExample).
				Returns(http.StatusUnauthorized, "Unauthorized", ""), ws,
		))

		ws.Route(addGetAllNamespacesListParams(
			ws.GET(gvr.Resource).
				Produces(mime.MIME_JSON, mime.MIME_YAML, mime.MIME_JSON_STREAM).
				Operation("list"+objKind+"ForAllNamespaces").
				To(Noop).Writes(listExample).
				Doc("Get a list of all "+objKind+" objects.").
				Metadata("kind", objKind).
				Returns(http.StatusOK, "OK", listExample).
				Returns(http.StatusUnauthorized, "Unauthorized", ""), ws,
		))

		ws.Route(addGetNamespacedListParams(
			ws.GET(ResourceBasePath(gvr, namespaced)).
				Produces(mime.MIME_JSON, mime.MIME_YAML, mime.MIME_JSON_STREAM).
				Operation("listNamespaced"+objKind).
				Writes(listExample).
				To(Noop).
				Doc("Get a list of "+objKind+" objects in a namespace.").
				Metadata("kind", objKind).
				Returns(http.StatusOK, "OK", listExample).
				Returns(http.StatusUnauthorized, "Unauthorized", ""), ws,
		))
	}

	if slice.ContainsString(actions, http.MethodPost) {
		ws.Route(addPostParams(
			ws.POST(ResourceBasePath(gvr, namespaced)).
				Produces(mime.MIME_JSON, mime.MIME_YAML).
				Consumes(mime.MIME_JSON, mime.MIME_YAML).
				Operation("createNamespaced"+objKind).
				To(Noop).Reads(objExample).Writes(objExample).
				Doc("Create a "+objKind+" object.").
				Metadata("kind", objKind).
				Returns(http.StatusOK, "OK", objExample).
				Returns(http.StatusCreated, "Created", objExample).
				Returns(http.StatusAccepted, "Accepted", objExample).
				Returns(http.StatusUnauthorized, "Unauthorized", ""), ws,
		))
	}

	if slice.ContainsString(actions, http.MethodPut) {
		ws.Route(addPutParams(
			ws.PUT(ResourcePath(gvr, namespaced)).
				Produces(mime.MIME_JSON, mime.MIME_YAML).
				Consumes(mime.MIME_JSON, mime.MIME_YAML).
				Operation("replaceNamespaced"+objKind).
				To(Noop).Reads(objExample).Writes(objExample).
				Doc("Update a "+objKind+" object.").
				Metadata("kind", objKind).
				Returns(http.StatusOK, "OK", objExample).
				Returns(http.StatusCreated, "Create", objExample).
				Returns(http.StatusUnauthorized, "Unauthorized", ""), ws,
		))
	}

	if slice.ContainsString(actions, http.MethodPatch) {
		ws.Route(addPatchParams(
			ws.PATCH(ResourcePath(gvr, namespaced)).
				Consumes(mime.MIME_JSON_PATCH, mime.MIME_MERGE_PATCH).
				Produces(mime.MIME_JSON).
				Operation("patchNamespaced"+objKind).
				To(Noop).
				Writes(objExample).Reads(metav1.Patch{}).
				Doc("Patch a "+objKind+" object.").
				Metadata("kind", objKind).
				Returns(http.StatusOK, "OK", objExample).
				Returns(http.StatusUnauthorized, "Unauthorized", ""), ws,
		))
	}

	if slice.ContainsString(actions, http.MethodDelete) {
		ws.Route(addDeleteParams(
			ws.DELETE(ResourcePath(gvr, namespaced)).
				Produces(mime.MIME_JSON, mime.MIME_YAML).
				Consumes(mime.MIME_JSON, mime.MIME_YAML).
				Operation("deleteNamespaced"+objKind).
				To(Noop).
				Reads(metav1.DeleteOptions{}).Writes(metav1.Status{}).
				Doc("Delete a "+objKind+" object.").
				Metadata("kind", objKind).
				Returns(http.StatusOK, "OK", metav1.Status{}).
				Returns(http.StatusUnauthorized, "Unauthorized", ""), ws,
		))
	}

}

func NewGroupWebService(gv schema.GroupVersion) *restful.WebService {
	ws := new(restful.WebService)
	ws.Path(GroupBasePath(gv))
	return ws
}

func addCollectionParams(builder *restful.RouteBuilder, ws *restful.WebService) *restful.RouteBuilder {
	return builder.Param(continueParam(ws)).
		Param(fieldSelectorParam(ws)).
		Param(includeUninitializedParam(ws)).
		Param(labelSelectorParam(ws)).
		Param(limitParam(ws)).
		Param(resourceVersionParam(ws)).
		Param(timeoutSecondsParam(ws)).
		Param(watchParam(ws))
}

func addGetAllNamespacesListParams(builder *restful.RouteBuilder, ws *restful.WebService) *restful.RouteBuilder {
	return addCollectionParams(builder, ws)
}

func addGetParams(builder *restful.RouteBuilder, ws *restful.WebService) *restful.RouteBuilder {
	return builder.Param(NameParam(ws)).
		Param(NamespaceParam(ws)).
		Param(exactParam(ws)).
		Param(exportParam(ws))
}

func addGetNamespacedListParams(builder *restful.RouteBuilder, ws *restful.WebService) *restful.RouteBuilder {
	return addCollectionParams(builder.Param(NamespaceParam(ws)), ws)
}

func addPostParams(builder *restful.RouteBuilder, ws *restful.WebService) *restful.RouteBuilder {
	return builder.Param(NamespaceParam(ws))
}

func addPutParams(builder *restful.RouteBuilder, ws *restful.WebService) *restful.RouteBuilder {
	return builder.Param(NamespaceParam(ws)).Param(NameParam(ws))
}

func addDeleteParams(builder *restful.RouteBuilder, ws *restful.WebService) *restful.RouteBuilder {
	return builder.Param(NamespaceParam(ws)).Param(NameParam(ws)).
		Param(gracePeriodSecondsParam(ws)).
		Param(orphanDependentsParam(ws)).
		Param(propagationPolicyParam(ws))
}

func addPatchParams(builder *restful.RouteBuilder, ws *restful.WebService) *restful.RouteBuilder {
	return builder.Param(NamespaceParam(ws)).Param(NameParam(ws))
}

func NameParam(ws *restful.WebService) *restful.Parameter {
	return ws.PathParameter("name", "Name of the resource").Required(true)
}

func NamespaceParam(ws *restful.WebService) *restful.Parameter {
	return ws.PathParameter("namespace", "Object name and auth scope, such as for teams and projects").Required(true)
}

func labelSelectorParam(ws *restful.WebService) *restful.Parameter {
	return ws.QueryParameter("labelSelector", "A selector to restrict the list of returned objects by their labels. Defaults to everything")
}

func fieldSelectorParam(ws *restful.WebService) *restful.Parameter {
	return ws.QueryParameter("fieldSelector", "A selector to restrict the list of returned objects by their fields. Defaults to everything.")
}

func resourceVersionParam(ws *restful.WebService) *restful.Parameter {
	return ws.QueryParameter("resourceVersion", "When specified with a watch call, shows changes that occur after that particular version of a resource. Defaults to changes from the beginning of history.")
}

func timeoutSecondsParam(ws *restful.WebService) *restful.Parameter {
	return ws.QueryParameter("timeoutSeconds", "TimeoutSeconds for the list/watch call.").DataType("integer")
}

func includeUninitializedParam(ws *restful.WebService) *restful.Parameter {
	return ws.QueryParameter("includeUninitialized", "If true, partially initialized resources are included in the response.").DataType("boolean")
}

func watchParam(ws *restful.WebService) *restful.Parameter {
	return ws.QueryParameter("watch", "Watch for changes to the described resources and return them as a stream of add, update, and remove notifications. Specify resourceVersion.").DataType("boolean")
}

func limitParam(ws *restful.WebService) *restful.Parameter {
	return ws.QueryParameter("limit", "limit is a maximum number of responses to return for a list call. If more items exist, the server will set the `continue` field on the list metadata to a value that can be used with the same initial query to retrieve the next set of results. Setting a limit may return fewer than the requested amount of items (up to zero items) in the event all requested objects are filtered out and clients should only use the presence of the continue field to determine whether more results are available. Servers may choose not to support the limit argument and will return all of the available results. If limit is specified and the continue field is empty, clients may assume that no more results are available. This field is not supported if watch is true.\n\nThe server guarantees that the objects returned when using continue will be identical to issuing a single list call without a limit - that is, no objects created, modified, or deleted after the first request is issued will be included in any subsequent continued requests. This is sometimes referred to as a consistent snapshot, and ensures that a client that is using limit to receive smaller chunks of a very large result can ensure they see all possible objects. If objects are updated during a chunked list the version of the object that was present at the time the first list result was calculated is returned.").DataType("integer")
}

func continueParam(ws *restful.WebService) *restful.Parameter {
	return ws.QueryParameter("continue", "The continue option should be set when retrieving more results from the server. Since this value is server defined, clients may only use the continue value from a previous query result with identical query parameters (except for the value of continue) and the server may reject a continue value it does not recognize. If the specified continue value is no longer valid whether due to expiration (generally five to fifteen minutes) or a configuration change on the server the server will respond with a 410 ResourceExpired error indicating the client must restart their list without the continue field. This field is not supported when watch is true. Clients may start a watch from the last resourceVersion value returned by the server and not miss any modifications.")
}

func exactParam(ws *restful.WebService) *restful.Parameter {
	return ws.QueryParameter("exact", "Should the export be exact. Exact export maintains cluster-specific fields like 'Namespace'.").DataType("boolean")
}

func exportParam(ws *restful.WebService) *restful.Parameter {
	return ws.QueryParameter("export", "Should this value be exported. Export strips fields that a user can not specify.").DataType("boolean")
}

func gracePeriodSecondsParam(ws *restful.WebService) *restful.Parameter {
	return ws.QueryParameter("gracePeriodSeconds", "The duration in seconds before the object should be deleted. Value must be non-negative integer. The value zero indicates delete immediately. If this value is nil, the default grace period for the specified type will be used. Defaults to a per object value if not specified. zero means delete immediately.").DataType("integer")
}

func orphanDependentsParam(ws *restful.WebService) *restful.Parameter {
	return ws.QueryParameter("orphanDependents", "Deprecated: please use the PropagationPolicy, this field will be deprecated in 1.7. Should the dependent objects be orphaned. If true/false, the \"orphan\" finalizer will be added to/removed from the object's finalizers list. Either this field or PropagationPolicy may be set, but not both.").DataType("boolean")
}

func propagationPolicyParam(ws *restful.WebService) *restful.Parameter {
	return ws.QueryParameter("propagationPolicy", "Whether and how garbage collection will be performed. Either this field or OrphanDependents may be set, but not both. The default policy is decided by the existing finalizer set in the metadata.finalizers and the resource-specific default policy. Acceptable values are: 'Orphan' - orphan the dependents; 'Background' - allow the garbage collector to delete the dependents in the background; 'Foreground' - a cascading policy that deletes all dependents in the foreground.")
}

func GroupBasePath(gvr schema.GroupVersion) string {
	return fmt.Sprintf("/apis/%s", gvr.Group)
}

func GroupVersionBasePath(gvr schema.GroupVersion) string {
	return fmt.Sprintf("/apis/%s/%s", gvr.Group, gvr.Version)
}

func ResourceBasePath(gvr schema.GroupVersionResource, namespaced bool) string {
	if namespaced {
		return fmt.Sprintf("/namespaces/{namespace:[a-z0-9][a-z0-9\\-]*}/%s", gvr.Resource)
	}
	return fmt.Sprintf("/%s", gvr.Resource)
}

func ResourcePath(gvr schema.GroupVersionResource, namespaced bool) string {
	if namespaced {
		return fmt.Sprintf("/namespaces/{namespace:[a-z0-9][a-z0-9\\-]*}/%s/{name:[a-z0-9][a-z0-9\\-]*}", gvr.Resource)
	}
	return fmt.Sprintf("/%s/{name:[a-z0-9][a-z0-9\\-]*}", gvr.Resource)
}

func Noop(request *restful.Request, response *restful.Response) {}