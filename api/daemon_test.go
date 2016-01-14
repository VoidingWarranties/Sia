package api

import (
	"reflect"
	"testing"

	"github.com/NebulousLabs/Sia/build"
	"github.com/NebulousLabs/Sia/types"
)

// TestVersion checks that the version returned by the daemon is correct.
func TestVersion(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}
	st, err := createServerTester("TestVersion")
	if err != nil {
		t.Fatal(err)
	}
	defer st.server.Close()
	var dv DaemonVersion
	st.getAPI("/daemon/version", &dv)
	if dv.Version != build.Version {
		t.Fatalf("/daemon/version reporting bad version: expected %v, got %v", build.Version, dv.Version)
	}
}

// TestConstants checks that the constants returned by the daemon are correct.
func TestConstants(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}
	st, err := createServerTester("TestConstants")
	if err != nil {
		t.Fatal(err)
	}
	defer st.server.Close()
	var sc SiaConstants
	st.getAPI("/daemon/constants", &sc)
	// Fail if fields are added or removed from SiaConstants to ensure we are
	// testing all fields. If this test fails then a field comparison test
	// needs to be added or removed below.
	if numConstants := reflect.TypeOf(sc).NumField(); numConstants != 16 {
		t.Errorf("daemon/constants reporting an incorrect number of constants: expected %v, got %v", 16, numConstants)
	}
	// Test all fields in SiaConstants.
	if sc.GenesisTimestamp != types.GenesisTimestamp {
		t.Errorf("/daemon/constants reporting bad genesis timestamp: expected %v, got %v", types.GenesisTimestamp, sc.GenesisTimestamp)
	}
	if sc.BlockSizeLimit != types.BlockSizeLimit {
		t.Errorf("/daemon/constants reporting bad block size limit: expected %v, got %v", types.BlockSizeLimit, sc.BlockSizeLimit)
	}
	if sc.BlockFrequency != types.BlockFrequency {
		t.Errorf("/daemon/constants reporting bad block frequency: expected %v, got %v", types.BlockFrequency, sc.BlockFrequency)
	}
	if sc.TargetWindow != types.TargetWindow {
		t.Errorf("/daemon/constants reporting bad target window: expected %v, got %v", types.TargetWindow, sc.TargetWindow)
	}
	if sc.MedianTimestampWindow != types.MedianTimestampWindow {
		t.Errorf("/daemon/constants reporting bad median timestamp window: expected %v, got %v", types.MedianTimestampWindow, sc.MedianTimestampWindow)
	}
	if sc.FutureThreshold != types.FutureThreshold {
		t.Errorf("/daemon/constants reporting bad future threshold: expected %v, got %v", types.FutureThreshold, sc.FutureThreshold)
	}
	if sc.SiafundCount.Cmp(types.SiafundCount) != 0 {
		t.Errorf("/daemon/constants reporting bad siafund count: expected %v, got %v", types.SiafundCount, sc.SiafundCount)
	}
	if sc.SiafundPortion.Cmp(types.SiafundPortion) != 0 {
		t.Errorf("/daemon/constants reporting bad siafund portion: expected %v, got %v", types.SiafundPortion, sc.SiafundPortion)
	}
	if sc.MaturityDelay != types.MaturityDelay {
		t.Errorf("/daemon/constants reporting bad maturity delay: expected %v, got %v", types.MaturityDelay, sc.MaturityDelay)
	}
	if sc.InitialCoinbase != types.InitialCoinbase {
		t.Errorf("/daemon/constants reporting bad initial coinbase: expected %v, got %v", types.InitialCoinbase, sc.InitialCoinbase)
	}
	if sc.MinimumCoinbase != types.MinimumCoinbase {
		t.Errorf("/daemon/constants reporting bad minimum coinbase: expected %v, got %v", types.MinimumCoinbase, sc.MinimumCoinbase)
	}
	if sc.RootTarget != types.RootTarget {
		t.Errorf("/daemon/constants reporting bad root target: expected %v, got %v", types.RootTarget, sc.RootTarget)
	}
	if sc.RootDepth != types.RootDepth {
		t.Errorf("/daemon/constants reporting bad root depth: expected %v, got %v", types.RootDepth, sc.RootDepth)
	}
	if sc.MaxAdjustmentUp.Cmp(types.MaxAdjustmentUp) != 0 {
		t.Errorf("/daemon/constants reporting bad max adjustment up: expected %v, got %v", types.MaxAdjustmentUp, sc.MaxAdjustmentUp)
	}
	if sc.MaxAdjustmentDown.Cmp(types.MaxAdjustmentDown) != 0 {
		t.Errorf("/daemon/constants reporting bad max adjustment down: expected %v, got %v", types.MaxAdjustmentDown, sc.MaxAdjustmentDown)
	}
	if sc.SiacoinPrecision.Cmp(types.SiacoinPrecision) != 0 {
		t.Errorf("/daemon/constants reporting bad siacoin precisiondown: expected %v, got %v", types.SiacoinPrecision, sc.SiacoinPrecision)
	}
}

// TestStop checks that a call to /daemon/stop returns a response with Success == true
// TestStop does not check that the daemon stopped.
func TestStop(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}
	st, err := createServerTester("TestStop")
	if err != nil {
		t.Fatal(err)
	}
	defer st.server.Close()
	var response struct{ Success bool }
	st.getAPI("/daemon/stop", &response)
	if response.Success != true {
		t.Fatal("/daemon/stop not reporting success")
	}
	// TODO: check that the daemon has stopped.
}
