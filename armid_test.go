package armid

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestResourceId_String(t *testing.T) {
	cases := []struct {
		name   string
		input  ResourceId
		expect string
	}{
		{
			name:   "Tenant",
			input:  &TenantId{},
			expect: "/",
		},
		{
			name:   "Subscription",
			input:  &SubscriptionId{Id: "sub1"},
			expect: "/subscriptions/sub1",
		},
		{
			name:   "Resource Group",
			input:  &ResourceGroup{SubscriptionId: "sub1", Name: "rg1"},
			expect: "/subscriptions/sub1/resourceGroups/rg1",
		},
		{
			name:   "Management Group",
			input:  &ManagementGroup{Name: "mg1"},
			expect: "/providers/Microsoft.Management/managementGroups/mg1",
		},
		{
			name: "Scoped Resource under tenant",
			input: &ScopedResourceId{
				AttrParentScope: &TenantId{},
				AttrProvider:    "Microsoft.Foo",
				AttrTypes:       []string{"foos", "bars"},
				AttrNames:       []string{"foo1", "bar1"},
			},
			expect: "/providers/Microsoft.Foo/foos/foo1/bars/bar1",
		},
		{
			name: "Scoped Resource under subscription",
			input: &ScopedResourceId{
				AttrParentScope: &SubscriptionId{Id: "sub1"},
				AttrProvider:    "Microsoft.Foo",
				AttrTypes:       []string{"foos", "bars"},
				AttrNames:       []string{"foo1", "bar1"},
			},
			expect: "/subscriptions/sub1/providers/Microsoft.Foo/foos/foo1/bars/bar1",
		},
		{
			name: "Scoped Resource under resource group",
			input: &ScopedResourceId{
				AttrParentScope: &ResourceGroup{SubscriptionId: "sub1", Name: "rg1"},
				AttrProvider:    "Microsoft.Foo",
				AttrTypes:       []string{"foos", "bars"},
				AttrNames:       []string{"foo1", "bar1"},
			},
			expect: "/subscriptions/sub1/resourceGroups/rg1/providers/Microsoft.Foo/foos/foo1/bars/bar1",
		},
		{
			name: "Scoped Resource under management group",
			input: &ScopedResourceId{
				AttrParentScope: &ManagementGroup{Name: "mg1"},
				AttrProvider:    "Microsoft.Foo",
				AttrTypes:       []string{"foos", "bars"},
				AttrNames:       []string{"foo1", "bar1"},
			},
			expect: "/providers/Microsoft.Management/managementGroups/mg1/providers/Microsoft.Foo/foos/foo1/bars/bar1",
		},
		{
			name: "Scoped Resource under another scoped resource which under tenant",
			input: &ScopedResourceId{
				AttrParentScope: &ScopedResourceId{
					AttrParentScope: &TenantId{},
					AttrProvider:    "Microsoft.Foo",
					AttrTypes:       []string{"foos", "bars"},
					AttrNames:       []string{"foo1", "bar1"},
				},
				AttrProvider: "Microsoft.Baz",
				AttrTypes:    []string{"bazs"},
				AttrNames:    []string{"baz1"},
			},
			expect: "/providers/Microsoft.Foo/foos/foo1/bars/bar1/providers/Microsoft.Baz/bazs/baz1",
		},
		{
			name: "Scoped Resource under another scoped resource which under subscription",
			input: &ScopedResourceId{
				AttrParentScope: &ScopedResourceId{
					AttrParentScope: &SubscriptionId{Id: "sub1"},
					AttrProvider:    "Microsoft.Foo",
					AttrTypes:       []string{"foos", "bars"},
					AttrNames:       []string{"foo1", "bar1"},
				},
				AttrProvider: "Microsoft.Baz",
				AttrTypes:    []string{"bazs"},
				AttrNames:    []string{"baz1"},
			},
			expect: "/subscriptions/sub1/providers/Microsoft.Foo/foos/foo1/bars/bar1/providers/Microsoft.Baz/bazs/baz1",
		},
		{
			name: "Scoped Resource under another scoped resource which under resource group",
			input: &ScopedResourceId{
				AttrParentScope: &ScopedResourceId{
					AttrParentScope: &ResourceGroup{SubscriptionId: "sub1", Name: "rg1"},
					AttrProvider:    "Microsoft.Foo",
					AttrTypes:       []string{"foos", "bars"},
					AttrNames:       []string{"foo1", "bar1"},
				},
				AttrProvider: "Microsoft.Baz",
				AttrTypes:    []string{"bazs"},
				AttrNames:    []string{"baz1"},
			},
			expect: "/subscriptions/sub1/resourceGroups/rg1/providers/Microsoft.Foo/foos/foo1/bars/bar1/providers/Microsoft.Baz/bazs/baz1",
		},
		{
			name: "Scoped Resource under another scoped resource which under management group",
			input: &ScopedResourceId{
				AttrParentScope: &ScopedResourceId{
					AttrParentScope: &ManagementGroup{Name: "mg1"},
					AttrProvider:    "Microsoft.Foo",
					AttrTypes:       []string{"foos", "bars"},
					AttrNames:       []string{"foo1", "bar1"},
				},
				AttrProvider: "Microsoft.Baz",
				AttrTypes:    []string{"bazs"},
				AttrNames:    []string{"baz1"},
			},
			expect: "/providers/Microsoft.Management/managementGroups/mg1/providers/Microsoft.Foo/foos/foo1/bars/bar1/providers/Microsoft.Baz/bazs/baz1",
		},
		{
			name: "Subscription scope level resource",
			input: &SubscriptionId{
				Id:        "sub1",
				AttrTypes: []string{"tagNames", "tagValues"},
				AttrNames: []string{"name1", "value1"},
			},
			expect: "/subscriptions/sub1/tagNames/name1/tagValues/value1",
		},
		{
			name: "Resource group scope level resource",
			input: &ResourceGroup{
				SubscriptionId: "sub1",
				Name:           "rg1",
				AttrTypes:      []string{"deployments"},
				AttrNames:      []string{"deploy1"},
			},
			expect: "/subscriptions/sub1/resourceGroups/rg1/deployments/deploy1",
		},
		{
			name: "Mgmt group scope level resource",
			input: &ManagementGroup{
				Name:      "group1",
				AttrTypes: []string{"foos"},
				AttrNames: []string{"foo1"},
			},
			expect: "/providers/Microsoft.Management/managementGroups/group1/foos/foo1",
		},
		{
			name: `RP level resource`,
			input: &ScopedResourceId{
				AttrParentScope: &TenantId{},
				AttrProvider:    "Microsoft.Foo",
				AttrTypes:       []string{},
				AttrNames:       []string{},
			},
			expect: "/providers/Microsoft.Foo",
		},
		{
			name: `RP level resource under another RP level resource`,
			input: &ScopedResourceId{
				AttrParentScope: &ScopedResourceId{
					AttrParentScope: &TenantId{},
					AttrProvider:    "Microsoft.Foo",
					AttrTypes:       []string{},
					AttrNames:       []string{},
				},
				AttrProvider: "Microsoft.Bar",
			},
			expect: "/providers/Microsoft.Foo/providers/Microsoft.Bar",
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.expect, tt.input.String())
		})
	}
}

