package main

import (
	"github.com/brocaar/lorawan"
)

var GenAppKey = lorawan.AES128Key{0x2B, 0x7E, 0x15, 0x16, 0x28, 0xAE, 0xD2, 0xA6, 0xAB, 0xF7, 0x15, 0x88, 0x09, 0xCF, 0x4F, 0x3C}

var ApplicationId = "6e5d05ec-24a3-4644-812e-f3f38eb4359e"

//var DevEuis = []string{"0080e1150530fef1", "0080e115000acbcf", "0080e115000acc1b", "0080e1150530f6c7", "0080e1150530fcba"}

//var DevEuis = []string{"0080e1150530fcba"}

var DevEuis = []string{"0080e1150530fef1", "0080e115000acbcf", "0080e115000acc1b"}
