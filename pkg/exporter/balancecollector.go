package exporter

import (
	"github.com/ktsstudio/selectel-exporter/pkg/selapi"
	"github.com/prometheus/client_golang/prometheus"
	"log"
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
	primary primaryMetrics
	storage storageMetrics
	vpc vpcMetrics
	vmware vmwareMetrics
}

func registerGauge(name string, project selapi.Project) prometheus.Gauge {
	g := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: name,
		ConstLabels: prometheus.Labels{"project": project.Name},
	})
	prometheus.MustRegister(g)
	return g
}

func NewBalanceCollector(project selapi.Project) *balanceCollector {
	c := &balanceCollector{}

	c.primary.main = registerGauge("selectel_billing_primary_main", project)
	c.primary.bonus = registerGauge("selectel_billing_primary_bonus", project)
	c.primary.vkRub = registerGauge("selectel_billing_primary_vk_rub", project)
	c.primary.ref = registerGauge("selectel_billing_primary_ref", project)
	c.primary.holdMain = registerGauge("selectel_billing_primary_hold_main", project)
	c.primary.holdBonus = registerGauge("selectel_billing_primary_hold_bonus", project)
	c.primary.holdVkRub = registerGauge("selectel_billing_primary_hold_vk_rub", project)

	c.storage.main = registerGauge("selectel_billing_storage_main", project)
	c.storage.bonus = registerGauge("selectel_billing_storage_bonus", project)
	c.storage.vkRub = registerGauge("selectel_billing_storage_vk_rub", project)
	c.storage.debt = registerGauge("selectel_billing_storage_debt", project)
	c.storage.sum = registerGauge("selectel_billing_storage_sum", project)

	c.vpc.main = registerGauge("selectel_billing_vpc_main", project)
	c.vpc.bonus = registerGauge("selectel_billing_vpc_bonus", project)
	c.vpc.vkRub = registerGauge("selectel_billing_vpc_vk_rub", project)
	c.vpc.debt = registerGauge("selectel_billing_vpc_debt", project)
	c.vpc.sum = registerGauge("selectel_billing_vpc_sum", project)

	c.vmware.main = registerGauge("selectel_billing_vmware_main", project)
	c.vmware.bonus = registerGauge("selectel_billing_vmware_bonus", project)
	c.vmware.vkRub = registerGauge("selectel_billing_vmware_vk_rub", project)
	c.vmware.debt = registerGauge("selectel_billing_vmware_debt", project)
	c.vmware.sum = registerGauge("selectel_billing_vmware_sum", project)

	return c
}

func (c *balanceCollector) Collect(e *exporter) error {
	log.Println("collect balance metrics")

	res, err := selapi.FetchBalance(e.token)
	if err != nil {
		return err
	}

	c.primary.main.Set(float64(res.Data.Primary.Main))
	c.primary.bonus.Set(float64(res.Data.Primary.Bonus))
	c.primary.vkRub.Set(float64(res.Data.Primary.VkRub))
	c.primary.ref.Set(float64(res.Data.Primary.Ref))
	c.primary.holdMain.Set(float64(res.Data.Primary.Hold.Main))
	c.primary.holdBonus.Set(float64(res.Data.Primary.Hold.Bonus))
	c.primary.holdVkRub.Set(float64(res.Data.Primary.Hold.VkRub))

	c.storage.main.Set(float64(res.Data.Storage.Main))
	c.storage.bonus.Set(float64(res.Data.Storage.Bonus))
	c.storage.vkRub.Set(float64(res.Data.Storage.VkRub))
	c.storage.debt.Set(float64(res.Data.Storage.Debt))
	c.storage.sum.Set(float64(res.Data.Storage.Sum))

	c.vpc.main.Set(float64(res.Data.Vpc.Main))
	c.vpc.bonus.Set(float64(res.Data.Vpc.Bonus))
	c.vpc.vkRub.Set(float64(res.Data.Vpc.VkRub))
	c.vpc.debt.Set(float64(res.Data.Vpc.Debt))
	c.vpc.sum.Set(float64(res.Data.Vpc.Sum))

	c.vmware.main.Set(float64(res.Data.Vmware.Main))
	c.vmware.bonus.Set(float64(res.Data.Vmware.Bonus))
	c.vmware.vkRub.Set(float64(res.Data.Vmware.VkRub))
	c.vmware.debt.Set(float64(res.Data.Vmware.Debt))
	c.vmware.sum.Set(float64(res.Data.Vmware.Sum))

	return nil
}
