package exporter

import (
	"fmt"
	"github.com/ktsstudio/selectel-exporter/pkg/selapi"
	"github.com/prometheus/client_golang/prometheus"
)

type primaryMetrics struct {
	main prometheus.Gauge
	bonus prometheus.Gauge
	vkRub prometheus.Gauge
	ref prometheus.Gauge
	holdMain prometheus.Gauge
	holdBonus prometheus.Gauge
	holdVkRub prometheus.Gauge
}

type storageMetrics struct {
	main  prometheus.Gauge
	bonus prometheus.Gauge
	vkRub prometheus.Gauge
	debt  prometheus.Gauge
	sum   prometheus.Gauge
}

type vpcMetrics storageMetrics
type vmwareMetrics storageMetrics

type balanceCollector struct {
	project selapi.Project

	primary primaryMetrics
	storage storageMetrics
	vpc vpcMetrics
	vmware vmwareMetrics
}

func registerGauge(account, name string, project selapi.Project) prometheus.Gauge {
	g := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: name,
		ConstLabels: prometheus.Labels{"project": project.Name, "account": account},
	})
	Registry.MustRegister(g)
	return g
}

func registerPrimaryGauge(name string, project selapi.Project) prometheus.Gauge {
	return registerGauge("primary", name, project)
}

func registerStorageGauge(name string, project selapi.Project) prometheus.Gauge {
	return registerGauge("storage", name, project)
}

func registerVpcGauge(name string, project selapi.Project) prometheus.Gauge {
	return registerGauge("vpc", name, project)
}

func registerVmwareGauge(name string, project selapi.Project) prometheus.Gauge {
	return registerGauge("vmware", name, project)
}

func NewBalanceCollector(project selapi.Project) *balanceCollector {
	c := &balanceCollector{project: project}

	c.primary.main = registerPrimaryGauge("selectel_billing_main", project)
	c.primary.bonus = registerPrimaryGauge("selectel_billing_bonus", project)
	c.primary.vkRub = registerPrimaryGauge("selectel_billing_vk_rub", project)
	c.primary.ref = registerPrimaryGauge("selectel_billing_ref", project)
	c.primary.holdMain = registerPrimaryGauge("selectel_billing_hold_main", project)
	c.primary.holdBonus = registerPrimaryGauge("selectel_billing_hold_bonus", project)
	c.primary.holdVkRub = registerPrimaryGauge("selectel_billing_hold_vk_rub", project)

	c.storage.main = registerStorageGauge("selectel_billing_main", project)
	c.storage.bonus = registerStorageGauge("selectel_billing_bonus", project)
	c.storage.vkRub = registerStorageGauge("selectel_billing_vk_rub", project)
	c.storage.debt = registerStorageGauge("selectel_billing_debt", project)
	c.storage.sum = registerStorageGauge("selectel_billing_sum", project)

	c.vpc.main = registerVpcGauge("selectel_billing_main", project)
	c.vpc.bonus = registerVpcGauge("selectel_billing_bonus", project)
	c.vpc.vkRub = registerVpcGauge("selectel_billing_vk_rub", project)
	c.vpc.debt = registerVpcGauge("selectel_billing_debt", project)
	c.vpc.sum = registerVpcGauge("selectel_billing_sum", project)

	c.vmware.main = registerVmwareGauge("selectel_billing_main", project)
	c.vmware.bonus = registerVmwareGauge("selectel_billing_bonus", project)
	c.vmware.vkRub = registerVmwareGauge("selectel_billing_vk_rub", project)
	c.vmware.debt = registerVmwareGauge("selectel_billing_debt", project)
	c.vmware.sum = registerVmwareGauge("selectel_billing_sum", project)

	return c
}

func (col *balanceCollector) GetInfo() string {
	return fmt.Sprintf("project: %s - collect balance metrics", col.project.Name)
}

func (col *balanceCollector) Collect(e *exporter) error {
	res, err := selapi.FetchBalance(e.token)
	if err != nil {
		return err
	}

	col.primary.main.Set(float64(res.Data.Primary.Main))
	col.primary.bonus.Set(float64(res.Data.Primary.Bonus))
	col.primary.vkRub.Set(float64(res.Data.Primary.VkRub))
	col.primary.ref.Set(float64(res.Data.Primary.Ref))
	col.primary.holdMain.Set(float64(res.Data.Primary.Hold.Main))
	col.primary.holdBonus.Set(float64(res.Data.Primary.Hold.Bonus))
	col.primary.holdVkRub.Set(float64(res.Data.Primary.Hold.VkRub))

	col.storage.main.Set(float64(res.Data.Storage.Main))
	col.storage.bonus.Set(float64(res.Data.Storage.Bonus))
	col.storage.vkRub.Set(float64(res.Data.Storage.VkRub))
	col.storage.debt.Set(float64(res.Data.Storage.Debt))
	col.storage.sum.Set(float64(res.Data.Storage.Sum))

	col.vpc.main.Set(float64(res.Data.Vpc.Main))
	col.vpc.bonus.Set(float64(res.Data.Vpc.Bonus))
	col.vpc.vkRub.Set(float64(res.Data.Vpc.VkRub))
	col.vpc.debt.Set(float64(res.Data.Vpc.Debt))
	col.vpc.sum.Set(float64(res.Data.Vpc.Sum))

	col.vmware.main.Set(float64(res.Data.Vmware.Main))
	col.vmware.bonus.Set(float64(res.Data.Vmware.Bonus))
	col.vmware.vkRub.Set(float64(res.Data.Vmware.VkRub))
	col.vmware.debt.Set(float64(res.Data.Vmware.Debt))
	col.vmware.sum.Set(float64(res.Data.Vmware.Sum))

	return nil
}