func TestResourceId_TypeString(t *testing.T) {
	cases := []struct {
		name   string
		input  ResourceId
		expect string
	}{
		{
			name:   "Tenant",
			input:  &TenantId{},
			expect: "",
		},
		{
			name:   "Subscription",
			input:  &SubscriptionId{Id: "sub1"},
			expect: "Microsoft.Resources/subscriptions",
		},
		{
			name:   "Resource Group",
			input:  &ResourceGroup{SubscriptionId: "sub1", Name: "rg1"},
			expect: "Microsoft.Resources/subscriptions/resourceGroups",
		},
		{
			name:   "Management Group",
			input:  &ManagementGroup{Name: "mg1"},
			expect: "Microsoft.Management/managementGroups",
		},
		{
			name: "Scoped Resource under tenant",
			input: &ScopedResourceId{
				AttrParentScope: &TenantId{},
				AttrProvider:    "Microsoft.Foo",
				AttrTypes:       []string{"foos", "bars"},
				AttrNames:       []string{"foo1", "bar1"},
			},
			expect: "Microsoft.Foo/foos/bars",
		},
		{
			name: "Scoped Resource under subscription",
			input: &ScopedResourceId{
				AttrParentScope: &SubscriptionId{Id: "sub1"},
				AttrProvider:    "Microsoft.Foo",
				AttrTypes:       []string{"foos", "bars"},
				AttrNames:       []string{"foo1", "bar1"},
			},
			expect: "Microsoft.Foo/foos/bars",
		},
		{
			name: "Scoped Resource under resource group",
			input: &ScopedResourceId{
				AttrParentScope: &ResourceGroup{SubscriptionId: "sub1", Name: "rg1"},
				AttrProvider:    "Microsoft.Foo",
				AttrTypes:       []string{"foos", "bars"},
				AttrNames:       []string{"foo1", "bar1"},
			},
			expect: "Microsoft.Foo/foos/bars",
		},
		{
			name: "Scoped Resource under management group",
			input: &ScopedResourceId{
				AttrParentScope: &ManagementGroup{Name: "mg1"},
				AttrProvider:    "Microsoft.Foo",
				AttrTypes:       []string{"foos", "bars"},
				AttrNames:       []string{"foo1", "bar1"},
			},
			expect: "Microsoft.Foo/foos/bars",
		},
		{
			name: "Scoped Resource under another scoped resource which under tenant",
			input: &ScopedResourceId{
				AttrParentScope: &ScopedResourceId{
					AttrParentScope: &TenantId{},
					AttrProvider:    "Microsoft.Foo",
					AttrTypes:       []string{"foos", "bars"},
					AttrNames:       []string{"foo1", "bar1"},
				},
				AttrProvider: "Microsoft.Baz",
				AttrTypes:    []string{"bazs"},
				AttrNames:    []string{"baz1"},
			},
			expect: "Microsoft.Baz/bazs",
		},
		{
			name: "Subscription scope level resource",
			input: &SubscriptionId{
				Id:        "sub1",
				AttrTypes: []string{"tagNames", "tagValues"},
				AttrNames: []string{"name1", "value1"},
			},
			expect: "Microsoft.Resources/subscriptions/tagNames/tagValues",
		},
		{
			name: "Resource group scope level resource",
			input: &ResourceGroup{
				SubscriptionId: "sub1",
				Name:           "rg1",
				AttrTypes:      []string{"deployments"},
				AttrNames:      []string{"deploy1"},
			},
			expect: "Microsoft.Resources/subscriptions/resourceGroups/deployments",
		},
		{
			name: "Mgmt group scope level resource",
			input: &ManagementGroup{
				Name:      "group1",
				AttrTypes: []string{"foos"},
				AttrNames: []string{"foo1"},
			},
			expect: "Microsoft.Management/managementGroups/foos",
		},
		{
			name: `RP level resource`,
			input: &ScopedResourceId{
				AttrParentScope: &TenantId{},
				AttrProvider:    "Microsoft.Foo",
				AttrTypes:       []string{},
				AttrNames:       []string{},
			},
			expect: "Microsoft.Foo",
		},
		{
			name: `RP level resource under another RP level resource`,
			input: &ScopedResourceId{
				AttrParentScope: &ScopedResourceId{
					AttrParentScope: &TenantId{},
					AttrProvider:    "Microsoft.Foo",
					AttrTypes:       []string{},
					AttrNames:       []string{},
				},
				AttrProvider: "Microsoft.Bar",
			},
			expect: "Microsoft.Bar",
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.expect, tt.input.TypeString())
		})
	}
}

func TestResourceId_RootScope(t *testing.T) {
	cases := []struct {
		name   string
		input  ResourceId
		expect RootScope
	}{
		{
			name:   "Tenant",
			input:  &TenantId{},
			expect: &TenantId{},
		},
		{
			name:   "Subscription",
			input:  &SubscriptionId{Id: "sub1"},
			expect: &SubscriptionId{Id: "sub1"},
		},
		{
			name:   "Resource Group",
			input:  &ResourceGroup{SubscriptionId: "sub1", Name: "rg1"},
			expect: &ResourceGroup{SubscriptionId: "sub1", Name: "rg1"},
		},
		{
			name:   "Management Group",
			input:  &ManagementGroup{Name: "mg1"},
			expect: &ManagementGroup{Name: "mg1"},
		},
		{
			name: "Root Scoped Resource under tenant",
			input: &ScopedResourceId{
				AttrParentScope: &TenantId{},
				AttrProvider:    "Microsoft.Foo",
				AttrTypes:       []string{"foos"},
				AttrNames:       []string{"foo1"},
			},
			expect: &TenantId{},
		},
		{
			name: "Child Scoped Resource under tenant",
			input: &ScopedResourceId{
				AttrParentScope: &TenantId{},
				AttrProvider:    "Microsoft.Foo",
				AttrTypes:       []string{"foos", "bars"},
				AttrNames:       []string{"foo1", "bar1"},
			},
			expect: &TenantId{},
		},
		{
			name: "Child Scoped Resource under resource group",
			input: &ScopedResourceId{
				AttrParentScope: &ResourceGroup{SubscriptionId: "sub1", Name: "rg1"},
				AttrProvider:    "Microsoft.Foo",
				AttrTypes:       []string{"foos", "bars"},
				AttrNames:       []string{"foo1", "bar1"},
			},
			expect: &ResourceGroup{SubscriptionId: "sub1", Name: "rg1"},
		},
		{
			name: "Subscription scope level resource",
			input: &SubscriptionId{
				Id:        "sub1",
				AttrTypes: []string{"tagNames", "tagValues"},
				AttrNames: []string{"name1", "value1"},
			},
			expect: &SubscriptionId{
				Id:        "sub1",
				AttrTypes: []string{"tagNames", "tagValues"},
				AttrNames: []string{"name1", "value1"},
			},
		},
		{
			name: "Resource group scope level resource",
			input: &ResourceGroup{
				SubscriptionId: "sub1",
				Name:           "rg1",
				AttrTypes:      []string{"deployments"},
				AttrNames:      []string{"deploy1"},
			},
			expect: &ResourceGroup{
				SubscriptionId: "sub1",
				Name:           "rg1",
				AttrTypes:      []string{"deployments"},
				AttrNames:      []string{"deploy1"},
			},
		},
		{
			name: "Mgmt group scope level resource",
			input: &ManagementGroup{
				Name:      "group1",
				AttrTypes: []string{"foos"},
				AttrNames: []string{"foo1"},
			},
			expect: &ManagementGroup{
				Name:      "group1",
				AttrTypes: []string{"foos"},
				AttrNames: []string{"foo1"},
			},
		},
		{
			name: `RP level resource`,
			input: &ScopedResourceId{
				AttrParentScope: &TenantId{},
				AttrProvider:    "Microsoft.Foo",
			},
			expect: &TenantId{},
		},
		{
			name: `RP level resource under another RP level resource`,
			input: &ScopedResourceId{
				AttrParentScope: &ScopedResourceId{
					AttrParentScope: &TenantId{},
					AttrProvider:    "Microsoft.Foo",
				},
				AttrProvider: "Microsoft.Bar",
			},
			expect: &TenantId{},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.expect, tt.input.RootScope())
		})
	}
}

