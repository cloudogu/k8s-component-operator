package ecosystem

import (
	"context"
	"fmt"
	"time"

	v1 "github.com/cloudogu/k8s-component-operator/pkg/api/v1"
	"github.com/cloudogu/retry-lib/retry"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

type ComponentInterface interface {
	// Create takes the representation of a component and creates it.  Returns the server's representation of the component, and an error, if there is any.
	Create(ctx context.Context, component *v1.Component, opts metav1.CreateOptions) (*v1.Component, error)

	// Update takes the representation of a component and updates it. Returns the server's representation of the component, and an error, if there is any.
	Update(ctx context.Context, component *v1.Component, opts metav1.UpdateOptions) (*v1.Component, error)

	// UpdateExpectedComponentVersion saves the given version as the expected version in the component and
	// returns the server's representation of the component, and an error, if there is any.
	UpdateExpectedComponentVersion(ctx context.Context, componentName, version string) (*v1.Component, error)

	// UpdateStatus was generated because the type contains a Status member.
	UpdateStatus(ctx context.Context, component *v1.Component, opts metav1.UpdateOptions) (*v1.Component, error)

	// UpdateStatusInstalling sets the status of the component to "installing".
	UpdateStatusInstalling(ctx context.Context, component *v1.Component) (*v1.Component, error)

	// UpdateStatusInstalled sets the status of the component to "installed".
	UpdateStatusInstalled(ctx context.Context, component *v1.Component) (*v1.Component, error)

	// UpdateStatusUpgrading sets the status of the component to "upgrading".
	UpdateStatusUpgrading(ctx context.Context, component *v1.Component) (*v1.Component, error)

	// UpdateStatusDeleting sets the status of the component to "deleting".
	UpdateStatusDeleting(ctx context.Context, component *v1.Component) (*v1.Component, error)

	// UpdateStatusNotInstalled sets the status of the component to "".
	UpdateStatusNotInstalled(ctx context.Context, component *v1.Component) (*v1.Component, error)

	// Delete takes name of the component and deletes it. Returns an error if one occurs.
	Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error

	// DeleteCollection deletes a collection of objects.
	DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error

	// Get takes name of the component, and returns the corresponding component object, and an error if there is any.
	Get(ctx context.Context, name string, opts metav1.GetOptions) (*v1.Component, error)

	// List takes label and field selectors, and returns the list of Components that match those selectors.
	List(ctx context.Context, opts metav1.ListOptions) (*v1.ComponentList, error)

	// Watch returns a watch.Interface that watches the requested components.
	Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error)

	// Patch applies the patch and returns the patched component.
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.Component, err error)

	// AddFinalizer adds the given finalizer to the component.
	AddFinalizer(ctx context.Context, component *v1.Component, finalizer string) (*v1.Component, error)

	// RemoveFinalizer removes the given finalizer to the component.
	RemoveFinalizer(ctx context.Context, component *v1.Component, finalizer string) (*v1.Component, error)
}

type componentClient struct {
	client rest.Interface
	ns     string
}

// UpdateStatusInstalling sets the status of the component to "installing".
func (client *componentClient) UpdateStatusInstalling(ctx context.Context, component *v1.Component) (*v1.Component, error) {
	return client.updateStatusWithRetry(ctx, component, v1.ComponentStatusInstalling)
}

// UpdateStatusInstalled sets the status of the component to "installed".
func (client *componentClient) UpdateStatusInstalled(ctx context.Context, component *v1.Component) (*v1.Component, error) {
	return client.updateStatusWithRetry(ctx, component, v1.ComponentStatusInstalled)
}

// UpdateStatusUpgrading sets the status of the component to "upgrading".
func (client *componentClient) UpdateStatusUpgrading(ctx context.Context, component *v1.Component) (*v1.Component, error) {
	return client.updateStatusWithRetry(ctx, component, v1.ComponentStatusUpgrading)
}

// UpdateStatusDeleting sets the status of the component to "deleting".
func (client *componentClient) UpdateStatusDeleting(ctx context.Context, component *v1.Component) (*v1.Component, error) {
	return client.updateStatusWithRetry(ctx, component, v1.ComponentStatusDeleting)
}

// UpdateStatusNotInstalled sets the status of the component to "".
func (client *componentClient) UpdateStatusNotInstalled(ctx context.Context, component *v1.Component) (*v1.Component, error) {
	return client.updateStatusWithRetry(ctx, component, v1.ComponentStatusNotInstalled)
}

