// Package apiserver does all of the work necessary to create a iam APIServer.
package apiserver

import "github.com/rose839/IAM/pkg/app"

const commandDesc = `The IAM API server validates and configures data
for the api objects which include users, policies, secrets, and
others. The API Server services REST operations to do the api objects management.`

func NewApp(basename string) *app.App {

}