func TestResourceId_Parent(t *testing.T) {
	cases := []struct {
		name   string
		input  ResourceId
		expect ResourceId
	}{
		{
			name:   "Tenant",
			input:  &TenantId{},
			expect: nil,
		},
		{
			name:   "Subscription",
			input:  &SubscriptionId{Id: "sub1"},
			expect: nil,
		},
		{
			name:   "Resource Group",
			input:  &ResourceGroup{SubscriptionId: "sub1", Name: "rg1"},
			expect: nil,
		},
		{
			name:   "Management Group",
			input:  &ManagementGroup{Name: "mg1"},
			expect: nil,
		},
		{
			name: "Root Scoped Resource under tenant",
			input: &ScopedResourceId{
				AttrParentScope: &TenantId{},
				AttrProvider:    "Microsoft.Foo",
				AttrTypes:       []string{"foos"},
				AttrNames:       []string{"foo1"},
			},
			expect: &ScopedResourceId{
				AttrParentScope: &TenantId{},
				AttrProvider:    "Microsoft.Foo",
				AttrTypes:       []string{},
				AttrNames:       []string{},
			},
		},
		{
			name: "Child Scoped Resource under tenant",
			input: &ScopedResourceId{
				AttrParentScope: &TenantId{},
				AttrProvider:    "Microsoft.Foo",
				AttrTypes:       []string{"foos", "bars"},
				AttrNames:       []string{"foo1", "bar1"},
			},
			expect: &ScopedResourceId{
				AttrParentScope: &TenantId{},
				AttrProvider:    "Microsoft.Foo",
				AttrTypes:       []string{"foos"},
				AttrNames:       []string{"foo1"},
			},
		},
		{
			name: "Subscription scope level resource",
			input: &SubscriptionId{
				Id:        "sub1",
				AttrTypes: []string{"tagNames", "tagValues"},
				AttrNames: []string{"name1", "value1"},
			},
			expect: &SubscriptionId{
				Id:        "sub1",
				AttrTypes: []string{"tagNames"},
				AttrNames: []string{"name1"},
			},
		},
		{
			name: "Resource group scope level resource",
			input: &ResourceGroup{
				SubscriptionId: "sub1",
				Name:           "rg1",
				AttrTypes:      []string{"deployments"},
				AttrNames:      []string{"deploy1"},
			},
			expect: &ResourceGroup{
				SubscriptionId: "sub1",
				Name:           "rg1",
				AttrTypes:      []string{},
				AttrNames:      []string{},
			},
		},
		{
			name: "Mgmt group scope level resource",
			input: &ManagementGroup{
				Name:      "group1",
				AttrTypes: []string{"foos"},
				AttrNames: []string{"foo1"},
			},
			expect: &ManagementGroup{
				Name:      "group1",
				AttrTypes: []string{},
				AttrNames: []string{},
			},
		},
		{
			name: `RP level resource`,
			input: &ScopedResourceId{
				AttrParentScope: &TenantId{},
				AttrProvider:    "Microsoft.Foo",
			},
			expect: nil,
		},
		{
			name: `RP level resource under another RP level resource`,
			input: &ScopedResourceId{
				AttrParentScope: &ScopedResourceId{
					AttrParentScope: &TenantId{},
					AttrProvider:    "Microsoft.Foo",
				},
				AttrProvider: "Microsoft.Bar",
			},
			expect: nil,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.expect, tt.input.Parent())
		})
	}
}

func TestResourceId_Equal(t *testing.T) {
	cases := []struct {
		name   string
		id     ResourceId
		oid    ResourceId
		expect bool
	}{
		{
			name:   "Tenant equals to Tenant",
			id:     &TenantId{},
			oid:    &TenantId{},
			expect: true,
		},
		{
			name:   "Tenant not equals to Subscription",
			id:     &TenantId{},
			oid:    &SubscriptionId{Id: "sub1"},
			expect: false,
		},
		{
			name:   "Subscription equals to subscription with same id",
			id:     &SubscriptionId{Id: "sub1"},
			oid:    &SubscriptionId{Id: "sub1"},
			expect: true,
		},
		{
			name:   "Subscription not equals to subscription with different id",
			id:     &SubscriptionId{Id: "sub1"},
			oid:    &SubscriptionId{Id: "sub2"},
			expect: false,
		},
		{
			name:   "Resource Group equals to Resource Group with same subscription id and resource group name",
			id:     &ResourceGroup{SubscriptionId: "sub1", Name: "rg1"},
			oid:    &ResourceGroup{SubscriptionId: "sub1", Name: "rg1"},
			expect: true,
		},
		{
			name:   "Resource Group not equals to Resource Group with different subscription id and resource group name",
			id:     &ResourceGroup{SubscriptionId: "sub1", Name: "rg1"},
			oid:    &ResourceGroup{SubscriptionId: "sub2", Name: "rg2"},
			expect: false,
		},
		{
			name:   "Management Group equals to Management Group with same name",
			id:     &ManagementGroup{Name: "mg1"},
			oid:    &ManagementGroup{Name: "mg1"},
			expect: true,
		},
		{
			name:   "Management Group not equals to Management Group with different name",
			id:     &ManagementGroup{Name: "mg1"},
			oid:    &ManagementGroup{Name: "mg2"},
			expect: false,
		},
		{
			name: "Root Scoped Resource under tenant equals to itself",
			id: &ScopedResourceId{
				AttrParentScope: &TenantId{},
				AttrProvider:    "Microsoft.Foo",
				AttrTypes:       []string{"foos"},
				AttrNames:       []string{"foo1"},
			},
			oid: &ScopedResourceId{
				AttrParentScope: &TenantId{},
				AttrProvider:    "Microsoft.Foo",
				AttrTypes:       []string{"foos"},
				AttrNames:       []string{"foo1"},
			},
			expect: true,
		},
		{
			name: "Root Scoped Resource under tenant equals to itself with different casing",
			id: &ScopedResourceId{
				AttrParentScope: &TenantId{},
				AttrProvider:    "MICROSOFT.FOO",
				AttrTypes:       []string{"FOOS"},
				AttrNames:       []string{"FOO1"},
			},
			oid: &ScopedResourceId{
				AttrParentScope: &TenantId{},
				AttrProvider:    "Microsoft.Foo",
				AttrTypes:       []string{"foos"},
				AttrNames:       []string{"foo1"},
			},
			expect: true,
		},
		{
			name: "Root Scoped Resource under tenant not equals to different resource id",
			id: &ScopedResourceId{
				AttrParentScope: &TenantId{},
				AttrProvider:    "Microsoft.Foo",
				AttrTypes:       []string{"foos"},
				AttrNames:       []string{"foo1"},
			},
			oid: &ScopedResourceId{
				AttrParentScope: &TenantId{},
				AttrProvider:    "Microsoft.Foo",
				AttrTypes:       []string{"bars"},
				AttrNames:       []string{"bar1"},
			},
			expect: false,
		},
		{
			name: "Subscription scope level resource",
			id: &SubscriptionId{
				Id:        "sub1",
				AttrTypes: []string{"tagNames", "tagValues"},
				AttrNames: []string{"name1", "value1"},
			},
			oid: &SubscriptionId{
				Id:        "sub1",
				AttrTypes: []string{"tagNames", "tagValues"},
				AttrNames: []string{"name1", "value1"},
			},
			expect: true,
		},
		{
			name: "Resource group scope level resource",
			id: &ResourceGroup{
				SubscriptionId: "sub1",
				Name:           "rg1",
				AttrTypes:      []string{"deployments"},
				AttrNames:      []string{"deploy1"},
			},
			oid: &ResourceGroup{
				SubscriptionId: "sub1",
				Name:           "rg1",
				AttrTypes:      []string{"deployments"},
				AttrNames:      []string{"deploy1"},
			},
			expect: true,
		},
		{
			name: "Mgmt group scope level resource",
			id: &ManagementGroup{
				Name:      "group1",
				AttrTypes: []string{"foos"},
				AttrNames: []string{"foo1"},
			},
			oid: &ManagementGroup{
				Name:      "group1",
				AttrTypes: []string{"foos"},
				AttrNames: []string{"foo1"},
			},
			expect: true,
		},
		{
			name: `RP level resource`,
			id: &ScopedResourceId{
				AttrParentScope: &TenantId{},
				AttrProvider:    "Microsoft.Foo",
			},
			oid: &ScopedResourceId{
				AttrParentScope: &TenantId{},
				AttrProvider:    "Microsoft.Foo",
			},
			expect: true,
		},
		{
			name: `RP level resource under another RP level resource`,
			id: &ScopedResourceId{
				AttrParentScope: &ScopedResourceId{
					AttrParentScope: &TenantId{},
					AttrProvider:    "Microsoft.Foo",
				},
				AttrProvider: "Microsoft.Bar",
			},
			oid: &ScopedResourceId{
				AttrParentScope: &ScopedResourceId{
					AttrParentScope: &TenantId{},
					AttrProvider:    "Microsoft.Foo",
				},
				AttrProvider: "Microsoft.Bar",
			},
			expect: true,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.expect, tt.id.Equal(tt.oid))
		})
	}
}

