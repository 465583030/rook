package smoke

import (
	"testing"

	"github.com/coreos/pkg/capnslog"
	"github.com/rook/rook/tests/framework/clients"
	"github.com/rook/rook/tests/framework/enums"
	"github.com/rook/rook/tests/framework/installer"
	"github.com/rook/rook/tests/framework/utils"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"time"
)

var (
	hslogger = capnslog.NewPackageLogger("github.com/rook/rook", "helmSmokeTest")
)

func TestHelmSmokeSuite(t *testing.T) {
	suite.Run(t, new(HelmSmokeSuite))
}

type HelmSmokeSuite struct {
	suite.Suite
	helper    *clients.TestClient
	k8sh      *utils.K8sHelper
	installer *installer.InstallHelper
	hh        *utils.HelmHelper
}

func (hs *HelmSmokeSuite) SetupSuite() {
	kh, err := utils.CreateK8sHelper()
	require.NoError(hs.T(), err)

	hs.k8sh = kh
	hs.hh = utils.NewHelmHelper()

	hs.installer = installer.NewK8sRookhelper(kh.Clientset)

	err = hs.installer.CreateK8sRookOperatorViaHelm(defaultRookNamespace)
	require.NoError(hs.T(), err)

	require.True(hs.T(), kh.IsPodInExpectedState("rook-operator", defaultRookNamespace, "Running"),
		"Make sure rook-operator is in running state")

	time.Sleep(10 * time.Second)

	err = hs.installer.CreateK8sRookCluster(defaultRookNamespace)
	require.NoError(hs.T(), err)

	err = hs.installer.CreateK8sRookToolbox(defaultRookNamespace)
	require.NoError(hs.T(), err)

	hs.helper, err = clients.CreateTestClient(enums.Kubernetes, kh, defaultRookNamespace)
	require.Nil(hs.T(), err)
}

func (hs *HelmSmokeSuite) TearDownSuite() {
	if hs.T().Failed() {
		gatherAllRookLogs(hs.k8sh, hs.Suite, defaultRookNamespace, defaultRookNamespace)

	}
	hs.installer.UninstallRookFromK8s(defaultRookNamespace, true)

}

//Test to make sure all rook components are installed and Running
func (hs *HelmSmokeSuite) TestRookInstallViaHelm() {
	checkIfRookClusterIsInstalled(hs.k8sh, hs.Suite, defaultRookNamespace, defaultRookNamespace)

}

//Test BlockCreation on Rook that was installed via Helm
func (hs *HelmSmokeSuite) TestBlockStoreOnRookInstalledViaHelm() {
	runBlockE2ETestLite(hs.helper, hs.k8sh, hs.Suite, defaultRookNamespace)

}

//Test File System Creation on Rook that was installed via helm
func (hs *HelmSmokeSuite) TestFileStoreOnRookInstalledViaHelm() {
	runFileE2ETestLite(hs.helper, hs.k8sh, hs.Suite, defaultRookNamespace, "testFs")

}

//Test Object StoreCreation on Rook that was installed via helm
func (hs *HelmSmokeSuite) TestObjectStoreOnRookInstalledViaHelm() {
	runObjectE2ETestLite(hs.helper, hs.k8sh, hs.Suite, defaultRookNamespace, "default", 3)

}
