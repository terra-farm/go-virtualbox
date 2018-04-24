package virtualbox

import (
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/golang/mock/gomock"
)

var (
	MockCtrl       *gomock.Controller
	ManageMock     *MockCommand
	VM             string
	TestDataFolder string
)

func ReadTestData(file string) string {
	out, err := ioutil.ReadFile(path.Join("testdata", file))
	if err != nil {
		panic("No such file :testdata/" + file)
	}
	return string(out)
}

func Setup(t *testing.T) {
	Verbose = true

	VM = os.Getenv("TEST_VM")
	if len(VM) <= 0 {
		VM = "go-virtualbox"
		t.Log("Missing TEST_VM environment variable")
	}
	t.Logf("Using VM='%s'\n", VM)

	MockCtrl = gomock.NewController(t)
	if len(os.Getenv("TEST_MOCK_VBM")) > 0 {
		ManageMock = NewMockCommand(MockCtrl)
		Manage = ManageMock
	}
	t.Logf("Using VBoxManage='%T'", Manage)
	t.Logf("ManageMock=%v (type=%T)", ManageMock, ManageMock)
}

func Teardown() {
	defer MockCtrl.Finish()
}

func TestVBMOut(t *testing.T) {
	Setup(t)

	t.Logf("VM=%s", VM)
	t.Logf("ManageMock=%v (type=%T)", ManageMock, ManageMock)
	if ManageMock != nil {
		var out = "\"Ubuntu\" {2e16b1fc-aaaa-4a7a-a9a1-e89a8bde7874}\n" +
			"\"go-virtualbox\" {def44546-aaaa-4902-8d15-b91c99c80cbc}"
		ManageMock.EXPECT().runOut("list", "vms").Return(out, nil)
	}
	b, err := Manage.runOut("list", "vms")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%s", b)

	Teardown()
}