func TestResourceId_ScopeEqual(t *testing.T) {
	cases := []struct {
		name   string
		id     ResourceId
		oid    ResourceId
		expect bool
	}{
		{
			name:   "Tenant equals scope to Tenant",
			id:     &TenantId{},
			oid:    &TenantId{},
			expect: true,
		},
		{
			name:   "Tenant not equals scope to Subscription",
			id:     &TenantId{},
			oid:    &SubscriptionId{Id: "sub1"},
			expect: false,
		},
		{
			name:   "Subscription equals scope to subscription with different id",
			id:     &SubscriptionId{Id: "sub1"},
			oid:    &SubscriptionId{Id: "sub2"},
			expect: true,
		},
		{
			name:   "Resource Group equals scope to Resource Group with different subscription id and resource group name",
			id:     &ResourceGroup{SubscriptionId: "sub1", Name: "rg1"},
			oid:    &ResourceGroup{SubscriptionId: "sub2", Name: "rg2"},
			expect: true,
		},
		{
			name:   "Management Group equals scope to Management Group with different name",
			id:     &ManagementGroup{Name: "mg1"},
			oid:    &ManagementGroup{Name: "mg2"},
			expect: true,
		},
		{
			name: "Root Scoped Resource under tenant equals scopes to different sub-type name",
			id: &ScopedResourceId{
				AttrParentScope: &TenantId{},
				AttrProvider:    "Microsoft.Foo",
				AttrTypes:       []string{"foos"},
				AttrNames:       []string{"foo1"},
			},
			oid: &ScopedResourceId{
				AttrParentScope: &TenantId{},
				AttrProvider:    "Microsoft.Foo",
				AttrTypes:       []string{"foos"},
				AttrNames:       []string{"foo2"},
			},
			expect: true,
		},
		{
			name: "Parent Scoped Resource under tenant not equals scopes to different sub-type type",
			id: &ScopedResourceId{
				AttrParentScope: &TenantId{},
				AttrProvider:    "Microsoft.Foo",
				AttrTypes:       []string{"foos"},
				AttrNames:       []string{"foo1"},
			},
			oid: &ScopedResourceId{
				AttrParentScope: &TenantId{},
				AttrProvider:    "Microsoft.Foo",
				AttrTypes:       []string{"bars"},
				AttrNames:       []string{"bar1"},
			},
			expect: false,
		},
		{
			name: "Parent Scoped Resource under tenant not equals scopes to its child Scoped Resource",
			id: &ScopedResourceId{
				AttrParentScope: &TenantId{},
				AttrProvider:    "Microsoft.Foo",
				AttrTypes:       []string{"foos"},
				AttrNames:       []string{"foo1"},
			},
			oid: &ScopedResourceId{
				AttrParentScope: &TenantId{},
				AttrProvider:    "Microsoft.Foo",
				AttrTypes:       []string{"foos", "bars"},
				AttrNames:       []string{"foo1", "bar1"},
			},
			expect: false,
		},
		{
			name: "Subscription scope level resource",
			id: &SubscriptionId{
				Id:        "sub1",
				AttrTypes: []string{"tagNames", "tagValues"},
				AttrNames: []string{"name1", "value1"},
			},
			oid: &SubscriptionId{
				Id:        "sub2",
				AttrTypes: []string{"tagNames", "tagValues"},
				AttrNames: []string{"name2", "value2"},
			},
			expect: true,
		},
		{
			name: "Resource group scope level resource",
			id: &ResourceGroup{
				SubscriptionId: "sub1",
				Name:           "rg1",
				AttrTypes:      []string{"deployments"},
				AttrNames:      []string{"deploy1"},
			},
			oid: &ResourceGroup{
				SubscriptionId: "sub2",
				Name:           "rg2",
				AttrTypes:      []string{"deployments"},
				AttrNames:      []string{"deploy2"},
			},
			expect: true,
		},
		{
			name: "Mgmt group scope level resource",
			id: &ManagementGroup{
				Name:      "group1",
				AttrTypes: []string{"foos"},
				AttrNames: []string{"foo1"},
			},
			oid: &ManagementGroup{
				Name:      "group2",
				AttrTypes: []string{"foos"},
				AttrNames: []string{"foo2"},
			},
			expect: true,
		},
		{
			name: `RP level resource`,
			id: &ScopedResourceId{
				AttrParentScope: &TenantId{},
				AttrProvider:    "Microsoft.Foo",
			},
			oid: &ScopedResourceId{
				AttrParentScope: &TenantId{},
				AttrProvider:    "Microsoft.Foo",
			},
			expect: true,
		},
		{
			name: `RP level resource under another RP level resource`,
			id: &ScopedResourceId{
				AttrParentScope: &ScopedResourceId{
					AttrParentScope: &TenantId{},
					AttrProvider:    "Microsoft.Foo",
				},
				AttrProvider: "Microsoft.Bar",
			},
			oid: &ScopedResourceId{
				AttrParentScope: &ScopedResourceId{
					AttrParentScope: &TenantId{},
					AttrProvider:    "Microsoft.Foo",
				},
				AttrProvider: "Microsoft.Bar",
			},
			expect: true,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.expect, tt.id.ScopeEqual(tt.oid))
		})
	}
}

