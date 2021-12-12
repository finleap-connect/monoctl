package output

import (
	"context"
	"fmt"
	"github.com/finleap-connect/monoskope/pkg/api/domain/eventdata"
	"github.com/finleap-connect/monoskope/pkg/api/domain/projections"
	esApi "github.com/finleap-connect/monoskope/pkg/api/eventsourcing"
	"github.com/finleap-connect/monoskope/pkg/domain/constants/events"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"strings"
	"time"
)

type HumanReadableEvent struct {
	When   string
	Issuer string
	IssuerId string
	EventType string
	Details string

	event *esApi.Event
	ctx context.Context
	conn *grpc.ClientConn
}


func NewHumanReadableEvent(ctx context.Context, conn *grpc.ClientConn, event *esApi.Event) *HumanReadableEvent {
	hre := &HumanReadableEvent{
		When: event.Timestamp.AsTime().Format(time.RFC822),
		Issuer: event.Metadata["x-auth-email"],
		IssuerId: event.AggregateId,
		EventType: event.Type,

		event: event,
		ctx: ctx,
		conn: conn,
	}
	hre.formatDetails()
	return hre
}

func (hre *HumanReadableEvent) formatDetails() {
	switch hre.EventType {
	case events.UserDeleted.String(): hre.formatDetailsUserDeleted()
	case events.UserRoleBindingDeleted.String(): hre.formatDetailsUserRoleBindingDeleted()
	case events.ClusterDeleted.String(): hre.formatDetailsClusterDeleted()
	case events.TenantDeleted.String(): hre.formatDetailsTenantDeleted()
	case events.TenantClusterBindingDeleted.String(): hre.formatDetailsTenantClusterBindingDeleted()
	}
	if hre.Details != "" {return} // break here if deletion event


	ed, ok := toPortoFromEventData(hre.event.Data)
	if !ok {
		return
	}
	switch ed.(type) {
	case *eventdata.UserCreated: hre.formatDetailsUserCreated(ed.(*eventdata.UserCreated))
	case *eventdata.UserRoleAdded: hre.formatDetailsUserRoleAdded(ed.(*eventdata.UserRoleAdded))
	case *eventdata.ClusterCreated: hre.formatDetailsClusterCreated(ed.(*eventdata.ClusterCreated))
	case *eventdata.ClusterCreatedV2: hre.formatDetailsClusterCreatedV2(ed.(*eventdata.ClusterCreatedV2))
	case *eventdata.ClusterBootstrapTokenCreated: hre.formatDetailsClusterBootstrapTokenCreated(ed.(*eventdata.ClusterBootstrapTokenCreated))
	case *eventdata.ClusterUpdated: hre.formatDetailsClusterUpdated(ed.(*eventdata.ClusterUpdated))
	case *eventdata.TenantCreated: hre.formatDetailsTenantCreated(ed.(*eventdata.TenantCreated))
	case *eventdata.TenantClusterBindingCreated: hre.formatDetailsTenantClusterBindingCreated(ed.(*eventdata.TenantClusterBindingCreated))
	case *eventdata.TenantUpdated: hre.formatDetailsTenantUpdated(ed.(*eventdata.TenantUpdated))
	case *eventdata.CertificateRequested: hre.formatDetailsCertificateRequested(ed.(*eventdata.CertificateRequested))
	case *eventdata.CertificateIssued: hre.formatDetailsCertificateIssued(ed.(*eventdata.CertificateIssued))
	}
}


// TODO: find a pattern to split message based (user, cluster, tenant ...etc)
func (hre *HumanReadableEvent) formatDetailsUserCreated(eventData *eventdata.UserCreated) {
	hre.Details = fmt.Sprintf("%s created user %s", hre.Issuer, eventData.Email)
}

func (hre *HumanReadableEvent) formatDetailsUserRoleAdded(eventData *eventdata.UserRoleAdded) {
	user := hre.getUserById(hre.event.AggregateId)
	hre.Details = fmt.Sprintf("%s assigned the role “%s” for scope “%s” to user %s",
		hre.Issuer, eventData.Role, eventData.Scope, user.Email)
}

func (hre *HumanReadableEvent) formatDetailsUserDeleted() {
	user := hre.getUserById(hre.event.AggregateId)
	hre.Details = fmt.Sprintf("%s deleted %s", hre.Issuer, user.Email)
}

func (hre *HumanReadableEvent) formatDetailsUserRoleBindingDeleted() {
	// TODO: maybe a client?
	urb := &projections.UserRoleBinding{}
	//urb, err := doApi.NewUserClient(hre.conn).GetRoleBindingsById(hre.ctx, &wrapperspb.StringValue{Value: hre.event.AggregateId})
	//if err != nil {
	//	urb = &projections.UserRoleBinding{}
	//}
	user := hre.getUserById(urb.UserId)
	hre.Details = fmt.Sprintf("%s removed the role “%s” for scope “%s” from user %s",
		hre.Issuer, urb.Role, urb.Scope, user.Email)
}

func (hre *HumanReadableEvent) formatDetailsClusterCreated(eventData *eventdata.ClusterCreated) {
	hre.Details = fmt.Sprintf("%s created cluster %s", hre.Issuer, eventData.Name)
}

func (hre *HumanReadableEvent) formatDetailsClusterCreatedV2(eventData *eventdata.ClusterCreatedV2) {
	hre.Details = fmt.Sprintf("%s created cluster %s", hre.Issuer, eventData.Name)
}

