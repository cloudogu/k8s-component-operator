package ecosystem

import (
	"context"
	"fmt"
	v1 "github.com/cloudogu/k8s-component-operator/pkg/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"time"
)

type ComponentInterface interface {
	Create(ctx context.Context, component *v1.Component, opts metav1.CreateOptions) (*v1.Component, error)
	Update(ctx context.Context, component *v1.Component, opts metav1.UpdateOptions) (*v1.Component, error)
	UpdateStatus(ctx context.Context, component *v1.Component, opts metav1.UpdateOptions) (*v1.Component, error)
	UpdateStatusInstalling(ctx context.Context, component *v1.Component) (*v1.Component, error)
	UpdateStatusInstalled(ctx context.Context, component *v1.Component) (*v1.Component, error)
	UpdateStatusUpgrading(ctx context.Context, component *v1.Component) (*v1.Component, error)
	UpdateStatusDeleting(ctx context.Context, component *v1.Component) (*v1.Component, error)
	Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error
	Get(ctx context.Context, name string, opts metav1.GetOptions) (*v1.Component, error)
	List(ctx context.Context, opts metav1.ListOptions) (*v1.ComponentList, error)
	Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.Component, err error)
	AddFinalizer(ctx context.Context, component *v1.Component, finalizer string) (*v1.Component, error)
	RemoveFinalizer(ctx context.Context, component *v1.Component, finalizer string) (*v1.Component, error)
}

type componentClient struct {
	client rest.Interface
	ns     string
}

// UpdateStatusInstalling sets the status of the component to "installing".
func (d *componentClient) UpdateStatusInstalling(ctx context.Context, component *v1.Component) (*v1.Component, error) {
	component.Status = v1.ComponentStatus{Status: v1.ComponentStatusInstalling}
	return d.UpdateStatus(ctx, component, metav1.UpdateOptions{})
}

// UpdateStatusInstalled sets the status of the component to "installed".
func (d *componentClient) UpdateStatusInstalled(ctx context.Context, component *v1.Component) (*v1.Component, error) {
	component.Status = v1.ComponentStatus{Status: v1.ComponentStatusInstalled}
	return d.UpdateStatus(ctx, component, metav1.UpdateOptions{})
}

// UpdateStatusUpgrading sets the status of the component to "upgrading".
func (d *componentClient) UpdateStatusUpgrading(ctx context.Context, component *v1.Component) (*v1.Component, error) {
	component.Status = v1.ComponentStatus{Status: v1.ComponentStatusUpgrading}
	return d.UpdateStatus(ctx, component, metav1.UpdateOptions{})
}

// UpdateStatusDeleting sets the status of the component to "deleting".
func (d *componentClient) UpdateStatusDeleting(ctx context.Context, component *v1.Component) (*v1.Component, error) {
	component.Status = v1.ComponentStatus{Status: v1.ComponentStatusDeleting}
	return d.UpdateStatus(ctx, component, metav1.UpdateOptions{})
}

// AddFinalizer adds the given finalizer to the component.
func (d *componentClient) AddFinalizer(ctx context.Context, component *v1.Component, finalizer string) (*v1.Component, error) {
	controllerutil.AddFinalizer(component, finalizer)
	result, err := d.Update(ctx, component, metav1.UpdateOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to add finalizer %s to component: %w", finalizer, err)
	}

	return result, nil
}

// RemoveFinalizer removes the given finalizer to the component.
func (d *componentClient) RemoveFinalizer(ctx context.Context, component *v1.Component, finalizer string) (*v1.Component, error) {
	controllerutil.RemoveFinalizer(component, finalizer)
	result, err := d.Update(ctx, component, metav1.UpdateOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to remove finalizer %s to component: %w", finalizer, err)
	}

	return result, err
}

// Get takes name of the component, and returns the corresponding component object, and an error if there is any.
func (d *componentClient) Get(ctx context.Context, name string, options metav1.GetOptions) (result *v1.Component, err error) {
	result = &v1.Component{}
	err = d.client.Get().
		Namespace(d.ns).
		Resource("components").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of Components that match those selectors.
func (d *componentClient) List(ctx context.Context, opts metav1.ListOptions) (result *v1.ComponentList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1.ComponentList{}
	err = d.client.Get().
		Namespace(d.ns).
		Resource("components").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested components.
func (d *componentClient) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return d.client.Get().
		Namespace(d.ns).
		Resource("components").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a component and creates it.  Returns the server's representation of the component, and an error, if there is any.
func (d *componentClient) Create(ctx context.Context, component *v1.Component, opts metav1.CreateOptions) (result *v1.Component, err error) {
	result = &v1.Component{}
	err = d.client.Post().
		Namespace(d.ns).
		Resource("components").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(component).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a component and updates it. Returns the server's representation of the component, and an error, if there is any.
func (d *componentClient) Update(ctx context.Context, component *v1.Component, opts metav1.UpdateOptions) (result *v1.Component, err error) {
	result = &v1.Component{}
	err = d.client.Put().
		Namespace(d.ns).
		Resource("components").
		Name(component.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(component).
		Do(ctx).
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (d *componentClient) UpdateStatus(ctx context.Context, component *v1.Component, opts metav1.UpdateOptions) (result *v1.Component, err error) {
	result = &v1.Component{}
	err = d.client.Put().
		Namespace(d.ns).
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
func (d *componentClient) Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	return d.client.Delete().
		Namespace(d.ns).
		Resource("components").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (d *componentClient) DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return d.client.Delete().
		Namespace(d.ns).
		Resource("components").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched component.
func (d *componentClient) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.Component, err error) {
	result = &v1.Component{}
	err = d.client.Patch(pt).
		Namespace(d.ns).
		Resource("components").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}