func TestResourceId_ScopeString(t *testing.T) {
	cases := []struct {
		name   string
		id     ResourceId
		expect string
	}{
		{
			name:   "Tenant",
			id:     &TenantId{},
			expect: "/",
		},
		{
			name:   "Subscription",
			id:     &SubscriptionId{Id: "sub1"},
			expect: "/subscriptions",
		},
		{
			name:   "Resource Group",
			id:     &ResourceGroup{SubscriptionId: "sub1", Name: "rg1"},
			expect: "/subscriptions/resourceGroups",
		},
		{
			name:   "Management Group",
			id:     &ManagementGroup{Name: "mg1"},
			expect: "/Microsoft.Management/managementGroups",
		},
		{
			name: "Root Scoped Resource under tenant",
			id: &ScopedResourceId{
				AttrParentScope: &TenantId{},
				AttrProvider:    "Microsoft.Foo",
				AttrTypes:       []string{"foos"},
				AttrNames:       []string{"foo1"},
			},
			expect: "/Microsoft.Foo/foos",
		},
		{
			name: "Child Scoped Resource under resource group with",
			id: &ScopedResourceId{
				AttrParentScope: &ResourceGroup{
					SubscriptionId: "sub1",
					Name:           "rg1",
				},
				AttrProvider: "Microsoft.Foo",
				AttrTypes:    []string{"foos", "bars"},
				AttrNames:    []string{"foo1", "bar1"},
			},
			expect: "/subscriptions/resourceGroups/Microsoft.Foo/foos/bars",
		},
		{
			name: "Root Scoped Resource under resource group",
			id: &ScopedResourceId{
				AttrParentScope: &ResourceGroup{
					SubscriptionId: "sub1",
					Name:           "rg1",
				},
				AttrProvider: "Microsoft.Foo",
				AttrTypes:    []string{"foos"},
				AttrNames:    []string{"foo1"},
			},
			expect: "/subscriptions/resourceGroups/Microsoft.Foo/foos",
		},
		{
			name: "Subscription scope level resource",
			id: &SubscriptionId{
				Id:        "sub1",
				AttrTypes: []string{"tagNames", "tagValues"},
				AttrNames: []string{"name1", "value1"},
			},
			expect: "/subscriptions/tagNames/tagValues",
		},
		{
			name: "Resource group scope level resource",
			id: &ResourceGroup{
				SubscriptionId: "sub1",
				Name:           "rg1",
				AttrTypes:      []string{"deployments"},
				AttrNames:      []string{"deploy1"},
			},
			expect: "/subscriptions/resourceGroups/deployments",
		},
		{
			name: "Mgmt group scope level resource",
			id: &ManagementGroup{
				Name:      "group1",
				AttrTypes: []string{"foos"},
				AttrNames: []string{"foo1"},
			},
			expect: "/Microsoft.Management/managementGroups/foos",
		},
		{
			name: `RP level resource`,
			id: &ScopedResourceId{
				AttrParentScope: &TenantId{},
				AttrProvider:    "Microsoft.Foo",
			},
			expect: "/Microsoft.Foo",
		},
		{
			name: `RP level resource under another RP level resource`,
			id: &ScopedResourceId{
				AttrParentScope: &ScopedResourceId{
					AttrParentScope: &TenantId{},
					AttrProvider:    "Microsoft.Foo",
				},
				AttrProvider: "Microsoft.Bar",
			},
			expect: "/Microsoft.Foo/Microsoft.Bar",
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.expect, tt.id.ScopeString())
		})
	}
}

func TestResourceId_RouteScopeString(t *testing.T) {
	cases := []struct {
		name   string
		id     ResourceId
		expect string
	}{
		{
			name:   "Tenant",
			id:     &TenantId{},
			expect: "/",
		},
		{
			name:   "Subscription",
			id:     &SubscriptionId{Id: "sub1"},
			expect: "/subscriptions",
		},
		{
			name:   "Resource Group",
			id:     &ResourceGroup{SubscriptionId: "sub1", Name: "rg1"},
			expect: "/subscriptions/resourceGroups",
		},
		{
			name:   "Management Group",
			id:     &ManagementGroup{Name: "mg1"},
			expect: "/Microsoft.Management/managementGroups",
		},
		{
			name: "Root Scoped Resource under tenant",
			id: &ScopedResourceId{
				AttrParentScope: &TenantId{},
				AttrProvider:    "Microsoft.Foo",
				AttrTypes:       []string{"foos"},
				AttrNames:       []string{"foo1"},
			},
			expect: "/Microsoft.Foo/foos",
		},
		{
			name: "Child Scoped Resource under resource group with",
			id: &ScopedResourceId{
				AttrParentScope: &ResourceGroup{
					SubscriptionId: "sub1",
					Name:           "rg1",
				},
				AttrProvider: "Microsoft.Foo",
				AttrTypes:    []string{"foos", "bars"},
				AttrNames:    []string{"foo1", "bar1"},
			},
			expect: "/Microsoft.Foo/foos/bars",
		},
		{
			name: "Root Scoped Resource under resource group",
			id: &ScopedResourceId{
				AttrParentScope: &ResourceGroup{
					SubscriptionId: "sub1",
					Name:           "rg1",
				},
				AttrProvider: "Microsoft.Foo",
				AttrTypes:    []string{"foos"},
				AttrNames:    []string{"foo1"},
			},
			expect: "/Microsoft.Foo/foos",
		},
		{
			name: "Root Scoped Resource under resource group",
			id: &ScopedResourceId{
				AttrParentScope: &ResourceGroup{
					SubscriptionId: "sub1",
					Name:           "rg1",
				},
				AttrProvider: "Microsoft.Foo",
				AttrTypes:    []string{"foos"},
				AttrNames:    []string{"foo1"},
			},
			expect: "/Microsoft.Foo/foos",
		},
		{
			name: "Subscription scope level resource",
			id: &SubscriptionId{
				Id:        "sub1",
				AttrTypes: []string{"tagNames", "tagValues"},
				AttrNames: []string{"name1", "value1"},
			},
			expect: "/subscriptions/tagNames/tagValues",
		},
		{
			name: "Resource group scope level resource",
			id: &ResourceGroup{
				SubscriptionId: "sub1",
				Name:           "rg1",
				AttrTypes:      []string{"deployments"},
				AttrNames:      []string{"deploy1"},
			},
			expect: "/subscriptions/resourceGroups/deployments",
		},
		{
			name: "Mgmt group scope level resource",
			id: &ManagementGroup{
				Name:      "group1",
				AttrTypes: []string{"foos"},
				AttrNames: []string{"foo1"},
			},
			expect: "/Microsoft.Management/managementGroups/foos",
		},
		{
			name: `RP level resource`,
			id: &ScopedResourceId{
				AttrParentScope: &TenantId{},
				AttrProvider:    "Microsoft.Foo",
			},
			expect: "/Microsoft.Foo",
		},
		{
			name: `RP level resource under another RP level resource`,
			id: &ScopedResourceId{
				AttrParentScope: &ScopedResourceId{
					AttrParentScope: &TenantId{},
					AttrProvider:    "Microsoft.Foo",
				},
				AttrProvider: "Microsoft.Bar",
			},
			expect: "/Microsoft.Bar",
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.expect, tt.id.RouteScopeString())
		})
	}
}

