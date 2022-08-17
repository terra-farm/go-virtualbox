package virtualbox

import (
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
)

func TestGuestProperty(t *testing.T) {
	Setup(t)

	t.Logf("ManageMock=%v (type=%T)", ManageMock, ManageMock)
	if ManageMock != nil {
		ManageMock.EXPECT().isGuest().Return(false)
		ManageMock.EXPECT().run("guestproperty", "set", VM, "test_key", "test_val").Return("", "", nil)
	}
	err := SetGuestProperty(VM, "test_key", "test_val")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("OK SetGuestProperty test_key=test_val")

	if ManageMock != nil {
		ManageMock.EXPECT().isGuest().Return(false)
		ManageMock.EXPECT().run("guestproperty", "get", VM, "test_key").Return("Value: test_val", "", nil).Times(1)
	}
	val, err := GetGuestProperty(VM, "test_key")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("val='%s'", val)
	if val != "test_val" {
		t.Fatal("Wrong value")
	}
	Debug("OK GetGuestProperty test_key=test_val")

	// Now deletes it...
	if ManageMock != nil {
		ManageMock.EXPECT().isGuest().Return(false)
		ManageMock.EXPECT().run("guestproperty", "delete", VM, "test_key").Return("", "", nil).Times(1)
	}
	err = DeleteGuestProperty(VM, "test_key")
	if err != nil {
		t.Fatal(err)
	}
	Debug("OK DeleteGuestProperty test_key")

	// ...and check that it is  no longer readable
	if ManageMock != nil {
		ManageMock.EXPECT().isGuest().Return(false)
		ManageMock.EXPECT().run("guestproperty", "get", VM, "test_key").Return("", "", errors.New("foo")).Times(1)
	}
	_, err = GetGuestProperty(VM, "test_key")
	if err == nil {
		t.Fatal(fmt.Errorf("Failed deleting guestproperty"))
	}
	Debug("OK GetGuestProperty test_key=empty")

	Teardown()
}

func TestWaitGuestProperty(t *testing.T) {
	Setup(t)

	keyE, valE := "test_key", "test_val1"
	if ManageMock != nil {
		waitGuestProperty1Out := ReadTestData("vboxmanage-guestproperty-wait-1.out")
		gomock.InOrder(
			ManageMock.EXPECT().isGuest().Return(false),
			ManageMock.EXPECT().run("guestproperty", "wait", VM, "test_*").Return(waitGuestProperty1Out, "", nil).Times(1),
		)
	} else {
		go func() {
			second := time.Second
			time.Sleep(1 * second)
			t.Logf(">>> key='%s', val='%s'", keyE, valE)
			_ = SetGuestProperty(VM, keyE, valE)
		}()
	}

	keyO, valO, err := WaitGuestProperty(VM, "test_*")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("<<< key='%s', val='%s'", keyO, valO)
	if keyE != keyO || valE != valO {
		t.Fatal(errors.New("unexpected key/val"))
	}

	Teardown()
}

func TestWaitGuestProperties(t *testing.T) {
	Setup(t)

	left := 2
	keyE, val1E, val2E := "test_key", "test_val1", "test_val2"

	if ManageMock != nil {
		waitGuestProperty1Out := ReadTestData("vboxmanage-guestproperty-wait-1.out")
		waitGuestProperty2Out := ReadTestData("vboxmanage-guestproperty-wait-2.out")
		gomock.InOrder(
			ManageMock.EXPECT().isGuest().Return(false),
			ManageMock.EXPECT().run("guestproperty", "wait", VM, "test_*").Return(waitGuestProperty1Out, "", nil).Times(1),
			ManageMock.EXPECT().isGuest().Return(false),
			ManageMock.EXPECT().run("guestproperty", "wait", VM, "test_*").Return(waitGuestProperty2Out, "", nil).Times(1),
			ManageMock.EXPECT().isGuest().Return(false),
			ManageMock.EXPECT().run("guestproperty", "wait", VM, "test_*").Return(waitGuestProperty1Out, "", nil).Times(1),
		)
	} else {
		go func() {
			second := time.Second

			time.Sleep(1 * second)
			t.Logf(">>> key='%s', val='%s'", keyE, val1E)
			_ = SetGuestProperty(VM, keyE, val1E)

			time.Sleep(1 * second)
			t.Logf(">>> key='%s', val='%s'", keyE, val2E)
			_ = SetGuestProperty(VM, keyE, val2E)

			time.Sleep(1 * second)
			t.Logf(">>> key='%s', val='%s'", keyE, val1E)
			_ = SetGuestProperty(VM, keyE, val1E)
		}()
	}

	props := "test_*"
	wg := new(sync.WaitGroup)
	done := make(chan bool)

	t.Logf("TestWaitGuestProperties(): will wait on '%s' for %d changes\n", props, left)
	propsC := WaitGuestProperties(VM, props, done, wg)

	t.Logf("TestWaitGuestProperties(): waiting on: %T(%v)\n", propsC, propsC)
	// for prop := range propsChan {
	ok := true
	for ; ok && left > 0; left-- {
		var prop GuestProperty
		t.Logf("TestWaitGuestProperties(): unstacking... (left=%d)\n", left)
		prop, ok = <-propsC
		t.Logf("TestWaitGuestProperties(): unstacked: %+v (left=%d)\n", prop, left)
	}
	t.Logf("TestWaitGuestProperties(): done...\n")
	close(done)
	t.Logf("TestWaitGuestProperties(): done... Ok\n")

	wg.Wait()
	t.Logf("TestWaitGuestProperties(): exiting\n")

	Teardown()
}

func TestWaitGuestPropertiesQuit(t *testing.T) {
	Setup(t)

	keyE, val1E := "test_key", "test_val1"

	if ManageMock != nil {
		waitGuestProperty1Out := ReadTestData("vboxmanage-guestproperty-wait-1.out")
		gomock.InOrder(
			ManageMock.EXPECT().isGuest().Return(false),
			ManageMock.EXPECT().run("guestproperty", "wait", VM, "test_*").Return(waitGuestProperty1Out, "", nil).Times(1),
		)
	} else {
		go func() {
			second := time.Second

			time.Sleep(1 * second)
			t.Logf(">>> key='%s', val='%s'", keyE, val1E)
			_ = SetGuestProperty(VM, keyE, val1E)
		}()
	}

	props := "test_*"
	wg := new(sync.WaitGroup)
	done := make(chan bool)

	t.Logf("TestWaitGuestProperties(): will wait on '%s'\n", props)
	propsC := WaitGuestProperties(VM, props, done, wg)
	t.Logf("TestWaitGuestProperties(): waiting on: %T(%v)\n", propsC, propsC)

	t.Logf("TestWaitGuestProperties(): done...\n")
	close(done)
	t.Logf("TestWaitGuestProperties(): done... Ok\n")

	wg.Wait()
	t.Logf("TestWaitGuestProperties(): exiting\n")

	Teardown()
}
