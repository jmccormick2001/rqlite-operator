package controller

import (
	"github.com/jmccormick2001/rqlite-operator/pkg/controller/rqcluster"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, rqcluster.Add)
}
