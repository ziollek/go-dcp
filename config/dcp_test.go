package config

import (
	"testing"
	"time"

	"github.com/Trendyol/go-dcp/helpers"
)

func TestDefaultConfig(t *testing.T) {
	config := Dcp{}
	config.ApplyDefaults()

	if config.ScopeName != "_default" {
		t.Errorf("ScopeName is not set to default")
	}

	if len(config.CollectionNames) != 1 {
		t.Errorf("CollectionNames length is not 1")
	}

	if config.CollectionNames[0] != "_default" {
		t.Errorf("CollectionNames is not set to default")
	}

	if config.Checkpoint.Type != CheckpointTypeAuto {
		t.Errorf("Checkpoint.Type is not set to auto")
	}

	if config.Dcp.Group.Membership.Type != MembershipTypeCouchbase {
		t.Errorf("Dcp.Group.Membership.Type is not set to couchbase")
	}
}

func TestGetCouchbaseMetadata(t *testing.T) {
	dcp := &Dcp{
		Metadata: Metadata{
			Config: map[string]string{
				CouchbaseMetadataBucketConfig: "mybucket",
				CouchbaseMetadataScopeConfig:  "myscope",
			},
		},
		BucketName: "mybucket2",
	}

	couchbaseMetadata := dcp.GetCouchbaseMetadata()

	expectedBucket := "mybucket"
	expectedScope := "myscope"
	expectedCollection := DefaultCollectionName
	expectedConnectionBufferSize := helpers.ResolveUnionIntOrStringValue("5mb")
	expectedConnectionTimeout := time.Minute

	if couchbaseMetadata.Bucket != expectedBucket {
		t.Errorf("Bucket is not set to expected value")
	}

	if couchbaseMetadata.Scope != expectedScope {
		t.Errorf("Scope is not set to expected value")
	}

	if couchbaseMetadata.Collection != expectedCollection {
		t.Errorf("Collection is not set to expected value")
	}

	if couchbaseMetadata.ConnectionBufferSize != uint(expectedConnectionBufferSize) {
		t.Errorf("ConnectionBufferSize is not set to expected value")
	}

	if couchbaseMetadata.ConnectionTimeout != expectedConnectionTimeout {
		t.Errorf("ConnectionTimeout is not set to expected value")
	}
}

func TestGetCouchbaseMembership(t *testing.T) {
	dcp := &Dcp{
		Dcp: ExternalDcp{
			Group: DCPGroup{
				Membership: DCPGroupMembership{
					Config: map[string]string{
						CouchbaseMembershipExpirySecondsConfig:     "120",
						CouchbaseMembershipHeartbeatIntervalConfig: "10s",
						CouchbaseMembershipMonitorIntervalConfig:   "30s",
						CouchbaseMembershipTimeoutConfig:           "30s",
					},
				},
			},
		},
	}

	couchbaseMembership := dcp.GetCouchbaseMembership()

	expectedExpiryDuration := uint32(120)
	expectedHeartbeatInterval := 10 * time.Second
	expectedHeartbeatTolerance := 60 * time.Second
	expectedMonitorInterval := 30 * time.Second
	expectedTimeout := 30 * time.Second

	if couchbaseMembership.ExpirySeconds != expectedExpiryDuration {
		t.Errorf("ExpiryDuration is not set to expected value")
	}

	if couchbaseMembership.HeartbeatInterval != expectedHeartbeatInterval {
		t.Errorf("HeartbeatInterval is not set to expected value")
	}

	if couchbaseMembership.HeartbeatToleranceDuration != expectedHeartbeatTolerance {
		t.Errorf("HeartbeatToleranceDuration is not set to expected value")
	}

	if couchbaseMembership.MonitorInterval != expectedMonitorInterval {
		t.Errorf("MonitorInterval is not set to expected value")
	}

	if couchbaseMembership.Timeout != expectedTimeout {
		t.Errorf("Timeout is not set to expected value")
	}
}

