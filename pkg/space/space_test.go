package space

import (
	"testing"

	toolchainv1alpha1 "github.com/codeready-toolchain/api/api/v1alpha1"
	spacerequesttest "github.com/codeready-toolchain/host-operator/test/spacerequest"
	"github.com/codeready-toolchain/toolchain-common/pkg/cluster"
	commoncluster "github.com/codeready-toolchain/toolchain-common/pkg/cluster"
	"github.com/codeready-toolchain/toolchain-common/pkg/test"
	commonsignup "github.com/codeready-toolchain/toolchain-common/pkg/test/usersignup"

	spacetest "github.com/codeready-toolchain/toolchain-common/pkg/test/space"

	"github.com/stretchr/testify/assert"
)

func TestNewSpace(t *testing.T) {
	// given
	userSignup := commonsignup.NewUserSignup()

	// when
	space := NewSpace(userSignup, test.MemberClusterName, "johny", "advanced")

	// then
	expectedSpace := spacetest.NewSpace(test.HostOperatorNs, "johny",
		spacetest.WithTierName("advanced"),
		spacetest.WithSpecTargetCluster("member-cluster"),
		spacetest.WithSpecTargetClusterRoles([]string{cluster.RoleLabel(cluster.Tenant)}),
		spacetest.WithCreatorLabel(userSignup.Name))
	assert.Equal(t, expectedSpace, space)
}

func TestNewSubSpace(t *testing.T) {
	// given
	srClusterRoles := []string{commoncluster.RoleLabel(commoncluster.Tenant)}
	sr := spacerequesttest.NewSpaceRequest("jane", "jane-tenant",
		spacerequesttest.WithTierName("appstudio"),
		spacerequesttest.WithTargetClusterRoles(srClusterRoles))
	parentSpace := spacetest.NewSpace(test.HostOperatorNs, "parentSpace")

	// when
	subSpace := NewSubSpace(sr, parentSpace)

	// then
	expectedSubSpace := spacetest.NewSpace(test.HostOperatorNs, SubSpaceName(parentSpace, sr),
		spacetest.WithSpecParentSpace(parentSpace.GetName()),
		spacetest.WithTierName("appstudio"),
		spacetest.WithSpecTargetClusterRoles([]string{cluster.RoleLabel(cluster.Tenant)}),
		spacetest.WithLabel(toolchainv1alpha1.SpaceRequestLabelKey, sr.GetName()),
		spacetest.WithLabel(toolchainv1alpha1.SpaceRequestNamespaceLabelKey, sr.GetNamespace()),
		spacetest.WithLabel(toolchainv1alpha1.ParentSpaceLabelKey, parentSpace.GetName()),
	)
	assert.Equal(t, expectedSubSpace, subSpace)
}
