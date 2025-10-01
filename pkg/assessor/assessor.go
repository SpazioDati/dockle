package assessor

import (
	"os"

	deckodertypes "github.com/goodwithtech/deckoder/types"

	"github.com/SpazioDati/dockle/pkg/assessor/cache"
	"github.com/SpazioDati/dockle/pkg/assessor/privilege"

	"github.com/SpazioDati/dockle/pkg/assessor/contentTrust"
	"github.com/SpazioDati/dockle/pkg/assessor/credential"
	"github.com/SpazioDati/dockle/pkg/assessor/hosts"

	"github.com/SpazioDati/dockle/pkg/assessor/group"
	"github.com/SpazioDati/dockle/pkg/assessor/manifest"
	"github.com/SpazioDati/dockle/pkg/assessor/passwd"
	"github.com/SpazioDati/dockle/pkg/assessor/user"

	"github.com/SpazioDati/dockle/pkg/log"
	"github.com/SpazioDati/dockle/pkg/types"
)

var assessors []Assessor

type Assessor interface {
	Assess(deckodertypes.FileMap) ([]*types.Assessment, error)
	RequiredFiles() []string
	RequiredExtensions() []string
	RequiredPermissions() []os.FileMode
}

func init() {
	RegisterAssessor(passwd.PasswdAssessor{})
	RegisterAssessor(privilege.PrivilegeAssessor{})
	RegisterAssessor(user.UserAssessor{})
	RegisterAssessor(group.GroupAssessor{})
	RegisterAssessor(hosts.HostsAssessor{})
	RegisterAssessor(credential.CredentialAssessor{})
	RegisterAssessor(manifest.ManifestAssessor{})
	RegisterAssessor(contentTrust.ContentTrustAssessor{})
	RegisterAssessor(cache.CacheAssessor{})
}

func GetAssessments(files deckodertypes.FileMap) (assessments []*types.Assessment) {
	for _, assessor := range assessors {
		results, err := assessor.Assess(files)
		if err != nil {
			log.Logger.Error(err)
		}
		assessments = append(assessments, results...)
	}
	return assessments
}

func RegisterAssessor(a Assessor) {
	assessors = append(assessors, a)
}

func LoadRequiredFiles() (filenames []string) {
	for _, assessor := range assessors {
		filenames = append(filenames, assessor.RequiredFiles()...)
	}
	return filenames
}

func LoadRequiredExtensions() (extensions []string) {
	for _, assessor := range assessors {
		extensions = append(extensions, assessor.RequiredExtensions()...)
	}
	return extensions
}

func LoadRequiredPermissions() (permissions []os.FileMode) {
	for _, assessor := range assessors {
		permissions = append(permissions, assessor.RequiredPermissions()...)
	}
	return permissions
}