func TestScopedResourceId_Normalize(t *testing.T) {
	cases := []struct {
		name     string
		id       ResourceId
		scopeStr string
		expect   string
		err      string
	}{
		{
			name: "Mismatch scope string",
			id: &ScopedResourceId{
				AttrParentScope: &TenantId{},
				AttrProvider:    "Microsoft.Foo",
				AttrTypes:       []string{"foos"},
				AttrNames:       []string{"foo1"},
			},
			scopeStr: "/Microsoft.Bar/foos",
			err:      `mismatch scope string ("/Microsoft.Bar/foos") for id "/providers/Microsoft.Foo/foos/foo1"`,
		},
		{
			name:     "Tenant root scope",
			id:       &TenantId{},
			scopeStr: "/",
			expect:   "/",
		},
		{
			name:     "Subscription root scope",
			id:       &SubscriptionId{Id: "sub1"},
			scopeStr: "/SUBSCRIPTIONS",
			expect:   "/SUBSCRIPTIONS/sub1",
		},
		{
			name: "Resource Group root scope",
			id: &ResourceGroup{
				SubscriptionId: "sub1",
				Name:           "grp1",
			},
			scopeStr: "/SUBSCRIPTIONS/RESOURCEGROUPS",
			expect:   "/SUBSCRIPTIONS/sub1/RESOURCEGROUPS/grp1",
		},
		{
			name: "Management Group root scope",
			id: &ManagementGroup{
				Name: "grp1",
			},
			scopeStr: "/MICROSOFT.MANAGEMENT/MANAGEMENTGROUPS",
			expect:   "/providers/MICROSOFT.MANAGEMENT/MANAGEMENTGROUPS/grp1",
		},
		{
			name: "Root Scoped Resource under tenant",
			id: &ScopedResourceId{
				AttrParentScope: &TenantId{},
				AttrProvider:    "MICROSOFT.Foo",
				AttrTypes:       []string{"FOOS"},
				AttrNames:       []string{"foo1"},
			},
			scopeStr: "/microsoft.foo/foos",
			expect:   "/providers/microsoft.foo/foos/foo1",
		},
		{
			name: "Root Scoped Resource under resource group (scopestr lower casing)",
			id: &ScopedResourceId{
				AttrParentScope: &ResourceGroup{
					SubscriptionId: "sub1",
					Name:           "rg1",
				},
				AttrProvider: "MICROSOFT.Foo",
				AttrTypes:    []string{"FOOS"},
				AttrNames:    []string{"foo1"},
			},
			scopeStr: "/subscriptions/resourcegroups/microsoft.foo/foos",
			expect:   "/subscriptions/sub1/resourcegroups/rg1/providers/microsoft.foo/foos/foo1",
		},
		{
			name: "Root Scoped Resource under resource group (scopestr upper casing)",
			id: &ScopedResourceId{
				AttrParentScope: &ResourceGroup{
					SubscriptionId: "sub1",
					Name:           "rg1",
				},
				AttrProvider: "MICROSOFT.Foo",
				AttrTypes:    []string{"FOOS"},
				AttrNames:    []string{"foo1"},
			},
			scopeStr: "/SUBSCRIPTIONS/RESOURCEGROUPS/MICROSOFT.FOO/FOOS",
			expect:   "/SUBSCRIPTIONS/sub1/RESOURCEGROUPS/rg1/providers/MICROSOFT.FOO/FOOS/foo1",
		},
		{
			name: "Subscription scope level resource",
			id: &SubscriptionId{
				Id:        "sub1",
				AttrTypes: []string{"tagNames", "tagValues"},
				AttrNames: []string{"name1", "value1"},
			},
			scopeStr: "/SUBSCRIPTIONS/TAGNAMES/TAGVALUES",
			expect:   "/SUBSCRIPTIONS/sub1/TAGNAMES/name1/TAGVALUES/value1",
		},
		{
			name: "Resource group scope level resource",
			id: &ResourceGroup{
				SubscriptionId: "sub1",
				Name:           "rg1",
				AttrTypes:      []string{"deployments"},
				AttrNames:      []string{"deploy1"},
			},
			scopeStr: "/SUBSCRIPTIONS/RESOURCEGROUPS/DEPLOYMENTS",
			expect:   "/SUBSCRIPTIONS/sub1/RESOURCEGROUPS/rg1/DEPLOYMENTS/deploy1",
		},
		{
			name: "Mgmt group scope level resource",
			id: &ManagementGroup{
				Name:      "group1",
				AttrTypes: []string{"foos"},
				AttrNames: []string{"foo1"},
			},
			scopeStr: "/MICROSOFT.MANAGEMENT/MANAGEMENTGROUPS/FOOS",
			expect:   "/providers/MICROSOFT.MANAGEMENT/MANAGEMENTGROUPS/group1/FOOS/foo1",
		},
		{
			name: `RP level resource`,
			id: &ScopedResourceId{
				AttrParentScope: &TenantId{},
				AttrProvider:    "Microsoft.Foo",
			},
			scopeStr: "/MICROSOFT.FOO",
			expect:   "/providers/MICROSOFT.FOO",
		},
		{
			name: `RP level resource under another RP level resource`,
			id: &ScopedResourceId{
				AttrParentScope: &ScopedResourceId{
					AttrParentScope: &TenantId{},
					AttrProvider:    "Microsoft.Foo",
				},
				AttrProvider: "Microsoft.Bar",
			},
			scopeStr: "/MICROSOFT.FOO/MICROSOFT.BAR",
			expect:   "/providers/MICROSOFT.FOO/providers/MICROSOFT.BAR",
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			id := tt.id
			err := id.Normalize(tt.scopeStr)
			if tt.err != "" {
				require.EqualError(t, err, tt.err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.expect, id.String())
		})
	}
}

