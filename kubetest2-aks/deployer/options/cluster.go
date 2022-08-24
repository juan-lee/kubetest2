package options

type ClusterOptions struct {
	Template      string `flag:"~template t" desc:"aks arm template."`
	ResourceGroup string `flag:"~resource-group g" desc:"azure resource group."`
}
