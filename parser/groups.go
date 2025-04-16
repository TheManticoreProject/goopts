package parser

import (
	"fmt"

	"github.com/TheManticoreProject/goopts/argumentgroup"
)

// NewArgumentGroup creates a new argument group with the specified name and adds it to the list of child groups.
//
// Parameters:
// - name: The name of the new argument group.
//
// Returns:
// A pointer to the newly created ArgumentGroup struct, which can be used to add arguments to the group.
//
// The function initializes a new ArgumentGroup with the provided name, appends it to the SubGroups slice,
// and returns a pointer to the newly added ArgumentGroup. This allows for further configuration and addition
// of arguments to the group.
func (ap *ArgumentsParser) NewArgumentGroup(name string) (*argumentgroup.ArgumentGroup, error) {
	// Initiate the map if it was not initialized yet
	if ap.Groups == nil {
		ap.Groups = make(map[string]*argumentgroup.ArgumentGroup)
	}

	if len(name) == 0 {
		return nil, fmt.Errorf("name of group cannot be empty, this is reserved for the default group")
	}

	group := argumentgroup.ArgumentGroup{
		Name: name,
		Type: argumentgroup.ARGUMENT_GROUP_TYPE_NORMAL,
	}

	// Add it to the Groups
	if _, exists := ap.Groups[group.Name]; !exists {
		ap.Groups[group.Name] = &group
	} else {
		return nil, fmt.Errorf("an argument group with name \"%s\" already exists", group.Name)
	}

	return &group, nil
}

// NewNotRequiredMutuallyExclusiveArgumentGroup creates a new argument group with the specified name
// and adds it to the list of child groups. This group allows at most one of the arguments
// within it to be set, but it is not mandatory to provide any.
//
// Parameters:
// - name: The name of the new argument group.
//
// Returns:
// A pointer to the newly created ArgumentGroup struct, which can be used to add arguments to the group,
// or an error if a group with the same name already exists.
//
// The function initializes a new ArgumentGroup with the provided name and the type
// ARGUMENT_GROUP_TYPE_NOT_REQUIRED_MUTUALLY_EXCLUSIVE. It appends the group to the Groups map
// and returns a pointer to the newly added ArgumentGroup. This allows for further configuration
// and addition of arguments to the group.
func (ap *ArgumentsParser) NewNotRequiredMutuallyExclusiveArgumentGroup(name string) (*argumentgroup.ArgumentGroup, error) {
	// Initiate the map if it was not initialized yet
	if ap.Groups == nil {
		ap.Groups = make(map[string]*argumentgroup.ArgumentGroup)
	}

	group := argumentgroup.ArgumentGroup{
		Name: name,
		Type: argumentgroup.ARGUMENT_GROUP_TYPE_NOT_REQUIRED_MUTUALLY_EXCLUSIVE,
	}

	// Add it to the Groups
	if _, exists := ap.Groups[group.Name]; !exists {
		ap.Groups[group.Name] = &group
	} else {
		return nil, fmt.Errorf("an argument group with name \"%s\" already exists", group.Name)
	}

	return &group, nil
}

// NewRequiredMutuallyExclusiveArgumentGroup creates a new argument group with the specified name
// and adds it to the list of child groups. This group enforces that exactly one of the arguments
// within it must be set; if none or more than one are provided, an error will be thrown.
//
// Parameters:
// - name: The name of the new argument group.
//
// Returns:
// A pointer to the newly created ArgumentGroup struct, which can be used to add arguments to the group,
// or an error if a group with the same name already exists.
//
// The function initializes a new ArgumentGroup with the provided name and the type
// ARGUMENT_GROUP_TYPE_REQUIRED_MUTUALLY_EXCLUSIVE. It adds the group to the Groups map
// and returns a pointer to the newly added ArgumentGroup. This allows for further configuration
// and addition of arguments to the group, ensuring that only one of the arguments is set during parsing.
func (ap *ArgumentsParser) NewRequiredMutuallyExclusiveArgumentGroup(name string) (*argumentgroup.ArgumentGroup, error) {
	// Initiate the map if it was not initialized yet
	if ap.Groups == nil {
		ap.Groups = make(map[string]*argumentgroup.ArgumentGroup)
	}

	group := argumentgroup.ArgumentGroup{
		Name: name,
		Type: argumentgroup.ARGUMENT_GROUP_TYPE_REQUIRED_MUTUALLY_EXCLUSIVE,
	}

	// Add it to the Groups
	if _, exists := ap.Groups[group.Name]; !exists {
		ap.Groups[group.Name] = &group
	} else {
		return nil, fmt.Errorf("an argument group with name \"%s\" already exists", group.Name)
	}

	return &group, nil
}

// NewRequiredMutuallyExclusiveArgumentGroup creates a new argument group with the specified name
// and adds it to the list of child groups. This group enforces that exactly one of the arguments
// within it must be set; if none or more than one are provided, an error will be thrown.
//
// Parameters:
// - name: The name of the new argument group.
//
// Returns:
// A pointer to the newly created ArgumentGroup struct, which can be used to add arguments to the group,
// or an error if a group with the same name already exists.
//
// The function initializes a new ArgumentGroup with the provided name and the type
// ARGUMENT_GROUP_TYPE_REQUIRED_MUTUALLY_EXCLUSIVE. It adds the group to the Groups map
// and returns a pointer to the newly added ArgumentGroup. This allows for further configuration
// and addition of arguments to the group, ensuring that only one of the arguments is set during parsing.
func (ap *ArgumentsParser) NewDependentArgumentGroup(name string) (*argumentgroup.ArgumentGroup, error) {
	// Initiate the map if it was not initialized yet
	if ap.Groups == nil {
		ap.Groups = make(map[string]*argumentgroup.ArgumentGroup)
	}

	group := argumentgroup.ArgumentGroup{
		Name: name,
		Type: argumentgroup.ARGUMENT_GROUP_TYPE_DEPENDENT,
	}

	// Add it to the Groups
	if _, exists := ap.Groups[group.Name]; !exists {
		ap.Groups[group.Name] = &group
	} else {
		return nil, fmt.Errorf("an argument group with name \"%s\" already exists", group.Name)
	}

	return &group, nil
}
