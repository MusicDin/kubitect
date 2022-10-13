package modelconfig

import v "cli/validation"

type Host struct {
	Name                 *string             `yaml:"name" opt:",id"`
	Default              *bool               `yaml:"default"`
	Connection           *Connection         `yaml:"connection"`
	MainResourcePoolPath *string             `yaml:"mainResourcePoolPath"`
	DataResourcePools    *[]DataResourcePool `yaml:"dataResourcePools"`
}

func (h Host) Validate() error {
	return v.Struct(&h,
		v.Field(&h.Name, v.Required(), v.AlphaNumericHypUS()),
		// v.Field(&h.Default), // TODO: only one can be true
		v.Field(&h.Connection),
		v.Field(&h.MainResourcePoolPath, v.OmitEmpty()), // TODO: validate dir path which does not have to exist
		v.Field(&h.DataResourcePools, v.OmitEmpty(), v.UniqueField("Name")),
	)
}
