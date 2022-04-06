package controller

import (
	"webank/AIDE/notebook-controller/pkg/controller/notebook"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, notebook.Add)
}
