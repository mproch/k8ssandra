package unit_test

import (
	"path/filepath"

	"github.com/gruntwork-io/terratest/modules/helm"
	api "github.com/k8ssandra/reaper-operator/api/v1alpha1"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Verify Reaper template", func() {
	var (
		helmChartPath string
		err           error
		reaper        *api.Reaper
	)

	BeforeEach(func() {
		helmChartPath, err = filepath.Abs(chartsPath)
		Expect(err).To(BeNil())
		reaper = &api.Reaper{}
	})

	AfterEach(func() {
		err = nil
	})

	renderTemplate := func(options *helm.Options) {
		renderedOutput := helm.RenderTemplate(
			GinkgoT(), options, helmChartPath, helmReleaseName,
			[]string{"templates/reaper.yaml"},
		)

		helm.UnmarshalK8SYaml(GinkgoT(), renderedOutput, reaper)
	}

	Context("by rendering it with options", func() {
		It("using only default options", func() {
			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
			}

			renderTemplate(options)
			Expect(string(reaper.Spec.ServerConfig.StorageType)).To(Equal("cassandra"))
			Expect(reaper.Kind).To(Equal("Reaper"))
		})

		It("changing datacenter name", func() {
			targetDcName := "reaper-dc"
			options := &helm.Options{
				SetStrValues:   map[string]string{"datacenterName": targetDcName},
				KubectlOptions: defaultKubeCtlOptions,
			}

			renderTemplate(options)
			Expect(reaper.Spec.ServerConfig.CassandraBackend.CassandraDatacenter.Name).To(Equal(targetDcName))
		})

		It("modifying autoscheduling option", func() {
			options := &helm.Options{
				SetStrValues:   map[string]string{"repair.reaper.autoschedule": "true"},
				KubectlOptions: defaultKubeCtlOptions,
			}

			renderTemplate(options)
			Expect(reaper.Spec.ServerConfig.AutoScheduling).ToNot(BeNil())
		})
	})
})