func TestDcp_GetFileMetadata(t *testing.T) {
	dcp := &Dcp{
		Metadata: Metadata{
			Config: map[string]string{
				FileMetadataFileNameConfig: "testfile.json",
			},
		},
	}

	metadata := dcp.GetFileMetadata()

	if metadata != "testfile.json" {
		t.Errorf("Metadata is not set to expected value")
	}
}

func TestApplyDefaultRollbackMitigation(t *testing.T) {
	c := &Dcp{
		RollbackMitigation: RollbackMitigation{},
	}
	c.applyDefaultRollbackMitigation()

	if c.RollbackMitigation.Interval != time.Second {
		t.Errorf("RollbackMitigation.Interval is not set to expected value")
	}

	if c.RollbackMitigation.ConfigWatchInterval != 10*time.Second {
		t.Errorf("RollbackMitigation.ConfigWatchInterval is not set to expected value")
	}
}

func TestDcpApplyDefaultCheckpoint(t *testing.T) {
	c := &Dcp{}
	c.applyDefaultCheckpoint()

	if c.Checkpoint.Interval != time.Minute {
		t.Errorf("Checkpoint.Interval is not set to expected value")
	}

	if c.Checkpoint.Timeout != time.Minute {
		t.Errorf("Checkpoint.Timeout is not set to expected value")
	}

	if c.Checkpoint.Type != CheckpointTypeAuto {
		t.Errorf("Checkpoint.Type is not set to expected value")
	}

	if c.Checkpoint.AutoReset != "earliest" {
		t.Errorf("Checkpoint.AutoReset is not set to expected value")
	}
}

func TestDcpApplyDefaultHealthCheck(t *testing.T) {
	c := &Dcp{}
	c.applyDefaultHealthCheck()

	if c.HealthCheck.Interval != time.Minute {
		t.Errorf("HealthCheck.Interval is not set to expected value")
	}

	if c.HealthCheck.Timeout != time.Minute {
		t.Errorf("HealthCheck.Timeout is not set to expected value")
	}

	if c.HealthCheck.Disabled {
		t.Errorf("HealthCheck.Disabled is not set to expected value")
	}
}

func TestDcpApplyDefaultGroupMembership(t *testing.T) {
	c := &Dcp{}
	c.applyDefaultGroupMembership()

	if c.Dcp.Group.Membership.RebalanceDelay != 30*time.Second {
		t.Errorf("Dcp.Group.Membership.RebalanceDelay is not set to expected value")
	}

	if c.Dcp.Group.Membership.TotalMembers != 1 {
		t.Errorf("Dcp.Group.Membership.TotalMembers is not set to expected value")
	}

	if c.Dcp.Group.Membership.MemberNumber != 1 {
		t.Errorf("Dcp.Group.Membership.MemberNumber is not set to expected value")
	}

	if c.Dcp.Group.Membership.Type != "couchbase" {
		t.Errorf("Dcp.Group.Membership.Type is not set to expected value")
	}
}

func TestDcpApplyDefaultConnectionTimeout(t *testing.T) {
	c := &Dcp{}
	c.applyDefaultConnectionTimeout()

	if c.Dcp.ConnectionTimeout != time.Minute {
		t.Errorf("Dcp.ConnectionTimeout is not set to expected value")
	}

	if c.ConnectionTimeout != time.Minute {
		t.Errorf("ConnectionTimeout is not set to expected value")
	}
}

func TestDcpApplyDefaultCollections(t *testing.T) {
	c := &Dcp{}
	c.applyDefaultCollections()

	if c.CollectionNames[0] != DefaultCollectionName {
		t.Errorf("CollectionNames is not set to expected value")
	}
}

func TestDcpApplyDefaultScopeName(t *testing.T) {
	c := &Dcp{}
	c.applyDefaultScopeName()

	if c.ScopeName != DefaultScopeName {
		t.Errorf("ScopeName is not set to expected value")
	}
}

