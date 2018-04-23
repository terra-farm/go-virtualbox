package virtualbox

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestGuestProperty(t *testing.T) {

	// Setup

	var vm = os.Getenv("TEST_VM")
	if len(vm) <= 0 {
		vm = "go-virtualbox"
		t.Logf("Missing TEST_VM environment variable")
	}
	t.Logf("Using VM='%s'", vm)

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	var ManageMock *MockCommand
	if len(os.Getenv("TEST_MOCK_VBM")) > 0 {
		ManageMock = NewMockCommand(mockCtrl)
		Manage = ManageMock
	}
	t.Logf("Using VBoxManage='%T'", Manage)

	// Tests
	if ManageMock != nil {
		ManageMock.EXPECT().run("guestproperty", "set", vm, "test_key", "test_val").Return(nil)
	}
	err := SetGuestProperty(vm, "test_key", "test_val")
	if err != nil {
		t.Fatal(err)
	}
	if Verbose {
		t.Logf("OK SetGuestProperty test_key=test_val")
	}

	if ManageMock != nil {
		ManageMock.EXPECT().runOut("guestproperty", "get", vm, "test_key").Return("Value: test_val", nil).Times(1)
	}
	val, err := GetGuestProperty(vm, "test_key")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("val='%s'", val)
	if val != "test_val" {
		t.Fatal("Wrong value")
	}
	if Verbose {
		t.Logf("OK GetGuestProperty test_key=test_val")
	}

	// Now deletes it...
	if ManageMock != nil {
		ManageMock.EXPECT().run("guestproperty", "delete", vm, "test_key").Return(nil).Times(1)
	}
	err = DeleteGuestProperty(vm, "test_key")
	if err != nil {
		t.Fatal(err)
	}
	if Verbose {
		t.Logf("OK DeleteGuestProperty test_key")
	}

	// ...and check that it is  no longer readable
	if ManageMock != nil {
		ManageMock.EXPECT().runOut("guestproperty", "get", vm, "test_key").Return("", errors.New("foo")).Times(1)
	}
	_, err = GetGuestProperty(vm, "test_key")
	if err == nil {
		t.Fatal(fmt.Errorf("Failed deleting guestproperty"))
	}
	if Verbose {
		t.Logf("OK GetGuestProperty test_key=empty")
	}

}