func (hre *HumanReadableEvent) formatDetailsClusterBootstrapTokenCreated(_ *eventdata.ClusterBootstrapTokenCreated) {
	hre.Details = fmt.Sprintf("%s created a cluster bootstrap token", hre.Issuer)
}

func (hre *HumanReadableEvent) formatDetailsClusterUpdated(eventData *eventdata.ClusterUpdated) {
	// TODO: how to get a projection of a specific version
	oldCluster := hre.getClusterById(hre.event.AggregateId)
	var details strings.Builder
	details.WriteString(fmt.Sprintf("%s updated the cluster", hre.Issuer))
	appendUpdate("Display name", eventData.DisplayName, oldCluster.DisplayName, &details)
	appendUpdate("API server address", eventData.ApiServerAddress, oldCluster.ApiServerAddress, &details)
	if len(eventData.CaCertificateBundle) != 0 {
		details.WriteString(fmt.Sprintf("\n- Certifcate to a new one"))
	}
	hre.Details = details.String()
}

func (hre *HumanReadableEvent) formatDetailsClusterDeleted() {
	cluster := hre.getClusterById(hre.event.AggregateId)
	hre.Details = fmt.Sprintf("%s deleted %s", hre.Issuer, cluster.DisplayName)
}

func (hre *HumanReadableEvent) formatDetailsTenantCreated(eventData *eventdata.TenantCreated) {
	hre.Details = fmt.Sprintf("%s created Tenant %s", hre.Issuer, eventData.Name)
}

func (hre *HumanReadableEvent) formatDetailsTenantClusterBindingCreated(eventData *eventdata.TenantClusterBindingCreated) {
	tenant := hre.getTenantById(eventData.TenantId)
	cluster := hre.getClusterById(eventData.ClusterId)
	hre.Details = fmt.Sprintf("%s binded tanent “%s” to cluster “%s”",
		hre.Issuer, tenant.Name, cluster.DisplayName)
}

func (hre *HumanReadableEvent) formatDetailsTenantUpdated(eventData *eventdata.TenantUpdated) {
	// TODO: how to get a projection of a specific version
	oldTenant := hre.getTenantById(hre.event.AggregateId)
	var details strings.Builder
	details.WriteString(fmt.Sprintf("%s updated the Tenant", hre.Issuer))
	appendUpdate("Name", eventData.Name.String(), oldTenant.Name, &details)
	hre.Details = details.String()
}

func (hre *HumanReadableEvent) formatDetailsTenantDeleted() {
	tenant := hre.getTenantById(hre.event.AggregateId)
	hre.Details = fmt.Sprintf("%s deleted tenant %s", hre.Issuer, tenant.Name)
}

func (hre *HumanReadableEvent) formatDetailsTenantClusterBindingDeleted() {
	// TODO: client?
	tcb := &projections.TenantClusterBinding{}
	//tenant, err := doApi.NewTenantClient(hre.conn).GetById(hre.ctx, &wrapperspb.StringValue{Value: hre.event.AggregateId})
	//if err != nil {
	//	tenant = &projections.Tenant{}
	//}
	tenant := hre.getTenantById(tcb.TenantId)
	cluster := hre.getClusterById(tcb.ClusterId)
	hre.Details = fmt.Sprintf("%s deleted the bound between cluster %s and tenant %s",
		hre.Issuer, cluster.DisplayName, tenant.Name)
}

func (hre *HumanReadableEvent) formatDetailsCertificateRequested(_ *eventdata.CertificateRequested) {
	hre.Details = fmt.Sprintf("%s requested a certificate", hre.Issuer)
}

func (hre *HumanReadableEvent) formatDetailsCertificateIssued(_ *eventdata.CertificateIssued) {
	// TODO: add context/reason/requester or how to connect both ends
	hre.Details = fmt.Sprintf("%s issued a certificate", hre.Issuer)
}


func (hre *HumanReadableEvent) getUserById(id string) *projections.User {
	return &projections.User{}
	//user, err := doApi.NewUserClient(hre.conn).GetById(hre.ctx, &wrapperspb.StringValue{Value: id})
	//if err != nil {
	//	user = &projections.User{}
	//}
	//return user
}

func (hre *HumanReadableEvent) getClusterById(id string) *projections.Cluster {
	return &projections.Cluster{}
	//cluster, err := doApi.NewClusterClient(hre.conn).GetById(hre.ctx, &wrapperspb.StringValue{Value: id})
	//if err != nil {
	//	cluster = &projections.Cluster{}
	//}
	//return cluster
}

func (hre *HumanReadableEvent) getTenantById(id string) *projections.Tenant {
	return &projections.Tenant{}
	//tenant, err := doApi.NewTenantClient(hre.conn).GetById(hre.ctx, &wrapperspb.StringValue{Value: id})
	//if err != nil {
	//	tenant = &projections.Tenant{}
	//}
	//return tenant
}

func appendUpdate(field string, update string, old string, strBuilder *strings.Builder) {
	if update != "" {
		strBuilder.WriteString(fmt.Sprintf("\n- %s to %s", field, update))
		if old != "" {
			strBuilder.WriteString(fmt.Sprintf(" from %s", old))
		}
	}
}

func toPortoFromEventData(eventData []byte) (proto.Message, bool) {
	porto := &anypb.Any{}
	if err := protojson.Unmarshal(eventData, porto); err != nil {
		return nil, true
	}
	ed, err := porto.UnmarshalNew()
	if err != nil {
		return nil, true
	}
	return ed, false
}