func (client *componentClient) updateStatusWithRetry(ctx context.Context, component *v1.Component, targetStatus string) (*v1.Component, error) {
	var resultComponent *v1.Component
	err := retry.OnConflict(func() error {
		updatedComponent, err := client.Get(ctx, component.GetName(), metav1.GetOptions{})
		if err != nil {
			return err
		}

		// do not overwrite the whole status, so we do not lose other values from the Status object
		// esp. a potentially set requeue time
		updatedComponent.Status.Status = targetStatus
		// Make sure the updated values are not lost on update
		updatedComponent.Spec.MappedValuesYamlOverwrite = component.Spec.MappedValuesYamlOverwrite
		resultComponent, err = client.UpdateStatus(ctx, updatedComponent, metav1.UpdateOptions{})
		return err
	})

	return resultComponent, err
}

// AddFinalizer adds the given finalizer to the component.
func (client *componentClient) AddFinalizer(ctx context.Context, component *v1.Component, finalizer string) (*v1.Component, error) {
	controllerutil.AddFinalizer(component, finalizer)
	result, err := client.Update(ctx, component, metav1.UpdateOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to add finalizer %s to component: %w", finalizer, err)
	}

	return result, nil
}

// RemoveFinalizer removes the given finalizer to the component.
func (client *componentClient) RemoveFinalizer(ctx context.Context, component *v1.Component, finalizer string) (*v1.Component, error) {
	controllerutil.RemoveFinalizer(component, finalizer)
	result, err := client.Update(ctx, component, metav1.UpdateOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to remove finalizer %s from component: %w", finalizer, err)
	}

	return result, err
}

// Get takes name of the component, and returns the corresponding component object, and an error if there is any.
func (client *componentClient) Get(ctx context.Context, name string, options metav1.GetOptions) (result *v1.Component, err error) {
	result = &v1.Component{}
	err = client.client.Get().
		Namespace(client.ns).
		Resource("components").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of Components that match those selectors.
func (client *componentClient) List(ctx context.Context, opts metav1.ListOptions) (result *v1.ComponentList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1.ComponentList{}
	err = client.client.Get().
		Namespace(client.ns).
		Resource("components").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested components.
func (client *componentClient) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return client.client.Get().
		Namespace(client.ns).
		Resource("components").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a component and creates it.  Returns the server's representation of the component, and an error, if there is any.
func (client *componentClient) Create(ctx context.Context, component *v1.Component, opts metav1.CreateOptions) (result *v1.Component, err error) {
	result = &v1.Component{}
	err = client.client.Post().
		Namespace(client.ns).
		Resource("components").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(component).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a component and updates it. Returns the server's representation of the component, and an error, if there is any.
func (client *componentClient) Update(ctx context.Context, component *v1.Component, opts metav1.UpdateOptions) (result *v1.Component, err error) {
	result = &v1.Component{}
	err = client.client.Put().
		Namespace(client.ns).
		Resource("components").
		Name(component.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(component).
		Do(ctx).
		Into(result)
	return
}

func (client *componentClient) UpdateExpectedComponentVersion(ctx context.Context, componentName, version string) (*v1.Component, error) {
	var updatedComponent *v1.Component
	err := retry.OnConflict(func() error {
		retryComponent, err := client.Get(ctx, componentName, metav1.GetOptions{})
		if err != nil {
			return fmt.Errorf("failed to get component %q for update: %w", componentName, err)
		}

		retryComponent.Spec.Version = version
		retryComponent, err = client.Update(ctx, retryComponent, metav1.UpdateOptions{})
		if err != nil {
			return err
		}
		updatedComponent = retryComponent
		return nil
	})
	if err != nil {
		return updatedComponent, fmt.Errorf("failed to update version in component %q: %w", componentName, err)
	}

	return updatedComponent, nil
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (client *componentClient) UpdateStatus(ctx context.Context, component *v1.Component, opts metav1.UpdateOptions) (result *v1.Component, err error) {
	result = &v1.Component{}
	err = client.client.Put().
		Namespace(client.ns).
		Resource("components").
		Name(component.Name).
		SubResource("status").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(component).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the component and deletes it. Returns an error if one occurs.
func (client *componentClient) Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	return client.client.Delete().
		Namespace(client.ns).
		Resource("components").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (client *componentClient) DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return client.client.Delete().
		Namespace(client.ns).
		Resource("components").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched component.
func (client *componentClient) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.Component, err error) {
	result = &v1.Component{}
	err = client.client.Patch(pt).
		Namespace(client.ns).
		Resource("components").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}
