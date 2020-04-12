/*
@Time : 2020/3/26 下午5:26
@Author : songxiuxuan
@File : version_test.go
@Software: GoLand
*/
package task

import (
	"testing"
)

func TestStringConverFloat64(t *testing.T) {
	var testCase = map[string]struct {
		Version float64
	}{
		"10.222": {
			Version: 10.222,
		},
		"92222": {
			Version: 92222,
		},
		"v1.0.222": {
			Version: 1.0222,
		},
		"v1.20.342.2.2.22": {
			Version: 1.203422222,
		},
		"": {
			Version: 0.0,
		},
		"*": {
			Version: 0.0,
		},
	}
	for str, caseStr := range testCase {
		ver, err := StringConverFloat64(str, ".", 64)
		if err != nil {
			t.Fatal(err)
		}

		if caseStr.Version != ver {
			t.Fatal("version not match", ver, "not eq", caseStr.Version)
		}
	}
}

func TestStringConverOperator(t *testing.T) {
	var testCase = map[string]struct {
		Operator  string
		Version1  string
		Version2  string
		Ver1Float float64
		Ver2Float float64
	}{
		">v1.33.92": {
			Operator:  ">",
			Version1:  "",
			Version2:  "v1.33.92",
			Ver1Float: 0.0,
			Ver2Float: 1.3392,
		},
		">=v132.00.22.222": {
			Operator:  ">=",
			Version1:  "",
			Version2:  "v132.00.22.222",
			Ver1Float: 0.0,
			Ver2Float: 132.0022222,
		},
		"<v1.0898.2": {
			Operator:  "<",
			Version1:  "",
			Version2:  "v1.0898.2",
			Ver1Float: 0.0,
			Ver2Float: 1.08982,
		},
		"<=v1.30.292": {
			Operator:  "<=",
			Version1:  "",
			Version2:  "v1.30.292",
			Ver1Float: 0.0,
			Ver2Float: 1.30292,
		},
		"=v9.0.2342": {
			Operator:  "=",
			Version1:  "",
			Version2:  "v9.0.2342",
			Ver1Float: 0.0,
			Ver2Float: 9.02342,
		},
		"v83.8.2322": {
			Operator:  "",
			Version1:  "",
			Version2:  "v83.8.2322",
			Ver1Float: 0.0,
			Ver2Float: 83.82322,
		},
		"~v83.8.8322": {
			Operator:  "~",
			Version1:  "",
			Version2:  "v83.8.8322",
			Ver1Float: 0.0,
			Ver2Float: 83.88322,
		},
		"v83.8.8322,v89.8.00": {
			Operator:  ",",
			Version1:  "v83.8.8322",
			Version2:  "v89.8.00",
			Ver1Float: 83.88322,
			Ver2Float: 89.8,
		},
		"*": {
			Operator:  "*",
			Version1:  "",
			Version2:  "",
			Ver1Float: 0.0,
			Ver2Float: 0.0,
		},
	}

	for str, caseStr := range testCase {
		va, op, vb := StringConverOperator(str)
		if caseStr.Operator != op {
			t.Fatal("operator not match", op, "not eq", caseStr.Operator)
		}
		if caseStr.Version1 != va {
			t.Fatal("version not match==>1", va, "not eq", caseStr.Version1)
		}
		if caseStr.Version2 != vb {
			t.Fatal("version not match==>2", vb, "not eq", caseStr.Version2)
		}

		ca, aerr := StringConverFloat64(va, ".", 64)
		if aerr != nil {
			t.Fatal(aerr)
		}
		cb, berr := StringConverFloat64(vb, ".", 64)
		if berr != nil {
			t.Fatal(berr)
		}

		if caseStr.Ver1Float != ca {
			t.Fatal("version not match==>1", ca, "not eq", caseStr.Ver1Float)
		}
		if caseStr.Ver2Float != cb {
			t.Fatal("version not match===>2", cb, "not eq", caseStr.Ver2Float)
		}
	}
}

func TestVersionComp(t *testing.T) {
	var testCase = map[string]struct {
		Version  []string
		VerFloat float64
	}{
		">v1.35": {
			Version: []string{"v1.38.9873", "v1.39.9873.222"},
		},
		">=v135.3": {
			Version: []string{"v135.30.22.222", "v135.3", "v135.3000"},
		},
		"<v1.22": {
			Version: []string{"v1.21.22", "v1.219.22", "v1.13.222", "v1.1.22.222"},
		},
		"<=v3.30": {
			Version: []string{"v3.2.292", "v3.30", "v3.29", "v3.300.00"},
		},
		"=v9.22": {
			Version: []string{"v9.22"},
		},
		"v8.8": {
			Version: []string{"v8.8"},
		},
		"~v7.8": {
			Version: []string{"v7.8.2987", "v7.8.9987", "v7.8.2987", "v7.8.3337"},
		},
	}

	for str, caseStr := range testCase {
		v1, op, v2 := StringConverOperator(str)
		vb, berr := StringConverFloat64(v1, ".", 64)
		if berr != nil {
			t.Fatal(berr)
		}
		vc, cerr := StringConverFloat64(v2, ".", 64)
		if cerr != nil {
			t.Fatal(cerr)
		}
		for _, ver := range caseStr.Version {
			va, err := StringConverFloat64(ver, ".", 64)
			if err != nil {
				t.Fatal(err)
			}
			if !VersionCompare(va, vb, op, vc) {
				t.Fatal("match faield", va, op, vb)
			}
		}
	}
}
