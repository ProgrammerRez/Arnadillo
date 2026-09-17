package main

// This type is to provide easy API Output for generic slices
type GenericSliceOutput struct{
	Items 	[]string	`json:"items"`
	Count	int			`json:"count"`
}