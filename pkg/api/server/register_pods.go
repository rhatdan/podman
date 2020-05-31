package server

import (
	"net/http"

	"github.com/containers/libpod/pkg/api/handlers/libpod"
	"github.com/gorilla/mux"
)

func (s *APIServer) registerPodsHandlers(r *mux.Router) error {
	// swagger:operation GET /libpod/pods/json pods ListPods
	// ---
	// summary: List pods
	// produces:
	// - application/json
	// parameters:
	// - in: query
	//   name: filters
	//   type: string
	//   description: needs description and plumbing for filters
	// responses:
	//   200:
	//     $ref: "#/responses/ListPodsResponse"
	//   400:
	//     $ref: "#/responses/BadParamError"
	//   500:
	//     $ref: "#/responses/InternalError"
	r.Handle(VersionedPath("/libpod/pods/json"), APIHandler(s.Context, libpod.Pods)).Methods(http.MethodGet)
	r.Handle(VersionedPath("/libpod/pods/create"), APIHandler(s.Context, libpod.PodCreate)).Methods(http.MethodPost)
	// swagger:operation POST /libpod/pods/prune pods PrunePods
	// ---
	// summary: Prune unused pods
	// produces:
	// - application/json
	// responses:
	//   200:
	//     description: tbd
	//     schema:
	//       type: object
	//       additionalProperties:
	//         type: string
	//   400:
	//     $ref: "#/responses/BadParamError"
	//   409:
	//     description: pod already exists
	//   500:
	//     $ref: "#/responses/InternalError"
	r.Handle(VersionedPath("/libpod/pods/prune"), APIHandler(s.Context, libpod.PodPrune)).Methods(http.MethodPost)
	// swagger:operation DELETE /libpod/pods/{name} pods removePod
	// ---
	// summary: Remove pod
	// produces:
	// - application/json
	// parameters:
	//  - in: path
	//    name: name
	//    type: string
	//    required: true
	//    description: the name or ID of the pod
	//  - in: query
	//    name: force
	//    type: boolean
	//    description : force removal of a running pod by first stopping all containers, then removing all containers in the pod
	// responses:
	//   204:
	//     description: no error
	//   400:
	//     $ref: "#/responses/BadParamError"
	//   404:
	//     $ref: "#/responses/NoSuchPod"
	//   500:
	//     $ref: "#/responses/InternalError"
	r.Handle(VersionedPath("/libpod/pods/{name}"), APIHandler(s.Context, libpod.PodDelete)).Methods(http.MethodDelete)
	// swagger:operation GET /libpod/pods/{name}/json pods inspectPod
	// ---
	// summary: Inspect pod
	// produces:
	// - application/json
	// parameters:
	//  - in: path
	//    name: name
	//    type: string
	//    required: true
	//    description: the name or ID of the pod
	// responses:
	//   200:
	//     $ref: "#/responses/InspectPodResponse"
	//   404:
	//      $ref: "#/responses/NoSuchPod"
	//   500:
	//      $ref: "#/responses/InternalError"
	r.Handle(VersionedPath("/libpod/pods/{name}/json"), APIHandler(s.Context, libpod.PodInspect)).Methods(http.MethodGet)
	// swagger:operation GET /libpod/pods/{name}/exists pods podExists
	// ---
	// summary: Pod exists
	// description: Check if a pod exists by name or ID
	// produces:
	// - application/json
	// parameters:
	//  - in: path
	//    name: name
	//    type: string
	//    required: true
	//    description: the name or ID of the pod
	// responses:
	//   204:
	//     description: pod exists
	//   404:
	//     $ref: "#/responses/NoSuchPod"
	//   500:
	//     $ref: "#/responses/InternalError"
	r.Handle(VersionedPath("/libpod/pods/{name}/exists"), APIHandler(s.Context, libpod.PodExists)).Methods(http.MethodGet)
	// swagger:operation POST /libpod/pods/{name}/kill pods killPod
	// ---
	// summary: Kill a pod
	// produces:
	// - application/json
	// parameters:
	//  - in: path
	//    name: name
	//    type: string
	//    required: true
	//    description: the name or ID of the pod
	//  - in: query
	//    name: signal
	//    type: string
	//    description: signal to be sent to pod
	//    default: SIGKILL
	// responses:
	//   204:
	//     description: no error
	//   400:
	//     $ref: "#/responses/BadParamError"
	//   404:
	//     $ref: "#/responses/NoSuchPod"
	//   409:
	//     $ref: "#/responses/ConflictError"
	//   500:
	//     $ref: "#/responses/InternalError"
	r.Handle(VersionedPath("/libpod/pods/{name}/kill"), APIHandler(s.Context, libpod.PodKill)).Methods(http.MethodPost)
	// swagger:operation POST /libpod/pods/{name}/pause pods pausePod
	// ---
	// summary: Pause a pod
	// description: Pause a pod
	// produces:
	// - application/json
	// parameters:
	//  - in: path
	//    name: name
	//    type: string
	//    required: true
	//    description: the name or ID of the pod
	// responses:
	//   204:
	//     description: no error
	//   404:
	//     $ref: "#/responses/NoSuchPod"
	//   500:
	//     $ref: "#/responses/InternalError"
	r.Handle(VersionedPath("/libpod/pods/{name}/pause"), APIHandler(s.Context, libpod.PodPause)).Methods(http.MethodPost)
	// swagger:operation POST /libpod/pods/{name}/restart pods restartPod
	// ---
	// summary: Restart a pod
	// produces:
	// - application/json
	// parameters:
	//  - in: path
	//    name: name
	//    type: string
	//    required: true
	//    description: the name or ID of the pod
	// responses:
	//   204:
	//     description: no error
	//   404:
	//     $ref: "#/responses/NoSuchPod"
	//   500:
	//     $ref: "#/responses/InternalError"
	r.Handle(VersionedPath("/libpod/pods/{name}/restart"), APIHandler(s.Context, libpod.PodRestart)).Methods(http.MethodPost)
	// swagger:operation POST /libpod/pods/{name}/start pods startPod
	// ---
	// summary: Start a pod
	// produces:
	// - application/json
	// parameters:
	//  - in: path
	//    name: name
	//    type: string
	//    required: true
	//    description: the name or ID of the pod
	// responses:
	//   204:
	//     description: no error
	//   304:
	//     $ref: "#/responses/PodAlreadyStartedError"
	//   404:
	//     $ref: "#/responses/NoSuchPod"
	//   500:
	//     $ref: "#/responses/InternalError"
	r.Handle(VersionedPath("/libpod/pods/{name}/start"), APIHandler(s.Context, libpod.PodStart)).Methods(http.MethodPost)
	// swagger:operation POST /libpod/pods/{name}/stop pods stopPod
	// ---
	// summary: Stop a pod
	// produces:
	// - application/json
	// parameters:
	//  - in: path
	//    name: name
	//    type: string
	//    required: true
	//    description: the name or ID of the pod
	//  - in: query
	//    name: t
	//    type: integer
	//    description: timeout
	// responses:
	//   204:
	//     description: no error
	//   304:
	//     $ref: "#/responses/PodAlreadyStoppedError"
	//   400:
	//     $ref: "#/responses/BadParamError"
	//   404:
	//     $ref: "#/responses/NoSuchPod"
	//   500:
	//     $ref: "#/responses/InternalError"
	r.Handle(VersionedPath("/libpod/pods/{name}/stop"), APIHandler(s.Context, libpod.PodStop)).Methods(http.MethodPost)
	// swagger:operation POST /libpod/pods/{name}/unpause pods unpausePod
	// ---
	// summary: Unpause a pod
	// produces:
	// - application/json
	// parameters:
	//  - in: path
	//    name: name
	//    type: string
	//    required: true
	//    description: the name or ID of the pod
	// responses:
	//   204:
	//     description: no error
	//   404:
	//     $ref: "#/responses/NoSuchPod"
	//   500:
	//     $ref: "#/responses/InternalError"
	r.Handle(VersionedPath("/libpod/pods/{name}/unpause"), APIHandler(s.Context, libpod.PodUnpause)).Methods(http.MethodPost)
	return nil
}