func TestScopedResourceId_Clone(t *testing.T) {
	cases := []struct {
		name   string
		id     ResourceId
		mutate func(ResourceId)
		expect ResourceId
	}{
		{
			name:   "Tenant root scope",
			id:     &TenantId{},
			mutate: nil,
			expect: &TenantId{},
		},
		{
			name: "Subscription root scope",
			id:   &SubscriptionId{Id: "sub1"},
			mutate: func(ri ResourceId) {
				id := ri.(*SubscriptionId)
				id.Id = "sub2"
			},
			expect: &SubscriptionId{Id: "sub1"},
		},
		{
			name: "Resource Group root scope",
			id: &ResourceGroup{
				SubscriptionId: "sub1",
				Name:           "grp1",
			},
			mutate: func(ri ResourceId) {
				id := ri.(*ResourceGroup)
				id.SubscriptionId = "sub2"
			},
			expect: &ResourceGroup{
				SubscriptionId: "sub1",
				Name:           "grp1",
			},
		},
		{
			name: "Management Group root scope",
			id: &ManagementGroup{
				Name: "grp1",
			},
			mutate: func(ri ResourceId) {
				id := ri.(*ManagementGroup)
				id.Name = "grp2"
			},
			expect: &ManagementGroup{
				Name: "grp1",
			},
		},
		{
			name: "Root Scoped Resource under tenant",
			id: &ScopedResourceId{
				AttrParentScope: &TenantId{},
				AttrProvider:    "Microsoft.Foo",
				AttrTypes:       []string{"foos"},
				AttrNames:       []string{"foo1"},
			},
			mutate: func(ri ResourceId) {
				id := ri.(*ScopedResourceId)
				id.AttrNames[0] = "foo2"
			},
			expect: &ScopedResourceId{
				AttrParentScope: &TenantId{},
				AttrProvider:    "Microsoft.Foo",
				AttrTypes:       []string{"foos"},
				AttrNames:       []string{"foo1"},
			},
		},
		{
			name: "Root Scoped Resource under resource group",
			id: &ScopedResourceId{
				AttrParentScope: &ResourceGroup{
					SubscriptionId: "sub1",
					Name:           "rg1",
				},
				AttrProvider: "Microsoft.Foo",
				AttrTypes:    []string{"foos"},
				AttrNames:    []string{"foo1"},
			},
			mutate: func(ri ResourceId) {
				id := ri.(*ScopedResourceId)
				id.AttrNames[0] = "foo2"
			},
			expect: &ScopedResourceId{
				AttrParentScope: &ResourceGroup{
					SubscriptionId: "sub1",
					Name:           "rg1",
				},
				AttrProvider: "Microsoft.Foo",
				AttrTypes:    []string{"foos"},
				AttrNames:    []string{"foo1"},
			},
		},
		{
			name: "Subscription scope level resource",
			id: &SubscriptionId{
				Id:        "sub1",
				AttrTypes: []string{"tagNames", "tagValues"},
				AttrNames: []string{"name1", "value1"},
			},
			expect: &SubscriptionId{
				Id:        "sub1",
				AttrTypes: []string{"tagNames", "tagValues"},
				AttrNames: []string{"name1", "value1"},
			},
		},
		{
			name: "Resource group scope level resource",
			id: &ResourceGroup{
				SubscriptionId: "sub1",
				Name:           "rg1",
				AttrTypes:      []string{"deployments"},
				AttrNames:      []string{"deploy1"},
			},
			expect: &ResourceGroup{
				SubscriptionId: "sub1",
				Name:           "rg1",
				AttrTypes:      []string{"deployments"},
				AttrNames:      []string{"deploy1"},
			},
		},
		{
			name: "Mgmt group scope level resource",
			id: &ManagementGroup{
				Name:      "group1",
				AttrTypes: []string{"foos"},
				AttrNames: []string{"foo1"},
			},
			expect: &ManagementGroup{
				Name:      "group1",
				AttrTypes: []string{"foos"},
				AttrNames: []string{"foo1"},
			},
		},
		{
			name: `RP level resource`,
			id: &ScopedResourceId{
				AttrParentScope: &TenantId{},
				AttrProvider:    "Microsoft.Foo",
			},
			expect: &ScopedResourceId{
				AttrParentScope: &TenantId{},
				AttrProvider:    "Microsoft.Foo",
			},
		},
		{
			name: `RP level resource under another RP level resource`,
			id: &ScopedResourceId{
				AttrParentScope: &ScopedResourceId{
					AttrParentScope: &TenantId{},
					AttrProvider:    "Microsoft.Foo",
				},
				AttrProvider: "Microsoft.Bar",
			},
			expect: &ScopedResourceId{
				AttrParentScope: &ScopedResourceId{
					AttrParentScope: &TenantId{},
					AttrProvider:    "Microsoft.Foo",
				},
				AttrProvider: "Microsoft.Bar",
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			id := tt.id
			cid := id.Clone()
			if tt.mutate != nil {
				tt.mutate(id)
			}
			require.Equal(t, tt.expect, cid)
		})
	}
}

func TestScopedResourceId_NormalizeRouteScope(t *testing.T) {
	cases := []struct {
		name     string
		id       ScopedResourceId
		scopeStr string
		expect   string
		err      string
	}{
		{
			name: "Mismatch scope string",
			id: ScopedResourceId{
				AttrParentScope: &TenantId{},
				AttrProvider:    "Microsoft.Foo",
				AttrTypes:       []string{"foos"},
				AttrNames:       []string{"foo1"},
			},
			scopeStr: "/Microsoft.Bar/foos",
			err:      `mismatch route scope string ("/Microsoft.Bar/foos") for id "/providers/Microsoft.Foo/foos/foo1"`,
		},
		{
			name: "Root Scoped Resource under tenant",
			id: ScopedResourceId{
				AttrParentScope: &TenantId{},
				AttrProvider:    "MICROSOFT.Foo",
				AttrTypes:       []string{"FOOS"},
				AttrNames:       []string{"foo1"},
			},
			scopeStr: "/microsoft.foo/foos",
			expect:   "/providers/microsoft.foo/foos/foo1",
		},
		{
			name: "Root Scoped Resource under resource group",
			id: ScopedResourceId{
				AttrParentScope: &ResourceGroup{
					SubscriptionId: "sub1",
					Name:           "rg1",
				},
				AttrProvider: "MICROSOFT.Foo",
				AttrTypes:    []string{"FOOS"},
				AttrNames:    []string{"foo1"},
			},
			scopeStr: "/microsoft.foo/foos",
			expect:   "/subscriptions/sub1/resourceGroups/rg1/providers/microsoft.foo/foos/foo1",
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			id := tt.id
			err := id.NormalizeRouteScope(tt.scopeStr)
			if tt.err != "" {
				require.EqualError(t, err, tt.err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.expect, id.String())
		})
	}
}

