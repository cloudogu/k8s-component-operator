package client

import (
	"context"
	"github.com/sirupsen/logrus"
	"helm.sh/helm/v3/pkg/cli"
	"os"
	"strconv"
	"time"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/release"
)

type provider struct {
	*action.Configuration
	plainHttp   bool
	insecureTls bool
}

var defaultRollbackTimeout = 15

func (p *provider) newInstall() installAction {
	installAction := action.NewInstall(p.Configuration)
	installAction.PlainHTTP = p.plainHttp
	installAction.InsecureSkipTLSverify = p.insecureTls
	return &install{Install: installAction}
}

func (p *provider) newUpgrade() upgradeAction {
	upgradeAction := action.NewUpgrade(p.Configuration)
	upgradeAction.PlainHTTP = p.plainHttp
	upgradeAction.InsecureSkipTLSverify = p.insecureTls
	return &upgrade{Upgrade: upgradeAction}
}

func (p *provider) newUninstall() uninstallAction {
	uninstallAction := action.NewUninstall(p.Configuration)
	return &uninstall{Uninstall: uninstallAction}
}

func (p *provider) newLocateChart() locateChartAction {
	dummyAction := action.NewInstall(p.Configuration)
	dummyAction.PlainHTTP = p.plainHttp
	dummyAction.InsecureSkipTLSverify = p.insecureTls
	return &locateChart{dummyAction: dummyAction}
}

func (p *provider) newListReleases() listReleasesAction {
	listAction := action.NewList(p.Configuration)
	return &listReleases{List: listAction}
}

func (p *provider) newGetReleaseValues() getReleaseValuesAction {
	getValuesAction := action.NewGetValues(p.Configuration)
	return &getReleaseValues{GetValues: getValuesAction}
}

func (p *provider) newGetRelease() getReleaseAction {
	getAction := action.NewGet(p.Configuration)
	return &getRelease{Get: getAction}
}

func (p *provider) newRollbackRelease() rollbackReleaseAction {
	rollbackAction := action.NewRollback(p.Configuration)
	const rollbackTimeoutEnv = "ROLLBACK_TIMEOUT"
	rollbackTimeoutString, found := os.LookupEnv(rollbackTimeoutEnv)
	rollbackTimeout, err := strconv.Atoi(rollbackTimeoutString)
	if !found || err != nil {
		logrus.Warningf("failed to read %s environment variable, using default value of %d", rollbackTimeoutEnv, defaultRollbackTimeout)
		rollbackTimeout = defaultRollbackTimeout
	}
	rollbackAction.Timeout = time.Duration(rollbackTimeout) * time.Minute
	return &rollbackRelease{Rollback: rollbackAction}
}

type install struct {
	*action.Install
}

func (i *install) install(ctx context.Context, chart *chart.Chart, values map[string]interface{}) (*release.Release, error) {
	return i.RunWithContext(ctx, chart, values)
}

func (i *install) raw() *action.Install {
	return i.Install
}

type upgrade struct {
	*action.Upgrade
}

func (u *upgrade) upgrade(ctx context.Context, releaseName string, chart *chart.Chart, values map[string]interface{}) (*release.Release, error) {
	return u.RunWithContext(ctx, releaseName, chart, values)
}

func (u *upgrade) raw() *action.Upgrade {
	return u.Upgrade
}

type uninstall struct {
	*action.Uninstall
}

func (u *uninstall) uninstall(releaseName string) (*release.UninstallReleaseResponse, error) {
	return u.Run(releaseName)
}

func (u *uninstall) raw() *action.Uninstall {
	return u.Uninstall
}

type locateChart struct {
	dummyAction *action.Install
}

func (l *locateChart) locateChart(name, version string, settings *cli.EnvSettings) (chartPath string, err error) {
	l.dummyAction.Version = version
	return l.dummyAction.ChartPathOptions.LocateChart(name, settings)
}

type listReleases struct {
	*action.List
}

func (l *listReleases) listReleases() ([]*release.Release, error) {
	return l.Run()
}

func (l *listReleases) raw() *action.List {
	return l.List
}

type getReleaseValues struct {
	*action.GetValues
}

func (v *getReleaseValues) getReleaseValues(releaseName string) (map[string]interface{}, error) {
	return v.Run(releaseName)
}

func (v *getReleaseValues) raw() *action.GetValues {
	return v.GetValues
}

type getRelease struct {
	*action.Get
}

func (g *getRelease) getRelease(releaseName string) (*release.Release, error) {
	return g.Run(releaseName)
}

func (g *getRelease) raw() *action.Get {
	return g.Get
}

type rollbackRelease struct {
	*action.Rollback
}

func (r *rollbackRelease) rollbackRelease(releaseName string) error {
	return r.Run(releaseName)
}

func (r *rollbackRelease) raw() *action.Rollback {
	return r.Rollback
}
