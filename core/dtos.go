package core

type DependencyTree struct {
	Items []DependencyTreeItem
}

type DependencyTreeItem struct {
	Item   Dependency
	Parent *Dependency
}

type Dependency struct {
	Name            string
	RequiredVersion string
	UsedVersion     string
}