func TestParseResourceId(t *testing.T) {
	cases := []struct {
		name   string
		input  string
		expect ResourceId
		err    string
	}{
		{
			name:   "Tenant",
			input:  "/",
			expect: &TenantId{},
		},
		{
			name:   "Subscription",
			input:  "/subscriptions/sub1",
			expect: &SubscriptionId{Id: "sub1"},
		},
		{
			name:   "Resource Group",
			input:  "/subscriptions/sub1/resourceGroups/rg1",
			expect: &ResourceGroup{SubscriptionId: "sub1", Name: "rg1"},
		},
		{
			name:   "Case-insensitive for resourceGroups",
			input:  "/SUBSCRIPTIONS/sub1/RESOURCEGROUPS/rg1",
			expect: &ResourceGroup{SubscriptionId: "sub1", Name: "rg1"},
		},
		{
			name:   "Management Group",
			input:  "/providers/Microsoft.Management/managementGroups/mg1",
			expect: &ManagementGroup{Name: "mg1"},
		},
		{
			name:   "Case-insensitive for managementGroup",
			input:  "/PROVIDERS/MICROSOFT.MANAGEMENT/MANAGEMENTGROUPS/mg1",
			expect: &ManagementGroup{Name: "mg1"},
		},
		{
			name:  "Scoped Resource under tenant",
			input: "/providers/Microsoft.Foo/foos/foo1/bars/bar1",
			expect: &ScopedResourceId{
				AttrParentScope: &TenantId{},
				AttrProvider:    "Microsoft.Foo",
				AttrTypes:       []string{"foos", "bars"},
				AttrNames:       []string{"foo1", "bar1"},
			},
		},
		{
			name:  "Case-insensitiev Scoped Resource under tenant",
			input: "/PROVIDERS/MICROSOFT.FOO/FOOS/foo1/BARS/bar1",
			expect: &ScopedResourceId{
				AttrParentScope: &TenantId{},
				AttrProvider:    "MICROSOFT.FOO",
				AttrTypes:       []string{"FOOS", "BARS"},
				AttrNames:       []string{"foo1", "bar1"},
			},
		},
		{
			name:  "Scoped Resource under subscription",
			input: "/subscriptions/sub1/providers/Microsoft.Foo/foos/foo1/bars/bar1",
			expect: &ScopedResourceId{
				AttrParentScope: &SubscriptionId{Id: "sub1"},
				AttrProvider:    "Microsoft.Foo",
				AttrTypes:       []string{"foos", "bars"},
				AttrNames:       []string{"foo1", "bar1"},
			},
		},
		{
			name:  "Scoped Resource under resource group",
			input: "/subscriptions/sub1/resourceGroups/rg1/providers/Microsoft.Foo/foos/foo1/bars/bar1",
			expect: &ScopedResourceId{
				AttrParentScope: &ResourceGroup{SubscriptionId: "sub1", Name: "rg1"},
				AttrProvider:    "Microsoft.Foo",
				AttrTypes:       []string{"foos", "bars"},
				AttrNames:       []string{"foo1", "bar1"},
			},
		},
		{
			name:  "Scoped Resource under management group",
			input: "/providers/Microsoft.Management/managementGroups/mg1/providers/Microsoft.Foo/foos/foo1/bars/bar1",
			expect: &ScopedResourceId{
				AttrParentScope: &ManagementGroup{Name: "mg1"},
				AttrProvider:    "Microsoft.Foo",
				AttrTypes:       []string{"foos", "bars"},
				AttrNames:       []string{"foo1", "bar1"},
			},
		},
		{
			name:  "Scoped Resource under another scoped resource which under tenant",
			input: "/providers/Microsoft.Foo/foos/foo1/bars/bar1/providers/Microsoft.Baz/bazs/baz1",
			expect: &ScopedResourceId{
				AttrParentScope: &ScopedResourceId{
					AttrParentScope: &TenantId{},
					AttrProvider:    "Microsoft.Foo",
					AttrTypes:       []string{"foos", "bars"},
					AttrNames:       []string{"foo1", "bar1"},
				},
				AttrProvider: "Microsoft.Baz",
				AttrTypes:    []string{"bazs"},
				AttrNames:    []string{"baz1"},
			},
		},
		{
			name:  "Scoped Resource under another scoped resource which under subscription",
			input: "/subscriptions/sub1/providers/Microsoft.Foo/foos/foo1/bars/bar1/providers/Microsoft.Baz/bazs/baz1",
			expect: &ScopedResourceId{
				AttrParentScope: &ScopedResourceId{
					AttrParentScope: &SubscriptionId{Id: "sub1"},
					AttrProvider:    "Microsoft.Foo",
					AttrTypes:       []string{"foos", "bars"},
					AttrNames:       []string{"foo1", "bar1"},
				},
				AttrProvider: "Microsoft.Baz",
				AttrTypes:    []string{"bazs"},
				AttrNames:    []string{"baz1"},
			},
		},
		{
			name:  "Scoped Resource under another scoped resource which under resource group",
			input: "/subscriptions/sub1/resourceGroups/rg1/providers/Microsoft.Foo/foos/foo1/bars/bar1/providers/Microsoft.Baz/bazs/baz1",
			expect: &ScopedResourceId{
				AttrParentScope: &ScopedResourceId{
					AttrParentScope: &ResourceGroup{SubscriptionId: "sub1", Name: "rg1"},
					AttrProvider:    "Microsoft.Foo",
					AttrTypes:       []string{"foos", "bars"},
					AttrNames:       []string{"foo1", "bar1"},
				},
				AttrProvider: "Microsoft.Baz",
				AttrTypes:    []string{"bazs"},
				AttrNames:    []string{"baz1"},
			},
		},
		{
			name:  "Scoped Resource under another scoped resource which under management group",
			input: "/providers/Microsoft.Management/managementGroups/mg1/providers/Microsoft.Foo/foos/foo1/bars/bar1/providers/Microsoft.Baz/bazs/baz1",
			expect: &ScopedResourceId{
				AttrParentScope: &ScopedResourceId{
					AttrParentScope: &ManagementGroup{Name: "mg1"},
					AttrProvider:    "Microsoft.Foo",
					AttrTypes:       []string{"foos", "bars"},
					AttrNames:       []string{"foo1", "bar1"},
				},
				AttrProvider: "Microsoft.Baz",
				AttrTypes:    []string{"bazs"},
				AttrNames:    []string{"baz1"},
			},
		},
		{
			name:  "Subscription scope level resource",
			input: "/subscriptions/sub1/tagNames/name1/tagValues/value1",
			expect: &SubscriptionId{
				Id:        "sub1",
				AttrTypes: []string{"tagNames", "tagValues"},
				AttrNames: []string{"name1", "value1"},
			},
		},
		{
			name:  "Resource group scope level resource",
			input: "/subscriptions/sub1/resourceGroups/rg1/deployments/deploy1",
			expect: &ResourceGroup{
				SubscriptionId: "sub1",
				Name:           "rg1",
				AttrTypes:      []string{"deployments"},
				AttrNames:      []string{"deploy1"},
			},
		},
		{
			name:  "Mgmt group scope level resource",
			input: "/providers/Microsoft.Management/managementGroups/group1/foos/foo1",
			expect: &ManagementGroup{
				Name:      "group1",
				AttrTypes: []string{"foos"},
				AttrNames: []string{"foo1"},
			},
		},
		{
			name:  `RP level resource`,
			input: "/providers/Microsoft.Foo",
			expect: &ScopedResourceId{
				AttrParentScope: &TenantId{},
				AttrProvider:    "Microsoft.Foo",
				AttrTypes:       []string{},
				AttrNames:       []string{},
			},
		},
		{
			name:  `RP level resource under another RP level resource`,
			input: "/providers/Microsoft.Foo/providers/Microsoft.Bar",
			expect: &ScopedResourceId{
				AttrParentScope: &ScopedResourceId{
					AttrParentScope: &TenantId{},
					AttrProvider:    "Microsoft.Foo",
					AttrTypes:       []string{},
					AttrNames:       []string{},
				},
				AttrProvider: "Microsoft.Bar",
				AttrTypes:    []string{},
				AttrNames:    []string{},
			},
		},
		{
			name:  "empty string",
			input: "",
			err:   `id should start with "/"`,
		},
		{
			name:  "id not starts with /",
			input: "foo",
			err:   `id should start with "/"`,
		},
		{
			name:  `id ends with "/"`,
			input: "/providers/",
			err:   `empty segment found behind 2th "/"`,
		},
		{
			name:  `id has empty segment in the middle "/"`,
			input: "/providers/Microsoft.Foo/foos//foo1",
			err:   `empty segment found behind 4th "/"`,
		},
		{
			name:  "invalid scope behind tenant scope",
			input: "/foo",
			err:   `extending for root level RP: missing resource type name after type foo`,
		},
		{
			name:  "invalid scope behind subscription scope",
			input: "/subscriptions/sub1/foo",
			err:   `extending for root level RP: missing resource type name after type foo`,
		},
		{
			name:  "invalid scope behind resource group scope",
			input: "/subscriptions/sub1/resourceGroups/rg1/foo",
			err:   `extending for root level RP: missing resource type name after type foo`,
		},
		{
			name:  "invalid scope behind management group scope",
			input: "/providers/Microsoft.Management/managementGroups/mg1/foo",
			err:   `extending for root level RP: missing resource type name after type foo`,
		},
		{
			name:  `missing provider namespace segment`,
			input: "/providers",
			err:   `missing provider namespace segment`,
		},
		{
			name:  `missing sub-type name`,
			input: "/providers/Microsoft.Foo/foos",
			err:   `extending for RP Microsoft.Foo: missing resource type name after type foos`,
		},
		{
			name:  `missing sub-type name in child`,
			input: "/providers/Microsoft.Foo/foos/foo1/bars",
			err:   `extending for RP Microsoft.Foo: missing resource type name after type bars`,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			id, err := ParseResourceId(tt.input)
			if tt.err != "" {
				require.EqualError(t, err, tt.err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.expect, id)
		})
	}
}
