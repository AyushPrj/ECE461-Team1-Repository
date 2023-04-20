/*
 * ECE 461 - Spring 2023 - Project 2
 *
 * API for ECE 461/Spring 2023/Project 2: A Trustworthy Module Registry
 *
 * API version: 3.0.2
 * Contact: davisjam@purdue.edu
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */
package swagger

// Authentication info for a user
type UserAuthenticationInfo struct {
	// Password for a user. Per the spec, this should be a \"strong\" password.
	Password string `json:"password"`
}