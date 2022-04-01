package awsplacementgroup

import (
	"context"
	"fmt"
	"os"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	configv1 "github.com/openshift/api/config/v1"
	machinev1 "github.com/openshift/api/machine/v1"

	resourcebuilder "github.com/openshift/machine-api-provider-aws/pkg/resourcebuilder"

	"github.com/golang/mock/gomock"
	mockaws "github.com/openshift/machine-api-provider-aws/pkg/client/mock"

	zaplogfmt "github.com/sykesm/zap-logfmt"
	uzap "go.uber.org/zap"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/envtest/komega"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

var _ = Describe("AWSPlacementGroupReconciler", func() {
	var namespaceName string
	var insfrastructureName string
	var mgrCancel context.CancelFunc
	var mgrDone chan struct{}
	var fakeRecorder *record.FakeRecorder

	BeforeEach(func() {
		By("Setting up a namespace for the test")
		namespace := &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{GenerateName: "pg-test-"}}
		Expect(k8sClient.Create(ctx, namespace)).To(Succeed())
		namespaceName = namespace.GetName()

		By("Setting up a new infrastructure for the test")
		// Create infrastructure object
		infra := resourcebuilder.Infrastructure().WithName("cluster").AsAWS("test", "eu-west-2").Build()
		infraCopy := infra.DeepCopy()
		Expect(k8sClient.Create(ctx, infra)).To(Succeed())

		// Get the infra back for its ResourceVersion
		gotInfra := &configv1.Infrastructure{}
		Expect(k8sClient.Get(ctx, client.ObjectKey{Name: "cluster"}, gotInfra)).To(Succeed())

		// Set the ResourceVersion on the cloned object and Update the Status
		infraCopy.SetResourceVersion(gotInfra.GetResourceVersion())
		Expect(k8sClient.Status().Update(ctx, infraCopy)).To(Succeed())
		insfrastructureName = infra.GetName()

		By("Setting up a manager and controller")
		mgr, err := manager.New(cfg, manager.Options{
			MetricsBindAddress: "0",
			Namespace:          namespaceName,
		})
		Expect(err).ToNot(HaveOccurred(), "Manager should be able to be created")

		// TODO(dam): REMOVE. HACK for Redirecting logging to stdout
		configLog := uzap.NewProductionEncoderConfig()
		logfmtEncoder := zaplogfmt.NewEncoder(configLog)
		logger := zap.New(zap.UseDevMode(true), zap.WriteTo(os.Stdout), zap.Encoder(logfmtEncoder))
		log.SetLogger(logger)

		// Set up a AWS client Mock
		mockCtrl := gomock.NewController(GinkgoT())
		mockAWSClient := mockaws.NewMockClient(mockCtrl)

		reconciler := &AWSPlacementGroupReconciler{
			Client:              mgr.GetClient(),
			Log:                 log.Log,
			ConfigManagedClient: mgr.GetClient(),
			AWSClient:           mockAWSClient,
		}
		Expect(reconciler.SetupWithManager(mgr, controller.Options{})).To(Succeed(), "Reconciler should be able to setup with manager")

		fakeRecorder = record.NewFakeRecorder(1)
		reconciler.recorder = fakeRecorder

		By("Starting the manager")
		var mgrCtx context.Context
		mgrCtx, mgrCancel = context.WithCancel(ctx)
		mgrDone = make(chan struct{})

		go func() {
			defer GinkgoRecover()
			defer close(mgrDone)

			Expect(mgr.Start(mgrCtx)).To(Succeed())
		}()

	})

	AfterEach(func() {
		By("Stopping the manager")
		mgrCancel()
		// Wait for the mgrDone to be closed, which will happen once the mgr has stopped
		<-mgrDone
		Expect(deleteAWSPlacementGroups(k8sClient, namespaceName)).To(Succeed())
		Expect(deleteInfrastructure(k8sClient, insfrastructureName)).To(Succeed())
	})

	Context("when a new Managed AWSPlacementGroup is created", func() {
		var pg *machinev1.AWSPlacementGroup
		// Create the AWSPlacementGroup just before each test so that we can set up
		// various test cases in BeforeEach blocks.
		JustBeforeEach(func() {
			pg = resourcebuilder.AWSPlacementGroup().WithGenerateName("pg-test-").WithNamespace(namespaceName).AsManaged().AsClusterType().Build()
			Expect(k8sClient.Create(ctx, pg)).Should(Succeed())
		})

		It("should add the awsplacementgroup.machine.openshift.io finalizer", func() {
			//spew.Dump(r)
			Eventually(komega.Object(pg), time.Second*10, time.Second*1).Should(HaveField("ObjectMeta.Finalizers", ContainElement(awsPlacementGroupFinalizer)))
		})
	})

	// Context("when a new Unmanaged AWSPlacementGroup is created", func() {
	// 	var pg *machinev1.AWSPlacementGroup
	// 	// Create the AWSPlacementGroup just before each test so that we can set up
	// 	// various test cases in BeforeEach blocks.
	// 	JustBeforeEach(func() {
	// 		pg = resourcebuilder.AWSPlacementGroup().WithGenerateName("pg-test-").WithNamespace(namespaceName).AsUnmanaged().Build()
	// 		Expect(k8sClient.Create(ctx, pg)).Should(Succeed())
	// 	})

	// 	It("should not add the awsplacementgroup.machine.openshift.io finalizer", func() {
	// 		Eventually(komega.Object(pg)).ShouldNot(HaveField("ObjectMeta.Finalizers", ContainElement(awsPlacementGroupFinalizer)))
	// 	})
	// })
})

func deleteAWSPlacementGroups(c client.Client, namespaceName string) error {
	placementGroups := &machinev1.AWSPlacementGroupList{}
	err := c.List(ctx, placementGroups, client.InNamespace(namespaceName))
	if err != nil {
		return err
	}

	for _, ms := range placementGroups.Items {
		err := c.Delete(ctx, &ms)
		if err != nil {
			return err
		}
	}

	Eventually(func() error {
		placementGroups := &machinev1.AWSPlacementGroupList{}
		err := c.List(ctx, placementGroups)
		if err != nil {
			return err
		}
		if len(placementGroups.Items) > 0 {
			return fmt.Errorf("AWSPlacementGroups not deleted")
		}
		return nil
	}, timeout).Should(Succeed())

	return nil
}

func deleteInfrastructure(c client.Client, insfrastructureName string) error {
	// TODO(dam): do this properly
	infra := &configv1.Infrastructure{}
	err := c.Get(ctx, client.ObjectKey{Name: insfrastructureName}, infra)
	if err != nil {
		return err
	}

	if err := c.Delete(ctx, infra); err != nil {
		return err
	}

	Eventually(func() error {
		infras := &configv1.InfrastructureList{}
		err := c.List(ctx, infras)
		if err != nil {
			return err
		}
		if len(infras.Items) > 0 {
			return fmt.Errorf("Infrastructures not deleted")
		}
		return nil
	}, timeout).Should(Succeed())

	return nil
}