func TestDcpApplyDefaultConnectionBufferSize(t *testing.T) {
	c := &Dcp{}
	c.applyDefaultConnectionBufferSize()

	if c.ConnectionBufferSize.(int) != 20971520 {
		t.Errorf("ConnectionBufferSize is not set to expected value")
	}
}

func TestDcpApplyDefaultMaxQueueSize(t *testing.T) {
	c := &Dcp{}
	c.applyDefaultMaxQueueSize()

	if c.MaxQueueSize != 2048 {
		t.Errorf("ConnectionBufferSize is not set to expected value")
	}
}

func TestDcpApplyDefaultMetrics(t *testing.T) {
	c := &Dcp{}
	c.applyDefaultMetrics()

	if c.Metric.Path != "/metrics" {
		t.Errorf("Metric.Path is not set to expected value")
	}
}

func TestDcpApplyDefaultAPI(t *testing.T) {
	c := &Dcp{}
	c.applyDefaultAPI()

	if c.API.Disabled {
		t.Errorf("API.Disabled is not set to expected value")
	}

	if c.API.Port != 8080 {
		t.Errorf("API.Port is not set to expected value")
	}
}

func TestDcpApplyDefaultLeaderElection(t *testing.T) {
	c := &Dcp{}
	c.applyDefaultLeaderElection()

	if c.LeaderElection.Enabled {
		t.Errorf("LeaderElection.Enabled is not set to expected value")
	}

	if c.LeaderElection.Type != "kubernetes" {
		t.Errorf("LeaderElection.Type is not set to expected value")
	}

	if c.LeaderElection.RPC.Port != 8081 {
		t.Errorf("LeaderElection.RPC.Port is not set to expected value")
	}
}

func TestDcpApplyDefaultDcp(t *testing.T) {
	c := &Dcp{}
	c.applyDefaultDcp()

	if c.Dcp.BufferSize.(int) != 16777216 {
		t.Errorf("Dcp.BufferSize is not set to expected value")
	}

	if c.Dcp.ConnectionBufferSize.(int) != 20971520 {
		t.Errorf("Dcp.ConnectionBufferSize is not set to expected value")
	}
}

func TestApplyDefaultMetadata(t *testing.T) {
	// Initialize a Dcp instance with no metadata
	c := &Dcp{
		BucketName: "my-bucket",
		Metadata:   Metadata{},
	}

	// Apply default metadata
	c.applyDefaultMetadata()

	// Check if the default metadata values were applied correctly
	if c.Metadata.Type != "couchbase" {
		t.Errorf("Metadata.Type is not set to expected value")
	}
}

func TestDcpMode(t *testing.T) {
	t.Run("it should return true when dcp mode finite", func(t *testing.T) {
		// Arrange
		dcp := &Dcp{
			Dcp: ExternalDcp{
				Mode: DcpModeFinite,
			},
		}

		expectedValue := true

		// Act
		actualValue := dcp.IsDcpModeFinite()

		// Assert
		if expectedValue != actualValue {
			t.Errorf("isDcpModeFinite check result. got %v want %v", actualValue, expectedValue)
		}
	})

	t.Run("it should return false when dcp mode infinite", func(t *testing.T) {
		// Arrange
		dcp := &Dcp{
			Dcp: ExternalDcp{
				Mode: DcpModeInfinite,
			},
		}

		expectedValue := false

		// Act
		actualValue := dcp.IsDcpModeFinite()

		// Assert
		if expectedValue != actualValue {
			t.Errorf("isDcpModeFinite check result. got %v want %v", actualValue, expectedValue)
		}
	})

	t.Run("it should return false when dcp mode empty", func(t *testing.T) {
		// Arrange
		dcp := &Dcp{
			Dcp: ExternalDcp{
				Mode: "",
			},
		}
		expectedValue := false

		// Act
		actualValue := dcp.IsDcpModeFinite()

		// Assert
		if expectedValue != actualValue {
			t.Errorf("isDcpModeFinite check result. got %v want %v", actualValue, expectedValue)
		}
	})
}
